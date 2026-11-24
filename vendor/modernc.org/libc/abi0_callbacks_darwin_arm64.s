// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The hand-written Y functions, see abi0_callbacks_darwin.go. In the ABI0
// frame of these functions 0(RSP) is the saved LR, so the arguments of the
// calls go from 8(RSP) up.

#include "funcdata.h"
#include "textflag.h"

// func Yatexit(tls *TLS, function uintptr) (r int32)
TEXT ·Yatexit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD function+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·yatexit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysignal(tls *TLS, signum int32, handler uintptr) (r uintptr)
TEXT ·Ysignal(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW signum+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD handler+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·ysignal(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ysigaction(tls *TLS, signum int32, act, oldact uintptr) (r int32)
TEXT ·Ysigaction(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW signum+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD act+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD oldact+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·ysigaction(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ypthread_create(tls *TLS, pThread, pAttr, startRoutine, arg uintptr) (r int32)
TEXT ·Ypthread_create(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pThread+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pAttr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD startRoutine+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD arg+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·ypthread_create(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ypthread_key_create(tls *TLS, pKey, destructor uintptr) (r int32)
TEXT ·Ypthread_key_create(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pKey+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD destructor+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·ypthread_key_create(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ypthread_atfork(tls *TLS, prepare, parent, child uintptr) (r int32)
TEXT ·Ypthread_atfork(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD prepare+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD parent+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD child+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·ypthread_atfork(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET
