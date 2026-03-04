// +build amd64
// +build !noasm

#include "textflag.h"

// func equal(a, b []byte) bool
TEXT ·equal(SB), NOSPLIT, $0-49
    MOVQ a_base+0(FP), SI
    MOVQ a_len+8(FP), BX
    MOVQ b_base+24(FP), DI

    // Check if lengths are the same (already done in Go)

    // Large data path (> 128 bytes)
    CMPQ BX, $128
    JGE loop_128_start

    // Medium data path (16 - 128 bytes)
    CMPQ BX, $16
    JGE loop_16_check

    // Small data path (< 16 bytes) - BRANCHLESS TAIL
    CMPQ BX, $0
    JE done_equal
    CMPQ BX, $8
    JGE tail_8
    CMPQ BX, $4
    JGE tail_4
    CMPQ BX, $1
    JGE tail_1_3
    JMP done_equal

tail_8:
    MOVQ (SI), AX
    MOVQ (DI), CX
    CMPQ AX, CX
    JNE not_equal
    // Overlapping read for the remaining bytes
    MOVQ -8(SI)(BX*1), AX
    MOVQ -8(DI)(BX*1), CX
    CMPQ AX, CX
    JNE not_equal
    JMP done_equal

tail_4:
    MOVL (SI), AX
    MOVL (DI), CX
    CMPL AX, CX
    JNE not_equal
    MOVL -4(SI)(BX*1), AX
    MOVL -4(DI)(BX*1), CX
    CMPL AX, CX
    JNE not_equal
    JMP done_equal

tail_1_3:
    MOVB (SI), AL
    MOVB (DI), CL
    CMPB AL, CL
    JNE not_equal
    MOVB -1(SI)(BX*1), AL
    MOVB -1(DI)(BX*1), CL
    CMPB AL, CL
    JNE not_equal
    // One more check for the middle byte if len == 3
    CMPQ BX, $3
    JNE done_equal
    MOVB 1(SI), AL
    MOVB 1(DI), CL
    CMPB AL, CL
    JNE not_equal
    JMP done_equal

loop_128_start:
    MOVQ SI, R8
    ADDQ BX, R8             // R8 = end of a
    PREFETCHT0 (SI)
    PREFETCHT0 (DI)

loop_128:
    VMOVDQU (SI), Y0
    VMOVDQU (DI), Y1
    VPCMPEQB Y0, Y1, Y4
    VMOVDQU 32(SI), Y2
    VMOVDQU 32(DI), Y3
    VPCMPEQB Y2, Y3, Y5
    VMOVDQU 64(SI), Y0
    VMOVDQU 64(DI), Y1
    VPCMPEQB Y0, Y1, Y6
    VMOVDQU 96(SI), Y2
    VMOVDQU 96(DI), Y3
    VPCMPEQB Y2, Y3, Y7

    VPAND Y4, Y5, Y8
    VPAND Y6, Y7, Y9
    VPAND Y8, Y9, Y10
    VPMOVMSKB Y10, DX
    CMPL DX, $0xFFFFFFFF
    JNE not_equal_avx

    ADDQ $128, SI
    ADDQ $128, DI
    SUBQ $128, BX
    CMPQ BX, $128
    JGE loop_128

    // Fallthrough to 16 check
    MOVQ R8, AX             // Restore end pointer

loop_16_check:
    MOVQ SI, AX
    ADDQ BX, AX             // Recalculate end for medium path
loop_16:
    CMPQ BX, $16
    JL tail_overlap_16
    VMOVDQU (SI), X0
    VMOVDQU (DI), X1
    VPCMPEQB X0, X1, X2
    VPMOVMSKB X2, DX
    CMPW DX, $0xFFFF
    JNE not_equal_avx
    ADDQ $16, SI
    ADDQ $16, DI
    SUBQ $16, BX
    JMP loop_16

tail_overlap_16:
    // One final 16-byte overlapping read if we have anything left
    CMPQ BX, $0
    JE done_equal
    VMOVDQU -16(SI)(BX*1), X0
    VMOVDQU -16(DI)(BX*1), X1
    VPCMPEQB X0, X1, X2
    VPMOVMSKB X2, DX
    CMPW DX, $0xFFFF
    JNE not_equal_avx
    JMP done_equal

not_equal_avx:
    VZEROUPPER
not_equal:
    MOVB $0, ret+48(FP)
    RET

done_equal:
    VZEROUPPER
    MOVB $1, ret+48(FP)
    RET
