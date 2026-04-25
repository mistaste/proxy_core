









package singbox_win

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
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




type winRoute struct {
	serverIP         string
	origGateway      string
	origInterfaceIdx string
	addedException   bool
}













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

	
	
	
	if serverIP != "" {
		gw, ifIdx, e := defaultRouteInfo(adapterName)
		if e == nil {
			routeInfo.origGateway = gw
			routeInfo.origInterfaceIdx = ifIdx
			
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


func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}






func buildConfig(adapter, socksHost string, socksPort, mtu int) []byte {
	
	
	
	
	return fmt.Appendf(nil, `{
  "log": {"level": "warn", "disabled": false},
  "inbounds": [{
    "type": "tun",
    "tag": "tun-in",
    "interface_name": %q,
    "address": ["100.100.0.2/24"],
    "mtu": %d,
    "auto_route": true,
    "strict_route": false,
    "stack": "system"
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






type routeCandidate struct {
	gateway  string
	ifaceName string
	metric   int
}

func defaultRouteInfo(ourAdapter string) (gateway, ifaceName string, err error) {
	cmd := exec.Command("route", "print", "0.0.0.0")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	
	var candidates []routeCandidate
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
			
			continue
		}
		metric := 9999
		if m, e := strconv.Atoi(fields[4]); e == nil {
			metric = m
		}
		candidates = append(candidates, routeCandidate{gateway: gw, ifaceName: ifName, metric: metric})
	}
	if len(candidates) == 0 {
		return "", "", fmt.Errorf("no original default route found")
	}
	
	
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].metric < candidates[j].metric
	})
	return candidates[0].gateway, candidates[0].ifaceName, nil
}

func interfaceNameByIP(ip string) (string, error) {
	
	
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
