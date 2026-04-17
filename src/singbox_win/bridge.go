//go:build windows

// Package singbox_win drives a minimal sing-box instance as a tun↔socks
// bridge on Windows, replacing tun2socks' gVisor engine. The only job of
// this bridge is to pull IP packets from a wintun adapter, terminate
// them with the system network stack (much faster than gVisor on
// Windows), and forward each TCP/UDP connection to a SOCKS5 proxy
// running locally. Xray stays as the SOCKS5 upstream and continues to
// handle VLESS/XHTTP/Reality transport — this package does not touch
// the proxy layer.
package singbox_win

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	singJSON "github.com/sagernet/sing/common/json"
)

var (
	mu        sync.Mutex
	started   bool
	instance  *box.Box
	routeInfo winRoute
)

// winRoute tracks the bypass /32 route we install for the VPN server
// IP so we can roll it back cleanly on Stop. The default route is
// owned by sing-box via auto_route=true and is torn down by Close().
type winRoute struct {
	serverIP         string
	origGateway      string
	origInterfaceIdx string
	addedException   bool
}

// StartBridge boots a sing-box instance with a single tun inbound
// (system stack) wired to a socks5 outbound pointing at the given
// localhost SOCKS5 port. It also installs a /32 bypass route for
// serverIP through the original default gateway so Xray's upstream
// connection to the VPN server does not loop back through the tun.
//
// adapterName: wintun interface name (e.g. "Guardex"). sing-box owns
// the adapter lifecycle — we no longer create it via wintun-go.
//
// socksAddr: "127.0.0.1:<port>" of the local Xray SOCKS5.
// serverIP:  public IP of the remote VPN server — /32 exception.
// mtu:       TUN MTU (0 → 1500).
func StartBridge(adapterName, socksAddr, serverIP string, mtu int) error {
	mu.Lock()
	defer mu.Unlock()

	if started {
		return fmt.Errorf("singbox bridge already started")
	}
	if adapterName == "" {
		return fmt.Errorf("adapterName is required")
	}
	if socksAddr == "" {
		return fmt.Errorf("socksAddr is required")
	}
	if mtu <= 0 {
		mtu = 1500
	}

	host, port, err := splitHostPort(socksAddr)
	if err != nil {
		return fmt.Errorf("parse socksAddr: %w", err)
	}

	cfg := buildConfig(adapterName, host, port, mtu)

	ctx := include.Context(context.Background())
	var opts option.Options
	decoder := singJSON.NewDecoderContext(ctx, bytes.NewReader(cfg))
	if err := decoder.Decode(&opts); err != nil {
		return fmt.Errorf("decode singbox config: %w", err)
	}

	b, err := box.New(box.Options{Options: opts, Context: ctx})
	if err != nil {
		return fmt.Errorf("singbox new: %w", err)
	}
	if err := b.Start(); err != nil {
		_ = b.Close()
		return fmt.Errorf("singbox start: %w", err)
	}

	instance = b
	routeInfo = winRoute{serverIP: serverIP}

	// Best-effort /32 exception for the remote VPN server IP so Xray's
	// upstream vless connection exits through the real NIC rather than
	// looping through our tun. Matches what libtun did previously.
	if serverIP != "" {
		gw, ifIdx, e := defaultRouteInfo(adapterName)
		if e == nil {
			routeInfo.origGateway = gw
			routeInfo.origInterfaceIdx = ifIdx
			// Scrub any stale /32 left over from a crashed run.
			_ = runNetsh("interface", "ipv4", "delete", "route",
				serverIP+"/32", "interface="+ifIdx)
			if e := runNetsh("interface", "ipv4", "add", "route",
				serverIP+"/32", "interface="+ifIdx,
				"nexthop="+gw, "metric=0", "store=active"); e == nil {
				routeInfo.addedException = true
			}
		}
	}

	started = true
	return nil
}

// StopBridge closes sing-box (which tears down the wintun adapter and
// removes the default route it installed) and removes the /32 bypass
// route for the VPN server.
func StopBridge() {
	mu.Lock()
	defer mu.Unlock()
	if !started {
		return
	}
	if instance != nil {
		_ = instance.Close()
		instance = nil
	}
	if routeInfo.addedException && routeInfo.serverIP != "" {
		_ = runNetsh("interface", "ipv4", "delete", "route",
			routeInfo.serverIP+"/32",
			"interface="+routeInfo.origInterfaceIdx)
	}
	routeInfo = winRoute{}
	started = false
}

// IsStarted reports whether the bridge is currently running.
func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}

// buildConfig returns a minimal sing-box JSON config with one tun
// inbound (system stack) and one socks outbound targeting our local
// Xray SOCKS5. auto_route=true makes sing-box install a default route
// through the tun. strict_route=false keeps Windows tolerant of IPv6
// and link-local traffic so nothing else dies.
func buildConfig(adapter, socksHost string, socksPort, mtu int) []byte {
	return fmt.Appendf(nil, `{
  "log": {"level": "warn", "disabled": false},
  "inbounds": [{
    "type": "tun",
    "tag": "tun-in",
    "interface_name": %q,
    "address": ["10.200.0.2/24"],
    "mtu": %d,
    "auto_route": true,
    "strict_route": false,
    "stack": "system",
    "sniff": true
  }],
  "outbounds": [{
    "type": "socks",
    "tag": "socks-out",
    "server": %q,
    "server_port": %d,
    "version": "5"
  }],
  "route": {
    "final": "socks-out",
    "auto_detect_interface": true
  }
}`, adapter, mtu, socksHost, socksPort)
}

// defaultRouteInfo returns the default gateway + interface name by
// parsing `route print 0.0.0.0`. The caller must pass its own
// adapter's name so we skip past it (sing-box may already have added
// the default route via the tun, and we want the *original* gateway,
// not the tun's).
func defaultRouteInfo(ourAdapter string) (gateway, ifaceName string, err error) {
	cmd := exec.Command("route", "print", "0.0.0.0")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	// Format: "   0.0.0.0        0.0.0.0      <gw>      <iface_ip>   <metric>"
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[0] != "0.0.0.0" || fields[1] != "0.0.0.0" {
			continue
		}
		gw := fields[2]
		ifIP := fields[3]
		ifName, e := interfaceNameByIP(ifIP)
		if e != nil {
			continue
		}
		if strings.EqualFold(ifName, ourAdapter) {
			// Skip our own tun's default route entry, keep looking.
			continue
		}
		return gw, ifName, nil
	}
	return "", "", fmt.Errorf("no original default route found")
}

func interfaceNameByIP(ip string) (string, error) {
	// Delegated to the stdlib helper in libtun_windows.go; reuse via
	// a small local shim to avoid cross-package imports.
	return interfaceNameByIPImpl(ip)
}

func runNetsh(args ...string) error {
	cmd := exec.Command("netsh", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("netsh %s: %v: %s",
			strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return nil
}

func splitHostPort(addr string) (string, int, error) {
	i := strings.LastIndex(addr, ":")
	if i < 0 {
		return "", 0, fmt.Errorf("missing port in %q", addr)
	}
	host := addr[:i]
	if host == "" {
		host = "127.0.0.1"
	}
	portStr := addr[i+1:]
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil || port <= 0 {
		return "", 0, fmt.Errorf("invalid port %q", portStr)
	}
	return host, port, nil
}
