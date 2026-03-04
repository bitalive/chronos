# Chronos Engine

**Chronos** is an ultra-high-performance, zero-allocation, sub-microsecond low-level execution engine for Go. It is designed to be the foundational core for extreme-performance systems like Bitalive, bypassing the standard Go heap, scheduler inefficiencies, and standard library overheads in favor of bare-metal, cache-line aligned, and NUMA-aware operations.

## Features

- **Zero-Allocation Hot Path**: Bypasses the Go Garbage Collector utilizing custom `mmap`-backed Arena allocators.
- **Microsecond Clock Precision**: Replaces `time.Now()` with a low-overhead TSC (Time Stamp Counter) based precision clock (~6ns overhead vs ~50ns).
- **Inlined Hashing**: Provides an ultra-minimal footprint `WyHash` implementation targeting <5ns latency for sub-16-byte keys.
- **Cache-Line Padding**: Mechanical sympathy patterns using `[64]byte` (or dynamic sizes) block alignments to prevent false-sharing on multi-core systems.
- **SIMD Optimized Operations**: Vectorized data manipulation routines. 
- **Advanced Linux I/O**: Direct integration with Kernel `epoll` and `io_uring` without intermediate wrapper allocations or context-switch overheads.

## Packages

- `mem/`: Zero-allocation memory management (Arena).
- `hash/`: High-performance inline hashing (`WyHash`).
- `sync/`: Synchronization primitives and TSC-based precision clock.
- `cpu/`: CPU topology and cache-line utilities.
- `simd/`: SIMD accelerated mathematical and string operations.
- `sys/`: Low-level system interactions including `io_uring` Polling.

## Architecture Guidelines

Chronos forces extreme "**Inlining Supremacy**" and "Atomic write singularity". 
- Functions meant for hot-loops are explicitly constructed within the Go's compiler budget (~80 AST nodes).
- Objects are rigidly cache-line separated utilizing struct padding, eliminating the dreaded `False-Sharing` cache bouncing problem on modern SMP architectures.
- All dynamic memory calls strictly interface with pre-allocated block buffers, rendering the application entirely invisible to the Golang Garbage Collector.

## Status

Chronos is currently in private Alpha/Beta, serving as the foundational I/O and Memory system for Bitalive Project.

## License

MIT License
