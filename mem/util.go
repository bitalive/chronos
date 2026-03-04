package mem

// Memset sets all bytes in buf to value.
func Memset(buf []byte, value byte) {
	for i := range buf {
		buf[i] = value
	}
}

// MemsetZero clears all bytes in buf to zero.
func MemsetZero(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

// Memcopy copies src to dst.
func Memcopy(dst, src []byte) int {
	return copy(dst, src)
}

// Clone creates a copy of the byte slice.
func Clone(b []byte) []byte {
	if b == nil {
		return nil
	}
	clone := make([]byte, len(b))
	copy(clone, b)
	return clone
}

// CloneToArena creates a copy of the byte slice in the arena.
func CloneToArena(arena *Arena, b []byte) []byte {
	return arena.Clone(b)
}

// Grow grows the slice to at least n bytes capacity.
func Grow(b []byte, n int) []byte {
	if cap(b) >= n {
		return b
	}
	newCap := cap(b) * 2
	if newCap < n {
		newCap = n
	}
	newSlice := make([]byte, len(b), newCap)
	copy(newSlice, b)
	return newSlice
}
