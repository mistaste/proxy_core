

package libtun

import (
	"fmt"

	"github.com/xjasonlyu/tun2socks/v2/engine"
)


func Start(tunFD int, proxyAddress string) error {
	mu.Lock()
	defer mu.Unlock()

	if started {
		return fmt.Errorf("tun2socks has already been started")
	}
	started = true
	key.Device = fmt.Sprintf("fd://%d", tunFD)
	key.Proxy = fmt.Sprintf("socks5://%s", proxyAddress)
	key.MTU = 1500
	key.LogLevel = "info"
	engine.Insert(key)
	engine.Start()
	return nil
}

func platformStopHook() {}
