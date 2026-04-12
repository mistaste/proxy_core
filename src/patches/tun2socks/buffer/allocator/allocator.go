package allocator

import (
	"errors"
	"math/bits"

	"github.com/xjasonlyu/tun2socks/v2/internal/pool"
)



type Allocator struct {
	buffers []*pool.Pool[[]byte]
}




func New() *Allocator {
	alloc := &Allocator{}
	alloc.buffers = make([]*pool.Pool[[]byte], 17) 
	for k := range alloc.buffers {
		i := k
		alloc.buffers[k] = pool.New(func() []byte {
			return make([]byte, 1<<uint32(i))
		})
	}
	return alloc
}


func (alloc *Allocator) Get(size int) []byte {
	if size <= 0 || size > 65536 {
		return nil
	}

	b := msb(size)
	if size == 1<<b {
		return alloc.buffers[b].Get()[:size]
	}

	return alloc.buffers[b+1].Get()[:size]
}



func (alloc *Allocator) Put(buf []byte) error {
	b := msb(cap(buf))
	if cap(buf) == 0 || cap(buf) > 65536 || cap(buf) != 1<<b {
		return errors.New("allocator Put() incorrect buffer size")
	}

	alloc.buffers[b].Put(buf)
	return nil
}


func msb(size int) uint16 {
	return uint16(bits.Len32(uint32(size)) - 1)
}
