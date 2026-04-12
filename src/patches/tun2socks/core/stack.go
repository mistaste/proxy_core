package core

import (
	"net/netip"

	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv6"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"gvisor.dev/gvisor/pkg/tcpip/transport/icmp"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
	"gvisor.dev/gvisor/pkg/tcpip/transport/udp"

	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	"github.com/xjasonlyu/tun2socks/v2/core/option"
)


type Config struct {
	
	
	LinkEndpoint stack.LinkEndpoint

	
	
	TransportHandler adapter.TransportHandler

	
	
	MulticastGroups []netip.Addr

	
	
	Options []option.Option
}


func CreateStack(cfg *Config) (*stack.Stack, error) {
	opts := []option.Option{option.WithDefault()}
	if len(opts) > 0 {
		opts = append(opts, cfg.Options...)
	}

	s := stack.New(stack.Options{
		NetworkProtocols: []stack.NetworkProtocolFactory{
			ipv4.NewProtocol,
			ipv6.NewProtocol,
		},
		TransportProtocols: []stack.TransportProtocolFactory{
			tcp.NewProtocol,
			udp.NewProtocol,
			icmp.NewProtocol4,
			icmp.NewProtocol6,
		},
	})

	
	nicID := s.NextNICID()

	opts = append(opts,
		
		
		
		
		withTCPHandler(cfg.TransportHandler.HandleTCP),
		withUDPHandler(cfg.TransportHandler.HandleUDP),

		
		withCreatingNIC(nicID, cfg.LinkEndpoint),

		
		
		
		
		
		
		
		withPromiscuousMode(nicID, nicPromiscuousModeEnabled),

		
		
		
		
		
		
		
		
		
		
		
		
		withSpoofing(nicID, nicSpoofingEnabled),

		
		
		withRouteTable(nicID),

		
		withMulticastGroups(nicID, cfg.MulticastGroups),
	)

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}
