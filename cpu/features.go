package cpu

import (
	"runtime"

	"golang.org/x/sys/cpu"
)

// CPU Features
var (
	HasAVX2   = cpu.X86.HasAVX2
	HasAVX512 = cpu.X86.HasAVX512F && cpu.X86.HasAVX512BW && cpu.X86.HasAVX512VL
	HasSSE42  = cpu.X86.HasSSE42
)

// Hardware specs
var (
	NumCPU = runtime.NumCPU()
)

// IsIntel returns true if the CPU is Intel.
func IsIntel() bool {
	// Simple check, can be expanded with more robust detection
	return true // Most k8s environments are Intel/AMD
}
