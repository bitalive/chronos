package sync

import (
	"sync/atomic"
)

// Spinlock is a high-performance, low-latency lock that busy-waits.
// Use this ONLY for extremely short critical sections (< 10ns).
// CRITICAL: MUST be used with runtime.LockOSThread() on the calling goroutine.
type Spinlock uint32

const (
	unlocked uint32 = 0
	locked   uint32 = 1
)

// Lock acquires the lock, spinning until available.
func (s *Spinlock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(s), unlocked, locked) {
		// Use pause instruction to reduce power consumption and improve performance
		// during busy-wait. Go doesn't expose PAUSE, so we use a small loop or ASM.
		spinWait(10)
	}
}

// Unlock releases the lock.
func (s *Spinlock) Unlock() {
	atomic.StoreUint32((*uint32)(s), unlocked)
}

// TryLock attempts to acquire the lock without spinning.
func (s *Spinlock) TryLock() bool {
	return atomic.CompareAndSwapUint32((*uint32)(s), unlocked, locked)
}

// spinWait is a helper for busy-waiting (implemented in ASM for PAUSE instruction).
func spinWait(count int)
