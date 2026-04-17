//go:build windows

package libtun

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"syscall"

	"github.com/xjasonlyu/tun2socks/v2/engine"
)

// windowsRouteState tracks routes added during Start so Stop can remove them.
type windowsRouteState struct {
	adapterName      string
	serverIP         string
	origGateway      string
	origInterfaceIdx string
	addedDefault     bool
	addedException   bool
}

var winState windowsRouteState

// Start is a no-op stub kept for API parity with the Unix build. The
// Windows flow is driven directly through the CGO export and uses
// StartWintun — the gRPC server path that takes a raw TUN fd does
// not apply on Windows because wintun is a userspace driver.
func Start(tunFD int, proxyAddress string) error {
	return fmt.Errorf("libtun.Start(fd) is not supported on Windows; use StartWintun")
}

// StartWintun initializes tun2socks using the wintun adapter with the given name,
// forwards its traffic into the given SOCKS5 proxy, and installs Windows
// routes so the adapter catches the default route while the proxy server
// itself is excepted back through the original gateway.
//
// adapterName: name of the wintun adapter to create (e.g. "Guardex").
// proxyAddress: "host:port" of local SOCKS5 proxy.
// serverIP: public IP of the remote VPN server — added as a /32 exception.
// mtu: link MTU, 0 = default 1500.
func StartWintun(adapterName, proxyAddress, serverIP string, mtu int) error {
	mu.Lock()
	defer mu.Unlock()

	if started {
		return fmt.Errorf("tun2socks has already been started")
	}
	if adapterName == "" {
		return fmt.Errorf("adapterName is required")
	}
	if mtu <= 0 {
		mtu = 1500
	}

	// tun2socks' Windows TUN driver (wintun under the hood) registers
	// under the "tun" scheme. "wintun://" is not a valid driver name.
	key.Device = fmt.Sprintf("tun://%s", adapterName)
	key.Proxy = fmt.Sprintf("socks5://%s", proxyAddress)
	key.MTU = mtu
	key.LogLevel = "info"
	engine.Insert(key)
	engine.Start()

	winState = windowsRouteState{adapterName: adapterName, serverIP: serverIP}

	gw, ifIdx, err := defaultRouteInfo()
	if err == nil {
		winState.origGateway = gw
		winState.origInterfaceIdx = ifIdx
	}

	if serverIP != "" && winState.origGateway != "" {
		// Clear any stale /32 exception left over from a previous run.
		_ = runNetsh("interface", "ipv4", "delete", "route",
			serverIP+"/32", "interface="+winState.origInterfaceIdx)
		if err := runNetsh("interface", "ipv4", "add", "route",
			serverIP+"/32", "interface="+winState.origInterfaceIdx,
			"nexthop="+winState.origGateway, "metric=1", "store=active"); err == nil {
			winState.addedException = true
		}
	}

	if err := configureAdapterIP(adapterName); err != nil {
		engine.Stop()
		platformRollback()
		return fmt.Errorf("configure adapter ip: %w", err)
	}

	// Clear any stale default route on our adapter (e.g. from a crashed
	// previous run that didn't get to rollback).
	_ = runNetsh("interface", "ipv4", "delete", "route", "0.0.0.0/0",
		"interface="+adapterName)
	if err := runNetsh("interface", "ipv4", "add", "route", "0.0.0.0/0",
		"interface="+adapterName, "metric=1", "store=active"); err != nil {
		engine.Stop()
		platformRollback()
		return fmt.Errorf("add default route: %w", err)
	}
	winState.addedDefault = true

	started = true
	return nil
}

// configureAdapterIP assigns a link-local-ish 10.200.0.2/24 IP + DNS to the
// wintun adapter. Uses the same ranges as common consumer VPN clients.
func configureAdapterIP(adapterName string) error {
	if err := runNetsh("interface", "ipv4", "set", "address",
		"name="+adapterName, "source=static", "addr=10.200.0.2",
		"mask=255.255.255.0", "gateway=10.200.0.1", "gwmetric=1"); err != nil {
		return err
	}
	_ = runNetsh("interface", "ipv4", "set", "dnsservers",
		"name="+adapterName, "source=static", "address=1.1.1.1",
		"register=none", "validate=no")
	_ = runNetsh("interface", "ipv4", "add", "dnsservers",
		"name="+adapterName, "address=8.8.8.8", "index=2",
		"validate=no")
	return nil
}

func platformStopHook() {
	platformRollback()
}

func platformRollback() {
	if winState.addedDefault {
		_ = runNetsh("interface", "ipv4", "delete", "route", "0.0.0.0/0",
			"interface="+winState.adapterName)
		winState.addedDefault = false
	}
	if winState.addedException && winState.serverIP != "" {
		_ = runNetsh("interface", "ipv4", "delete", "route",
			winState.serverIP+"/32", "interface="+winState.origInterfaceIdx)
		winState.addedException = false
	}
	winState = windowsRouteState{}
}

// defaultRouteInfo returns the current default gateway + interface index
// by parsing `route print 0.0.0.0`.
func defaultRouteInfo() (gateway, ifaceIdx string, err error) {
	cmd := exec.Command("route", "print", "0.0.0.0")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 5 && fields[0] == "0.0.0.0" && fields[1] == "0.0.0.0" {
			gw := fields[2]
			if net.ParseIP(gw) == nil {
				continue
			}
			// Interface is fields[3] (an IP of the iface). Resolve its index
			// via its IP — netsh accepts "interface=<name>" so we'll look up
			// the adapter name matching that IP below.
			ifName, e := interfaceNameByIP(fields[3])
			if e != nil {
				return gw, fields[3], nil
			}
			return gw, ifName, nil
		}
	}
	return "", "", fmt.Errorf("no default route found")
}

func interfaceNameByIP(ip string) (string, error) {
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

func runNetsh(args ...string) error {
	cmd := exec.Command("netsh", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("netsh %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return nil
}
