package mem

import (
	"unsafe"
)

// Copy8 copies 8 bytes from src to dst using a single 64-bit MOV instruction.
// Highly optimized for small keys.
func Copy8(dst, src []byte) {
	_ = dst[7]
	_ = src[7]
	*(*uint64)(unsafe.Pointer(&dst[0])) = *(*uint64)(unsafe.Pointer(&src[0]))
}

// Copy16 copies 16 bytes from src to dst using two 64-bit MOV instructions.
func Copy16(dst, src []byte) {
	_ = dst[15]
	_ = src[15]
	*(*uint64)(unsafe.Pointer(&dst[0])) = *(*uint64)(unsafe.Pointer(&src[0]))
	*(*uint64)(unsafe.Pointer(&dst[8])) = *(*uint64)(unsafe.Pointer(&src[8]))
}

// MemoryBarrier ensures all memory stores are visible to other cores.
// This is a placeholder for potential ASM implementation (SFENCE/MFENCE).
func MemoryBarrier() {
	// runtime.KeepAlive or atomic operations can sometimes act as a barrier
}
