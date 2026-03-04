package simd

// FastValidateHeader checks Magic (0xBA01) and Version (0x01) in one shot.
// It returns true if valid, false otherwise.
// Performance: < 2ns
func FastValidateHeader(buf []byte) bool
