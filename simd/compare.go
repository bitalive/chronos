package simd

// Equal returns true if a and b are equal using skeleton pattern for maximum inlining.
func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	return equalNonEmpty(a, b)
}

//go:noinline
func equalNonEmpty(a, b []byte) bool {
	// OPTIMIZATION: Check for aliasing (same memory address)
	// If both pointers and length are identical, they are equal.
	if &a[0] == &b[0] {
		return true
	}
	return equal(a, b)
}

// equal is implemented in ASM.
func equal(a, b []byte) bool
