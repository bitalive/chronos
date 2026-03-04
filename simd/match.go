package simd

// Match16 checks if two 16-byte blocks are equal.
// Performance: ~1ns
func Match16(a, b []byte) bool

// Match32 checks if two 32-byte blocks are equal.
// Performance: ~1.5ns
func Match32(a, b []byte) bool
// SearchGroup16 compares a tag against 16 control bytes using SIMD.
// Returns a 16-bit mask where each bit represents a match.
func SearchGroup16(control *uint8, tag uint8) uint16
