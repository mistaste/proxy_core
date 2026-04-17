

package libtun

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"syscall"

	"github.com/xjasonlyu/tun2socks/v2/engine"
)


type windowsRouteState struct {
	adapterName      string
	serverIP         string
	origGateway      string
	origInterfaceIdx string
	addedDefault     bool
	addedException   bool
}

var winState windowsRouteState





func Start(tunFD int, proxyAddress string) error {
	return fmt.Errorf("libtun.Start(fd) is not supported on Windows; use StartWintun")
}










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

	key.Device = fmt.Sprintf("wintun://%s", adapterName)
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
