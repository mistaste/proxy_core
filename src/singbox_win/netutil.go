//go:build windows

package singbox_win

import (
	"fmt"
	"net"
	"strings"
)

// interfaceNameByIPImpl finds the adapter name whose primary IPv4
// matches ip. Used to resolve the friendly interface name that netsh
// accepts when installing the /32 server exception route.
func interfaceNameByIPImpl(ip string) (string, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, ifc := range ifs {
		addrs, _ := ifc.Addrs()
		for _, a := range addrs {
			if strings.HasPrefix(a.String(), ip+"/") {
				return ifc.Name, nil
			}
		}
	}
	return "", fmt.Errorf("no interface with ip %s", ip)
}
