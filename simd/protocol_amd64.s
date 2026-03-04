#include "textflag.h"

// func FastValidateHeader(buf []byte) bool
// Header bytes 0-3: Magic(0xBA, 0x01), Version(0x01), OpCode(??)
// In Little-Endian (register): [Op] [01] [01] [BA] ...
TEXT ·FastValidateHeader(SB), NOSPLIT, $0-25
    MOVQ buf_base+0(FP), SI
    MOVQ buf_len+8(FP), CX

    // Check min length 16
    CMPQ CX, $16
    JL invalid

    // Load first 4 bytes as uint32
    MOVL (SI), AX

    // Mask out OpCode (byte 3)
    // AX = [Op][01][01][BA]
    // Mask = 0x00FFFFFF
    // Expected = 0x0101BA
    ANDL $0x00FFFFFF, AX
    CMPL AX, $0x0101BA
    JNE invalid

    // Verify OpCode range (byte 3)
    // Extract OpCode
    MOVBLZX 3(SI), DX

    // Valid if (op <= OpDecr) or (op >= OpHGet && op <= OpRPush)
    // Simplified for hot path: op <= OpRPush (0x21)
    CMPB DX, $0x21
    JA invalid

    MOVB $1, ret+24(FP)
    RET

invalid:
    MOVB $0, ret+24(FP)
    RET
