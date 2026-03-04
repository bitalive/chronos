package mem

import (
	"fmt"
	"unsafe"

	"github.com/bitalive/chronos/cpu"
	"github.com/bitalive/chronos/hash"
	"golang.org/x/sys/unix"
)

// Arena is a bump allocator for fast, sequential memory allocation.
// It uses mmap to allocate GC-invisible memory outside the Go heap.
type Arena struct {
	free      uintptr  // Current allocation pointer
	end       uintptr  // End of current buffer
	buffer    []byte   // Current buffer for GC visibility
	allocated uint64   // Total bytes in past segments
	numAllocs uint64   // Number of segments
	segments  [][]byte // All segments
	useHuge   bool
}

// Flags for mmap
const (
	MAP_HUGETLB = 0x40000 // Huge page support
)

// NewArena creates a new arena with the specified initial capacity.
func NewArena(capacity int) *Arena {
	return NewArenaExt(capacity, false)
}

// NewArenaExt creates an arena with custom options (e.g., Huge Pages).
func NewArenaExt(capacity int, useHuge bool) *Arena {
	pageSize := 4096
	flags := unix.MAP_ANON | unix.MAP_PRIVATE
	if useHuge {
		pageSize = 2 * 1024 * 1024 // 2MB
		flags |= MAP_HUGETLB
	}

	capacity = (capacity + pageSize - 1) &^ (pageSize - 1)

	segment, err := unix.Mmap(
		-1, 0, capacity,
		unix.PROT_READ|unix.PROT_WRITE,
		flags,
	)
	if err != nil {
		if useHuge {
			return NewArenaExt(capacity, false)
		}
		panic(fmt.Sprintf("nano/mem: mmap failed: %v", err))
	}

	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(segment)))
	return &Arena{
		buffer:   segment,
		free:     ptr,
		end:      ptr + uintptr(capacity),
		segments: [][]byte{segment},
		useHuge:  useHuge,
	}
}

// Alloc allocates size bytes from the arena with 8-byte alignment.
func (a *Arena) Alloc(size int) []byte {
	sz := uintptr((size + 7) &^ 7)
	if a.free+sz > a.end {
		return a.allocSlow(int(sz))
	}
	f := a.free
	a.free = f + sz
	return unsafe.Slice((*byte)(unsafe.Pointer(f)), size)
}

//go:noinline
func (a *Arena) allocSlow(size int) []byte {
	if a.buffer != nil {
		a.allocated += uint64(a.free - uintptr(unsafe.Pointer(unsafe.SliceData(a.buffer))))
	}
	a.grow(size)
	a.numAllocs++
	return a.Alloc(size)
}

// AllocAligned allocates size bytes with explicit alignment.
func (a *Arena) AllocAligned(size int, alignment uintptr) []byte {
	if size <= 0 {
		return nil
	}
	ptr := a.free
	aligned := (ptr + alignment - 1) &^ (alignment - 1)
	next := aligned + uintptr(size)
	if next <= a.end {
		a.free = next
		return unsafe.Slice((*byte)(unsafe.Pointer(aligned)), size)
	}
	return a.allocAlignedSlow(size, alignment)
}

//go:noinline
func (a *Arena) allocAlignedSlow(size int, alignment uintptr) []byte {
	if a.buffer != nil {
		a.allocated += uint64(a.free - uintptr(unsafe.Pointer(unsafe.SliceData(a.buffer))))
	}
	a.grow(size + int(alignment))
	a.numAllocs++
	return a.AllocAligned(size, alignment)
}

// AllocCacheAligned allocates size bytes aligned to a CPU cache line (64 bytes).
func (a *Arena) AllocCacheAligned(size int) []byte {
	return a.AllocAligned(size, cpu.CacheLineSize)
}

func (a *Arena) grow(minSize int) {
	currCap := 0
	if a.buffer != nil {
		currCap = len(a.buffer)
	}
	segSize := currCap
	if minSize*2 > segSize {
		segSize = minSize * 2
	}

	const pageSize = 4096
	segSize = (segSize + pageSize - 1) &^ (pageSize - 1)

	flags := unix.MAP_ANON | unix.MAP_PRIVATE
	if a.useHuge {
		flags |= MAP_HUGETLB
	}

	segment, err := unix.Mmap(
		-1, 0, segSize,
		unix.PROT_READ|unix.PROT_WRITE,
		flags,
	)
	if err != nil {
		panic(fmt.Sprintf("nano/mem: mmap failed: %v", err))
	}

	a.segments = append(a.segments, segment)
	a.buffer = segment
	a.free = uintptr(unsafe.Pointer(unsafe.SliceData(segment)))
	a.end = a.free + uintptr(segSize)
}

// Reset clears the arena.
func (a *Arena) Reset() {
	a.allocated = 0
	a.numAllocs = 0

	if len(a.segments) > 1 {
		for i := 1; i < len(a.segments); i++ {
			unix.Munmap(a.segments[i])
		}
		a.segments = a.segments[:1]
		a.buffer = a.segments[0]
	}

	if a.buffer != nil {
		a.free = uintptr(unsafe.Pointer(unsafe.SliceData(a.buffer)))
		a.end = a.free + uintptr(len(a.buffer))
	}
}

// IntegrityHash computes a WyHash of all data currently allocated in the arena.
// This is used for background verification and data integrity checks.
func (a *Arena) IntegrityHash() uint64 {
	var combined uint64
	seed := uint64(0)

	// Combine hashes of all segments up to the current allocation pointer (free)
	for i, seg := range a.segments {
		target := seg
		if i == len(a.segments)-1 {
			// Last segment: Hash only up to free pointer
			ptr := uintptr(unsafe.Pointer(unsafe.SliceData(seg)))
			allocatedInLast := a.free - ptr
			if allocatedInLast == 0 {
				continue
			}
			target = seg[:allocatedInLast]
		}

		if len(target) > 0 {
			res := hash.WyHash(target, seed)
			combined ^= res
			seed = res
		}
	}
	return combined
}

// Close releases memory.
func (a *Arena) Close() error {
	for _, seg := range a.segments {
		if err := unix.Munmap(seg); err != nil {
			return err
		}
	}
	a.segments = nil
	a.buffer = nil
	a.free = 0
	a.end = 0
	return nil
}

// Segments returns all allocated memory segments in the arena.
func (a *Arena) Segments() [][]byte {
	return a.segments
}

// Clone copies a byte slice into arena memory.
func (a *Arena) Clone(src []byte) []byte {
	n := len(src)
	if n == 0 {
		return nil
	}
	if n <= 8 {
		return a.cloneSmall(src)
	}
	dst := a.Alloc(n)
	copy(dst, src)
	return dst
}

//go:noinline
func (a *Arena) cloneSmall(src []byte) []byte {
	n := len(src)
	dst := a.Alloc(n)
	switch n {
	case 1:
		dst[0] = src[0]
	case 2:
		*(*uint16)(unsafe.Pointer(&dst[0])) = *(*uint16)(unsafe.Pointer(&src[0]))
	case 3:
		*(*uint16)(unsafe.Pointer(&dst[0])) = *(*uint16)(unsafe.Pointer(&src[0]))
		dst[2] = src[2]
	case 4:
		*(*uint32)(unsafe.Pointer(&dst[0])) = *(*uint32)(unsafe.Pointer(&src[0]))
	case 5:
		*(*uint32)(unsafe.Pointer(&dst[0])) = *(*uint32)(unsafe.Pointer(&src[0]))
		dst[4] = src[4]
	case 6:
		*(*uint32)(unsafe.Pointer(&dst[0])) = *(*uint32)(unsafe.Pointer(&src[0]))
		*(*uint16)(unsafe.Pointer(&dst[4])) = *(*uint16)(unsafe.Pointer(&src[4]))
	case 7:
		*(*uint32)(unsafe.Pointer(&dst[0])) = *(*uint32)(unsafe.Pointer(&src[0]))
		*(*uint16)(unsafe.Pointer(&dst[4])) = *(*uint16)(unsafe.Pointer(&src[4]))
		dst[6] = src[6]
	case 8:
		*(*uint64)(unsafe.Pointer(&dst[0])) = *(*uint64)(unsafe.Pointer(&src[0]))
	}
	return dst
}

// AllocCommand allocates a zero-initialized buffer for a Command.
func (a *Arena) AllocCommand() unsafe.Pointer {
	// Command struct is approximately 64 bytes
	const commandSize = 64
	buf := a.Alloc(commandSize)
	return unsafe.Pointer(&buf[0])
}
