// +build amd64
// +build !noasm

#include "textflag.h"

// func RDTSC() uint64
TEXT ·RDTSC(SB), NOSPLIT, $0-8
    RDTSC
    SHLQ $32, DX
    ORQ DX, AX
    MOVQ AX, ret+0(FP)
    RET

// func RDTSCP() (tsc uint64, aux uint32)
TEXT ·RDTSCP(SB), NOSPLIT, $0-16
    RDTSCP
    SHLQ $32, DX
    ORQ DX, AX
    MOVQ AX, tsc+0(FP)
    MOVL CX, aux+8(FP)
    RET
