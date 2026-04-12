package core

import (
	"fmt"
	"net/netip"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv6"
	"gvisor.dev/gvisor/pkg/tcpip/stack"

	"github.com/xjasonlyu/tun2socks/v2/core/option"
)

const (
	
	
	nicPromiscuousModeEnabled = true

	
	
	nicSpoofingEnabled = true
)


func withCreatingNIC(nicID tcpip.NICID, ep stack.LinkEndpoint) option.Option {
	return func(s *stack.Stack) error {
		if err := s.CreateNICWithOptions(nicID, ep,
			stack.NICOptions{
				Disabled: false,
				
				
				
				QDisc: nil,
			}); err != nil {
			return fmt.Errorf("create NIC: %s", err)
		}
		return nil
	}
}


func withPromiscuousMode(nicID tcpip.NICID, v bool) option.Option {
	return func(s *stack.Stack) error {
		if err := s.SetPromiscuousMode(nicID, v); err != nil {
			return fmt.Errorf("set promiscuous mode: %s", err)
		}
		return nil
	}
}



func withSpoofing(nicID tcpip.NICID, v bool) option.Option {
	return func(s *stack.Stack) error {
		if err := s.SetSpoofing(nicID, v); err != nil {
			return fmt.Errorf("set spoofing: %s", err)
		}
		return nil
	}
}


func withMulticastGroups(nicID tcpip.NICID, multicastGroups []netip.Addr) option.Option {
	return func(s *stack.Stack) error {
		if len(multicastGroups) == 0 {
			return nil
		}
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		s.AddProtocolAddress(
			nicID,
			tcpip.ProtocolAddress{
				Protocol: ipv4.ProtocolNumber,
				AddressWithPrefix: tcpip.AddressWithPrefix{
					Address:   tcpip.AddrFrom4([4]byte{0x0a, 0, 0, 0x01}),
					PrefixLen: 8,
				},
			},
			stack.AddressProperties{PEB: stack.CanBePrimaryEndpoint},
		)
		s.AddProtocolAddress(
			nicID,
			tcpip.ProtocolAddress{
				Protocol: ipv6.ProtocolNumber,
				AddressWithPrefix: tcpip.AddressWithPrefix{
					Address:   tcpip.AddrFrom16([16]byte{0xfd, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}),
					PrefixLen: 8,
				},
			},
			stack.AddressProperties{PEB: stack.CanBePrimaryEndpoint},
		)
		for _, multicastGroup := range multicastGroups {
			var err tcpip.Error
			switch {
			case multicastGroup.Is4():
				err = s.JoinGroup(ipv4.ProtocolNumber, nicID, tcpip.AddrFrom4(multicastGroup.As4()))
			case multicastGroup.Is6():
				err = s.JoinGroup(ipv6.ProtocolNumber, nicID, tcpip.AddrFrom16(multicastGroup.As16()))
			}
			if err != nil {
				return fmt.Errorf("join multicast group: %s", err)
			}
		}
		return nil
	}
}
