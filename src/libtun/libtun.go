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

// Stop stops the tun2socks engine and resets state.
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

// IsStarted checks if tun2socks has been started.
func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}
