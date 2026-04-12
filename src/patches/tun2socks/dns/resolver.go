package dns

import (
	"net"

	"github.com/xjasonlyu/tun2socks/v2/dialer"
)

func init() {
	
	
	net.DefaultResolver.PreferGo = true
	net.DefaultResolver.Dial = dialer.DefaultDialer.DialContext
}
