
package buffer

import (
	"github.com/xjasonlyu/tun2socks/v2/buffer/allocator"
)

const (
	
	MaxSegmentSize = (1 << 16) - 1

	
	
	
	
	RelayBufferSize = 20 << 10
)

var _allocator = allocator.New()


func Get(size int) []byte {
	return _allocator.Get(size)
}


func Put(buf []byte) error {
	return _allocator.Put(buf)
}
