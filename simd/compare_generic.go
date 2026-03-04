//go:build !amd64 || noasm
// +build !amd64 noasm

package simd

import "bytes"

func equal(a, b []byte) bool {
	return bytes.Equal(a, b)
}
