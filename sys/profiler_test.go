package sys

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestProfiler(t *testing.T) {
	p := NewProfiler(10)
	p.Mark("Start")
	time.Sleep(10 * time.Millisecond)
	p.Mark("End")

	p.Result()
}

func TestCalibration(t *testing.T) {
	Calibrate()
	t.Logf("TSC Overhead: %d cycles", atomic.LoadUint64(&tscOverhead))
	t.Logf("Cycles/Nano (fixed-10): %d", atomic.LoadUint64(&cyclesPerNano))

	s := RDTSC()
	time.Sleep(10 * time.Millisecond)
	e := RDTSC()

	nanos := PreciseDuration(s, e)
	t.Logf("Measured 10ms as: %d ns", nanos)
}
