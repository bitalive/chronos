package sync

import (
	"sync/atomic"

	"github.com/bitalive/chronos/cpu"
)

// PaddedUint64 is a uint64 value padded to a full CPU cache line (64 bytes).
// Use this for per-core metrics to prevent false sharing.
type PaddedUint64 struct {
	Value uint64
	_     [(cpu.CacheLineSize - 8)]byte // Padding to 64 bytes
}

// Add adds delta to the value and returns the new value.
func (p *PaddedUint64) Add(delta uint64) uint64 {
	return atomic.AddUint64(&p.Value, delta)
}

// Load loads the current value.
func (p *PaddedUint64) Load() uint64 {
	return atomic.LoadUint64(&p.Value)
}

// Store stores the new value.
func (p *PaddedUint64) Store(val uint64) {
	atomic.StoreUint64(&p.Value, val)
}

// PaddedInt64 is an int64 value padded to a full CPU cache line.
type PaddedInt64 struct {
	Value int64
	_     [(cpu.CacheLineSize - 8)]byte
}

// Add adds delta to the value and returns the new value.
func (p *PaddedInt64) Add(delta int64) int64 {
	return atomic.AddInt64(&p.Value, delta)
}

// Load loads the current value.
func (p *PaddedInt64) Load() int64 {
	return atomic.LoadInt64(&p.Value)
}

// Store stores the new value.
func (p *PaddedInt64) Store(val int64) {
	atomic.StoreInt64(&p.Value, val)
}
