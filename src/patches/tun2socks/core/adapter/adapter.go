package adapter

import (
	"net"

	"gvisor.dev/gvisor/pkg/tcpip/stack"
)


type TCPConn interface {
	net.Conn

	
	ID() *stack.TransportEndpointID
}


type UDPConn interface {
	net.Conn
	net.PacketConn

	
	ID() *stack.TransportEndpointID
}
