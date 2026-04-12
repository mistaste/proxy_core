package tunnel

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	"github.com/xjasonlyu/tun2socks/v2/proxy"
	"github.com/xjasonlyu/tun2socks/v2/tunnel/statistic"
)

const (
	
	tcpConnectTimeout = 5 * time.Second
	
	tcpWaitTimeout = 60 * time.Second
	
	udpSessionTimeout = 60 * time.Second
)

var _ adapter.TransportHandler = (*Tunnel)(nil)

type Tunnel struct {
	
	tcpQueue chan adapter.TCPConn
	udpQueue chan adapter.UDPConn

	
	udpTimeout *atomic.Duration

	
	dialerMu sync.RWMutex
	dialer   proxy.Dialer

	
	manager *statistic.Manager

	procOnce   sync.Once
	procCancel context.CancelFunc
}

func New(dialer proxy.Dialer, manager *statistic.Manager) *Tunnel {
	return &Tunnel{
		tcpQueue:   make(chan adapter.TCPConn),
		udpQueue:   make(chan adapter.UDPConn),
		udpTimeout: atomic.NewDuration(udpSessionTimeout),
		dialer:     dialer,
		manager:    manager,
		procCancel: func() {  },
	}
}


func (t *Tunnel) TCPIn() chan<- adapter.TCPConn {
	return t.tcpQueue
}


func (t *Tunnel) UDPIn() chan<- adapter.UDPConn {
	return t.udpQueue
}

func (t *Tunnel) HandleTCP(conn adapter.TCPConn) {
	t.TCPIn() <- conn
}

func (t *Tunnel) HandleUDP(conn adapter.UDPConn) {
	t.UDPIn() <- conn
}

func (t *Tunnel) process(ctx context.Context) {
	for {
		select {
		case conn := <-t.tcpQueue:
			go t.handleTCPConn(conn)
		case conn := <-t.udpQueue:
			go t.handleUDPConn(conn)
		case <-ctx.Done():
			return
		}
	}
}


func (t *Tunnel) ProcessAsync() {
	t.procOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		t.procCancel = cancel
		go t.process(ctx)
	})
}


func (t *Tunnel) Close() {
	t.procCancel()
}

func (t *Tunnel) Dialer() proxy.Dialer {
	t.dialerMu.RLock()
	d := t.dialer
	t.dialerMu.RUnlock()
	return d
}

func (t *Tunnel) SetDialer(dialer proxy.Dialer) {
	t.dialerMu.Lock()
	t.dialer = dialer
	t.dialerMu.Unlock()
}

func (t *Tunnel) SetUDPTimeout(timeout time.Duration) {
	t.udpTimeout.Store(timeout)
}
