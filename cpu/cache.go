package cpu

import (
	"unsafe"
)

// CacheLineSize is the size of a CPU cache line in bytes.
// Modern x86-64 and ARM64 CPUs use 64-byte cache lines.
const CacheLineSize = 64

// CacheLinePadding is a helper type for padding structs to cache line boundaries.
type CacheLinePadding [CacheLineSize]byte

// Padding calculates the number of bytes needed to align size to cache line boundary.
func Padding(size uintptr) uintptr {
	remainder := size % CacheLineSize
	if remainder == 0 {
		return 0
	}
	return CacheLineSize - remainder
}

// IsAligned checks if a pointer is aligned to cache line boundary.
func IsAligned(ptr unsafe.Pointer) bool {
	return uintptr(ptr)%CacheLineSize == 0
}

// AlignedSize returns the size rounded up to the next cache line boundary.
func AlignedSize(size uintptr) uintptr {
	return (size + CacheLineSize - 1) &^ (CacheLineSize - 1)
}

// Offset returns the offset within a cache line for a given pointer.
func Offset(ptr unsafe.Pointer) uintptr {
	return uintptr(ptr) % CacheLineSize
}
