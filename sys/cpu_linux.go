//go:build linux

package sys

import (
	"fmt"
	"runtime"

	"golang.org/x/sys/unix"
)

// PinToCore pins the current goroutine's thread to a specific CPU core.
// This is critical for cache locality and preventing context switches.
// This function calls runtime.LockOSThread() internally.
func PinToCore(coreID int) error {
	numCPU := runtime.NumCPU()
	if coreID < 0 || coreID >= numCPU {
		return fmt.Errorf("nano/sys: invalid core ID %d (available: 0-%d)", coreID, numCPU-1)
	}

	runtime.LockOSThread()

	var cpuSet unix.CPUSet
	cpuSet.Zero()
	cpuSet.Set(coreID)

	// 0 means current thread
	if err := unix.SchedSetaffinity(0, &cpuSet); err != nil {
		runtime.UnlockOSThread()
		return fmt.Errorf("nano/sys: failed to set CPU affinity: %w", err)
	}

	return nil
}

// GetAffinity returns the CPU cores that the current thread is allowed to run on.
func GetAffinity() ([]int, error) {
	var cpuSet unix.CPUSet
	if err := unix.SchedGetaffinity(0, &cpuSet); err != nil {
		return nil, fmt.Errorf("nano/sys: failed to get CPU affinity: %w", err)
	}

	var cores []int
	for i := 0; i < runtime.NumCPU(); i++ {
		if cpuSet.IsSet(i) {
			cores = append(cores, i)
		}
	}
	return cores, nil
}
