package sys

import (
	"fmt"
	"sync/atomic"
)

// ProfilerEntry is a single point of measurement.
type ProfilerEntry struct {
	Name string
	TSC  uint64
}

// NanoProfiler is a zero-allocation, fixed-size path recorder.
type NanoProfiler struct {
	entries []ProfilerEntry
	count   uint32
	max     uint32
}

// NewProfiler creates a profiler with pre-allocated entries.
func NewProfiler(maxEntries int) *NanoProfiler {
	return &NanoProfiler{
		entries: make([]ProfilerEntry, maxEntries),
		max:     uint32(maxEntries),
	}
}

// Mark records a timestamp with a name.
// This is designed to be extremely fast.
func (p *NanoProfiler) Mark(name string) {
	idx := atomic.AddUint32(&p.count, 1) - 1
	if idx < p.max {
		p.entries[idx].Name = name
		p.entries[idx].TSC = RDTSC()
	}
}

// Result calculates and prints the duration between marks.
func (p *NanoProfiler) Result() {
	count := atomic.LoadUint32(&p.count)
	if count > p.max {
		count = p.max
	}

	fmt.Println("\n--- NANO PROFILER REPORT ---")
	for i := uint32(1); i < count; i++ {
		prev := p.entries[i-1]
		curr := p.entries[i]
		duration := PreciseDuration(prev.TSC, curr.TSC)
		fmt.Printf("[%d] %s -> %s: %d ns\n", i, prev.Name, curr.Name, duration)
	}
	fmt.Println("----------------------------")
}

// GetCount returns the current number of marks.
func (p *NanoProfiler) GetCount() uint32 {
	return atomic.LoadUint32(&p.count)
}

// Reset clears the profiler.
func (p *NanoProfiler) Reset() {
	atomic.StoreUint32(&p.count, 0)
}
