package mem

import (
	"sync"
)

// BufferPool manages reusable byte buffers with size-based bucketing.
type BufferPool struct {
	pools []*sync.Pool // One pool per size class
	sizes []int        // Size class boundaries
}

// Standard size classes (powers of 2 for efficient bucketing)
var defaultSizes = []int{
	4 * 1024,   // 4KB  - Small values
	16 * 1024,  // 16KB - Medium values
	64 * 1024,  // 64KB - Large values
	256 * 1024, // 256KB - Very large values
}

// NewBufferPool creates a buffer pool with default size classes.
func NewBufferPool() *BufferPool {
	return NewBufferPoolWithSizes(defaultSizes)
}

// NewBufferPoolWithSizes creates a buffer pool with custom size classes.
func NewBufferPoolWithSizes(sizes []int) *BufferPool {
	pools := make([]*sync.Pool, len(sizes))

	for i, size := range sizes {
		sz := size
		pools[i] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, sz)
			},
		}
	}

	return &BufferPool{
		pools: pools,
		sizes: sizes,
	}
}

// Get retrieves a buffer of at least the requested size.
func (p *BufferPool) Get(size int) []byte {
	idx := p.findBucket(size)
	if idx < 0 {
		return make([]byte, 0, size)
	}
	buf := p.pools[idx].Get().([]byte)
	return buf[:0]
}

// Put returns a buffer to the pool for reuse.
func (p *BufferPool) Put(buf []byte) {
	cap := cap(buf)
	idx := p.findBucket(cap)
	if idx < 0 {
		return
	}
	if cap > p.sizes[idx] {
		if idx+1 < len(p.sizes) && cap <= p.sizes[idx+1] {
			idx++
		} else {
			return
		}
	}
	buf = buf[:0]
	p.pools[idx].Put(buf)
}

func (p *BufferPool) findBucket(size int) int {
	for i, bucketSize := range p.sizes {
		if size <= bucketSize {
			return i
		}
	}
	return -1
}

// Global default pool
var defaultPool = NewBufferPool()

func Get(size int) []byte {
	return defaultPool.Get(size)
}

func Put(buf []byte) {
	defaultPool.Put(buf)
}
