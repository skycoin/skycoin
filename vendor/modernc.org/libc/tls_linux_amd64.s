#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"

// TLSAlloc and TLSFree inline the fast paths of (*TLS).Alloc and (*TLS).Free
// and fall back to calling them. Both must stay splittable: callers link
// against them assuming the usual split prologue, which bounds the stack they
// need to the 8 bytes of the return address. The slow paths use CALL, not a
// tail JMP, because the assembler would mark a $0 frame function without a
// CALL NOSPLIT. TestTLSStubsSplittable checks it.

// func TLSAlloc(p0 *TLS, p1 int) uintptr
TEXT ·TLSAlloc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ p0+0(FP), AX
	MOVQ TLS_sp(AX), CX
	CMPQ CX, (TLS_stack+8)(AX)		// len(tls.stack)
	JCC slow				// sp >= len(tls.stack), unsigned
	MOVQ CX, DX
	SHLQ $4, DX				// sp*unsafe.Sizeof(tlsStackSlot{})
	ADDQ TLS_stack(AX), DX			// &tls.stack[sp]
	MOVQ p1+8(FP), BX
	CMPQ tlsStackSlot_sz(DX), BX
	JCS slow				// tls.stack[sp].sz < Tsize_t(p1)
	INCQ CX
	MOVQ CX, TLS_sp(AX)
	MOVQ tlsStackSlot_p(DX), AX
	MOVQ AX, ret+16(FP)
	RET
slow:
	MOVQ AX, 0(SP)
	MOVQ p1+8(FP), AX
	MOVQ AX, 8(SP)
	CALL ·tlsAlloc(SB)
	MOVQ 16(SP), AX
	MOVQ AX, ret+16(FP)
	RET

// func TLSFree(p0 *TLS, p1 int)
TEXT ·TLSFree(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ p0+0(FP), AX
	CMPB TLS_checkSignals(AX), $0
	JNE slow
	DECQ TLS_sp(AX)
	RET
slow:
	MOVQ AX, 0(SP)
	MOVQ p1+8(FP), AX
	MOVQ AX, 8(SP)
	CALL ·tlsFree(SB)
	RET

TEXT ·TLSAllocaEntry(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ p0+0(FP), AX
	MOVQ AX, 0(SP)
	CALL ·tlsAllocaEntry(SB)
	RET

TEXT ·TLSAllocaExit(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ p0+0(FP), AX
	MOVQ AX, 0(SP)
	CALL ·tlsAllocaExit(SB)
	RET
