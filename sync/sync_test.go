package sync

import (
	"sync"
	"testing"

	"github.com/bitalive/chronos/sys"
)

func TestSpinlock(t *testing.T) {
	var lock Spinlock
	var counter int

	const numGoroutines = 10
	const numIterations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				lock.Lock()
				counter++
				lock.Unlock()
			}
		}()
	}

	wg.Wait()

	if counter != numGoroutines*numIterations {
		t.Errorf("Counter value wrong: got %d, want %d", counter, numGoroutines*numIterations)
	}
}

func BenchmarkPrecisionClock(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sys.RDTSC()
	}
}
