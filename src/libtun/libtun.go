package libtun

import (
	"sync"

	"github.com/xjasonlyu/tun2socks/v2/engine"
)

var (
	key     = new(engine.Key)
	started bool
	mu      sync.Mutex
)


func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if !started {
		return
	}
	engine.Stop()
	started = false
	key = new(engine.Key)
	platformStopHook()
}


func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}
