package option

import (
	"fmt"

	"golang.org/x/time/rate"
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv6"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
)

const (
	
	defaultTimeToLive uint8 = 64

	
	
	ipForwardingEnabled = true

	
	
	icmpBurst = 50

	
	
	icmpLimit rate.Limit = 1000

	
	
	tcpCongestionControlAlgorithm = "reno" 

	
	
	tcpDelayEnabled = false

	
	
	tcpModerateReceiveBufferEnabled = false

	
	
	tcpSACKEnabled = true

	
	tcpRecovery = tcpip.TCPRACKLossDetection

	
	tcpMinBufferSize = tcp.MinBufferSize

	
	tcpMaxBufferSize = tcp.MaxBufferSize

	
	
	tcpDefaultSendBufferSize = tcp.DefaultSendBufferSize

	
	
	tcpDefaultReceiveBufferSize = tcp.DefaultReceiveBufferSize
)

type Option func(*stack.Stack) error


func WithDefault() Option {
	return func(s *stack.Stack) error {
		opts := []Option{
			WithDefaultTTL(defaultTimeToLive),
			WithForwarding(ipForwardingEnabled),

			
			WithICMPBurst(icmpBurst), WithICMPLimit(icmpLimit),

			
			
			
			
			
			WithTCPSendBufferSizeRange(tcpMinBufferSize, tcpDefaultSendBufferSize, tcpMaxBufferSize),
			WithTCPReceiveBufferSizeRange(tcpMinBufferSize, tcpDefaultReceiveBufferSize, tcpMaxBufferSize),

			WithTCPCongestionControl(tcpCongestionControlAlgorithm),
			WithTCPDelay(tcpDelayEnabled),

			
			
			WithTCPModerateReceiveBuffer(tcpModerateReceiveBufferEnabled),

			
			
			WithTCPSACKEnabled(tcpSACKEnabled),

			
			
			
			
			
			
			
			
			
			WithTCPRecovery(tcpRecovery),
		}

		for _, opt := range opts {
			if err := opt(s); err != nil {
				return err
			}
		}

		return nil
	}
}


func WithDefaultTTL(ttl uint8) Option {
	return func(s *stack.Stack) error {
		opt := tcpip.DefaultTTLOption(ttl)
		if err := s.SetNetworkProtocolOption(ipv4.ProtocolNumber, &opt); err != nil {
			return fmt.Errorf("set ipv4 default TTL: %s", err)
		}
		if err := s.SetNetworkProtocolOption(ipv6.ProtocolNumber, &opt); err != nil {
			return fmt.Errorf("set ipv6 default TTL: %s", err)
		}
		return nil
	}
}


func WithForwarding(v bool) Option {
	return func(s *stack.Stack) error {
		if err := s.SetForwardingDefaultAndAllNICs(ipv4.ProtocolNumber, v); err != nil {
			return fmt.Errorf("set ipv4 forwarding: %s", err)
		}
		if err := s.SetForwardingDefaultAndAllNICs(ipv6.ProtocolNumber, v); err != nil {
			return fmt.Errorf("set ipv6 forwarding: %s", err)
		}
		return nil
	}
}



func WithICMPBurst(burst int) Option {
	return func(s *stack.Stack) error {
		s.SetICMPBurst(burst)
		return nil
	}
}



func WithICMPLimit(limit rate.Limit) Option {
	return func(s *stack.Stack) error {
		s.SetICMPLimit(limit)
		return nil
	}
}


func WithTCPSendBufferSize(size int) Option {
	return func(s *stack.Stack) error {
		sndOpt := tcpip.TCPSendBufferSizeRangeOption{Min: tcpMinBufferSize, Default: size, Max: tcpMaxBufferSize}
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &sndOpt); err != nil {
			return fmt.Errorf("set TCP send buffer size range: %s", err)
		}
		return nil
	}
}


func WithTCPSendBufferSizeRange(a, b, c int) Option {
	return func(s *stack.Stack) error {
		sndOpt := tcpip.TCPSendBufferSizeRangeOption{Min: a, Default: b, Max: c}
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &sndOpt); err != nil {
			return fmt.Errorf("set TCP send buffer size range: %s", err)
		}
		return nil
	}
}


func WithTCPReceiveBufferSize(size int) Option {
	return func(s *stack.Stack) error {
		rcvOpt := tcpip.TCPReceiveBufferSizeRangeOption{Min: tcpMinBufferSize, Default: size, Max: tcpMaxBufferSize}
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &rcvOpt); err != nil {
			return fmt.Errorf("set TCP receive buffer size range: %s", err)
		}
		return nil
	}
}


func WithTCPReceiveBufferSizeRange(a, b, c int) Option {
	return func(s *stack.Stack) error {
		rcvOpt := tcpip.TCPReceiveBufferSizeRangeOption{Min: a, Default: b, Max: c}
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &rcvOpt); err != nil {
			return fmt.Errorf("set TCP receive buffer size range: %s", err)
		}
		return nil
	}
}


func WithTCPCongestionControl(cc string) Option {
	return func(s *stack.Stack) error {
		opt := tcpip.CongestionControlOption(cc)
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &opt); err != nil {
			return fmt.Errorf("set TCP congestion control algorithm: %s", err)
		}
		return nil
	}
}


func WithTCPDelay(v bool) Option {
	return func(s *stack.Stack) error {
		opt := tcpip.TCPDelayEnabled(v)
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &opt); err != nil {
			return fmt.Errorf("set TCP delay: %s", err)
		}
		return nil
	}
}


func WithTCPModerateReceiveBuffer(v bool) Option {
	return func(s *stack.Stack) error {
		opt := tcpip.TCPModerateReceiveBufferOption(v)
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &opt); err != nil {
			return fmt.Errorf("set TCP moderate receive buffer: %s", err)
		}
		return nil
	}
}


func WithTCPSACKEnabled(v bool) Option {
	return func(s *stack.Stack) error {
		opt := tcpip.TCPSACKEnabled(v)
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &opt); err != nil {
			return fmt.Errorf("set TCP SACK: %s", err)
		}
		return nil
	}
}


func WithTCPRecovery(v tcpip.TCPRecovery) Option {
	return func(s *stack.Stack) error {
		if err := s.SetTransportProtocolOption(tcp.ProtocolNumber, &v); err != nil {
			return fmt.Errorf("set TCP Recovery: %s", err)
		}
		return nil
	}
}
