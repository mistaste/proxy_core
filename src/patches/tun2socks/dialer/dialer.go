package dialer

import (
	"context"
	"net"
	"syscall"

	"go.uber.org/atomic"
)


var DefaultDialer = &Dialer{
	InterfaceName:  atomic.NewString(""),
	InterfaceIndex: atomic.NewInt32(0),
	RoutingMark:    atomic.NewInt32(0),
}

type Dialer struct {
	InterfaceName  *atomic.String
	InterfaceIndex *atomic.Int32
	RoutingMark    *atomic.Int32
}

type Options struct {
	
	
	
	InterfaceName string

	
	
	
	InterfaceIndex int

	
	
	
	RoutingMark int
}


func DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return DefaultDialer.DialContext(ctx, network, address)
}


func ListenPacket(network, address string) (net.PacketConn, error) {
	return DefaultDialer.ListenPacket(network, address)
}

func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.DialContextWithOptions(ctx, network, address, &Options{
		InterfaceName:  d.InterfaceName.Load(),
		InterfaceIndex: int(d.InterfaceIndex.Load()),
		RoutingMark:    int(d.RoutingMark.Load()),
	})
}

func (_ *Dialer) DialContextWithOptions(ctx context.Context, network, address string, opts *Options) (net.Conn, error) {
	d := &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return setSocketOptions(network, address, c, opts)
		},
	}
	return d.DialContext(ctx, network, address)
}

func (d *Dialer) ListenPacket(network, address string) (net.PacketConn, error) {
	return d.ListenPacketWithOptions(network, address, &Options{
		InterfaceName:  d.InterfaceName.Load(),
		InterfaceIndex: int(d.InterfaceIndex.Load()),
		RoutingMark:    int(d.RoutingMark.Load()),
	})
}

func (_ *Dialer) ListenPacketWithOptions(network, address string, opts *Options) (net.PacketConn, error) {
	lc := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return setSocketOptions(network, address, c, opts)
		},
	}
	return lc.ListenPacket(context.Background(), network, address)
}
