package mem

import "unsafe"

//go:noinline
func slow() {}

type Foo struct{ free uintptr }

func (f *Foo) Alloc(size int) []byte {
	ptr := f.free
	if size > 0 {
		slow()
	}
	f.free = ptr + uintptr(size)
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), size)
}
