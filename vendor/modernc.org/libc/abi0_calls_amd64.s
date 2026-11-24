// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

// The calls of ABI0 functions, see abi0_callbacks.go.

#include "funcdata.h"
#include "textflag.h"

// func abi0CallV(tls *TLS, entry uintptr)
TEXT ·abi0CallV(SB),$8-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ entry+8(FP), AX
	CALL AX
	RET

// func abi0CallI(tls *TLS, a int32, entry uintptr)
TEXT ·abi0CallI(SB),$16-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVL a+8(FP), AX
	MOVL AX, 8(SP)
	MOVQ entry+16(FP), AX
	CALL AX
	RET

// func abi0CallIPP(tls *TLS, a int32, b, c, entry uintptr)
TEXT ·abi0CallIPP(SB),$32-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVL a+8(FP), AX
	MOVL AX, 8(SP)
	MOVQ b+16(FP), AX
	MOVQ AX, 16(SP)
	MOVQ c+24(FP), AX
	MOVQ AX, 24(SP)
	MOVQ entry+32(FP), AX
	CALL AX
	RET

// func abi0CallP(tls *TLS, a, entry uintptr)
TEXT ·abi0CallP(SB),$16-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ a+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ entry+16(FP), AX
	CALL AX
	RET

// func abi0CallP_P(tls *TLS, a, entry uintptr) (r uintptr)
TEXT ·abi0CallP_P(SB),$24-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ a+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ entry+16(FP), AX
	CALL AX
	MOVQ 16(SP), AX
	MOVQ AX, r+24(FP)
	RET

// func abi0CallP_I(tls *TLS, a, entry uintptr) (r int32)
TEXT ·abi0CallP_I(SB),$24-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ a+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ entry+16(FP), AX
	CALL AX
	MOVL 16(SP), AX
	MOVL AX, r+24(FP)
	RET

// func abi0CallPPS_S(tls *TLS, a, b uintptr, c uint64, entry uintptr) (r int64)
TEXT ·abi0CallPPS_S(SB),$40-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ a+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ b+16(FP), AX
	MOVQ AX, 16(SP)
	MOVQ c+24(FP), AX
	MOVQ AX, 24(SP)
	MOVQ entry+32(FP), AX
	CALL AX
	MOVQ 32(SP), AX
	MOVQ AX, r+40(FP)
	RET

// func abi0CallPPI_I(tls *TLS, a, b uintptr, c int32, entry uintptr) (r int32)
TEXT ·abi0CallPPI_I(SB),$40-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ a+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ b+16(FP), AX
	MOVQ AX, 16(SP)
	MOVL c+24(FP), AX
	MOVL AX, 24(SP)
	MOVQ entry+32(FP), AX
	CALL AX
	MOVL 32(SP), AX
	MOVL AX, r+40(FP)
	RET
