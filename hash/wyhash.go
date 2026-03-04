package hash

import (
	"math/bits"
	"unsafe"
)

const (
	wyp0 = 0xa0761d6478bd642f
	wyp1 = 0xe7037ed1a0b428db
	wyp2 = 0x8ebc6af09c88c6e3
	wyp3 = 0x589965cc75374cc3
	wyp4 = 0x1d8e4e27c47d124f
)

// WyHash calculates a 64-bit hash with ultra-minimal skeleton for maximum inlining.
func WyHash(key []byte, seed uint64) uint64 {
	if len(key) == 0 {
		return seed ^ wyp0
	}
	return wyHashNonEmpty(key, seed)
}

// Sum64String calculates WyHash for string keys with ultra-minimal skeleton.
func Sum64String(key string, seed uint64) uint64 {
	if len(key) == 0 {
		return seed ^ wyp0
	}
	return sum64StringNonEmpty(key, seed)
}

//go:noinline
func wyHashNonEmpty(key []byte, seed uint64) uint64 {
	length := uint64(len(key))
	p := unsafe.Pointer(unsafe.SliceData(key))
	if length <= 16 {
		return wyHash1to16(p, length, seed)
	}
	return wyHashLarge(key, seed, length, p)
}

//go:noinline
func sum64StringNonEmpty(key string, seed uint64) uint64 {
	length := uint64(len(key))
	p := unsafe.Pointer(unsafe.StringData(key))
	if length <= 16 {
		return wyHash1to16(p, length, seed)
	}
	return wyHashLarge(unsafe.Slice((*byte)(p), length), seed, length, p)
}

func wyHash1to16(p unsafe.Pointer, length, seed uint64) uint64 {
	if length <= 8 {
		if length <= 3 {
			a := uint64(*(*byte)(p))<<16 | uint64(*(*byte)(unsafe.Add(p, length>>1)))<<8 | uint64(*(*byte)(unsafe.Add(p, length-1)))
			return wymum(wymum(a^wyp0, seed)^wyp1, length^wyp4)
		}
		a := read32(p)
		b := read32(unsafe.Add(p, length-4))
		return wymum(wymum(a^wyp0, b^seed)^wyp1, length^wyp4)
	}
	a := read64(p)
	b := read64(unsafe.Add(p, length-8))
	return wymum(wymum(a^wyp0, b^seed)^wyp1, length^wyp4)
}

//go:noinline
func wyHashLarge(key []byte, seed, length uint64, p unsafe.Pointer) uint64 {
	var a, b uint64
	if length <= 24 {
		a = read64(p)
		b = read64(unsafe.Add(p, 8))
		c := read64(unsafe.Add(p, length-8))
		return wymum(wymum(a^wyp0, b^wyp1)^wymum(c^wyp2, seed^wyp3), length^wyp4)
	}
	if length <= 32 {
		a = read64(p)
		b = read64(unsafe.Add(p, 8))
		c := read64(unsafe.Add(p, 16))
		d := read64(unsafe.Add(p, length-8))
		return wymum(wymum(a^wyp0, b^wyp1)^wymum(c^wyp2, d^wyp3), length^wyp4)
	}
	idx := uint64(0)
	see1 := seed
	for idx+32 <= length {
		seed = wymum(read64(unsafe.Add(p, idx))^wyp0, read64(unsafe.Add(p, idx+8))^wyp1) ^ seed
		see1 = wymum(read64(unsafe.Add(p, idx+16))^wyp2, read64(unsafe.Add(p, idx+24))^wyp3) ^ see1
		idx += 32
	}
	switch {
	case length-idx == 0:
		a, b = 0, 0
	case length-idx <= 8:
		a, b = read32(unsafe.Add(p, idx)), read32(unsafe.Add(p, length-4))
	case length-idx <= 16:
		a, b = read64(unsafe.Add(p, idx)), read64(unsafe.Add(p, length-8))
	case length-idx <= 24:
		seed ^= wymum(read64(unsafe.Add(p, idx))^wyp0, read64(unsafe.Add(p, idx+8))^wyp1)
		a, b = read64(unsafe.Add(p, length-8)), 0
	default: // length-idx <= 31
		seed ^= wymum(read64(unsafe.Add(p, idx))^wyp0, read64(unsafe.Add(p, idx+8))^wyp1)
		a, b = read64(unsafe.Add(p, idx+16)), read64(unsafe.Add(p, length-8))
	}
	return wymum(seed^see1^wymum(a^wyp0, b^wyp1), length^wyp4)
}

func wymum(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}

func read32(p unsafe.Pointer) uint64 {
	return uint64(*(*uint32)(p))
}

func read64(p unsafe.Pointer) uint64 {
	return *(*uint64)(p)
}
