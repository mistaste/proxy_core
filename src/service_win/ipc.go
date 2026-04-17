//go:build windows

// Package service_win implements the privileged Windows service that owns
// the wintun adapter and installs routing table entries. The main Flutter
// application, running without elevation, talks to this service via a
// named pipe so only the service itself ever needs Admin/SYSTEM rights.
package service_win

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Microsoft/go-winio"

	"segment/libtun"
)

// PipeName is the named pipe address the service listens on. The Flutter
// plugin connects to the same path.
const PipeName = `\\.\pipe\guardex-svc`

// pipeSDDL grants the local Administrators group full access and any
// Interactive user (i.e. a logged-in desktop session on the same machine)
// read+write. This mirrors what WireGuard-like services use.
//
//	D:  DACL
//	A;;GA;;;BA  Allow, Generic All, Built-in Admins
//	A;;GRGW;;;IU Allow, Generic Read+Write, Interactive Users
const pipeSDDL = "D:P(A;;GA;;;BA)(A;;GRGW;;;IU)"

// Request is the wire format received from the client. A single request
// is one line of JSON terminated by \n.
type Request struct {
	ID     uint64          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// Response is the wire format sent back to the client. One line of JSON.
type Response struct {
	ID     uint64      `json:"id"`
	OK     bool        `json:"ok"`
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

// StartVPNParams are the arguments for the start_vpn method.
type StartVPNParams struct {
	Adapter  string `json:"adapter"`
	Proxy    string `json:"proxy"`
	Server   string `json:"server"`
	MTU      int    `json:"mtu"`
}

// Server owns the pipe listener and serializes incoming commands on a
// single mutex so wintun state transitions stay consistent.
type Server struct {
	mu sync.Mutex
	ln net.Listener
}

// NewServer creates a pipe server ready to be started.
func NewServer() *Server { return &Server{} }

// Run listens on the pipe until ctx is cancelled. Each accepted connection
// is served in its own goroutine but all libtun mutations are globally
// serialized via s.mu.
func (s *Server) Run(ctx context.Context) error {
	cfg := &winio.PipeConfig{
		SecurityDescriptor: pipeSDDL,
		MessageMode:        false,
		InputBufferSize:    4096,
		OutputBufferSize:   4096,
	}
	ln, err := winio.ListenPipe(PipeName, cfg)
	if err != nil {
		return fmt.Errorf("listen pipe %s: %w", PipeName, err)
	}
	s.ln = ln

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("accept: %v", err)
			continue
		}
		go s.serveConn(conn)
	}
}

func (s *Server) serveConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return
		}
		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			writeResp(w, Response{ID: req.ID, OK: false, Error: "bad json: " + err.Error()})
			continue
		}
		resp := s.safeDispatch(&req)
		writeResp(w, resp)
	}
}

func (s *Server) safeDispatch(req *Request) (resp Response) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in dispatch method=%s: %v", req.Method, r)
			resp = Response{ID: req.ID, OK: false, Error: fmt.Sprintf("panic: %v", r)}
		}
	}()
	return s.dispatch(req)
}

func writeResp(w *bufio.Writer, resp Response) {
	b, _ := json.Marshal(resp)
	_, _ = w.Write(b)
	_, _ = w.Write([]byte{'\n'})
	_ = w.Flush()
}

func (s *Server) dispatch(req *Request) Response {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch req.Method {
	case "ping":
		return Response{ID: req.ID, OK: true, Result: "pong"}

	case "start_vpn":
		var p StartVPNParams
		if len(req.Params) > 0 {
			if err := json.Unmarshal(req.Params, &p); err != nil {
				return Response{ID: req.ID, OK: false, Error: "bad params: " + err.Error()}
			}
		}
		if p.Adapter == "" {
			p.Adapter = "Guardex"
		}
		if p.MTU == 0 {
			p.MTU = 1500
		}
		if p.Proxy == "" {
			return Response{ID: req.ID, OK: false, Error: "proxy is required"}
		}
		if err := libtun.StartWintun(p.Adapter, p.Proxy, p.Server, p.MTU); err != nil {
			return Response{ID: req.ID, OK: false, Error: err.Error()}
		}
		return Response{ID: req.ID, OK: true}

	case "stop_vpn":
		libtun.Stop()
		return Response{ID: req.ID, OK: true}

	case "is_running":
		return Response{ID: req.ID, OK: true, Result: libtun.IsStarted()}

	default:
		return Response{ID: req.ID, OK: false, Error: "unknown method: " + req.Method}
	}
}

// StopActiveSession is invoked by the service lifecycle on shutdown so
// the TUN is torn down cleanly before the process exits.
func StopActiveSession() {
	if libtun.IsStarted() {
		libtun.Stop()
	}
}
