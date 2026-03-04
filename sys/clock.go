package sys

import (
	"sync/atomic"
	"time"
)

var (
	// cyclesPerNano is scaled by 1024 (<< 10) to maintain precision.
	// 2048 roughly equals 2.0 GHz (a safe baseline for Server processors like Xeon/EPYC).
	// This will be accurately calibrated during package init().
	cyclesPerNano uint64 = 2048
	tscOverhead   uint64 = 0 // Overhead of calling RDTSC
)

func init() {
	Calibrate()
}

// RDTSC reads the current CPU time-stamp counter.
func RDTSC() uint64

// RDTSCP reads the TSC and processor ID.
func RDTSCP() (tsc uint64, aux uint32)

// Calibrate calculates the CPU frequency and measurement overhead.
func Calibrate() {
	// 1. Measure overhead
	start := RDTSC()
	for i := 0; i < 1000; i++ {
		_ = RDTSC()
	}
	end := RDTSC()
	atomic.StoreUint64(&tscOverhead, (end-start)/1000)

	// 2. Calibrate frequency
	t1 := time.Now()
	c1 := RDTSC()
	time.Sleep(10 * time.Millisecond)
	t2 := time.Now()
	c2 := RDTSC()

	nanos := uint64(t2.Sub(t1).Nanoseconds())
	cycles := c2 - c1
	if nanos > 0 {
		atomic.StoreUint64(&cyclesPerNano, (cycles << 10 / nanos))
	}
}

// CyclesToNano converts CPU cycles to nanoseconds.
func CyclesToNano(cycles uint64) uint64 {
	cpn := atomic.LoadUint64(&cyclesPerNano)
	if cpn == 0 {
		return 0
	}
	return (cycles << 10) / cpn
}

// PreciseDuration returns the "clean" duration between two TSC readings.
func PreciseDuration(start, end uint64) uint64 {
	overhead := atomic.LoadUint64(&tscOverhead)
	if end <= start+overhead {
		return 0
	}
	return CyclesToNano(end - start - overhead)
}
