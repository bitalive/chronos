package sys

import (
	"syscall"
)

// RawSyscall6 is a wrapper for syscall.RawSyscall6 to avoid Go scheduler overhead.
// Use only for non-blocking syscalls.
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	return syscall.RawSyscall6(trap, a1, a2, a3, a4, a5, a6)
}

// LockOSThread is a wrapper for runtime.LockOSThread but could be optimized
// with ASM if we need even lower overhead or specific CPU pinning.
// For now, it serves as a placeholder for the library structure.
func LockOSThread() {
	// runtime.LockOSThread()
}
