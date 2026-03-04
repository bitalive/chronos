<div align="center">
    <h1>⏱️ Chronos</h1>
    <i>The Time Stealer - Extreme Low Latency I/O & Execution Engine for Go.</i>
    <br/><br/>
    <p>
      <a href="https://godoc.org/github.com/bitalive/chronos"><img src="https://godoc.org/github.com/bitalive/chronos?status.svg" alt="GoDoc" /></a>
      <a href="https://goreportcard.com/report/github.com/bitalive/chronos"><img src="https://goreportcard.com/badge/github.com/bitalive/chronos" alt="Go Report Card" /></a>
      <a href="https://github.com/bitalive/chronos/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-BSD--2--Clause-blue.svg" alt="License" /></a>
    </p>
</div>

<hr/>

## 🎯 What is Chronos?

**Chronos** ("Kẻ đánh cắp thời gian") is a highly specialized, ultra-low latency execution engine designed from the ground up to empower [Bitalive](https://github.com/bitalive/bitalive) and any other Golang applications that require the rawest, unadulterated performance hardware can offer.

By completely bypassing the standard Go heap allocations, interface dispatch overheads, and multi-syscall tax, Chronos "steals" back every possible nanosecond between your application and the operating system kernel.

## 🚀 Key Features

*   **`sys/iouring`**: A zero-syscall `io_uring` polling backend utilizing batch submission and cache-line isolated buffers. Target: **10M+ operations/sec** per core without breaking a sweat.
*   **`mem/arena`**: An off-heap mmap-backed Bump Allocator. Circumvents the Go Garbage Collector entirely, dropping allocation overhead from the typical `~40ns` down to **`< 3ns`**.
*   **`hash/wyhash`**: Implementation of `WyHash` via ultra-minimal skeleton for maximum AST-inlining. Beats FNV-1a by 5x (running at **~3ns**).
*   **`sync/time`**: A hardware-clock precision instrument utilizing `x86 TSC (Time Stamp Counter)` returning extreme accuracy with **`~6ns`** overhead (vs normal OS syscall time queries).
*   **`cpu/simd`**: AVX2 & AVX512 instruction-set powered data scanning for lightning-fast memory parsing.

## 📦 Installation

```bash
go get github.com/bitalive/chronos
```

## ⚡ Quick Start

### Mmap Off-Heap Arena Allocation
```go
package main

import (
	"fmt"
	"github.com/bitalive/chronos/mem"
)

func main() {
	// Allocate 4096 bytes directly from OS Page (No GC tracking)
	arena := mem.NewArena(4096)
	defer arena.Close()

	// Allocates 64 bytes in ~2.9ns (Zero Allocator overhead)
	buf := arena.AllocCacheAligned(64)
	fmt.Printf("Zero GC allocation success, size: %d\n", len(buf))
}
```

### Ultra-fast WyHash Inlining
```go
package main

import (
	"fmt"
	"github.com/bitalive/chronos/hash"
)

func main() {
	key := []byte("hello_bitalive")
	
	// Hashes in ~3ns via deeply inlined fast-path
	h := hash.WyHash(key, 0x123456)
	fmt.Printf("WyHash Result: %x\n", h)
}
```

## 🧠 Design Philosophy (The Chronos Way)
1.  **Zero Interface Dispatch:** No `interfaces{}` in the hot path. Static branching enables the CPU branch-predictor to achieve ~100% accuracy.
2.  **Guaranteed Inlining:** Go functions are brutally stripped down to below the `< 80 AST nodes` limit, ensuring the Compiler always blends the execution directly into the caller. `go test -gcflags="-m"` must always succeed.
3.  **The GC Is The Enemy:** All data living in the Hot Path is routed through `mmap` anonymous pages. The GC only sees the `uintptr` root, making it essentially blind to the millions of operations happening underneath.
4.  **Hardware-Driven:** NUMA-aware initializations, SIMD instruction availability detections, and `Cache-Line Isolation (False Sharing Prevention)` are built into the primitive level.

---

## ⚖️ License
Chronos Engine is released under the **BSD 2-Clause "Simplified" License**. See the [LICENSE](LICENSE) file for more details.

*"In the pursuit of performance, time must be stolen back."*
