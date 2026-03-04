// +build amd64
// +build !noasm

#include "textflag.h"

// func spinWait(count int)
TEXT ·spinWait(SB), NOSPLIT, $0-8
    MOVQ count+0(FP), CX
spin:
    PAUSE
    LOOP spin
    RET
