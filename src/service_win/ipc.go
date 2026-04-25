





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
	singbox "segment/singbox_win"
)



const PipeName = `\\.\pipe\guardex-svc`








const pipeSDDL = "D:P(A;;GA;;;BA)(A;;GRGW;;;IU)"



type Request struct {
	ID     uint64          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}


type Response struct {
	ID     uint64      `json:"id"`
	OK     bool        `json:"ok"`
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}


type StartVPNParams struct {
	Adapter  string `json:"adapter"`
	Proxy    string `json:"proxy"`
	Server   string `json:"server"`
	MTU      int    `json:"mtu"`
}



type Server struct {
	mu sync.Mutex
	ln net.Listener
}


func NewServer() *Server { return &Server{} }




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

	log.Printf("pipe listener ready on %s", PipeName)
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("accept: %v", err)
			continue
		}
		log.Printf("pipe client connected")
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
	log.Printf("rpc req id=%d method=%s params=%s", req.ID, req.Method, string(req.Params))
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in dispatch method=%s: %v", req.Method, r)
			resp = Response{ID: req.ID, OK: false, Error: fmt.Sprintf("panic: %v", r)}
		}
	}()
	resp = s.dispatch(req)
	log.Printf("rpc resp id=%d ok=%v err=%s", resp.ID, resp.OK, resp.Error)
	return resp
}

func writeResp(w *bufio.Writer, resp Response) {
	b, _ := json.Marshal(resp)
	_, _ = w.Write(b)
	_, _ = w.Write([]byte{'\n'})
	_ = w.Flush()
}

func (s *Server) dispatch(req *Request) Response {
	
	
	switch req.Method {
	case "ping":
		return Response{ID: req.ID, OK: true, Result: "pong"}
	case "is_running":
		return Response{ID: req.ID, OK: true, Result: singbox.IsStarted()}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch req.Method {
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
		
		
		
		if err := singbox.StartBridge(p.Adapter, p.Proxy, p.Server, p.MTU); err != nil {
			return Response{ID: req.ID, OK: false, Error: err.Error()}
		}
		return Response{ID: req.ID, OK: true}

	case "stop_vpn":
		singbox.StopBridge()
		return Response{ID: req.ID, OK: true}

	default:
		return Response{ID: req.ID, OK: false, Error: "unknown method: " + req.Method}
	}
}



func StopActiveSession() {
	if singbox.IsStarted() {
		singbox.StopBridge()
	}
	
	
	
	if libtun.IsStarted() {
		libtun.Stop()
	}
}
