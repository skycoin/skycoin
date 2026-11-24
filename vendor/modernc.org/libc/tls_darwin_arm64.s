// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "funcdata.h"
#include "textflag.h"

// The ABI0 entry points of the TLS, see tls_abi0.go. Unlike on linux, they
// only call the Go functions. Like there, they must stay splittable, see
// tls_linux_amd64.s. In their ABI0 frame 0(RSP) is the saved LR, so the
// arguments of the calls go from 8(RSP) up.

// func TLSAlloc(p0 *TLS, p1 int) uintptr
TEXT ·TLSAlloc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD p0+0(FP), R0
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

// See tls_linux_arm64.s. Nothing calls abi0Xalloca.
TEXT abi0Xalloca<>(SB),NOSPLIT|NOFRAME,$0-0
	JMP ·Xalloca(SB)
