// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The hand-written Y functions, see abi0_callbacks_linux.go. In the ABI0
// frame of these functions 0(RSP) is the saved LR, so the arguments of the
// calls go from 8(RSP) up.

#include "funcdata.h"
#include "textflag.h"

// func Yatexit(tls *TLS, func_ uintptr) (r int32)
TEXT ·Yatexit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD func_+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·yatexit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yat_quick_exit(tls *TLS, __ccgo_fp_func uintptr) (r1 int32)
TEXT ·Yat_quick_exit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD __ccgo_fp_func+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·yat_quick_exit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
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

// func Ysigaction(tls *TLS, sig int32, sa uintptr, old uintptr) (r int32)
TEXT ·Ysigaction(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sig+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD sa+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD old+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·ysigaction(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ypthread_create(tls *TLS, res, attrp, entry, arg uintptr) (_2 int32)
TEXT ·Ypthread_create(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD res+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD attrp+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD entry+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD arg+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·ypthread_create(SB)
	MOVW 48(RSP), R0
	MOVW R0, _2+40(FP)
	RET

// func Ypthread_key_create(tls *TLS, k uintptr, dtor uintptr) (_3 int32)
TEXT ·Ypthread_key_create(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD k+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD dtor+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·ypthread_key_create(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y_pthread_cleanup_push(tls *TLS, cb, f, x uintptr)
TEXT ·Y_pthread_cleanup_push(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD cb+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD x+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·y_pthread_cleanup_push(SB)
	RET

// func Yfopencookie(tls *TLS, cookie uintptr, mode uintptr, iofuncs Tcookie_io_functions_t) (r uintptr)
TEXT ·Yfopencookie(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD cookie+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD iofuncs_Fread+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD iofuncs_Fwrite+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD iofuncs_Fseek+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD iofuncs_Fclose1+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·yfopencookie(SB)
	MOVD 64(RSP), R0
	MOVD R0, r+56(FP)
	RET

// func Ypthread_atfork(tls *TLS, prepare, parent, child uintptr) (_2 int32)
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
	MOVW R0, _2+32(FP)
	RET
