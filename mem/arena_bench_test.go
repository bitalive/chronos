package mem

import (
	"testing"
)

var sink []byte

func BenchmarkArenaAlloc(b *testing.B) {
	arena := NewArena(1024 * 1024) // 1MB
	defer arena.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = arena.Alloc(64)
	}
}

func BenchmarkArenaAllocAligned(b *testing.B) {
	arena := NewArena(1024 * 1024) // 1MB
	defer arena.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = arena.AllocAligned(64, 64)
	}
}

func BenchmarkHeapAlloc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = make([]byte, 64)
	}
}