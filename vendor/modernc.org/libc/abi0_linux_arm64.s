// Code generated for linux/arm64 by 'qbecc --abi0wrap .', DO NOT EDIT.

#include "funcdata.h"
#include "textflag.h"

// func Y_Exit(tls *TLS, ec int32)
TEXT ·Y_Exit(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ec+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X_Exit(SB)
	RET

// func Y_IO_feof_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Y_IO_feof_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X_IO_feof_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y_IO_ferror_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Y_IO_ferror_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X_IO_ferror_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y_IO_getc(tls *TLS, f1 uintptr) (r int32)
TEXT ·Y_IO_getc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X_IO_getc(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y_IO_getc_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Y_IO_getc_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X_IO_getc_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y_IO_putc(tls *TLS, c1 int32, f1 uintptr) (r int32)
TEXT ·Y_IO_putc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c1+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X_IO_putc(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y_IO_putc_unlocked(tls *TLS, c int32, f uintptr) (r int32)
TEXT ·Y_IO_putc_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X_IO_putc_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y___errno_location(tls *TLS) (r uintptr)
TEXT ·Y___errno_location(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X___errno_location(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__aio_close(tls *TLS, fd int32) (_2 int32)
TEXT ·Y__aio_close(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__aio_close(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__asctime_r(tls *TLS, tm uintptr, buf uintptr) (r uintptr)
TEXT ·Y__asctime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tm+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__asctime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__assert_fail(tls *TLS, expr uintptr, file uintptr, line int32, func1 uintptr)
TEXT ·Y__assert_fail(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD expr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD file+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW line+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD func1+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__assert_fail(SB)
	RET

// func Y__atomic_compare_exchangeInt16(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)
TEXT ·Y__atomic_compare_exchangeInt16(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW weak+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW success+36(FP), R0
	MOVW R0, 44(RSP)
	MOVW failure+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__atomic_compare_exchangeInt16(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__atomic_compare_exchangeInt32(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)
TEXT ·Y__atomic_compare_exchangeInt32(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW weak+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW success+36(FP), R0
	MOVW R0, 44(RSP)
	MOVW failure+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__atomic_compare_exchangeInt32(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__atomic_compare_exchangeInt64(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)
TEXT ·Y__atomic_compare_exchangeInt64(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW weak+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW success+36(FP), R0
	MOVW R0, 44(RSP)
	MOVW failure+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__atomic_compare_exchangeInt64(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__atomic_compare_exchangeInt8(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)
TEXT ·Y__atomic_compare_exchangeInt8(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW weak+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW success+36(FP), R0
	MOVW R0, 44(RSP)
	MOVW failure+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__atomic_compare_exchangeInt8(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__atomic_compare_exchangeUint16(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)
TEXT ·Y__atomic_compare_exchangeUint16(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW weak+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW success+36(FP), R0
	MOVW R0, 44(RSP)
	MOVW failure+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__atomic_compare_exchangeUint16(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__atomic_compare_exchangeUint32(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)
TEXT ·Y__atomic_compare_exchangeUint32(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW weak+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW success+36(FP), R0
	MOVW R0, 44(RSP)
	MOVW failure+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__atomic_compare_exchangeUint32(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__atomic_compare_exchangeUint64(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)
TEXT ·Y__atomic_compare_exchangeUint64(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW weak+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW success+36(FP), R0
	MOVW R0, 44(RSP)
	MOVW failure+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__atomic_compare_exchangeUint64(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__atomic_compare_exchangeUint8(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)
TEXT ·Y__atomic_compare_exchangeUint8(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW weak+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW success+36(FP), R0
	MOVW R0, 44(RSP)
	MOVW failure+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__atomic_compare_exchangeUint8(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__atomic_exchangeInt16(t *TLS, ptr, val, ret uintptr, _ int32)
TEXT ·Y__atomic_exchangeInt16(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ret+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__atomic_exchangeInt16(SB)
	RET

// func Y__atomic_exchangeInt32(t *TLS, ptr, val, ret uintptr, _ int32)
TEXT ·Y__atomic_exchangeInt32(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ret+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__atomic_exchangeInt32(SB)
	RET

// func Y__atomic_exchangeInt64(t *TLS, ptr, val, ret uintptr, _ int32)
TEXT ·Y__atomic_exchangeInt64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ret+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__atomic_exchangeInt64(SB)
	RET

// func Y__atomic_exchangeInt8(t *TLS, ptr, val, ret uintptr, _ int32)
TEXT ·Y__atomic_exchangeInt8(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ret+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__atomic_exchangeInt8(SB)
	RET

// func Y__atomic_exchangeUint16(t *TLS, ptr, val, ret uintptr, _ int32)
TEXT ·Y__atomic_exchangeUint16(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ret+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__atomic_exchangeUint16(SB)
	RET

// func Y__atomic_exchangeUint32(t *TLS, ptr, val, ret uintptr, _ int32)
TEXT ·Y__atomic_exchangeUint32(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ret+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__atomic_exchangeUint32(SB)
	RET

// func Y__atomic_exchangeUint64(t *TLS, ptr, val, ret uintptr, _ int32)
TEXT ·Y__atomic_exchangeUint64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ret+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__atomic_exchangeUint64(SB)
	RET

// func Y__atomic_exchangeUint8(t *TLS, ptr, val, ret uintptr, _ int32)
TEXT ·Y__atomic_exchangeUint8(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ret+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__atomic_exchangeUint8(SB)
	RET

// func Y__atomic_fetch_addInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__atomic_fetch_addInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_addInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_addInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__atomic_fetch_addInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_addInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_addInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__atomic_fetch_addInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_addInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_addInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__atomic_fetch_addInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_addInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_addUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__atomic_fetch_addUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_addUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_addUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__atomic_fetch_addUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_addUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_addUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__atomic_fetch_addUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_addUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_addUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__atomic_fetch_addUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_addUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_andInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__atomic_fetch_andInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_andInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_andInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__atomic_fetch_andInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_andInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_andInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__atomic_fetch_andInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_andInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_andInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__atomic_fetch_andInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_andInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_andUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__atomic_fetch_andUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_andUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_andUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__atomic_fetch_andUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_andUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_andUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__atomic_fetch_andUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_andUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_andUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__atomic_fetch_andUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_andUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_orInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__atomic_fetch_orInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_orInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_orInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__atomic_fetch_orInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_orInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_orInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__atomic_fetch_orInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_orInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_orInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__atomic_fetch_orInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_orInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_orUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__atomic_fetch_orUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_orUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_orUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__atomic_fetch_orUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_orUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_orUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__atomic_fetch_orUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_orUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_orUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__atomic_fetch_orUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_orUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_subInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__atomic_fetch_subInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_subInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_subInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__atomic_fetch_subInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_subInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_subInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__atomic_fetch_subInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_subInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_subInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__atomic_fetch_subInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_subInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_subUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__atomic_fetch_subUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_subUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_subUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__atomic_fetch_subUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_subUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_subUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__atomic_fetch_subUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_subUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_subUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__atomic_fetch_subUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_subUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_xorInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__atomic_fetch_xorInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_xorInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_xorInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__atomic_fetch_xorInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_xorInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_xorInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__atomic_fetch_xorInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_xorInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_xorInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__atomic_fetch_xorInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_xorInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_fetch_xorUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__atomic_fetch_xorUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_xorUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__atomic_fetch_xorUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__atomic_fetch_xorUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_xorUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__atomic_fetch_xorUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__atomic_fetch_xorUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_fetch_xorUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__atomic_fetch_xorUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__atomic_fetch_xorUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__atomic_fetch_xorUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__atomic_loadFloat32(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadFloat32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadFloat32(SB)
	RET

// func Y__atomic_loadFloat64(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadFloat64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadFloat64(SB)
	RET

// func Y__atomic_loadInt16(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadInt16(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadInt16(SB)
	RET

// func Y__atomic_loadInt32(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadInt32(SB)
	RET

// func Y__atomic_loadInt64(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadInt64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadInt64(SB)
	RET

// func Y__atomic_loadInt8(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadInt8(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadInt8(SB)
	RET

// func Y__atomic_loadUint16(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadUint16(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadUint16(SB)
	RET

// func Y__atomic_loadUint32(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadUint32(SB)
	RET

// func Y__atomic_loadUint64(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadUint64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadUint64(SB)
	RET

// func Y__atomic_loadUint8(t *TLS, ptr, ret uintptr, memorder int32)
TEXT ·Y__atomic_loadUint8(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ret+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_loadUint8(SB)
	RET

// func Y__atomic_storeFloat32(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeFloat32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeFloat32(SB)
	RET

// func Y__atomic_storeFloat64(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeFloat64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeFloat64(SB)
	RET

// func Y__atomic_storeInt16(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeInt16(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeInt16(SB)
	RET

// func Y__atomic_storeInt32(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeInt32(SB)
	RET

// func Y__atomic_storeInt64(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeInt64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeInt64(SB)
	RET

// func Y__atomic_storeInt8(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeInt8(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeInt8(SB)
	RET

// func Y__atomic_storeUint16(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeUint16(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeUint16(SB)
	RET

// func Y__atomic_storeUint32(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeUint32(SB)
	RET

// func Y__atomic_storeUint64(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeUint64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeUint64(SB)
	RET

// func Y__atomic_storeUint8(t *TLS, ptr, val uintptr, memorder int32)
TEXT ·Y__atomic_storeUint8(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__atomic_storeUint8(SB)
	RET

// func Y__block_all_sigs(tls *TLS, set uintptr)
TEXT ·Y__block_all_sigs(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__block_all_sigs(SB)
	RET

// func Y__block_app_sigs(tls *TLS, set uintptr)
TEXT ·Y__block_app_sigs(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__block_app_sigs(SB)
	RET

// func Y__builtin___memcpy_chk(t *TLS, dest, src uintptr, n, os Tsize_t) (r uintptr)
TEXT ·Y__builtin___memcpy_chk(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD os+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__builtin___memcpy_chk(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Y__builtin___memmove_chk(t *TLS, dest, src uintptr, n, os Tsize_t) (_3 uintptr)
TEXT ·Y__builtin___memmove_chk(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD os+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__builtin___memmove_chk(SB)
	MOVD 48(RSP), R0
	MOVD R0, _3+40(FP)
	RET

// func Y__builtin___memset_chk(t *TLS, s uintptr, c int32, n, os Tsize_t) (_4 uintptr)
TEXT ·Y__builtin___memset_chk(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD os+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__builtin___memset_chk(SB)
	MOVD 48(RSP), R0
	MOVD R0, _4+40(FP)
	RET

// func Y__builtin___snprintf_chk(t *TLS, str uintptr, maxlen Tsize_t, flag int32, os Tsize_t, format, args uintptr) (r int32)
TEXT ·Y__builtin___snprintf_chk(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD maxlen+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flag+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD os+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD format+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD args+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·X__builtin___snprintf_chk(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Y__builtin___sprintf_chk(t *TLS, s uintptr, flag int32, os Tsize_t, format, args uintptr) (r int32)
TEXT ·Y__builtin___sprintf_chk(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flag+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD os+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD format+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD args+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__builtin___sprintf_chk(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Y__builtin___strcat_chk(t *TLS, dest, src uintptr, os Tsize_t) (r uintptr)
TEXT ·Y__builtin___strcat_chk(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD os+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin___strcat_chk(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__builtin___strcpy_chk(t *TLS, dest, src uintptr, os Tsize_t) (_3 uintptr)
TEXT ·Y__builtin___strcpy_chk(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD os+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin___strcpy_chk(SB)
	MOVD 40(RSP), R0
	MOVD R0, _3+32(FP)
	RET

// func Y__builtin___strncat_chk(t *TLS, dest, src uintptr, n, destlen Tsize_t) (r uintptr)
TEXT ·Y__builtin___strncat_chk(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD destlen+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__builtin___strncat_chk(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Y__builtin___strncpy_chk(t *TLS, dest, src uintptr, n, os Tsize_t) (r uintptr)
TEXT ·Y__builtin___strncpy_chk(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD os+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__builtin___strncpy_chk(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Y__builtin___vsnprintf_chk(t *TLS, str uintptr, maxlen Tsize_t, flag int32, os Tsize_t, format, args uintptr) (r int32)
TEXT ·Y__builtin___vsnprintf_chk(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD maxlen+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flag+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD os+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD format+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD args+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·X__builtin___vsnprintf_chk(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Y__builtin___vsprintf_chk(t *TLS, s uintptr, flag int32, os Tsize_t, format, args uintptr) (r int32)
TEXT ·Y__builtin___vsprintf_chk(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flag+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD os+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD format+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD args+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__builtin___vsprintf_chk(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Y__builtin_abort(t *TLS)
TEXT ·Y__builtin_abort(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_abort(SB)
	RET

// func Y__builtin_abs(t *TLS, j int32) (_2 int32)
TEXT ·Y__builtin_abs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW j+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_abs(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_acosh(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_acosh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_acosh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_add_overflowInt64(t *TLS, a, b int64, res uintptr) (_3 int32)
TEXT ·Y__builtin_add_overflowInt64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD res+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_add_overflowInt64(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Y__builtin_add_overflowUint32(t *TLS, a, b uint32, res uintptr) (_3 int32)
TEXT ·Y__builtin_add_overflowUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW a+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW b+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD res+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_add_overflowUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__builtin_add_overflowUint64(t *TLS, a, b uint64, res uintptr) (_3 int32)
TEXT ·Y__builtin_add_overflowUint64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD res+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_add_overflowUint64(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Y__builtin_alloca(tls *TLS, size Tsize_t) (_2 uintptr)
TEXT ·Y__builtin_alloca(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD size+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_alloca(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_asinh(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_asinh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_asinh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_asinhf(t *TLS, x float32) (_2 float32)
TEXT ·Y__builtin_asinhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_asinhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_atanh(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_atanh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_atanh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_bswap16(t *TLS, x uint16) (_2 uint16)
TEXT ·Y__builtin_bswap16(SB),$24-18
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVH x+8(FP), R0
	MOVH R0, 16(RSP)
	CALL ·X__builtin_bswap16(SB)
	MOVH 24(RSP), R0
	MOVH R0, _2+16(FP)
	RET

// func Y__builtin_bswap32(t *TLS, x uint32) (_2 uint32)
TEXT ·Y__builtin_bswap32(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_bswap32(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_bswap64(t *TLS, x uint64) (_2 uint64)
TEXT ·Y__builtin_bswap64(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_bswap64(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_bzero(t *TLS, s uintptr, n Tsize_t)
TEXT ·Y__builtin_bzero(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_bzero(SB)
	RET

// func Y__builtin_cbrt(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_cbrt(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_cbrt(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_clz(t *TLS, n uint32) (_2 int32)
TEXT ·Y__builtin_clz(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_clz(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_clzl(t *TLS, n ulong) (_2 int32)
TEXT ·Y__builtin_clzl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_clzl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_clzll(t *TLS, n uint64) (_2 int32)
TEXT ·Y__builtin_clzll(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_clzll(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_copysign(t *TLS, x, y float64) (_2 float64)
TEXT ·Y__builtin_copysign(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_copysign(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__builtin_copysignf(t *TLS, x, y float32) (_2 float32)
TEXT ·Y__builtin_copysignf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·X__builtin_copysignf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_copysignl(t *TLS, x, y float64) (_2 float64)
TEXT ·Y__builtin_copysignl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_copysignl(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__builtin_ctz(t *TLS, n uint32) (_2 int32)
TEXT ·Y__builtin_ctz(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_ctz(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_ctzl(tls *TLS, x ulong) (_2 int32)
TEXT ·Y__builtin_ctzl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_ctzl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_ctzll(tls *TLS, x uint64) (_2 int32)
TEXT ·Y__builtin_ctzll(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_ctzll(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_erf(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_erf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_erf(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_erfc(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_erfc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_erfc(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_exit(t *TLS, status int32)
TEXT ·Y__builtin_exit(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW status+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_exit(SB)
	RET

// func Y__builtin_exp2(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_exp2(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_exp2(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_expect(t *TLS, exp, c long) (_2 long)
TEXT ·Y__builtin_expect(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD exp+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD c+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_expect(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__builtin_expm1(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_expm1(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_expm1(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_fabs(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_fabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_fabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_fabsf(t *TLS, x float32) (_2 float32)
TEXT ·Y__builtin_fabsf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_fabsf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_fabsl(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_fabsl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_fabsl(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_fdim(t *TLS, x, y float64) (_2 float64)
TEXT ·Y__builtin_fdim(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_fdim(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__builtin_ffs(tls *TLS, i int32) (r int32)
TEXT ·Y__builtin_ffs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW i+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_ffs(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__builtin_flt_rounds(t *TLS) (_1 int32)
TEXT ·Y__builtin_flt_rounds(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_flt_rounds(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Y__builtin_fma(tls *TLS, x, y, z float64) (r float64)
TEXT ·Y__builtin_fma(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD z+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_fma(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__builtin_fmaf(t *TLS, x, y, z float32) (_2 float32)
TEXT ·Y__builtin_fmaf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW z+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__builtin_fmaf(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Y__builtin_fmax(tls *TLS, x float64, y float64) (r float64)
TEXT ·Y__builtin_fmax(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_fmax(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__builtin_fmin(tls *TLS, x float64, y float64) (r float64)
TEXT ·Y__builtin_fmin(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_fmin(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__builtin_free(t *TLS, ptr uintptr)
TEXT ·Y__builtin_free(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_free(SB)
	RET

// func Y__builtin_getentropy(t *TLS, buf uintptr, n Tsize_t) (_3 int32)
TEXT ·Y__builtin_getentropy(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_getentropy(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__builtin_huge_val(t *TLS) (_1 float64)
TEXT ·Y__builtin_huge_val(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_huge_val(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__builtin_huge_valf(t *TLS) (_1 float32)
TEXT ·Y__builtin_huge_valf(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_huge_valf(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Y__builtin_hypot(tls *TLS, x float64, y float64) (r float64)
TEXT ·Y__builtin_hypot(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_hypot(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__builtin_ilogb(tls *TLS, x float64) (_2 int32)
TEXT ·Y__builtin_ilogb(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_ilogb(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_ilogbf(tls *TLS, x float32) (_2 int32)
TEXT ·Y__builtin_ilogbf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_ilogbf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_ilogbl(tls *TLS, x float64) (_2 int32)
TEXT ·Y__builtin_ilogbl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_ilogbl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_inf(t *TLS) (_1 float64)
TEXT ·Y__builtin_inf(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_inf(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__builtin_inff(tls *TLS) (_1 float32)
TEXT ·Y__builtin_inff(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_inff(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Y__builtin_infl(t *TLS) (_1 float64)
TEXT ·Y__builtin_infl(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_infl(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__builtin_isblank(tls *TLS, c int32) (r int32)
TEXT ·Y__builtin_isblank(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_isblank(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__builtin_isnan(t *TLS, x float64) (_2 int32)
TEXT ·Y__builtin_isnan(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_isnan(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_isnanf(t *TLS, x float32) (_2 int32)
TEXT ·Y__builtin_isnanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_isnanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_isnanl(t *TLS, x float64) (_2 int32)
TEXT ·Y__builtin_isnanl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_isnanl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_isprint(tls *TLS, c int32) (r int32)
TEXT ·Y__builtin_isprint(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_isprint(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__builtin_isunordered(t *TLS, a, b float64) (_2 int32)
TEXT ·Y__builtin_isunordered(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_isunordered(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Y__builtin_lgamma(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_lgamma(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_lgamma(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_llabs(tls *TLS, a int64) (_2 int64)
TEXT ·Y__builtin_llabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_llabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_log1p(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_log1p(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_log1p(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_log1pf(t *TLS, x float32) (_2 float32)
TEXT ·Y__builtin_log1pf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_log1pf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_log2(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_log2(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_log2(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_logb(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_logb(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_logb(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_lrint(tls *TLS, x float64) (r long)
TEXT ·Y__builtin_lrint(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_lrint(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__builtin_lrintf(tls *TLS, x float32) (r long)
TEXT ·Y__builtin_lrintf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_lrintf(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__builtin_lround(tls *TLS, x float64) (r long)
TEXT ·Y__builtin_lround(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_lround(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__builtin_malloc(t *TLS, size Tsize_t) (_2 uintptr)
TEXT ·Y__builtin_malloc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD size+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_malloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_memcmp(t *TLS, s1, s2 uintptr, n Tsize_t) (_3 int32)
TEXT ·Y__builtin_memcmp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_memcmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Y__builtin_memcpy(t *TLS, dest, src uintptr, n Tsize_t) (r uintptr)
TEXT ·Y__builtin_memcpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_memcpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__builtin_memset(t *TLS, s uintptr, c int32, n Tsize_t) (_4 uintptr)
TEXT ·Y__builtin_memset(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_memset(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Y__builtin_mmap(t *TLS, addr uintptr, length Tsize_t, prot, flags, fd int32, offset Toff_t) (_5 uintptr)
TEXT ·Y__builtin_mmap(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD length+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW prot+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW fd+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD offset+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__builtin_mmap(SB)
	MOVD 56(RSP), R0
	MOVD R0, _5+48(FP)
	RET

// func Y__builtin_mul_overflowInt64(t *TLS, a, b int64, res uintptr) (_3 int32)
TEXT ·Y__builtin_mul_overflowInt64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD res+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_mul_overflowInt64(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Y__builtin_mul_overflowUint128(t *TLS, a, b Uint128, res uintptr) (_3 int32)
TEXT ·Y__builtin_mul_overflowUint128(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a_Lo+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a_Hi+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD b_Lo+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD b_Hi+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD res+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__builtin_mul_overflowUint128(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Y__builtin_mul_overflowUint64(t *TLS, a, b uint64, res uintptr) (_3 int32)
TEXT ·Y__builtin_mul_overflowUint64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD res+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_mul_overflowUint64(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Y__builtin_nan(t *TLS, s uintptr) (_2 float64)
TEXT ·Y__builtin_nan(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_nan(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_nanf(tls *TLS, s uintptr) (_2 float32)
TEXT ·Y__builtin_nanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_nanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_nanl(t *TLS, s uintptr) (_2 float64)
TEXT ·Y__builtin_nanl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_nanl(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_nextafter(t *TLS, x, y float64) (_2 float64)
TEXT ·Y__builtin_nextafter(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_nextafter(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__builtin_nextafterf(t *TLS, x, y float32) (_2 float32)
TEXT ·Y__builtin_nextafterf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·X__builtin_nextafterf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_nexttoward(t *TLS, x, y float64) (_2 float64)
TEXT ·Y__builtin_nexttoward(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_nexttoward(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__builtin_object_size(t *TLS, p uintptr, typ int32) (_3 Tsize_t)
TEXT ·Y__builtin_object_size(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW typ+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__builtin_object_size(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Y__builtin_popcount(t *TLS, x uint32) (_2 int32)
TEXT ·Y__builtin_popcount(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_popcount(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_popcountl(t *TLS, x ulong) (_2 int32)
TEXT ·Y__builtin_popcountl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_popcountl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_popcountll(t *TLS, x uint64) (_2 int32)
TEXT ·Y__builtin_popcountll(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_popcountll(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__builtin_prefetch(t *TLS, addr, args uintptr)
TEXT ·Y__builtin_prefetch(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD args+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_prefetch(SB)
	RET

// func Y__builtin_printf(tls *TLS, fmt uintptr, va uintptr) (r int32)
TEXT ·Y__builtin_printf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_printf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__builtin_remainder(t *TLS, x, y float64) (_2 float64)
TEXT ·Y__builtin_remainder(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_remainder(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__builtin_remquo(t *TLS, x, y float64, quo uintptr) (_3 float64)
TEXT ·Y__builtin_remquo(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD quo+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_remquo(SB)
	MOVD 40(RSP), R0
	MOVD R0, _3+32(FP)
	RET

// func Y__builtin_rint(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_rint(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_rint(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_rintf(tls *TLS, x float32) (r float32)
TEXT ·Y__builtin_rintf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_rintf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__builtin_round(tls *TLS, x float64) (r float64)
TEXT ·Y__builtin_round(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_round(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__builtin_roundf(tls *TLS, x float32) (r float32)
TEXT ·Y__builtin_roundf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_roundf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__builtin_scalbn(t *TLS, x float64, n int32) (_3 float64)
TEXT ·Y__builtin_scalbn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__builtin_scalbn(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Y__builtin_snprintf(t *TLS, str uintptr, size Tsize_t, format, args uintptr) (_4 int32)
TEXT ·Y__builtin_snprintf(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD format+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD args+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__builtin_snprintf(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__builtin_sprintf(t *TLS, str, format, args uintptr) (r int32)
TEXT ·Y__builtin_sprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD format+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD args+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_sprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__builtin_strchr(t *TLS, s uintptr, c int32) (_3 uintptr)
TEXT ·Y__builtin_strchr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__builtin_strchr(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Y__builtin_strcmp(t *TLS, s1, s2 uintptr) (_2 int32)
TEXT ·Y__builtin_strcmp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s2+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_strcmp(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Y__builtin_strcpy(t *TLS, dest, src uintptr) (_2 uintptr)
TEXT ·Y__builtin_strcpy(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_strcpy(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__builtin_strlen(t *TLS, s uintptr) (_2 Tsize_t)
TEXT ·Y__builtin_strlen(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_strlen(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_sub_overflowInt64(t *TLS, a, b int64, res uintptr) (_3 int32)
TEXT ·Y__builtin_sub_overflowInt64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD res+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_sub_overflowInt64(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Y__builtin_tgamma(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_tgamma(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_tgamma(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_trap(t *TLS)
TEXT ·Y__builtin_trap(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_trap(SB)
	RET

// func Y__builtin_trunc(tls *TLS, x float64) (r float64)
TEXT ·Y__builtin_trunc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_trunc(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__builtin_unreachable(t *TLS)
TEXT ·Y__builtin_unreachable(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__builtin_unreachable(SB)
	RET

// func Y__builtin_vsnprintf(t *TLS, str uintptr, size Tsize_t, format, va uintptr) (_4 int32)
TEXT ·Y__builtin_vsnprintf(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD format+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD va+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__builtin_vsnprintf(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__c11_atomic_compare_exchange_strongInt16(t *TLS, ptr, expected uintptr, desired int16, success, failure int32) (_4 int32)
TEXT ·Y__c11_atomic_compare_exchange_strongInt16(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVH desired+24(FP), R0
	MOVH R0, 32(RSP)
	MOVW success+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW failure+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__c11_atomic_compare_exchange_strongInt16(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__c11_atomic_compare_exchange_strongInt32(t *TLS, ptr, expected uintptr, desired, success, failure int32) (_3 int32)
TEXT ·Y__c11_atomic_compare_exchange_strongInt32(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW desired+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW success+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW failure+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__c11_atomic_compare_exchange_strongInt32(SB)
	MOVW 48(RSP), R0
	MOVW R0, _3+40(FP)
	RET

// func Y__c11_atomic_compare_exchange_strongInt64(t *TLS, ptr, expected uintptr, desired int64, success, failure int32) (_4 int32)
TEXT ·Y__c11_atomic_compare_exchange_strongInt64(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW success+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW failure+36(FP), R0
	MOVW R0, 44(RSP)
	CALL ·X__c11_atomic_compare_exchange_strongInt64(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__c11_atomic_compare_exchange_strongInt8(t *TLS, ptr, expected uintptr, desired int8, success, failure int32) (_4 int32)
TEXT ·Y__c11_atomic_compare_exchange_strongInt8(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVB desired+24(FP), R0
	MOVB R0, 32(RSP)
	MOVW success+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW failure+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__c11_atomic_compare_exchange_strongInt8(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__c11_atomic_compare_exchange_strongUint16(t *TLS, ptr, expected uintptr, desired uint16, success, failure int32) (_4 int32)
TEXT ·Y__c11_atomic_compare_exchange_strongUint16(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVH desired+24(FP), R0
	MOVH R0, 32(RSP)
	MOVW success+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW failure+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__c11_atomic_compare_exchange_strongUint16(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__c11_atomic_compare_exchange_strongUint32(t *TLS, ptr, expected uintptr, desired uint32, success, failure int32) (_4 int32)
TEXT ·Y__c11_atomic_compare_exchange_strongUint32(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW desired+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW success+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW failure+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__c11_atomic_compare_exchange_strongUint32(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__c11_atomic_compare_exchange_strongUint64(t *TLS, ptr, expected uintptr, desired uint64, success, failure int32) (_4 int32)
TEXT ·Y__c11_atomic_compare_exchange_strongUint64(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD desired+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW success+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW failure+36(FP), R0
	MOVW R0, 44(RSP)
	CALL ·X__c11_atomic_compare_exchange_strongUint64(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__c11_atomic_compare_exchange_strongUint8(t *TLS, ptr, expected uintptr, desired uint8, success, failure int32) (_4 int32)
TEXT ·Y__c11_atomic_compare_exchange_strongUint8(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD expected+16(FP), R0
	MOVD R0, 24(RSP)
	MOVB desired+24(FP), R0
	MOVB R0, 32(RSP)
	MOVW success+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW failure+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__c11_atomic_compare_exchange_strongUint8(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__c11_atomic_exchangeInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__c11_atomic_exchangeInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_exchangeInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_exchangeInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__c11_atomic_exchangeInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_exchangeInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_exchangeInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__c11_atomic_exchangeInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_exchangeInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_exchangeInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__c11_atomic_exchangeInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_exchangeInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_exchangeUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__c11_atomic_exchangeUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_exchangeUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_exchangeUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__c11_atomic_exchangeUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_exchangeUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_exchangeUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__c11_atomic_exchangeUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_exchangeUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_exchangeUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__c11_atomic_exchangeUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_exchangeUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_addInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__c11_atomic_fetch_addInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_addInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_addInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__c11_atomic_fetch_addInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_addInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_addInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__c11_atomic_fetch_addInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_addInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_addInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__c11_atomic_fetch_addInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_addInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_addUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__c11_atomic_fetch_addUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_addUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_addUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__c11_atomic_fetch_addUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_addUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_addUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__c11_atomic_fetch_addUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_addUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_addUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__c11_atomic_fetch_addUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_addUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_andInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__c11_atomic_fetch_andInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_andInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_andInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__c11_atomic_fetch_andInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_andInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_andInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__c11_atomic_fetch_andInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_andInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_andInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__c11_atomic_fetch_andInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_andInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_andUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__c11_atomic_fetch_andUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_andUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_andUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__c11_atomic_fetch_andUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_andUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_andUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__c11_atomic_fetch_andUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_andUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_andUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__c11_atomic_fetch_andUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_andUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_orInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__c11_atomic_fetch_orInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_orInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_orInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__c11_atomic_fetch_orInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_orInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_orInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__c11_atomic_fetch_orInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_orInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_orInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__c11_atomic_fetch_orInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_orInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_orUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__c11_atomic_fetch_orUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_orUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_orUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__c11_atomic_fetch_orUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_orUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_orUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__c11_atomic_fetch_orUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_orUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_orUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__c11_atomic_fetch_orUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_orUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_subInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__c11_atomic_fetch_subInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_subInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_subInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__c11_atomic_fetch_subInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_subInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_subInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__c11_atomic_fetch_subInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_subInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_subInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__c11_atomic_fetch_subInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_subInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_subUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__c11_atomic_fetch_subUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_subUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_subUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__c11_atomic_fetch_subUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_subUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_subUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__c11_atomic_fetch_subUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_subUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_subUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__c11_atomic_fetch_subUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_subUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_xorInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)
TEXT ·Y__c11_atomic_fetch_xorInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_xorInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_xorInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)
TEXT ·Y__c11_atomic_fetch_xorInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_xorInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_xorInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)
TEXT ·Y__c11_atomic_fetch_xorInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_xorInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_xorInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)
TEXT ·Y__c11_atomic_fetch_xorInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_xorInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_xorUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)
TEXT ·Y__c11_atomic_fetch_xorUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_xorUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_xorUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)
TEXT ·Y__c11_atomic_fetch_xorUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_xorUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_fetch_xorUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)
TEXT ·Y__c11_atomic_fetch_xorUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_fetch_xorUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__c11_atomic_fetch_xorUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)
TEXT ·Y__c11_atomic_fetch_xorUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW _+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_fetch_xorUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_loadInt16(t *TLS, ptr uintptr, memorder int32) (r int16)
TEXT ·Y__c11_atomic_loadInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW memorder+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__c11_atomic_loadInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_loadInt32(t *TLS, ptr uintptr, memorder int32) (r int32)
TEXT ·Y__c11_atomic_loadInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW memorder+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__c11_atomic_loadInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_loadInt64(t *TLS, ptr uintptr, memorder int32) (r int64)
TEXT ·Y__c11_atomic_loadInt64(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW memorder+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__c11_atomic_loadInt64(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__c11_atomic_loadInt8(t *TLS, ptr uintptr, memorder int32) (r int8)
TEXT ·Y__c11_atomic_loadInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW memorder+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__c11_atomic_loadInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_loadUint16(t *TLS, ptr uintptr, memorder int32) (r uint16)
TEXT ·Y__c11_atomic_loadUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW memorder+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__c11_atomic_loadUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__c11_atomic_loadUint32(t *TLS, ptr uintptr, memorder int32) (r uint32)
TEXT ·Y__c11_atomic_loadUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW memorder+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__c11_atomic_loadUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__c11_atomic_loadUint64(t *TLS, ptr uintptr, memorder int32) (r uint64)
TEXT ·Y__c11_atomic_loadUint64(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW memorder+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__c11_atomic_loadUint64(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__c11_atomic_loadUint8(t *TLS, ptr uintptr, memorder int32) (r uint8)
TEXT ·Y__c11_atomic_loadUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW memorder+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__c11_atomic_loadUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__c11_atomic_storeInt16(t *TLS, ptr uintptr, val int16, memorder int32)
TEXT ·Y__c11_atomic_storeInt16(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW memorder+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_storeInt16(SB)
	RET

// func Y__c11_atomic_storeInt32(t *TLS, ptr uintptr, val int32, memorder int32)
TEXT ·Y__c11_atomic_storeInt32(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW memorder+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_storeInt32(SB)
	RET

// func Y__c11_atomic_storeInt64(t *TLS, ptr uintptr, val int64, memorder int32)
TEXT ·Y__c11_atomic_storeInt64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_storeInt64(SB)
	RET

// func Y__c11_atomic_storeInt8(t *TLS, ptr uintptr, val int8, memorder int32)
TEXT ·Y__c11_atomic_storeInt8(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW memorder+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_storeInt8(SB)
	RET

// func Y__c11_atomic_storeUint16(t *TLS, ptr uintptr, val uint16, memorder int32)
TEXT ·Y__c11_atomic_storeUint16(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH val+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW memorder+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_storeUint16(SB)
	RET

// func Y__c11_atomic_storeUint32(t *TLS, ptr uintptr, val uint32, memorder int32)
TEXT ·Y__c11_atomic_storeUint32(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW memorder+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_storeUint32(SB)
	RET

// func Y__c11_atomic_storeUint64(t *TLS, ptr uintptr, val uint64, memorder int32)
TEXT ·Y__c11_atomic_storeUint64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW memorder+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__c11_atomic_storeUint64(SB)
	RET

// func Y__c11_atomic_storeUint8(t *TLS, ptr uintptr, val uint8, memorder int32)
TEXT ·Y__c11_atomic_storeUint8(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB val+16(FP), R0
	MOVB R0, 24(RSP)
	MOVW memorder+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__c11_atomic_storeUint8(SB)
	RET

// func Y__ccgo_dmesg(t *TLS, fmt uintptr, va uintptr)
TEXT ·Y__ccgo_dmesg(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__ccgo_dmesg(SB)
	RET

// func Y__ccgo_getMutexType(tls *TLS, m uintptr) (_2 int32)
TEXT ·Y__ccgo_getMutexType(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__ccgo_getMutexType(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__ccgo_in6addr_anyp(t *TLS) (_1 uintptr)
TEXT ·Y__ccgo_in6addr_anyp(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__ccgo_in6addr_anyp(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__ccgo_pthreadAttrGetDetachState(tls *TLS, a uintptr) (_2 int32)
TEXT ·Y__ccgo_pthreadAttrGetDetachState(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__ccgo_pthreadAttrGetDetachState(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__ccgo_pthreadMutexattrGettype(tls *TLS, a uintptr) (_2 int32)
TEXT ·Y__ccgo_pthreadMutexattrGettype(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__ccgo_pthreadMutexattrGettype(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__ccgo_sqlite3_log(t *TLS, iErrCode int32, zFormat uintptr, args uintptr)
TEXT ·Y__ccgo_sqlite3_log(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW iErrCode+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD zFormat+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD args+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__ccgo_sqlite3_log(SB)
	RET

// func Y__clock_gettime(tls *TLS, clk Tclockid_t, ts uintptr) (r1 int32)
TEXT ·Y__clock_gettime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clk+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD ts+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__clock_gettime(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Y__clock_nanosleep(tls *TLS, clk Tclockid_t, flags int32, req uintptr, rem uintptr) (r int32)
TEXT ·Y__clock_nanosleep(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clk+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW flags+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD req+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD rem+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__clock_nanosleep(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__cmsg_nxthdr(t *TLS, msgh, cmsg uintptr) (_2 uintptr)
TEXT ·Y__cmsg_nxthdr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msgh+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD cmsg+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__cmsg_nxthdr(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__convert_scm_timestamps(tls *TLS, msg uintptr, csize Tsocklen_t)
TEXT ·Y__convert_scm_timestamps(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW csize+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__convert_scm_timestamps(SB)
	RET

// func Y__cos(tls *TLS, x float64, y float64) (r1 float64)
TEXT ·Y__cos(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__cos(SB)
	MOVD 32(RSP), R0
	MOVD R0, r1+24(FP)
	RET

// func Y__cosdf(tls *TLS, x float64) (r1 float32)
TEXT ·Y__cosdf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__cosdf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Y__crypt_blowfish(tls *TLS, key uintptr, setting uintptr, output uintptr) (r uintptr)
TEXT ·Y__crypt_blowfish(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD setting+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD output+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__crypt_blowfish(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__crypt_des(tls *TLS, key uintptr, setting uintptr, output uintptr) (r uintptr)
TEXT ·Y__crypt_des(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD setting+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD output+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__crypt_des(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__crypt_md5(tls *TLS, key uintptr, setting uintptr, output uintptr) (r uintptr)
TEXT ·Y__crypt_md5(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD setting+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD output+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__crypt_md5(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__crypt_r(tls *TLS, key uintptr, salt uintptr, data uintptr) (r uintptr)
TEXT ·Y__crypt_r(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD salt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD data+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__crypt_r(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__crypt_sha256(tls *TLS, key uintptr, setting uintptr, output uintptr) (r uintptr)
TEXT ·Y__crypt_sha256(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD setting+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD output+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__crypt_sha256(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__crypt_sha512(tls *TLS, key uintptr, setting uintptr, output uintptr) (r uintptr)
TEXT ·Y__crypt_sha512(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD setting+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD output+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__crypt_sha512(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__ctype_b_loc(tls *TLS) (r uintptr)
TEXT ·Y__ctype_b_loc(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__ctype_b_loc(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__ctype_get_mb_cur_max(tls *TLS) (r Tsize_t)
TEXT ·Y__ctype_get_mb_cur_max(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__ctype_get_mb_cur_max(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__ctype_tolower_loc(tls *TLS) (r uintptr)
TEXT ·Y__ctype_tolower_loc(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__ctype_tolower_loc(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__ctype_toupper_loc(tls *TLS) (r uintptr)
TEXT ·Y__ctype_toupper_loc(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__ctype_toupper_loc(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__des_setkey(tls *TLS, key uintptr, ekey uintptr)
TEXT ·Y__des_setkey(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ekey+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__des_setkey(SB)
	RET

// func Y__dn_expand(tls *TLS, base uintptr, end uintptr, src uintptr, dest uintptr, space int32) (r int32)
TEXT ·Y__dn_expand(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD base+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD end+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD src+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD dest+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW space+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·X__dn_expand(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Y__dns_parse(tls *TLS, r uintptr, rlen int32, __ccgo_fp_callback uintptr, ctx uintptr) (r1 int32)
TEXT ·Y__dns_parse(SB),$56-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 56(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 56(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_callback+24(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal___dns_parse_2(SB)	// Create the closure for calling __ccgo_fp_callback
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD r+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW rlen+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD 56(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 32(RSP)
	MOVD ctx+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__dns_parse(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0___dns_parse_2(SB),$72-68
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW _2+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD _3+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _4+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD _5+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW _6+48(FP), R0
	MOVW R0, 56(RSP)
	MOVD __ccgo_fp+56(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 64(RSP), R0
	MOVW R0, _7+64(FP)
	RET

// func Y__do_des(tls *TLS, l_in Tuint32_t, r_in Tuint32_t, l_out uintptr, r_out uintptr, count Tuint32_t, saltbits Tuint32_t, ekey uintptr)
TEXT ·Y__do_des(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW l_in+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW r_in+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD l_out+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD r_out+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW count+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW saltbits+36(FP), R0
	MOVW R0, 44(RSP)
	MOVD ekey+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__do_des(SB)
	RET

// func Y__do_orphaned_stdio_locks(tls *TLS)
TEXT ·Y__do_orphaned_stdio_locks(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__do_orphaned_stdio_locks(SB)
	RET

// func Y__dup3(tls *TLS, old int32, new1 int32, flags int32) (r1 int32)
TEXT ·Y__dup3(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW old+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW new1+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__dup3(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Y__duplocale(tls *TLS, old Tlocale_t) (r Tlocale_t)
TEXT ·Y__duplocale(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD old+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__duplocale(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__env_rm_add(tls *TLS, old uintptr, new1 uintptr)
TEXT ·Y__env_rm_add(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD old+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD new1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__env_rm_add(SB)
	RET

// func Y__errno_location(tls *TLS) (r uintptr)
TEXT ·Y__errno_location(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__errno_location(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__execvpe(tls *TLS, file uintptr, argv uintptr, envp uintptr) (r int32)
TEXT ·Y__execvpe(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD file+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD envp+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__execvpe(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__expo2(tls *TLS, x float64, sign float64) (r float64)
TEXT ·Y__expo2(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sign+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__expo2(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__expo2f(tls *TLS, x float32, sign float32) (r float32)
TEXT ·Y__expo2f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW sign+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·X__expo2f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fbufsize(tls *TLS, f uintptr) (r Tsize_t)
TEXT ·Y__fbufsize(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fbufsize(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__fclose_ca(tls *TLS, f uintptr) (r int32)
TEXT ·Y__fclose_ca(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fclose_ca(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fdopen(tls *TLS, fd int32, mode uintptr) (r uintptr)
TEXT ·Y__fdopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__fdopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__fesetround(tls *TLS, r int32) (r1 int32)
TEXT ·Y__fesetround(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW r+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__fesetround(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Y__fgetwc_unlocked(tls *TLS, f uintptr) (r Twint_t)
TEXT ·Y__fgetwc_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fgetwc_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__flbf(tls *TLS, f uintptr) (r int32)
TEXT ·Y__flbf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__flbf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__floatscan(tls *TLS, f uintptr, prec int32, pok int32) (r float64)
TEXT ·Y__floatscan(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW prec+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW pok+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__floatscan(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__flt_rounds(t *TLS) (_1 int32)
TEXT ·Y__flt_rounds(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__flt_rounds(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Y__fmodeflags(tls *TLS, mode uintptr) (r int32)
TEXT ·Y__fmodeflags(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD mode+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fmodeflags(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fopen_rb_ca(tls *TLS, filename uintptr, f uintptr, buf uintptr, len1 Tsize_t) (r uintptr)
TEXT ·Y__fopen_rb_ca(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD len1+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__fopen_rb_ca(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Y__fpclassify(tls *TLS, x float64) (r int32)
TEXT ·Y__fpclassify(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fpclassify(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fpclassifyf(tls *TLS, x float32) (r int32)
TEXT ·Y__fpclassifyf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__fpclassifyf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fpclassifyl(tls *TLS, x float64) (r int32)
TEXT ·Y__fpclassifyl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fpclassifyl(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fpending(tls *TLS, f uintptr) (r Tsize_t)
TEXT ·Y__fpending(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fpending(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__fpurge(tls *TLS, f uintptr) (r int32)
TEXT ·Y__fpurge(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fpurge(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fputwc_unlocked(tls *TLS, c Twchar_t, f uintptr) (r Twint_t)
TEXT ·Y__fputwc_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__fputwc_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__freadable(tls *TLS, f uintptr) (r int32)
TEXT ·Y__freadable(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__freadable(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__freadahead(tls *TLS, f uintptr) (r Tsize_t)
TEXT ·Y__freadahead(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__freadahead(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__freading(tls *TLS, f uintptr) (r int32)
TEXT ·Y__freading(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__freading(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__freadptr(tls *TLS, f uintptr, sizep uintptr) (r uintptr)
TEXT ·Y__freadptr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sizep+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__freadptr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__freadptrinc(tls *TLS, f uintptr, inc Tsize_t)
TEXT ·Y__freadptrinc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD inc+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__freadptrinc(SB)
	RET

// func Y__freelocale(tls *TLS, l Tlocale_t)
TEXT ·Y__freelocale(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__freelocale(SB)
	RET

// func Y__fseeko(tls *TLS, f uintptr, off Toff_t, whence int32) (r int32)
TEXT ·Y__fseeko(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD off+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__fseeko(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__fseeko_unlocked(tls *TLS, f uintptr, off Toff_t, whence int32) (r int32)
TEXT ·Y__fseeko_unlocked(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD off+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__fseeko_unlocked(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__fseterr(tls *TLS, f uintptr)
TEXT ·Y__fseterr(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fseterr(SB)
	RET

// func Y__fsetlocking(tls *TLS, f uintptr, type1 int32) (r int32)
TEXT ·Y__fsetlocking(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW type1+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__fsetlocking(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__fstat(tls *TLS, fd int32, st uintptr) (r int32)
TEXT ·Y__fstat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD st+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__fstat(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__fstatat(tls *TLS, fd int32, path uintptr, st uintptr, flag int32) (r int32)
TEXT ·Y__fstatat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD st+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flag+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__fstatat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Y__ftello(tls *TLS, f uintptr) (r Toff_t)
TEXT ·Y__ftello(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__ftello(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__ftello_unlocked(tls *TLS, f uintptr) (r Toff_t)
TEXT ·Y__ftello_unlocked(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__ftello_unlocked(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__funcs_on_quick_exit(tls *TLS)
TEXT ·Y__funcs_on_quick_exit(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__funcs_on_quick_exit(SB)
	RET

// func Y__futimesat(tls *TLS, dirfd int32, pathname uintptr, times uintptr) (r int32)
TEXT ·Y__futimesat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD times+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__futimesat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__fwritable(tls *TLS, f uintptr) (r int32)
TEXT ·Y__fwritable(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fwritable(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fwritex(tls *TLS, s uintptr, l Tsize_t, f uintptr) (r Tsize_t)
TEXT ·Y__fwritex(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__fwritex(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__fwriting(tls *TLS, f uintptr) (r int32)
TEXT ·Y__fwriting(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fwriting(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fxstat(tls *TLS, ver int32, fd int32, buf uintptr) (r int32)
TEXT ·Y__fxstat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ver+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW fd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__fxstat(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__fxstatat(tls *TLS, ver int32, fd int32, path uintptr, buf uintptr, flag int32) (r int32)
TEXT ·Y__fxstatat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ver+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW fd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flag+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__fxstatat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Y__get_handler_set(tls *TLS, set uintptr)
TEXT ·Y__get_handler_set(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__get_handler_set(SB)
	RET

// func Y__get_locale(tls *TLS, cat int32, val uintptr) (r uintptr)
TEXT ·Y__get_locale(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW cat+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__get_locale(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__get_resolv_conf(tls *TLS, conf uintptr, search uintptr, search_sz Tsize_t) (r int32)
TEXT ·Y__get_resolv_conf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD conf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD search+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD search_sz+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__get_resolv_conf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__getauxval(tls *TLS, item uint64) (r uint64)
TEXT ·Y__getauxval(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD item+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__getauxval(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__getdelim(tls *TLS, s uintptr, n uintptr, delim int32, f uintptr) (r Tssize_t)
TEXT ·Y__getdelim(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW delim+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD f+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__getdelim(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Y__getgr_a(tls *TLS, name uintptr, gid Tgid_t, gr uintptr, buf uintptr, size uintptr, mem uintptr, nmem uintptr, res uintptr) (r int32)
TEXT ·Y__getgr_a(SB),$80-76
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW gid+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD gr+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buf+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD size+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD mem+48(FP), R0
	MOVD R0, 56(RSP)
	MOVD nmem+56(FP), R0
	MOVD R0, 64(RSP)
	MOVD res+64(FP), R0
	MOVD R0, 72(RSP)
	CALL ·X__getgr_a(SB)
	MOVW 80(RSP), R0
	MOVW R0, r+72(FP)
	RET

// func Y__getgrent_a(tls *TLS, f uintptr, gr uintptr, line uintptr, size uintptr, mem uintptr, nmem uintptr, res uintptr) (r int32)
TEXT ·Y__getgrent_a(SB),$72-68
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD gr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD line+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD mem+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD nmem+48(FP), R0
	MOVD R0, 56(RSP)
	MOVD res+56(FP), R0
	MOVD R0, 64(RSP)
	CALL ·X__getgrent_a(SB)
	MOVW 72(RSP), R0
	MOVW R0, r+64(FP)
	RET

// func Y__getopt_msg(tls *TLS, a uintptr, b uintptr, c uintptr, l Tsize_t)
TEXT ·Y__getopt_msg(SB),$40-40
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
	MOVD l+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__getopt_msg(SB)
	RET

// func Y__getpw_a(tls *TLS, name uintptr, uid Tuid_t, pw uintptr, buf uintptr, size uintptr, res uintptr) (r int32)
TEXT ·Y__getpw_a(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW uid+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD pw+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buf+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD size+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD res+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·X__getpw_a(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Y__getpwent_a(tls *TLS, f uintptr, pw uintptr, line uintptr, size uintptr, res uintptr) (r int32)
TEXT ·Y__getpwent_a(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pw+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD line+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD res+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__getpwent_a(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Y__gettextdomain(tls *TLS) (r uintptr)
TEXT ·Y__gettextdomain(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__gettextdomain(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__gmtime_r(tls *TLS, t uintptr, tm uintptr) (r uintptr)
TEXT ·Y__gmtime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tm+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__gmtime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__h_errno_location(tls *TLS) (r uintptr)
TEXT ·Y__h_errno_location(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__h_errno_location(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__inet_aton(tls *TLS, s0 uintptr, dest uintptr) (r int32)
TEXT ·Y__inet_aton(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s0+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD dest+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__inet_aton(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__init_ssp(tls *TLS, entropy uintptr)
TEXT ·Y__init_ssp(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD entropy+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__init_ssp(SB)
	RET

// func Y__intscan(tls *TLS, f uintptr, base uint32, pok int32, lim uint64) (r uint64)
TEXT ·Y__intscan(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW base+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW pok+20(FP), R0
	MOVW R0, 28(RSP)
	MOVD lim+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__intscan(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__isalnum_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isalnum_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isalnum_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isalpha_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isalpha_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isalpha_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isblank_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isblank_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isblank_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iscntrl_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__iscntrl_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iscntrl_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isdigit_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isdigit_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isdigit_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isfinite(tls *TLS, d float64) (_2 int32)
TEXT ·Y__isfinite(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__isfinite(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__isfinitef(tls *TLS, f float32) (_2 int32)
TEXT ·Y__isfinitef(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW f+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__isfinitef(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__isfinitel(tls *TLS, d float64) (_2 int32)
TEXT ·Y__isfinitel(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__isfinitel(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__isgraph_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isgraph_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isgraph_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__islower_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__islower_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__islower_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isnan(t *TLS, x float64) (_2 int32)
TEXT ·Y__isnan(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__isnan(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__isnanf(t *TLS, arg float32) (_2 int32)
TEXT ·Y__isnanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW arg+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__isnanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__isnanl(t *TLS, arg float64) (_2 int32)
TEXT ·Y__isnanl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD arg+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__isnanl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__isoc99_fscanf(tls *TLS, f uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Y__isoc99_fscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_fscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__isoc99_fwscanf(tls *TLS, f uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Y__isoc99_fwscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_fwscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__isoc99_scanf(tls *TLS, fmt uintptr, va uintptr) (r int32)
TEXT ·Y__isoc99_scanf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isoc99_scanf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isoc99_sscanf(tls *TLS, s uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Y__isoc99_sscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_sscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__isoc99_swscanf(tls *TLS, s uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Y__isoc99_swscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_swscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__isoc99_vfscanf(tls *TLS, f uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Y__isoc99_vfscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_vfscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__isoc99_vfwscanf(tls *TLS, f uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Y__isoc99_vfwscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_vfwscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__isoc99_vscanf(tls *TLS, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Y__isoc99_vscanf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isoc99_vscanf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isoc99_vsscanf(tls *TLS, s uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Y__isoc99_vsscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_vsscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__isoc99_vswscanf(tls *TLS, s uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Y__isoc99_vswscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_vswscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__isoc99_vwscanf(tls *TLS, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Y__isoc99_vwscanf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isoc99_vwscanf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isoc99_wscanf(tls *TLS, fmt uintptr, va uintptr) (r int32)
TEXT ·Y__isoc99_wscanf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isoc99_wscanf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isprint_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isprint_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isprint_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__ispunct_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__ispunct_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__ispunct_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isspace_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isspace_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isspace_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isupper_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isupper_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isupper_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswalnum_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswalnum_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswalnum_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswalpha_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswalpha_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswalpha_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswblank_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswblank_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswblank_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswcntrl_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswcntrl_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswcntrl_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswctype_l(tls *TLS, c Twint_t, t Twctype_t, l Tlocale_t) (r int32)
TEXT ·Y__iswctype_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD t+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__iswctype_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__iswdigit_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswdigit_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswdigit_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswgraph_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswgraph_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswgraph_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswlower_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswlower_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswlower_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswprint_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswprint_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswprint_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswpunct_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswpunct_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswpunct_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswspace_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswspace_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswspace_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswupper_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswupper_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswupper_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__iswxdigit_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Y__iswxdigit_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__iswxdigit_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__isxdigit_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__isxdigit_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isxdigit_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__lctrans(tls *TLS, msg uintptr, lm uintptr) (r uintptr)
TEXT ·Y__lctrans(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD lm+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__lctrans(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__lctrans_cur(tls *TLS, msg uintptr) (r uintptr)
TEXT ·Y__lctrans_cur(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msg+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__lctrans_cur(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__lctrans_impl(tls *TLS, msg uintptr, lm uintptr) (r uintptr)
TEXT ·Y__lctrans_impl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD lm+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__lctrans_impl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__ldexp_cexp(tls *TLS, z complex128, expt int32) (r complex128)
TEXT ·Y__ldexp_cexp(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW expt+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__ldexp_cexp(SB)
	MOVD 40(RSP), R0
	MOVD R0, r_real+32(FP)
	MOVD 48(RSP), R0
	MOVD R0, r_imag+40(FP)
	RET

// func Y__ldexp_cexpf(tls *TLS, z complex64, expt int32) (r complex64)
TEXT ·Y__ldexp_cexpf(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW expt+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__ldexp_cexpf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r_real+24(FP)
	MOVW 36(RSP), R0
	MOVW R0, r_imag+28(FP)
	RET

// func Y__lgamma_r(tls *TLS, x float64, signgamp uintptr) (r1 float64)
TEXT ·Y__lgamma_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD signgamp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__lgamma_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r1+24(FP)
	RET

// func Y__lgammaf_r(tls *TLS, x float32, signgamp uintptr) (r1 float32)
TEXT ·Y__lgammaf_r(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD signgamp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__lgammaf_r(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Y__lgammal_r(tls *TLS, x float64, sg uintptr) (r float64)
TEXT ·Y__lgammal_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sg+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__lgammal_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__libc_current_sigrtmax(tls *TLS) (r int32)
TEXT ·Y__libc_current_sigrtmax(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__libc_current_sigrtmax(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Y__libc_current_sigrtmin(tls *TLS) (r int32)
TEXT ·Y__libc_current_sigrtmin(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__libc_current_sigrtmin(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Y__loc_is_allocated(tls *TLS, loc Tlocale_t) (r int32)
TEXT ·Y__loc_is_allocated(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD loc+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__loc_is_allocated(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__localtime_r(tls *TLS, t uintptr, tm uintptr) (r uintptr)
TEXT ·Y__localtime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tm+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__localtime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__lockfile(tls *TLS, file uintptr) (_2 int32)
TEXT ·Y__lockfile(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD file+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__lockfile(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__lookup_ipliteral(tls *TLS, buf uintptr, name uintptr, family int32) (r int32)
TEXT ·Y__lookup_ipliteral(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW family+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__lookup_ipliteral(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__lookup_name(tls *TLS, buf uintptr, canon uintptr, name uintptr, family int32, flags int32) (r int32)
TEXT ·Y__lookup_name(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD canon+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD name+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW family+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW flags+36(FP), R0
	MOVW R0, 44(RSP)
	CALL ·X__lookup_name(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Y__lookup_serv(tls *TLS, buf uintptr, name uintptr, proto int32, socktype int32, flags int32) (r int32)
TEXT ·Y__lookup_serv(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW proto+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW socktype+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__lookup_serv(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Y__lseek(tls *TLS, fd int32, offset Toff_t, whence int32) (r Toff_t)
TEXT ·Y__lseek(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD offset+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__lseek(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__lsysinfo(tls *TLS, info uintptr) (r int32)
TEXT ·Y__lsysinfo(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD info+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__lsysinfo(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__lxstat(tls *TLS, ver int32, path uintptr, buf uintptr) (r int32)
TEXT ·Y__lxstat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ver+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__lxstat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__madvise(tls *TLS, addr uintptr, len1 Tsize_t, advice int32) (r int32)
TEXT ·Y__madvise(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW advice+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__madvise(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__map_file(tls *TLS, pathname uintptr, size uintptr) (r uintptr)
TEXT ·Y__map_file(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__map_file(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__math_divzero(tls *TLS, sign Tuint32_t) (r float64)
TEXT ·Y__math_divzero(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sign+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__math_divzero(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__math_divzerof(tls *TLS, sign Tuint32_t) (r float32)
TEXT ·Y__math_divzerof(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sign+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__math_divzerof(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__math_invalid(tls *TLS, x float64) (r float64)
TEXT ·Y__math_invalid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__math_invalid(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__math_invalidf(tls *TLS, x float32) (r float32)
TEXT ·Y__math_invalidf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__math_invalidf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__math_oflow(tls *TLS, sign Tuint32_t) (r float64)
TEXT ·Y__math_oflow(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sign+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__math_oflow(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__math_oflowf(tls *TLS, sign Tuint32_t) (r float32)
TEXT ·Y__math_oflowf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sign+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__math_oflowf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__math_uflow(tls *TLS, sign Tuint32_t) (r float64)
TEXT ·Y__math_uflow(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sign+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__math_uflow(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__math_uflowf(tls *TLS, sign Tuint32_t) (r float32)
TEXT ·Y__math_uflowf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sign+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__math_uflowf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__math_xflow(tls *TLS, sign Tuint32_t, y2 float64) (r float64)
TEXT ·Y__math_xflow(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sign+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD y2+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__math_xflow(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__math_xflowf(tls *TLS, sign Tuint32_t, y2 float32) (r float32)
TEXT ·Y__math_xflowf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sign+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y2+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·X__math_xflowf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__memrchr(tls *TLS, m uintptr, c int32, n Tsize_t) (r uintptr)
TEXT ·Y__memrchr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__memrchr(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__mkostemps(tls *TLS, template uintptr, len1 int32, flags int32) (r int32)
TEXT ·Y__mkostemps(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW len1+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW flags+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__mkostemps(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__mmap(tls *TLS, start uintptr, len1 Tsize_t, prot int32, flags int32, fd int32, off Toff_t) (r uintptr)
TEXT ·Y__mmap(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD start+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW prot+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW fd+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD off+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__mmap(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Y__mo_lookup(tls *TLS, p uintptr, size Tsize_t, s uintptr) (r uintptr)
TEXT ·Y__mo_lookup(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD s+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__mo_lookup(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__month_to_secs(tls *TLS, month int32, is_leap int32) (r int32)
TEXT ·Y__month_to_secs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW month+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW is_leap+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·X__month_to_secs(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__mprotect(tls *TLS, addr uintptr, len1 Tsize_t, prot int32) (r int32)
TEXT ·Y__mprotect(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW prot+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__mprotect(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__mremap(tls *TLS, old_addr uintptr, old_len Tsize_t, new_len Tsize_t, flags int32, va uintptr) (r uintptr)
TEXT ·Y__mremap(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD old_addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD old_len+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD new_len+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD va+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__mremap(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Y__munmap(tls *TLS, start uintptr, len1 Tsize_t) (r int32)
TEXT ·Y__munmap(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD start+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__munmap(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__newlocale(tls *TLS, mask int32, name uintptr, loc Tlocale_t) (r Tlocale_t)
TEXT ·Y__newlocale(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW mask+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD loc+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__newlocale(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__nl_langinfo(tls *TLS, item Tnl_item) (r uintptr)
TEXT ·Y__nl_langinfo(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW item+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__nl_langinfo(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__nl_langinfo_l(tls *TLS, item Tnl_item, loc Tlocale_t) (r uintptr)
TEXT ·Y__nl_langinfo_l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW item+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD loc+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__nl_langinfo_l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__nscd_query(tls *TLS, req Tint32_t, key uintptr, buf uintptr, len1 Tsize_t, swap uintptr) (r uintptr)
TEXT ·Y__nscd_query(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW req+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD key+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD len1+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD swap+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__nscd_query(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Y__ofl_add(tls *TLS, f uintptr) (r uintptr)
TEXT ·Y__ofl_add(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__ofl_add(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__ofl_lock(tls *TLS) (r uintptr)
TEXT ·Y__ofl_lock(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__ofl_lock(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__ofl_unlock(tls *TLS)
TEXT ·Y__ofl_unlock(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__ofl_unlock(SB)
	RET

// func Y__overflow(tls *TLS, f uintptr, _c int32) (r int32)
TEXT ·Y__overflow(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW _c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__overflow(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__pleval(tls *TLS, s uintptr, n uint64) (r uint64)
TEXT ·Y__pleval(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__pleval(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__posix_getopt(tls *TLS, argc int32, argv uintptr, optstring uintptr) (r int32)
TEXT ·Y__posix_getopt(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW argc+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD optstring+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__posix_getopt(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__procfdname(tls *TLS, buf uintptr, fd uint32)
TEXT ·Y__procfdname(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW fd+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__procfdname(SB)
	RET

// func Y__pthread_testcancel(tls *TLS)
TEXT ·Y__pthread_testcancel(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__pthread_testcancel(SB)
	RET

// func Y__ptsname_r(tls *TLS, fd int32, buf uintptr, len1 Tsize_t) (r int32)
TEXT ·Y__ptsname_r(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__ptsname_r(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__putenv(tls *TLS, s uintptr, l Tsize_t, r uintptr) (r1 int32)
TEXT ·Y__putenv(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD r+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__putenv(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Y__qsort_r(tls *TLS, base uintptr, nel Tsize_t, width Tsize_t, __ccgo_fp_cmp Tcmpfun, arg uintptr)
TEXT ·Y__qsort_r(SB),$56-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 56(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 56(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cmp+32(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal___qsort_r_3(SB)	// Create the closure for calling __ccgo_fp_cmp
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD base+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD nel+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD width+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD 56(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 40(RSP)
	MOVD arg+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__qsort_r(SB)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0___qsort_r_3(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD _3+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD __ccgo_fp+32(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 40(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y__rand48_step(tls *TLS, xi uintptr, lc uintptr) (r Tuint64_t)
TEXT ·Y__rand48_step(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD xi+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD lc+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__rand48_step(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__register_locked_file(tls *TLS, f uintptr, self Tpthread_t)
TEXT ·Y__register_locked_file(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD self+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__register_locked_file(SB)
	RET

// func Y__rem_pio2(tls *TLS, x float64, y uintptr) (r1 int32)
TEXT ·Y__rem_pio2(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__rem_pio2(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Y__rem_pio2_large(tls *TLS, x uintptr, y uintptr, e0 int32, nx int32, prec int32) (r int32)
TEXT ·Y__rem_pio2_large(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW e0+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW nx+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW prec+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__rem_pio2_large(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Y__rem_pio2f(tls *TLS, x float32, y uintptr) (r int32)
TEXT ·Y__rem_pio2f(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__rem_pio2f(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__res_mkquery(tls *TLS, op int32, dname uintptr, class int32, type1 int32, data uintptr, datalen int32, newrr uintptr, buf uintptr, buflen int32) (r int32)
TEXT ·Y__res_mkquery(SB),$80-76
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW op+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD dname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW class+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW type1+28(FP), R0
	MOVW R0, 36(RSP)
	MOVD data+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW datalen+40(FP), R0
	MOVW R0, 48(RSP)
	MOVD newrr+48(FP), R0
	MOVD R0, 56(RSP)
	MOVD buf+56(FP), R0
	MOVD R0, 64(RSP)
	MOVW buflen+64(FP), R0
	MOVW R0, 72(RSP)
	CALL ·X__res_mkquery(SB)
	MOVW 80(RSP), R0
	MOVW R0, r+72(FP)
	RET

// func Y__res_msend(tls *TLS, nqueries int32, queries uintptr, qlens uintptr, answers uintptr, alens uintptr, asize int32) (r int32)
TEXT ·Y__res_msend(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW nqueries+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD queries+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD qlens+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD answers+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD alens+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW asize+48(FP), R0
	MOVW R0, 56(RSP)
	CALL ·X__res_msend(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Y__res_msend_rc(tls *TLS, nqueries int32, queries uintptr, qlens uintptr, answers uintptr, alens uintptr, asize int32, conf uintptr) (r1 int32)
TEXT ·Y__res_msend_rc(SB),$72-68
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW nqueries+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD queries+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD qlens+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD answers+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD alens+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW asize+48(FP), R0
	MOVW R0, 56(RSP)
	MOVD conf+56(FP), R0
	MOVD R0, 64(RSP)
	CALL ·X__res_msend_rc(SB)
	MOVW 72(RSP), R0
	MOVW R0, r1+64(FP)
	RET

// func Y__res_send(tls *TLS, _msg uintptr, _msglen int32, _answer uintptr, _anslen int32) (r1 int32)
TEXT ·Y__res_send(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _msg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW _msglen+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD _answer+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _anslen+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·X__res_send(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Y__res_state(tls *TLS) (r uintptr)
TEXT ·Y__res_state(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__res_state(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Y__reset_tls(tls *TLS)
TEXT ·Y__reset_tls(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__reset_tls(SB)
	RET

// func Y__restore(tls *TLS)
TEXT ·Y__restore(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__restore(SB)
	RET

// func Y__restore_rt(tls *TLS)
TEXT ·Y__restore_rt(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__restore_rt(SB)
	RET

// func Y__restore_sigs(tls *TLS, set uintptr)
TEXT ·Y__restore_sigs(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__restore_sigs(SB)
	RET

// func Y__rtnetlink_enumerate(tls *TLS, link_af int32, addr_af int32, __ccgo_fp_cb uintptr, ctx uintptr) (r1 int32)
TEXT ·Y__rtnetlink_enumerate(SB),$48-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 48(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 48(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cb+16(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal___rtnetlink_enumerate_2(SB)	// Create the closure for calling __ccgo_fp_cb
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW link_af+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW addr_af+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD 48(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 24(RSP)
	MOVD ctx+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__rtnetlink_enumerate(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0___rtnetlink_enumerate_2(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Y__sched_cpucount(tls *TLS, size Tsize_t, set uintptr) (r int32)
TEXT ·Y__sched_cpucount(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD size+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD set+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__sched_cpucount(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__secs_to_tm(tls *TLS, t int64, tm uintptr) (r int32)
TEXT ·Y__secs_to_tm(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tm+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__secs_to_tm(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__secs_to_zone(tls *TLS, t int64, local int32, isdst uintptr, offset uintptr, oppoff uintptr, zonename uintptr)
TEXT ·Y__secs_to_zone(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW local+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD isdst+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD offset+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD oppoff+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD zonename+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·X__secs_to_zone(SB)
	RET

// func Y__setxid(tls *TLS, nr int32, id int32, eid int32, sid int32) (r int32)
TEXT ·Y__setxid(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW nr+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW id+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW eid+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW sid+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__setxid(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__shgetc(tls *TLS, f uintptr) (r int32)
TEXT ·Y__shgetc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__shgetc(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__shlim(tls *TLS, f uintptr, lim Toff_t)
TEXT ·Y__shlim(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD lim+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__shlim(SB)
	RET

// func Y__shm_mapname(tls *TLS, name uintptr, buf uintptr) (r uintptr)
TEXT ·Y__shm_mapname(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__shm_mapname(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__sigaction(tls *TLS, sig int32, sa uintptr, old uintptr) (r1 int32)
TEXT ·Y__sigaction(SB),$40-36
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
	CALL ·X__sigaction(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Y__signbit(tls *TLS, x float64) (r int32)
TEXT ·Y__signbit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__signbit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__signbitf(tls *TLS, x float32) (r int32)
TEXT ·Y__signbitf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__signbitf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__signbitl(tls *TLS, x float64) (r int32)
TEXT ·Y__signbitl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__signbitl(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__sigsetjmp_tail(tls *TLS, jb uintptr, ret int32) (r int32)
TEXT ·Y__sigsetjmp_tail(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD jb+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW ret+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__sigsetjmp_tail(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__sin(tls *TLS, x float64, y float64, iy int32) (r1 float64)
TEXT ·Y__sin(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW iy+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__sin(SB)
	MOVD 40(RSP), R0
	MOVD R0, r1+32(FP)
	RET

// func Y__sindf(tls *TLS, x float64) (r1 float32)
TEXT ·Y__sindf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__sindf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Y__stack_chk_fail(tls *TLS)
TEXT ·Y__stack_chk_fail(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__stack_chk_fail(SB)
	RET

// func Y__stack_chk_fail_local(tls *TLS)
TEXT ·Y__stack_chk_fail_local(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__stack_chk_fail_local(SB)
	RET

// func Y__stdio_close(tls *TLS, f uintptr) (r int32)
TEXT ·Y__stdio_close(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__stdio_close(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__stdio_exit(tls *TLS)
TEXT ·Y__stdio_exit(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__stdio_exit(SB)
	RET

// func Y__stdio_exit_needed(tls *TLS)
TEXT ·Y__stdio_exit_needed(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__stdio_exit_needed(SB)
	RET

// func Y__stdio_read(tls *TLS, f uintptr, buf uintptr, len1 Tsize_t) (r Tsize_t)
TEXT ·Y__stdio_read(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__stdio_read(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__stdio_seek(tls *TLS, f uintptr, off Toff_t, whence int32) (r Toff_t)
TEXT ·Y__stdio_seek(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD off+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__stdio_seek(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__stdio_write(tls *TLS, f uintptr, buf uintptr, len1 Tsize_t) (r Tsize_t)
TEXT ·Y__stdio_write(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__stdio_write(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__stdout_write(tls *TLS, f uintptr, buf uintptr, len1 Tsize_t) (r Tsize_t)
TEXT ·Y__stdout_write(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__stdout_write(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__stpcpy(tls *TLS, d uintptr, s uintptr) (r uintptr)
TEXT ·Y__stpcpy(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__stpcpy(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__stpncpy(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Y__stpncpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__stpncpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strcasecmp_l(tls *TLS, l uintptr, r uintptr, loc Tlocale_t) (r1 int32)
TEXT ·Y__strcasecmp_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD loc+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__strcasecmp_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Y__strchrnul(tls *TLS, s uintptr, c int32) (r uintptr)
TEXT ·Y__strchrnul(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__strchrnul(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__strcoll_l(tls *TLS, l uintptr, r uintptr, loc Tlocale_t) (r1 int32)
TEXT ·Y__strcoll_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD loc+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__strcoll_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Y__strerror_l(tls *TLS, e int32, loc Tlocale_t) (r uintptr)
TEXT ·Y__strerror_l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW e+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD loc+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__strerror_l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__strftime_fmt_1(tls *TLS, s uintptr, l uintptr, f int32, tm uintptr, loc Tlocale_t, pad int32) (_7 uintptr)
TEXT ·Y__strftime_fmt_1(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW f+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD tm+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD loc+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW pad+48(FP), R0
	MOVW R0, 56(RSP)
	CALL ·X__strftime_fmt_1(SB)
	MOVD 64(RSP), R0
	MOVD R0, _7+56(FP)
	RET

// func Y__strftime_l(tls *TLS, s uintptr, n Tsize_t, f uintptr, tm uintptr, loc Tlocale_t) (r Tsize_t)
TEXT ·Y__strftime_l(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD tm+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD loc+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__strftime_l(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Y__strncasecmp_l(tls *TLS, l uintptr, r uintptr, n Tsize_t, loc Tlocale_t) (r1 int32)
TEXT ·Y__strncasecmp_l(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD loc+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__strncasecmp_l(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Y__strtod_l(tls *TLS, s uintptr, p uintptr, l Tlocale_t) (r float64)
TEXT ·Y__strtod_l(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__strtod_l(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strtof_l(tls *TLS, s uintptr, p uintptr, l Tlocale_t) (r float32)
TEXT ·Y__strtof_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__strtof_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__strtoimax_internal(tls *TLS, s uintptr, p uintptr, base int32) (r Tintmax_t)
TEXT ·Y__strtoimax_internal(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__strtoimax_internal(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strtol_internal(tls *TLS, s uintptr, p uintptr, base int32) (r int64)
TEXT ·Y__strtol_internal(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__strtol_internal(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strtold_l(tls *TLS, s uintptr, p uintptr, l Tlocale_t) (r float64)
TEXT ·Y__strtold_l(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__strtold_l(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strtoll_internal(tls *TLS, s uintptr, p uintptr, base int32) (r int64)
TEXT ·Y__strtoll_internal(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__strtoll_internal(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strtoul_internal(tls *TLS, s uintptr, p uintptr, base int32) (r uint64)
TEXT ·Y__strtoul_internal(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__strtoul_internal(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strtoull_internal(tls *TLS, s uintptr, p uintptr, base int32) (r uint64)
TEXT ·Y__strtoull_internal(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__strtoull_internal(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strtoumax_internal(tls *TLS, s uintptr, p uintptr, base int32) (r Tuintmax_t)
TEXT ·Y__strtoumax_internal(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__strtoumax_internal(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__strxfrm_l(tls *TLS, dest uintptr, src uintptr, n Tsize_t, loc Tlocale_t) (r Tsize_t)
TEXT ·Y__strxfrm_l(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD loc+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__strxfrm_l(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Y__sync_add_and_fetch_uint32(t *TLS, p uintptr, v uint32) (_3 uint32)
TEXT ·Y__sync_add_and_fetch_uint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW v+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__sync_add_and_fetch_uint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__sync_sub_and_fetch_uint32(t *TLS, p uintptr, v uint32) (_3 uint32)
TEXT ·Y__sync_sub_and_fetch_uint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW v+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__sync_sub_and_fetch_uint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__sync_synchronize(t *TLS)
TEXT ·Y__sync_synchronize(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__sync_synchronize(SB)
	RET

// func Y__sync_val_compare_and_swapInt16(t *TLS, ptr uintptr, oldval, newval int16) (r int16)
TEXT ·Y__sync_val_compare_and_swapInt16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH oldval+16(FP), R0
	MOVH R0, 24(RSP)
	MOVH newval+18(FP), R0
	MOVH R0, 26(RSP)
	CALL ·X__sync_val_compare_and_swapInt16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__sync_val_compare_and_swapInt32(t *TLS, ptr uintptr, oldval, newval int32) (r int32)
TEXT ·Y__sync_val_compare_and_swapInt32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW oldval+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW newval+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__sync_val_compare_and_swapInt32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__sync_val_compare_and_swapInt64(t *TLS, ptr uintptr, oldval, newval int64) (r int64)
TEXT ·Y__sync_val_compare_and_swapInt64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD oldval+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD newval+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__sync_val_compare_and_swapInt64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__sync_val_compare_and_swapInt8(t *TLS, ptr uintptr, oldval, newval int8) (r int8)
TEXT ·Y__sync_val_compare_and_swapInt8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB oldval+16(FP), R0
	MOVB R0, 24(RSP)
	MOVB newval+17(FP), R0
	MOVB R0, 25(RSP)
	CALL ·X__sync_val_compare_and_swapInt8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__sync_val_compare_and_swapUint16(t *TLS, ptr uintptr, oldval, newval uint16) (r uint16)
TEXT ·Y__sync_val_compare_and_swapUint16(SB),$32-26
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH oldval+16(FP), R0
	MOVH R0, 24(RSP)
	MOVH newval+18(FP), R0
	MOVH R0, 26(RSP)
	CALL ·X__sync_val_compare_and_swapUint16(SB)
	MOVH 32(RSP), R0
	MOVH R0, r+24(FP)
	RET

// func Y__sync_val_compare_and_swapUint32(t *TLS, ptr uintptr, oldval, newval uint32) (r uint32)
TEXT ·Y__sync_val_compare_and_swapUint32(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW oldval+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW newval+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·X__sync_val_compare_and_swapUint32(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__sync_val_compare_and_swapUint64(t *TLS, ptr uintptr, oldval, newval uint64) (r uint64)
TEXT ·Y__sync_val_compare_and_swapUint64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD oldval+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD newval+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__sync_val_compare_and_swapUint64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Y__sync_val_compare_and_swapUint8(t *TLS, ptr uintptr, oldval, newval uint8) (r uint8)
TEXT ·Y__sync_val_compare_and_swapUint8(SB),$32-25
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVB oldval+16(FP), R0
	MOVB R0, 24(RSP)
	MOVB newval+17(FP), R0
	MOVB R0, 25(RSP)
	CALL ·X__sync_val_compare_and_swapUint8(SB)
	MOVB 32(RSP), R0
	MOVB R0, r+24(FP)
	RET

// func Y__syscall0(tls *TLS, n long) (_2 long)
TEXT ·Y__syscall0(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__syscall0(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__syscall1(tls *TLS, n, a1 long) (_2 long)
TEXT ·Y__syscall1(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__syscall1(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Y__syscall2(tls *TLS, n, a1, a2 long) (_2 long)
TEXT ·Y__syscall2(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD a2+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__syscall2(SB)
	MOVD 40(RSP), R0
	MOVD R0, _2+32(FP)
	RET

// func Y__syscall3(tls *TLS, n, a1, a2, a3 long) (_2 long)
TEXT ·Y__syscall3(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD a2+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD a3+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__syscall3(SB)
	MOVD 48(RSP), R0
	MOVD R0, _2+40(FP)
	RET

// func Y__syscall4(tls *TLS, n, a1, a2, a3, a4 long) (_2 long)
TEXT ·Y__syscall4(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD a2+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD a3+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD a4+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__syscall4(SB)
	MOVD 56(RSP), R0
	MOVD R0, _2+48(FP)
	RET

// func Y__syscall5(tls *TLS, n, a1, a2, a3, a4, a5 long) (_2 long)
TEXT ·Y__syscall5(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD a2+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD a3+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD a4+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD a5+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·X__syscall5(SB)
	MOVD 64(RSP), R0
	MOVD R0, _2+56(FP)
	RET

// func Y__syscall6(tls *TLS, n, a1, a2, a3, a4, a5, a6 long) (_2 long)
TEXT ·Y__syscall6(SB),$72-72
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD a2+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD a3+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD a4+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD a5+48(FP), R0
	MOVD R0, 56(RSP)
	MOVD a6+56(FP), R0
	MOVD R0, 64(RSP)
	CALL ·X__syscall6(SB)
	MOVD 72(RSP), R0
	MOVD R0, _2+64(FP)
	RET

// func Y__syscall_ret(tls *TLS, r uint64) (r1 int64)
TEXT ·Y__syscall_ret(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD r+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__syscall_ret(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Y__tan(tls *TLS, x float64, y float64, odd int32) (r1 float64)
TEXT ·Y__tan(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW odd+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__tan(SB)
	MOVD 40(RSP), R0
	MOVD R0, r1+32(FP)
	RET

// func Y__tandf(tls *TLS, x float64, odd int32) (r1 float32)
TEXT ·Y__tandf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW odd+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X__tandf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Y__tm_to_secs(tls *TLS, tm uintptr) (r int64)
TEXT ·Y__tm_to_secs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tm+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__tm_to_secs(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__tm_to_tzname(tls *TLS, tm uintptr) (r uintptr)
TEXT ·Y__tm_to_tzname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tm+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__tm_to_tzname(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__tolower_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__tolower_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__tolower_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__toread(tls *TLS, f uintptr) (r int32)
TEXT ·Y__toread(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__toread(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__toread_needs_stdio_exit(tls *TLS)
TEXT ·Y__toread_needs_stdio_exit(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__toread_needs_stdio_exit(SB)
	RET

// func Y__toupper_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Y__toupper_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__toupper_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__towctrans_l(tls *TLS, c Twint_t, t Twctrans_t, l Tlocale_t) (r Twint_t)
TEXT ·Y__towctrans_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD t+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__towctrans_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__towlower_l(tls *TLS, c Twint_t, l Tlocale_t) (r Twint_t)
TEXT ·Y__towlower_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__towlower_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__towrite(tls *TLS, f uintptr) (r int32)
TEXT ·Y__towrite(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__towrite(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__towrite_needs_stdio_exit(tls *TLS)
TEXT ·Y__towrite_needs_stdio_exit(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__towrite_needs_stdio_exit(SB)
	RET

// func Y__towupper_l(tls *TLS, c Twint_t, l Tlocale_t) (r Twint_t)
TEXT ·Y__towupper_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__towupper_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Y__tre_mem_alloc_impl(tls *TLS, mem Ttre_mem_t, provided int32, provided_block uintptr, zero int32, size Tsize_t) (r uintptr)
TEXT ·Y__tre_mem_alloc_impl(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD mem+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW provided+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD provided_block+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW zero+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD size+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__tre_mem_alloc_impl(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Y__tre_mem_destroy(tls *TLS, mem Ttre_mem_t)
TEXT ·Y__tre_mem_destroy(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD mem+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__tre_mem_destroy(SB)
	RET

// func Y__tre_mem_new_impl(tls *TLS, provided int32, provided_block uintptr) (r Ttre_mem_t)
TEXT ·Y__tre_mem_new_impl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW provided+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD provided_block+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__tre_mem_new_impl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__tsearch_balance(tls *TLS, p uintptr) (r int32)
TEXT ·Y__tsearch_balance(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__tsearch_balance(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__uflow(tls *TLS, f uintptr) (r int32)
TEXT ·Y__uflow(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__uflow(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__unlist_locked_file(tls *TLS, f uintptr)
TEXT ·Y__unlist_locked_file(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__unlist_locked_file(SB)
	RET

// func Y__unlockfile(tls *TLS, file uintptr)
TEXT ·Y__unlockfile(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD file+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__unlockfile(SB)
	RET

// func Y__uselocale(tls *TLS, new1 Tlocale_t) (r Tlocale_t)
TEXT ·Y__uselocale(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD new1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__uselocale(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__vm_wait(tls *TLS)
TEXT ·Y__vm_wait(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__vm_wait(SB)
	RET

// func Y__wcscoll_l(tls *TLS, l uintptr, r uintptr, locale Tlocale_t) (r1 int32)
TEXT ·Y__wcscoll_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD locale+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__wcscoll_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Y__wcsftime_l(tls *TLS, s uintptr, n Tsize_t, f uintptr, tm uintptr, loc Tlocale_t) (r Tsize_t)
TEXT ·Y__wcsftime_l(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD tm+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD loc+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·X__wcsftime_l(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Y__wcsxfrm_l(tls *TLS, dest uintptr, src uintptr, n Tsize_t, loc Tlocale_t) (r Tsize_t)
TEXT ·Y__wcsxfrm_l(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD loc+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__wcsxfrm_l(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Y__wctrans_l(tls *TLS, s uintptr, l Tlocale_t) (r Twctrans_t)
TEXT ·Y__wctrans_l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__wctrans_l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__wctype_l(tls *TLS, s uintptr, l Tlocale_t) (r Twctype_t)
TEXT ·Y__wctype_l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__wctype_l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Y__xmknod(tls *TLS, ver int32, path uintptr, mode Tmode_t, dev uintptr) (r int32)
TEXT ·Y__xmknod(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ver+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mode+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD dev+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__xmknod(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Y__xmknodat(tls *TLS, ver int32, fd int32, path uintptr, mode Tmode_t, dev uintptr) (r int32)
TEXT ·Y__xmknodat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ver+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW fd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mode+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD dev+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__xmknodat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Y__xpg_basename(tls *TLS, s uintptr) (r uintptr)
TEXT ·Y__xpg_basename(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__xpg_basename(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Y__xpg_strerror_r(tls *TLS, err int32, buf uintptr, buflen Tsize_t) (r int32)
TEXT ·Y__xpg_strerror_r(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW err+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buflen+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__xpg_strerror_r(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__xstat(tls *TLS, ver int32, path uintptr, buf uintptr) (r int32)
TEXT ·Y__xstat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ver+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__xstat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Y__year_to_secs(tls *TLS, year int64, is_leap uintptr) (_3 int64)
TEXT ·Y__year_to_secs(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD year+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD is_leap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__year_to_secs(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Y_exit(tls *TLS, status int32)
TEXT ·Y_exit(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW status+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X_exit(SB)
	RET

// func Y_flushlbf(tls *TLS)
TEXT ·Y_flushlbf(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X_flushlbf(SB)
	RET

// func Y_longjmp(t *TLS, env uintptr, val int32)
TEXT ·Y_longjmp(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD env+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X_longjmp(SB)
	RET

// func Y_obstack_begin(t *TLS, obstack uintptr, size, alignment int32, chunkfun, freefun uintptr) (_4 int32)
TEXT ·Y_obstack_begin(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD obstack+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW size+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW alignment+20(FP), R0
	MOVW R0, 28(RSP)
	MOVD chunkfun+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD freefun+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X_obstack_begin(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Y_obstack_newchunk(t *TLS, obstack uintptr, length int32) (_3 int32)
TEXT ·Y_obstack_newchunk(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD obstack+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW length+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X_obstack_newchunk(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y_pthread_cleanup_pop(tls *TLS, _ uintptr, run int32)
TEXT ·Y_pthread_cleanup_pop(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW run+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·X_pthread_cleanup_pop(SB)
	RET

// func Y_setjmp(t *TLS, env uintptr) (_2 int32)
TEXT ·Y_setjmp(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD env+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X_setjmp(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ya64l(tls *TLS, s uintptr) (r int64)
TEXT ·Ya64l(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xa64l(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yabort(tls *TLS)
TEXT ·Yabort(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xabort(SB)
	RET

// func Yabs(tls *TLS, a int32) (r int32)
TEXT ·Yabs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW a+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xabs(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yaccept(tls *TLS, fd int32, addr uintptr, len1 uintptr) (r1 int32)
TEXT ·Yaccept(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xaccept(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Yaccept4(tls *TLS, fd int32, addr uintptr, len1 uintptr, flg int32) (r1 int32)
TEXT ·Yaccept4(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flg+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xaccept4(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Yaccess(tls *TLS, filename uintptr, amode int32) (r int32)
TEXT ·Yaccess(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW amode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xaccess(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yacct(tls *TLS, filename uintptr) (r int32)
TEXT ·Yacct(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xacct(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yacos(tls *TLS, x float64) (r float64)
TEXT ·Yacos(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xacos(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yacosf(tls *TLS, x float32) (r float32)
TEXT ·Yacosf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xacosf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yacosh(tls *TLS, x float64) (r float64)
TEXT ·Yacosh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xacosh(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yacoshf(tls *TLS, x float32) (r float32)
TEXT ·Yacoshf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xacoshf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yacoshl(tls *TLS, x float64) (r float64)
TEXT ·Yacoshl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xacoshl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yacosl(tls *TLS, x float64) (r float64)
TEXT ·Yacosl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xacosl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yaddmntent(tls *TLS, f uintptr, mnt uintptr) (r int32)
TEXT ·Yaddmntent(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mnt+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xaddmntent(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yadjtime(tls *TLS, in uintptr, out uintptr) (r int32)
TEXT ·Yadjtime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD in+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD out+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xadjtime(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yadjtimex(tls *TLS, tx uintptr) (r int32)
TEXT ·Yadjtimex(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tx+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xadjtimex(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yalarm(tls *TLS, seconds uint32) (r uint32)
TEXT ·Yalarm(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW seconds+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xalarm(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yalloca(tls *TLS, size Tsize_t) (_2 uintptr)
TEXT ·Yalloca(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD size+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xalloca(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yalphasort(tls *TLS, a uintptr, b uintptr) (r int32)
TEXT ·Yalphasort(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xalphasort(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yasctime(tls *TLS, tm uintptr) (r uintptr)
TEXT ·Yasctime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tm+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xasctime(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yasctime_r(tls *TLS, tm uintptr, buf uintptr) (r uintptr)
TEXT ·Yasctime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tm+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xasctime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yasin(tls *TLS, x float64) (r1 float64)
TEXT ·Yasin(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xasin(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Yasinf(tls *TLS, x float32) (r float32)
TEXT ·Yasinf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xasinf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yasinh(tls *TLS, x3 float64) (r float64)
TEXT ·Yasinh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xasinh(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yasinhf(tls *TLS, x3 float32) (r float32)
TEXT ·Yasinhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xasinhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yasinhl(tls *TLS, x float64) (r float64)
TEXT ·Yasinhl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xasinhl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yasinl(tls *TLS, x float64) (r float64)
TEXT ·Yasinl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xasinl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yasprintf(tls *TLS, s uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Yasprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xasprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yatan(tls *TLS, x3 float64) (r float64)
TEXT ·Yatan(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatan(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yatan2(tls *TLS, y float64, x float64) (r float64)
TEXT ·Yatan2(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD y+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD x+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xatan2(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yatan2f(tls *TLS, y float32, x float32) (r float32)
TEXT ·Yatan2f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW y+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW x+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xatan2f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yatan2l(tls *TLS, y float64, x float64) (r float64)
TEXT ·Yatan2l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD y+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD x+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xatan2l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yatanf(tls *TLS, x3 float32) (r float32)
TEXT ·Yatanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xatanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yatanh(tls *TLS, x3 float64) (r float64)
TEXT ·Yatanh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatanh(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yatanhf(tls *TLS, x3 float32) (r float32)
TEXT ·Yatanhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xatanhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yatanhl(tls *TLS, x float64) (r float64)
TEXT ·Yatanhl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatanhl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yatanl(tls *TLS, x float64) (r float64)
TEXT ·Yatanl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatanl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yatof(tls *TLS, s uintptr) (r float64)
TEXT ·Yatof(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatof(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yatoi(tls *TLS, s uintptr) (r int32)
TEXT ·Yatoi(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatoi(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yatol(tls *TLS, s uintptr) (r int64)
TEXT ·Yatol(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatol(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yatoll(tls *TLS, s uintptr) (r int64)
TEXT ·Yatoll(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatoll(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ybacktrace(t *TLS, buf uintptr, size int32) (_3 int32)
TEXT ·Ybacktrace(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW size+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xbacktrace(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ybacktrace_symbols_fd(t *TLS, buffer uintptr, size, fd int32)
TEXT ·Ybacktrace_symbols_fd(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buffer+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW size+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW fd+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xbacktrace_symbols_fd(SB)
	RET

// func Ybasename(tls *TLS, s uintptr) (r uintptr)
TEXT ·Ybasename(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xbasename(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ybcmp(tls *TLS, s1 uintptr, s2 uintptr, n Tsize_t) (r int32)
TEXT ·Ybcmp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xbcmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ybcopy(tls *TLS, s1 uintptr, s2 uintptr, n Tsize_t)
TEXT ·Ybcopy(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xbcopy(SB)
	RET

// func Ybind(tls *TLS, fd int32, addr uintptr, len1 Tsocklen_t) (r1 int32)
TEXT ·Ybind(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW len1+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xbind(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ybind_textdomain_codeset(tls *TLS, domainname uintptr, codeset uintptr) (r uintptr)
TEXT ·Ybind_textdomain_codeset(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD domainname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD codeset+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xbind_textdomain_codeset(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ybindtextdomain(tls *TLS, domainname uintptr, dirname uintptr) (r uintptr)
TEXT ·Ybindtextdomain(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD domainname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD dirname+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xbindtextdomain(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ybrk(tls *TLS, end uintptr) (r int32)
TEXT ·Ybrk(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD end+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xbrk(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ybsearch(tls *TLS, key uintptr, base uintptr, nel Tsize_t, width Tsize_t, __ccgo_fp_cmp uintptr) (r uintptr)
TEXT ·Ybsearch(SB),$64-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 64(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 64(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cmp+40(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_bsearch_4(SB)	// Create the closure for calling __ccgo_fp_cmp
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD base+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nel+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD width+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD 64(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 48(RSP)
	CALL ·Xbsearch(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_bsearch_4(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ybtowc(tls *TLS, c int32) (r Twint_t)
TEXT ·Ybtowc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xbtowc(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ybzero(tls *TLS, s uintptr, n Tsize_t)
TEXT ·Ybzero(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xbzero(SB)
	RET

// func Yc16rtomb(tls *TLS, s uintptr, c16 Tchar16_t, ps uintptr) (r Tsize_t)
TEXT ·Yc16rtomb(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH c16+16(FP), R0
	MOVH R0, 24(RSP)
	MOVD ps+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xc16rtomb(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yc32rtomb(tls *TLS, s uintptr, c32 Tchar32_t, ps uintptr) (r Tsize_t)
TEXT ·Yc32rtomb(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c32+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD ps+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xc32rtomb(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ycabs(tls *TLS, z complex128) (r float64)
TEXT ·Ycabs(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcabs(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycabsf(tls *TLS, z complex64) (r float32)
TEXT ·Ycabsf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcabsf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycabsl(tls *TLS, z complex128) (r float64)
TEXT ·Ycabsl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcabsl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycacos(tls *TLS, z complex128) (r complex128)
TEXT ·Ycacos(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcacos(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycacosf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycacosf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcacosf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycacosh(tls *TLS, z complex128) (r complex128)
TEXT ·Ycacosh(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcacosh(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycacoshf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycacoshf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcacoshf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycacoshl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycacoshl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcacoshl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycacosl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycacosl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcacosl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycalloc(tls *TLS, m Tsize_t, n Tsize_t) (r uintptr)
TEXT ·Ycalloc(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcalloc(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycapget(tls *TLS, a uintptr, b uintptr) (r int32)
TEXT ·Ycapget(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcapget(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ycapset(tls *TLS, a uintptr, b uintptr) (r int32)
TEXT ·Ycapset(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcapset(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ycarg(tls *TLS, z complex128) (r float64)
TEXT ·Ycarg(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcarg(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycargf(tls *TLS, z complex64) (r float32)
TEXT ·Ycargf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcargf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycargl(tls *TLS, z complex128) (r float64)
TEXT ·Ycargl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcargl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycasin(tls *TLS, z complex128) (r1 complex128)
TEXT ·Ycasin(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcasin(SB)
	MOVD 32(RSP), R0
	MOVD R0, r1_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r1_imag+32(FP)
	RET

// func Ycasinf(tls *TLS, z complex64) (r1 complex64)
TEXT ·Ycasinf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcasinf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r1_imag+20(FP)
	RET

// func Ycasinh(tls *TLS, z complex128) (r complex128)
TEXT ·Ycasinh(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcasinh(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycasinhf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycasinhf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcasinhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycasinhl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycasinhl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcasinhl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycasinl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycasinl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcasinl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycatan(tls *TLS, z complex128) (r complex128)
TEXT ·Ycatan(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcatan(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycatanf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycatanf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcatanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycatanh(tls *TLS, z complex128) (r complex128)
TEXT ·Ycatanh(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcatanh(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycatanhf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycatanhf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcatanhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycatanhl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycatanhl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcatanhl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycatanl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycatanl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcatanl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycatclose(tls *TLS, catd Tnl_catd) (r int32)
TEXT ·Ycatclose(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD catd+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcatclose(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycatgets(tls *TLS, catd Tnl_catd, set_id int32, msg_id int32, s uintptr) (r uintptr)
TEXT ·Ycatgets(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD catd+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW set_id+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW msg_id+20(FP), R0
	MOVW R0, 28(RSP)
	MOVD s+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xcatgets(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ycatopen(tls *TLS, name uintptr, oflag int32) (r Tnl_catd)
TEXT ·Ycatopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW oflag+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xcatopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycbrt(tls *TLS, x float64) (r1 float64)
TEXT ·Ycbrt(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcbrt(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Ycbrtf(tls *TLS, x float32) (r1 float32)
TEXT ·Ycbrtf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xcbrtf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Ycbrtl(tls *TLS, x float64) (r float64)
TEXT ·Ycbrtl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcbrtl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yccos(tls *TLS, z complex128) (r complex128)
TEXT ·Yccos(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xccos(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yccosf(tls *TLS, z complex64) (r complex64)
TEXT ·Yccosf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xccosf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Yccosh(tls *TLS, z complex128) (r complex128)
TEXT ·Yccosh(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xccosh(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yccoshf(tls *TLS, z complex64) (r complex64)
TEXT ·Yccoshf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xccoshf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Yccoshl(tls *TLS, z complex128) (r complex128)
TEXT ·Yccoshl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xccoshl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yccosl(tls *TLS, z complex128) (r complex128)
TEXT ·Yccosl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xccosl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yceil(tls *TLS, x3 float64) (r float64)
TEXT ·Yceil(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xceil(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yceilf(tls *TLS, x3 float32) (r float32)
TEXT ·Yceilf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xceilf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yceill(tls *TLS, x float64) (r float64)
TEXT ·Yceill(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xceill(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ycexp(tls *TLS, z complex128) (r complex128)
TEXT ·Ycexp(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcexp(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycexpf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycexpf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcexpf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycexpl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycexpl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcexpl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycfgetispeed(tls *TLS, tio uintptr) (r Tspeed_t)
TEXT ·Ycfgetispeed(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tio+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcfgetispeed(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycfgetospeed(tls *TLS, tio uintptr) (r Tspeed_t)
TEXT ·Ycfgetospeed(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tio+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcfgetospeed(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycfmakeraw(tls *TLS, t uintptr)
TEXT ·Ycfmakeraw(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcfmakeraw(SB)
	RET

// func Ycfsetispeed(tls *TLS, tio uintptr, speed Tspeed_t) (r int32)
TEXT ·Ycfsetispeed(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tio+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW speed+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xcfsetispeed(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ycfsetospeed(tls *TLS, tio uintptr, speed Tspeed_t) (r int32)
TEXT ·Ycfsetospeed(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tio+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW speed+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xcfsetospeed(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ycfsetspeed(tls *TLS, tio uintptr, speed Tspeed_t) (r int32)
TEXT ·Ycfsetspeed(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tio+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW speed+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xcfsetspeed(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ychdir(tls *TLS, path uintptr) (r int32)
TEXT ·Ychdir(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xchdir(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ychmod(tls *TLS, path uintptr, mode Tmode_t) (r int32)
TEXT ·Ychmod(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW mode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xchmod(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ychown(tls *TLS, path uintptr, uid Tuid_t, gid Tgid_t) (r int32)
TEXT ·Ychown(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW uid+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW gid+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xchown(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ychroot(tls *TLS, path uintptr) (r int32)
TEXT ·Ychroot(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xchroot(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycimag(tls *TLS, z complex128) (r float64)
TEXT ·Ycimag(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcimag(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycimagf(tls *TLS, z complex64) (r float32)
TEXT ·Ycimagf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcimagf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycimagl(tls *TLS, z complex128) (r float64)
TEXT ·Ycimagl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcimagl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yclearenv(tls *TLS) (r int32)
TEXT ·Yclearenv(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xclearenv(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yclearerr(tls *TLS, f uintptr)
TEXT ·Yclearerr(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xclearerr(SB)
	RET

// func Yclearerr_unlocked(tls *TLS, f uintptr)
TEXT ·Yclearerr_unlocked(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xclearerr_unlocked(SB)
	RET

// func Yclock(tls *TLS) (r Tclock_t)
TEXT ·Yclock(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xclock(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Yclock_adjtime(tls *TLS, clock_id Tclockid_t, utx uintptr) (r1 int32)
TEXT ·Yclock_adjtime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clock_id+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD utx+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xclock_adjtime(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Yclock_getcpuclockid(tls *TLS, pid Tpid_t, clk uintptr) (r int32)
TEXT ·Yclock_getcpuclockid(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD clk+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xclock_getcpuclockid(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yclock_getres(tls *TLS, clk Tclockid_t, ts uintptr) (r int32)
TEXT ·Yclock_getres(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clk+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD ts+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xclock_getres(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yclock_gettime(tls *TLS, clk Tclockid_t, ts uintptr) (r int32)
TEXT ·Yclock_gettime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clk+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD ts+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xclock_gettime(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yclock_nanosleep(tls *TLS, clk Tclockid_t, flags int32, req uintptr, rem uintptr) (r int32)
TEXT ·Yclock_nanosleep(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clk+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW flags+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD req+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD rem+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xclock_nanosleep(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yclock_settime(tls *TLS, clk Tclockid_t, ts uintptr) (r int32)
TEXT ·Yclock_settime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clk+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD ts+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xclock_settime(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yclog(tls *TLS, z complex128) (r1 complex128)
TEXT ·Yclog(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xclog(SB)
	MOVD 32(RSP), R0
	MOVD R0, r1_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r1_imag+32(FP)
	RET

// func Yclogf(tls *TLS, z complex64) (r1 complex64)
TEXT ·Yclogf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xclogf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r1_imag+20(FP)
	RET

// func Yclogl(tls *TLS, z complex128) (r complex128)
TEXT ·Yclogl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xclogl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yclose(tls *TLS, fd int32) (r1 int32)
TEXT ·Yclose(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xclose(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yclosedir(tls *TLS, dir uintptr) (r int32)
TEXT ·Yclosedir(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xclosedir(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycloselog(tls *TLS)
TEXT ·Ycloselog(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xcloselog(SB)
	RET

// func Yconfstr(tls *TLS, name int32, buf uintptr, len1 Tsize_t) (r Tsize_t)
TEXT ·Yconfstr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW name+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xconfstr(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yconj(tls *TLS, z complex128) (r complex128)
TEXT ·Yconj(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xconj(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yconjf(tls *TLS, z complex64) (r complex64)
TEXT ·Yconjf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xconjf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Yconjl(tls *TLS, z complex128) (r complex128)
TEXT ·Yconjl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xconjl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yconnect(tls *TLS, fd int32, addr uintptr, len1 Tsocklen_t) (r1 int32)
TEXT ·Yconnect(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW len1+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xconnect(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ycopy_file_range(tls *TLS, fd_in int32, off_in uintptr, fd_out int32, off_out uintptr, len1 Tsize_t, flags uint32) (r Tssize_t)
TEXT ·Ycopy_file_range(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd_in+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD off_in+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW fd_out+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD off_out+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD len1+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW flags+48(FP), R0
	MOVW R0, 56(RSP)
	CALL ·Xcopy_file_range(SB)
	MOVD 64(RSP), R0
	MOVD R0, r+56(FP)
	RET

// func Ycopysign(tls *TLS, x float64, y float64) (r float64)
TEXT ·Ycopysign(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcopysign(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycopysignf(tls *TLS, x float32, y float32) (r float32)
TEXT ·Ycopysignf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcopysignf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycopysignl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Ycopysignl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcopysignl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycos(tls *TLS, x3 float64) (r float64)
TEXT ·Ycos(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcos(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ycosf(tls *TLS, x3 float32) (r float32)
TEXT ·Ycosf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xcosf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycosh(tls *TLS, x3 float64) (r float64)
TEXT ·Ycosh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcosh(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ycoshf(tls *TLS, x3 float32) (r float32)
TEXT ·Ycoshf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xcoshf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycoshl(tls *TLS, x float64) (r float64)
TEXT ·Ycoshl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcoshl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ycosl(tls *TLS, x float64) (r float64)
TEXT ·Ycosl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcosl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ycpow(tls *TLS, z complex128, c complex128) (r complex128)
TEXT ·Ycpow(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD c_real+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD c_imag+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xcpow(SB)
	MOVD 48(RSP), R0
	MOVD R0, r_real+40(FP)
	MOVD 56(RSP), R0
	MOVD R0, r_imag+48(FP)
	RET

// func Ycpowf(tls *TLS, z complex64, c complex64) (r complex64)
TEXT ·Ycpowf(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW c_real+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW c_imag+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xcpowf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r_real+24(FP)
	MOVW 36(RSP), R0
	MOVW R0, r_imag+28(FP)
	RET

// func Ycpowl(tls *TLS, z complex128, c complex128) (r complex128)
TEXT ·Ycpowl(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD c_real+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD c_imag+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xcpowl(SB)
	MOVD 48(RSP), R0
	MOVD R0, r_real+40(FP)
	MOVD 56(RSP), R0
	MOVD R0, r_imag+48(FP)
	RET

// func Ycproj(tls *TLS, z complex128) (r complex128)
TEXT ·Ycproj(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcproj(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycprojf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycprojf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcprojf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycprojl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycprojl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcprojl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycreal(tls *TLS, z complex128) (r float64)
TEXT ·Ycreal(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcreal(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycrealf(tls *TLS, z complex64) (r float32)
TEXT ·Ycrealf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcrealf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ycreall(tls *TLS, z complex128) (r float64)
TEXT ·Ycreall(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcreall(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycreat(tls *TLS, filename uintptr, mode Tmode_t) (r int32)
TEXT ·Ycreat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW mode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xcreat(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ycrypt(tls *TLS, key uintptr, salt uintptr) (r uintptr)
TEXT ·Ycrypt(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD salt+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcrypt(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycrypt_r(tls *TLS, key uintptr, salt uintptr, data uintptr) (r uintptr)
TEXT ·Ycrypt_r(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD salt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD data+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xcrypt_r(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ycsin(tls *TLS, z complex128) (r complex128)
TEXT ·Ycsin(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcsin(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycsinf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycsinf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcsinf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycsinh(tls *TLS, z complex128) (r complex128)
TEXT ·Ycsinh(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcsinh(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycsinhf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycsinhf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcsinhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycsinhl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycsinhl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcsinhl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycsinl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycsinl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcsinl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycsqrt(tls *TLS, z complex128) (r complex128)
TEXT ·Ycsqrt(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcsqrt(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Ycsqrtf(tls *TLS, z complex64) (r complex64)
TEXT ·Ycsqrtf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcsqrtf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Ycsqrtl(tls *TLS, z complex128) (r complex128)
TEXT ·Ycsqrtl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcsqrtl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yctan(tls *TLS, z complex128) (r complex128)
TEXT ·Yctan(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xctan(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yctanf(tls *TLS, z complex64) (r complex64)
TEXT ·Yctanf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xctanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Yctanh(tls *TLS, z complex128) (r complex128)
TEXT ·Yctanh(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xctanh(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yctanhf(tls *TLS, z complex64) (r complex64)
TEXT ·Yctanhf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW z_real+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW z_imag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xctanhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_real+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_imag+20(FP)
	RET

// func Yctanhl(tls *TLS, z complex128) (r complex128)
TEXT ·Yctanhl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xctanhl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yctanl(tls *TLS, z complex128) (r complex128)
TEXT ·Yctanl(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD z_real+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD z_imag+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xctanl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_real+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_imag+32(FP)
	RET

// func Yctermid(tls *TLS, s uintptr) (r uintptr)
TEXT ·Yctermid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xctermid(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yctime(tls *TLS, t uintptr) (r uintptr)
TEXT ·Yctime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xctime(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yctime_r(tls *TLS, t uintptr, buf uintptr) (r uintptr)
TEXT ·Yctime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xctime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ycuserid(tls *TLS, buf uintptr) (r uintptr)
TEXT ·Ycuserid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcuserid(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ydcgettext(tls *TLS, domainname uintptr, msgid uintptr, category int32) (r uintptr)
TEXT ·Ydcgettext(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD domainname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD msgid+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW category+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xdcgettext(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ydcngettext(tls *TLS, domainname uintptr, msgid1 uintptr, msgid2 uintptr, n uint64, category int32) (r1 uintptr)
TEXT ·Ydcngettext(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD domainname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD msgid1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD msgid2+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD n+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW category+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xdcngettext(SB)
	MOVD 56(RSP), R0
	MOVD R0, r1+48(FP)
	RET

// func Ydelete_module(tls *TLS, a uintptr, b uint32) (r int32)
TEXT ·Ydelete_module(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW b+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xdelete_module(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ydgettext(tls *TLS, domainname uintptr, msgid uintptr) (r uintptr)
TEXT ·Ydgettext(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD domainname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD msgid+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xdgettext(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ydifftime(tls *TLS, t1 Ttime_t, t0 Ttime_t) (r float64)
TEXT ·Ydifftime(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD t0+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xdifftime(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ydirfd(tls *TLS, d uintptr) (r int32)
TEXT ·Ydirfd(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xdirfd(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ydirname(tls *TLS, s uintptr) (r uintptr)
TEXT ·Ydirname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xdirname(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ydiv(tls *TLS, num int32, den int32) (r Tdiv_t)
TEXT ·Ydiv(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW num+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW den+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xdiv(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_Fquot+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, r_Frem+20(FP)
	RET

// func Ydlclose(t *TLS, handle uintptr) (_2 int32)
TEXT ·Ydlclose(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD handle+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xdlclose(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ydlerror(t *TLS) (_1 uintptr)
TEXT ·Ydlerror(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xdlerror(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Ydlopen(t *TLS, filename uintptr, flags int32) (_3 uintptr)
TEXT ·Ydlopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xdlopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Ydlsym(t *TLS, handle, symbol uintptr) (_2 uintptr)
TEXT ·Ydlsym(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD handle+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD symbol+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xdlsym(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ydn_comp(tls *TLS, src uintptr, dst uintptr, space int32, dnptrs uintptr, lastdnptr uintptr) (r int32)
TEXT ·Ydn_comp(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD src+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD dst+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW space+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD dnptrs+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD lastdnptr+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xdn_comp(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ydn_expand(tls *TLS, base uintptr, end uintptr, src uintptr, dest uintptr, space int32) (r int32)
TEXT ·Ydn_expand(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD base+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD end+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD src+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD dest+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW space+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xdn_expand(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ydn_skipname(tls *TLS, s uintptr, end uintptr) (r int32)
TEXT ·Ydn_skipname(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD end+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xdn_skipname(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ydngettext(tls *TLS, domainname uintptr, msgid1 uintptr, msgid2 uintptr, n uint64) (r uintptr)
TEXT ·Ydngettext(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD domainname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD msgid1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD msgid2+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD n+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xdngettext(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ydprintf(tls *TLS, fd int32, fmt uintptr, va uintptr) (r int32)
TEXT ·Ydprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xdprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ydrand48(tls *TLS) (r float64)
TEXT ·Ydrand48(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xdrand48(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ydrem(tls *TLS, x float64, y float64) (r float64)
TEXT ·Ydrem(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xdrem(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ydremf(tls *TLS, x float32, y float32) (r float32)
TEXT ·Ydremf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xdremf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ydup(tls *TLS, fd int32) (r int32)
TEXT ·Ydup(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xdup(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ydup2(tls *TLS, old int32, new1 int32) (r1 int32)
TEXT ·Ydup2(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW old+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW new1+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xdup2(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Ydup3(tls *TLS, old int32, new1 int32, flags int32) (r int32)
TEXT ·Ydup3(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW old+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW new1+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xdup3(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yduplocale(tls *TLS, old Tlocale_t) (r Tlocale_t)
TEXT ·Yduplocale(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD old+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xduplocale(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yeaccess(tls *TLS, filename uintptr, amode int32) (r int32)
TEXT ·Yeaccess(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW amode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xeaccess(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yecvt(tls *TLS, x float64, n int32, dp uintptr, sign uintptr) (r uintptr)
TEXT ·Yecvt(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD dp+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD sign+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xecvt(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yencrypt(tls *TLS, block uintptr, edflag int32)
TEXT ·Yencrypt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD block+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW edflag+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xencrypt(SB)
	RET

// func Yendgrent(tls *TLS)
TEXT ·Yendgrent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendgrent(SB)
	RET

// func Yendhostent(tls *TLS)
TEXT ·Yendhostent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendhostent(SB)
	RET

// func Yendmntent(tls *TLS, f uintptr) (r int32)
TEXT ·Yendmntent(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xendmntent(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yendnetent(tls *TLS)
TEXT ·Yendnetent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendnetent(SB)
	RET

// func Yendprotoent(tls *TLS)
TEXT ·Yendprotoent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendprotoent(SB)
	RET

// func Yendpwent(tls *TLS)
TEXT ·Yendpwent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendpwent(SB)
	RET

// func Yendservent(tls *TLS)
TEXT ·Yendservent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendservent(SB)
	RET

// func Yendspent(tls *TLS)
TEXT ·Yendspent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendspent(SB)
	RET

// func Yendusershell(tls *TLS)
TEXT ·Yendusershell(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendusershell(SB)
	RET

// func Yendutent(tls *TLS)
TEXT ·Yendutent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendutent(SB)
	RET

// func Yendutxent(tls *TLS)
TEXT ·Yendutxent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendutxent(SB)
	RET

// func Yepoll_create(tls *TLS, size int32) (r int32)
TEXT ·Yepoll_create(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW size+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xepoll_create(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yepoll_create1(tls *TLS, flags int32) (r1 int32)
TEXT ·Yepoll_create1(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW flags+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xepoll_create1(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yepoll_ctl(tls *TLS, fd int32, op int32, fd2 int32, ev uintptr) (r int32)
TEXT ·Yepoll_ctl(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW op+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW fd2+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD ev+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xepoll_ctl(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yepoll_pwait(tls *TLS, fd int32, ev uintptr, cnt int32, to int32, sigs uintptr) (r1 int32)
TEXT ·Yepoll_pwait(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD ev+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW cnt+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW to+28(FP), R0
	MOVW R0, 36(RSP)
	MOVD sigs+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xepoll_pwait(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Yepoll_wait(tls *TLS, fd int32, ev uintptr, cnt int32, to int32) (r int32)
TEXT ·Yepoll_wait(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD ev+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW cnt+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW to+28(FP), R0
	MOVW R0, 36(RSP)
	CALL ·Xepoll_wait(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yerand48(tls *TLS, s uintptr) (r float64)
TEXT ·Yerand48(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xerand48(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yerf(tls *TLS, x float64) (r1 float64)
TEXT ·Yerf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xerf(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Yerfc(tls *TLS, x float64) (r1 float64)
TEXT ·Yerfc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xerfc(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Yerfcf(tls *TLS, x float32) (r1 float32)
TEXT ·Yerfcf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xerfcf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yerfcl(tls *TLS, x float64) (r float64)
TEXT ·Yerfcl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xerfcl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yerff(tls *TLS, x float32) (r1 float32)
TEXT ·Yerff(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xerff(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yerfl(tls *TLS, x float64) (r float64)
TEXT ·Yerfl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xerfl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yerr(tls *TLS, status int32, fmt uintptr, va uintptr)
TEXT ·Yerr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW status+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xerr(SB)
	RET

// func Yerrx(tls *TLS, status int32, fmt uintptr, va uintptr)
TEXT ·Yerrx(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW status+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xerrx(SB)
	RET

// func Yether_aton(tls *TLS, x uintptr) (r uintptr)
TEXT ·Yether_aton(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xether_aton(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yether_aton_r(tls *TLS, x uintptr, p_a uintptr) (r uintptr)
TEXT ·Yether_aton_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p_a+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xether_aton_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yether_hostton(tls *TLS, hostname uintptr, e uintptr) (r int32)
TEXT ·Yether_hostton(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD hostname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD e+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xether_hostton(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yether_line(tls *TLS, l uintptr, e uintptr, hostname uintptr) (r int32)
TEXT ·Yether_line(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD e+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD hostname+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xether_line(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yether_ntoa(tls *TLS, p_a uintptr) (r uintptr)
TEXT ·Yether_ntoa(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p_a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xether_ntoa(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yether_ntoa_r(tls *TLS, p_a uintptr, x uintptr) (r uintptr)
TEXT ·Yether_ntoa_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p_a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD x+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xether_ntoa_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yether_ntohost(tls *TLS, hostname uintptr, e uintptr) (r int32)
TEXT ·Yether_ntohost(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD hostname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD e+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xether_ntohost(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yeuidaccess(tls *TLS, filename uintptr, amode int32) (r int32)
TEXT ·Yeuidaccess(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW amode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xeuidaccess(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yeventfd(tls *TLS, count uint32, flags int32) (r1 int32)
TEXT ·Yeventfd(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW count+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW flags+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xeventfd(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yeventfd_read(tls *TLS, fd int32, value uintptr) (r int32)
TEXT ·Yeventfd_read(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD value+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xeventfd_read(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yeventfd_write(tls *TLS, fd int32, _value Teventfd_t) (r int32)
TEXT ·Yeventfd_write(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD _value+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xeventfd_write(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yexecl(tls *TLS, path uintptr, argv0 uintptr, va uintptr) (r int32)
TEXT ·Yexecl(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD argv0+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xexecl(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yexecle(tls *TLS, path uintptr, argv0 uintptr, va uintptr) (r int32)
TEXT ·Yexecle(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD argv0+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xexecle(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yexeclp(tls *TLS, file uintptr, argv0 uintptr, va uintptr) (r int32)
TEXT ·Yexeclp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD file+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD argv0+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xexeclp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yexecv(tls *TLS, path uintptr, argv uintptr) (r int32)
TEXT ·Yexecv(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xexecv(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yexecve(tls *TLS, path uintptr, argv uintptr, envp uintptr) (r int32)
TEXT ·Yexecve(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD envp+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xexecve(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yexecvp(tls *TLS, file uintptr, argv uintptr) (r int32)
TEXT ·Yexecvp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD file+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xexecvp(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yexecvpe(tls *TLS, file uintptr, argv uintptr, envp uintptr) (r int32)
TEXT ·Yexecvpe(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD file+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD envp+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xexecvpe(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yexit(tls *TLS, code int32)
TEXT ·Yexit(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW code+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xexit(SB)
	RET

// func Yexp(tls *TLS, x1 float64) (r1 float64)
TEXT ·Yexp(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexp(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Yexp10(tls *TLS, x float64) (r float64)
TEXT ·Yexp10(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexp10(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yexp10f(tls *TLS, x float32) (r float32)
TEXT ·Yexp10f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xexp10f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yexp10l(tls *TLS, x float64) (r float64)
TEXT ·Yexp10l(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexp10l(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yexp2(tls *TLS, x1 float64) (r1 float64)
TEXT ·Yexp2(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexp2(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Yexp2f(tls *TLS, x2 float32) (r1 float32)
TEXT ·Yexp2f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x2+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xexp2f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yexp2l(tls *TLS, x float64) (r float64)
TEXT ·Yexp2l(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexp2l(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yexpf(tls *TLS, x2 float32) (r1 float32)
TEXT ·Yexpf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x2+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xexpf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yexpl(tls *TLS, x float64) (r float64)
TEXT ·Yexpl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexpl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yexplicit_bzero(tls *TLS, d uintptr, n Tsize_t)
TEXT ·Yexplicit_bzero(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xexplicit_bzero(SB)
	RET

// func Yexpm1(tls *TLS, x3 float64) (r float64)
TEXT ·Yexpm1(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexpm1(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yexpm1f(tls *TLS, x3 float32) (r float32)
TEXT ·Yexpm1f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xexpm1f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yexpm1l(tls *TLS, x float64) (r float64)
TEXT ·Yexpm1l(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexpm1l(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfabs(tls *TLS, x float64) (r float64)
TEXT ·Yfabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfabsf(tls *TLS, x float32) (r float32)
TEXT ·Yfabsf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfabsf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfabsl(tls *TLS, x float64) (r float64)
TEXT ·Yfabsl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfabsl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfaccessat(tls *TLS, fd int32, filename uintptr, amode int32, flag int32) (r int32)
TEXT ·Yfaccessat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD filename+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW amode+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flag+28(FP), R0
	MOVW R0, 36(RSP)
	CALL ·Xfaccessat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfallocate(tls *TLS, fd int32, mode int32, base Toff_t, len1 Toff_t) (r int32)
TEXT ·Yfallocate(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW mode+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD base+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfallocate(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfanotify_init(tls *TLS, flags uint32, event_f_flags uint32) (r int32)
TEXT ·Yfanotify_init(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW flags+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW event_f_flags+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfanotify_init(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfanotify_mark(tls *TLS, fanotify_fd int32, flags uint32, mask uint64, dfd int32, pathname uintptr) (r int32)
TEXT ·Yfanotify_mark(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fanotify_fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW flags+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD mask+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW dfd+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD pathname+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfanotify_mark(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yfchdir(tls *TLS, fd int32) (r int32)
TEXT ·Yfchdir(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfchdir(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfchmod(tls *TLS, fd int32, mode Tmode_t) (r int32)
TEXT ·Yfchmod(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW mode+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfchmod(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfchmodat(tls *TLS, fd int32, path uintptr, mode Tmode_t, flag int32) (r int32)
TEXT ·Yfchmodat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mode+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flag+28(FP), R0
	MOVW R0, 36(RSP)
	CALL ·Xfchmodat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfchown(tls *TLS, fd int32, uid Tuid_t, gid Tgid_t) (r int32)
TEXT ·Yfchown(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW uid+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW gid+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xfchown(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfchownat(tls *TLS, fd int32, path uintptr, uid Tuid_t, gid Tgid_t, flag int32) (r int32)
TEXT ·Yfchownat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW uid+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW gid+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW flag+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xfchownat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yfclose(tls *TLS, f uintptr) (r1 int32)
TEXT ·Yfclose(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfclose(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yfcntl(tls *TLS, fd int32, cmd int32, va uintptr) (r int32)
TEXT ·Yfcntl(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW cmd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfcntl(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfcntl64(tls *TLS, fd int32, cmd int32, va uintptr) (r int32)
TEXT ·Yfcntl64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW cmd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfcntl64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfcvt(tls *TLS, x float64, n int32, dp uintptr, sign uintptr) (r uintptr)
TEXT ·Yfcvt(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD dp+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD sign+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfcvt(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yfdatasync(tls *TLS, fd int32) (r int32)
TEXT ·Yfdatasync(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfdatasync(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfdim(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yfdim(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfdim(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfdimf(tls *TLS, x float32, y float32) (r float32)
TEXT ·Yfdimf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfdimf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfdiml(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yfdiml(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfdiml(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfdopen(tls *TLS, fd int32, mode uintptr) (r uintptr)
TEXT ·Yfdopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfdopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfdopendir(tls *TLS, fd int32) (r uintptr)
TEXT ·Yfdopendir(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfdopendir(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfeclearexcept(tls *TLS, mask int32) (r int32)
TEXT ·Yfeclearexcept(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW mask+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfeclearexcept(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfegetenv(tls *TLS, envp uintptr) (r int32)
TEXT ·Yfegetenv(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD envp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfegetenv(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfegetround(tls *TLS) (r int32)
TEXT ·Yfegetround(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xfegetround(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yfeof(tls *TLS, f uintptr) (r int32)
TEXT ·Yfeof(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfeof(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfeof_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Yfeof_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfeof_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yferaiseexcept(tls *TLS, mask int32) (r int32)
TEXT ·Yferaiseexcept(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW mask+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xferaiseexcept(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yferror(tls *TLS, f uintptr) (r int32)
TEXT ·Yferror(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xferror(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yferror_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Yferror_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xferror_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfesetenv(tls *TLS, envp uintptr) (r int32)
TEXT ·Yfesetenv(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD envp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfesetenv(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfetestexcept(tls *TLS, mask int32) (r int32)
TEXT ·Yfetestexcept(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW mask+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfetestexcept(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfexecve(tls *TLS, fd int32, argv uintptr, envp uintptr) (r1 int32)
TEXT ·Yfexecve(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD envp+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfexecve(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Yfflush(tls *TLS, f uintptr) (r1 int32)
TEXT ·Yfflush(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfflush(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yfflush_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Yfflush_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfflush_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yffs(tls *TLS, i int32) (r int32)
TEXT ·Yffs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW i+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xffs(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yffsl(tls *TLS, i int64) (r int32)
TEXT ·Yffsl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD i+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xffsl(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yffsll(tls *TLS, i int64) (r int32)
TEXT ·Yffsll(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD i+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xffsll(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfgetc(tls *TLS, f1 uintptr) (r int32)
TEXT ·Yfgetc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfgetc(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfgetc_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Yfgetc_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfgetc_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfgetgrent(tls *TLS, f uintptr) (r uintptr)
TEXT ·Yfgetgrent(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfgetgrent(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfgetln(tls *TLS, f uintptr, plen uintptr) (r uintptr)
TEXT ·Yfgetln(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD plen+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfgetln(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfgetpos(tls *TLS, f uintptr, pos uintptr) (r int32)
TEXT ·Yfgetpos(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pos+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfgetpos(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfgetpwent(tls *TLS, f uintptr) (r uintptr)
TEXT ·Yfgetpwent(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfgetpwent(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfgets(tls *TLS, s uintptr, n int32, f uintptr) (r uintptr)
TEXT ·Yfgets(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfgets(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfgets_unlocked(tls *TLS, s uintptr, n int32, f uintptr) (r uintptr)
TEXT ·Yfgets_unlocked(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfgets_unlocked(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfgetwc(tls *TLS, f uintptr) (r Twint_t)
TEXT ·Yfgetwc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfgetwc(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfgetwc_unlocked(tls *TLS, f uintptr) (r Twint_t)
TEXT ·Yfgetwc_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfgetwc_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfgetws(tls *TLS, s uintptr, n int32, f uintptr) (r uintptr)
TEXT ·Yfgetws(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfgetws(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfgetws_unlocked(tls *TLS, s uintptr, n int32, f uintptr) (r uintptr)
TEXT ·Yfgetws_unlocked(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfgetws_unlocked(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfgetxattr(tls *TLS, filedes int32, name uintptr, value uintptr, size Tsize_t) (r Tssize_t)
TEXT ·Yfgetxattr(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW filedes+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD value+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfgetxattr(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yfileno(tls *TLS, f uintptr) (r int32)
TEXT ·Yfileno(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfileno(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfileno_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Yfileno_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfileno_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfinite(tls *TLS, x float64) (r int32)
TEXT ·Yfinite(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfinite(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfinitef(tls *TLS, x float32) (r int32)
TEXT ·Yfinitef(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfinitef(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yflistxattr(tls *TLS, filedes int32, list uintptr, size Tsize_t) (r Tssize_t)
TEXT ·Yflistxattr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW filedes+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD list+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xflistxattr(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yflock(tls *TLS, fd int32, op int32) (r int32)
TEXT ·Yflock(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW op+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xflock(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yflockfile(tls *TLS, f uintptr)
TEXT ·Yflockfile(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xflockfile(SB)
	RET

// func Yfloor(tls *TLS, x3 float64) (r float64)
TEXT ·Yfloor(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfloor(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfloorf(tls *TLS, x3 float32) (r float32)
TEXT ·Yfloorf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfloorf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfloorl(tls *TLS, x float64) (r float64)
TEXT ·Yfloorl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfloorl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfma(tls *TLS, x float64, y float64, z float64) (r float64)
TEXT ·Yfma(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD z+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfma(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfmaf(t *TLS, x, y, z float32) (r float32)
TEXT ·Yfmaf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW z+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xfmaf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfmal(tls *TLS, x float64, y float64, z float64) (r float64)
TEXT ·Yfmal(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD z+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfmal(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfmax(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yfmax(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfmax(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfmaxf(tls *TLS, x float32, y float32) (r float32)
TEXT ·Yfmaxf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfmaxf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfmaxl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yfmaxl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfmaxl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfmemopen(tls *TLS, buf uintptr, size Tsize_t, mode uintptr) (r uintptr)
TEXT ·Yfmemopen(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD mode+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfmemopen(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfmin(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yfmin(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfmin(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfminf(tls *TLS, x float32, y float32) (r float32)
TEXT ·Yfminf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfminf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfminl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yfminl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfminl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfmod(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yfmod(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfmod(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfmodf(tls *TLS, x float32, y float32) (r float32)
TEXT ·Yfmodf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfmodf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfmodl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yfmodl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfmodl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfmtmsg(tls *TLS, classification int64, label uintptr, severity int32, text uintptr, action uintptr, tag uintptr) (r int32)
TEXT ·Yfmtmsg(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD classification+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD label+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW severity+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD text+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD action+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD tag+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xfmtmsg(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Yfnmatch(tls *TLS, pat uintptr, str uintptr, flags int32) (r int32)
TEXT ·Yfnmatch(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pat+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD str+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xfnmatch(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfopen(tls *TLS, filename uintptr, mode uintptr) (r uintptr)
TEXT ·Yfopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfopen64(tls *TLS, filename uintptr, mode uintptr) (r uintptr)
TEXT ·Yfopen64(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfopen64(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfork(t *TLS) (_1 int32)
TEXT ·Yfork(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xfork(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Yfpathconf(tls *TLS, fd int32, name int32) (r int64)
TEXT ·Yfpathconf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW name+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfpathconf(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yfprintf(tls *TLS, f uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Yfprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfpurge(tls *TLS, f uintptr) (r int32)
TEXT ·Yfpurge(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfpurge(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfputc(tls *TLS, c1 int32, f1 uintptr) (r int32)
TEXT ·Yfputc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c1+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputc(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfputc_unlocked(tls *TLS, c int32, f uintptr) (r int32)
TEXT ·Yfputc_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputc_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfputs(tls *TLS, s uintptr, f uintptr) (r int32)
TEXT ·Yfputs(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputs(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfputs_unlocked(tls *TLS, s uintptr, f uintptr) (r int32)
TEXT ·Yfputs_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputs_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfputwc(tls *TLS, c Twchar_t, f uintptr) (r Twint_t)
TEXT ·Yfputwc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputwc(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfputwc_unlocked(tls *TLS, c Twchar_t, f uintptr) (r Twint_t)
TEXT ·Yfputwc_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputwc_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfputws(tls *TLS, _ws uintptr, f uintptr) (r int32)
TEXT ·Yfputws(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _ws+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputws(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfputws_unlocked(tls *TLS, _ws uintptr, f uintptr) (r int32)
TEXT ·Yfputws_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _ws+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputws_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfread(tls *TLS, destv uintptr, size Tsize_t, nmemb Tsize_t, f uintptr) (r Tsize_t)
TEXT ·Yfread(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD destv+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nmemb+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD f+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfread(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yfread_unlocked(tls *TLS, destv uintptr, size Tsize_t, nmemb Tsize_t, f uintptr) (r Tsize_t)
TEXT ·Yfread_unlocked(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD destv+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nmemb+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD f+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfread_unlocked(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yfree(tls *TLS, p uintptr)
TEXT ·Yfree(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfree(SB)
	RET

// func Yfreeaddrinfo(tls *TLS, p uintptr)
TEXT ·Yfreeaddrinfo(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfreeaddrinfo(SB)
	RET

// func Yfreeifaddrs(tls *TLS, ifp uintptr)
TEXT ·Yfreeifaddrs(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ifp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfreeifaddrs(SB)
	RET

// func Yfreelocale(tls *TLS, l Tlocale_t)
TEXT ·Yfreelocale(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfreelocale(SB)
	RET

// func Yfremovexattr(tls *TLS, fd int32, name uintptr) (r int32)
TEXT ·Yfremovexattr(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfremovexattr(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfreopen(tls *TLS, filename uintptr, mode uintptr, f uintptr) (r uintptr)
TEXT ·Yfreopen(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfreopen(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfrexp(tls *TLS, x float64, e uintptr) (r float64)
TEXT ·Yfrexp(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD e+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfrexp(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfrexpf(tls *TLS, x float32, e uintptr) (r float32)
TEXT ·Yfrexpf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD e+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfrexpf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfrexpl(tls *TLS, x float64, e uintptr) (r float64)
TEXT ·Yfrexpl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD e+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfrexpl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yfscanf(tls *TLS, f uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Yfscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfseek(tls *TLS, f uintptr, off int64, whence int32) (r int32)
TEXT ·Yfseek(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD off+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xfseek(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfseeko(tls *TLS, f uintptr, off Toff_t, whence int32) (r int32)
TEXT ·Yfseeko(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD off+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xfseeko(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfsetpos(tls *TLS, f uintptr, pos uintptr) (r int32)
TEXT ·Yfsetpos(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pos+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfsetpos(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfsetxattr(tls *TLS, filedes int32, name uintptr, value uintptr, size Tsize_t, flags int32) (r int32)
TEXT ·Yfsetxattr(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW filedes+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD value+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xfsetxattr(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Yfstat(tls *TLS, fd int32, st uintptr) (r int32)
TEXT ·Yfstat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD st+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfstat(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfstat64(tls *TLS, fd int32, st uintptr) (r int32)
TEXT ·Yfstat64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD st+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfstat64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfstatat(tls *TLS, fd int32, path uintptr, st uintptr, flag int32) (r int32)
TEXT ·Yfstatat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD st+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flag+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xfstatat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yfstatfs(tls *TLS, fd int32, buf uintptr) (r int32)
TEXT ·Yfstatfs(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfstatfs(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfstatvfs(tls *TLS, fd int32, buf uintptr) (r int32)
TEXT ·Yfstatvfs(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfstatvfs(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfsync(tls *TLS, fd int32) (r int32)
TEXT ·Yfsync(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfsync(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yftell(tls *TLS, f uintptr) (r int64)
TEXT ·Yftell(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xftell(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yftello(tls *TLS, f uintptr) (r Toff_t)
TEXT ·Yftello(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xftello(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yftime(tls *TLS, tp uintptr) (r int32)
TEXT ·Yftime(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xftime(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yftok(tls *TLS, path uintptr, id int32) (r Tkey_t)
TEXT ·Yftok(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW id+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xftok(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yftruncate(tls *TLS, fd int32, length Toff_t) (r int32)
TEXT ·Yftruncate(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD length+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xftruncate(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yftruncate64(tls *TLS, fd int32, length Toff_t) (r int32)
TEXT ·Yftruncate64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD length+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xftruncate64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yftrylockfile(tls *TLS, f uintptr) (r int32)
TEXT ·Yftrylockfile(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xftrylockfile(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yfts64_close(t *TLS, ftsp uintptr) (_2 int32)
TEXT ·Yfts64_close(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ftsp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfts64_close(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yfts64_open(t *TLS, path_argv uintptr, options int32, compar uintptr) (_4 uintptr)
TEXT ·Yfts64_open(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path_argv+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW options+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD compar+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfts64_open(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yfts64_read(t *TLS, ftsp uintptr) (_2 uintptr)
TEXT ·Yfts64_read(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ftsp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfts64_read(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yfts_close(t *TLS, ftsp uintptr) (_2 int32)
TEXT ·Yfts_close(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ftsp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfts_close(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yfts_open(t *TLS, path_argv uintptr, options int32, compar uintptr) (_4 uintptr)
TEXT ·Yfts_open(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path_argv+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW options+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD compar+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfts_open(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yfts_read(t *TLS, ftsp uintptr) (_2 uintptr)
TEXT ·Yfts_read(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ftsp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfts_read(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yftw(tls *TLS, path uintptr, __ccgo_fp_fn uintptr, fd_limit int32) (r int32)
TEXT ·Yftw(SB),$48-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 48(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 48(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_fn+16(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_ftw_1(SB)	// Create the closure for calling __ccgo_fp_fn
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD 48(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 24(RSP)
	MOVW fd_limit+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xftw(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_ftw_1(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _3+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD __ccgo_fp+32(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 40(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Yfunlockfile(tls *TLS, f uintptr)
TEXT ·Yfunlockfile(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfunlockfile(SB)
	RET

// func Yfutimens(tls *TLS, fd int32, times uintptr) (r int32)
TEXT ·Yfutimens(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD times+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfutimens(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfutimes(tls *TLS, fd int32, tv uintptr) (r int32)
TEXT ·Yfutimes(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD tv+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfutimes(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfutimesat(tls *TLS, dirfd int32, pathname uintptr, times uintptr) (r int32)
TEXT ·Yfutimesat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD times+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfutimesat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfwide(tls *TLS, f uintptr, mode int32) (r int32)
TEXT ·Yfwide(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW mode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xfwide(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfwprintf(tls *TLS, f uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Yfwprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfwprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yfwrite(tls *TLS, src uintptr, size Tsize_t, nmemb Tsize_t, f uintptr) (r Tsize_t)
TEXT ·Yfwrite(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD src+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nmemb+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD f+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfwrite(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yfwrite_unlocked(tls *TLS, src uintptr, size Tsize_t, nmemb Tsize_t, f uintptr) (r Tsize_t)
TEXT ·Yfwrite_unlocked(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD src+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nmemb+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD f+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfwrite_unlocked(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yfwscanf(tls *TLS, f uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Yfwscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfwscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ygai_strerror(tls *TLS, ecode int32) (r uintptr)
TEXT ·Ygai_strerror(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ecode+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgai_strerror(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygcvt(tls *TLS, x float64, n int32, b uintptr) (r uintptr)
TEXT ·Ygcvt(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD b+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgcvt(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yget_avphys_pages(tls *TLS) (r int64)
TEXT ·Yget_avphys_pages(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xget_avphys_pages(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Yget_current_dir_name(tls *TLS) (r uintptr)
TEXT ·Yget_current_dir_name(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xget_current_dir_name(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Yget_nprocs(tls *TLS) (r int32)
TEXT ·Yget_nprocs(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xget_nprocs(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yget_nprocs_conf(tls *TLS) (r int32)
TEXT ·Yget_nprocs_conf(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xget_nprocs_conf(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yget_phys_pages(tls *TLS) (r int64)
TEXT ·Yget_phys_pages(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xget_phys_pages(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetaddrinfo(tls *TLS, host uintptr, serv uintptr, hint uintptr, res uintptr) (r1 int32)
TEXT ·Ygetaddrinfo(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD host+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD serv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD hint+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD res+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xgetaddrinfo(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Ygetauxval(tls *TLS, item uint64) (r uint64)
TEXT ·Ygetauxval(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD item+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetauxval(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetc(tls *TLS, f1 uintptr) (r int32)
TEXT ·Ygetc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetc(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ygetc_unlocked(tls *TLS, f uintptr) (r int32)
TEXT ·Ygetc_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetc_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ygetchar(tls *TLS) (r int32)
TEXT ·Ygetchar(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetchar(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetchar_unlocked(tls *TLS) (r int32)
TEXT ·Ygetchar_unlocked(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetchar_unlocked(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetcwd(tls *TLS, buf uintptr, size Tsize_t) (r uintptr)
TEXT ·Ygetcwd(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetcwd(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ygetdate(tls *TLS, s uintptr) (r uintptr)
TEXT ·Ygetdate(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetdate(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetdelim(tls *TLS, s uintptr, n uintptr, delim int32, f uintptr) (r Tssize_t)
TEXT ·Ygetdelim(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW delim+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD f+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xgetdelim(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ygetdents(tls *TLS, fd int32, buf uintptr, len1 Tsize_t) (r int32)
TEXT ·Ygetdents(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetdents(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ygetdomainname(tls *TLS, name uintptr, len1 Tsize_t) (r int32)
TEXT ·Ygetdomainname(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetdomainname(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygetdtablesize(tls *TLS) (r int32)
TEXT ·Ygetdtablesize(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetdtablesize(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetegid(tls *TLS) (r Tgid_t)
TEXT ·Ygetegid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetegid(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetentropy(tls *TLS, buffer uintptr, len1 Tsize_t) (r int32)
TEXT ·Ygetentropy(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buffer+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetentropy(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygetenv(tls *TLS, name uintptr) (r uintptr)
TEXT ·Ygetenv(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetenv(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygeteuid(tls *TLS) (r Tuid_t)
TEXT ·Ygeteuid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgeteuid(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetgid(tls *TLS) (r Tgid_t)
TEXT ·Ygetgid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetgid(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetgrent(tls *TLS) (r uintptr)
TEXT ·Ygetgrent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetgrent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetgrgid(tls *TLS, gid Tgid_t) (r uintptr)
TEXT ·Ygetgrgid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW gid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgetgrgid(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetgrgid_r(tls *TLS, gid Tgid_t, gr uintptr, buf uintptr, size Tsize_t, res uintptr) (r int32)
TEXT ·Ygetgrgid_r(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW gid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD gr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD res+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetgrgid_r(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ygetgrnam(tls *TLS, name uintptr) (r uintptr)
TEXT ·Ygetgrnam(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetgrnam(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetgrnam_r(tls *TLS, name uintptr, gr uintptr, buf uintptr, size Tsize_t, res uintptr) (r int32)
TEXT ·Ygetgrnam_r(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD gr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD res+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetgrnam_r(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ygetgrouplist(tls *TLS, user uintptr, gid Tgid_t, groups uintptr, ngroups uintptr) (r int32)
TEXT ·Ygetgrouplist(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD user+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW gid+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD groups+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD ngroups+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xgetgrouplist(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ygetgroups(tls *TLS, count int32, list uintptr) (r int32)
TEXT ·Ygetgroups(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW count+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD list+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetgroups(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygethostbyaddr(tls *TLS, a uintptr, l Tsocklen_t, af int32) (r uintptr)
TEXT ·Ygethostbyaddr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW l+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW af+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xgethostbyaddr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ygethostbyaddr_r(tls *TLS, a uintptr, l Tsocklen_t, af int32, h uintptr, buf uintptr, buflen Tsize_t, res uintptr, err uintptr) (r int32)
TEXT ·Ygethostbyaddr_r(SB),$72-68
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW l+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW af+20(FP), R0
	MOVW R0, 28(RSP)
	MOVD h+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buf+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD buflen+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD res+48(FP), R0
	MOVD R0, 56(RSP)
	MOVD err+56(FP), R0
	MOVD R0, 64(RSP)
	CALL ·Xgethostbyaddr_r(SB)
	MOVW 72(RSP), R0
	MOVW R0, r+64(FP)
	RET

// func Ygethostbyname(tls *TLS, name uintptr) (r uintptr)
TEXT ·Ygethostbyname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgethostbyname(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygethostbyname2(tls *TLS, name uintptr, af int32) (r uintptr)
TEXT ·Ygethostbyname2(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW af+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xgethostbyname2(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ygethostbyname2_r(tls *TLS, name uintptr, af int32, h uintptr, buf uintptr, buflen Tsize_t, res uintptr, err uintptr) (r int32)
TEXT ·Ygethostbyname2_r(SB),$72-68
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW af+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD h+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buf+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD buflen+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD res+48(FP), R0
	MOVD R0, 56(RSP)
	MOVD err+56(FP), R0
	MOVD R0, 64(RSP)
	CALL ·Xgethostbyname2_r(SB)
	MOVW 72(RSP), R0
	MOVW R0, r+64(FP)
	RET

// func Ygethostbyname_r(tls *TLS, name uintptr, h uintptr, buf uintptr, buflen Tsize_t, res uintptr, err uintptr) (r int32)
TEXT ·Ygethostbyname_r(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD h+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buflen+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD res+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD err+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xgethostbyname_r(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Ygethostent(tls *TLS) (r uintptr)
TEXT ·Ygethostent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgethostent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygethostid(tls *TLS) (r int64)
TEXT ·Ygethostid(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgethostid(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygethostname(tls *TLS, name uintptr, len1 Tsize_t) (r int32)
TEXT ·Ygethostname(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgethostname(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygetifaddrs(tls *TLS, ifap uintptr) (r1 int32)
TEXT ·Ygetifaddrs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ifap+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetifaddrs(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Ygetitimer(tls *TLS, which int32, old uintptr) (r1 int32)
TEXT ·Ygetitimer(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW which+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD old+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetitimer(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ygetline(tls *TLS, s uintptr, n uintptr, f uintptr) (r Tssize_t)
TEXT ·Ygetline(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetline(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ygetloadavg(tls *TLS, a uintptr, n int32) (r int32)
TEXT ·Ygetloadavg(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xgetloadavg(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygetlogin(tls *TLS) (r uintptr)
TEXT ·Ygetlogin(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetlogin(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetlogin_r(tls *TLS, name uintptr, size Tsize_t) (r int32)
TEXT ·Ygetlogin_r(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetlogin_r(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygetmntent(tls *TLS, f uintptr) (r uintptr)
TEXT ·Ygetmntent(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetmntent(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetmntent_r(tls *TLS, f uintptr, mnt uintptr, linebuf uintptr, buflen int32) (r uintptr)
TEXT ·Ygetmntent_r(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mnt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD linebuf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW buflen+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xgetmntent_r(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ygetnameinfo(tls *TLS, sa uintptr, sl Tsocklen_t, node uintptr, nodelen Tsocklen_t, serv uintptr, servlen Tsocklen_t, flags int32) (r int32)
TEXT ·Ygetnameinfo(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sa+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW sl+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD node+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW nodelen+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD serv+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW servlen+48(FP), R0
	MOVW R0, 56(RSP)
	MOVW flags+52(FP), R0
	MOVW R0, 60(RSP)
	CALL ·Xgetnameinfo(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Ygetnetbyaddr(tls *TLS, net Tuint32_t, type1 int32) (r uintptr)
TEXT ·Ygetnetbyaddr(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW net+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW type1+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xgetnetbyaddr(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetnetbyname(tls *TLS, name uintptr) (r uintptr)
TEXT ·Ygetnetbyname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetnetbyname(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetnetent(tls *TLS) (r uintptr)
TEXT ·Ygetnetent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetnetent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetopt(tls *TLS, argc int32, argv uintptr, optstring uintptr) (r int32)
TEXT ·Ygetopt(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW argc+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD optstring+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetopt(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ygetopt_long(tls *TLS, argc int32, argv uintptr, optstring uintptr, longopts uintptr, idx uintptr) (r int32)
TEXT ·Ygetopt_long(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW argc+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD optstring+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD longopts+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD idx+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetopt_long(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ygetopt_long_only(tls *TLS, argc int32, argv uintptr, optstring uintptr, longopts uintptr, idx uintptr) (r int32)
TEXT ·Ygetopt_long_only(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW argc+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD argv+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD optstring+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD longopts+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD idx+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetopt_long_only(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ygetpagesize(tls *TLS) (r int32)
TEXT ·Ygetpagesize(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetpagesize(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetpass(tls *TLS, prompt uintptr) (r uintptr)
TEXT ·Ygetpass(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD prompt+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetpass(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetpeername(tls *TLS, fd int32, addr uintptr, len1 uintptr) (r1 int32)
TEXT ·Ygetpeername(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetpeername(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ygetpgid(tls *TLS, pid Tpid_t) (r Tpid_t)
TEXT ·Ygetpgid(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgetpgid(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ygetpgrp(tls *TLS) (r Tpid_t)
TEXT ·Ygetpgrp(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetpgrp(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetpid(tls *TLS) (r Tpid_t)
TEXT ·Ygetpid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetpid(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetppid(tls *TLS) (r Tpid_t)
TEXT ·Ygetppid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetppid(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetpriority(tls *TLS, which int32, who Tid_t) (r int32)
TEXT ·Ygetpriority(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW which+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW who+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xgetpriority(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ygetprotobyname(tls *TLS, name uintptr) (r uintptr)
TEXT ·Ygetprotobyname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetprotobyname(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetprotobynumber(tls *TLS, num int32) (r uintptr)
TEXT ·Ygetprotobynumber(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW num+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgetprotobynumber(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetprotoent(tls *TLS) (r uintptr)
TEXT ·Ygetprotoent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetprotoent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetpwent(tls *TLS) (r uintptr)
TEXT ·Ygetpwent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetpwent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetpwnam(tls *TLS, name uintptr) (r uintptr)
TEXT ·Ygetpwnam(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetpwnam(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetpwnam_r(tls *TLS, name uintptr, pw uintptr, buf uintptr, size Tsize_t, res uintptr) (r int32)
TEXT ·Ygetpwnam_r(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pw+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD res+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetpwnam_r(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ygetpwuid(tls *TLS, uid Tuid_t) (r uintptr)
TEXT ·Ygetpwuid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW uid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgetpwuid(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetpwuid_r(tls *TLS, uid Tuid_t, pw uintptr, buf uintptr, size Tsize_t, res uintptr) (r int32)
TEXT ·Ygetpwuid_r(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW uid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pw+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD res+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetpwuid_r(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ygetrandom(tls *TLS, buf uintptr, buflen Tsize_t, flags uint32) (r Tssize_t)
TEXT ·Ygetrandom(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buflen+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xgetrandom(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ygetresgid(tls *TLS, rgid uintptr, egid uintptr, sgid uintptr) (r int32)
TEXT ·Ygetresgid(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD rgid+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD egid+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD sgid+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetresgid(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ygetresuid(tls *TLS, ruid uintptr, euid uintptr, suid uintptr) (r int32)
TEXT ·Ygetresuid(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ruid+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD euid+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD suid+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetresuid(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ygetrlimit(tls *TLS, resource int32, rlim uintptr) (r int32)
TEXT ·Ygetrlimit(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW resource+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD rlim+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetrlimit(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygetrlimit64(tls *TLS, resource int32, rlim uintptr) (r int32)
TEXT ·Ygetrlimit64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW resource+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD rlim+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetrlimit64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygetrusage(tls *TLS, who int32, ru uintptr) (r1 int32)
TEXT ·Ygetrusage(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW who+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD ru+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetrusage(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ygets(tls *TLS, s uintptr) (r uintptr)
TEXT ·Ygets(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgets(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetservbyname(tls *TLS, name uintptr, prots uintptr) (r uintptr)
TEXT ·Ygetservbyname(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD prots+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetservbyname(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ygetservbyname_r(tls *TLS, name uintptr, prots uintptr, se uintptr, buf uintptr, buflen Tsize_t, res uintptr) (r int32)
TEXT ·Ygetservbyname_r(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD prots+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD se+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buf+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD buflen+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD res+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xgetservbyname_r(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Ygetservent(tls *TLS) (r uintptr)
TEXT ·Ygetservent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetservent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetsid(tls *TLS, pid Tpid_t) (r Tpid_t)
TEXT ·Ygetsid(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgetsid(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ygetsockname(tls *TLS, fd int32, addr uintptr, len1 uintptr) (r1 int32)
TEXT ·Ygetsockname(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetsockname(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ygetsockopt(tls *TLS, fd int32, level int32, optname int32, optval uintptr, optlen uintptr) (r2 int32)
TEXT ·Ygetsockopt(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW level+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW optname+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD optval+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD optlen+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xgetsockopt(SB)
	MOVW 48(RSP), R0
	MOVW R0, r2+40(FP)
	RET

// func Ygetspent(tls *TLS) (r uintptr)
TEXT ·Ygetspent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetspent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetsubopt(tls *TLS, opt uintptr, keys uintptr, val uintptr) (r int32)
TEXT ·Ygetsubopt(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD opt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD keys+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD val+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetsubopt(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ygettext(tls *TLS, msgid uintptr) (r uintptr)
TEXT ·Ygettext(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msgid+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgettext(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygettimeofday(tls *TLS, tv uintptr, tz uintptr) (r int32)
TEXT ·Ygettimeofday(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tv+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tz+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgettimeofday(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ygetuid(tls *TLS) (r Tuid_t)
TEXT ·Ygetuid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetuid(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetusershell(tls *TLS) (r uintptr)
TEXT ·Ygetusershell(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetusershell(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetutent(tls *TLS) (r uintptr)
TEXT ·Ygetutent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetutent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetutid(tls *TLS, ut uintptr) (r uintptr)
TEXT ·Ygetutid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ut+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetutid(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetutline(tls *TLS, ut uintptr) (r uintptr)
TEXT ·Ygetutline(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ut+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetutline(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetutxent(tls *TLS) (r uintptr)
TEXT ·Ygetutxent(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetutxent(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ygetutxid(tls *TLS, ut uintptr) (r uintptr)
TEXT ·Ygetutxid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ut+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetutxid(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetutxline(tls *TLS, ut uintptr) (r uintptr)
TEXT ·Ygetutxline(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ut+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetutxline(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygetw(tls *TLS, f uintptr) (r int32)
TEXT ·Ygetw(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetw(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ygetwc(tls *TLS, f uintptr) (r Twint_t)
TEXT ·Ygetwc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetwc(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ygetwc_unlocked(tls *TLS, f uintptr) (r Twint_t)
TEXT ·Ygetwc_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetwc_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ygetwchar(tls *TLS) (r Twint_t)
TEXT ·Ygetwchar(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetwchar(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetwchar_unlocked(tls *TLS) (r Twint_t)
TEXT ·Ygetwchar_unlocked(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetwchar_unlocked(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ygetxattr(tls *TLS, path uintptr, name uintptr, value uintptr, size Tsize_t) (r Tssize_t)
TEXT ·Ygetxattr(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD value+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xgetxattr(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yglob(tls *TLS, pat uintptr, flags int32, __ccgo_fp_errfunc uintptr, g_ uintptr) (r int32)
TEXT ·Yglob(SB),$56-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 56(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 56(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_errfunc+24(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_glob_2(SB)	// Create the closure for calling __ccgo_fp_errfunc
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pat+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD 56(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 32(RSP)
	MOVD g_+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xglob(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_glob_2(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW _2+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Yglobfree(tls *TLS, g_ uintptr)
TEXT ·Yglobfree(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD g_+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xglobfree(SB)
	RET

// func Ygmtime(tls *TLS, t uintptr) (r uintptr)
TEXT ·Ygmtime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgmtime(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ygmtime_r(tls *TLS, t uintptr, tm uintptr) (r uintptr)
TEXT ·Ygmtime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tm+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgmtime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ygrantpt(tls *TLS, fd int32) (r int32)
TEXT ·Ygrantpt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgrantpt(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yhasmntopt(tls *TLS, mnt uintptr, opt uintptr) (r uintptr)
TEXT ·Yhasmntopt(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD mnt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD opt+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xhasmntopt(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yhcreate(tls *TLS, nel Tsize_t) (r int32)
TEXT ·Yhcreate(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD nel+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xhcreate(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yhdestroy(tls *TLS)
TEXT ·Yhdestroy(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xhdestroy(SB)
	RET

// func Yherror(tls *TLS, msg uintptr)
TEXT ·Yherror(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msg+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xherror(SB)
	RET

// func Yhsearch(tls *TLS, item TENTRY, action TACTION) (r uintptr)
TEXT ·Yhsearch(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD item_Fkey+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD item_Fdata+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW action+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xhsearch(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yhstrerror(tls *TLS, ecode int32) (r uintptr)
TEXT ·Yhstrerror(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW ecode+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xhstrerror(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yhtonl(tls *TLS, n Tuint32_t) (r Tuint32_t)
TEXT ·Yhtonl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xhtonl(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yhtons(tls *TLS, n Tuint16_t) (r Tuint16_t)
TEXT ·Yhtons(SB),$24-18
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVH n+8(FP), R0
	MOVH R0, 16(RSP)
	CALL ·Xhtons(SB)
	MOVH 24(RSP), R0
	MOVH R0, r+16(FP)
	RET

// func Yhypot(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yhypot(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xhypot(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yhypotf(tls *TLS, x float32, y float32) (r float32)
TEXT ·Yhypotf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xhypotf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yhypotl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yhypotl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xhypotl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yiconv(tls *TLS, cd Ticonv_t, in uintptr, inb uintptr, out uintptr, outb uintptr) (_6 Tsize_t)
TEXT ·Yiconv(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD cd+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD in+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD inb+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD out+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD outb+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xiconv(SB)
	MOVD 56(RSP), R0
	MOVD R0, _6+48(FP)
	RET

// func Yiconv_close(tls *TLS, cd Ticonv_t) (r int32)
TEXT ·Yiconv_close(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD cd+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xiconv_close(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiconv_open(tls *TLS, to uintptr, from uintptr) (r Ticonv_t)
TEXT ·Yiconv_open(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD to+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD from+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiconv_open(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yif_freenameindex(tls *TLS, idx uintptr)
TEXT ·Yif_freenameindex(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD idx+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xif_freenameindex(SB)
	RET

// func Yif_indextoname(tls *TLS, index uint32, name uintptr) (r1 uintptr)
TEXT ·Yif_indextoname(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW index+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xif_indextoname(SB)
	MOVD 32(RSP), R0
	MOVD R0, r1+24(FP)
	RET

// func Yif_nameindex(tls *TLS) (r uintptr)
TEXT ·Yif_nameindex(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xif_nameindex(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Yif_nametoindex(tls *TLS, name uintptr) (r1 uint32)
TEXT ·Yif_nametoindex(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xif_nametoindex(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yilogb(tls *TLS, x3 float64) (r int32)
TEXT ·Yilogb(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xilogb(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yilogbf(tls *TLS, x3 float32) (r int32)
TEXT ·Yilogbf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xilogbf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yilogbl(tls *TLS, x float64) (r int32)
TEXT ·Yilogbl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xilogbl(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yimaxabs(tls *TLS, a Tintmax_t) (r Tintmax_t)
TEXT ·Yimaxabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Ximaxabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yimaxdiv(tls *TLS, num Tintmax_t, den Tintmax_t) (r Timaxdiv_t)
TEXT ·Yimaxdiv(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD num+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD den+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Ximaxdiv(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_Fquot+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_Frem+32(FP)
	RET

// func Yindex(tls *TLS, s uintptr, c int32) (r uintptr)
TEXT ·Yindex(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xindex(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yinet_addr(tls *TLS, p uintptr) (r Tin_addr_t)
TEXT ·Yinet_addr(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xinet_addr(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yinet_aton(tls *TLS, s0 uintptr, dest uintptr) (r int32)
TEXT ·Yinet_aton(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s0+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD dest+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xinet_aton(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yinet_lnaof(tls *TLS, in Tin_addr) (r Tin_addr_t)
TEXT ·Yinet_lnaof(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW in_Fs_addr+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xinet_lnaof(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yinet_makeaddr(tls *TLS, n Tin_addr_t, h Tin_addr_t) (r Tin_addr)
TEXT ·Yinet_makeaddr(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW h+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xinet_makeaddr(SB)
	MOVW 24(RSP), R0
	MOVW R0, r_Fs_addr+16(FP)
	RET

// func Yinet_netof(tls *TLS, in Tin_addr) (r Tin_addr_t)
TEXT ·Yinet_netof(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW in_Fs_addr+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xinet_netof(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yinet_network(tls *TLS, p uintptr) (r Tin_addr_t)
TEXT ·Yinet_network(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xinet_network(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yinet_ntoa(tls *TLS, in Tin_addr) (r uintptr)
TEXT ·Yinet_ntoa(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW in_Fs_addr+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xinet_ntoa(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yinet_ntop(tls *TLS, af int32, a0 uintptr, s uintptr, l Tsocklen_t) (r uintptr)
TEXT ·Yinet_ntop(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW af+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD a0+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD s+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW l+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xinet_ntop(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yinet_pton(tls *TLS, af int32, s uintptr, a0 uintptr) (r int32)
TEXT ·Yinet_pton(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW af+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD a0+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xinet_pton(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yinit_module(tls *TLS, a uintptr, b uint64, c uintptr) (r int32)
TEXT ·Yinit_module(SB),$40-36
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
	CALL ·Xinit_module(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yinitstate(tls *TLS, seed uint32, state uintptr, size Tsize_t) (r uintptr)
TEXT ·Yinitstate(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW seed+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD state+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xinitstate(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yinitstate_r(t *TLS, seed uint32, statebuf uintptr, statelen Tsize_t, buf uintptr) (_5 int32)
TEXT ·Yinitstate_r(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW seed+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD statebuf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD statelen+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buf+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xinitstate_r(SB)
	MOVW 48(RSP), R0
	MOVW R0, _5+40(FP)
	RET

// func Yinotify_add_watch(tls *TLS, fd int32, pathname uintptr, mask Tuint32_t) (r int32)
TEXT ·Yinotify_add_watch(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mask+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xinotify_add_watch(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yinotify_init(tls *TLS) (r int32)
TEXT ·Yinotify_init(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xinotify_init(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yinotify_init1(tls *TLS, flags int32) (r1 int32)
TEXT ·Yinotify_init1(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW flags+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xinotify_init1(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yinotify_rm_watch(tls *TLS, fd int32, wd int32) (r int32)
TEXT ·Yinotify_rm_watch(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW wd+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xinotify_rm_watch(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yinsque(tls *TLS, element uintptr, pred uintptr)
TEXT ·Yinsque(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD element+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pred+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xinsque(SB)
	RET

// func Yioctl(tls *TLS, fd int32, req int32, va uintptr) (r1 int32)
TEXT ·Yioctl(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW req+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xioctl(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Yisalnum(tls *TLS, c int32) (r int32)
TEXT ·Yisalnum(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisalnum(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisalnum_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisalnum_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisalnum_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yisalpha(tls *TLS, c int32) (r int32)
TEXT ·Yisalpha(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisalpha(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisalpha_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisalpha_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisalpha_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yisascii(tls *TLS, c int32) (r int32)
TEXT ·Yisascii(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisascii(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisastream(tls *TLS, fd int32) (r int32)
TEXT ·Yisastream(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisastream(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisatty(tls *TLS, fd int32) (r1 int32)
TEXT ·Yisatty(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisatty(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yisblank(tls *TLS, c int32) (r int32)
TEXT ·Yisblank(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisblank(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisblank_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisblank_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisblank_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiscntrl(tls *TLS, c int32) (r int32)
TEXT ·Yiscntrl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiscntrl(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiscntrl_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yiscntrl_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiscntrl_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yisdigit(tls *TLS, c int32) (r int32)
TEXT ·Yisdigit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisdigit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisdigit_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisdigit_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisdigit_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yisgraph(tls *TLS, c int32) (r int32)
TEXT ·Yisgraph(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisgraph(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisgraph_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisgraph_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisgraph_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yislower(tls *TLS, c int32) (r int32)
TEXT ·Yislower(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xislower(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yislower_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yislower_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xislower_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yisnan(t *TLS, x float64) (_2 int32)
TEXT ·Yisnan(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xisnan(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisnanf(t *TLS, arg float32) (_2 int32)
TEXT ·Yisnanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW arg+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisnanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisnanl(t *TLS, arg float64) (_2 int32)
TEXT ·Yisnanl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD arg+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xisnanl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisprint(tls *TLS, c int32) (r int32)
TEXT ·Yisprint(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisprint(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisprint_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisprint_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisprint_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yispunct(tls *TLS, c int32) (r int32)
TEXT ·Yispunct(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xispunct(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yispunct_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yispunct_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xispunct_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yissetugid(tls *TLS) (r int32)
TEXT ·Yissetugid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xissetugid(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yisspace(tls *TLS, c int32) (r int32)
TEXT ·Yisspace(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisspace(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisspace_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisspace_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisspace_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yisupper(tls *TLS, c int32) (r int32)
TEXT ·Yisupper(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisupper(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisupper_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisupper_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisupper_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswalnum(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswalnum(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswalnum(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswalnum_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswalnum_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswalnum_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswalpha(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswalpha(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswalpha(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswalpha_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswalpha_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswalpha_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswblank(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswblank(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswblank(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswblank_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswblank_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswblank_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswcntrl(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswcntrl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswcntrl(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswcntrl_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswcntrl_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswcntrl_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswctype(tls *TLS, wc Twint_t, type1 Twctype_t) (r int32)
TEXT ·Yiswctype(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD type1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswctype(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswctype_l(tls *TLS, c Twint_t, t Twctype_t, l Tlocale_t) (r int32)
TEXT ·Yiswctype_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD t+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xiswctype_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yiswdigit(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswdigit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswdigit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswdigit_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswdigit_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswdigit_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswgraph(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswgraph(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswgraph(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswgraph_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswgraph_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswgraph_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswlower(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswlower(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswlower(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswlower_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswlower_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswlower_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswprint(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswprint(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswprint(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswprint_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswprint_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswprint_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswpunct(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswpunct(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswpunct(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswpunct_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswpunct_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswpunct_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswspace(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswspace(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswspace(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswspace_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswspace_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswspace_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswupper(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswupper(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswupper(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswupper_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswupper_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswupper_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yiswxdigit(tls *TLS, wc Twint_t) (r int32)
TEXT ·Yiswxdigit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswxdigit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yiswxdigit_l(tls *TLS, c Twint_t, l Tlocale_t) (r int32)
TEXT ·Yiswxdigit_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xiswxdigit_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yisxdigit(tls *TLS, c int32) (r int32)
TEXT ·Yisxdigit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisxdigit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yisxdigit_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Yisxdigit_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xisxdigit_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yj0(tls *TLS, x float64) (r1 float64)
TEXT ·Yj0(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xj0(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Yj0f(tls *TLS, x float32) (r1 float32)
TEXT ·Yj0f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xj0f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yj1(tls *TLS, x float64) (r1 float64)
TEXT ·Yj1(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xj1(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Yj1f(tls *TLS, x float32) (r1 float32)
TEXT ·Yj1f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xj1f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yjn(tls *TLS, n int32, x float64) (r float64)
TEXT ·Yjn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD x+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xjn(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yjnf(tls *TLS, n int32, x float32) (r float32)
TEXT ·Yjnf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW x+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xjnf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yjrand48(tls *TLS, s uintptr) (r int64)
TEXT ·Yjrand48(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xjrand48(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ykill(tls *TLS, pid Tpid_t, sig int32) (r int32)
TEXT ·Ykill(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW sig+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xkill(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ykillpg(tls *TLS, pgid Tpid_t, sig int32) (r int32)
TEXT ·Ykillpg(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pgid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW sig+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xkillpg(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yklogctl(tls *TLS, type1 int32, buf uintptr, len1 int32) (r int32)
TEXT ·Yklogctl(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW type1+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW len1+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xklogctl(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yl64a(tls *TLS, x0 int64) (r uintptr)
TEXT ·Yl64a(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x0+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xl64a(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylabs(tls *TLS, a int64) (r int64)
TEXT ·Ylabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylchmod(tls *TLS, path uintptr, mode Tmode_t) (r int32)
TEXT ·Ylchmod(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW mode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xlchmod(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ylchown(tls *TLS, path uintptr, uid Tuid_t, gid Tgid_t) (r int32)
TEXT ·Ylchown(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW uid+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW gid+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xlchown(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ylckpwdf(tls *TLS) (r int32)
TEXT ·Ylckpwdf(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xlckpwdf(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ylcong48(tls *TLS, p uintptr)
TEXT ·Ylcong48(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlcong48(SB)
	RET

// func Yldexp(tls *TLS, x float64, n int32) (r float64)
TEXT ·Yldexp(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xldexp(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yldexpf(tls *TLS, x float32, n int32) (r float32)
TEXT ·Yldexpf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW n+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xldexpf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yldexpl(tls *TLS, x float64, n int32) (r float64)
TEXT ·Yldexpl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xldexpl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yldiv(tls *TLS, num int64, den int64) (r Tldiv_t)
TEXT ·Yldiv(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD num+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD den+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xldiv(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_Fquot+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_Frem+32(FP)
	RET

// func Ylfind(tls *TLS, key uintptr, base uintptr, nelp uintptr, width Tsize_t, __ccgo_fp_compar uintptr) (r uintptr)
TEXT ·Ylfind(SB),$64-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 64(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 64(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_compar+40(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_lfind_4(SB)	// Create the closure for calling __ccgo_fp_compar
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD base+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nelp+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD width+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD 64(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 48(RSP)
	CALL ·Xlfind(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_lfind_4(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ylgamma(tls *TLS, x float64) (r float64)
TEXT ·Ylgamma(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlgamma(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylgamma_r(tls *TLS, x float64, signgamp uintptr) (r float64)
TEXT ·Ylgamma_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD signgamp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlgamma_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ylgammaf(tls *TLS, x float32) (r float32)
TEXT ·Ylgammaf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlgammaf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ylgammaf_r(tls *TLS, x float32, signgamp uintptr) (r float32)
TEXT ·Ylgammaf_r(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD signgamp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlgammaf_r(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ylgammal(tls *TLS, x float64) (r float64)
TEXT ·Ylgammal(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlgammal(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylgammal_r(tls *TLS, x float64, sg uintptr) (r float64)
TEXT ·Ylgammal_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sg+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlgammal_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ylgetxattr(tls *TLS, path uintptr, name uintptr, value uintptr, size Tsize_t) (r Tssize_t)
TEXT ·Ylgetxattr(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD value+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xlgetxattr(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ylink(tls *TLS, existing uintptr, new1 uintptr) (r int32)
TEXT ·Ylink(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD existing+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD new1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlink(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ylinkat(tls *TLS, fd1 int32, existing uintptr, fd2 int32, new1 uintptr, flag int32) (r int32)
TEXT ·Ylinkat(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd1+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD existing+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW fd2+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD new1+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flag+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xlinkat(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ylisten(tls *TLS, fd int32, backlog int32) (r1 int32)
TEXT ·Ylisten(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW backlog+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xlisten(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Ylistxattr(tls *TLS, path uintptr, list uintptr, size Tsize_t) (r Tssize_t)
TEXT ·Ylistxattr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD list+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xlistxattr(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yllabs(tls *TLS, a int64) (r int64)
TEXT ·Yllabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xllabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylldiv(tls *TLS, num int64, den int64) (r Tlldiv_t)
TEXT ·Ylldiv(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD num+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD den+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlldiv(SB)
	MOVD 32(RSP), R0
	MOVD R0, r_Fquot+24(FP)
	MOVD 40(RSP), R0
	MOVD R0, r_Frem+32(FP)
	RET

// func Yllistxattr(tls *TLS, path uintptr, list uintptr, size Tsize_t) (r Tssize_t)
TEXT ·Yllistxattr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD list+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xllistxattr(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yllrint(tls *TLS, x float64) (r int64)
TEXT ·Yllrint(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xllrint(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yllrintf(tls *TLS, x float32) (r int64)
TEXT ·Yllrintf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xllrintf(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yllrintl(tls *TLS, x float64) (r int64)
TEXT ·Yllrintl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xllrintl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yllround(tls *TLS, x float64) (r int64)
TEXT ·Yllround(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xllround(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yllroundf(tls *TLS, x float32) (r int64)
TEXT ·Yllroundf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xllroundf(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yllroundl(tls *TLS, x float64) (r int64)
TEXT ·Yllroundl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xllroundl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylocaleconv(tls *TLS) (r uintptr)
TEXT ·Ylocaleconv(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xlocaleconv(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ylocaltime(tls *TLS, t uintptr) (r uintptr)
TEXT ·Ylocaltime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlocaltime(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylocaltime_r(tls *TLS, t uintptr, tm uintptr) (r uintptr)
TEXT ·Ylocaltime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tm+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlocaltime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ylockf(tls *TLS, fd int32, op int32, size Toff_t) (r int32)
TEXT ·Ylockf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW op+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlockf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ylog(tls *TLS, x1 float64) (r1 float64)
TEXT ·Ylog(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Ylog10(tls *TLS, x float64) (r float64)
TEXT ·Ylog10(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog10(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylog10f(tls *TLS, x float32) (r float32)
TEXT ·Ylog10f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlog10f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ylog10l(tls *TLS, x float64) (r float64)
TEXT ·Ylog10l(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog10l(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylog1p(tls *TLS, x3 float64) (r float64)
TEXT ·Ylog1p(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog1p(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylog1pf(tls *TLS, x3 float32) (r float32)
TEXT ·Ylog1pf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlog1pf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ylog1pl(tls *TLS, x float64) (r float64)
TEXT ·Ylog1pl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog1pl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylog2(tls *TLS, x1 float64) (r1 float64)
TEXT ·Ylog2(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog2(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Ylog2f(tls *TLS, x1 float32) (r1 float32)
TEXT ·Ylog2f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x1+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlog2f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Ylog2l(tls *TLS, x float64) (r float64)
TEXT ·Ylog2l(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog2l(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylogb(tls *TLS, x float64) (r float64)
TEXT ·Ylogb(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlogb(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylogbf(tls *TLS, x float32) (r float32)
TEXT ·Ylogbf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlogbf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ylogbl(tls *TLS, x float64) (r float64)
TEXT ·Ylogbl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlogbl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylogf(tls *TLS, x1 float32) (r1 float32)
TEXT ·Ylogf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x1+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlogf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Ylogin_tty(tls *TLS, fd int32) (r int32)
TEXT ·Ylogin_tty(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlogin_tty(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ylogl(tls *TLS, x float64) (r float64)
TEXT ·Ylogl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlogl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylongjmp(t *TLS, env uintptr, val int32)
TEXT ·Ylongjmp(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD env+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW val+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xlongjmp(SB)
	RET

// func Ylrand48(tls *TLS) (r int64)
TEXT ·Ylrand48(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xlrand48(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ylremovexattr(tls *TLS, path uintptr, name uintptr) (r int32)
TEXT ·Ylremovexattr(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlremovexattr(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ylrint(tls *TLS, x float64) (r int64)
TEXT ·Ylrint(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlrint(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylrintf(tls *TLS, x float32) (r int64)
TEXT ·Ylrintf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlrintf(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylrintl(tls *TLS, x float64) (r int64)
TEXT ·Ylrintl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlrintl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylround(tls *TLS, x float64) (r int64)
TEXT ·Ylround(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlround(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylroundf(tls *TLS, x float32) (r int64)
TEXT ·Ylroundf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlroundf(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylroundl(tls *TLS, x float64) (r int64)
TEXT ·Ylroundl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlroundl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ylsearch(tls *TLS, key uintptr, base uintptr, nelp uintptr, width Tsize_t, __ccgo_fp_compar uintptr) (r uintptr)
TEXT ·Ylsearch(SB),$64-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 64(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 64(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_compar+40(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_lsearch_4(SB)	// Create the closure for calling __ccgo_fp_compar
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD base+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nelp+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD width+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD 64(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 48(RSP)
	CALL ·Xlsearch(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_lsearch_4(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ylseek(tls *TLS, fd int32, offset Toff_t, whence int32) (r Toff_t)
TEXT ·Ylseek(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD offset+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xlseek(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ylseek64(tls *TLS, fd int32, offset Toff_t, whence int32) (r Toff_t)
TEXT ·Ylseek64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD offset+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xlseek64(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ylsetxattr(tls *TLS, path uintptr, name uintptr, value uintptr, size Tsize_t, flags int32) (r int32)
TEXT ·Ylsetxattr(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD value+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xlsetxattr(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ylstat(tls *TLS, path uintptr, buf uintptr) (r int32)
TEXT ·Ylstat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlstat(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ylstat64(tls *TLS, path uintptr, buf uintptr) (r int32)
TEXT ·Ylstat64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlstat64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ylutimes(tls *TLS, filename uintptr, tv uintptr) (r int32)
TEXT ·Ylutimes(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tv+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlutimes(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymadvise(tls *TLS, addr uintptr, len1 Tsize_t, advice int32) (r int32)
TEXT ·Ymadvise(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW advice+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xmadvise(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ymalloc(tls *TLS, n Tsize_t) (r uintptr)
TEXT ·Ymalloc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmalloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ymalloc_usable_size(tls *TLS, p uintptr) (r Tsize_t)
TEXT ·Ymalloc_usable_size(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmalloc_usable_size(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ymblen(tls *TLS, s uintptr, n Tsize_t) (r int32)
TEXT ·Ymblen(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmblen(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymbrlen(tls *TLS, s uintptr, n Tsize_t, st uintptr) (r Tsize_t)
TEXT ·Ymbrlen(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD st+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmbrlen(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymbrtoc16(tls *TLS, pc16 uintptr, s uintptr, n Tsize_t, ps uintptr) (r Tsize_t)
TEXT ·Ymbrtoc16(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pc16+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD ps+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xmbrtoc16(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ymbrtoc32(tls *TLS, pc32 uintptr, s uintptr, n Tsize_t, ps uintptr) (r Tsize_t)
TEXT ·Ymbrtoc32(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pc32+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD ps+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xmbrtoc32(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ymbrtowc(tls *TLS, wc uintptr, src uintptr, n Tsize_t, st uintptr) (_5 Tsize_t)
TEXT ·Ymbrtowc(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD wc+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD st+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xmbrtowc(SB)
	MOVD 48(RSP), R0
	MOVD R0, _5+40(FP)
	RET

// func Ymbsinit(tls *TLS, st uintptr) (r int32)
TEXT ·Ymbsinit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD st+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmbsinit(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ymbsnrtowcs(tls *TLS, wcs uintptr, src uintptr, n Tsize_t, wn Tsize_t, st uintptr) (r Tsize_t)
TEXT ·Ymbsnrtowcs(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD wcs+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD wn+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD st+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xmbsnrtowcs(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ymbsrtowcs(tls *TLS, ws uintptr, src uintptr, wn Tsize_t, st uintptr) (r Tsize_t)
TEXT ·Ymbsrtowcs(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ws+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD wn+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD st+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xmbsrtowcs(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ymbstowcs(tls *TLS, ws uintptr, _s uintptr, wn Tsize_t) (r Tsize_t)
TEXT ·Ymbstowcs(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ws+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD wn+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmbstowcs(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymbtowc(tls *TLS, wc uintptr, src uintptr, n Tsize_t) (_4 int32)
TEXT ·Ymbtowc(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD wc+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmbtowc(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Ymemccpy(tls *TLS, dest uintptr, src uintptr, c int32, n Tsize_t) (r uintptr)
TEXT ·Ymemccpy(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW c+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD n+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xmemccpy(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ymemchr(tls *TLS, src uintptr, c int32, n Tsize_t) (r uintptr)
TEXT ·Ymemchr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD src+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmemchr(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymemcmp(tls *TLS, vl uintptr, vr uintptr, n Tsize_t) (r1 int32)
TEXT ·Ymemcmp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD vl+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD vr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmemcmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ymemcpy(tls *TLS, dest uintptr, src uintptr, n Tsize_t) (r uintptr)
TEXT ·Ymemcpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmemcpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymemfd_create(tls *TLS, name uintptr, flags uint32) (r int32)
TEXT ·Ymemfd_create(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmemfd_create(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymemmem(tls *TLS, h0 uintptr, k Tsize_t, n0 uintptr, l Tsize_t) (r uintptr)
TEXT ·Ymemmem(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD h0+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD k+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n0+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD l+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xmemmem(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ymemmove(tls *TLS, dest uintptr, src uintptr, n Tsize_t) (r uintptr)
TEXT ·Ymemmove(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmemmove(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymempcpy(tls *TLS, dest uintptr, src uintptr, n Tsize_t) (r uintptr)
TEXT ·Ymempcpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmempcpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymemrchr(tls *TLS, m uintptr, c int32, n Tsize_t) (r uintptr)
TEXT ·Ymemrchr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmemrchr(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymemset(tls *TLS, dest uintptr, c int32, n Tsize_t) (r uintptr)
TEXT ·Ymemset(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmemset(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymincore(tls *TLS, addr uintptr, len1 Tsize_t, vec uintptr) (r int32)
TEXT ·Ymincore(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD vec+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmincore(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ymkdir(tls *TLS, path uintptr, mode Tmode_t) (r int32)
TEXT ·Ymkdir(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW mode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmkdir(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymkdirat(tls *TLS, fd int32, path uintptr, mode Tmode_t) (r int32)
TEXT ·Ymkdirat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mode+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xmkdirat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ymkdtemp(tls *TLS, template uintptr) (r uintptr)
TEXT ·Ymkdtemp(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmkdtemp(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ymkfifo(tls *TLS, path uintptr, mode Tmode_t) (r int32)
TEXT ·Ymkfifo(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW mode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmkfifo(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymkfifoat(tls *TLS, fd int32, path uintptr, mode Tmode_t) (r int32)
TEXT ·Ymkfifoat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mode+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xmkfifoat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ymknod(tls *TLS, path uintptr, mode Tmode_t, dev Tdev_t) (r int32)
TEXT ·Ymknod(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW mode+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD dev+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmknod(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ymknodat(tls *TLS, fd int32, path uintptr, mode Tmode_t, dev Tdev_t) (r int32)
TEXT ·Ymknodat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mode+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD dev+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xmknodat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ymkostemp(tls *TLS, template uintptr, flags int32) (r int32)
TEXT ·Ymkostemp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmkostemp(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymkostemps(tls *TLS, template uintptr, len1 int32, flags int32) (r int32)
TEXT ·Ymkostemps(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW len1+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW flags+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xmkostemps(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymkstemp(tls *TLS, template uintptr) (r int32)
TEXT ·Ymkstemp(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmkstemp(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ymkstemp64(tls *TLS, template uintptr) (r int32)
TEXT ·Ymkstemp64(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmkstemp64(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ymkstemps(tls *TLS, template uintptr, len1 int32) (r int32)
TEXT ·Ymkstemps(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW len1+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmkstemps(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymkstemps64(tls *TLS, template uintptr, len1 int32) (r int32)
TEXT ·Ymkstemps64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW len1+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmkstemps64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymktemp(tls *TLS, template uintptr) (r uintptr)
TEXT ·Ymktemp(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmktemp(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ymktime(tls *TLS, tm uintptr) (r Ttime_t)
TEXT ·Ymktime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tm+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmktime(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ymlock(tls *TLS, addr uintptr, len1 Tsize_t) (r int32)
TEXT ·Ymlock(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmlock(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymlock2(tls *TLS, addr uintptr, len1 Tsize_t, flags uint32) (r int32)
TEXT ·Ymlock2(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xmlock2(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ymlockall(tls *TLS, flags int32) (r int32)
TEXT ·Ymlockall(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW flags+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xmlockall(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ymmap(tls *TLS, start uintptr, len1 Tsize_t, prot int32, flags int32, fd int32, off Toff_t) (r uintptr)
TEXT ·Ymmap(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD start+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW prot+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW fd+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD off+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xmmap(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ymmap64(tls *TLS, start uintptr, len1 Tsize_t, prot int32, flags int32, fd int32, off Toff_t) (r uintptr)
TEXT ·Ymmap64(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD start+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW prot+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW fd+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD off+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xmmap64(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ymodf(tls *TLS, x float64, iptr uintptr) (r float64)
TEXT ·Ymodf(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD iptr+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmodf(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ymodff(tls *TLS, x float32, iptr uintptr) (r float32)
TEXT ·Ymodff(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iptr+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmodff(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymodfl(tls *TLS, x float64, iptr uintptr) (r1 float64)
TEXT ·Ymodfl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD iptr+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmodfl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r1+24(FP)
	RET

// func Ymount(tls *TLS, special uintptr, dir uintptr, fstype uintptr, flags uint64, data uintptr) (r int32)
TEXT ·Ymount(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD special+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD dir+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD fstype+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD flags+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD data+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xmount(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ymprotect(tls *TLS, addr uintptr, len1 Tsize_t, prot int32) (r int32)
TEXT ·Ymprotect(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW prot+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xmprotect(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ymrand48(tls *TLS) (r int64)
TEXT ·Ymrand48(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xmrand48(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ymremap(tls *TLS, old_addr uintptr, old_len Tsize_t, new_len Tsize_t, flags int32, va uintptr) (r uintptr)
TEXT ·Ymremap(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD old_addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD old_len+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD new_len+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD va+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xmremap(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ymsgctl(tls *TLS, q int32, cmd int32, buf uintptr) (r1 int32)
TEXT ·Ymsgctl(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW q+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW cmd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmsgctl(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ymsgget(tls *TLS, k Tkey_t, flag int32) (r int32)
TEXT ·Ymsgget(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW k+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW flag+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xmsgget(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ymsgrcv(tls *TLS, q int32, m uintptr, len1 Tsize_t, type1 int64, flag int32) (r Tssize_t)
TEXT ·Ymsgrcv(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW q+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD m+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD type1+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flag+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xmsgrcv(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ymsgsnd(tls *TLS, q int32, m uintptr, len1 Tsize_t, flag int32) (r int32)
TEXT ·Ymsgsnd(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW q+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD m+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flag+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xmsgsnd(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ymsync(tls *TLS, start uintptr, len1 Tsize_t, flags int32) (r int32)
TEXT ·Ymsync(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD start+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xmsync(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ymunlock(tls *TLS, addr uintptr, len1 Tsize_t) (r int32)
TEXT ·Ymunlock(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmunlock(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ymunlockall(tls *TLS) (r int32)
TEXT ·Ymunlockall(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xmunlockall(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ymunmap(tls *TLS, start uintptr, len1 Tsize_t) (r int32)
TEXT ·Ymunmap(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD start+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmunmap(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yname_to_handle_at(tls *TLS, dirfd int32, pathname uintptr, handle uintptr, mount_id uintptr, flags int32) (r int32)
TEXT ·Yname_to_handle_at(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD handle+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD mount_id+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xname_to_handle_at(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ynan(tls *TLS, s uintptr) (r float64)
TEXT ·Ynan(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xnan(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ynanf(tls *TLS, s uintptr) (r float32)
TEXT ·Ynanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xnanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ynanl(tls *TLS, s uintptr) (r float64)
TEXT ·Ynanl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xnanl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ynanosleep(tls *TLS, req uintptr, rem uintptr) (r int32)
TEXT ·Ynanosleep(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD req+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD rem+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnanosleep(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ynewlocale(tls *TLS, mask int32, name uintptr, loc Tlocale_t) (r Tlocale_t)
TEXT ·Ynewlocale(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW mask+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD loc+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xnewlocale(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ynextafter(tls *TLS, x3 float64, y3 float64) (r float64)
TEXT ·Ynextafter(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y3+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnextafter(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ynextafterf(tls *TLS, x3 float32, y3 float32) (r float32)
TEXT ·Ynextafterf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y3+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xnextafterf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ynextafterl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Ynextafterl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnextafterl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ynexttoward(tls *TLS, x float64, y float64) (r float64)
TEXT ·Ynexttoward(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnexttoward(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ynexttowardf(tls *TLS, x3 float32, y3 float64) (r float32)
TEXT ·Ynexttowardf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD y3+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnexttowardf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ynexttowardl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Ynexttowardl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnexttowardl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ynftw(tls *TLS, path uintptr, __ccgo_fp_fn uintptr, fd_limit int32, flags int32) (r1 int32)
TEXT ·Ynftw(SB),$48-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 48(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 48(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_fn+16(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_nftw_1(SB)	// Create the closure for calling __ccgo_fp_fn
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD 48(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 24(RSP)
	MOVW fd_limit+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	CALL ·Xnftw(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_nftw_1(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _3+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD _4+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD __ccgo_fp+40(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 48(RSP), R0
	MOVW R0, _5+48(FP)
	RET

// func Yngettext(tls *TLS, msgid1 uintptr, msgid2 uintptr, n uint64) (r uintptr)
TEXT ·Yngettext(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msgid1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD msgid2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xngettext(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ynice(tls *TLS, inc int32) (r int32)
TEXT ·Ynice(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW inc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xnice(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ynl_langinfo(tls *TLS, item Tnl_item) (r uintptr)
TEXT ·Ynl_langinfo(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW item+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xnl_langinfo(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ynl_langinfo_l(tls *TLS, item Tnl_item, loc Tlocale_t) (r uintptr)
TEXT ·Ynl_langinfo_l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW item+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD loc+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnl_langinfo_l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ynrand48(tls *TLS, s uintptr) (r int64)
TEXT ·Ynrand48(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xnrand48(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yns_get16(tls *TLS, cp uintptr) (r uint32)
TEXT ·Yns_get16(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD cp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xns_get16(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yns_get32(tls *TLS, cp uintptr) (r uint64)
TEXT ·Yns_get32(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD cp+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xns_get32(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yns_initparse(tls *TLS, msg uintptr, msglen int32, handle uintptr) (r1 int32)
TEXT ·Yns_initparse(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW msglen+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD handle+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xns_initparse(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Yns_name_uncompress(tls *TLS, msg uintptr, eom uintptr, src uintptr, dst uintptr, dstsiz Tsize_t) (r1 int32)
TEXT ·Yns_name_uncompress(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD eom+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD src+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD dst+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD dstsiz+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xns_name_uncompress(SB)
	MOVW 56(RSP), R0
	MOVW R0, r1+48(FP)
	RET

// func Yns_parserr(tls *TLS, handle uintptr, section Tns_sect, rrnum int32, rr uintptr) (r1 int32)
TEXT ·Yns_parserr(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD handle+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW section+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW rrnum+20(FP), R0
	MOVW R0, 28(RSP)
	MOVD rr+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xns_parserr(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Yns_put16(tls *TLS, s uint32, cp uintptr)
TEXT ·Yns_put16(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW s+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD cp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xns_put16(SB)
	RET

// func Yns_put32(tls *TLS, l uint64, cp uintptr)
TEXT ·Yns_put32(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD cp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xns_put32(SB)
	RET

// func Yns_skiprr(tls *TLS, ptr uintptr, eom uintptr, section Tns_sect, count int32) (r1 int32)
TEXT ·Yns_skiprr(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD eom+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW section+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW count+28(FP), R0
	MOVW R0, 36(RSP)
	CALL ·Xns_skiprr(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Yntohl(tls *TLS, n Tuint32_t) (r Tuint32_t)
TEXT ·Yntohl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xntohl(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yntohs(tls *TLS, n Tuint16_t) (r Tuint16_t)
TEXT ·Yntohs(SB),$24-18
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVH n+8(FP), R0
	MOVH R0, 16(RSP)
	CALL ·Xntohs(SB)
	MOVH 24(RSP), R0
	MOVH R0, r+16(FP)
	RET

// func Yobstack_free(t *TLS, obstack, obj uintptr)
TEXT ·Yobstack_free(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD obstack+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD obj+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xobstack_free(SB)
	RET

// func Yobstack_vprintf(t *TLS, obstack, template, va uintptr) (_2 int32)
TEXT ·Yobstack_vprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD obstack+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD template+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xobstack_vprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Yopen(tls *TLS, filename uintptr, flags int32, va uintptr) (r int32)
TEXT ·Yopen(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xopen(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yopen64(tls *TLS, filename uintptr, flags int32, va uintptr) (r int32)
TEXT ·Yopen64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xopen64(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yopen_by_handle_at(tls *TLS, mount_fd int32, handle uintptr, flags int32) (r int32)
TEXT ·Yopen_by_handle_at(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW mount_fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD handle+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xopen_by_handle_at(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yopen_memstream(tls *TLS, bufp uintptr, sizep uintptr) (r uintptr)
TEXT ·Yopen_memstream(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD bufp+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sizep+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xopen_memstream(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yopen_wmemstream(tls *TLS, bufp uintptr, sizep uintptr) (r uintptr)
TEXT ·Yopen_wmemstream(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD bufp+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sizep+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xopen_wmemstream(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yopenat(tls *TLS, fd int32, filename uintptr, flags int32, va uintptr) (r int32)
TEXT ·Yopenat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD filename+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD va+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xopenat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yopendir(tls *TLS, name uintptr) (r uintptr)
TEXT ·Yopendir(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xopendir(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yopenlog(tls *TLS, ident uintptr, opt int32, facility int32)
TEXT ·Yopenlog(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ident+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW opt+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW facility+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xopenlog(SB)
	RET

// func Yopenpty(tls *TLS, pm uintptr, ps uintptr, name uintptr, tio uintptr, ws uintptr) (r int32)
TEXT ·Yopenpty(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pm+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ps+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD name+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD tio+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD ws+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xopenpty(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ypathconf(tls *TLS, path uintptr, name int32) (r int64)
TEXT ·Ypathconf(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW name+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpathconf(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ypause(tls *TLS) (r int32)
TEXT ·Ypause(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xpause(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ypclose(tls *TLS, f uintptr) (r1 int32)
TEXT ·Ypclose(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpclose(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yperror(tls *TLS, msg uintptr)
TEXT ·Yperror(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD msg+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xperror(SB)
	RET

// func Ypersonality(tls *TLS, persona uint64) (r int32)
TEXT ·Ypersonality(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD persona+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpersonality(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ypipe(tls *TLS, fd uintptr) (r int32)
TEXT ·Ypipe(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fd+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpipe(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ypipe2(tls *TLS, fd uintptr, flag int32) (r int32)
TEXT ·Ypipe2(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fd+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flag+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpipe2(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ypivot_root(tls *TLS, new1 uintptr, old uintptr) (r int32)
TEXT ·Ypivot_root(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD new1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD old+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpivot_root(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ypoll(tls *TLS, fds uintptr, n Tnfds_t, timeout int32) (r int32)
TEXT ·Ypoll(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fds+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW timeout+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xpoll(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ypopen(t *TLS, command, type1 uintptr) (_2 uintptr)
TEXT ·Ypopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD command+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD type1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Yposix_close(tls *TLS, fd int32, flags int32) (r int32)
TEXT ·Yposix_close(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW flags+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xposix_close(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yposix_fadvise(tls *TLS, fd int32, base Toff_t, len1 Toff_t, advice int32) (r int32)
TEXT ·Yposix_fadvise(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD base+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW advice+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xposix_fadvise(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yposix_fallocate(tls *TLS, fd int32, base Toff_t, len1 Toff_t) (r int32)
TEXT ·Yposix_fallocate(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD base+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xposix_fallocate(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yposix_madvise(tls *TLS, addr uintptr, len1 Tsize_t, advice int32) (r int32)
TEXT ·Yposix_madvise(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW advice+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xposix_madvise(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yposix_openpt(tls *TLS, flags int32) (r1 int32)
TEXT ·Yposix_openpt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW flags+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xposix_openpt(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yposix_spawn_file_actions_addchdir_np(tls *TLS, fa uintptr, path uintptr) (r int32)
TEXT ·Yposix_spawn_file_actions_addchdir_np(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fa+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawn_file_actions_addchdir_np(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawn_file_actions_addclose(tls *TLS, fa uintptr, fd int32) (r int32)
TEXT ·Yposix_spawn_file_actions_addclose(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fa+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW fd+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xposix_spawn_file_actions_addclose(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawn_file_actions_adddup2(tls *TLS, fa uintptr, srcfd int32, fd int32) (r int32)
TEXT ·Yposix_spawn_file_actions_adddup2(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fa+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW srcfd+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW fd+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xposix_spawn_file_actions_adddup2(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawn_file_actions_addfchdir_np(tls *TLS, fa uintptr, fd int32) (r int32)
TEXT ·Yposix_spawn_file_actions_addfchdir_np(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fa+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW fd+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xposix_spawn_file_actions_addfchdir_np(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawn_file_actions_addopen(tls *TLS, fa uintptr, fd int32, path uintptr, flags int32, mode Tmode_t) (r int32)
TEXT ·Yposix_spawn_file_actions_addopen(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fa+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW fd+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD path+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	MOVW mode+36(FP), R0
	MOVW R0, 44(RSP)
	CALL ·Xposix_spawn_file_actions_addopen(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yposix_spawn_file_actions_destroy(tls *TLS, fa uintptr) (r int32)
TEXT ·Yposix_spawn_file_actions_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fa+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xposix_spawn_file_actions_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yposix_spawn_file_actions_init(tls *TLS, fa uintptr) (r int32)
TEXT ·Yposix_spawn_file_actions_init(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fa+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xposix_spawn_file_actions_init(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yposix_spawnattr_destroy(tls *TLS, attr uintptr) (r int32)
TEXT ·Yposix_spawnattr_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xposix_spawnattr_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yposix_spawnattr_getflags(tls *TLS, attr uintptr, flags uintptr) (r int32)
TEXT ·Yposix_spawnattr_getflags(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD flags+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_getflags(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_getpgroup(tls *TLS, attr uintptr, pgrp uintptr) (r int32)
TEXT ·Yposix_spawnattr_getpgroup(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pgrp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_getpgroup(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_getschedparam(tls *TLS, attr uintptr, schedparam uintptr) (r int32)
TEXT ·Yposix_spawnattr_getschedparam(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD schedparam+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_getschedparam(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_getschedpolicy(tls *TLS, attr uintptr, policy uintptr) (r int32)
TEXT ·Yposix_spawnattr_getschedpolicy(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD policy+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_getschedpolicy(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_getsigdefault(tls *TLS, attr uintptr, def uintptr) (r int32)
TEXT ·Yposix_spawnattr_getsigdefault(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD def+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_getsigdefault(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_getsigmask(tls *TLS, attr uintptr, mask uintptr) (r int32)
TEXT ·Yposix_spawnattr_getsigmask(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mask+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_getsigmask(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_init(tls *TLS, attr uintptr) (r int32)
TEXT ·Yposix_spawnattr_init(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xposix_spawnattr_init(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yposix_spawnattr_setflags(tls *TLS, attr uintptr, flags int16) (r int32)
TEXT ·Yposix_spawnattr_setflags(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH flags+16(FP), R0
	MOVH R0, 24(RSP)
	CALL ·Xposix_spawnattr_setflags(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_setpgroup(tls *TLS, attr uintptr, pgrp Tpid_t) (r int32)
TEXT ·Yposix_spawnattr_setpgroup(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW pgrp+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xposix_spawnattr_setpgroup(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_setschedparam(tls *TLS, attr uintptr, schedparam uintptr) (r int32)
TEXT ·Yposix_spawnattr_setschedparam(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD schedparam+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_setschedparam(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_setschedpolicy(tls *TLS, attr uintptr, policy int32) (r int32)
TEXT ·Yposix_spawnattr_setschedpolicy(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW policy+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xposix_spawnattr_setschedpolicy(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_setsigdefault(tls *TLS, attr uintptr, def uintptr) (r int32)
TEXT ·Yposix_spawnattr_setsigdefault(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD def+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_setsigdefault(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yposix_spawnattr_setsigmask(tls *TLS, attr uintptr, mask uintptr) (r int32)
TEXT ·Yposix_spawnattr_setsigmask(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mask+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xposix_spawnattr_setsigmask(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ypow(tls *TLS, x1 float64, y1 float64) (r float64)
TEXT ·Ypow(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpow(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ypow10(tls *TLS, x float64) (r float64)
TEXT ·Ypow10(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpow10(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ypow10f(tls *TLS, x float32) (r float32)
TEXT ·Ypow10f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xpow10f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ypow10l(tls *TLS, x float64) (r float64)
TEXT ·Ypow10l(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpow10l(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ypowf(tls *TLS, x1 float32, y1 float32) (r float32)
TEXT ·Ypowf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x1+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y1+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xpowf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ypowl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Ypowl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpowl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yppoll(tls *TLS, fds uintptr, n Tnfds_t, to uintptr, mask uintptr) (r int32)
TEXT ·Yppoll(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fds+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD to+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD mask+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xppoll(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yprctl(tls *TLS, op int32, va uintptr) (r int32)
TEXT ·Yprctl(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW op+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xprctl(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ypread(tls *TLS, fd int32, buf uintptr, size Tsize_t, ofs Toff_t) (r Tssize_t)
TEXT ·Ypread(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD ofs+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xpread(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ypreadv(tls *TLS, fd int32, iov uintptr, count int32, ofs Toff_t) (r Tssize_t)
TEXT ·Ypreadv(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW count+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD ofs+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xpreadv(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ypreadv2(tls *TLS, fd int32, iov uintptr, count int32, ofs Toff_t, flags int32) (r Tssize_t)
TEXT ·Ypreadv2(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW count+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD ofs+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xpreadv2(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Yprintf(tls *TLS, fmt uintptr, va uintptr) (r int32)
TEXT ·Yprintf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xprintf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yprlimit(tls *TLS, pid Tpid_t, resource int32, new_limit uintptr, old_limit uintptr) (r1 int32)
TEXT ·Yprlimit(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW resource+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD new_limit+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD old_limit+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xprlimit(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Yprocess_vm_readv(tls *TLS, pid Tpid_t, lvec uintptr, liovcnt uint64, rvec uintptr, riovcnt uint64, flags uint64) (r Tssize_t)
TEXT ·Yprocess_vm_readv(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD lvec+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD liovcnt+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD rvec+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD riovcnt+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD flags+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xprocess_vm_readv(SB)
	MOVD 64(RSP), R0
	MOVD R0, r+56(FP)
	RET

// func Yprocess_vm_writev(tls *TLS, pid Tpid_t, lvec uintptr, liovcnt uint64, rvec uintptr, riovcnt uint64, flags uint64) (r Tssize_t)
TEXT ·Yprocess_vm_writev(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD lvec+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD liovcnt+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD rvec+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD riovcnt+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD flags+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xprocess_vm_writev(SB)
	MOVD 64(RSP), R0
	MOVD R0, r+56(FP)
	RET

// func Ypselect(tls *TLS, n int32, rfds uintptr, wfds uintptr, efds uintptr, ts uintptr, mask uintptr) (r int32)
TEXT ·Ypselect(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD rfds+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD wfds+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD efds+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD ts+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD mask+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xpselect(SB)
	MOVW 64(RSP), R0
	MOVW R0, r+56(FP)
	RET

// func Ypsiginfo(tls *TLS, si uintptr, msg uintptr)
TEXT ·Ypsiginfo(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD si+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD msg+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpsiginfo(SB)
	RET

// func Ypsignal(tls *TLS, sig int32, msg uintptr)
TEXT ·Ypsignal(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sig+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD msg+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpsignal(SB)
	RET

// func Ypthread_attr_destroy(tls *TLS, a uintptr) (_2 int32)
TEXT ·Ypthread_attr_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_attr_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_attr_getdetachstate(tls *TLS, a uintptr, state uintptr) (_3 int32)
TEXT ·Ypthread_attr_getdetachstate(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD state+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_attr_getdetachstate(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_attr_init(tls *TLS, a uintptr) (_2 int32)
TEXT ·Ypthread_attr_init(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_attr_init(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_attr_setdetachstate(tls *TLS, a uintptr, state int32) (r int32)
TEXT ·Ypthread_attr_setdetachstate(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW state+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpthread_attr_setdetachstate(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ypthread_attr_setscope(tls *TLS, a uintptr, scope int32) (_3 int32)
TEXT ·Ypthread_attr_setscope(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW scope+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpthread_attr_setscope(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_attr_setstacksize(tls *TLS, a uintptr, stacksite Tsize_t) (_3 int32)
TEXT ·Ypthread_attr_setstacksize(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD stacksite+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_attr_setstacksize(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_barrier_destroy(tls *TLS, barrier uintptr) (_2 int32)
TEXT ·Ypthread_barrier_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD barrier+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_barrier_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_barrier_init(tls *TLS, barrier, attr uintptr, count uint32) (_3 int32)
TEXT ·Ypthread_barrier_init(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD barrier+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD attr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW count+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xpthread_barrier_init(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ypthread_barrier_wait(tls *TLS, barrier uintptr) (_2 int32)
TEXT ·Ypthread_barrier_wait(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD barrier+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_barrier_wait(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_cleanup_pop(tls *TLS, run int32)
TEXT ·Ypthread_cleanup_pop(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW run+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xpthread_cleanup_pop(SB)
	RET

// func Ypthread_cleanup_push(tls *TLS, f, x uintptr)
TEXT ·Ypthread_cleanup_push(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD x+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_cleanup_push(SB)
	RET

// func Ypthread_cond_broadcast(tls *TLS, c uintptr) (_2 int32)
TEXT ·Ypthread_cond_broadcast(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD c+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_cond_broadcast(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_cond_destroy(tls *TLS, c uintptr) (_2 int32)
TEXT ·Ypthread_cond_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD c+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_cond_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_cond_init(tls *TLS, c, a uintptr) (_2 int32)
TEXT ·Ypthread_cond_init(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD c+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_cond_init(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypthread_cond_signal(tls *TLS, c uintptr) (_2 int32)
TEXT ·Ypthread_cond_signal(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD c+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_cond_signal(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_cond_timedwait(tls *TLS, c, m, ts uintptr) (r int32)
TEXT ·Ypthread_cond_timedwait(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD c+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD m+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ts+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xpthread_cond_timedwait(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ypthread_cond_wait(tls *TLS, c, m uintptr) (_2 int32)
TEXT ·Ypthread_cond_wait(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD c+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD m+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_cond_wait(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypthread_condattr_destroy(tls *TLS, a uintptr) (_2 int32)
TEXT ·Ypthread_condattr_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_condattr_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_condattr_getclock(tls *TLS, a uintptr, clk uintptr) (_3 int32)
TEXT ·Ypthread_condattr_getclock(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD clk+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_condattr_getclock(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_condattr_getpshared(tls *TLS, a uintptr, pshared uintptr) (_3 int32)
TEXT ·Ypthread_condattr_getpshared(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pshared+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_condattr_getpshared(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_condattr_init(tls *TLS, a uintptr) (_2 int32)
TEXT ·Ypthread_condattr_init(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_condattr_init(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_condattr_setclock(tls *TLS, a uintptr, clk Tclockid_t) (_3 int32)
TEXT ·Ypthread_condattr_setclock(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW clk+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpthread_condattr_setclock(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_condattr_setpshared(tls *TLS, a uintptr, pshared int32) (_3 int32)
TEXT ·Ypthread_condattr_setpshared(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW pshared+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpthread_condattr_setpshared(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_detach(tls *TLS, t uintptr) (_2 int32)
TEXT ·Ypthread_detach(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_detach(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_equal(tls *TLS, t, u uintptr) (_2 int32)
TEXT ·Ypthread_equal(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD u+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_equal(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypthread_exit(tls *TLS, result uintptr)
TEXT ·Ypthread_exit(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD result+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_exit(SB)
	RET

// func Ypthread_getaffinity_np(tls *TLS, td Tpthread_t, size Tsize_t, set uintptr) (r int32)
TEXT ·Ypthread_getaffinity_np(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD td+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD set+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xpthread_getaffinity_np(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ypthread_getcpuclockid(tls *TLS, t uintptr, clockid uintptr) (_3 int32)
TEXT ·Ypthread_getcpuclockid(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD clockid+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_getcpuclockid(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_getname_np(tls *TLS, thread uintptr, name uintptr, len1 Tsize_t) (_4 int32)
TEXT ·Ypthread_getname_np(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD thread+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xpthread_getname_np(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Ypthread_getspecific(tls *TLS, k Tpthread_key_t) (_2 uintptr)
TEXT ·Ypthread_getspecific(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW k+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xpthread_getspecific(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ypthread_join(tls *TLS, t Tpthread_t, res uintptr) (r int32)
TEXT ·Ypthread_join(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD res+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_join(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ypthread_key_delete(tls *TLS, k Tpthread_key_t) (_2 int32)
TEXT ·Ypthread_key_delete(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW k+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xpthread_key_delete(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_kill(tls *TLS, t Tpthread_t, sig int32) (r1 int32)
TEXT ·Ypthread_kill(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW sig+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpthread_kill(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ypthread_mutex_destroy(tls *TLS, m uintptr) (_2 int32)
TEXT ·Ypthread_mutex_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutex_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutex_init(tls *TLS, m, a uintptr) (_2 int32)
TEXT ·Ypthread_mutex_init(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD a+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_mutex_init(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypthread_mutex_lock(tls *TLS, m uintptr) (_2 int32)
TEXT ·Ypthread_mutex_lock(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutex_lock(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutex_trylock(tls *TLS, m uintptr) (_2 int32)
TEXT ·Ypthread_mutex_trylock(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutex_trylock(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutex_unlock(tls *TLS, m uintptr) (_2 int32)
TEXT ·Ypthread_mutex_unlock(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD m+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutex_unlock(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutexattr_destroy(tls *TLS, a uintptr) (_2 int32)
TEXT ·Ypthread_mutexattr_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutexattr_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutexattr_init(tls *TLS, a uintptr) (_2 int32)
TEXT ·Ypthread_mutexattr_init(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutexattr_init(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutexattr_settype(tls *TLS, a uintptr, typ int32) (_3 int32)
TEXT ·Ypthread_mutexattr_settype(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW typ+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpthread_mutexattr_settype(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_self(tls *TLS) (_1 uintptr)
TEXT ·Ypthread_self(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xpthread_self(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Ypthread_setaffinity_np(tls *TLS, td Tpthread_t, size Tsize_t, set uintptr) (r int32)
TEXT ·Ypthread_setaffinity_np(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD td+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD set+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xpthread_setaffinity_np(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ypthread_setcancelstate(tls *TLS, new int32, old uintptr) (_3 int32)
TEXT ·Ypthread_setcancelstate(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW new+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD old+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_setcancelstate(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_setname_np(tls *TLS, thread uintptr, name uintptr) (_3 int32)
TEXT ·Ypthread_setname_np(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD thread+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_setname_np(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_setspecific(tls *TLS, k Tpthread_key_t, x uintptr) (_3 int32)
TEXT ·Ypthread_setspecific(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW k+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD x+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_setspecific(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_sigmask(tls *TLS, now int32, set, old uintptr) (_3 int32)
TEXT ·Ypthread_sigmask(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW now+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD set+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD old+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xpthread_sigmask(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ypthread_testcancel(tls *TLS)
TEXT ·Ypthread_testcancel(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xpthread_testcancel(SB)
	RET

// func Yptrace(tls *TLS, req int32, va uintptr) (r int64)
TEXT ·Yptrace(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW req+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xptrace(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yptsname(tls *TLS, fd int32) (r uintptr)
TEXT ·Yptsname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xptsname(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yptsname_r(tls *TLS, fd int32, buf uintptr, len1 Tsize_t) (r int32)
TEXT ·Yptsname_r(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xptsname_r(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yputc(tls *TLS, c1 int32, f1 uintptr) (r int32)
TEXT ·Yputc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c1+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputc(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yputc_unlocked(tls *TLS, c int32, f uintptr) (r int32)
TEXT ·Yputc_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputc_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yputchar(tls *TLS, c1 int32) (r int32)
TEXT ·Yputchar(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c1+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xputchar(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yputchar_unlocked(tls *TLS, c int32) (r int32)
TEXT ·Yputchar_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xputchar_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yputenv(tls *TLS, s uintptr) (r int32)
TEXT ·Yputenv(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xputenv(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yputgrent(tls *TLS, gr uintptr, f uintptr) (r1 int32)
TEXT ·Yputgrent(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD gr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputgrent(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Yputpwent(tls *TLS, pw uintptr, f uintptr) (r int32)
TEXT ·Yputpwent(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pw+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputpwent(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yputs(tls *TLS, s uintptr) (r1 int32)
TEXT ·Yputs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xputs(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yputspent(tls *TLS, sp uintptr, f uintptr) (r int32)
TEXT ·Yputspent(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sp+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputspent(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ypututline(tls *TLS, ut uintptr) (r uintptr)
TEXT ·Ypututline(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ut+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpututline(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ypututxline(tls *TLS, ut uintptr) (r uintptr)
TEXT ·Ypututxline(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ut+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpututxline(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yputw(tls *TLS, _x int32, f uintptr) (r int32)
TEXT ·Yputw(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputw(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yputwc(tls *TLS, c Twchar_t, f uintptr) (r Twint_t)
TEXT ·Yputwc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputwc(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yputwc_unlocked(tls *TLS, c Twchar_t, f uintptr) (r Twint_t)
TEXT ·Yputwc_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputwc_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yputwchar(tls *TLS, c Twchar_t) (r Twint_t)
TEXT ·Yputwchar(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xputwchar(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yputwchar_unlocked(tls *TLS, c Twchar_t) (r Twint_t)
TEXT ·Yputwchar_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xputwchar_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ypwrite(tls *TLS, fd int32, buf uintptr, size Tsize_t, ofs Toff_t) (r Tssize_t)
TEXT ·Ypwrite(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD ofs+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xpwrite(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ypwritev(tls *TLS, fd int32, iov uintptr, count int32, ofs Toff_t) (r Tssize_t)
TEXT ·Ypwritev(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW count+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD ofs+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xpwritev(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ypwritev2(tls *TLS, fd int32, iov uintptr, count int32, ofs Toff_t, flags int32) (r Tssize_t)
TEXT ·Ypwritev2(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW count+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD ofs+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xpwritev2(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Yqsort(tls *TLS, base uintptr, nel Tsize_t, width Tsize_t, __ccgo_fp_cmp Tcmpfun)
TEXT ·Yqsort(SB),$48-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 48(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 48(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cmp+32(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_qsort_3(SB)	// Create the closure for calling __ccgo_fp_cmp
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD base+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD nel+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD width+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD 48(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 40(RSP)
	CALL ·Xqsort(SB)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_qsort_3(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Yqsort_r(tls *TLS, base uintptr, nel Tsize_t, width Tsize_t, __ccgo_fp_cmp Tcmpfun, arg uintptr)
TEXT ·Yqsort_r(SB),$56-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 56(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 56(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cmp+32(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_qsort_r_3(SB)	// Create the closure for calling __ccgo_fp_cmp
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD base+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD nel+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD width+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD 56(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 40(RSP)
	MOVD arg+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xqsort_r(SB)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_qsort_r_3(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD _3+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD __ccgo_fp+32(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 40(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Yquick_exit(tls *TLS, code int32)
TEXT ·Yquick_exit(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW code+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xquick_exit(SB)
	RET

// func Yquotactl(tls *TLS, cmd int32, special uintptr, id int32, addr uintptr) (r int32)
TEXT ·Yquotactl(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW cmd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD special+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW id+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD addr+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xquotactl(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yraise(tls *TLS, sig int32) (r int32)
TEXT ·Yraise(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sig+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xraise(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yrand(tls *TLS) (r int32)
TEXT ·Yrand(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xrand(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yrand_r(tls *TLS, seed uintptr) (r int32)
TEXT ·Yrand_r(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD seed+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrand_r(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yrandom(tls *TLS) (r int64)
TEXT ·Yrandom(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xrandom(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Yrandom_r(t *TLS, buf, result uintptr) (_2 int32)
TEXT ·Yrandom_r(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD result+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xrandom_r(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Yread(tls *TLS, fd int32, buf uintptr, count Tsize_t) (r Tssize_t)
TEXT ·Yread(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD count+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xread(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yreadahead(tls *TLS, fd int32, pos Toff_t, len1 Tsize_t) (r Tssize_t)
TEXT ·Yreadahead(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pos+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xreadahead(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yreaddir(tls *TLS, dir uintptr) (r uintptr)
TEXT ·Yreaddir(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xreaddir(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yreaddir64(tls *TLS, dir uintptr) (r uintptr)
TEXT ·Yreaddir64(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xreaddir64(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yreaddir_r(tls *TLS, dir uintptr, buf uintptr, result uintptr) (r int32)
TEXT ·Yreaddir_r(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD result+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xreaddir_r(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yreadlink(tls *TLS, path uintptr, buf uintptr, bufsize Tsize_t) (r1 Tssize_t)
TEXT ·Yreadlink(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD bufsize+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xreadlink(SB)
	MOVD 40(RSP), R0
	MOVD R0, r1+32(FP)
	RET

// func Yreadlinkat(tls *TLS, fd int32, path uintptr, buf uintptr, bufsize Tsize_t) (r1 Tssize_t)
TEXT ·Yreadlinkat(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD bufsize+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xreadlinkat(SB)
	MOVD 48(RSP), R0
	MOVD R0, r1+40(FP)
	RET

// func Yreadv(tls *TLS, fd int32, iov uintptr, count int32) (r Tssize_t)
TEXT ·Yreadv(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW count+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xreadv(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yrealloc(tls *TLS, p uintptr, n Tsize_t) (r uintptr)
TEXT ·Yrealloc(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xrealloc(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yreallocarray(tls *TLS, ptr uintptr, m Tsize_t, n Tsize_t) (r uintptr)
TEXT ·Yreallocarray(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD m+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xreallocarray(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yrealpath(tls *TLS, filename uintptr, resolved uintptr) (r uintptr)
TEXT ·Yrealpath(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD resolved+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xrealpath(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yreboot(tls *TLS, type1 int32) (r int32)
TEXT ·Yreboot(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW type1+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xreboot(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yrecv(tls *TLS, fd int32, buf uintptr, len1 Tsize_t, flags int32) (r Tssize_t)
TEXT ·Yrecv(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xrecv(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yrecvfrom(tls *TLS, fd int32, buf uintptr, len1 Tsize_t, flags int32, addr uintptr, alen uintptr) (r1 Tssize_t)
TEXT ·Yrecvfrom(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD addr+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD alen+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xrecvfrom(SB)
	MOVD 64(RSP), R0
	MOVD R0, r1+56(FP)
	RET

// func Yrecvmmsg(tls *TLS, fd int32, msgvec uintptr, vlen uint32, flags uint32, timeout uintptr) (r int32)
TEXT ·Yrecvmmsg(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD msgvec+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW vlen+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	MOVD timeout+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xrecvmmsg(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yrecvmsg(tls *TLS, fd int32, msg uintptr, flags int32) (r2 Tssize_t)
TEXT ·Yrecvmsg(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD msg+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xrecvmsg(SB)
	MOVD 40(RSP), R0
	MOVD R0, r2+32(FP)
	RET

// func Yregcomp(tls *TLS, preg uintptr, regex uintptr, cflags int32) (r int32)
TEXT ·Yregcomp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD preg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD regex+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW cflags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xregcomp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yregerror(tls *TLS, e int32, preg uintptr, buf uintptr, size Tsize_t) (r Tsize_t)
TEXT ·Yregerror(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW e+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD preg+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xregerror(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yregexec(tls *TLS, preg uintptr, string1 uintptr, nmatch Tsize_t, pmatch uintptr, eflags int32) (r int32)
TEXT ·Yregexec(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD preg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD string1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nmatch+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD pmatch+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW eflags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xregexec(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Yregfree(tls *TLS, preg uintptr)
TEXT ·Yregfree(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD preg+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xregfree(SB)
	RET

// func Yremainder(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yremainder(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xremainder(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yremainderf(tls *TLS, x float32, y float32) (r float32)
TEXT ·Yremainderf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xremainderf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yremainderl(tls *TLS, x float64, y float64) (r float64)
TEXT ·Yremainderl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xremainderl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yremap_file_pages(tls *TLS, addr uintptr, size Tsize_t, prot int32, pgoff Tsize_t, flags int32) (r int32)
TEXT ·Yremap_file_pages(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW prot+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD pgoff+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xremap_file_pages(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Yremove(tls *TLS, path uintptr) (r1 int32)
TEXT ·Yremove(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xremove(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Yremovexattr(tls *TLS, path uintptr, name uintptr) (r int32)
TEXT ·Yremovexattr(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xremovexattr(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yremque(tls *TLS, element uintptr)
TEXT ·Yremque(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD element+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xremque(SB)
	RET

// func Yremquo(tls *TLS, x float64, y float64, quo uintptr) (r float64)
TEXT ·Yremquo(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD quo+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xremquo(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yremquof(tls *TLS, x float32, y float32, quo uintptr) (r float32)
TEXT ·Yremquof(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD quo+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xremquof(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yremquol(tls *TLS, x float64, y float64, quo uintptr) (r float64)
TEXT ·Yremquol(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD quo+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xremquol(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yrename(tls *TLS, old uintptr, new1 uintptr) (r int32)
TEXT ·Yrename(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD old+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD new1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xrename(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yrenameat(tls *TLS, oldfd int32, old uintptr, newfd int32, new1 uintptr) (r int32)
TEXT ·Yrenameat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW oldfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD old+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW newfd+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD new1+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xrenameat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yrenameat2(t *TLS, olddirfd int32, oldpath uintptr, newdirfd int32, newpath uintptr, flags int32) (_6 int32)
TEXT ·Yrenameat2(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW olddirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD oldpath+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW newdirfd+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD newpath+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xrenameat2(SB)
	MOVW 56(RSP), R0
	MOVW R0, _6+48(FP)
	RET

// func Yres_init(tls *TLS) (r int32)
TEXT ·Yres_init(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xres_init(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yres_mkquery(tls *TLS, op int32, dname uintptr, class int32, type1 int32, data uintptr, datalen int32, newrr uintptr, buf uintptr, buflen int32) (r int32)
TEXT ·Yres_mkquery(SB),$80-76
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW op+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD dname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW class+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW type1+28(FP), R0
	MOVW R0, 36(RSP)
	MOVD data+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW datalen+40(FP), R0
	MOVW R0, 48(RSP)
	MOVD newrr+48(FP), R0
	MOVD R0, 56(RSP)
	MOVD buf+56(FP), R0
	MOVD R0, 64(RSP)
	MOVW buflen+64(FP), R0
	MOVW R0, 72(RSP)
	CALL ·Xres_mkquery(SB)
	MOVW 80(RSP), R0
	MOVW R0, r+72(FP)
	RET

// func Yres_send(tls *TLS, _msg uintptr, _msglen int32, _answer uintptr, _anslen int32) (r int32)
TEXT ·Yres_send(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _msg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW _msglen+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD _answer+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW _anslen+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xres_send(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yrewind(tls *TLS, f uintptr)
TEXT ·Yrewind(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrewind(SB)
	RET

// func Yrewinddir(tls *TLS, dir uintptr)
TEXT ·Yrewinddir(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrewinddir(SB)
	RET

// func Yrindex(tls *TLS, s uintptr, c int32) (r uintptr)
TEXT ·Yrindex(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xrindex(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yrint(tls *TLS, x float64) (r float64)
TEXT ·Yrint(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrint(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yrintf(tls *TLS, x float32) (r float32)
TEXT ·Yrintf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xrintf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yrintl(tls *TLS, x float64) (r float64)
TEXT ·Yrintl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrintl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yrmdir(tls *TLS, path uintptr) (r int32)
TEXT ·Yrmdir(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrmdir(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yround(tls *TLS, x3 float64) (r float64)
TEXT ·Yround(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xround(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yroundf(tls *TLS, x3 float32) (r float32)
TEXT ·Yroundf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xroundf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yroundl(tls *TLS, x float64) (r float64)
TEXT ·Yroundl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xroundl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysbrk(tls *TLS, inc Tintptr_t) (r uintptr)
TEXT ·Ysbrk(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD inc+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsbrk(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yscalb(tls *TLS, x float64, fn float64) (r float64)
TEXT ·Yscalb(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fn+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xscalb(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yscalbf(tls *TLS, x float32, fn float32) (r float32)
TEXT ·Yscalbf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW fn+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xscalbf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yscalbln(tls *TLS, x float64, n int64) (r float64)
TEXT ·Yscalbln(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xscalbln(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yscalblnf(tls *TLS, x float32, n int64) (r float32)
TEXT ·Yscalblnf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xscalblnf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yscalblnl(tls *TLS, x float64, n int64) (r float64)
TEXT ·Yscalblnl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xscalblnl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yscalbn(tls *TLS, x float64, n int32) (r float64)
TEXT ·Yscalbn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xscalbn(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yscalbnf(tls *TLS, x float32, n int32) (r float32)
TEXT ·Yscalbnf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW n+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xscalbnf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yscalbnl(tls *TLS, x float64, n int32) (r float64)
TEXT ·Yscalbnl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW n+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xscalbnl(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yscandir(tls *TLS, path uintptr, res uintptr, __ccgo_fp_sel uintptr, __ccgo_fp_cmp uintptr) (r int32)
TEXT ·Yscandir(SB),$56-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $32, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 56(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 56(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_sel+24(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_scandir_2(SB)	// Create the closure for calling __ccgo_fp_sel
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 56(RSP), R0
	ADD $16, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cmp+32(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_scandir_3(SB)	// Create the closure for calling __ccgo_fp_cmp
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD res+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD 56(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 32(RSP)
	MOVD 56(RSP), R0 // Trampolines[0]
	ADD $16, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 40(RSP)
	CALL ·Xscandir(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $32, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_scandir_2(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp+16(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 24(RSP), R0
	MOVW R0, _2+24(FP)
	RET

TEXT ·__ccgo_abi0_scandir_3(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Yscanf(tls *TLS, fmt uintptr, va uintptr) (r int32)
TEXT ·Yscanf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xscanf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysched_get_priority_max(tls *TLS, policy int32) (r int32)
TEXT ·Ysched_get_priority_max(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW policy+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsched_get_priority_max(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysched_get_priority_min(tls *TLS, policy int32) (r int32)
TEXT ·Ysched_get_priority_min(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW policy+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsched_get_priority_min(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysched_getaffinity(tls *TLS, tid Tpid_t, size Tsize_t, set uintptr) (r int32)
TEXT ·Ysched_getaffinity(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW tid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD set+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsched_getaffinity(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysched_getparam(tls *TLS, pid Tpid_t, param uintptr) (r int32)
TEXT ·Ysched_getparam(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD param+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsched_getparam(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysched_getscheduler(tls *TLS, pid Tpid_t) (r int32)
TEXT ·Ysched_getscheduler(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsched_getscheduler(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysched_rr_get_interval(tls *TLS, pid Tpid_t, ts uintptr) (r int32)
TEXT ·Ysched_rr_get_interval(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD ts+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsched_rr_get_interval(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysched_setaffinity(tls *TLS, tid Tpid_t, size Tsize_t, set uintptr) (r int32)
TEXT ·Ysched_setaffinity(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW tid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD set+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsched_setaffinity(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysched_setparam(tls *TLS, pid Tpid_t, param uintptr) (r int32)
TEXT ·Ysched_setparam(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD param+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsched_setparam(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysched_setscheduler(tls *TLS, pid Tpid_t, sched int32, param uintptr) (r int32)
TEXT ·Ysched_setscheduler(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW sched+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD param+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsched_setscheduler(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysched_yield(tls *TLS) (_1 int32)
TEXT ·Ysched_yield(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsched_yield(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Ysecure_getenv(tls *TLS, name uintptr) (r uintptr)
TEXT ·Ysecure_getenv(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsecure_getenv(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yseed48(tls *TLS, s uintptr) (r uintptr)
TEXT ·Yseed48(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xseed48(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yseekdir(tls *TLS, dir uintptr, off int64)
TEXT ·Yseekdir(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD off+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xseekdir(SB)
	RET

// func Yselect(tls *TLS, n int32, rfds uintptr, wfds uintptr, efds uintptr, tv uintptr) (r int32)
TEXT ·Yselect(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD rfds+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD wfds+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD efds+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD tv+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xselect(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Ysem_close(tls *TLS, sem uintptr) (r int32)
TEXT ·Ysem_close(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sem+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsem_close(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysem_destroy(tls *TLS, sem uintptr) (r int32)
TEXT ·Ysem_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sem+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsem_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysem_getvalue(tls *TLS, sem uintptr, valp uintptr) (r int32)
TEXT ·Ysem_getvalue(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sem+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD valp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsem_getvalue(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysem_init(tls *TLS, sem uintptr, pshared int32, value uint32) (r int32)
TEXT ·Ysem_init(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sem+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW pshared+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW value+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xsem_init(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysem_open(tls *TLS, name uintptr, flags int32, va uintptr) (r uintptr)
TEXT ·Ysem_open(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsem_open(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ysem_post(tls *TLS, sem uintptr) (r int32)
TEXT ·Ysem_post(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sem+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsem_post(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysem_timedwait(tls *TLS, sem uintptr, at uintptr) (r1 int32)
TEXT ·Ysem_timedwait(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sem+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD at+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsem_timedwait(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ysem_trywait(tls *TLS, sem uintptr) (r int32)
TEXT ·Ysem_trywait(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sem+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsem_trywait(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysem_unlink(tls *TLS, name uintptr) (r int32)
TEXT ·Ysem_unlink(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsem_unlink(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysem_wait(tls *TLS, sem uintptr) (r int32)
TEXT ·Ysem_wait(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sem+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsem_wait(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysemctl(tls *TLS, id int32, num int32, cmd int32, va uintptr) (r1 int32)
TEXT ·Ysemctl(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW id+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW num+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW cmd+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsemctl(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ysemget(tls *TLS, key Tkey_t, n int32, fl int32) (r int32)
TEXT ·Ysemget(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW key+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW n+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW fl+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsemget(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysemop(tls *TLS, id int32, buf uintptr, n Tsize_t) (r int32)
TEXT ·Ysemop(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW id+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsemop(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysemtimedop(tls *TLS, id int32, buf uintptr, n Tsize_t, ts uintptr) (r int32)
TEXT ·Ysemtimedop(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW id+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD ts+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xsemtimedop(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ysend(tls *TLS, fd int32, buf uintptr, len1 Tsize_t, flags int32) (r Tssize_t)
TEXT ·Ysend(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xsend(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ysendfile(tls *TLS, out_fd int32, in_fd int32, ofs uintptr, count Tsize_t) (r Tssize_t)
TEXT ·Ysendfile(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW out_fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW in_fd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD ofs+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD count+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsendfile(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ysendmmsg(tls *TLS, fd int32, msgvec uintptr, vlen uint32, flags uint32) (r1 int32)
TEXT ·Ysendmmsg(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD msgvec+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW vlen+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	CALL ·Xsendmmsg(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ysendmsg(tls *TLS, fd int32, msg uintptr, flags int32) (r1 Tssize_t)
TEXT ·Ysendmsg(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD msg+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xsendmsg(SB)
	MOVD 40(RSP), R0
	MOVD R0, r1+32(FP)
	RET

// func Ysendto(tls *TLS, fd int32, buf uintptr, len1 Tsize_t, flags int32, addr uintptr, alen Tsocklen_t) (r1 Tssize_t)
TEXT ·Ysendto(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD addr+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW alen+48(FP), R0
	MOVW R0, 56(RSP)
	CALL ·Xsendto(SB)
	MOVD 64(RSP), R0
	MOVD R0, r1+56(FP)
	RET

// func Ysetbuf(tls *TLS, f uintptr, buf uintptr)
TEXT ·Ysetbuf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetbuf(SB)
	RET

// func Ysetbuffer(tls *TLS, f uintptr, buf uintptr, size Tsize_t)
TEXT ·Ysetbuffer(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsetbuffer(SB)
	RET

// func Ysetdomainname(tls *TLS, name uintptr, len1 Tsize_t) (r int32)
TEXT ·Ysetdomainname(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetdomainname(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysetenv(tls *TLS, var1 uintptr, value uintptr, overwrite int32) (r int32)
TEXT ·Ysetenv(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD var1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD value+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW overwrite+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xsetenv(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysetfsgid(tls *TLS, gid Tgid_t) (r int32)
TEXT ·Ysetfsgid(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW gid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsetfsgid(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysetfsuid(tls *TLS, uid Tuid_t) (r int32)
TEXT ·Ysetfsuid(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW uid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsetfsuid(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysetgid(tls *TLS, gid Tgid_t) (r int32)
TEXT ·Ysetgid(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW gid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsetgid(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysetgrent(tls *TLS)
TEXT ·Ysetgrent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetgrent(SB)
	RET

// func Ysethostent(tls *TLS, x int32)
TEXT ·Ysethostent(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsethostent(SB)
	RET

// func Ysethostname(tls *TLS, name uintptr, len1 Tsize_t) (r int32)
TEXT ·Ysethostname(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsethostname(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysetitimer(tls *TLS, which int32, new1 uintptr, old uintptr) (r1 int32)
TEXT ·Ysetitimer(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW which+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD new1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD old+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsetitimer(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ysetjmp(t *TLS, env uintptr) (_2 int32)
TEXT ·Ysetjmp(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD env+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsetjmp(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ysetkey(tls *TLS, key uintptr)
TEXT ·Ysetkey(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsetkey(SB)
	RET

// func Ysetlinebuf(tls *TLS, f uintptr)
TEXT ·Ysetlinebuf(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsetlinebuf(SB)
	RET

// func Ysetlocale(tls *TLS, cat int32, name uintptr) (r uintptr)
TEXT ·Ysetlocale(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW cat+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetlocale(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ysetlogmask(tls *TLS, maskpri int32) (r int32)
TEXT ·Ysetlogmask(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW maskpri+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsetlogmask(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysetmntent(tls *TLS, name uintptr, mode uintptr) (r uintptr)
TEXT ·Ysetmntent(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetmntent(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ysetnetent(tls *TLS, x int32)
TEXT ·Ysetnetent(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsetnetent(SB)
	RET

// func Ysetns(tls *TLS, fd int32, nstype int32) (r int32)
TEXT ·Ysetns(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW nstype+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xsetns(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysetpgid(tls *TLS, pid Tpid_t, pgid Tpid_t) (r int32)
TEXT ·Ysetpgid(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW pgid+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xsetpgid(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysetpgrp(tls *TLS) (r Tpid_t)
TEXT ·Ysetpgrp(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetpgrp(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ysetpriority(tls *TLS, which int32, who Tid_t, prio int32) (r int32)
TEXT ·Ysetpriority(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW which+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW who+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW prio+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsetpriority(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysetprotoent(tls *TLS, stayopen int32)
TEXT ·Ysetprotoent(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW stayopen+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsetprotoent(SB)
	RET

// func Ysetpwent(tls *TLS)
TEXT ·Ysetpwent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetpwent(SB)
	RET

// func Ysetrlimit(tls *TLS, resource int32, rlim uintptr) (r int32)
TEXT ·Ysetrlimit(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW resource+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD rlim+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetrlimit(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysetrlimit64(tls *TLS, resource int32, rlim uintptr) (r int32)
TEXT ·Ysetrlimit64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW resource+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD rlim+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetrlimit64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysetservent(tls *TLS, stayopen int32)
TEXT ·Ysetservent(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW stayopen+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsetservent(SB)
	RET

// func Ysetsid(tls *TLS) (r Tpid_t)
TEXT ·Ysetsid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetsid(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Ysetsockopt(tls *TLS, fd int32, level int32, optname int32, optval uintptr, optlen Tsocklen_t) (r2 int32)
TEXT ·Ysetsockopt(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW level+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW optname+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD optval+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW optlen+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xsetsockopt(SB)
	MOVW 48(RSP), R0
	MOVW R0, r2+40(FP)
	RET

// func Ysetspent(tls *TLS)
TEXT ·Ysetspent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetspent(SB)
	RET

// func Ysetstate(tls *TLS, state uintptr) (r uintptr)
TEXT ·Ysetstate(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD state+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsetstate(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysettimeofday(tls *TLS, tv uintptr, tz uintptr) (r int32)
TEXT ·Ysettimeofday(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tv+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tz+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsettimeofday(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysetuid(tls *TLS, uid Tuid_t) (r int32)
TEXT ·Ysetuid(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW uid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsetuid(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysetusershell(tls *TLS)
TEXT ·Ysetusershell(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetusershell(SB)
	RET

// func Ysetutent(tls *TLS)
TEXT ·Ysetutent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetutent(SB)
	RET

// func Ysetutxent(tls *TLS)
TEXT ·Ysetutxent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetutxent(SB)
	RET

// func Ysetvbuf(tls *TLS, f uintptr, buf uintptr, type1 int32, size Tsize_t) (r int32)
TEXT ·Ysetvbuf(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW type1+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xsetvbuf(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ysetxattr(tls *TLS, path uintptr, name uintptr, value uintptr, size Tsize_t, flags int32) (r int32)
TEXT ·Ysetxattr(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD value+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW flags+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xsetxattr(SB)
	MOVW 56(RSP), R0
	MOVW R0, r+48(FP)
	RET

// func Yshm_open(tls *TLS, name uintptr, flag int32, mode Tmode_t) (r int32)
TEXT ·Yshm_open(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flag+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW mode+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xshm_open(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yshm_unlink(tls *TLS, name uintptr) (r int32)
TEXT ·Yshm_unlink(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xshm_unlink(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yshmat(tls *TLS, id int32, addr uintptr, flag int32) (r uintptr)
TEXT ·Yshmat(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW id+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flag+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xshmat(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yshmctl(tls *TLS, id int32, cmd int32, buf uintptr) (r1 int32)
TEXT ·Yshmctl(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW id+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW cmd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xshmctl(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Yshmdt(tls *TLS, addr uintptr) (r int32)
TEXT ·Yshmdt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xshmdt(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yshmget(tls *TLS, key Tkey_t, size Tsize_t, flag int32) (r int32)
TEXT ·Yshmget(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW key+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flag+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xshmget(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yshutdown(tls *TLS, fd int32, how int32) (r1 int32)
TEXT ·Yshutdown(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW how+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xshutdown(SB)
	MOVW 24(RSP), R0
	MOVW R0, r1+16(FP)
	RET

// func Ysigaddset(tls *TLS, set uintptr, sig int32) (r int32)
TEXT ·Ysigaddset(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW sig+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsigaddset(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysigaltstack(tls *TLS, ss uintptr, old uintptr) (r int32)
TEXT ·Ysigaltstack(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ss+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD old+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsigaltstack(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysigandset(tls *TLS, dest uintptr, left uintptr, right uintptr) (r1 int32)
TEXT ·Ysigandset(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD left+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD right+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsigandset(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ysigdelset(tls *TLS, set uintptr, sig int32) (r int32)
TEXT ·Ysigdelset(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW sig+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsigdelset(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysigemptyset(tls *TLS, set uintptr) (r int32)
TEXT ·Ysigemptyset(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsigemptyset(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysigfillset(tls *TLS, set uintptr) (r int32)
TEXT ·Ysigfillset(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsigfillset(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysigisemptyset(tls *TLS, set uintptr) (r int32)
TEXT ·Ysigisemptyset(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsigisemptyset(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysigismember(tls *TLS, set uintptr, sig int32) (r int32)
TEXT ·Ysigismember(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW sig+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsigismember(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysignalfd(tls *TLS, fd int32, sigs uintptr, flags int32) (r int32)
TEXT ·Ysignalfd(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD sigs+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xsignalfd(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysignificand(tls *TLS, x float64) (r float64)
TEXT ·Ysignificand(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsignificand(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysignificandf(tls *TLS, x float32) (r float32)
TEXT ·Ysignificandf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsignificandf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysigorset(tls *TLS, dest uintptr, left uintptr, right uintptr) (r1 int32)
TEXT ·Ysigorset(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD left+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD right+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsigorset(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ysigpending(tls *TLS, set uintptr) (r int32)
TEXT ·Ysigpending(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsigpending(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysigprocmask(tls *TLS, how int32, set uintptr, old uintptr) (r1 int32)
TEXT ·Ysigprocmask(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW how+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD set+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD old+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsigprocmask(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ysigqueue(tls *TLS, pid Tpid_t, sig int32, value Tsigval) (r1 int32)
TEXT ·Ysigqueue(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW sig+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW value_Fsival_int+16(FP), R0
	MOVW R0, 24(RSP)
	MOVB value_F__ccgo_pad2_0+20(FP), R0
	MOVB R0, 28(RSP)
	MOVB value_F__ccgo_pad2_1+21(FP), R0
	MOVB R0, 29(RSP)
	MOVB value_F__ccgo_pad2_2+22(FP), R0
	MOVB R0, 30(RSP)
	MOVB value_F__ccgo_pad2_3+23(FP), R0
	MOVB R0, 31(RSP)
	CALL ·Xsigqueue(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ysigsuspend(tls *TLS, mask uintptr) (r int32)
TEXT ·Ysigsuspend(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD mask+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsigsuspend(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysigtimedwait(tls *TLS, mask uintptr, si uintptr, timeout uintptr) (r int32)
TEXT ·Ysigtimedwait(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD mask+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD si+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD timeout+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsigtimedwait(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysigwait(tls *TLS, mask uintptr, sig uintptr) (r int32)
TEXT ·Ysigwait(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD mask+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sig+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsigwait(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysigwaitinfo(tls *TLS, mask uintptr, si uintptr) (r int32)
TEXT ·Ysigwaitinfo(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD mask+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD si+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsigwaitinfo(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysin(tls *TLS, x3 float64) (r float64)
TEXT ·Ysin(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsin(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysincos(tls *TLS, x3 float64, sin uintptr, cos uintptr)
TEXT ·Ysincos(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sin+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD cos+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsincos(SB)
	RET

// func Ysincosf(tls *TLS, x3 float32, sin uintptr, cos uintptr)
TEXT ·Ysincosf(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD sin+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD cos+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsincosf(SB)
	RET

// func Ysincosl(tls *TLS, x float64, sin uintptr, cos uintptr)
TEXT ·Ysincosl(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sin+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD cos+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsincosl(SB)
	RET

// func Ysinf(tls *TLS, x3 float32) (r float32)
TEXT ·Ysinf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsinf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysinh(tls *TLS, x float64) (r float64)
TEXT ·Ysinh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsinh(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysinhf(tls *TLS, x float32) (r float32)
TEXT ·Ysinhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsinhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysinhl(tls *TLS, x float64) (r float64)
TEXT ·Ysinhl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsinhl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysinl(tls *TLS, x float64) (r float64)
TEXT ·Ysinl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsinl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysleep(tls *TLS, seconds uint32) (r uint32)
TEXT ·Ysleep(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW seconds+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsleep(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysnprintf(tls *TLS, s uintptr, n Tsize_t, fmt uintptr, va uintptr) (r int32)
TEXT ·Ysnprintf(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD fmt+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD va+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xsnprintf(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ysockatmark(tls *TLS, s int32) (r int32)
TEXT ·Ysockatmark(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW s+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsockatmark(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysocket(tls *TLS, domain int32, type1 int32, protocol int32) (r1 int32)
TEXT ·Ysocket(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW domain+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW type1+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW protocol+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsocket(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ysocketpair(tls *TLS, domain int32, type1 int32, protocol int32, fd uintptr) (r2 int32)
TEXT ·Ysocketpair(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW domain+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW type1+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW protocol+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD fd+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsocketpair(SB)
	MOVW 40(RSP), R0
	MOVW R0, r2+32(FP)
	RET

// func Ysplice(tls *TLS, fd_in int32, off_in uintptr, fd_out int32, off_out uintptr, len1 Tsize_t, flags uint32) (r Tssize_t)
TEXT ·Ysplice(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd_in+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD off_in+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW fd_out+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD off_out+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD len1+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW flags+48(FP), R0
	MOVW R0, 56(RSP)
	CALL ·Xsplice(SB)
	MOVD 64(RSP), R0
	MOVD R0, r+56(FP)
	RET

// func Ysprintf(tls *TLS, s uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Ysprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysqrt(tls *TLS, x float64) (r float64)
TEXT ·Ysqrt(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsqrt(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysqrtf(tls *TLS, x float32) (r float32)
TEXT ·Ysqrtf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsqrtf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysqrtl(tls *TLS, x float64) (r float64)
TEXT ·Ysqrtl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsqrtl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysrand(tls *TLS, s uint32)
TEXT ·Ysrand(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW s+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsrand(SB)
	RET

// func Ysrand48(tls *TLS, seed int64)
TEXT ·Ysrand48(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD seed+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsrand48(SB)
	RET

// func Ysrandom(tls *TLS, seed uint32)
TEXT ·Ysrandom(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW seed+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsrandom(SB)
	RET

// func Ysscanf(tls *TLS, s uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Ysscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ystat(tls *TLS, path uintptr, buf uintptr) (r int32)
TEXT ·Ystat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstat(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ystat64(tls *TLS, path uintptr, buf uintptr) (r int32)
TEXT ·Ystat64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstat64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ystatfs(tls *TLS, path uintptr, buf uintptr) (_3 int32)
TEXT ·Ystatfs(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstatfs(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ystatvfs(tls *TLS, path uintptr, buf uintptr) (r int32)
TEXT ·Ystatvfs(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstatvfs(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ystatx(tls *TLS, dirfd int32, path uintptr, flags int32, mask uint32, stx uintptr) (r int32)
TEXT ·Ystatx(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW mask+28(FP), R0
	MOVW R0, 36(RSP)
	MOVD stx+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xstatx(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ystime(tls *TLS, t uintptr) (r int32)
TEXT ·Ystime(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xstime(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ystpcpy(tls *TLS, d uintptr, s uintptr) (r uintptr)
TEXT ·Ystpcpy(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstpcpy(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystpncpy(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ystpncpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstpncpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrcasecmp(tls *TLS, _l uintptr, _r uintptr) (r1 int32)
TEXT ·Ystrcasecmp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _r+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcasecmp(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ystrcasecmp_l(tls *TLS, l uintptr, r uintptr, loc Tlocale_t) (r1 int32)
TEXT ·Ystrcasecmp_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD loc+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrcasecmp_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ystrcasestr(tls *TLS, h uintptr, n uintptr) (r uintptr)
TEXT ·Ystrcasestr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD h+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcasestr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrcat(tls *TLS, dest uintptr, src uintptr) (r uintptr)
TEXT ·Ystrcat(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcat(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrchr(tls *TLS, s uintptr, c int32) (r1 uintptr)
TEXT ·Ystrchr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xstrchr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r1+24(FP)
	RET

// func Ystrchrnul(tls *TLS, s uintptr, c int32) (r uintptr)
TEXT ·Ystrchrnul(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xstrchrnul(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrcmp(tls *TLS, l uintptr, r uintptr) (r1 int32)
TEXT ·Ystrcmp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcmp(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ystrcoll(tls *TLS, l uintptr, r uintptr) (r1 int32)
TEXT ·Ystrcoll(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcoll(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ystrcoll_l(tls *TLS, l uintptr, r uintptr, loc Tlocale_t) (r1 int32)
TEXT ·Ystrcoll_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD loc+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrcoll_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ystrcpy(tls *TLS, dest uintptr, src uintptr) (r uintptr)
TEXT ·Ystrcpy(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcpy(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrcspn(tls *TLS, s uintptr, c uintptr) (r Tsize_t)
TEXT ·Ystrcspn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD c+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcspn(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrdup(tls *TLS, s uintptr) (r uintptr)
TEXT ·Ystrdup(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xstrdup(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ystrerror(tls *TLS, e int32) (r uintptr)
TEXT ·Ystrerror(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW e+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xstrerror(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ystrerror_l(tls *TLS, e int32, loc Tlocale_t) (r uintptr)
TEXT ·Ystrerror_l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW e+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD loc+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrerror_l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrerror_r(tls *TLS, err int32, buf uintptr, buflen Tsize_t) (r int32)
TEXT ·Ystrerror_r(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW err+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buflen+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrerror_r(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ystrfmon(tls *TLS, s uintptr, n Tsize_t, fmt uintptr, va uintptr) (r Tssize_t)
TEXT ·Ystrfmon(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD fmt+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD va+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xstrfmon(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ystrfmon_l(tls *TLS, s uintptr, n Tsize_t, loc Tlocale_t, fmt uintptr, va uintptr) (r Tssize_t)
TEXT ·Ystrfmon_l(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD loc+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD fmt+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD va+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xstrfmon_l(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ystrftime(tls *TLS, s uintptr, n Tsize_t, f uintptr, tm uintptr) (r Tsize_t)
TEXT ·Ystrftime(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD tm+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xstrftime(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ystrftime_l(tls *TLS, s uintptr, n Tsize_t, f uintptr, tm uintptr, loc Tlocale_t) (r Tsize_t)
TEXT ·Ystrftime_l(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD tm+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD loc+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xstrftime_l(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ystrlcat(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r Tsize_t)
TEXT ·Ystrlcat(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrlcat(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrlcpy(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r Tsize_t)
TEXT ·Ystrlcpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrlcpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrlen(t *TLS, s uintptr) (r Tsize_t)
TEXT ·Ystrlen(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xstrlen(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ystrncasecmp(tls *TLS, _l uintptr, _r uintptr, n Tsize_t) (r1 int32)
TEXT ·Ystrncasecmp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrncasecmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ystrncasecmp_l(tls *TLS, l uintptr, r uintptr, n Tsize_t, loc Tlocale_t) (r1 int32)
TEXT ·Ystrncasecmp_l(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD loc+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xstrncasecmp_l(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Ystrncat(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ystrncat(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrncat(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrncmp(tls *TLS, _l uintptr, _r uintptr, n Tsize_t) (r1 int32)
TEXT ·Ystrncmp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrncmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ystrncpy(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ystrncpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrncpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrndup(tls *TLS, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ystrndup(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrndup(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrnlen(tls *TLS, s uintptr, n Tsize_t) (r Tsize_t)
TEXT ·Ystrnlen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrnlen(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrpbrk(tls *TLS, s uintptr, b uintptr) (r uintptr)
TEXT ·Ystrpbrk(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrpbrk(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrptime(tls *TLS, s uintptr, f uintptr, tm uintptr) (r uintptr)
TEXT ·Ystrptime(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD tm+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrptime(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrrchr(tls *TLS, s uintptr, c int32) (r uintptr)
TEXT ·Ystrrchr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xstrrchr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrsep(tls *TLS, str uintptr, sep uintptr) (r uintptr)
TEXT ·Ystrsep(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sep+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrsep(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrsignal(tls *TLS, signum int32) (r uintptr)
TEXT ·Ystrsignal(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW signum+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xstrsignal(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ystrspn(tls *TLS, s uintptr, c uintptr) (r Tsize_t)
TEXT ·Ystrspn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD c+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrspn(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrstr(tls *TLS, h uintptr, n uintptr) (r uintptr)
TEXT ·Ystrstr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD h+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrstr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrtod(tls *TLS, s uintptr, p uintptr) (r float64)
TEXT ·Ystrtod(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrtod(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrtod_l(tls *TLS, s uintptr, p uintptr, l Tlocale_t) (r float64)
TEXT ·Ystrtod_l(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrtod_l(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrtof(tls *TLS, s uintptr, p uintptr) (r float32)
TEXT ·Ystrtof(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrtof(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ystrtof_l(tls *TLS, s uintptr, p uintptr, l Tlocale_t) (r float32)
TEXT ·Ystrtof_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrtof_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ystrtoimax(tls *TLS, s uintptr, p uintptr, base int32) (r Tintmax_t)
TEXT ·Ystrtoimax(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xstrtoimax(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrtok(tls *TLS, s uintptr, sep uintptr) (r uintptr)
TEXT ·Ystrtok(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sep+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrtok(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrtok_r(tls *TLS, s uintptr, sep uintptr, p uintptr) (r uintptr)
TEXT ·Ystrtok_r(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sep+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD p+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrtok_r(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrtol(tls *TLS, s uintptr, p uintptr, base int32) (r int64)
TEXT ·Ystrtol(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xstrtol(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrtold(tls *TLS, s uintptr, p uintptr) (r float64)
TEXT ·Ystrtold(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrtold(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrtold_l(tls *TLS, s uintptr, p uintptr, l Tlocale_t) (r float64)
TEXT ·Ystrtold_l(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrtold_l(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrtoll(tls *TLS, s uintptr, p uintptr, base int32) (r int64)
TEXT ·Ystrtoll(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xstrtoll(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrtoul(tls *TLS, s uintptr, p uintptr, base int32) (r uint64)
TEXT ·Ystrtoul(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xstrtoul(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrtoull(tls *TLS, s uintptr, p uintptr, base int32) (r uint64)
TEXT ·Ystrtoull(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xstrtoull(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrtoumax(tls *TLS, s uintptr, p uintptr, base int32) (r Tuintmax_t)
TEXT ·Ystrtoumax(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xstrtoumax(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrverscmp(tls *TLS, l0 uintptr, r0 uintptr) (r1 int32)
TEXT ·Ystrverscmp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l0+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r0+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrverscmp(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ystrxfrm(tls *TLS, dest uintptr, src uintptr, n Tsize_t) (r Tsize_t)
TEXT ·Ystrxfrm(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrxfrm(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrxfrm_l(tls *TLS, dest uintptr, src uintptr, n Tsize_t, loc Tlocale_t) (r Tsize_t)
TEXT ·Ystrxfrm_l(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD loc+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xstrxfrm_l(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yswab(tls *TLS, _src uintptr, _dest uintptr, n Tssize_t)
TEXT ·Yswab(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _src+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _dest+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xswab(SB)
	RET

// func Yswapoff(tls *TLS, path uintptr) (r int32)
TEXT ·Yswapoff(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xswapoff(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yswapon(tls *TLS, path uintptr, flags int32) (r int32)
TEXT ·Yswapon(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xswapon(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yswprintf(tls *TLS, s uintptr, n Tsize_t, fmt uintptr, va uintptr) (r int32)
TEXT ·Yswprintf(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD fmt+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD va+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xswprintf(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yswscanf(tls *TLS, s uintptr, fmt uintptr, va uintptr) (r int32)
TEXT ·Yswscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xswscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysymlink(tls *TLS, existing uintptr, new1 uintptr) (r int32)
TEXT ·Ysymlink(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD existing+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD new1+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsymlink(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ysymlinkat(tls *TLS, existing uintptr, fd int32, new1 uintptr) (r int32)
TEXT ·Ysymlinkat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD existing+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW fd+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD new1+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsymlinkat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysync(tls *TLS)
TEXT ·Ysync(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsync(SB)
	RET

// func Ysync_file_range(tls *TLS, fd int32, pos Toff_t, len1 Toff_t, flags uint32) (r int32)
TEXT ·Ysync_file_range(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pos+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len1+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xsync_file_range(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ysyncfs(tls *TLS, fd int32) (r int32)
TEXT ·Ysyncfs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsyncfs(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysyscall(tls *TLS, n int64, va uintptr) (r int64)
TEXT ·Ysyscall(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsyscall(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ysysconf(tls *TLS, name int32) (r int64)
TEXT ·Ysysconf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW name+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsysconf(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ysysctlbyname(t *TLS, name, oldp, oldlenp, newp uintptr, newlen Tsize_t) (_3 int32)
TEXT ·Ysysctlbyname(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD oldp+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD oldlenp+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD newp+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD newlen+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xsysctlbyname(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Ysysinfo(tls *TLS, info uintptr) (r int32)
TEXT ·Ysysinfo(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD info+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsysinfo(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ysyslog(tls *TLS, priority int32, message uintptr, va uintptr)
TEXT ·Ysyslog(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW priority+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD message+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsyslog(SB)
	RET

// func Ysystem(t *TLS, command uintptr) (_2 int32)
TEXT ·Ysystem(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD command+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsystem(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ytan(tls *TLS, x3 float64) (r float64)
TEXT ·Ytan(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtan(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytanf(tls *TLS, x3 float32) (r float32)
TEXT ·Ytanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytanh(tls *TLS, x3 float64) (r float64)
TEXT ·Ytanh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtanh(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytanhf(tls *TLS, x3 float32) (r float32)
TEXT ·Ytanhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtanhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytanhl(tls *TLS, x float64) (r float64)
TEXT ·Ytanhl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtanhl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytanl(tls *TLS, x float64) (r float64)
TEXT ·Ytanl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtanl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytcdrain(tls *TLS, fd int32) (r int32)
TEXT ·Ytcdrain(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtcdrain(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytcflow(tls *TLS, fd int32, action int32) (r int32)
TEXT ·Ytcflow(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW action+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xtcflow(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytcflush(tls *TLS, fd int32, queue int32) (r int32)
TEXT ·Ytcflush(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW queue+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xtcflush(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytcgetattr(tls *TLS, fd int32, tio uintptr) (r int32)
TEXT ·Ytcgetattr(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD tio+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtcgetattr(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytcgetpgrp(tls *TLS, fd int32) (r Tpid_t)
TEXT ·Ytcgetpgrp(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtcgetpgrp(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytcgetsid(tls *TLS, fd int32) (r Tpid_t)
TEXT ·Ytcgetsid(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtcgetsid(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytcgetwinsize(tls *TLS, fd int32, wsz uintptr) (r int32)
TEXT ·Ytcgetwinsize(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD wsz+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtcgetwinsize(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytcsendbreak(tls *TLS, fd int32, dur int32) (r int32)
TEXT ·Ytcsendbreak(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW dur+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xtcsendbreak(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytcsetattr(tls *TLS, fd int32, act int32, tio uintptr) (r int32)
TEXT ·Ytcsetattr(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW act+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD tio+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtcsetattr(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytcsetpgrp(tls *TLS, fd int32, pgrp Tpid_t) (r int32)
TEXT ·Ytcsetpgrp(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW pgrp+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xtcsetpgrp(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytcsetwinsize(tls *TLS, fd int32, wsz uintptr) (r int32)
TEXT ·Ytcsetwinsize(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD wsz+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtcsetwinsize(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytdelete(tls *TLS, key uintptr, rootp uintptr, __ccgo_fp_cmp uintptr) (r uintptr)
TEXT ·Ytdelete(SB),$48-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 48(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 48(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cmp+24(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_tdelete_2(SB)	// Create the closure for calling __ccgo_fp_cmp
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD rootp+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD 48(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 32(RSP)
	CALL ·Xtdelete(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_tdelete_2(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ytdestroy(tls *TLS, root uintptr, __ccgo_fp_freekey uintptr)
TEXT ·Ytdestroy(SB),$32-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 32(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 32(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_freekey+16(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_tdestroy_1(SB)	// Create the closure for calling __ccgo_fp_freekey
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD root+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD 32(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 24(RSP)
	CALL ·Xtdestroy(SB)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_tdestroy_1(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp+16(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	RET

// func Ytee(tls *TLS, src int32, dest int32, len1 Tsize_t, flags uint32) (r Tssize_t)
TEXT ·Ytee(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW src+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW dest+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD len1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xtee(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ytelldir(tls *TLS, dir uintptr) (r int64)
TEXT ·Ytelldir(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtelldir(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytempnam(tls *TLS, dir uintptr, pfx uintptr) (r1 uintptr)
TEXT ·Ytempnam(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pfx+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtempnam(SB)
	MOVD 32(RSP), R0
	MOVD R0, r1+24(FP)
	RET

// func Ytextdomain(tls *TLS, domainname uintptr) (r uintptr)
TEXT ·Ytextdomain(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD domainname+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtextdomain(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytfind(tls *TLS, key uintptr, rootp uintptr, __ccgo_fp_cmp uintptr) (r uintptr)
TEXT ·Ytfind(SB),$48-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 48(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 48(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cmp+24(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_tfind_2(SB)	// Create the closure for calling __ccgo_fp_cmp
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD rootp+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD 48(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 32(RSP)
	CALL ·Xtfind(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_tfind_2(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ytgamma(tls *TLS, x3 float64) (r1 float64)
TEXT ·Ytgamma(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtgamma(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Ytgammaf(tls *TLS, x float32) (r float32)
TEXT ·Ytgammaf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtgammaf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytgammal(tls *TLS, x float64) (r float64)
TEXT ·Ytgammal(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtgammal(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytime(tls *TLS, t uintptr) (r Ttime_t)
TEXT ·Ytime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtime(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytimegm(tls *TLS, tm uintptr) (r Ttime_t)
TEXT ·Ytimegm(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tm+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtimegm(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytimer_delete(tls *TLS, t Ttimer_t) (r int32)
TEXT ·Ytimer_delete(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtimer_delete(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytimer_getoverrun(tls *TLS, t Ttimer_t) (r int32)
TEXT ·Ytimer_getoverrun(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtimer_getoverrun(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytimer_gettime(tls *TLS, t Ttimer_t, val uintptr) (r int32)
TEXT ·Ytimer_gettime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD val+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtimer_gettime(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytimer_settime(tls *TLS, t Ttimer_t, flags int32, val uintptr, old uintptr) (r int32)
TEXT ·Ytimer_settime(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD val+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD old+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xtimer_settime(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ytimerfd_create(tls *TLS, clockid int32, flags int32) (r int32)
TEXT ·Ytimerfd_create(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clockid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW flags+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xtimerfd_create(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytimerfd_gettime(tls *TLS, fd int32, cur uintptr) (r int32)
TEXT ·Ytimerfd_gettime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD cur+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtimerfd_gettime(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytimerfd_settime(tls *TLS, fd int32, flags int32, new1 uintptr, old uintptr) (r int32)
TEXT ·Ytimerfd_settime(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW flags+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD new1+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD old+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xtimerfd_settime(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ytimes(tls *TLS, tms uintptr) (r Tclock_t)
TEXT ·Ytimes(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tms+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtimes(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytimespec_get(tls *TLS, ts uintptr, base int32) (r int32)
TEXT ·Ytimespec_get(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ts+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW base+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xtimespec_get(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytmpfile(tls *TLS) (r uintptr)
TEXT ·Ytmpfile(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xtmpfile(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Ytmpnam(tls *TLS, buf uintptr) (r1 uintptr)
TEXT ·Ytmpnam(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtmpnam(SB)
	MOVD 24(RSP), R0
	MOVD R0, r1+16(FP)
	RET

// func Ytoascii(tls *TLS, c int32) (r int32)
TEXT ·Ytoascii(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtoascii(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytolower(tls *TLS, c int32) (r int32)
TEXT ·Ytolower(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtolower(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytolower_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Ytolower_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtolower_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytoupper(tls *TLS, c int32) (r int32)
TEXT ·Ytoupper(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtoupper(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytoupper_l(tls *TLS, c int32, l Tlocale_t) (r int32)
TEXT ·Ytoupper_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtoupper_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytowctrans(tls *TLS, wc Twint_t, trans Twctrans_t) (r Twint_t)
TEXT ·Ytowctrans(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD trans+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtowctrans(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytowctrans_l(tls *TLS, c Twint_t, t Twctrans_t, l Tlocale_t) (r Twint_t)
TEXT ·Ytowctrans_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD t+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD l+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xtowctrans_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ytowlower(tls *TLS, wc Twint_t) (r Twint_t)
TEXT ·Ytowlower(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtowlower(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytowlower_l(tls *TLS, c Twint_t, l Tlocale_t) (r Twint_t)
TEXT ·Ytowlower_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtowlower_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytowupper(tls *TLS, wc Twint_t) (r Twint_t)
TEXT ·Ytowupper(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtowupper(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytowupper_l(tls *TLS, c Twint_t, l Tlocale_t) (r Twint_t)
TEXT ·Ytowupper_l(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtowupper_l(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytrunc(tls *TLS, x3 float64) (r float64)
TEXT ·Ytrunc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x3+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtrunc(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytruncate(tls *TLS, path uintptr, length Toff_t) (r int32)
TEXT ·Ytruncate(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD length+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtruncate(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ytruncf(tls *TLS, x3 float32) (r float32)
TEXT ·Ytruncf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x3+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtruncf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytruncl(tls *TLS, x float64) (r float64)
TEXT ·Ytruncl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtruncl(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytsearch(tls *TLS, key uintptr, rootp uintptr, __ccgo_fp_cmp uintptr) (r1 uintptr)
TEXT ·Ytsearch(SB),$48-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 48(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 48(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_cmp+24(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_tsearch_2(SB)	// Create the closure for calling __ccgo_fp_cmp
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD rootp+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD 48(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 32(RSP)
	CALL ·Xtsearch(SB)
	MOVD 40(RSP), R0
	MOVD R0, r1+32(FP)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_tsearch_2(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD _2+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	MOVW 32(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Yttyname(tls *TLS, fd int32) (r uintptr)
TEXT ·Yttyname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xttyname(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yttyname_r(tls *TLS, fd int32, name uintptr, size Tsize_t) (r int32)
TEXT ·Yttyname_r(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xttyname_r(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ytwalk(tls *TLS, root uintptr, __ccgo_fp_action uintptr)
TEXT ·Ytwalk(SB),$32-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 32(RSP) // Trampolines[0]
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 32(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD __ccgo_fp_action+16(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_twalk_1(SB)	// Create the closure for calling __ccgo_fp_action
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD root+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD 32(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 24(RSP)
	CALL ·Xtwalk(SB)
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSFree(SB)
	RET

TEXT ·__ccgo_abi0_twalk_1(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _0+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD _1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW _2+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW _3+20(FP), R0
	MOVW R0, 28(RSP)
	MOVD __ccgo_fp+24(FP), R0
	CALL (R0)	// Call the ABI0 code ptr
	RET

// func Ytzset(tls *TLS)
TEXT ·Ytzset(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xtzset(SB)
	RET

// func Yualarm(tls *TLS, value uint32, interval uint32) (r uint32)
TEXT ·Yualarm(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW value+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW interval+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xualarm(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yulckpwdf(tls *TLS) (r int32)
TEXT ·Yulckpwdf(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xulckpwdf(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yulimit(tls *TLS, cmd int32, va uintptr) (r int64)
TEXT ·Yulimit(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW cmd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xulimit(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yumask(tls *TLS, mode Tmode_t) (r Tmode_t)
TEXT ·Yumask(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW mode+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xumask(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yumount(tls *TLS, special uintptr) (r int32)
TEXT ·Yumount(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD special+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xumount(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yumount2(tls *TLS, special uintptr, flags int32) (r int32)
TEXT ·Yumount2(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD special+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xumount2(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yuname(tls *TLS, uts uintptr) (r int32)
TEXT ·Yuname(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD uts+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xuname(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yungetc(tls *TLS, c int32, f uintptr) (r int32)
TEXT ·Yungetc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xungetc(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yungetwc(tls *TLS, c Twint_t, f uintptr) (r Twint_t)
TEXT ·Yungetwc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xungetwc(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yunlink(tls *TLS, path uintptr) (r int32)
TEXT ·Yunlink(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xunlink(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yunlinkat(tls *TLS, fd int32, path uintptr, flag int32) (r int32)
TEXT ·Yunlinkat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flag+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xunlinkat(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yunlockpt(tls *TLS, fd int32) (r int32)
TEXT ·Yunlockpt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xunlockpt(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yunsetenv(tls *TLS, name uintptr) (r int32)
TEXT ·Yunsetenv(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xunsetenv(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yunshare(tls *TLS, flags int32) (r int32)
TEXT ·Yunshare(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW flags+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xunshare(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yupdwtmp(tls *TLS, f uintptr, u uintptr)
TEXT ·Yupdwtmp(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD u+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xupdwtmp(SB)
	RET

// func Yupdwtmpx(tls *TLS, f uintptr, u uintptr)
TEXT ·Yupdwtmpx(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD u+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xupdwtmpx(SB)
	RET

// func Yuselocale(tls *TLS, new1 Tlocale_t) (r Tlocale_t)
TEXT ·Yuselocale(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD new1+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xuselocale(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yusleep(tls *TLS, useconds uint32) (r int32)
TEXT ·Yusleep(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW useconds+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xusleep(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yutime(tls *TLS, path uintptr, times uintptr) (r int32)
TEXT ·Yutime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD times+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xutime(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yutimensat(tls *TLS, fd int32, path uintptr, times uintptr, flags int32) (r1 int32)
TEXT ·Yutimensat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD path+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD times+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xutimensat(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Yutimes(tls *TLS, path uintptr, times uintptr) (r int32)
TEXT ·Yutimes(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD times+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xutimes(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yuuid_copy(t *TLS, dst, src uintptr)
TEXT ·Yuuid_copy(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dst+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xuuid_copy(SB)
	RET

// func Yuuid_generate_random(t *TLS, out uintptr)
TEXT ·Yuuid_generate_random(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD out+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xuuid_generate_random(SB)
	RET

// func Yuuid_parse(t *TLS, in uintptr, uu uintptr) (_3 int32)
TEXT ·Yuuid_parse(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD in+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD uu+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xuuid_parse(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yuuid_unparse(t *TLS, uu, out uintptr)
TEXT ·Yuuid_unparse(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD uu+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD out+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xuuid_unparse(SB)
	RET

// func Yvasprintf(tls *TLS, s uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvasprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvasprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yvdprintf(tls *TLS, fd int32, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvdprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvdprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yverr(tls *TLS, status int32, fmt uintptr, ap Tva_list)
TEXT ·Yverr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW status+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xverr(SB)
	RET

// func Yverrx(tls *TLS, status int32, fmt uintptr, ap Tva_list)
TEXT ·Yverrx(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW status+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xverrx(SB)
	RET

// func Yversionsort(tls *TLS, a uintptr, b uintptr) (r int32)
TEXT ·Yversionsort(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xversionsort(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yvfork(tls *TLS) (r Tpid_t)
TEXT ·Yvfork(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xvfork(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yvfprintf(tls *TLS, f uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvfprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvfprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yvfscanf(tls *TLS, f uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvfscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvfscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yvfwprintf(tls *TLS, f uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvfwprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvfwprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yvfwscanf(tls *TLS, f uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvfwscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvfwscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yvhangup(tls *TLS) (r int32)
TEXT ·Yvhangup(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xvhangup(SB)
	MOVW 16(RSP), R0
	MOVW R0, r+8(FP)
	RET

// func Yvmsplice(tls *TLS, fd int32, iov uintptr, cnt Tsize_t, flags uint32) (r Tssize_t)
TEXT ·Yvmsplice(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD cnt+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xvmsplice(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Yvprintf(tls *TLS, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvprintf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xvprintf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yvscanf(tls *TLS, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvscanf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xvscanf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yvsnprintf(tls *TLS, s uintptr, n Tsize_t, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvsnprintf(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD fmt+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD ap+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xvsnprintf(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Yvsprintf(tls *TLS, s uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvsprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvsprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yvsscanf(tls *TLS, s uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvsscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvsscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yvswprintf(tls *TLS, s uintptr, n Tsize_t, fmt uintptr, ap Tva_list) (r1 int32)
TEXT ·Yvswprintf(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD fmt+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD ap+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xvswprintf(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Yvswscanf(tls *TLS, s uintptr, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvswscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvswscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Yvwarn(tls *TLS, fmt uintptr, ap Tva_list)
TEXT ·Yvwarn(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xvwarn(SB)
	RET

// func Yvwarnx(tls *TLS, fmt uintptr, ap Tva_list)
TEXT ·Yvwarnx(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xvwarnx(SB)
	RET

// func Yvwprintf(tls *TLS, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvwprintf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xvwprintf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yvwscanf(tls *TLS, fmt uintptr, ap Tva_list) (r int32)
TEXT ·Yvwscanf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xvwscanf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ywait(tls *TLS, status uintptr) (r Tpid_t)
TEXT ·Ywait(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD status+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xwait(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ywait3(tls *TLS, status uintptr, options int32, usage uintptr) (r Tpid_t)
TEXT ·Ywait3(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD status+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW options+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD usage+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwait3(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ywait4(tls *TLS, pid Tpid_t, status uintptr, options int32, ru uintptr) (r1 Tpid_t)
TEXT ·Ywait4(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD status+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW options+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD ru+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xwait4(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Ywaitid(tls *TLS, type1 Tidtype_t, id Tid_t, info uintptr, options int32) (r int32)
TEXT ·Ywaitid(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW type1+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW id+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD info+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW options+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwaitid(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ywaitpid(tls *TLS, pid Tpid_t, status uintptr, options int32) (r Tpid_t)
TEXT ·Ywaitpid(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD status+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW options+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwaitpid(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ywarn(tls *TLS, fmt uintptr, va uintptr)
TEXT ·Ywarn(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwarn(SB)
	RET

// func Ywarnx(tls *TLS, fmt uintptr, va uintptr)
TEXT ·Ywarnx(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwarnx(SB)
	RET

// func Ywcpcpy(tls *TLS, d uintptr, s uintptr) (r uintptr)
TEXT ·Ywcpcpy(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcpcpy(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcpncpy(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ywcpncpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcpncpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcrtomb(tls *TLS, s uintptr, wc Twchar_t, st uintptr) (r Tsize_t)
TEXT ·Ywcrtomb(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW wc+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD st+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcrtomb(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcscasecmp(tls *TLS, l uintptr, r uintptr) (r1 int32)
TEXT ·Ywcscasecmp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcscasecmp(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ywcscasecmp_l(tls *TLS, l uintptr, r uintptr, locale Tlocale_t) (r1 int32)
TEXT ·Ywcscasecmp_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD locale+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcscasecmp_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ywcscat(tls *TLS, dest uintptr, src uintptr) (r uintptr)
TEXT ·Ywcscat(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcscat(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcschr(tls *TLS, s uintptr, c Twchar_t) (r uintptr)
TEXT ·Ywcschr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xwcschr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcscmp(tls *TLS, l uintptr, r uintptr) (r1 int32)
TEXT ·Ywcscmp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcscmp(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ywcscoll(tls *TLS, l uintptr, r uintptr) (r1 int32)
TEXT ·Ywcscoll(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcscoll(SB)
	MOVW 32(RSP), R0
	MOVW R0, r1+24(FP)
	RET

// func Ywcscoll_l(tls *TLS, l uintptr, r uintptr, locale Tlocale_t) (r1 int32)
TEXT ·Ywcscoll_l(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD locale+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcscoll_l(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ywcscpy(tls *TLS, d uintptr, s uintptr) (r uintptr)
TEXT ·Ywcscpy(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcscpy(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcscspn(tls *TLS, s uintptr, c uintptr) (r Tsize_t)
TEXT ·Ywcscspn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD c+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcscspn(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcsdup(tls *TLS, s uintptr) (r uintptr)
TEXT ·Ywcsdup(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xwcsdup(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ywcsftime(tls *TLS, wcs uintptr, n Tsize_t, f uintptr, tm uintptr) (r Tsize_t)
TEXT ·Ywcsftime(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD wcs+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD tm+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xwcsftime(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ywcsftime_l(tls *TLS, s uintptr, n Tsize_t, f uintptr, tm uintptr, loc Tlocale_t) (r Tsize_t)
TEXT ·Ywcsftime_l(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD f+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD tm+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD loc+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xwcsftime_l(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ywcslen(tls *TLS, s uintptr) (r Tsize_t)
TEXT ·Ywcslen(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xwcslen(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ywcsncasecmp(tls *TLS, l uintptr, r uintptr, n Tsize_t) (r1 int32)
TEXT ·Ywcsncasecmp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcsncasecmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ywcsncasecmp_l(tls *TLS, l uintptr, r uintptr, n Tsize_t, locale Tlocale_t) (r1 int32)
TEXT ·Ywcsncasecmp_l(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD locale+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xwcsncasecmp_l(SB)
	MOVW 48(RSP), R0
	MOVW R0, r1+40(FP)
	RET

// func Ywcsncat(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ywcsncat(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcsncat(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcsncmp(tls *TLS, l uintptr, r uintptr, n Tsize_t) (r1 int32)
TEXT ·Ywcsncmp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcsncmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ywcsncpy(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ywcsncpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcsncpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcsnlen(tls *TLS, s uintptr, n Tsize_t) (r Tsize_t)
TEXT ·Ywcsnlen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcsnlen(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcsnrtombs(tls *TLS, dst uintptr, wcs uintptr, wn Tsize_t, n Tsize_t, st uintptr) (r Tsize_t)
TEXT ·Ywcsnrtombs(SB),$56-56
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dst+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD wcs+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD wn+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD n+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD st+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xwcsnrtombs(SB)
	MOVD 56(RSP), R0
	MOVD R0, r+48(FP)
	RET

// func Ywcspbrk(tls *TLS, s uintptr, b uintptr) (r uintptr)
TEXT ·Ywcspbrk(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD b+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcspbrk(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcsrchr(tls *TLS, s uintptr, c Twchar_t) (r uintptr)
TEXT ·Ywcsrchr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xwcsrchr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcsrtombs(tls *TLS, s uintptr, ws uintptr, n Tsize_t, st uintptr) (r Tsize_t)
TEXT ·Ywcsrtombs(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ws+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD st+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xwcsrtombs(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ywcsspn(tls *TLS, s uintptr, c uintptr) (r Tsize_t)
TEXT ·Ywcsspn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD c+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcsspn(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcsstr(tls *TLS, h uintptr, n uintptr) (r uintptr)
TEXT ·Ywcsstr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD h+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcsstr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcstod(tls *TLS, s uintptr, p uintptr) (r float64)
TEXT ·Ywcstod(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcstod(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcstof(tls *TLS, s uintptr, p uintptr) (r float32)
TEXT ·Ywcstof(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcstof(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ywcstoimax(tls *TLS, s uintptr, p uintptr, base int32) (r Tintmax_t)
TEXT ·Ywcstoimax(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwcstoimax(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcstok(tls *TLS, s uintptr, sep uintptr, p uintptr) (r uintptr)
TEXT ·Ywcstok(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD sep+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD p+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcstok(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcstol(tls *TLS, s uintptr, p uintptr, base int32) (r int64)
TEXT ·Ywcstol(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwcstol(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcstold(tls *TLS, s uintptr, p uintptr) (r float64)
TEXT ·Ywcstold(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcstold(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcstoll(tls *TLS, s uintptr, p uintptr, base int32) (r int64)
TEXT ·Ywcstoll(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwcstoll(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcstombs(tls *TLS, s uintptr, ws uintptr, n Tsize_t) (r Tsize_t)
TEXT ·Ywcstombs(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ws+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcstombs(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcstoul(tls *TLS, s uintptr, p uintptr, base int32) (r uint64)
TEXT ·Ywcstoul(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwcstoul(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcstoull(tls *TLS, s uintptr, p uintptr, base int32) (r uint64)
TEXT ·Ywcstoull(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwcstoull(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcstoumax(tls *TLS, s uintptr, p uintptr, base int32) (r Tuintmax_t)
TEXT ·Ywcstoumax(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW base+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwcstoumax(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcswcs(tls *TLS, haystack uintptr, needle uintptr) (r uintptr)
TEXT ·Ywcswcs(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD haystack+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD needle+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcswcs(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcswidth(tls *TLS, wcs uintptr, n Tsize_t) (r int32)
TEXT ·Ywcswidth(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD wcs+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwcswidth(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ywcsxfrm(tls *TLS, dest uintptr, src uintptr, n Tsize_t) (r Tsize_t)
TEXT ·Ywcsxfrm(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwcsxfrm(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywcsxfrm_l(tls *TLS, dest uintptr, src uintptr, n Tsize_t, loc Tlocale_t) (r Tsize_t)
TEXT ·Ywcsxfrm_l(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD loc+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xwcsxfrm_l(SB)
	MOVD 48(RSP), R0
	MOVD R0, r+40(FP)
	RET

// func Ywctob(tls *TLS, c Twint_t) (r int32)
TEXT ·Ywctob(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xwctob(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ywctomb(tls *TLS, s uintptr, wc Twchar_t) (r int32)
TEXT ·Ywctomb(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW wc+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xwctomb(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ywctrans(tls *TLS, class uintptr) (r Twctrans_t)
TEXT ·Ywctrans(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD class+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xwctrans(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ywctrans_l(tls *TLS, s uintptr, l Tlocale_t) (r Twctrans_t)
TEXT ·Ywctrans_l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwctrans_l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywctype(tls *TLS, s uintptr) (r Twctype_t)
TEXT ·Ywctype(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xwctype(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ywctype_l(tls *TLS, s uintptr, l Tlocale_t) (r Twctype_t)
TEXT ·Ywctype_l(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD l+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwctype_l(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ywcwidth(tls *TLS, wc Twchar_t) (r int32)
TEXT ·Ywcwidth(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xwcwidth(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ywmemchr(tls *TLS, s uintptr, c Twchar_t, n Tsize_t) (r uintptr)
TEXT ·Ywmemchr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwmemchr(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywmemcmp(tls *TLS, l uintptr, r uintptr, n Tsize_t) (r1 int32)
TEXT ·Ywmemcmp(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD l+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD r+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwmemcmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, r1+32(FP)
	RET

// func Ywmemcpy(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ywmemcpy(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwmemcpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywmemmove(tls *TLS, d uintptr, s uintptr, n Tsize_t) (r uintptr)
TEXT ·Ywmemmove(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwmemmove(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywmemset(tls *TLS, d uintptr, c Twchar_t, n Tsize_t) (r uintptr)
TEXT ·Ywmemset(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD d+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwmemset(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywprintf(tls *TLS, fmt uintptr, va uintptr) (r int32)
TEXT ·Ywprintf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwprintf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Ywrite(tls *TLS, fd int32, buf uintptr, count Tsize_t) (r Tssize_t)
TEXT ·Ywrite(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD count+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwrite(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywritev(tls *TLS, fd int32, iov uintptr, count int32) (r Tssize_t)
TEXT ·Ywritev(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW count+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwritev(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ywscanf(tls *TLS, fmt uintptr, va uintptr) (r int32)
TEXT ·Ywscanf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fmt+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD va+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xwscanf(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yy0(tls *TLS, x float64) (r float64)
TEXT ·Yy0(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xy0(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yy0f(tls *TLS, x float32) (r float32)
TEXT ·Yy0f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xy0f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yy1(tls *TLS, x float64) (r float64)
TEXT ·Yy1(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xy1(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Yy1f(tls *TLS, x float32) (r float32)
TEXT ·Yy1f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xy1f(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yyn(tls *TLS, n int32, x float64) (r float64)
TEXT ·Yyn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD x+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xyn(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yynf(tls *TLS, n int32, x float32) (r float32)
TEXT ·Yynf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW n+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW x+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xynf(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET
