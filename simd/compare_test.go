package simd

import (
	"bytes"
	"testing"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		a, b []byte
		want bool
	}{
		{[]byte(""), []byte(""), true},
		{[]byte("a"), []byte("a"), true},
		{[]byte("a"), []byte("b"), false},
		{[]byte("hello"), []byte("hello"), true},
		{[]byte("hello world"), []byte("hello world"), true},
		{[]byte("hello world!"), []byte("hello world?"), false},
		{make([]byte, 17), make([]byte, 17), true},
		{append(make([]byte, 16), 'a'), append(make([]byte, 16), 'b'), false},
	}

	for _, tt := range tests {
		if got := Equal(tt.a, tt.b); got != tt.want {
			t.Errorf("Equal(%q, %q) = %v; want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func BenchmarkEqual_Standard(b *testing.B) {
	s1 := make([]byte, 1024)
	s2 := make([]byte, 1024)
	copy(s2, s1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytes.Equal(s1, s2)
	}
}

func BenchmarkEqual_Nano(b *testing.B) {
	s1 := make([]byte, 1024)
	s2 := make([]byte, 1024)
	copy(s2, s1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(s1, s2)
	}
}

func BenchmarkEqual_Standard_Large(b *testing.B) {
	size := 64 * 1024 // 64KB
	s1 := make([]byte, size)
	s2 := make([]byte, size)
	copy(s2, s1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytes.Equal(s1, s2)
	}
}

func BenchmarkEqual_Nano_Large(b *testing.B) {
	size := 64 * 1024 // 64KB
	s1 := make([]byte, size)
	s2 := make([]byte, size)
	copy(s2, s1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(s1, s2)
	}
}
