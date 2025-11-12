package libtun

import (
	"fmt"
	"sync"

	"github.com/xjasonlyu/tun2socks/v2/engine"
)

var (
	key     = new(engine.Key)
	started bool 
	mu      sync.Mutex
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


func Stop() {
	mu.Lock()
	defer mu.Unlock()
	engine.Stop()
	started = false       
	key = new(engine.Key) 
}


func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}
