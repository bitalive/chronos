package hash

import (
	"testing"
)

// WyHash correctness tests - verify consistency and basic functionality
func TestWyHashCorrectness(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		seed  uint64
	}{
		{
			name:  "empty",
			input: []byte(""),
			seed:  0,
		},
		{
			name:  "single_byte_a",
			input: []byte("a"),
			seed:  0,
		},
		{
			name:  "three_bytes_abc",
			input: []byte("abc"),
			seed:  0,
		},
		{
			name:  "eight_bytes",
			input: []byte("12345678"),
			seed:  0,
		},
		{
			name:  "sixteen_bytes",
			input: []byte("1234567890123456"),
			seed:  0,
		},
		{
			name:  "bitalive_key_10",
			input: []byte("user:12345"),
			seed:  0x123456789abcdef0,
		},
		{
			name:  "bitalive_key_16",
			input: []byte("session:abcdef01"),
			seed:  0x123456789abcdef0,
		},
		{
			name:  "large_key_32",
			input: []byte("this_is_a_32_byte_test_key_12345"),
			seed:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test consistency - same input should always produce same output
			result1 := WyHash(tt.input, tt.seed)
			result2 := WyHash(tt.input, tt.seed)
			
			if result1 != result2 {
				t.Errorf("WyHash(%q, %x) inconsistent: %x != %x", tt.input, tt.seed, result1, result2)
			}
			
			// Log the actual values for reference
			t.Logf("WyHash(%q, %x) = %x", tt.input, tt.seed, result1)
		})
	}
}

// Test Sum64String function consistency with WyHash
func TestSum64StringConsistency(t *testing.T) {
	tests := []struct {
		name  string
		input string
		seed  uint64
	}{
		{
			name:  "empty_string",
			input: "",
			seed:  0,
		},
		{
			name:  "single_char_a",
			input: "a",
			seed:  0,
		},
		{
			name:  "three_chars_abc",
			input: "abc",
			seed:  0,
		},
		{
			name:  "bitalive_string_key",
			input: "user:12345",
			seed:  0x123456789abcdef0,
		},
		{
			name:  "sixteen_char_key",
			input: "session:abcdef01",
			seed:  0x123456789abcdef0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stringResult := Sum64String(tt.input, tt.seed)
			byteResult := WyHash([]byte(tt.input), tt.seed)
			
			if stringResult != byteResult {
				t.Errorf("Sum64String(%q, %x) = %x, but WyHash([]byte(%q), %x) = %x", 
					tt.input, tt.seed, stringResult, tt.input, tt.seed, byteResult)
			}
			
			t.Logf("Sum64String(%q, %x) = %x", tt.input, tt.seed, stringResult)
		})
	}
}

// Benchmark comparison between WyHash and Sum64String
func BenchmarkWyHashVsString(b *testing.B) {
	key := "user:session:12345678"
	keyBytes := []byte(key)
	seed := uint64(0x123456789abcdef0)

	b.Run("WyHash_bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = WyHash(keyBytes, seed)
		}
	})

	b.Run("Sum64String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Sum64String(key, seed)
		}
	})
}

// Test edge cases and boundary conditions
func TestWyHashEdgeCases(t *testing.T) {
	// Test various key lengths around boundaries
	lengths := []int{0, 1, 2, 3, 4, 7, 8, 9, 15, 16, 17, 23, 24, 25, 31, 32, 33, 63, 64, 65}
	
	for _, length := range lengths {
		key := make([]byte, length)
		for i := range key {
			key[i] = byte(i % 256)
		}
		
		// Test with different seeds
		seeds := []uint64{0, 1, 0xffffffffffffffff, 0x123456789abcdef0}
		for _, seed := range seeds {
			result1 := WyHash(key, seed)
			result2 := WyHash(key, seed)
			
			if result1 != result2 {
				t.Errorf("WyHash inconsistent for length %d, seed %x", length, seed)
			}
			
			// Test string version consistency for reasonable lengths
			if length <= 64 {
				stringResult := Sum64String(string(key), seed)
				if result1 != stringResult {
					t.Errorf("WyHash/Sum64String mismatch for length %d, seed %x: %x != %x", 
						length, seed, result1, stringResult)
				}
			}
		}
	}
}

// Test empty key special case
func TestWyHashEmpty(t *testing.T) {
	seeds := []uint64{0, 1, 0xffffffffffffffff, 0x123456789abcdef0}
	
	for _, seed := range seeds {
		expected := seed ^ wyp0
		
		result1 := WyHash([]byte{}, seed)
		result2 := WyHash(nil, seed)
		result3 := Sum64String("", seed)
		
		if result1 != expected {
			t.Errorf("WyHash([]byte{}, %x) = %x, want %x", seed, result1, expected)
		}
		if result2 != expected {
			t.Errorf("WyHash(nil, %x) = %x, want %x", seed, result2, expected)
		}
		if result3 != expected {
			t.Errorf("Sum64String(\"\", %x) = %x, want %x", seed, result3, expected)
		}
	}
}

// Benchmark specifically for Bitalive's common key patterns
func BenchmarkBitaliveKeys(b *testing.B) {
	seed := uint64(0x123456789abcdef0)
	
	// Common Bitalive key patterns
	keys := []struct {
		name string
		key  string
	}{
		{"user_id_10", "user:12345"},
		{"session_16", "session:abcdef01"},
		{"cache_key_14", "cache:data:123"},
		{"token_key_20", "token:jwt:1234567890"},
	}
	
	for _, k := range keys {
		keyBytes := []byte(k.key)
		
		b.Run(k.name+"_bytes", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = WyHash(keyBytes, seed)
			}
		})
		
		b.Run(k.name+"_string", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Sum64String(k.key, seed)
			}
		})
	}
}

// Benchmark to verify inlining effectiveness
func BenchmarkInliningEffectiveness(b *testing.B) {
	// Test empty keys (should be ultra-fast due to inlining)
	b.Run("empty_key", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = WyHash(nil, 0)
		}
	})
	
	// Test 8-byte keys (common size, should benefit from inlining)
	key8 := []byte("12345678")
	b.Run("8byte_key", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = WyHash(key8, 0)
		}
	})
	
	// Test 16-byte keys (still in fast path)
	key16 := []byte("1234567890123456")
	b.Run("16byte_key", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = WyHash(key16, 0)
		}
	})
	
	// Test 24-byte keys (goes to slow path)
	key24 := []byte("123456789012345678901234")
	b.Run("24byte_key", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = WyHash(key24, 0)
		}
	})
}