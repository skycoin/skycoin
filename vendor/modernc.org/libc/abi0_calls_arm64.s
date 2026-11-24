// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || darwin

// The calls of ABI0 functions, see abi0_callbacks.go. In the ABI0 frame of
// these functions 0(RSP) is the saved LR, so the arguments of the calls go
// from 8(RSP) up.

#include "funcdata.h"
#include "textflag.h"

// func abi0CallV(tls *TLS, entry uintptr)
TEXT ·abi0CallV(SB),$8-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD entry+8(FP), R0
	CALL (R0)
	RET

// func abi0CallI(tls *TLS, a int32, entry uintptr)
TEXT ·abi0CallI(SB),$16-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW a+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD entry+16(FP), R0
	CALL (R0)
	RET

// func abi0CallIPP(tls *TLS, a int32, b, c, entry uintptr)
TEXT ·abi0CallIPP(SB),$32-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW a+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD c+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD entry+32(FP), R0
	CALL (R0)
	RET

// func abi0CallP(tls *TLS, a, entry uintptr)
TEXT ·abi0CallP(SB),$16-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD entry+16(FP), R0
	CALL (R0)
	RET

// func abi0CallP_P(tls *TLS, a, entry uintptr) (r uintptr)
TEXT ·abi0CallP_P(SB),$24-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD entry+16(FP), R0
	CALL (R0)
	MOVD 24(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func abi0CallP_I(tls *TLS, a, entry uintptr) (r int32)
TEXT ·abi0CallP_I(SB),$24-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD entry+16(FP), R0
	CALL (R0)
	MOVW 24(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func abi0CallPPS_S(tls *TLS, a, b uintptr, c uint64, entry uintptr) (r int64)
TEXT ·abi0CallPPS_S(SB),$40-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD c+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD entry+32(FP), R0
	CALL (R0)
	MOVD 40(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func abi0CallPPI_I(tls *TLS, a, b uintptr, c int32, entry uintptr) (r int32)
TEXT ·abi0CallPPI_I(SB),$40-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW c+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD entry+32(FP), R0
	CALL (R0)
	MOVW 40(RSP), R0
	MOVW R0, r+40(FP)
	RET
