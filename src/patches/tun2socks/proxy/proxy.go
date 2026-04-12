
package proxy

import (
	"context"
	"net"
	"time"

	M "github.com/xjasonlyu/tun2socks/v2/metadata"
	"github.com/xjasonlyu/tun2socks/v2/proxy/proto"
)

const (
	tcpConnectTimeout = 5 * time.Second
)

var _defaultDialer Dialer = &Base{}

type Dialer interface {
	DialContext(context.Context, *M.Metadata) (net.Conn, error)
	DialUDP(*M.Metadata) (net.PacketConn, error)
}

type Proxy interface {
	Dialer
	Addr() string
	Proto() proto.Proto
}


func SetDialer(d Dialer) {
	_defaultDialer = d
}


func Dial(metadata *M.Metadata) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpConnectTimeout)
	defer cancel()
	return _defaultDialer.DialContext(ctx, metadata)
}


func DialContext(ctx context.Context, metadata *M.Metadata) (net.Conn, error) {
	return _defaultDialer.DialContext(ctx, metadata)
}


func DialUDP(metadata *M.Metadata) (net.PacketConn, error) {
	return _defaultDialer.DialUDP(metadata)
}
