package tunnel

import (
	"sync"

	"github.com/xjasonlyu/tun2socks/v2/proxy"
	"github.com/xjasonlyu/tun2socks/v2/tunnel/statistic"
)

var (
	_globalMu sync.RWMutex
	_globalT  *Tunnel
)

func init() {
	ReplaceGlobal(New(&proxy.Base{}, statistic.DefaultManager))
	T().ProcessAsync()
}



func T() *Tunnel {
	_globalMu.RLock()
	t := _globalT
	_globalMu.RUnlock()
	return t
}



func ReplaceGlobal(t *Tunnel) func() {
	_globalMu.Lock()
	prev := _globalT
	_globalT = t
	_globalMu.Unlock()
	return func() { ReplaceGlobal(prev) }
}
