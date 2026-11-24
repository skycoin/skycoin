// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"

// The arm64 version of tls_linux_amd64.s, see there. In the ABI0 frame of
// these functions 0(RSP) is the saved LR, so the arguments of the calls go
// from 8(RSP) up.

// func TLSAlloc(p0 *TLS, p1 int) uintptr
TEXT ·TLSAlloc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD p0+0(FP), R0
	MOVD TLS_sp(R0), R1
	MOVD (TLS_stack+8)(R0), R2	// len(tls.stack)
	CMP R2, R1
	BHS slow			// sp >= len(tls.stack), unsigned
	MOVD TLS_stack(R0), R3
	ADD R1<<4, R3, R3		// &tls.stack[sp], sp*unsafe.Sizeof(tlsStackSlot{})
	MOVD p1+8(FP), R4
	MOVD tlsStackSlot_sz(R3), R5
	CMP R4, R5
	BLO slow			// tls.stack[sp].sz < Tsize_t(p1)
	ADD $1, R1
	MOVD R1, TLS_sp(R0)
	MOVD tlsStackSlot_p(R3), R0
	MOVD R0, ret+16(FP)
	RET
slow:
	MOVD R0, 8(RSP)
	MOVD p1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·tlsAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, ret+16(FP)
	RET

// func TLSFree(p0 *TLS, p1 int)
TEXT ·TLSFree(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD p0+0(FP), R0
	MOVBU TLS_checkSignals(R0), R1
	CBNZ R1, slow
	MOVD TLS_sp(R0), R1
	SUB $1, R1
	MOVD R1, TLS_sp(R0)
	RET
slow:
	MOVD R0, 8(RSP)
	MOVD p1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·tlsFree(SB)
	RET

TEXT ·TLSAllocaEntry(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD p0+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·tlsAllocaEntry(SB)
	RET

TEXT ·TLSAllocaExit(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD p0+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·tlsAllocaExit(SB)
	RET

// ABI0 code calls Xalloca(tls *TLS, n Tsize_t) uintptr for the allocs of qbe's
// goabi0 output. Go gives a function an ABI0 wrapper only if assembly in its
// own package refers to it. On linux/amd64 the generated Y functions do, on
// arm64 this does. Nothing calls abi0Xalloca.
TEXT abi0Xalloca<>(SB),NOSPLIT|NOFRAME,$0-0
	JMP ·Xalloca(SB)
