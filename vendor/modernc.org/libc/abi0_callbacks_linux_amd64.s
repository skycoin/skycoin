// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The hand-written Y functions and the calls of ABI0 callbacks, see
// abi0_callbacks_linux_amd64.go.

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

// func abi0CallPPS_S(tls *TLS, a, b uintptr, c Tsize_t, entry uintptr) (r Tssize_t)
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

// func Yatexit(tls *TLS, func_ uintptr) (r int32)
TEXT ·Yatexit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ func_+8(FP), AX
	MOVQ AX, 8(SP)
	CALL ·yatexit(SB)
	MOVL 16(SP), AX
	MOVL AX, r+16(FP)
	RET

// func Yat_quick_exit(tls *TLS, __ccgo_fp_func uintptr) (r1 int32)
TEXT ·Yat_quick_exit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ __ccgo_fp_func+8(FP), AX
	MOVQ AX, 8(SP)
	CALL ·yat_quick_exit(SB)
	MOVL 16(SP), AX
	MOVL AX, r1+16(FP)
	RET

// func Ysignal(tls *TLS, signum int32, handler uintptr) (r uintptr)
TEXT ·Ysignal(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVL signum+8(FP), AX
	MOVL AX, 8(SP)
	MOVQ handler+16(FP), AX
	MOVQ AX, 16(SP)
	CALL ·ysignal(SB)
	MOVQ 24(SP), AX
	MOVQ AX, r+24(FP)
	RET

// func Ysigaction(tls *TLS, sig int32, sa uintptr, old uintptr) (r int32)
TEXT ·Ysigaction(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVL sig+8(FP), AX
	MOVL AX, 8(SP)
	MOVQ sa+16(FP), AX
	MOVQ AX, 16(SP)
	MOVQ old+24(FP), AX
	MOVQ AX, 24(SP)
	CALL ·ysigaction(SB)
	MOVL 32(SP), AX
	MOVL AX, r+32(FP)
	RET

// func Ypthread_create(tls *TLS, res, attrp, entry, arg uintptr) (_2 int32)
TEXT ·Ypthread_create(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ res+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ attrp+16(FP), AX
	MOVQ AX, 16(SP)
	MOVQ entry+24(FP), AX
	MOVQ AX, 24(SP)
	MOVQ arg+32(FP), AX
	MOVQ AX, 32(SP)
	CALL ·ypthread_create(SB)
	MOVL 40(SP), AX
	MOVL AX, _2+40(FP)
	RET

// func Ypthread_key_create(tls *TLS, k uintptr, dtor uintptr) (_3 int32)
TEXT ·Ypthread_key_create(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ k+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ dtor+16(FP), AX
	MOVQ AX, 16(SP)
	CALL ·ypthread_key_create(SB)
	MOVL 24(SP), AX
	MOVL AX, _3+24(FP)
	RET

// func Y_pthread_cleanup_push(tls *TLS, cb, f, x uintptr)
TEXT ·Y_pthread_cleanup_push(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ cb+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ f+16(FP), AX
	MOVQ AX, 16(SP)
	MOVQ x+24(FP), AX
	MOVQ AX, 24(SP)
	CALL ·y_pthread_cleanup_push(SB)
	RET

// func Yfopencookie(tls *TLS, cookie uintptr, mode uintptr, iofuncs Tcookie_io_functions_t) (r uintptr)
TEXT ·Yfopencookie(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ cookie+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ mode+16(FP), AX
	MOVQ AX, 16(SP)
	MOVQ iofuncs_Fread+24(FP), AX
	MOVQ AX, 24(SP)
	MOVQ iofuncs_Fwrite+32(FP), AX
	MOVQ AX, 32(SP)
	MOVQ iofuncs_Fseek+40(FP), AX
	MOVQ AX, 40(SP)
	MOVQ iofuncs_Fclose1+48(FP), AX
	MOVQ AX, 48(SP)
	CALL ·yfopencookie(SB)
	MOVQ 56(SP), AX
	MOVQ AX, r+56(FP)
	RET

// func Ypthread_atfork(tls *TLS, prepare, parent, child uintptr) (_2 int32)
TEXT ·Ypthread_atfork(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVQ tls+0(FP), AX
	MOVQ AX, 0(SP)
	MOVQ prepare+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ parent+16(FP), AX
	MOVQ AX, 16(SP)
	MOVQ child+24(FP), AX
	MOVQ AX, 24(SP)
	CALL ·ypthread_atfork(SB)
	MOVL 32(SP), AX
	MOVL AX, _2+32(FP)
	RET
