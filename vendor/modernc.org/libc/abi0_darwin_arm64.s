// Code generated for darwin/arm64 by 'qbecc --abi0wrap .', DO NOT EDIT.

#include "funcdata.h"
#include "textflag.h"

// func Y_IO_putc(t *TLS, c int32, fp uintptr) (_3 int32)
TEXT ·Y_IO_putc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD fp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X_IO_putc(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y_NSGetEnviron(t *TLS) (_1 uintptr)
TEXT ·Y_NSGetEnviron(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X_NSGetEnviron(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y___errno_location(t *TLS) (_1 uintptr)
TEXT ·Y___errno_location(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X___errno_location(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__assert_fail(t *TLS, assertion, file uintptr, line uint32, function uintptr)
TEXT ·Y__assert_fail(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD assertion+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD file+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW line+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD function+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__assert_fail(SB)
	RET

// func Y__assert_rtn(t *TLS, function, file uintptr, line int32, assertion uintptr)
TEXT ·Y__assert_rtn(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD function+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD file+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW line+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD assertion+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·X__assert_rtn(SB)
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

// func Y__builtin___memcpy_chk(t *TLS, dest, src uintptr, n, os types.Size_t) (r uintptr)
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

// func Y__builtin___memmove_chk(t *TLS, dest, src uintptr, n, os types.Size_t) (_3 uintptr)
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

// func Y__builtin___memset_chk(t *TLS, s uintptr, c int32, n, os types.Size_t) (_4 uintptr)
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

// func Y__builtin___snprintf_chk(t *TLS, str uintptr, maxlen types.Size_t, flag int32, os types.Size_t, format, args uintptr) (r int32)
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

// func Y__builtin___sprintf_chk(t *TLS, s uintptr, flag int32, os types.Size_t, format, args uintptr) (r int32)
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

// func Y__builtin___strcat_chk(t *TLS, dest, src uintptr, os types.Size_t) (r uintptr)
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

// func Y__builtin___strcpy_chk(t *TLS, dest, src uintptr, os types.Size_t) (_3 uintptr)
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

// func Y__builtin___strncpy_chk(t *TLS, dest, src uintptr, n, os types.Size_t) (r uintptr)
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

// func Y__builtin___vsnprintf_chk(t *TLS, str uintptr, maxlen types.Size_t, flag int32, os types.Size_t, format, args uintptr) (r int32)
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

// func Y__builtin___vsprintf_chk(t *TLS, s uintptr, flag int32, os types.Size_t, format, args uintptr) (r int32)
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

// func Y__builtin_alloca(tls *TLS, size size_t) (_2 uintptr)
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

// func Y__builtin_bzero(t *TLS, s uintptr, n types.Size_t)
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

// func Y__builtin_fma(t *TLS, x, y, z float64) (_2 float64)
TEXT ·Y__builtin_fma(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD z+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__builtin_fma(SB)
	MOVD 40(RSP), R0
	MOVD R0, _2+32(FP)
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

// func Y__builtin_getentropy(t *TLS, buf uintptr, n types.Size_t) (_3 int32)
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

// func Y__builtin_hypot(t *TLS, x float64, y float64) (r float64)
TEXT ·Y__builtin_hypot(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
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

// func Y__builtin_inff(t *TLS) (_1 float32)
TEXT ·Y__builtin_inff(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
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

// func Y__builtin_malloc(t *TLS, size types.Size_t) (_2 uintptr)
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

// func Y__builtin_memcmp(t *TLS, s1, s2 uintptr, n types.Size_t) (_3 int32)
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

// func Y__builtin_memcpy(t *TLS, dest, src uintptr, n types.Size_t) (r uintptr)
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

// func Y__builtin_memset(t *TLS, s uintptr, c int32, n types.Size_t) (_4 uintptr)
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

// func Y__builtin_mmap(t *TLS, addr uintptr, length types.Size_t, prot, flags, fd int32, offset types.Off_t) (_5 uintptr)
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

// func Y__builtin_nanf(t *TLS, s uintptr) (_2 float32)
TEXT ·Y__builtin_nanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
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

// func Y__builtin_object_size(t *TLS, p uintptr, typ int32) (_3 types.Size_t)
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

// func Y__builtin_printf(t *TLS, s, args uintptr) (_2 int32)
TEXT ·Y__builtin_printf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD args+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__builtin_printf(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
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

// func Y__builtin_round(t *TLS, x float64) (_2 float64)
TEXT ·Y__builtin_round(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__builtin_round(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Y__builtin_roundf(t *TLS, x float32) (_2 float32)
TEXT ·Y__builtin_roundf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__builtin_roundf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
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

// func Y__builtin_snprintf(t *TLS, str uintptr, size types.Size_t, format, args uintptr) (_4 int32)
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

// func Y__builtin_strlen(t *TLS, s uintptr) (_2 types.Size_t)
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

// func Y__builtin_vsnprintf(t *TLS, str uintptr, size types.Size_t, format, va uintptr) (_4 int32)
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

// func Y__ctype_get_mb_cur_max(t *TLS) (_1 types.Size_t)
TEXT ·Y__ctype_get_mb_cur_max(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__ctype_get_mb_cur_max(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__darwin_check_fd_set_overflow(tls *TLS, _ int32, _ uintptr, _ int32) (_4 int32)
TEXT ·Y__darwin_check_fd_set_overflow(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD _+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW _+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·X__darwin_check_fd_set_overflow(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Y__darwin_fd_clr(tls *TLS, _fd int32, _p uintptr) (_3 int32)
TEXT ·Y__darwin_fd_clr(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD _p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__darwin_fd_clr(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__darwin_fd_isset(tls *TLS, _fd int32, _p uintptr) (_3 int32)
TEXT ·Y__darwin_fd_isset(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD _p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__darwin_fd_isset(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__darwin_fd_set(tls *TLS, _fd int32, _p uintptr) (_3 int32)
TEXT ·Y__darwin_fd_set(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD _p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__darwin_fd_set(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__env_rm_add(tls *TLS, old uintptr, new uintptr)
TEXT ·Y__env_rm_add(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD old+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD new+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__env_rm_add(SB)
	RET

// func Y__errno_location(t *TLS) (_1 uintptr)
TEXT ·Y__errno_location(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__errno_location(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__error(t *TLS) (_1 uintptr)
TEXT ·Y__error(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__error(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__floatscan(tls *TLS, f uintptr, prec int32, pok int32) (_4 float64)
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
	MOVD R0, _4+24(FP)
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

// func Y__fpclassify(tls *TLS, x float64) (_2 int32)
TEXT ·Y__fpclassify(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fpclassify(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__fpclassifyd(tls *TLS, x float64) (r int32)
TEXT ·Y__fpclassifyd(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fpclassifyd(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Y__fpclassifyf(tls *TLS, x float32) (_2 int32)
TEXT ·Y__fpclassifyf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__fpclassifyf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__fpclassifyl(tls *TLS, x float64) (_2 int32)
TEXT ·Y__fpclassifyl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__fpclassifyl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__h_errno_location(tls *TLS) (_1 uintptr)
TEXT ·Y__h_errno_location(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__h_errno_location(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Y__inet_aton(tls *TLS, s0 uintptr, dest uintptr) (_3 int32)
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
	MOVW R0, _3+24(FP)
	RET

// func Y__inline_isnand(t *TLS, x float64) (_2 int32)
TEXT ·Y__inline_isnand(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__inline_isnand(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__inline_isnanf(t *TLS, x float32) (_2 int32)
TEXT ·Y__inline_isnanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__inline_isnanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__inline_isnanl(t *TLS, x float64) (_2 int32)
TEXT ·Y__inline_isnanl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__inline_isnanl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__intscan(tls *TLS, f uintptr, base uint32, pok int32, lim uint64) (_5 uint64)
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
	MOVD R0, _5+32(FP)
	RET

// func Y__isctype(tls *TLS, _c int32, _f uint64) (_3 int32)
TEXT ·Y__isctype(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD _f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__isctype(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
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

// func Y__isnan(t *TLS, arg float64) (_2 int32)
TEXT ·Y__isnan(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD arg+8(FP), R0
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

// func Y__isoc99_sscanf(t *TLS, str, format, va uintptr) (_2 int32)
TEXT ·Y__isoc99_sscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD format+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__isoc99_sscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Y__istype(tls *TLS, _c int32, _f uint64) (_3 int32)
TEXT ·Y__istype(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD _f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__istype(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__lookup_ipliteral(tls *TLS, buf uintptr, name uintptr, family int32) (_4 int32)
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
	MOVW R0, _4+32(FP)
	RET

// func Y__lookup_name(tls *TLS, buf uintptr, canon uintptr, name uintptr, family int32, flags int32) (_6 int32)
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
	MOVW R0, _6+40(FP)
	RET

// func Y__lookup_serv(tls *TLS, buf uintptr, name uintptr, proto int32, socktype int32, flags int32) (_6 int32)
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
	MOVW R0, _6+40(FP)
	RET

// func Y__maskrune(tls *TLS, _c int32, _f uint64) (_3 int32)
TEXT ·Y__maskrune(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD _f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__maskrune(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Y__putenv(tls *TLS, s uintptr, l size_t, r uintptr) (_4 int32)
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
	MOVW R0, _4+32(FP)
	RET

// func Y__shgetc(tls *TLS, f uintptr) (_2 int32)
TEXT ·Y__shgetc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__shgetc(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__shlim(tls *TLS, f uintptr, lim off_t)
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

// func Y__sincos_stret(t *TLS, f float64) (_2 struct{F__sinval float64; F__cosval float64})
TEXT ·Y__sincos_stret(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__sincos_stret(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2_F__sinval+16(FP)
	MOVD 32(RSP), R0
	MOVD R0, _2_F__cosval+24(FP)
	RET

// func Y__sincosf_stret(t *TLS, f float32) (_2 struct{F__sinval float32; F__cosval float32})
TEXT ·Y__sincosf_stret(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW f+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__sincosf_stret(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2_F__sinval+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, _2_F__cosval+20(FP)
	RET

// func Y__sincospi_stret(t *TLS, f float64) (_2 struct{F__sinval float64; F__cosval float64})
TEXT ·Y__sincospi_stret(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__sincospi_stret(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2_F__sinval+16(FP)
	MOVD 32(RSP), R0
	MOVD R0, _2_F__cosval+24(FP)
	RET

// func Y__sincospif_stret(t *TLS, f float32) (_2 struct{F__sinval float32; F__cosval float32})
TEXT ·Y__sincospif_stret(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW f+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__sincospif_stret(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2_F__sinval+16(FP)
	MOVW 28(RSP), R0
	MOVW R0, _2_F__cosval+20(FP)
	RET

// func Y__srget(t *TLS, f uintptr) (_2 int32)
TEXT ·Y__srget(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__srget(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__strchrnul(tls *TLS, s uintptr, c int32) (_3 uintptr)
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
	MOVD R0, _3+24(FP)
	RET

// func Y__strncasecmp_l(tls *TLS, l uintptr, r uintptr, n size_t, loc locale_t) (_5 int32)
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
	MOVW R0, _5+40(FP)
	RET

// func Y__svfscanf(t *TLS, f uintptr, p, q uintptr) (_3 int32)
TEXT ·Y__svfscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD p+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD q+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·X__svfscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Y__swbuf(t *TLS, i int32, f uintptr) (_3 int32)
TEXT ·Y__swbuf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW i+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD f+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·X__swbuf(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
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

// func Y__tolower(tls *TLS, c int32) (_2 int32)
TEXT ·Y__tolower(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__tolower(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__toread(tls *TLS, f uintptr) (_2 int32)
TEXT ·Y__toread(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__toread(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__toread_needs_stdio_exit(tls *TLS)
TEXT ·Y__toread_needs_stdio_exit(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·X__toread_needs_stdio_exit(SB)
	RET

// func Y__toupper(tls *TLS, c int32) (_2 int32)
TEXT ·Y__toupper(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__toupper(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__uflow(tls *TLS, f uintptr) (_2 int32)
TEXT ·Y__uflow(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·X__uflow(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y__wcwidth(tls *TLS, _c int32) (_2 int32)
TEXT ·Y__wcwidth(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X__wcwidth(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Y_exit(_ *TLS, status int32)
TEXT ·Y_exit(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW status+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·X_exit(SB)
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

// func Yabort(t *TLS)
TEXT ·Yabort(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xabort(SB)
	RET

// func Yabs(t *TLS, j int32) (_2 int32)
TEXT ·Yabs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW j+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xabs(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yaccept(t *TLS, sockfd int32, addr uintptr, addrlen uintptr) (_4 int32)
TEXT ·Yaccept(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD addrlen+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xaccept(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Yaccess(t *TLS, pathname uintptr, mode int32) (_3 int32)
TEXT ·Yaccess(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW mode+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xaccess(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yacos(t *TLS, x float64) (_2 float64)
TEXT ·Yacos(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xacos(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yacosf(t *TLS, x float32) (_2 float32)
TEXT ·Yacosf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xacosf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yacosh(t *TLS, x float64) (_2 float64)
TEXT ·Yacosh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xacosh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yacoshf(t *TLS, x float32) (_2 float32)
TEXT ·Yacoshf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xacoshf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yalarm(t *TLS, seconds uint32) (_2 uint32)
TEXT ·Yalarm(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW seconds+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xalarm(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yalloca(tls *TLS, size size_t) (_2 uintptr)
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

// func Yarc4random_buf(t *TLS, buf uintptr, buflen size_t)
TEXT ·Yarc4random_buf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buflen+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xarc4random_buf(SB)
	RET

// func Yasin(t *TLS, x float64) (_2 float64)
TEXT ·Yasin(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xasin(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yasinf(t *TLS, x float32) (_2 float32)
TEXT ·Yasinf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xasinf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yasinh(t *TLS, x float64) (_2 float64)
TEXT ·Yasinh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xasinh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yasinhf(t *TLS, x float32) (_2 float32)
TEXT ·Yasinhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xasinhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yatan(t *TLS, x float64) (_2 float64)
TEXT ·Yatan(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatan(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yatan2(t *TLS, x, y float64) (_2 float64)
TEXT ·Yatan2(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xatan2(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Yatan2f(t *TLS, x, y float32) (_2 float32)
TEXT ·Yatan2f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xatan2f(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yatanf(t *TLS, x float32) (_2 float32)
TEXT ·Yatanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xatanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yatanh(t *TLS, x float64) (_2 float64)
TEXT ·Yatanh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatanh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yatanhf(t *TLS, x float32) (_2 float32)
TEXT ·Yatanhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xatanhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yatof(t *TLS, nptr uintptr) (_2 float64)
TEXT ·Yatof(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD nptr+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatof(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yatoi(t *TLS, nptr uintptr) (_2 int32)
TEXT ·Yatoi(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD nptr+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatoi(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yatol(t *TLS, nptr uintptr) (_2 long)
TEXT ·Yatol(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD nptr+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xatol(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ybind(t *TLS, sockfd int32, addr uintptr, addrlen uint32) (_4 int32)
TEXT ·Ybind(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW addrlen+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xbind(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Ybsearch(tls *TLS, key uintptr, base uintptr, nel size_t, width size_t, cmp uintptr) (_6 uintptr)
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
	MOVD cmp+40(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_bsearch_4(SB)	// Create the closure for calling cmp
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
	MOVD R0, _6+48(FP)
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

// func Ybzero(t *TLS, s uintptr, n types.Size_t)
TEXT ·Ybzero(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xbzero(SB)
	RET

// func Ycalloc(t *TLS, n, size types.Size_t) (_2 uintptr)
TEXT ·Ycalloc(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcalloc(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
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

// func Yceil(t *TLS, x float64) (_2 float64)
TEXT ·Yceil(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xceil(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yceilf(t *TLS, x float32) (_2 float32)
TEXT ·Yceilf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xceilf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ycfgetospeed(t *TLS, termios_p uintptr) (_2 termios.Speed_t)
TEXT ·Ycfgetospeed(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD termios_p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcfgetospeed(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ychdir(t *TLS, path uintptr) (_2 int32)
TEXT ·Ychdir(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xchdir(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ychflags(t *TLS, path uintptr, flags uint32) (_3 int32)
TEXT ·Ychflags(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xchflags(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ychmod(t *TLS, pathname uintptr, mode types.Mode_t) (_3 int32)
TEXT ·Ychmod(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH mode+16(FP), R0
	MOVH R0, 24(RSP)
	CALL ·Xchmod(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ychown(t *TLS, pathname uintptr, owner types.Uid_t, group types.Gid_t) (_4 int32)
TEXT ·Ychown(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW owner+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW group+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xchown(SB)
	MOVW 32(RSP), R0
	MOVW R0, _4+24(FP)
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

// func Yclock(t *TLS) (r time.Clock_t)
TEXT ·Yclock(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xclock(SB)
	MOVD 16(RSP), R0
	MOVD R0, r+8(FP)
	RET

// func Yclock_gettime(t *TLS, clk_id int32, tp uintptr) (_3 int32)
TEXT ·Yclock_gettime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW clk_id+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD tp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xclock_gettime(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yclose(t *TLS, fd int32) (_2 int32)
TEXT ·Yclose(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xclose(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yclosedir(t *TLS, dir uintptr) (_2 int32)
TEXT ·Yclosedir(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xclosedir(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yconfstr(t *TLS, name int32, buf uintptr, size types.Size_t) (_4 types.Size_t)
TEXT ·Yconfstr(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW name+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xconfstr(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yconnect(t *TLS, sockfd int32, addr uintptr, addrlen uint32) (_4 int32)
TEXT ·Yconnect(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW addrlen+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xconnect(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Ycopysign(t *TLS, x, y float64) (_2 float64)
TEXT ·Ycopysign(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xcopysign(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ycopysignf(t *TLS, x, y float32) (_2 float32)
TEXT ·Ycopysignf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xcopysignf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ycopysignl(tls *TLS, x float64, y float64) (_3 float64)
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
	MOVD R0, _3+24(FP)
	RET

// func Ycos(t *TLS, x float64) (_2 float64)
TEXT ·Ycos(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcos(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ycosf(t *TLS, x float32) (_2 float32)
TEXT ·Ycosf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xcosf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ycosh(t *TLS, x float64) (_2 float64)
TEXT ·Ycosh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xcosh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ycoshf(t *TLS, x float32) (_2 float32)
TEXT ·Ycoshf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xcoshf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yctime(t *TLS, timep uintptr) (_2 uintptr)
TEXT ·Yctime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD timep+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xctime(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yctime_r(t *TLS, timep, buf uintptr) (_2 uintptr)
TEXT ·Yctime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD timep+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xctime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ydifftime(t *TLS, time1, time0 types.Time_t) (_2 float64)
TEXT ·Ydifftime(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD time1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD time0+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xdifftime(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ydigittoint(tls *TLS, _c int32) (_2 int32)
TEXT ·Ydigittoint(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xdigittoint(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ydiv(tls *TLS, num, den int32) (r Tdiv_t)
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

// func Ydup2(t *TLS, oldfd, newfd int32) (_2 int32)
TEXT ·Ydup2(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW oldfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW newfd+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xdup2(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yendpwent(t *TLS)
TEXT ·Yendpwent(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xendpwent(SB)
	RET

// func Yerf(t *TLS, x float64) (_2 float64)
TEXT ·Yerf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xerf(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yerfc(t *TLS, x float64) (_2 float64)
TEXT ·Yerfc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xerfc(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
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

// func Yexit(t *TLS, status int32)
TEXT ·Yexit(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW status+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xexit(SB)
	RET

// func Yexp(t *TLS, x float64) (_2 float64)
TEXT ·Yexp(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexp(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yexp2(t *TLS, x float64) (_2 float64)
TEXT ·Yexp2(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xexp2(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yexpf(t *TLS, x float32) (_2 float32)
TEXT ·Yexpf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xexpf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
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

// func Yfabs(t *TLS, x float64) (_2 float64)
TEXT ·Yfabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yfabsf(t *TLS, x float32) (_2 float32)
TEXT ·Yfabsf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfabsf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yfabsl(tls *TLS, x float64) (_2 float64)
TEXT ·Yfabsl(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfabsl(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yfaccessat(t *TLS, dirfd int32, pathname uintptr, mode int32, flags int32) (_5 int32)
TEXT ·Yfaccessat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mode+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	CALL ·Xfaccessat(SB)
	MOVW 40(RSP), R0
	MOVW R0, _5+32(FP)
	RET

// func Yfchmod(t *TLS, fd int32, mode types.Mode_t) (_3 int32)
TEXT ·Yfchmod(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVH mode+12(FP), R0
	MOVH R0, 20(RSP)
	CALL ·Xfchmod(SB)
	MOVW 24(RSP), R0
	MOVW R0, _3+16(FP)
	RET

// func Yfchmodat(t *TLS, dirfd int32, pathname uintptr, mode types.Mode_t, flags int32) (_5 int32)
TEXT ·Yfchmodat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVH mode+24(FP), R0
	MOVH R0, 32(RSP)
	MOVW flags+28(FP), R0
	MOVW R0, 36(RSP)
	CALL ·Xfchmodat(SB)
	MOVW 40(RSP), R0
	MOVW R0, _5+32(FP)
	RET

// func Yfchown(t *TLS, fd int32, owner types.Uid_t, group types.Gid_t) (_4 int32)
TEXT ·Yfchown(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW owner+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW group+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xfchown(SB)
	MOVW 32(RSP), R0
	MOVW R0, _4+24(FP)
	RET

// func Yfchownat(t *TLS, dirfd int32, pathname uintptr, owner types.Uid_t, group types.Gid_t, flags int32) (_6 int32)
TEXT ·Yfchownat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW owner+24(FP), R0
	MOVW R0, 32(RSP)
	MOVW group+28(FP), R0
	MOVW R0, 36(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xfchownat(SB)
	MOVW 48(RSP), R0
	MOVW R0, _6+40(FP)
	RET

// func Yfclose(t *TLS, stream uintptr) (_2 int32)
TEXT ·Yfclose(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfclose(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yfcntl(t *TLS, fd, cmd int32, args uintptr) (_3 int32)
TEXT ·Yfcntl(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW cmd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD args+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfcntl(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yfcntl64(t *TLS, fd, cmd int32, args uintptr) (r int32)
TEXT ·Yfcntl64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW cmd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD args+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfcntl64(SB)
	MOVW 32(RSP), R0
	MOVW R0, r+24(FP)
	RET

// func Yfdim(t *TLS, x, y float64) (_2 float64)
TEXT ·Yfdim(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfdim(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Yfdopen(t *TLS, fd int32, mode uintptr) (_3 uintptr)
TEXT ·Yfdopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfdopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Yfdopendir(t *TLS, fd int32) (_2 uintptr)
TEXT ·Yfdopendir(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfdopendir(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yfeof(t *TLS, f uintptr) (r int32)
TEXT ·Yfeof(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfeof(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Yferror(t *TLS, stream uintptr) (_2 int32)
TEXT ·Yferror(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xferror(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yfflush(t *TLS, stream uintptr) (_2 int32)
TEXT ·Yfflush(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfflush(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
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

// func Yfgetc(t *TLS, stream uintptr) (_2 int32)
TEXT ·Yfgetc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfgetc(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yfgets(t *TLS, s uintptr, size int32, stream uintptr) (_4 uintptr)
TEXT ·Yfgets(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW size+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD stream+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfgets(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yfileno(t *TLS, stream uintptr) (_2 int32)
TEXT ·Yfileno(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfileno(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yflock(t *TLS, fd, operation int32) (_2 int32)
TEXT ·Yflock(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW operation+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xflock(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yflockfile(t *TLS, filehandle uintptr)
TEXT ·Yflockfile(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filehandle+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xflockfile(SB)
	RET

// func Yfloor(t *TLS, x float64) (_2 float64)
TEXT ·Yfloor(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfloor(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yfloorf(t *TLS, x float32) (_2 float32)
TEXT ·Yfloorf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfloorf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yfma(t *TLS, x, y, z float64) (_2 float64)
TEXT ·Yfma(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD z+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfma(SB)
	MOVD 40(RSP), R0
	MOVD R0, _2+32(FP)
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

// func Yfmod(t *TLS, x, y float64) (_2 float64)
TEXT ·Yfmod(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfmod(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Yfmodf(t *TLS, x, y float32) (_2 float32)
TEXT ·Yfmodf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfmodf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yfmodl(tls *TLS, x float64, y float64) (_3 float64)
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
	MOVD R0, _3+24(FP)
	RET

// func Yfopen(t *TLS, pathname, mode uintptr) (_2 uintptr)
TEXT ·Yfopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Yfopen64(t *TLS, pathname, mode uintptr) (_2 uintptr)
TEXT ·Yfopen64(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfopen64(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
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

// func Yfpathconf(t *TLS, fd int32, name int32) (_3 long)
TEXT ·Yfpathconf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW name+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xfpathconf(SB)
	MOVD 24(RSP), R0
	MOVD R0, _3+16(FP)
	RET

// func Yfprintf(t *TLS, stream, format, args uintptr) (_2 int32)
TEXT ·Yfprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD format+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD args+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Yfputc(t *TLS, c int32, stream uintptr) (_3 int32)
TEXT ·Yfputc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD stream+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputc(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yfputs(t *TLS, s, stream uintptr) (_2 int32)
TEXT ·Yfputs(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD stream+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfputs(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Yfread(t *TLS, ptr uintptr, size, nmemb types.Size_t, stream uintptr) (_4 types.Size_t)
TEXT ·Yfread(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nmemb+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD stream+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfread(SB)
	MOVD 48(RSP), R0
	MOVD R0, _4+40(FP)
	RET

// func Yfree(t *TLS, p uintptr)
TEXT ·Yfree(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
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

// func Yfreopen(t *TLS, pathname, mode, stream uintptr) (r uintptr)
TEXT ·Yfreopen(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD stream+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfreopen(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Yfrexp(t *TLS, x float64, exp uintptr) (_3 float64)
TEXT ·Yfrexp(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD exp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfrexp(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Yfrexpf(t *TLS, x float32, exp uintptr) (_3 float32)
TEXT ·Yfrexpf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD exp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfrexpf(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yfscanf(t *TLS, stream, format, va uintptr) (_2 int32)
TEXT ·Yfscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD format+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xfscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Yfsctl(t *TLS, path uintptr, request ulong, data uintptr, options uint32) (_5 int32)
TEXT ·Yfsctl(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD request+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD data+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW options+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xfsctl(SB)
	MOVW 48(RSP), R0
	MOVW R0, _5+40(FP)
	RET

// func Yfseek(t *TLS, stream uintptr, offset long, whence int32) (_4 int32)
TEXT ·Yfseek(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD offset+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xfseek(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Yfseeko(t *TLS, stream uintptr, offset types.Off_t, whence int32) (_4 int32)
TEXT ·Yfseeko(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD offset+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xfseeko(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Yfstat(t *TLS, fd int32, statbuf uintptr) (_3 int32)
TEXT ·Yfstat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD statbuf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfstat(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yfstat64(t *TLS, fd int32, statbuf uintptr) (_3 int32)
TEXT ·Yfstat64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD statbuf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfstat64(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yfstatat(t *TLS, dirfd int32, pathname, statbuf uintptr, flags int32) (_4 int32)
TEXT ·Yfstatat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD statbuf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xfstatat(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Yfstatfs(t *TLS, fd int32, buf uintptr) (_3 int32)
TEXT ·Yfstatfs(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfstatfs(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yfsync(t *TLS, fd int32) (_2 int32)
TEXT ·Yfsync(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xfsync(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yftell(t *TLS, stream uintptr) (_2 long)
TEXT ·Yftell(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xftell(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yftello(t *TLS, stream uintptr) (_2 types.Off_t)
TEXT ·Yftello(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xftello(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yftruncate(t *TLS, fd int32, length types.Off_t) (_3 int32)
TEXT ·Yftruncate(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD length+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xftruncate(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yftrylockfile(t *TLS, filehandle uintptr) (_2 int32)
TEXT ·Yftrylockfile(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filehandle+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xftrylockfile(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
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

// func Yftw(tls *TLS, path uintptr, fn uintptr, fd_limit int32) (r int32)
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
	MOVD fn+16(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_ftw_1(SB)	// Create the closure for calling fn
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

// func Yfunlockfile(t *TLS, filehandle uintptr)
TEXT ·Yfunlockfile(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filehandle+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xfunlockfile(SB)
	RET

// func Yfutimens(t *TLS, fd int32, times uintptr) (_3 int32)
TEXT ·Yfutimens(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD times+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfutimens(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yfutimes(t *TLS, fd int32, tv uintptr) (_3 int32)
TEXT ·Yfutimes(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD tv+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xfutimes(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yfwrite(t *TLS, ptr uintptr, size, nmemb types.Size_t, stream uintptr) (_4 types.Size_t)
TEXT ·Yfwrite(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD nmemb+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD stream+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xfwrite(SB)
	MOVD 48(RSP), R0
	MOVD R0, _4+40(FP)
	RET

// func Ygai_strerror(t *TLS, errcode int32) (_2 uintptr)
TEXT ·Ygai_strerror(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW errcode+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgai_strerror(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ygetaddrinfo(tls *TLS, host uintptr, serv uintptr, hint uintptr, res uintptr) (_5 int32)
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
	MOVW R0, _5+40(FP)
	RET

// func Ygetattrlist(t *TLS, path, attrList, attrBuf uintptr, attrBufSize types.Size_t, options uint32) (_4 int32)
TEXT ·Ygetattrlist(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD attrList+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD attrBuf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD attrBufSize+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW options+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xgetattrlist(SB)
	MOVW 56(RSP), R0
	MOVW R0, _4+48(FP)
	RET

// func Ygetc(t *TLS, stream uintptr) (_2 int32)
TEXT ·Ygetc(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetc(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ygetc_unlocked(t *TLS, stream uintptr) (_2 int32)
TEXT ·Ygetc_unlocked(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetc_unlocked(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ygetcwd(t *TLS, buf uintptr, size types.Size_t) (_3 uintptr)
TEXT ·Ygetcwd(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetcwd(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Ygetegid(t *TLS) (_1 types.Gid_t)
TEXT ·Ygetegid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetegid(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Ygetentropy(t *TLS, buffer uintptr, length size_t) (_3 int32)
TEXT ·Ygetentropy(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buffer+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD length+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetentropy(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ygetenv(t *TLS, name uintptr) (_2 uintptr)
TEXT ·Ygetenv(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetenv(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ygeteuid(t *TLS) (_1 types.Uid_t)
TEXT ·Ygeteuid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgeteuid(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Ygetgid(t *TLS) (_1 types.Gid_t)
TEXT ·Ygetgid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetgid(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Ygetgrgid(t *TLS, gid uint32) (_2 uintptr)
TEXT ·Ygetgrgid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW gid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgetgrgid(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ygetgrgid_r(t *TLS, gid uint32, pGrp, buf uintptr, buflen types.Size_t, result uintptr) (_5 int32)
TEXT ·Ygetgrgid_r(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW gid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pGrp+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buflen+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD result+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetgrgid_r(SB)
	MOVW 56(RSP), R0
	MOVW R0, _5+48(FP)
	RET

// func Ygetgrnam(t *TLS, name uintptr) (_2 uintptr)
TEXT ·Ygetgrnam(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetgrnam(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ygetgrnam_r(t *TLS, name, pGrp, buf uintptr, buflen types.Size_t, result uintptr) (_4 int32)
TEXT ·Ygetgrnam_r(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pGrp+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buflen+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD result+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetgrnam_r(SB)
	MOVW 56(RSP), R0
	MOVW R0, _4+48(FP)
	RET

// func Ygethostbyaddr(tls *TLS, a uintptr, l socklen_t, af int32) (_4 uintptr)
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
	MOVD R0, _4+24(FP)
	RET

// func Ygethostbyaddr_r(tls *TLS, a uintptr, l socklen_t, af int32, h uintptr, buf uintptr, buflen size_t, res uintptr, err uintptr) (_9 int32)
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
	MOVW R0, _9+64(FP)
	RET

// func Ygethostbyname(tls *TLS, name uintptr) (_2 uintptr)
TEXT ·Ygethostbyname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgethostbyname(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ygethostbyname2(tls *TLS, name uintptr, af int32) (_3 uintptr)
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
	MOVD R0, _3+24(FP)
	RET

// func Ygethostbyname2_r(tls *TLS, name uintptr, af int32, h uintptr, buf uintptr, buflen size_t, res uintptr, err uintptr) (_8 int32)
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
	MOVW R0, _8+64(FP)
	RET

// func Ygethostname(t *TLS, name uintptr, slen types.Size_t) (_3 int32)
TEXT ·Ygethostname(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD slen+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgethostname(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ygethostuuid(t *TLS, id uintptr, wait uintptr) (_3 int32)
TEXT ·Ygethostuuid(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD id+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD wait+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgethostuuid(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ygetloadavg(t *TLS, loadavg uintptr, nelem int32) (_3 int32)
TEXT ·Ygetloadavg(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD loadavg+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW nelem+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xgetloadavg(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ygetlogin_r(t *TLS, buf uintptr, bufsize types.Size_t) (_3 int32)
TEXT ·Ygetlogin_r(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD bufsize+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetlogin_r(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ygetnameinfo(tls *TLS, sa1 uintptr, sl socklen_t, node uintptr, nodelen socklen_t, serv uintptr, servlen socklen_t, flags int32) (_8 int32)
TEXT ·Ygetnameinfo(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD sa1+8(FP), R0
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
	MOVW R0, _8+56(FP)
	RET

// func Ygetpeername(t *TLS, sockfd int32, addr uintptr, addrlen uintptr) (_4 int32)
TEXT ·Ygetpeername(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD addrlen+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetpeername(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Ygetpid(t *TLS) (_1 int32)
TEXT ·Ygetpid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetpid(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Ygetprogname(t *TLS) (_1 uintptr)
TEXT ·Ygetprogname(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetprogname(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Ygetpwnam(t *TLS, name uintptr) (_2 uintptr)
TEXT ·Ygetpwnam(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xgetpwnam(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ygetpwnam_r(t *TLS, name, cpwd, buf uintptr, buflen types.Size_t, result uintptr) (_4 int32)
TEXT ·Ygetpwnam_r(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD cpwd+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buflen+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD result+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetpwnam_r(SB)
	MOVW 56(RSP), R0
	MOVW R0, _4+48(FP)
	RET

// func Ygetpwuid(t *TLS, uid uint32) (_2 uintptr)
TEXT ·Ygetpwuid(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW uid+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgetpwuid(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ygetpwuid_r(t *TLS, uid types.Uid_t, cpwd, buf uintptr, buflen types.Size_t, result uintptr) (_5 int32)
TEXT ·Ygetpwuid_r(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW uid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD cpwd+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD buflen+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD result+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xgetpwuid_r(SB)
	MOVW 56(RSP), R0
	MOVW R0, _5+48(FP)
	RET

// func Ygetresgid(t *TLS, rgid, egid, sgid uintptr) (_2 int32)
TEXT ·Ygetresgid(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD rgid+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD egid+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD sgid+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetresgid(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Ygetresuid(t *TLS, ruid, euid, suid uintptr) (_2 int32)
TEXT ·Ygetresuid(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ruid+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD euid+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD suid+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetresuid(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Ygetrlimit(t *TLS, resource int32, rlim uintptr) (_3 int32)
TEXT ·Ygetrlimit(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW resource+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD rlim+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetrlimit(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ygetrusage(t *TLS, who int32, usage uintptr) (_3 int32)
TEXT ·Ygetrusage(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW who+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD usage+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetrusage(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ygetservbyname(t *TLS, name, proto uintptr) (_2 uintptr)
TEXT ·Ygetservbyname(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD proto+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgetservbyname(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ygetsockname(t *TLS, sockfd int32, addr, addrlen uintptr) (_3 int32)
TEXT ·Ygetsockname(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD addr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD addrlen+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xgetsockname(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ygetsockopt(t *TLS, sockfd, level, optname int32, optval, optlen uintptr) (_3 int32)
TEXT ·Ygetsockopt(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
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
	MOVW R0, _3+40(FP)
	RET

// func Ygettimeofday(t *TLS, tv, tz uintptr) (_2 int32)
TEXT ·Ygettimeofday(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tv+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD tz+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xgettimeofday(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ygetuid(t *TLS) (_1 types.Uid_t)
TEXT ·Ygetuid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xgetuid(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
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

// func Ygrantpt(t *TLS, fd int32) (_2 int32)
TEXT ·Ygrantpt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xgrantpt(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yhtonl(t *TLS, hostlong uint32) (_2 uint32)
TEXT ·Yhtonl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW hostlong+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xhtonl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yhtons(t *TLS, hostshort uint16) (_2 uint16)
TEXT ·Yhtons(SB),$24-18
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVH hostshort+8(FP), R0
	MOVH R0, 16(RSP)
	CALL ·Xhtons(SB)
	MOVH 24(RSP), R0
	MOVH R0, _2+16(FP)
	RET

// func Yhypot(t *TLS, x, y float64) (_2 float64)
TEXT ·Yhypot(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xhypot(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Yhypotf(t *TLS, x, y float32) (_2 float32)
TEXT ·Yhypotf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xhypotf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yilogb(t *TLS, x float64) (_2 int32)
TEXT ·Yilogb(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xilogb(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yinet_ntoa(t *TLS, in1 in.In_addr) (_2 uintptr)
TEXT ·Yinet_ntoa(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW in1_Fs_addr+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xinet_ntoa(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yinet_ntop(tls *TLS, af int32, a0 uintptr, s uintptr, l socklen_t) (_5 uintptr)
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
	MOVD R0, _5+40(FP)
	RET

// func Yinet_pton(tls *TLS, af int32, s uintptr, a0 uintptr) (_4 int32)
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
	MOVW R0, _4+32(FP)
	RET

// func Yinitstate(t *TLS, seed uint32, statebuf uintptr, statelen types.Size_t) (_4 uintptr)
TEXT ·Yinitstate(SB),$40-40
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
	CALL ·Xinitstate(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yinitstate_r(t *TLS, seed uint32, statebuf uintptr, statelen types.Size_t, buf uintptr) (_5 int32)
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

// func Yioctl(t *TLS, fd int32, request ulong, va uintptr) (_4 int32)
TEXT ·Yioctl(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD request+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xioctl(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Yisalnum(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisalnum(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisalnum(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisalpha(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisalpha(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisalpha(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisascii(t *TLS, c int32) (_2 int32)
TEXT ·Yisascii(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisascii(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisatty(t *TLS, fd int32) (_2 int32)
TEXT ·Yisatty(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisatty(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisblank(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisblank(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisblank(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yiscntrl(tls *TLS, _c int32) (_2 int32)
TEXT ·Yiscntrl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiscntrl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisdigit(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisdigit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisdigit(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisgraph(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisgraph(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisgraph(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yishexnumber(tls *TLS, _c int32) (_2 int32)
TEXT ·Yishexnumber(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xishexnumber(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisideogram(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisideogram(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisideogram(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yislower(tls *TLS, _c int32) (_2 int32)
TEXT ·Yislower(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xislower(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
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

// func Yisnanf(t *TLS, x float32) (_2 int32)
TEXT ·Yisnanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisnanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisnanl(t *TLS, x float64) (_2 int32)
TEXT ·Yisnanl(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xisnanl(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisnumber(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisnumber(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisnumber(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisphonogram(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisphonogram(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisphonogram(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisprint(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisprint(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisprint(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yispunct(tls *TLS, _c int32) (_2 int32)
TEXT ·Yispunct(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xispunct(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisrune(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisrune(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisrune(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yissetugid(t *TLS) (_1 int32)
TEXT ·Yissetugid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xissetugid(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Yisspace(t *TLS, c int32) (_2 int32)
TEXT ·Yisspace(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisspace(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisspecial(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisspecial(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisspecial(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisupper(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisupper(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisupper(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yiswalnum(t *TLS, wc wctype.Wint_t) (_2 int32)
TEXT ·Yiswalnum(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswalnum(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yiswspace(t *TLS, wc wctype.Wint_t) (_2 int32)
TEXT ·Yiswspace(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW wc+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xiswspace(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yisxdigit(tls *TLS, _c int32) (_2 int32)
TEXT ·Yisxdigit(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xisxdigit(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ykevent(t *TLS, kq int32, changelist uintptr, nchanges int32, eventlist uintptr, nevents int32, timeout uintptr) (_7 int32)
TEXT ·Ykevent(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW kq+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD changelist+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW nchanges+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD eventlist+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW nevents+40(FP), R0
	MOVW R0, 48(RSP)
	MOVD timeout+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xkevent(SB)
	MOVW 64(RSP), R0
	MOVW R0, _7+56(FP)
	RET

// func Ykill(t *TLS, pid types.Pid_t, sig int32) (_3 int32)
TEXT ·Ykill(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW sig+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xkill(SB)
	MOVW 24(RSP), R0
	MOVW R0, _3+16(FP)
	RET

// func Ykqueue(t *TLS) (_1 int32)
TEXT ·Ykqueue(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xkqueue(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Ylabs(t *TLS, j long) (_2 long)
TEXT ·Ylabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD j+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ylchmod(t *TLS, path uintptr, mode types.Mode_t) (_3 int32)
TEXT ·Ylchmod(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH mode+16(FP), R0
	MOVH R0, 24(RSP)
	CALL ·Xlchmod(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yldexp(t *TLS, x float64, exp int32) (_3 float64)
TEXT ·Yldexp(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW exp+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xldexp(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Yldiv(tls *TLS, num, den long) (r Tldiv_t)
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

// func Ylgamma(t *TLS, x float64) (_2 float64)
TEXT ·Ylgamma(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlgamma(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ylink(t *TLS, oldpath, newpath uintptr) (_2 int32)
TEXT ·Ylink(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD oldpath+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD newpath+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlink(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ylinkat(t *TLS, olddirfd int32, oldpath uintptr, newdirfd int32, newpath uintptr, flags int32) (_6 int32)
TEXT ·Ylinkat(SB),$56-52
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
	CALL ·Xlinkat(SB)
	MOVW 56(RSP), R0
	MOVW R0, _6+48(FP)
	RET

// func Ylisten(t *TLS, sockfd, backlog int32) (_2 int32)
TEXT ·Ylisten(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW backlog+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xlisten(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yllabs(tls *TLS, a int64) (_2 int64)
TEXT ·Yllabs(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xllabs(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ylocaleconv(tls *TLS) (_1 uintptr)
TEXT ·Ylocaleconv(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xlocaleconv(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Ylocaltime(_ *TLS, timep uintptr) (_2 uintptr)
TEXT ·Ylocaltime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD timep+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlocaltime(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ylocaltime_r(_ *TLS, timep, result uintptr) (_2 uintptr)
TEXT ·Ylocaltime_r(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD _+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD timep+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD result+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlocaltime_r(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ylog(t *TLS, x float64) (_2 float64)
TEXT ·Ylog(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ylog10(t *TLS, x float64) (_2 float64)
TEXT ·Ylog10(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog10(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ylog10f(t *TLS, x float32) (_2 float32)
TEXT ·Ylog10f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlog10f(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
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

// func Ylog2(t *TLS, x float64) (_2 float64)
TEXT ·Ylog2(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlog2(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ylog2f(t *TLS, x float32) (_2 float32)
TEXT ·Ylog2f(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlog2f(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ylogb(t *TLS, x float64) (_2 float64)
TEXT ·Ylogb(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xlogb(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ylogf(t *TLS, x float32) (_2 float32)
TEXT ·Ylogf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xlogf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
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

// func Ylrand48(t *TLS) (_1 long)
TEXT ·Ylrand48(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xlrand48(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Ylrint(tls *TLS, x float64) (r long)
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

// func Ylrintf(tls *TLS, x float32) (r long)
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

// func Ylround(tls *TLS, x float64) (r long)
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

// func Ylseek(t *TLS, fd int32, offset types.Off_t, whence int32) (_4 types.Off_t)
TEXT ·Ylseek(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD offset+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xlseek(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Ylseek64(t *TLS, fd int32, offset types.Off_t, whence int32) (_4 types.Off_t)
TEXT ·Ylseek64(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD offset+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW whence+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xlseek64(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Ylstat(t *TLS, pathname, statbuf uintptr) (_2 int32)
TEXT ·Ylstat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD statbuf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlstat(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ylstat64(t *TLS, pathname, statbuf uintptr) (_2 int32)
TEXT ·Ylstat64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD statbuf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xlstat64(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ymach_absolute_time(t *TLS) (_1 uint64)
TEXT ·Ymach_absolute_time(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xmach_absolute_time(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Ymach_timebase_info(t *TLS, info uintptr) (_2 int32)
TEXT ·Ymach_timebase_info(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD info+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmach_timebase_info(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ymalloc(t *TLS, n types.Size_t) (_2 uintptr)
TEXT ·Ymalloc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD n+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmalloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ymalloc_size(t *TLS, p uintptr) (r types.Size_t)
TEXT ·Ymalloc_size(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD p+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmalloc_size(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ymalloc_usable_size(tls *TLS, p uintptr) (r types.Size_t)
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

// func Ymblen(t *TLS, s uintptr, n types.Size_t) (_3 int32)
TEXT ·Ymblen(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD n+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmblen(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ymbstowcs(t *TLS, dest, src uintptr, n types.Size_t) (_3 types.Size_t)
TEXT ·Ymbstowcs(SB),$40-40
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
	CALL ·Xmbstowcs(SB)
	MOVD 40(RSP), R0
	MOVD R0, _3+32(FP)
	RET

// func Ymbtowc(t *TLS, pwc, s uintptr, n types.Size_t) (_3 int32)
TEXT ·Ymbtowc(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pwc+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmbtowc(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ymemchr(t *TLS, s uintptr, c int32, n types.Size_t) (_4 uintptr)
TEXT ·Ymemchr(SB),$40-40
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
	CALL ·Xmemchr(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Ymemcmp(t *TLS, s1, s2 uintptr, n types.Size_t) (_3 int32)
TEXT ·Ymemcmp(SB),$40-36
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
	CALL ·Xmemcmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ymemcpy(t *TLS, dest, src uintptr, n types.Size_t) (r uintptr)
TEXT ·Ymemcpy(SB),$40-40
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
	CALL ·Xmemcpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ymemmove(t *TLS, dest, src uintptr, n types.Size_t) (_3 uintptr)
TEXT ·Ymemmove(SB),$40-40
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
	CALL ·Xmemmove(SB)
	MOVD 40(RSP), R0
	MOVD R0, _3+32(FP)
	RET

// func Ymemset(t *TLS, dest uintptr, c int32, n types.Size_t) (_4 uintptr)
TEXT ·Ymemset(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD n+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xmemset(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Ymkdir(t *TLS, path uintptr, mode types.Mode_t) (_3 int32)
TEXT ·Ymkdir(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH mode+16(FP), R0
	MOVH R0, 24(RSP)
	CALL ·Xmkdir(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ymkdirat(t *TLS, dirfd int32, pathname uintptr, mode types.Mode_t) (_4 int32)
TEXT ·Ymkdirat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVH mode+24(FP), R0
	MOVH R0, 32(RSP)
	CALL ·Xmkdirat(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Ymkfifo(t *TLS, pathname uintptr, mode types.Mode_t) (_3 int32)
TEXT ·Ymkfifo(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH mode+16(FP), R0
	MOVH R0, 24(RSP)
	CALL ·Xmkfifo(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ymknod(t *TLS, pathname uintptr, mode types.Mode_t, dev types.Dev_t) (_4 int32)
TEXT ·Ymknod(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVH mode+16(FP), R0
	MOVH R0, 24(RSP)
	MOVW dev+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xmknod(SB)
	MOVW 32(RSP), R0
	MOVW R0, _4+24(FP)
	RET

// func Ymkostemp(t *TLS, template uintptr, flags int32) (_3 int32)
TEXT ·Ymkostemp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmkostemp(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ymkstemp(t *TLS, template uintptr) (_2 int32)
TEXT ·Ymkstemp(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmkstemp(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ymkstemp64(t *TLS, template uintptr) (_2 int32)
TEXT ·Ymkstemp64(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmkstemp64(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ymkstemps(t *TLS, template uintptr, suffixlen int32) (_3 int32)
TEXT ·Ymkstemps(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW suffixlen+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmkstemps(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ymkstemps64(t *TLS, template uintptr, suffixlen int32) (_3 int32)
TEXT ·Ymkstemps64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD template+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW suffixlen+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xmkstemps64(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ymktime(t *TLS, ptm uintptr) (r time.Time_t)
TEXT ·Ymktime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptm+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xmktime(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ymmap(t *TLS, addr uintptr, length types.Size_t, prot, flags, fd int32, offset types.Off_t) (_5 uintptr)
TEXT ·Ymmap(SB),$56-56
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
	CALL ·Xmmap(SB)
	MOVD 56(RSP), R0
	MOVD R0, _5+48(FP)
	RET

// func Ymodf(t *TLS, x float64, iptr uintptr) (_3 float64)
TEXT ·Ymodf(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD iptr+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmodf(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Ymodff(t *TLS, x float32, iptr uintptr) (_3 float32)
TEXT ·Ymodff(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iptr+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmodff(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ymunmap(t *TLS, addr uintptr, length types.Size_t) (_3 int32)
TEXT ·Ymunmap(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD addr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD length+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xmunmap(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ynanf(tls *TLS, s uintptr) (_2 float32)
TEXT ·Ynanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xnanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ynanosleep(t *TLS, req, rem uintptr) (_2 int32)
TEXT ·Ynanosleep(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD req+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD rem+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnanosleep(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ynextafter(t *TLS, x, y float64) (_2 float64)
TEXT ·Ynextafter(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnextafter(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ynextafterf(t *TLS, x, y float32) (_2 float32)
TEXT ·Ynextafterf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xnextafterf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ynexttoward(t *TLS, x, y float64) (_2 float64)
TEXT ·Ynexttoward(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xnexttoward(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ynl_langinfo(t *TLS, item langinfo.Nl_item) (_2 uintptr)
TEXT ·Ynl_langinfo(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW item+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xnl_langinfo(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yntohs(t *TLS, netshort uint16) (_2 uint16)
TEXT ·Yntohs(SB),$24-18
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVH netshort+8(FP), R0
	MOVH R0, 16(RSP)
	CALL ·Xntohs(SB)
	MOVH 24(RSP), R0
	MOVH R0, _2+16(FP)
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

// func Yopen(t *TLS, pathname uintptr, flags int32, args uintptr) (_4 int32)
TEXT ·Yopen(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD args+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xopen(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Yopen64(t *TLS, pathname uintptr, flags int32, args uintptr) (_4 int32)
TEXT ·Yopen64(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW flags+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD args+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xopen64(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Yopenat(t *TLS, dirfd int32, pathname uintptr, flags int32, args uintptr) (_5 int32)
TEXT ·Yopenat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD args+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xopenat(SB)
	MOVW 48(RSP), R0
	MOVW R0, _5+40(FP)
	RET

// func Yopendir(t *TLS, name uintptr) (_2 uintptr)
TEXT ·Yopendir(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xopendir(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yopenpty(t *TLS, amaster, aslave, name, termp, winp uintptr) (_2 int32)
TEXT ·Yopenpty(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD amaster+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD aslave+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD name+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD termp+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD winp+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xopenpty(SB)
	MOVW 56(RSP), R0
	MOVW R0, _2+48(FP)
	RET

// func Ypathconf(t *TLS, path uintptr, name int32) (_3 long)
TEXT ·Ypathconf(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW name+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpathconf(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Ypause(t *TLS) (_1 int32)
TEXT ·Ypause(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xpause(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Ypclose(t *TLS, stream uintptr) (r int32)
TEXT ·Ypclose(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpclose(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
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

// func Ypipe(t *TLS, pipefd uintptr) (_2 int32)
TEXT ·Ypipe(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pipefd+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpipe(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypoll(t *TLS, fds uintptr, nfds poll.Nfds_t, timeout int32) (_4 int32)
TEXT ·Ypoll(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD fds+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW nfds+16(FP), R0
	MOVW R0, 24(RSP)
	MOVW timeout+20(FP), R0
	MOVW R0, 28(RSP)
	CALL ·Xpoll(SB)
	MOVW 32(RSP), R0
	MOVW R0, _4+24(FP)
	RET

// func Ypopen(t *TLS, command, mode uintptr) (r uintptr)
TEXT ·Ypopen(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD command+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD mode+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpopen(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Yposix_fadvise(t *TLS, fd int32, offset, len types.Off_t, advice int32) (_4 int32)
TEXT ·Yposix_fadvise(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD offset+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW advice+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xposix_fadvise(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Yposix_openpt(t *TLS, flags int32) (_2 int32)
TEXT ·Yposix_openpt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW flags+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xposix_openpt(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypow(t *TLS, x, y float64) (_2 float64)
TEXT ·Ypow(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpow(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ypowf(t *TLS, x, y float32) (_2 float32)
TEXT ·Ypowf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW y+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xpowf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypread(t *TLS, fd int32, buf uintptr, count types.Size_t, offset types.Off_t) (_5 types.Ssize_t)
TEXT ·Ypread(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD count+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD offset+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xpread(SB)
	MOVD 48(RSP), R0
	MOVD R0, _5+40(FP)
	RET

// func Yprintf(t *TLS, format, args uintptr) (_2 int32)
TEXT ·Yprintf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD format+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD args+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xprintf(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypselect(t *TLS, nfds int32, readfds, writefds, exceptfds, timeout, sigmask uintptr) (_3 int32)
TEXT ·Ypselect(SB),$64-60
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW nfds+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD readfds+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD writefds+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD exceptfds+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD timeout+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD sigmask+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xpselect(SB)
	MOVW 64(RSP), R0
	MOVW R0, _3+56(FP)
	RET

// func Ypthread_attr_destroy(t *TLS, pAttr uintptr) (_2 int32)
TEXT ·Ypthread_attr_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pAttr+8(FP), R0
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

// func Ypthread_attr_init(t *TLS, pAttr uintptr) (_2 int32)
TEXT ·Ypthread_attr_init(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pAttr+8(FP), R0
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

// func Ypthread_attr_setscope(t *TLS, pAttr uintptr, contentionScope int32) (_3 int32)
TEXT ·Ypthread_attr_setscope(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pAttr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW contentionScope+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpthread_attr_setscope(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_attr_setstacksize(t *TLS, attr uintptr, stackSize types.Size_t) (_3 int32)
TEXT ·Ypthread_attr_setstacksize(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD attr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD stackSize+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_attr_setstacksize(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_cond_broadcast(t *TLS, pCond uintptr) (_2 int32)
TEXT ·Ypthread_cond_broadcast(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pCond+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_cond_broadcast(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_cond_destroy(t *TLS, pCond uintptr) (_2 int32)
TEXT ·Ypthread_cond_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pCond+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_cond_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_cond_init(t *TLS, pCond, pAttr uintptr) (_2 int32)
TEXT ·Ypthread_cond_init(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pCond+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pAttr+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_cond_init(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypthread_cond_signal(t *TLS, pCond uintptr) (_2 int32)
TEXT ·Ypthread_cond_signal(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pCond+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_cond_signal(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_cond_timedwait(t *TLS, pCond, pMutex, pAbsTime uintptr) (_2 int32)
TEXT ·Ypthread_cond_timedwait(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pCond+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pMutex+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD pAbsTime+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xpthread_cond_timedwait(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Ypthread_cond_wait(t *TLS, pCond, pMutex uintptr) (_2 int32)
TEXT ·Ypthread_cond_wait(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pCond+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pMutex+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_cond_wait(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypthread_detach(t *TLS, thread pthread.Pthread_t) (_2 int32)
TEXT ·Ypthread_detach(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD thread+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_detach(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_equal(t *TLS, t1, t2 pthread.Pthread_t) (_2 int32)
TEXT ·Ypthread_equal(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD t1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD t2+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_equal(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypthread_exit(t *TLS, value uintptr)
TEXT ·Ypthread_exit(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD value+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_exit(SB)
	RET

// func Ypthread_getname_np(t *TLS, thread pthread.Pthread_t, name uintptr, len1 types.Size_t) (_4 int32)
TEXT ·Ypthread_getname_np(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
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

// func Ypthread_getspecific(t *TLS, key pthread.Pthread_key_t) (_2 uintptr)
TEXT ·Ypthread_getspecific(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_getspecific(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ypthread_join(t *TLS, thread pthread.Pthread_t, pValue uintptr) (_3 int32)
TEXT ·Ypthread_join(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD thread+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pValue+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_join(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_key_delete(t *TLS, key pthread.Pthread_key_t) (_2 int32)
TEXT ·Ypthread_key_delete(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_key_delete(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutex_destroy(t *TLS, pMutex uintptr) (_2 int32)
TEXT ·Ypthread_mutex_destroy(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pMutex+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutex_destroy(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutex_init(t *TLS, pMutex, pAttr uintptr) (_2 int32)
TEXT ·Ypthread_mutex_init(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pMutex+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD pAttr+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xpthread_mutex_init(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ypthread_mutex_lock(t *TLS, pMutex uintptr) (_2 int32)
TEXT ·Ypthread_mutex_lock(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pMutex+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutex_lock(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutex_trylock(t *TLS, pMutex uintptr) (_2 int32)
TEXT ·Ypthread_mutex_trylock(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pMutex+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_mutex_trylock(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_mutex_unlock(t *TLS, pMutex uintptr) (_2 int32)
TEXT ·Ypthread_mutex_unlock(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pMutex+8(FP), R0
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

// func Ypthread_mutexattr_settype(tls *TLS, a uintptr, type1 int32) (_3 int32)
TEXT ·Ypthread_mutexattr_settype(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD a+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW type1+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xpthread_mutexattr_settype(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ypthread_self(t *TLS) (_1 pthread.Pthread_t)
TEXT ·Ypthread_self(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xpthread_self(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Ypthread_setname_np(t *TLS, name uintptr) (_2 int32)
TEXT ·Ypthread_setname_np(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xpthread_setname_np(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypthread_setspecific(t *TLS, key pthread.Pthread_key_t, value uintptr) (_3 int32)
TEXT ·Ypthread_setspecific(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD key+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD value+16(FP), R0
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

// func Yptsname(t *TLS, fd int32) (_2 uintptr)
TEXT ·Yptsname(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xptsname(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yptsname_r(t *TLS, fd int32, buf uintptr, buflen types.Size_t) (_4 int32)
TEXT ·Yptsname_r(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buflen+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xptsname_r(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Yputc(t *TLS, c int32, fp uintptr) (_3 int32)
TEXT ·Yputc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD fp+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputc(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yputc_unlocked(t *TLS, c int32, stream uintptr) (_3 int32)
TEXT ·Yputc_unlocked(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD stream+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xputc_unlocked(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yputchar(t *TLS, c int32) (_2 int32)
TEXT ·Yputchar(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xputchar(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yputenv(tls *TLS, s uintptr) (_2 int32)
TEXT ·Yputenv(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xputenv(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yputs(t *TLS, s uintptr) (_2 int32)
TEXT ·Yputs(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xputs(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ypwrite(t *TLS, fd int32, buf uintptr, count types.Size_t, offset types.Off_t) (_5 types.Ssize_t)
TEXT ·Ypwrite(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD count+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD offset+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xpwrite(SB)
	MOVD 48(RSP), R0
	MOVD R0, _5+40(FP)
	RET

// func Yqsort(t *TLS, base uintptr, nmemb, size types.Size_t, compar uintptr)
TEXT ·Yqsort(SB),$48-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0	// alloc all ABI trampolines
	MOVD R0, 8(RSP)
	MOVD $16, R0
	MOVD R0, 16(RSP)	// 16*(number of func ptrs in signature)
	CALL modernc·org∕libc·TLSAlloc(SB)
	MOVD 24(RSP), R0
	MOVD R0, 48(RSP) // Trampolines[0]
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD 48(RSP), R0
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 16(RSP)
	MOVD compar+32(FP), R0	// ABI0 code ptr
	MOVD R0, 24(RSP)
	CALL ·__ccgo_abiInternal_qsort_3(SB)	// Create the closure for calling compar
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD base+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD nmemb+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD 48(RSP), R0 // Trampolines[0]
	ADD $0, R0	// 16*(0-based ordinal number of the func ptr in signature)
	MOVD R0, 40(RSP)
	CALL ·Xqsort(SB)
	MOVD t+0(FP), R0
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

// func Yraise(t *TLS, sig int32) (_2 int32)
TEXT ·Yraise(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sig+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xraise(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yrand(t *TLS) (_1 int32)
TEXT ·Yrand(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xrand(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Yrand_r(tls *TLS, seed uintptr) (_2 int32)
TEXT ·Yrand_r(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD seed+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrand_r(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yrandom(t *TLS) (_1 long)
TEXT ·Yrandom(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xrandom(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
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

// func Yread(t *TLS, fd int32, buf uintptr, count types.Size_t) (_4 types.Ssize_t)
TEXT ·Yread(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD count+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xread(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yreaddir(t *TLS, dir uintptr) (_2 uintptr)
TEXT ·Yreaddir(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dir+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xreaddir(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yreadlink(t *TLS, path, buf uintptr, bufsize types.Size_t) (_3 types.Ssize_t)
TEXT ·Yreadlink(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD bufsize+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xreadlink(SB)
	MOVD 40(RSP), R0
	MOVD R0, _3+32(FP)
	RET

// func Yreadlinkat(t *TLS, dirfd int32, pathname, buf uintptr, bufsiz types.Size_t) (_4 types.Ssize_t)
TEXT ·Yreadlinkat(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD bufsiz+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xreadlinkat(SB)
	MOVD 48(RSP), R0
	MOVD R0, _4+40(FP)
	RET

// func Yreadv(t *TLS, fd int32, iov uintptr, iovcnt int32) (_4 types.Ssize_t)
TEXT ·Yreadv(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW iovcnt+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xreadv(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yrealloc(t *TLS, ptr uintptr, size types.Size_t) (_3 uintptr)
TEXT ·Yrealloc(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD size+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xrealloc(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Yreallocarray(t *TLS, ptr uintptr, nmemb, size size_t) (_3 uintptr)
TEXT ·Yreallocarray(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD ptr+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD nmemb+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD size+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xreallocarray(SB)
	MOVD 40(RSP), R0
	MOVD R0, _3+32(FP)
	RET

// func Yrealpath(t *TLS, path, resolved_path uintptr) (_2 uintptr)
TEXT ·Yrealpath(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD resolved_path+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xrealpath(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Yrecv(t *TLS, sockfd int32, buf uintptr, len types.Size_t, flags int32) (_5 types.Ssize_t)
TEXT ·Yrecv(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xrecv(SB)
	MOVD 48(RSP), R0
	MOVD R0, _5+40(FP)
	RET

// func Yrecvfrom(t *TLS, sockfd int32, buf uintptr, len types.Size_t, flags int32, src_addr, addrlen uintptr) (_6 types.Ssize_t)
TEXT ·Yrecvfrom(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD src_addr+40(FP), R0
	MOVD R0, 48(RSP)
	MOVD addrlen+48(FP), R0
	MOVD R0, 56(RSP)
	CALL ·Xrecvfrom(SB)
	MOVD 64(RSP), R0
	MOVD R0, _6+56(FP)
	RET

// func Yrecvmsg(t *TLS, sockfd int32, msg uintptr, flags int32) (_4 types.Ssize_t)
TEXT ·Yrecvmsg(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD msg+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xrecvmsg(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yremainder(t *TLS, x, y float64) (_2 float64)
TEXT ·Yremainder(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD y+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xremainder(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Yremove(t *TLS, pathname uintptr) (_2 int32)
TEXT ·Yremove(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xremove(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yremquo(t *TLS, x, y float64, quo uintptr) (_3 float64)
TEXT ·Yremquo(SB),$40-40
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
	CALL ·Xremquo(SB)
	MOVD 40(RSP), R0
	MOVD R0, _3+32(FP)
	RET

// func Yrename(t *TLS, oldpath, newpath uintptr) (_2 int32)
TEXT ·Yrename(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD oldpath+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD newpath+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xrename(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Yrenameat(t *TLS, olddirfd int32, oldpath uintptr, newdirfd int32, newpath uintptr) (_5 int32)
TEXT ·Yrenameat(SB),$48-44
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
	CALL ·Xrenameat(SB)
	MOVW 48(RSP), R0
	MOVW R0, _5+40(FP)
	RET

// func Yrewind(t *TLS, stream uintptr)
TEXT ·Yrewind(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrewind(SB)
	RET

// func Yrewinddir(tls *TLS, f uintptr)
TEXT ·Yrewinddir(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD f+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrewinddir(SB)
	RET

// func Yrindex(t *TLS, s uintptr, c int32) (_3 uintptr)
TEXT ·Yrindex(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xrindex(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Yrint(tls *TLS, x float64) (_2 float64)
TEXT ·Yrint(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrint(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
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

// func Yrmdir(t *TLS, pathname uintptr) (_2 int32)
TEXT ·Yrmdir(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xrmdir(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yround(t *TLS, x float64) (_2 float64)
TEXT ·Yround(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xround(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Yroundf(t *TLS, x float32) (_2 float32)
TEXT ·Yroundf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xroundf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yscalbn(tls *TLS, x float64, n int32) (_3 float64)
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
	MOVD R0, _3+24(FP)
	RET

// func Yscalbnl(tls *TLS, x float64, n int32) (_3 float64)
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
	MOVD R0, _3+24(FP)
	RET

// func Ysched_yield(t *TLS) (_1 int32)
TEXT ·Ysched_yield(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsched_yield(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Yselect(t *TLS, nfds int32, readfds, writefds, exceptfds, timeout uintptr) (_3 int32)
TEXT ·Yselect(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW nfds+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD readfds+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD writefds+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD exceptfds+32(FP), R0
	MOVD R0, 40(RSP)
	MOVD timeout+40(FP), R0
	MOVD R0, 48(RSP)
	CALL ·Xselect(SB)
	MOVW 56(RSP), R0
	MOVW R0, _3+48(FP)
	RET

// func Ysend(t *TLS, sockfd int32, buf uintptr, len types.Size_t, flags int32) (_5 types.Ssize_t)
TEXT ·Ysend(SB),$48-48
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xsend(SB)
	MOVD 48(RSP), R0
	MOVD R0, _5+40(FP)
	RET

// func Ysendmsg(t *TLS, sockfd int32, msg uintptr, flags int32) (_4 types.Ssize_t)
TEXT ·Ysendmsg(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD msg+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xsendmsg(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Ysendto(t *TLS, sockfd int32, buf uintptr, len types.Size_t, flags int32, src_addr uintptr, addrlen socklen_t) (_7 types.Ssize_t)
TEXT ·Ysendto(SB),$64-64
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD len+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	MOVD src_addr+40(FP), R0
	MOVD R0, 48(RSP)
	MOVW addrlen+48(FP), R0
	MOVW R0, 56(RSP)
	CALL ·Xsendto(SB)
	MOVD 64(RSP), R0
	MOVD R0, _7+56(FP)
	RET

// func Ysetattrlist(t *TLS, path, attrList, attrBuf uintptr, attrBufSize types.Size_t, options uint32) (_4 int32)
TEXT ·Ysetattrlist(SB),$56-52
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD attrList+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD attrBuf+24(FP), R0
	MOVD R0, 32(RSP)
	MOVD attrBufSize+32(FP), R0
	MOVD R0, 40(RSP)
	MOVW options+40(FP), R0
	MOVW R0, 48(RSP)
	CALL ·Xsetattrlist(SB)
	MOVW 56(RSP), R0
	MOVW R0, _4+48(FP)
	RET

// func Ysetbuf(t *TLS, stream, buf uintptr)
TEXT ·Ysetbuf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetbuf(SB)
	RET

// func Ysetenv(tls *TLS, var1 uintptr, value uintptr, overwrite int32) (_4 int32)
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
	MOVW R0, _4+32(FP)
	RET

// func Ysetgid(tls *TLS, gid uint32) (r int32)
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

// func Ysetlocale(t *TLS, category int32, name uintptr) (r uintptr)
TEXT ·Ysetlocale(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW category+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD name+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetlocale(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ysetrlimit(t *TLS, resource int32, rlim uintptr) (_3 int32)
TEXT ·Ysetrlimit(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW resource+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD rlim+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsetrlimit(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ysetsid(t *TLS) (_1 types.Pid_t)
TEXT ·Ysetsid(SB),$16-12
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsetsid(SB)
	MOVW 16(RSP), R0
	MOVW R0, _1+8(FP)
	RET

// func Ysetsockopt(t *TLS, sockfd, level, optname int32, optval uintptr, optlen socket.Socklen_t) (_4 int32)
TEXT ·Ysetsockopt(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
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
	MOVW R0, _4+40(FP)
	RET

// func Ysetstate(t *TLS, state uintptr) (_2 uintptr)
TEXT ·Ysetstate(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD state+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsetstate(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ysetuid(tls *TLS, uid uint32) (r int32)
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

// func Ysetvbuf(t *TLS, stream, buf uintptr, mode int32, size types.Size_t) (_4 int32)
TEXT ·Ysetvbuf(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW mode+24(FP), R0
	MOVW R0, 32(RSP)
	MOVD size+32(FP), R0
	MOVD R0, 40(RSP)
	CALL ·Xsetvbuf(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Yshmat(t *TLS, shmid int32, shmaddr uintptr, shmflg int32) (_4 uintptr)
TEXT ·Yshmat(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW shmid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD shmaddr+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW shmflg+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xshmat(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Yshmctl(t *TLS, shmid, cmd int32, buf uintptr) (_3 int32)
TEXT ·Yshmctl(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW shmid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW cmd+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xshmctl(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yshmdt(t *TLS, shmaddr uintptr) (_2 int32)
TEXT ·Yshmdt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD shmaddr+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xshmdt(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yshutdown(t *TLS, sockfd, how int32) (_2 int32)
TEXT ·Yshutdown(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW sockfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW how+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xshutdown(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ysigaddset(t *TLS, set uintptr, signum int32) (_3 int32)
TEXT ·Ysigaddset(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW signum+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsigaddset(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ysigdelset(t *TLS, set uintptr, signum int32) (_3 int32)
TEXT ·Ysigdelset(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW signum+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsigdelset(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ysigemptyset(t *TLS, set uintptr) (_2 int32)
TEXT ·Ysigemptyset(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsigemptyset(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ysigfillset(t *TLS, set uintptr) (_2 int32)
TEXT ·Ysigfillset(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsigfillset(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ysigismember(t *TLS, set uintptr, signum int32) (_3 int32)
TEXT ·Ysigismember(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD set+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW signum+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsigismember(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ysin(t *TLS, x float64) (_2 float64)
TEXT ·Ysin(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsin(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ysinf(t *TLS, x float32) (_2 float32)
TEXT ·Ysinf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsinf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ysinh(t *TLS, x float64) (_2 float64)
TEXT ·Ysinh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsinh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ysinhf(t *TLS, x float32) (_2 float32)
TEXT ·Ysinhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsinhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ysleep(t *TLS, seconds uint32) (_2 uint32)
TEXT ·Ysleep(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW seconds+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsleep(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ysnprintf(t *TLS, str uintptr, size types.Size_t, format, args uintptr) (r int32)
TEXT ·Ysnprintf(SB),$48-44
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
	CALL ·Xsnprintf(SB)
	MOVW 48(RSP), R0
	MOVW R0, r+40(FP)
	RET

// func Ysocket(t *TLS, domain, type1, protocol int32) (_2 int32)
TEXT ·Ysocket(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW domain+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW type1+12(FP), R0
	MOVW R0, 20(RSP)
	MOVW protocol+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xsocket(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ysprintf(t *TLS, str, format, args uintptr) (r int32)
TEXT ·Ysprintf(SB),$40-36
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
	CALL ·Xsprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ysqrt(t *TLS, x float64) (_2 float64)
TEXT ·Ysqrt(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsqrt(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ysqrtf(t *TLS, x float32) (_2 float32)
TEXT ·Ysqrtf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsqrtf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ysrand48(t *TLS, seedval long)
TEXT ·Ysrand48(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD seedval+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsrand48(SB)
	RET

// func Ysrandomdev(t *TLS)
TEXT ·Ysrandomdev(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xsrandomdev(SB)
	RET

// func Ysscanf(t *TLS, str, format, va uintptr) (_2 int32)
TEXT ·Ysscanf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD format+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsscanf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Ystat(t *TLS, pathname, statbuf uintptr) (_2 int32)
TEXT ·Ystat(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD statbuf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstat(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ystat64(t *TLS, pathname, statbuf uintptr) (_2 int32)
TEXT ·Ystat64(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD statbuf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstat64(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ystatfs(t *TLS, path uintptr, buf uintptr) (_3 int32)
TEXT ·Ystatfs(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstatfs(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ystrcasecmp(t *TLS, s1, s2 uintptr) (_2 int32)
TEXT ·Ystrcasecmp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s2+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcasecmp(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ystrcat(t *TLS, dest, src uintptr) (r uintptr)
TEXT ·Ystrcat(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcat(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrchr(t *TLS, s uintptr, c int32) (_3 uintptr)
TEXT ·Ystrchr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xstrchr(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
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

// func Ystrcmp(t *TLS, s1, s2 uintptr) (_2 int32)
TEXT ·Ystrcmp(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s2+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcmp(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ystrcoll(t *TLS, s1, s2 uintptr) (_2 int32)
TEXT ·Ystrcoll(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s1+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD s2+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcoll(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ystrcpy(t *TLS, dest, src uintptr) (r uintptr)
TEXT ·Ystrcpy(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD dest+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD src+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcpy(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrcspn(t *TLS, s, reject uintptr) (r types.Size_t)
TEXT ·Ystrcspn(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD reject+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrcspn(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrdup(tls *TLS, s uintptr) (_2 uintptr)
TEXT ·Ystrdup(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xstrdup(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ystrerror(t *TLS, errnum int32) (_2 uintptr)
TEXT ·Ystrerror(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW errnum+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xstrerror(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ystrerror_r(t *TLS, errnum int32, buf uintptr, buflen types.Size_t) (r int32)
TEXT ·Ystrerror_r(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW errnum+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD buflen+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xstrerror_r(SB)
	MOVW 40(RSP), R0
	MOVW R0, r+32(FP)
	RET

// func Ystrftime(tls *TLS, s uintptr, n size_t, f uintptr, tm uintptr) (r size_t)
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

// func Ystrlcat(tls *TLS, d uintptr, s uintptr, n size_t) (_4 size_t)
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
	MOVD R0, _4+32(FP)
	RET

// func Ystrlcpy(tls *TLS, d uintptr, s uintptr, n size_t) (_4 size_t)
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
	MOVD R0, _4+32(FP)
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

// func Ystrncasecmp(tls *TLS, _l uintptr, _r uintptr, n size_t) (_4 int32)
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
	MOVW R0, _4+32(FP)
	RET

// func Ystrncat(tls *TLS, d uintptr, s uintptr, n size_t) (_4 uintptr)
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
	MOVD R0, _4+32(FP)
	RET

// func Ystrncmp(t *TLS, s1, s2 uintptr, n types.Size_t) (_3 int32)
TEXT ·Ystrncmp(SB),$40-36
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
	CALL ·Xstrncmp(SB)
	MOVW 40(RSP), R0
	MOVW R0, _3+32(FP)
	RET

// func Ystrncpy(t *TLS, dest, src uintptr, n types.Size_t) (r uintptr)
TEXT ·Ystrncpy(SB),$40-40
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
	CALL ·Xstrncpy(SB)
	MOVD 40(RSP), R0
	MOVD R0, r+32(FP)
	RET

// func Ystrnlen(tls *TLS, s uintptr, n size_t) (_3 size_t)
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
	MOVD R0, _3+24(FP)
	RET

// func Ystrpbrk(t *TLS, s, accept uintptr) (_2 uintptr)
TEXT ·Ystrpbrk(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD accept+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrpbrk(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ystrrchr(t *TLS, s uintptr, c int32) (r uintptr)
TEXT ·Ystrrchr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW c+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xstrrchr(SB)
	MOVD 32(RSP), R0
	MOVD R0, r+24(FP)
	RET

// func Ystrspn(tls *TLS, s uintptr, c uintptr) (_3 size_t)
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
	MOVD R0, _3+24(FP)
	RET

// func Ystrstr(t *TLS, haystack, needle uintptr) (_2 uintptr)
TEXT ·Ystrstr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD haystack+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD needle+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xstrstr(SB)
	MOVD 32(RSP), R0
	MOVD R0, _2+24(FP)
	RET

// func Ystrtod(tls *TLS, s uintptr, p uintptr) (_3 float64)
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
	MOVD R0, _3+24(FP)
	RET

// func Ystrtof(tls *TLS, s uintptr, p uintptr) (_3 float32)
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
	MOVW R0, _3+24(FP)
	RET

// func Ystrtoimax(tls *TLS, s uintptr, p uintptr, base int32) (_4 intmax_t)
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
	MOVD R0, _4+32(FP)
	RET

// func Ystrtok(tls *TLS, s uintptr, sep uintptr) (_3 uintptr)
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
	MOVD R0, _3+24(FP)
	RET

// func Ystrtol(tls *TLS, s uintptr, p uintptr, base int32) (_4 int64)
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
	MOVD R0, _4+32(FP)
	RET

// func Ystrtold(tls *TLS, s uintptr, p uintptr) (_3 float64)
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
	MOVD R0, _3+24(FP)
	RET

// func Ystrtoll(tls *TLS, s uintptr, p uintptr, base int32) (_4 int64)
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
	MOVD R0, _4+32(FP)
	RET

// func Ystrtoul(tls *TLS, s uintptr, p uintptr, base int32) (_4 uint64)
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
	MOVD R0, _4+32(FP)
	RET

// func Ystrtoull(tls *TLS, s uintptr, p uintptr, base int32) (_4 uint64)
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
	MOVD R0, _4+32(FP)
	RET

// func Ystrtoumax(tls *TLS, s uintptr, p uintptr, base int32) (_4 uintmax_t)
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
	MOVD R0, _4+32(FP)
	RET

// func Ysymlink(t *TLS, target, linkpath uintptr) (_2 int32)
TEXT ·Ysymlink(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD target+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD linkpath+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xsymlink(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Ysymlinkat(t *TLS, target uintptr, newdirfd int32, linkpath uintptr) (_4 int32)
TEXT ·Ysymlinkat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD target+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW newdirfd+16(FP), R0
	MOVW R0, 24(RSP)
	MOVD linkpath+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xsymlinkat(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Ysysconf(t *TLS, name int32) (_2 long)
TEXT ·Ysysconf(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW name+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xsysconf(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ysysctlbyname(t *TLS, name, oldp, oldlenp, newp uintptr, newlen types.Size_t) (_3 int32)
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

// func Ysystem(t *TLS, command uintptr) (r int32)
TEXT ·Ysystem(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD command+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xsystem(SB)
	MOVW 24(RSP), R0
	MOVW R0, r+16(FP)
	RET

// func Ytan(t *TLS, x float64) (_2 float64)
TEXT ·Ytan(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtan(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ytanf(t *TLS, x float32) (_2 float32)
TEXT ·Ytanf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtanf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ytanh(t *TLS, x float64) (_2 float64)
TEXT ·Ytanh(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtanh(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ytanhf(t *TLS, x float32) (_2 float32)
TEXT ·Ytanhf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtanhf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ytcgetattr(t *TLS, fd int32, termios_p uintptr) (_3 int32)
TEXT ·Ytcgetattr(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD termios_p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtcgetattr(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ytcsendbreak(t *TLS, fd, duration int32) (_2 int32)
TEXT ·Ytcsendbreak(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW duration+12(FP), R0
	MOVW R0, 20(RSP)
	CALL ·Xtcsendbreak(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ytcsetattr(t *TLS, fd, optional_actions int32, termios_p uintptr) (_3 int32)
TEXT ·Ytcsetattr(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVW optional_actions+12(FP), R0
	MOVW R0, 20(RSP)
	MOVD termios_p+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtcsetattr(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ytgamma(t *TLS, x float64) (_2 float64)
TEXT ·Ytgamma(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtgamma(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ytime(t *TLS, tloc uintptr) (_2 types.Time_t)
TEXT ·Ytime(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD tloc+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtime(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ytmpfile(t *TLS) (_1 uintptr)
TEXT ·Ytmpfile(SB),$16-16
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xtmpfile(SB)
	MOVD 16(RSP), R0
	MOVD R0, _1+8(FP)
	RET

// func Ytmpnam(t *TLS, s uintptr) (r uintptr)
TEXT ·Ytmpnam(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtmpnam(SB)
	MOVD 24(RSP), R0
	MOVD R0, r+16(FP)
	RET

// func Ytoascii(tls *TLS, _c int32) (_2 int32)
TEXT ·Ytoascii(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW _c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtoascii(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ytolower(t *TLS, c int32) (_2 int32)
TEXT ·Ytolower(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtolower(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ytoupper(t *TLS, c int32) (_2 int32)
TEXT ·Ytoupper(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtoupper(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ytrunc(t *TLS, x float64) (_2 float64)
TEXT ·Ytrunc(SB),$24-24
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD x+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xtrunc(SB)
	MOVD 24(RSP), R0
	MOVD R0, _2+16(FP)
	RET

// func Ytruncate(t *TLS, path uintptr, length types.Off_t) (_3 int32)
TEXT ·Ytruncate(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD path+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD length+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xtruncate(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ytruncf(t *TLS, x float32) (_2 float32)
TEXT ·Ytruncf(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW x+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xtruncf(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ytzset(t *TLS)
TEXT ·Ytzset(SB),$8-8
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	CALL ·Xtzset(SB)
	RET

// func Yumask(t *TLS, mask types.Mode_t) (_2 types.Mode_t)
TEXT ·Yumask(SB),$24-18
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVH mask+8(FP), R0
	MOVH R0, 16(RSP)
	CALL ·Xumask(SB)
	MOVH 24(RSP), R0
	MOVH R0, _2+16(FP)
	RET

// func Yuname(t *TLS, buf uintptr) (_2 int32)
TEXT ·Yuname(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD buf+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xuname(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yungetc(t *TLS, c int32, stream uintptr) (_3 int32)
TEXT ·Yungetc(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD stream+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xungetc(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Yunlink(t *TLS, pathname uintptr) (_2 int32)
TEXT ·Yunlink(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD pathname+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xunlink(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yunlinkat(t *TLS, dirfd int32, pathname uintptr, flags int32) (_4 int32)
TEXT ·Yunlinkat(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW flags+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xunlinkat(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Yunlockpt(t *TLS, fd int32) (_2 int32)
TEXT ·Yunlockpt(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xunlockpt(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yunsetenv(tls *TLS, name uintptr) (_2 int32)
TEXT ·Yunsetenv(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD tls+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD name+8(FP), R0
	MOVD R0, 16(RSP)
	CALL ·Xunsetenv(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yusleep(t *TLS, usec types.Useconds_t) (_2 int32)
TEXT ·Yusleep(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW usec+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xusleep(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Yutime(t *TLS, filename, times uintptr) (_2 int32)
TEXT ·Yutime(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD times+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xutime(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Yutimensat(t *TLS, dirfd int32, pathname, times uintptr, flags int32) (_4 int32)
TEXT ·Yutimensat(SB),$48-44
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW dirfd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD pathname+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD times+24(FP), R0
	MOVD R0, 32(RSP)
	MOVW flags+32(FP), R0
	MOVW R0, 40(RSP)
	CALL ·Xutimensat(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Yutimes(t *TLS, filename, times uintptr) (_2 int32)
TEXT ·Yutimes(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD filename+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD times+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xutimes(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
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

// func Yvasprintf(t *TLS, strp, fmt, ap uintptr) (_2 int32)
TEXT ·Yvasprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD strp+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD fmt+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvasprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Yvfprintf(t *TLS, stream, format, ap uintptr) (_2 int32)
TEXT ·Yvfprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD stream+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD format+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD ap+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvfprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Yvprintf(t *TLS, s, ap uintptr) (_2 int32)
TEXT ·Yvprintf(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD ap+16(FP), R0
	MOVD R0, 24(RSP)
	CALL ·Xvprintf(SB)
	MOVW 32(RSP), R0
	MOVW R0, _2+24(FP)
	RET

// func Yvsnprintf(t *TLS, str uintptr, size types.Size_t, format, va uintptr) (_4 int32)
TEXT ·Yvsnprintf(SB),$48-44
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
	CALL ·Xvsnprintf(SB)
	MOVW 48(RSP), R0
	MOVW R0, _4+40(FP)
	RET

// func Yvsprintf(t *TLS, str, format, va uintptr) (_2 int32)
TEXT ·Yvsprintf(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD str+8(FP), R0
	MOVD R0, 16(RSP)
	MOVD format+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD va+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xvsprintf(SB)
	MOVW 40(RSP), R0
	MOVW R0, _2+32(FP)
	RET

// func Ywaitpid(t *TLS, pid types.Pid_t, wstatus uintptr, optname int32) (_4 types.Pid_t)
TEXT ·Ywaitpid(SB),$40-36
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW pid+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD wstatus+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW optname+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwaitpid(SB)
	MOVW 40(RSP), R0
	MOVW R0, _4+32(FP)
	RET

// func Ywcschr(t *TLS, wcs uintptr, wc wchar_t) (_3 uintptr)
TEXT ·Ywcschr(SB),$32-32
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD wcs+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW wc+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xwcschr(SB)
	MOVD 32(RSP), R0
	MOVD R0, _3+24(FP)
	RET

// func Ywctomb(t *TLS, s uintptr, wc wchar_t) (_3 int32)
TEXT ·Ywctomb(SB),$32-28
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVD s+8(FP), R0
	MOVD R0, 16(RSP)
	MOVW wc+16(FP), R0
	MOVW R0, 24(RSP)
	CALL ·Xwctomb(SB)
	MOVW 32(RSP), R0
	MOVW R0, _3+24(FP)
	RET

// func Ywcwidth(t *TLS, c wchar_t) (_2 int32)
TEXT ·Ywcwidth(SB),$24-20
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW c+8(FP), R0
	MOVW R0, 16(RSP)
	CALL ·Xwcwidth(SB)
	MOVW 24(RSP), R0
	MOVW R0, _2+16(FP)
	RET

// func Ywrite(t *TLS, fd int32, buf uintptr, count types.Size_t) (_4 types.Ssize_t)
TEXT ·Ywrite(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD buf+16(FP), R0
	MOVD R0, 24(RSP)
	MOVD count+24(FP), R0
	MOVD R0, 32(RSP)
	CALL ·Xwrite(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET

// func Ywritev(t *TLS, fd int32, iov uintptr, iovcnt int32) (_4 types.Ssize_t)
TEXT ·Ywritev(SB),$40-40
	GO_ARGS
	NO_LOCAL_POINTERS
	MOVD t+0(FP), R0
	MOVD R0, 8(RSP)
	MOVW fd+8(FP), R0
	MOVW R0, 16(RSP)
	MOVD iov+16(FP), R0
	MOVD R0, 24(RSP)
	MOVW iovcnt+24(FP), R0
	MOVW R0, 32(RSP)
	CALL ·Xwritev(SB)
	MOVD 40(RSP), R0
	MOVD R0, _4+32(FP)
	RET
