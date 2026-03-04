package hash

import (
	"crypto/sha256"
	"hash/fnv"
	"testing"
)

func BenchmarkWyHash(b *testing.B) {
	key := []byte("this is a test key of 32 bytes--")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WyHash(key, 0)
	}
}

func BenchmarkFNV1a(b *testing.B) {
	key := []byte("this is a test key of 32 bytes--")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := fnv.New64a()
		h.Write(key)
		_ = h.Sum64()
	}
}

func BenchmarkSHA256(b *testing.B) {
	key := []byte("this is a test key of 32 bytes--")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sha256.Sum256(key)
	}
}

func BenchmarkWyHashShort(b *testing.B) {
	key := []byte("short")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WyHash(key, 0)
	}
}

func BenchmarkFNV1aShort(b *testing.B) {
	key := []byte("short")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := fnv.New64a()
		h.Write(key)
		_ = h.Sum64()
	}
}
