#include "textflag.h"

// func Match16(a, b []byte) bool
TEXT ·Match16(SB), NOSPLIT, $0-49
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI

    // Load 16 bytes into XMM registers
    MOVOU (SI), X0
    MOVOU (DI), X1

    // Compare
    PCMPEQB X1, X0
    PMOVMSKB X0, AX

    // If all bits set (0xFFFF), they are equal
    CMPW AX, $0xFFFF
    SETEQ ret+48(FP)
    RET

// func Match32(a, b []byte) bool
TEXT ·Match32(SB), NOSPLIT, $0-49
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI

    // Load 32 bytes into YMM registers (AVX2)
    VMOVDQU (SI), Y0
    VMOVDQU (DI), Y1

    // Compare
    VPCMPEQB Y1, Y0, Y2
    VPMOVMSKB Y2, AX

    // If all bits set (0xFFFFFFFF), they are equal
    CMPL AX, $0xFFFFFFFF
    SETEQ ret+48(FP)
    VZEROUPPER
    RET

// func SearchGroup16(control *uint8, tag uint8) uint16
TEXT ·SearchGroup16(SB), NOSPLIT, $0-18
    MOVQ control+0(FP), SI
    MOVBLZX tag+8(FP), AX

    // Broadcast tag to all bytes in X0
    MOVD AX, X0
    VPBROADCASTB X0, X0

    // Load 16 bytes of control
    VMOVDQU (SI), X1

    // Compare equal
    VPCMPEQB X0, X1, X1

    // Move mask
    VPMOVMSKB X1, AX

    MOVW AX, ret+16(FP)
    RET
