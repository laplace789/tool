package cursor

import (
	"golang.org/x/sys/cpu"
	"sync/atomic"
	"unsafe"
)

const CacheLinePadSize = unsafe.Sizeof(cpu.CacheLinePad{})

type Cursor struct {
	//cache 填滿 CacheLinePadSize - 8個,避免false sharing
	v uint64
	_ [CacheLinePadSize - 8]byte
}

func newCursor() *Cursor {
	return &Cursor{}
}

func (c *Cursor) Increment() uint64 {
	return atomic.AddUint64(&c.v, 1)
}

func (c *Cursor) AtomicLoad() uint64 {
	return atomic.LoadUint64(&c.v)
}

func (c *Cursor) load() uint64 {
	return c.v
}

func (c *Cursor) Store(expectVal, newVal uint64) bool {
	return atomic.CompareAndSwapUint64(&c.v, expectVal, newVal)
}
