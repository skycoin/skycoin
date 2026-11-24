// Code generated for darwin/arm64 by 'qbecc --abi0wrap .', DO NOT EDIT.

package libc

import (
	"modernc.org/libc/langinfo"
	"modernc.org/libc/netinet/in"
	"modernc.org/libc/poll"
	"modernc.org/libc/pthread"
	"modernc.org/libc/sys/socket"
	"modernc.org/libc/sys/types"
	"modernc.org/libc/termios"
	"modernc.org/libc/time"
	"modernc.org/libc/wctype"
	"unsafe"
)

var _ unsafe.Pointer

// Not wrapped:
//
//	X__sync_add_and_fetch: generic function
//	X__sync_sub_and_fetch: generic function
//	Xatexit: Y is hand-written
//	Xcopyfile: unsupported type []interface{}
//	Xpthread_atfork: Y is hand-written
//	Xpthread_create: Y is hand-written
//	Xpthread_key_create: Y is hand-written
//	Xsigaction: Y is hand-written
//	Xsignal: Y is hand-written

//go:noescape
func Y_IO_putc(t *TLS, c int32, fp uintptr) (_3 int32)

//go:noescape
func Y_NSGetEnviron(t *TLS) (_1 uintptr)

//go:noescape
func Y___errno_location(t *TLS) (_1 uintptr)

//go:noescape
func Y__assert_fail(t *TLS, assertion, file uintptr, line uint32, function uintptr)

//go:noescape
func Y__assert_rtn(t *TLS, function, file uintptr, line int32, assertion uintptr)

//go:noescape
func Y__atomic_compare_exchangeInt16(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)

//go:noescape
func Y__atomic_compare_exchangeInt32(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)

//go:noescape
func Y__atomic_compare_exchangeInt64(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)

//go:noescape
func Y__atomic_compare_exchangeInt8(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)

//go:noescape
func Y__atomic_compare_exchangeUint16(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)

//go:noescape
func Y__atomic_compare_exchangeUint32(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)

//go:noescape
func Y__atomic_compare_exchangeUint64(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)

//go:noescape
func Y__atomic_compare_exchangeUint8(t *TLS, ptr, expected, desired uintptr, weak, success, failure int32) (_3 int32)

//go:noescape
func Y__atomic_exchangeInt16(t *TLS, ptr, val, ret uintptr, _ int32)

//go:noescape
func Y__atomic_exchangeInt32(t *TLS, ptr, val, ret uintptr, _ int32)

//go:noescape
func Y__atomic_exchangeInt64(t *TLS, ptr, val, ret uintptr, _ int32)

//go:noescape
func Y__atomic_exchangeInt8(t *TLS, ptr, val, ret uintptr, _ int32)

//go:noescape
func Y__atomic_exchangeUint16(t *TLS, ptr, val, ret uintptr, _ int32)

//go:noescape
func Y__atomic_exchangeUint32(t *TLS, ptr, val, ret uintptr, _ int32)

//go:noescape
func Y__atomic_exchangeUint64(t *TLS, ptr, val, ret uintptr, _ int32)

//go:noescape
func Y__atomic_exchangeUint8(t *TLS, ptr, val, ret uintptr, _ int32)

//go:noescape
func Y__atomic_fetch_addInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__atomic_fetch_addInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__atomic_fetch_addInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__atomic_fetch_addInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__atomic_fetch_addUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__atomic_fetch_addUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__atomic_fetch_addUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__atomic_fetch_addUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__atomic_fetch_andInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__atomic_fetch_andInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__atomic_fetch_andInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__atomic_fetch_andInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__atomic_fetch_andUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__atomic_fetch_andUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__atomic_fetch_andUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__atomic_fetch_andUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__atomic_fetch_orInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__atomic_fetch_orInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__atomic_fetch_orInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__atomic_fetch_orInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__atomic_fetch_orUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__atomic_fetch_orUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__atomic_fetch_orUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__atomic_fetch_orUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__atomic_fetch_subInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__atomic_fetch_subInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__atomic_fetch_subInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__atomic_fetch_subInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__atomic_fetch_subUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__atomic_fetch_subUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__atomic_fetch_subUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__atomic_fetch_subUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__atomic_fetch_xorInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__atomic_fetch_xorInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__atomic_fetch_xorInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__atomic_fetch_xorInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__atomic_fetch_xorUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__atomic_fetch_xorUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__atomic_fetch_xorUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__atomic_fetch_xorUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__atomic_loadFloat32(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadFloat64(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadInt16(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadInt32(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadInt64(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadInt8(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadUint16(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadUint32(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadUint64(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_loadUint8(t *TLS, ptr, ret uintptr, memorder int32)

//go:noescape
func Y__atomic_storeFloat32(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeFloat64(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeInt16(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeInt32(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeInt64(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeInt8(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeUint16(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeUint32(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeUint64(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__atomic_storeUint8(t *TLS, ptr, val uintptr, memorder int32)

//go:noescape
func Y__builtin___memcpy_chk(t *TLS, dest, src uintptr, n, os types.Size_t) (r uintptr)

//go:noescape
func Y__builtin___memmove_chk(t *TLS, dest, src uintptr, n, os types.Size_t) (_3 uintptr)

//go:noescape
func Y__builtin___memset_chk(t *TLS, s uintptr, c int32, n, os types.Size_t) (_4 uintptr)

//go:noescape
func Y__builtin___snprintf_chk(t *TLS, str uintptr, maxlen types.Size_t, flag int32, os types.Size_t, format, args uintptr) (r int32)

//go:noescape
func Y__builtin___sprintf_chk(t *TLS, s uintptr, flag int32, os types.Size_t, format, args uintptr) (r int32)

//go:noescape
func Y__builtin___strcat_chk(t *TLS, dest, src uintptr, os types.Size_t) (r uintptr)

//go:noescape
func Y__builtin___strcpy_chk(t *TLS, dest, src uintptr, os types.Size_t) (_3 uintptr)

//go:noescape
func Y__builtin___strncat_chk(t *TLS, dest, src uintptr, n, destlen Tsize_t) (r uintptr)

//go:noescape
func Y__builtin___strncpy_chk(t *TLS, dest, src uintptr, n, os types.Size_t) (r uintptr)

//go:noescape
func Y__builtin___vsnprintf_chk(t *TLS, str uintptr, maxlen types.Size_t, flag int32, os types.Size_t, format, args uintptr) (r int32)

//go:noescape
func Y__builtin___vsprintf_chk(t *TLS, s uintptr, flag int32, os types.Size_t, format, args uintptr) (r int32)

//go:noescape
func Y__builtin_abort(t *TLS)

//go:noescape
func Y__builtin_abs(t *TLS, j int32) (_2 int32)

//go:noescape
func Y__builtin_acosh(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_add_overflowInt64(t *TLS, a, b int64, res uintptr) (_3 int32)

//go:noescape
func Y__builtin_add_overflowUint32(t *TLS, a, b uint32, res uintptr) (_3 int32)

//go:noescape
func Y__builtin_add_overflowUint64(t *TLS, a, b uint64, res uintptr) (_3 int32)

//go:noescape
func Y__builtin_alloca(tls *TLS, size size_t) (_2 uintptr)

//go:noescape
func Y__builtin_asinh(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_asinhf(t *TLS, x float32) (_2 float32)

//go:noescape
func Y__builtin_atanh(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_bswap16(t *TLS, x uint16) (_2 uint16)

//go:noescape
func Y__builtin_bswap32(t *TLS, x uint32) (_2 uint32)

//go:noescape
func Y__builtin_bswap64(t *TLS, x uint64) (_2 uint64)

//go:noescape
func Y__builtin_bzero(t *TLS, s uintptr, n types.Size_t)

//go:noescape
func Y__builtin_cbrt(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_clz(t *TLS, n uint32) (_2 int32)

//go:noescape
func Y__builtin_clzl(t *TLS, n ulong) (_2 int32)

//go:noescape
func Y__builtin_clzll(t *TLS, n uint64) (_2 int32)

//go:noescape
func Y__builtin_copysign(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Y__builtin_copysignf(t *TLS, x, y float32) (_2 float32)

//go:noescape
func Y__builtin_copysignl(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Y__builtin_ctz(t *TLS, n uint32) (_2 int32)

//go:noescape
func Y__builtin_ctzll(tls *TLS, x uint64) (_2 int32)

//go:noescape
func Y__builtin_erf(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_erfc(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_exit(t *TLS, status int32)

//go:noescape
func Y__builtin_exp2(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_expect(t *TLS, exp, c long) (_2 long)

//go:noescape
func Y__builtin_expm1(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_fabs(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_fabsf(t *TLS, x float32) (_2 float32)

//go:noescape
func Y__builtin_fabsl(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_fdim(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Y__builtin_ffs(tls *TLS, i int32) (r int32)

//go:noescape
func Y__builtin_flt_rounds(t *TLS) (_1 int32)

//go:noescape
func Y__builtin_fma(t *TLS, x, y, z float64) (_2 float64)

//go:noescape
func Y__builtin_fmaf(t *TLS, x, y, z float32) (_2 float32)

//go:noescape
func Y__builtin_fmax(tls *TLS, x float64, y float64) (r float64)

//go:noescape
func Y__builtin_fmin(tls *TLS, x float64, y float64) (r float64)

//go:noescape
func Y__builtin_free(t *TLS, ptr uintptr)

//go:noescape
func Y__builtin_getentropy(t *TLS, buf uintptr, n types.Size_t) (_3 int32)

//go:noescape
func Y__builtin_huge_val(t *TLS) (_1 float64)

//go:noescape
func Y__builtin_huge_valf(t *TLS) (_1 float32)

//go:noescape
func Y__builtin_hypot(t *TLS, x float64, y float64) (r float64)

//go:noescape
func Y__builtin_ilogb(tls *TLS, x float64) (_2 int32)

//go:noescape
func Y__builtin_ilogbf(tls *TLS, x float32) (_2 int32)

//go:noescape
func Y__builtin_ilogbl(tls *TLS, x float64) (_2 int32)

//go:noescape
func Y__builtin_inf(t *TLS) (_1 float64)

//go:noescape
func Y__builtin_inff(t *TLS) (_1 float32)

//go:noescape
func Y__builtin_infl(t *TLS) (_1 float64)

//go:noescape
func Y__builtin_isnan(t *TLS, x float64) (_2 int32)

//go:noescape
func Y__builtin_isunordered(t *TLS, a, b float64) (_2 int32)

//go:noescape
func Y__builtin_lgamma(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_llabs(tls *TLS, a int64) (_2 int64)

//go:noescape
func Y__builtin_log1p(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_log1pf(t *TLS, x float32) (_2 float32)

//go:noescape
func Y__builtin_log2(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_logb(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_lrint(tls *TLS, x float64) (r long)

//go:noescape
func Y__builtin_lrintf(tls *TLS, x float32) (r long)

//go:noescape
func Y__builtin_lround(tls *TLS, x float64) (r long)

//go:noescape
func Y__builtin_malloc(t *TLS, size types.Size_t) (_2 uintptr)

//go:noescape
func Y__builtin_memcmp(t *TLS, s1, s2 uintptr, n types.Size_t) (_3 int32)

//go:noescape
func Y__builtin_memcpy(t *TLS, dest, src uintptr, n types.Size_t) (r uintptr)

//go:noescape
func Y__builtin_memset(t *TLS, s uintptr, c int32, n types.Size_t) (_4 uintptr)

//go:noescape
func Y__builtin_mmap(t *TLS, addr uintptr, length types.Size_t, prot, flags, fd int32, offset types.Off_t) (_5 uintptr)

//go:noescape
func Y__builtin_mul_overflowInt64(t *TLS, a, b int64, res uintptr) (_3 int32)

//go:noescape
func Y__builtin_mul_overflowUint128(t *TLS, a, b Uint128, res uintptr) (_3 int32)

//go:noescape
func Y__builtin_mul_overflowUint64(t *TLS, a, b uint64, res uintptr) (_3 int32)

//go:noescape
func Y__builtin_nan(t *TLS, s uintptr) (_2 float64)

//go:noescape
func Y__builtin_nanf(t *TLS, s uintptr) (_2 float32)

//go:noescape
func Y__builtin_nanl(t *TLS, s uintptr) (_2 float64)

//go:noescape
func Y__builtin_nextafter(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Y__builtin_nextafterf(t *TLS, x, y float32) (_2 float32)

//go:noescape
func Y__builtin_nexttoward(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Y__builtin_object_size(t *TLS, p uintptr, typ int32) (_3 types.Size_t)

//go:noescape
func Y__builtin_popcount(t *TLS, x uint32) (_2 int32)

//go:noescape
func Y__builtin_popcountl(t *TLS, x ulong) (_2 int32)

//go:noescape
func Y__builtin_popcountll(t *TLS, x uint64) (_2 int32)

//go:noescape
func Y__builtin_prefetch(t *TLS, addr, args uintptr)

//go:noescape
func Y__builtin_printf(t *TLS, s, args uintptr) (_2 int32)

//go:noescape
func Y__builtin_remainder(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Y__builtin_remquo(t *TLS, x, y float64, quo uintptr) (_3 float64)

//go:noescape
func Y__builtin_rint(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_rintf(tls *TLS, x float32) (r float32)

//go:noescape
func Y__builtin_round(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_roundf(t *TLS, x float32) (_2 float32)

//go:noescape
func Y__builtin_scalbn(t *TLS, x float64, n int32) (_3 float64)

//go:noescape
func Y__builtin_snprintf(t *TLS, str uintptr, size types.Size_t, format, args uintptr) (_4 int32)

//go:noescape
func Y__builtin_sprintf(t *TLS, str, format, args uintptr) (r int32)

//go:noescape
func Y__builtin_strchr(t *TLS, s uintptr, c int32) (_3 uintptr)

//go:noescape
func Y__builtin_strcmp(t *TLS, s1, s2 uintptr) (_2 int32)

//go:noescape
func Y__builtin_strcpy(t *TLS, dest, src uintptr) (_2 uintptr)

//go:noescape
func Y__builtin_strlen(t *TLS, s uintptr) (_2 types.Size_t)

//go:noescape
func Y__builtin_sub_overflowInt64(t *TLS, a, b int64, res uintptr) (_3 int32)

//go:noescape
func Y__builtin_tgamma(t *TLS, x float64) (_2 float64)

//go:noescape
func Y__builtin_trap(t *TLS)

//go:noescape
func Y__builtin_trunc(tls *TLS, x float64) (r float64)

//go:noescape
func Y__builtin_unreachable(t *TLS)

//go:noescape
func Y__builtin_vsnprintf(t *TLS, str uintptr, size types.Size_t, format, va uintptr) (_4 int32)

//go:noescape
func Y__c11_atomic_compare_exchange_strongInt16(t *TLS, ptr, expected uintptr, desired int16, success, failure int32) (_4 int32)

//go:noescape
func Y__c11_atomic_compare_exchange_strongInt32(t *TLS, ptr, expected uintptr, desired, success, failure int32) (_3 int32)

//go:noescape
func Y__c11_atomic_compare_exchange_strongInt64(t *TLS, ptr, expected uintptr, desired int64, success, failure int32) (_4 int32)

//go:noescape
func Y__c11_atomic_compare_exchange_strongInt8(t *TLS, ptr, expected uintptr, desired int8, success, failure int32) (_4 int32)

//go:noescape
func Y__c11_atomic_compare_exchange_strongUint16(t *TLS, ptr, expected uintptr, desired uint16, success, failure int32) (_4 int32)

//go:noescape
func Y__c11_atomic_compare_exchange_strongUint32(t *TLS, ptr, expected uintptr, desired uint32, success, failure int32) (_4 int32)

//go:noescape
func Y__c11_atomic_compare_exchange_strongUint64(t *TLS, ptr, expected uintptr, desired uint64, success, failure int32) (_4 int32)

//go:noescape
func Y__c11_atomic_compare_exchange_strongUint8(t *TLS, ptr, expected uintptr, desired uint8, success, failure int32) (_4 int32)

//go:noescape
func Y__c11_atomic_exchangeInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__c11_atomic_exchangeInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__c11_atomic_exchangeInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__c11_atomic_exchangeInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__c11_atomic_exchangeUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__c11_atomic_exchangeUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__c11_atomic_exchangeUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__c11_atomic_exchangeUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__c11_atomic_fetch_addInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__c11_atomic_fetch_addInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__c11_atomic_fetch_addInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__c11_atomic_fetch_addInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__c11_atomic_fetch_addUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__c11_atomic_fetch_addUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__c11_atomic_fetch_addUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__c11_atomic_fetch_addUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__c11_atomic_fetch_andInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__c11_atomic_fetch_andInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__c11_atomic_fetch_andInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__c11_atomic_fetch_andInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__c11_atomic_fetch_andUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__c11_atomic_fetch_andUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__c11_atomic_fetch_andUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__c11_atomic_fetch_andUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__c11_atomic_fetch_orInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__c11_atomic_fetch_orInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__c11_atomic_fetch_orInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__c11_atomic_fetch_orInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__c11_atomic_fetch_orUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__c11_atomic_fetch_orUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__c11_atomic_fetch_orUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__c11_atomic_fetch_orUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__c11_atomic_fetch_subInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__c11_atomic_fetch_subInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__c11_atomic_fetch_subInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__c11_atomic_fetch_subInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__c11_atomic_fetch_subUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__c11_atomic_fetch_subUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__c11_atomic_fetch_subUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__c11_atomic_fetch_subUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__c11_atomic_fetch_xorInt16(t *TLS, ptr uintptr, val int16, _ int32) (r int16)

//go:noescape
func Y__c11_atomic_fetch_xorInt32(t *TLS, ptr uintptr, val int32, _ int32) (r int32)

//go:noescape
func Y__c11_atomic_fetch_xorInt64(t *TLS, ptr uintptr, val int64, _ int32) (r int64)

//go:noescape
func Y__c11_atomic_fetch_xorInt8(t *TLS, ptr uintptr, val int8, _ int32) (r int8)

//go:noescape
func Y__c11_atomic_fetch_xorUint16(t *TLS, ptr uintptr, val uint16, _ int32) (r uint16)

//go:noescape
func Y__c11_atomic_fetch_xorUint32(t *TLS, ptr uintptr, val uint32, _ int32) (r uint32)

//go:noescape
func Y__c11_atomic_fetch_xorUint64(t *TLS, ptr uintptr, val uint64, _ int32) (r uint64)

//go:noescape
func Y__c11_atomic_fetch_xorUint8(t *TLS, ptr uintptr, val uint8, _ int32) (r uint8)

//go:noescape
func Y__c11_atomic_loadInt16(t *TLS, ptr uintptr, memorder int32) (r int16)

//go:noescape
func Y__c11_atomic_loadInt32(t *TLS, ptr uintptr, memorder int32) (r int32)

//go:noescape
func Y__c11_atomic_loadInt64(t *TLS, ptr uintptr, memorder int32) (r int64)

//go:noescape
func Y__c11_atomic_loadInt8(t *TLS, ptr uintptr, memorder int32) (r int8)

//go:noescape
func Y__c11_atomic_loadUint16(t *TLS, ptr uintptr, memorder int32) (r uint16)

//go:noescape
func Y__c11_atomic_loadUint32(t *TLS, ptr uintptr, memorder int32) (r uint32)

//go:noescape
func Y__c11_atomic_loadUint64(t *TLS, ptr uintptr, memorder int32) (r uint64)

//go:noescape
func Y__c11_atomic_loadUint8(t *TLS, ptr uintptr, memorder int32) (r uint8)

//go:noescape
func Y__c11_atomic_storeInt16(t *TLS, ptr uintptr, val int16, memorder int32)

//go:noescape
func Y__c11_atomic_storeInt32(t *TLS, ptr uintptr, val int32, memorder int32)

//go:noescape
func Y__c11_atomic_storeInt64(t *TLS, ptr uintptr, val int64, memorder int32)

//go:noescape
func Y__c11_atomic_storeInt8(t *TLS, ptr uintptr, val int8, memorder int32)

//go:noescape
func Y__c11_atomic_storeUint16(t *TLS, ptr uintptr, val uint16, memorder int32)

//go:noescape
func Y__c11_atomic_storeUint32(t *TLS, ptr uintptr, val uint32, memorder int32)

//go:noescape
func Y__c11_atomic_storeUint64(t *TLS, ptr uintptr, val uint64, memorder int32)

//go:noescape
func Y__c11_atomic_storeUint8(t *TLS, ptr uintptr, val uint8, memorder int32)

//go:noescape
func Y__ccgo_dmesg(t *TLS, fmt uintptr, va uintptr)

//go:noescape
func Y__ccgo_getMutexType(tls *TLS, m uintptr) (_2 int32)

//go:noescape
func Y__ccgo_in6addr_anyp(t *TLS) (_1 uintptr)

//go:noescape
func Y__ccgo_pthreadAttrGetDetachState(tls *TLS, a uintptr) (_2 int32)

//go:noescape
func Y__ccgo_pthreadMutexattrGettype(tls *TLS, a uintptr) (_2 int32)

//go:noescape
func Y__ccgo_sqlite3_log(t *TLS, iErrCode int32, zFormat uintptr, args uintptr)

//go:noescape
func Y__cmsg_nxthdr(t *TLS, msgh, cmsg uintptr) (_2 uintptr)

//go:noescape
func Y__ctype_get_mb_cur_max(t *TLS) (_1 types.Size_t)

//go:noescape
func Y__darwin_check_fd_set_overflow(tls *TLS, _ int32, _ uintptr, _ int32) (_4 int32)

//go:noescape
func Y__darwin_fd_clr(tls *TLS, _fd int32, _p uintptr) (_3 int32)

//go:noescape
func Y__darwin_fd_isset(tls *TLS, _fd int32, _p uintptr) (_3 int32)

//go:noescape
func Y__darwin_fd_set(tls *TLS, _fd int32, _p uintptr) (_3 int32)

//go:noescape
func Y__env_rm_add(tls *TLS, old uintptr, new uintptr)

//go:noescape
func Y__errno_location(t *TLS) (_1 uintptr)

//go:noescape
func Y__error(t *TLS) (_1 uintptr)

//go:noescape
func Y__floatscan(tls *TLS, f uintptr, prec int32, pok int32) (_4 float64)

//go:noescape
func Y__flt_rounds(t *TLS) (_1 int32)

//go:noescape
func Y__fpclassify(tls *TLS, x float64) (_2 int32)

//go:noescape
func Y__fpclassifyd(tls *TLS, x float64) (r int32)

//go:noescape
func Y__fpclassifyf(tls *TLS, x float32) (_2 int32)

//go:noescape
func Y__fpclassifyl(tls *TLS, x float64) (_2 int32)

//go:noescape
func Y__h_errno_location(tls *TLS) (_1 uintptr)

//go:noescape
func Y__inet_aton(tls *TLS, s0 uintptr, dest uintptr) (_3 int32)

//go:noescape
func Y__inline_isnand(t *TLS, x float64) (_2 int32)

//go:noescape
func Y__inline_isnanf(t *TLS, x float32) (_2 int32)

//go:noescape
func Y__inline_isnanl(t *TLS, x float64) (_2 int32)

//go:noescape
func Y__intscan(tls *TLS, f uintptr, base uint32, pok int32, lim uint64) (_5 uint64)

//go:noescape
func Y__isctype(tls *TLS, _c int32, _f uint64) (_3 int32)

//go:noescape
func Y__isfinite(tls *TLS, d float64) (_2 int32)

//go:noescape
func Y__isfinitef(tls *TLS, f float32) (_2 int32)

//go:noescape
func Y__isfinitel(tls *TLS, d float64) (_2 int32)

//go:noescape
func Y__isnan(t *TLS, arg float64) (_2 int32)

//go:noescape
func Y__isnanf(t *TLS, arg float32) (_2 int32)

//go:noescape
func Y__isnanl(t *TLS, arg float64) (_2 int32)

//go:noescape
func Y__isoc99_sscanf(t *TLS, str, format, va uintptr) (_2 int32)

//go:noescape
func Y__istype(tls *TLS, _c int32, _f uint64) (_3 int32)

//go:noescape
func Y__lookup_ipliteral(tls *TLS, buf uintptr, name uintptr, family int32) (_4 int32)

//go:noescape
func Y__lookup_name(tls *TLS, buf uintptr, canon uintptr, name uintptr, family int32, flags int32) (_6 int32)

//go:noescape
func Y__lookup_serv(tls *TLS, buf uintptr, name uintptr, proto int32, socktype int32, flags int32) (_6 int32)

//go:noescape
func Y__maskrune(tls *TLS, _c int32, _f uint64) (_3 int32)

//go:noescape
func Y__putenv(tls *TLS, s uintptr, l size_t, r uintptr) (_4 int32)

//go:noescape
func Y__shgetc(tls *TLS, f uintptr) (_2 int32)

//go:noescape
func Y__shlim(tls *TLS, f uintptr, lim off_t)

//go:noescape
func Y__sincos_stret(t *TLS, f float64) (_2 struct {
	F__sinval float64
	F__cosval float64
})

//go:noescape
func Y__sincosf_stret(t *TLS, f float32) (_2 struct {
	F__sinval float32
	F__cosval float32
})

//go:noescape
func Y__sincospi_stret(t *TLS, f float64) (_2 struct {
	F__sinval float64
	F__cosval float64
})

//go:noescape
func Y__sincospif_stret(t *TLS, f float32) (_2 struct {
	F__sinval float32
	F__cosval float32
})

//go:noescape
func Y__srget(t *TLS, f uintptr) (_2 int32)

//go:noescape
func Y__strchrnul(tls *TLS, s uintptr, c int32) (_3 uintptr)

//go:noescape
func Y__strncasecmp_l(tls *TLS, l uintptr, r uintptr, n size_t, loc locale_t) (_5 int32)

//go:noescape
func Y__svfscanf(t *TLS, f uintptr, p, q uintptr) (_3 int32)

//go:noescape
func Y__swbuf(t *TLS, i int32, f uintptr) (_3 int32)

//go:noescape
func Y__sync_add_and_fetch_uint32(t *TLS, p uintptr, v uint32) (_3 uint32)

//go:noescape
func Y__sync_sub_and_fetch_uint32(t *TLS, p uintptr, v uint32) (_3 uint32)

//go:noescape
func Y__sync_synchronize(t *TLS)

//go:noescape
func Y__sync_val_compare_and_swapInt16(t *TLS, ptr uintptr, oldval, newval int16) (r int16)

//go:noescape
func Y__sync_val_compare_and_swapInt32(t *TLS, ptr uintptr, oldval, newval int32) (r int32)

//go:noescape
func Y__sync_val_compare_and_swapInt64(t *TLS, ptr uintptr, oldval, newval int64) (r int64)

//go:noescape
func Y__sync_val_compare_and_swapInt8(t *TLS, ptr uintptr, oldval, newval int8) (r int8)

//go:noescape
func Y__sync_val_compare_and_swapUint16(t *TLS, ptr uintptr, oldval, newval uint16) (r uint16)

//go:noescape
func Y__sync_val_compare_and_swapUint32(t *TLS, ptr uintptr, oldval, newval uint32) (r uint32)

//go:noescape
func Y__sync_val_compare_and_swapUint64(t *TLS, ptr uintptr, oldval, newval uint64) (r uint64)

//go:noescape
func Y__sync_val_compare_and_swapUint8(t *TLS, ptr uintptr, oldval, newval uint8) (r uint8)

//go:noescape
func Y__tolower(tls *TLS, c int32) (_2 int32)

//go:noescape
func Y__toread(tls *TLS, f uintptr) (_2 int32)

//go:noescape
func Y__toread_needs_stdio_exit(tls *TLS)

//go:noescape
func Y__toupper(tls *TLS, c int32) (_2 int32)

//go:noescape
func Y__uflow(tls *TLS, f uintptr) (_2 int32)

//go:noescape
func Y__wcwidth(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Y_exit(_ *TLS, status int32)

//go:noescape
func Y_longjmp(t *TLS, env uintptr, val int32)

//go:noescape
func Y_obstack_begin(t *TLS, obstack uintptr, size, alignment int32, chunkfun, freefun uintptr) (_4 int32)

//go:noescape
func Y_obstack_newchunk(t *TLS, obstack uintptr, length int32) (_3 int32)

//go:noescape
func Y_setjmp(t *TLS, env uintptr) (_2 int32)

//go:noescape
func Yabort(t *TLS)

//go:noescape
func Yabs(t *TLS, j int32) (_2 int32)

//go:noescape
func Yaccept(t *TLS, sockfd int32, addr uintptr, addrlen uintptr) (_4 int32)

//go:noescape
func Yaccess(t *TLS, pathname uintptr, mode int32) (_3 int32)

//go:noescape
func Yacos(t *TLS, x float64) (_2 float64)

//go:noescape
func Yacosf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yacosh(t *TLS, x float64) (_2 float64)

//go:noescape
func Yacoshf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yalarm(t *TLS, seconds uint32) (_2 uint32)

//go:noescape
func Yalloca(tls *TLS, size size_t) (_2 uintptr)

//go:noescape
func Yarc4random_buf(t *TLS, buf uintptr, buflen size_t)

//go:noescape
func Yasin(t *TLS, x float64) (_2 float64)

//go:noescape
func Yasinf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yasinh(t *TLS, x float64) (_2 float64)

//go:noescape
func Yasinhf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yatan(t *TLS, x float64) (_2 float64)

//go:noescape
func Yatan2(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Yatan2f(t *TLS, x, y float32) (_2 float32)

//go:noescape
func Yatanf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yatanh(t *TLS, x float64) (_2 float64)

//go:noescape
func Yatanhf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yatof(t *TLS, nptr uintptr) (_2 float64)

//go:noescape
func Yatoi(t *TLS, nptr uintptr) (_2 int32)

//go:noescape
func Yatol(t *TLS, nptr uintptr) (_2 long)

//go:noescape
func Ybind(t *TLS, sockfd int32, addr uintptr, addrlen uint32) (_4 int32)

//go:noescape
func Ybsearch(tls *TLS, key uintptr, base uintptr, nel size_t, width size_t, cmp uintptr) (_6 uintptr)

//go:noescape
func __ccgo_abi0_bsearch_4(_0 *TLS, _1 uintptr, _2 uintptr, __ccgo_fp uintptr) (_3 int32)

func __ccgo_abiInternal_bsearch_4(tls *TLS, dest, abi0CodePtr uintptr) {
	f := func(_0 *TLS, _1 uintptr, _2 uintptr) (_3 int32) {
		return __ccgo_abi0_bsearch_4(_0, _1, _2, abi0CodePtr)
	}
	*(*[2]uintptr)(unsafe.Pointer(dest)) = *(*[2]uintptr)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&f))))
}

//go:noescape
func Ybzero(t *TLS, s uintptr, n types.Size_t)

//go:noescape
func Ycalloc(t *TLS, n, size types.Size_t) (_2 uintptr)

//go:noescape
func Ycbrt(tls *TLS, x float64) (r1 float64)

//go:noescape
func Ycbrtf(tls *TLS, x float32) (r1 float32)

//go:noescape
func Ycbrtl(tls *TLS, x float64) (r float64)

//go:noescape
func Yceil(t *TLS, x float64) (_2 float64)

//go:noescape
func Yceilf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ycfgetospeed(t *TLS, termios_p uintptr) (_2 termios.Speed_t)

//go:noescape
func Ychdir(t *TLS, path uintptr) (_2 int32)

//go:noescape
func Ychflags(t *TLS, path uintptr, flags uint32) (_3 int32)

//go:noescape
func Ychmod(t *TLS, pathname uintptr, mode types.Mode_t) (_3 int32)

//go:noescape
func Ychown(t *TLS, pathname uintptr, owner types.Uid_t, group types.Gid_t) (_4 int32)

//go:noescape
func Yclearerr(tls *TLS, f uintptr)

//go:noescape
func Yclock(t *TLS) (r time.Clock_t)

//go:noescape
func Yclock_gettime(t *TLS, clk_id int32, tp uintptr) (_3 int32)

//go:noescape
func Yclose(t *TLS, fd int32) (_2 int32)

//go:noescape
func Yclosedir(t *TLS, dir uintptr) (_2 int32)

//go:noescape
func Yconfstr(t *TLS, name int32, buf uintptr, size types.Size_t) (_4 types.Size_t)

//go:noescape
func Yconnect(t *TLS, sockfd int32, addr uintptr, addrlen uint32) (_4 int32)

//go:noescape
func Ycopysign(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Ycopysignf(t *TLS, x, y float32) (_2 float32)

//go:noescape
func Ycopysignl(tls *TLS, x float64, y float64) (_3 float64)

//go:noescape
func Ycos(t *TLS, x float64) (_2 float64)

//go:noescape
func Ycosf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ycosh(t *TLS, x float64) (_2 float64)

//go:noescape
func Ycoshf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yctime(t *TLS, timep uintptr) (_2 uintptr)

//go:noescape
func Yctime_r(t *TLS, timep, buf uintptr) (_2 uintptr)

//go:noescape
func Ydifftime(t *TLS, time1, time0 types.Time_t) (_2 float64)

//go:noescape
func Ydigittoint(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Ydiv(tls *TLS, num, den int32) (r Tdiv_t)

//go:noescape
func Ydlclose(t *TLS, handle uintptr) (_2 int32)

//go:noescape
func Ydlerror(t *TLS) (_1 uintptr)

//go:noescape
func Ydlopen(t *TLS, filename uintptr, flags int32) (_3 uintptr)

//go:noescape
func Ydlsym(t *TLS, handle, symbol uintptr) (_2 uintptr)

//go:noescape
func Ydup(tls *TLS, fd int32) (r int32)

//go:noescape
func Ydup2(t *TLS, oldfd, newfd int32) (_2 int32)

//go:noescape
func Yendpwent(t *TLS)

//go:noescape
func Yerf(t *TLS, x float64) (_2 float64)

//go:noescape
func Yerfc(t *TLS, x float64) (_2 float64)

//go:noescape
func Yexecve(tls *TLS, path uintptr, argv uintptr, envp uintptr) (r int32)

//go:noescape
func Yexit(t *TLS, status int32)

//go:noescape
func Yexp(t *TLS, x float64) (_2 float64)

//go:noescape
func Yexp2(t *TLS, x float64) (_2 float64)

//go:noescape
func Yexpf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yexpm1(tls *TLS, x3 float64) (r float64)

//go:noescape
func Yexpm1f(tls *TLS, x3 float32) (r float32)

//go:noescape
func Yexpm1l(tls *TLS, x float64) (r float64)

//go:noescape
func Yfabs(t *TLS, x float64) (_2 float64)

//go:noescape
func Yfabsf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yfabsl(tls *TLS, x float64) (_2 float64)

//go:noescape
func Yfaccessat(t *TLS, dirfd int32, pathname uintptr, mode int32, flags int32) (_5 int32)

//go:noescape
func Yfchmod(t *TLS, fd int32, mode types.Mode_t) (_3 int32)

//go:noescape
func Yfchmodat(t *TLS, dirfd int32, pathname uintptr, mode types.Mode_t, flags int32) (_5 int32)

//go:noescape
func Yfchown(t *TLS, fd int32, owner types.Uid_t, group types.Gid_t) (_4 int32)

//go:noescape
func Yfchownat(t *TLS, dirfd int32, pathname uintptr, owner types.Uid_t, group types.Gid_t, flags int32) (_6 int32)

//go:noescape
func Yfclose(t *TLS, stream uintptr) (_2 int32)

//go:noescape
func Yfcntl(t *TLS, fd, cmd int32, args uintptr) (_3 int32)

//go:noescape
func Yfcntl64(t *TLS, fd, cmd int32, args uintptr) (r int32)

//go:noescape
func Yfdim(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Yfdopen(t *TLS, fd int32, mode uintptr) (_3 uintptr)

//go:noescape
func Yfdopendir(t *TLS, fd int32) (_2 uintptr)

//go:noescape
func Yfeof(t *TLS, f uintptr) (r int32)

//go:noescape
func Yferror(t *TLS, stream uintptr) (_2 int32)

//go:noescape
func Yfflush(t *TLS, stream uintptr) (_2 int32)

//go:noescape
func Yffs(tls *TLS, i int32) (r int32)

//go:noescape
func Yfgetc(t *TLS, stream uintptr) (_2 int32)

//go:noescape
func Yfgets(t *TLS, s uintptr, size int32, stream uintptr) (_4 uintptr)

//go:noescape
func Yfileno(t *TLS, stream uintptr) (_2 int32)

//go:noescape
func Yflock(t *TLS, fd, operation int32) (_2 int32)

//go:noescape
func Yflockfile(t *TLS, filehandle uintptr)

//go:noescape
func Yfloor(t *TLS, x float64) (_2 float64)

//go:noescape
func Yfloorf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yfma(t *TLS, x, y, z float64) (_2 float64)

//go:noescape
func Yfmaf(t *TLS, x, y, z float32) (r float32)

//go:noescape
func Yfmax(tls *TLS, x float64, y float64) (r float64)

//go:noescape
func Yfmaxf(tls *TLS, x float32, y float32) (r float32)

//go:noescape
func Yfmaxl(tls *TLS, x float64, y float64) (r float64)

//go:noescape
func Yfmin(tls *TLS, x float64, y float64) (r float64)

//go:noescape
func Yfminf(tls *TLS, x float32, y float32) (r float32)

//go:noescape
func Yfminl(tls *TLS, x float64, y float64) (r float64)

//go:noescape
func Yfmod(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Yfmodf(t *TLS, x, y float32) (_2 float32)

//go:noescape
func Yfmodl(tls *TLS, x float64, y float64) (_3 float64)

//go:noescape
func Yfopen(t *TLS, pathname, mode uintptr) (_2 uintptr)

//go:noescape
func Yfopen64(t *TLS, pathname, mode uintptr) (_2 uintptr)

//go:noescape
func Yfork(t *TLS) (_1 int32)

//go:noescape
func Yfpathconf(t *TLS, fd int32, name int32) (_3 long)

//go:noescape
func Yfprintf(t *TLS, stream, format, args uintptr) (_2 int32)

//go:noescape
func Yfputc(t *TLS, c int32, stream uintptr) (_3 int32)

//go:noescape
func Yfputs(t *TLS, s, stream uintptr) (_2 int32)

//go:noescape
func Yfread(t *TLS, ptr uintptr, size, nmemb types.Size_t, stream uintptr) (_4 types.Size_t)

//go:noescape
func Yfree(t *TLS, p uintptr)

//go:noescape
func Yfreeaddrinfo(tls *TLS, p uintptr)

//go:noescape
func Yfreopen(t *TLS, pathname, mode, stream uintptr) (r uintptr)

//go:noescape
func Yfrexp(t *TLS, x float64, exp uintptr) (_3 float64)

//go:noescape
func Yfrexpf(t *TLS, x float32, exp uintptr) (_3 float32)

//go:noescape
func Yfscanf(t *TLS, stream, format, va uintptr) (_2 int32)

//go:noescape
func Yfsctl(t *TLS, path uintptr, request ulong, data uintptr, options uint32) (_5 int32)

//go:noescape
func Yfseek(t *TLS, stream uintptr, offset long, whence int32) (_4 int32)

//go:noescape
func Yfseeko(t *TLS, stream uintptr, offset types.Off_t, whence int32) (_4 int32)

//go:noescape
func Yfstat(t *TLS, fd int32, statbuf uintptr) (_3 int32)

//go:noescape
func Yfstat64(t *TLS, fd int32, statbuf uintptr) (_3 int32)

//go:noescape
func Yfstatat(t *TLS, dirfd int32, pathname, statbuf uintptr, flags int32) (_4 int32)

//go:noescape
func Yfstatfs(t *TLS, fd int32, buf uintptr) (_3 int32)

//go:noescape
func Yfsync(t *TLS, fd int32) (_2 int32)

//go:noescape
func Yftell(t *TLS, stream uintptr) (_2 long)

//go:noescape
func Yftello(t *TLS, stream uintptr) (_2 types.Off_t)

//go:noescape
func Yftruncate(t *TLS, fd int32, length types.Off_t) (_3 int32)

//go:noescape
func Yftrylockfile(t *TLS, filehandle uintptr) (_2 int32)

//go:noescape
func Yfts_close(t *TLS, ftsp uintptr) (_2 int32)

//go:noescape
func Yfts_open(t *TLS, path_argv uintptr, options int32, compar uintptr) (_4 uintptr)

//go:noescape
func Yfts_read(t *TLS, ftsp uintptr) (_2 uintptr)

//go:noescape
func Yftw(tls *TLS, path uintptr, fn uintptr, fd_limit int32) (r int32)

//go:noescape
func __ccgo_abi0_ftw_1(_0 *TLS, _1 uintptr, _2 uintptr, _3 int32, __ccgo_fp uintptr) (_4 int32)

func __ccgo_abiInternal_ftw_1(tls *TLS, dest, abi0CodePtr uintptr) {
	f := func(_0 *TLS, _1 uintptr, _2 uintptr, _3 int32) (_4 int32) {
		return __ccgo_abi0_ftw_1(_0, _1, _2, _3, abi0CodePtr)
	}
	*(*[2]uintptr)(unsafe.Pointer(dest)) = *(*[2]uintptr)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&f))))
}

//go:noescape
func Yfunlockfile(t *TLS, filehandle uintptr)

//go:noescape
func Yfutimens(t *TLS, fd int32, times uintptr) (_3 int32)

//go:noescape
func Yfutimes(t *TLS, fd int32, tv uintptr) (_3 int32)

//go:noescape
func Yfwrite(t *TLS, ptr uintptr, size, nmemb types.Size_t, stream uintptr) (_4 types.Size_t)

//go:noescape
func Ygai_strerror(t *TLS, errcode int32) (_2 uintptr)

//go:noescape
func Ygetaddrinfo(tls *TLS, host uintptr, serv uintptr, hint uintptr, res uintptr) (_5 int32)

//go:noescape
func Ygetattrlist(t *TLS, path, attrList, attrBuf uintptr, attrBufSize types.Size_t, options uint32) (_4 int32)

//go:noescape
func Ygetc(t *TLS, stream uintptr) (_2 int32)

//go:noescape
func Ygetc_unlocked(t *TLS, stream uintptr) (_2 int32)

//go:noescape
func Ygetcwd(t *TLS, buf uintptr, size types.Size_t) (_3 uintptr)

//go:noescape
func Ygetegid(t *TLS) (_1 types.Gid_t)

//go:noescape
func Ygetentropy(t *TLS, buffer uintptr, length size_t) (_3 int32)

//go:noescape
func Ygetenv(t *TLS, name uintptr) (_2 uintptr)

//go:noescape
func Ygeteuid(t *TLS) (_1 types.Uid_t)

//go:noescape
func Ygetgid(t *TLS) (_1 types.Gid_t)

//go:noescape
func Ygetgrgid(t *TLS, gid uint32) (_2 uintptr)

//go:noescape
func Ygetgrgid_r(t *TLS, gid uint32, pGrp, buf uintptr, buflen types.Size_t, result uintptr) (_5 int32)

//go:noescape
func Ygetgrnam(t *TLS, name uintptr) (_2 uintptr)

//go:noescape
func Ygetgrnam_r(t *TLS, name, pGrp, buf uintptr, buflen types.Size_t, result uintptr) (_4 int32)

//go:noescape
func Ygethostbyaddr(tls *TLS, a uintptr, l socklen_t, af int32) (_4 uintptr)

//go:noescape
func Ygethostbyaddr_r(tls *TLS, a uintptr, l socklen_t, af int32, h uintptr, buf uintptr, buflen size_t, res uintptr, err uintptr) (_9 int32)

//go:noescape
func Ygethostbyname(tls *TLS, name uintptr) (_2 uintptr)

//go:noescape
func Ygethostbyname2(tls *TLS, name uintptr, af int32) (_3 uintptr)

//go:noescape
func Ygethostbyname2_r(tls *TLS, name uintptr, af int32, h uintptr, buf uintptr, buflen size_t, res uintptr, err uintptr) (_8 int32)

//go:noescape
func Ygethostname(t *TLS, name uintptr, slen types.Size_t) (_3 int32)

//go:noescape
func Ygethostuuid(t *TLS, id uintptr, wait uintptr) (_3 int32)

//go:noescape
func Ygetloadavg(t *TLS, loadavg uintptr, nelem int32) (_3 int32)

//go:noescape
func Ygetlogin_r(t *TLS, buf uintptr, bufsize types.Size_t) (_3 int32)

//go:noescape
func Ygetnameinfo(tls *TLS, sa1 uintptr, sl socklen_t, node uintptr, nodelen socklen_t, serv uintptr, servlen socklen_t, flags int32) (_8 int32)

//go:noescape
func Ygetpeername(t *TLS, sockfd int32, addr uintptr, addrlen uintptr) (_4 int32)

//go:noescape
func Ygetpid(t *TLS) (_1 int32)

//go:noescape
func Ygetprogname(t *TLS) (_1 uintptr)

//go:noescape
func Ygetpwnam(t *TLS, name uintptr) (_2 uintptr)

//go:noescape
func Ygetpwnam_r(t *TLS, name, cpwd, buf uintptr, buflen types.Size_t, result uintptr) (_4 int32)

//go:noescape
func Ygetpwuid(t *TLS, uid uint32) (_2 uintptr)

//go:noescape
func Ygetpwuid_r(t *TLS, uid types.Uid_t, cpwd, buf uintptr, buflen types.Size_t, result uintptr) (_5 int32)

//go:noescape
func Ygetresgid(t *TLS, rgid, egid, sgid uintptr) (_2 int32)

//go:noescape
func Ygetresuid(t *TLS, ruid, euid, suid uintptr) (_2 int32)

//go:noescape
func Ygetrlimit(t *TLS, resource int32, rlim uintptr) (_3 int32)

//go:noescape
func Ygetrusage(t *TLS, who int32, usage uintptr) (_3 int32)

//go:noescape
func Ygetservbyname(t *TLS, name, proto uintptr) (_2 uintptr)

//go:noescape
func Ygetsockname(t *TLS, sockfd int32, addr, addrlen uintptr) (_3 int32)

//go:noescape
func Ygetsockopt(t *TLS, sockfd, level, optname int32, optval, optlen uintptr) (_3 int32)

//go:noescape
func Ygettimeofday(t *TLS, tv, tz uintptr) (_2 int32)

//go:noescape
func Ygetuid(t *TLS) (_1 types.Uid_t)

//go:noescape
func Ygmtime(tls *TLS, t uintptr) (r uintptr)

//go:noescape
func Ygmtime_r(tls *TLS, t uintptr, tm uintptr) (r uintptr)

//go:noescape
func Ygrantpt(t *TLS, fd int32) (_2 int32)

//go:noescape
func Yhtonl(t *TLS, hostlong uint32) (_2 uint32)

//go:noescape
func Yhtons(t *TLS, hostshort uint16) (_2 uint16)

//go:noescape
func Yhypot(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Yhypotf(t *TLS, x, y float32) (_2 float32)

//go:noescape
func Yilogb(t *TLS, x float64) (_2 int32)

//go:noescape
func Yinet_ntoa(t *TLS, in1 in.In_addr) (_2 uintptr)

//go:noescape
func Yinet_ntop(tls *TLS, af int32, a0 uintptr, s uintptr, l socklen_t) (_5 uintptr)

//go:noescape
func Yinet_pton(tls *TLS, af int32, s uintptr, a0 uintptr) (_4 int32)

//go:noescape
func Yinitstate(t *TLS, seed uint32, statebuf uintptr, statelen types.Size_t) (_4 uintptr)

//go:noescape
func Yinitstate_r(t *TLS, seed uint32, statebuf uintptr, statelen types.Size_t, buf uintptr) (_5 int32)

//go:noescape
func Yioctl(t *TLS, fd int32, request ulong, va uintptr) (_4 int32)

//go:noescape
func Yisalnum(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisalpha(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisascii(t *TLS, c int32) (_2 int32)

//go:noescape
func Yisatty(t *TLS, fd int32) (_2 int32)

//go:noescape
func Yisblank(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yiscntrl(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisdigit(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisgraph(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yishexnumber(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisideogram(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yislower(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisnan(t *TLS, x float64) (_2 int32)

//go:noescape
func Yisnanf(t *TLS, x float32) (_2 int32)

//go:noescape
func Yisnanl(t *TLS, x float64) (_2 int32)

//go:noescape
func Yisnumber(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisphonogram(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisprint(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yispunct(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisrune(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yissetugid(t *TLS) (_1 int32)

//go:noescape
func Yisspace(t *TLS, c int32) (_2 int32)

//go:noescape
func Yisspecial(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yisupper(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Yiswalnum(t *TLS, wc wctype.Wint_t) (_2 int32)

//go:noescape
func Yiswspace(t *TLS, wc wctype.Wint_t) (_2 int32)

//go:noescape
func Yisxdigit(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Ykevent(t *TLS, kq int32, changelist uintptr, nchanges int32, eventlist uintptr, nevents int32, timeout uintptr) (_7 int32)

//go:noescape
func Ykill(t *TLS, pid types.Pid_t, sig int32) (_3 int32)

//go:noescape
func Ykqueue(t *TLS) (_1 int32)

//go:noescape
func Ylabs(t *TLS, j long) (_2 long)

//go:noescape
func Ylchmod(t *TLS, path uintptr, mode types.Mode_t) (_3 int32)

//go:noescape
func Yldexp(t *TLS, x float64, exp int32) (_3 float64)

//go:noescape
func Yldiv(tls *TLS, num, den long) (r Tldiv_t)

//go:noescape
func Ylgamma(t *TLS, x float64) (_2 float64)

//go:noescape
func Ylink(t *TLS, oldpath, newpath uintptr) (_2 int32)

//go:noescape
func Ylinkat(t *TLS, olddirfd int32, oldpath uintptr, newdirfd int32, newpath uintptr, flags int32) (_6 int32)

//go:noescape
func Ylisten(t *TLS, sockfd, backlog int32) (_2 int32)

//go:noescape
func Yllabs(tls *TLS, a int64) (_2 int64)

//go:noescape
func Ylocaleconv(tls *TLS) (_1 uintptr)

//go:noescape
func Ylocaltime(_ *TLS, timep uintptr) (_2 uintptr)

//go:noescape
func Ylocaltime_r(_ *TLS, timep, result uintptr) (_2 uintptr)

//go:noescape
func Ylog(t *TLS, x float64) (_2 float64)

//go:noescape
func Ylog10(t *TLS, x float64) (_2 float64)

//go:noescape
func Ylog10f(t *TLS, x float32) (_2 float32)

//go:noescape
func Ylog1p(tls *TLS, x3 float64) (r float64)

//go:noescape
func Ylog1pf(tls *TLS, x3 float32) (r float32)

//go:noescape
func Ylog1pl(tls *TLS, x float64) (r float64)

//go:noescape
func Ylog2(t *TLS, x float64) (_2 float64)

//go:noescape
func Ylog2f(t *TLS, x float32) (_2 float32)

//go:noescape
func Ylogb(t *TLS, x float64) (_2 float64)

//go:noescape
func Ylogf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ylongjmp(t *TLS, env uintptr, val int32)

//go:noescape
func Ylrand48(t *TLS) (_1 long)

//go:noescape
func Ylrint(tls *TLS, x float64) (r long)

//go:noescape
func Ylrintf(tls *TLS, x float32) (r long)

//go:noescape
func Ylround(tls *TLS, x float64) (r long)

//go:noescape
func Ylseek(t *TLS, fd int32, offset types.Off_t, whence int32) (_4 types.Off_t)

//go:noescape
func Ylseek64(t *TLS, fd int32, offset types.Off_t, whence int32) (_4 types.Off_t)

//go:noescape
func Ylstat(t *TLS, pathname, statbuf uintptr) (_2 int32)

//go:noescape
func Ylstat64(t *TLS, pathname, statbuf uintptr) (_2 int32)

//go:noescape
func Ymach_absolute_time(t *TLS) (_1 uint64)

//go:noescape
func Ymach_timebase_info(t *TLS, info uintptr) (_2 int32)

//go:noescape
func Ymalloc(t *TLS, n types.Size_t) (_2 uintptr)

//go:noescape
func Ymalloc_size(t *TLS, p uintptr) (r types.Size_t)

//go:noescape
func Ymalloc_usable_size(tls *TLS, p uintptr) (r types.Size_t)

//go:noescape
func Ymblen(t *TLS, s uintptr, n types.Size_t) (_3 int32)

//go:noescape
func Ymbstowcs(t *TLS, dest, src uintptr, n types.Size_t) (_3 types.Size_t)

//go:noescape
func Ymbtowc(t *TLS, pwc, s uintptr, n types.Size_t) (_3 int32)

//go:noescape
func Ymemchr(t *TLS, s uintptr, c int32, n types.Size_t) (_4 uintptr)

//go:noescape
func Ymemcmp(t *TLS, s1, s2 uintptr, n types.Size_t) (_3 int32)

//go:noescape
func Ymemcpy(t *TLS, dest, src uintptr, n types.Size_t) (r uintptr)

//go:noescape
func Ymemmove(t *TLS, dest, src uintptr, n types.Size_t) (_3 uintptr)

//go:noescape
func Ymemset(t *TLS, dest uintptr, c int32, n types.Size_t) (_4 uintptr)

//go:noescape
func Ymkdir(t *TLS, path uintptr, mode types.Mode_t) (_3 int32)

//go:noescape
func Ymkdirat(t *TLS, dirfd int32, pathname uintptr, mode types.Mode_t) (_4 int32)

//go:noescape
func Ymkfifo(t *TLS, pathname uintptr, mode types.Mode_t) (_3 int32)

//go:noescape
func Ymknod(t *TLS, pathname uintptr, mode types.Mode_t, dev types.Dev_t) (_4 int32)

//go:noescape
func Ymkostemp(t *TLS, template uintptr, flags int32) (_3 int32)

//go:noescape
func Ymkstemp(t *TLS, template uintptr) (_2 int32)

//go:noescape
func Ymkstemp64(t *TLS, template uintptr) (_2 int32)

//go:noescape
func Ymkstemps(t *TLS, template uintptr, suffixlen int32) (_3 int32)

//go:noescape
func Ymkstemps64(t *TLS, template uintptr, suffixlen int32) (_3 int32)

//go:noescape
func Ymktime(t *TLS, ptm uintptr) (r time.Time_t)

//go:noescape
func Ymmap(t *TLS, addr uintptr, length types.Size_t, prot, flags, fd int32, offset types.Off_t) (_5 uintptr)

//go:noescape
func Ymodf(t *TLS, x float64, iptr uintptr) (_3 float64)

//go:noescape
func Ymodff(t *TLS, x float32, iptr uintptr) (_3 float32)

//go:noescape
func Ymunmap(t *TLS, addr uintptr, length types.Size_t) (_3 int32)

//go:noescape
func Ynanf(tls *TLS, s uintptr) (_2 float32)

//go:noescape
func Ynanosleep(t *TLS, req, rem uintptr) (_2 int32)

//go:noescape
func Ynextafter(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Ynextafterf(t *TLS, x, y float32) (_2 float32)

//go:noescape
func Ynexttoward(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Ynl_langinfo(t *TLS, item langinfo.Nl_item) (_2 uintptr)

//go:noescape
func Yntohs(t *TLS, netshort uint16) (_2 uint16)

//go:noescape
func Yobstack_free(t *TLS, obstack, obj uintptr)

//go:noescape
func Yobstack_vprintf(t *TLS, obstack, template, va uintptr) (_2 int32)

//go:noescape
func Yopen(t *TLS, pathname uintptr, flags int32, args uintptr) (_4 int32)

//go:noescape
func Yopen64(t *TLS, pathname uintptr, flags int32, args uintptr) (_4 int32)

//go:noescape
func Yopenat(t *TLS, dirfd int32, pathname uintptr, flags int32, args uintptr) (_5 int32)

//go:noescape
func Yopendir(t *TLS, name uintptr) (_2 uintptr)

//go:noescape
func Yopenpty(t *TLS, amaster, aslave, name, termp, winp uintptr) (_2 int32)

//go:noescape
func Ypathconf(t *TLS, path uintptr, name int32) (_3 long)

//go:noescape
func Ypause(t *TLS) (_1 int32)

//go:noescape
func Ypclose(t *TLS, stream uintptr) (r int32)

//go:noescape
func Yperror(tls *TLS, msg uintptr)

//go:noescape
func Ypipe(t *TLS, pipefd uintptr) (_2 int32)

//go:noescape
func Ypoll(t *TLS, fds uintptr, nfds poll.Nfds_t, timeout int32) (_4 int32)

//go:noescape
func Ypopen(t *TLS, command, mode uintptr) (r uintptr)

//go:noescape
func Yposix_fadvise(t *TLS, fd int32, offset, len types.Off_t, advice int32) (_4 int32)

//go:noescape
func Yposix_openpt(t *TLS, flags int32) (_2 int32)

//go:noescape
func Ypow(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Ypowf(t *TLS, x, y float32) (_2 float32)

//go:noescape
func Ypread(t *TLS, fd int32, buf uintptr, count types.Size_t, offset types.Off_t) (_5 types.Ssize_t)

//go:noescape
func Yprintf(t *TLS, format, args uintptr) (_2 int32)

//go:noescape
func Ypselect(t *TLS, nfds int32, readfds, writefds, exceptfds, timeout, sigmask uintptr) (_3 int32)

//go:noescape
func Ypthread_attr_destroy(t *TLS, pAttr uintptr) (_2 int32)

//go:noescape
func Ypthread_attr_getdetachstate(tls *TLS, a uintptr, state uintptr) (_3 int32)

//go:noescape
func Ypthread_attr_init(t *TLS, pAttr uintptr) (_2 int32)

//go:noescape
func Ypthread_attr_setdetachstate(tls *TLS, a uintptr, state int32) (r int32)

//go:noescape
func Ypthread_attr_setscope(t *TLS, pAttr uintptr, contentionScope int32) (_3 int32)

//go:noescape
func Ypthread_attr_setstacksize(t *TLS, attr uintptr, stackSize types.Size_t) (_3 int32)

//go:noescape
func Ypthread_cond_broadcast(t *TLS, pCond uintptr) (_2 int32)

//go:noescape
func Ypthread_cond_destroy(t *TLS, pCond uintptr) (_2 int32)

//go:noescape
func Ypthread_cond_init(t *TLS, pCond, pAttr uintptr) (_2 int32)

//go:noescape
func Ypthread_cond_signal(t *TLS, pCond uintptr) (_2 int32)

//go:noescape
func Ypthread_cond_timedwait(t *TLS, pCond, pMutex, pAbsTime uintptr) (_2 int32)

//go:noescape
func Ypthread_cond_wait(t *TLS, pCond, pMutex uintptr) (_2 int32)

//go:noescape
func Ypthread_detach(t *TLS, thread pthread.Pthread_t) (_2 int32)

//go:noescape
func Ypthread_equal(t *TLS, t1, t2 pthread.Pthread_t) (_2 int32)

//go:noescape
func Ypthread_exit(t *TLS, value uintptr)

//go:noescape
func Ypthread_getname_np(t *TLS, thread pthread.Pthread_t, name uintptr, len1 types.Size_t) (_4 int32)

//go:noescape
func Ypthread_getspecific(t *TLS, key pthread.Pthread_key_t) (_2 uintptr)

//go:noescape
func Ypthread_join(t *TLS, thread pthread.Pthread_t, pValue uintptr) (_3 int32)

//go:noescape
func Ypthread_key_delete(t *TLS, key pthread.Pthread_key_t) (_2 int32)

//go:noescape
func Ypthread_mutex_destroy(t *TLS, pMutex uintptr) (_2 int32)

//go:noescape
func Ypthread_mutex_init(t *TLS, pMutex, pAttr uintptr) (_2 int32)

//go:noescape
func Ypthread_mutex_lock(t *TLS, pMutex uintptr) (_2 int32)

//go:noescape
func Ypthread_mutex_trylock(t *TLS, pMutex uintptr) (_2 int32)

//go:noescape
func Ypthread_mutex_unlock(t *TLS, pMutex uintptr) (_2 int32)

//go:noescape
func Ypthread_mutexattr_destroy(tls *TLS, a uintptr) (_2 int32)

//go:noescape
func Ypthread_mutexattr_init(tls *TLS, a uintptr) (_2 int32)

//go:noescape
func Ypthread_mutexattr_settype(tls *TLS, a uintptr, type1 int32) (_3 int32)

//go:noescape
func Ypthread_self(t *TLS) (_1 pthread.Pthread_t)

//go:noescape
func Ypthread_setname_np(t *TLS, name uintptr) (_2 int32)

//go:noescape
func Ypthread_setspecific(t *TLS, key pthread.Pthread_key_t, value uintptr) (_3 int32)

//go:noescape
func Ypthread_sigmask(tls *TLS, now int32, set, old uintptr) (_3 int32)

//go:noescape
func Yptsname(t *TLS, fd int32) (_2 uintptr)

//go:noescape
func Yptsname_r(t *TLS, fd int32, buf uintptr, buflen types.Size_t) (_4 int32)

//go:noescape
func Yputc(t *TLS, c int32, fp uintptr) (_3 int32)

//go:noescape
func Yputc_unlocked(t *TLS, c int32, stream uintptr) (_3 int32)

//go:noescape
func Yputchar(t *TLS, c int32) (_2 int32)

//go:noescape
func Yputenv(tls *TLS, s uintptr) (_2 int32)

//go:noescape
func Yputs(t *TLS, s uintptr) (_2 int32)

//go:noescape
func Ypwrite(t *TLS, fd int32, buf uintptr, count types.Size_t, offset types.Off_t) (_5 types.Ssize_t)

//go:noescape
func Yqsort(t *TLS, base uintptr, nmemb, size types.Size_t, compar uintptr)

//go:noescape
func __ccgo_abi0_qsort_3(_0 *TLS, _1 uintptr, _2 uintptr, __ccgo_fp uintptr) (_3 int32)

func __ccgo_abiInternal_qsort_3(tls *TLS, dest, abi0CodePtr uintptr) {
	f := func(_0 *TLS, _1 uintptr, _2 uintptr) (_3 int32) {
		return __ccgo_abi0_qsort_3(_0, _1, _2, abi0CodePtr)
	}
	*(*[2]uintptr)(unsafe.Pointer(dest)) = *(*[2]uintptr)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&f))))
}

//go:noescape
func Yraise(t *TLS, sig int32) (_2 int32)

//go:noescape
func Yrand(t *TLS) (_1 int32)

//go:noescape
func Yrand_r(tls *TLS, seed uintptr) (_2 int32)

//go:noescape
func Yrandom(t *TLS) (_1 long)

//go:noescape
func Yrandom_r(t *TLS, buf, result uintptr) (_2 int32)

//go:noescape
func Yread(t *TLS, fd int32, buf uintptr, count types.Size_t) (_4 types.Ssize_t)

//go:noescape
func Yreaddir(t *TLS, dir uintptr) (_2 uintptr)

//go:noescape
func Yreadlink(t *TLS, path, buf uintptr, bufsize types.Size_t) (_3 types.Ssize_t)

//go:noescape
func Yreadlinkat(t *TLS, dirfd int32, pathname, buf uintptr, bufsiz types.Size_t) (_4 types.Ssize_t)

//go:noescape
func Yreadv(t *TLS, fd int32, iov uintptr, iovcnt int32) (_4 types.Ssize_t)

//go:noescape
func Yrealloc(t *TLS, ptr uintptr, size types.Size_t) (_3 uintptr)

//go:noescape
func Yreallocarray(t *TLS, ptr uintptr, nmemb, size size_t) (_3 uintptr)

//go:noescape
func Yrealpath(t *TLS, path, resolved_path uintptr) (_2 uintptr)

//go:noescape
func Yrecv(t *TLS, sockfd int32, buf uintptr, len types.Size_t, flags int32) (_5 types.Ssize_t)

//go:noescape
func Yrecvfrom(t *TLS, sockfd int32, buf uintptr, len types.Size_t, flags int32, src_addr, addrlen uintptr) (_6 types.Ssize_t)

//go:noescape
func Yrecvmsg(t *TLS, sockfd int32, msg uintptr, flags int32) (_4 types.Ssize_t)

//go:noescape
func Yremainder(t *TLS, x, y float64) (_2 float64)

//go:noescape
func Yremove(t *TLS, pathname uintptr) (_2 int32)

//go:noescape
func Yremquo(t *TLS, x, y float64, quo uintptr) (_3 float64)

//go:noescape
func Yrename(t *TLS, oldpath, newpath uintptr) (_2 int32)

//go:noescape
func Yrenameat(t *TLS, olddirfd int32, oldpath uintptr, newdirfd int32, newpath uintptr) (_5 int32)

//go:noescape
func Yrewind(t *TLS, stream uintptr)

//go:noescape
func Yrewinddir(tls *TLS, f uintptr)

//go:noescape
func Yrindex(t *TLS, s uintptr, c int32) (_3 uintptr)

//go:noescape
func Yrint(tls *TLS, x float64) (_2 float64)

//go:noescape
func Yrintf(tls *TLS, x float32) (r float32)

//go:noescape
func Yrmdir(t *TLS, pathname uintptr) (_2 int32)

//go:noescape
func Yround(t *TLS, x float64) (_2 float64)

//go:noescape
func Yroundf(t *TLS, x float32) (_2 float32)

//go:noescape
func Yscalbn(tls *TLS, x float64, n int32) (_3 float64)

//go:noescape
func Yscalbnl(tls *TLS, x float64, n int32) (_3 float64)

//go:noescape
func Ysched_yield(t *TLS) (_1 int32)

//go:noescape
func Yselect(t *TLS, nfds int32, readfds, writefds, exceptfds, timeout uintptr) (_3 int32)

//go:noescape
func Ysend(t *TLS, sockfd int32, buf uintptr, len types.Size_t, flags int32) (_5 types.Ssize_t)

//go:noescape
func Ysendmsg(t *TLS, sockfd int32, msg uintptr, flags int32) (_4 types.Ssize_t)

//go:noescape
func Ysendto(t *TLS, sockfd int32, buf uintptr, len types.Size_t, flags int32, src_addr uintptr, addrlen socklen_t) (_7 types.Ssize_t)

//go:noescape
func Ysetattrlist(t *TLS, path, attrList, attrBuf uintptr, attrBufSize types.Size_t, options uint32) (_4 int32)

//go:noescape
func Ysetbuf(t *TLS, stream, buf uintptr)

//go:noescape
func Ysetenv(tls *TLS, var1 uintptr, value uintptr, overwrite int32) (_4 int32)

//go:noescape
func Ysetgid(tls *TLS, gid uint32) (r int32)

//go:noescape
func Ysetjmp(t *TLS, env uintptr) (_2 int32)

//go:noescape
func Ysetlocale(t *TLS, category int32, name uintptr) (r uintptr)

//go:noescape
func Ysetrlimit(t *TLS, resource int32, rlim uintptr) (_3 int32)

//go:noescape
func Ysetsid(t *TLS) (_1 types.Pid_t)

//go:noescape
func Ysetsockopt(t *TLS, sockfd, level, optname int32, optval uintptr, optlen socket.Socklen_t) (_4 int32)

//go:noescape
func Ysetstate(t *TLS, state uintptr) (_2 uintptr)

//go:noescape
func Ysetuid(tls *TLS, uid uint32) (r int32)

//go:noescape
func Ysetvbuf(t *TLS, stream, buf uintptr, mode int32, size types.Size_t) (_4 int32)

//go:noescape
func Yshmat(t *TLS, shmid int32, shmaddr uintptr, shmflg int32) (_4 uintptr)

//go:noescape
func Yshmctl(t *TLS, shmid, cmd int32, buf uintptr) (_3 int32)

//go:noescape
func Yshmdt(t *TLS, shmaddr uintptr) (_2 int32)

//go:noescape
func Yshutdown(t *TLS, sockfd, how int32) (_2 int32)

//go:noescape
func Ysigaddset(t *TLS, set uintptr, signum int32) (_3 int32)

//go:noescape
func Ysigdelset(t *TLS, set uintptr, signum int32) (_3 int32)

//go:noescape
func Ysigemptyset(t *TLS, set uintptr) (_2 int32)

//go:noescape
func Ysigfillset(t *TLS, set uintptr) (_2 int32)

//go:noescape
func Ysigismember(t *TLS, set uintptr, signum int32) (_3 int32)

//go:noescape
func Ysin(t *TLS, x float64) (_2 float64)

//go:noescape
func Ysinf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ysinh(t *TLS, x float64) (_2 float64)

//go:noescape
func Ysinhf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ysleep(t *TLS, seconds uint32) (_2 uint32)

//go:noescape
func Ysnprintf(t *TLS, str uintptr, size types.Size_t, format, args uintptr) (r int32)

//go:noescape
func Ysocket(t *TLS, domain, type1, protocol int32) (_2 int32)

//go:noescape
func Ysprintf(t *TLS, str, format, args uintptr) (r int32)

//go:noescape
func Ysqrt(t *TLS, x float64) (_2 float64)

//go:noescape
func Ysqrtf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ysrand48(t *TLS, seedval long)

//go:noescape
func Ysrandomdev(t *TLS)

//go:noescape
func Ysscanf(t *TLS, str, format, va uintptr) (_2 int32)

//go:noescape
func Ystat(t *TLS, pathname, statbuf uintptr) (_2 int32)

//go:noescape
func Ystat64(t *TLS, pathname, statbuf uintptr) (_2 int32)

//go:noescape
func Ystatfs(t *TLS, path uintptr, buf uintptr) (_3 int32)

//go:noescape
func Ystrcasecmp(t *TLS, s1, s2 uintptr) (_2 int32)

//go:noescape
func Ystrcat(t *TLS, dest, src uintptr) (r uintptr)

//go:noescape
func Ystrchr(t *TLS, s uintptr, c int32) (_3 uintptr)

//go:noescape
func Ystrchrnul(tls *TLS, s uintptr, c int32) (r uintptr)

//go:noescape
func Ystrcmp(t *TLS, s1, s2 uintptr) (_2 int32)

//go:noescape
func Ystrcoll(t *TLS, s1, s2 uintptr) (_2 int32)

//go:noescape
func Ystrcpy(t *TLS, dest, src uintptr) (r uintptr)

//go:noescape
func Ystrcspn(t *TLS, s, reject uintptr) (r types.Size_t)

//go:noescape
func Ystrdup(tls *TLS, s uintptr) (_2 uintptr)

//go:noescape
func Ystrerror(t *TLS, errnum int32) (_2 uintptr)

//go:noescape
func Ystrerror_r(t *TLS, errnum int32, buf uintptr, buflen types.Size_t) (r int32)

//go:noescape
func Ystrftime(tls *TLS, s uintptr, n size_t, f uintptr, tm uintptr) (r size_t)

//go:noescape
func Ystrlcat(tls *TLS, d uintptr, s uintptr, n size_t) (_4 size_t)

//go:noescape
func Ystrlcpy(tls *TLS, d uintptr, s uintptr, n size_t) (_4 size_t)

//go:noescape
func Ystrlen(t *TLS, s uintptr) (r Tsize_t)

//go:noescape
func Ystrncasecmp(tls *TLS, _l uintptr, _r uintptr, n size_t) (_4 int32)

//go:noescape
func Ystrncat(tls *TLS, d uintptr, s uintptr, n size_t) (_4 uintptr)

//go:noescape
func Ystrncmp(t *TLS, s1, s2 uintptr, n types.Size_t) (_3 int32)

//go:noescape
func Ystrncpy(t *TLS, dest, src uintptr, n types.Size_t) (r uintptr)

//go:noescape
func Ystrnlen(tls *TLS, s uintptr, n size_t) (_3 size_t)

//go:noescape
func Ystrpbrk(t *TLS, s, accept uintptr) (_2 uintptr)

//go:noescape
func Ystrrchr(t *TLS, s uintptr, c int32) (r uintptr)

//go:noescape
func Ystrspn(tls *TLS, s uintptr, c uintptr) (_3 size_t)

//go:noescape
func Ystrstr(t *TLS, haystack, needle uintptr) (_2 uintptr)

//go:noescape
func Ystrtod(tls *TLS, s uintptr, p uintptr) (_3 float64)

//go:noescape
func Ystrtof(tls *TLS, s uintptr, p uintptr) (_3 float32)

//go:noescape
func Ystrtoimax(tls *TLS, s uintptr, p uintptr, base int32) (_4 intmax_t)

//go:noescape
func Ystrtok(tls *TLS, s uintptr, sep uintptr) (_3 uintptr)

//go:noescape
func Ystrtol(tls *TLS, s uintptr, p uintptr, base int32) (_4 int64)

//go:noescape
func Ystrtold(tls *TLS, s uintptr, p uintptr) (_3 float64)

//go:noescape
func Ystrtoll(tls *TLS, s uintptr, p uintptr, base int32) (_4 int64)

//go:noescape
func Ystrtoul(tls *TLS, s uintptr, p uintptr, base int32) (_4 uint64)

//go:noescape
func Ystrtoull(tls *TLS, s uintptr, p uintptr, base int32) (_4 uint64)

//go:noescape
func Ystrtoumax(tls *TLS, s uintptr, p uintptr, base int32) (_4 uintmax_t)

//go:noescape
func Ysymlink(t *TLS, target, linkpath uintptr) (_2 int32)

//go:noescape
func Ysymlinkat(t *TLS, target uintptr, newdirfd int32, linkpath uintptr) (_4 int32)

//go:noescape
func Ysysconf(t *TLS, name int32) (_2 long)

//go:noescape
func Ysysctlbyname(t *TLS, name, oldp, oldlenp, newp uintptr, newlen types.Size_t) (_3 int32)

//go:noescape
func Ysystem(t *TLS, command uintptr) (r int32)

//go:noescape
func Ytan(t *TLS, x float64) (_2 float64)

//go:noescape
func Ytanf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ytanh(t *TLS, x float64) (_2 float64)

//go:noescape
func Ytanhf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ytcgetattr(t *TLS, fd int32, termios_p uintptr) (_3 int32)

//go:noescape
func Ytcsendbreak(t *TLS, fd, duration int32) (_2 int32)

//go:noescape
func Ytcsetattr(t *TLS, fd, optional_actions int32, termios_p uintptr) (_3 int32)

//go:noescape
func Ytgamma(t *TLS, x float64) (_2 float64)

//go:noescape
func Ytime(t *TLS, tloc uintptr) (_2 types.Time_t)

//go:noescape
func Ytmpfile(t *TLS) (_1 uintptr)

//go:noescape
func Ytmpnam(t *TLS, s uintptr) (r uintptr)

//go:noescape
func Ytoascii(tls *TLS, _c int32) (_2 int32)

//go:noescape
func Ytolower(t *TLS, c int32) (_2 int32)

//go:noescape
func Ytoupper(t *TLS, c int32) (_2 int32)

//go:noescape
func Ytrunc(t *TLS, x float64) (_2 float64)

//go:noescape
func Ytruncate(t *TLS, path uintptr, length types.Off_t) (_3 int32)

//go:noescape
func Ytruncf(t *TLS, x float32) (_2 float32)

//go:noescape
func Ytzset(t *TLS)

//go:noescape
func Yumask(t *TLS, mask types.Mode_t) (_2 types.Mode_t)

//go:noescape
func Yuname(t *TLS, buf uintptr) (_2 int32)

//go:noescape
func Yungetc(t *TLS, c int32, stream uintptr) (_3 int32)

//go:noescape
func Yunlink(t *TLS, pathname uintptr) (_2 int32)

//go:noescape
func Yunlinkat(t *TLS, dirfd int32, pathname uintptr, flags int32) (_4 int32)

//go:noescape
func Yunlockpt(t *TLS, fd int32) (_2 int32)

//go:noescape
func Yunsetenv(tls *TLS, name uintptr) (_2 int32)

//go:noescape
func Yusleep(t *TLS, usec types.Useconds_t) (_2 int32)

//go:noescape
func Yutime(t *TLS, filename, times uintptr) (_2 int32)

//go:noescape
func Yutimensat(t *TLS, dirfd int32, pathname, times uintptr, flags int32) (_4 int32)

//go:noescape
func Yutimes(t *TLS, filename, times uintptr) (_2 int32)

//go:noescape
func Yuuid_copy(t *TLS, dst, src uintptr)

//go:noescape
func Yuuid_generate_random(t *TLS, out uintptr)

//go:noescape
func Yuuid_parse(t *TLS, in uintptr, uu uintptr) (_3 int32)

//go:noescape
func Yuuid_unparse(t *TLS, uu, out uintptr)

//go:noescape
func Yvasprintf(t *TLS, strp, fmt, ap uintptr) (_2 int32)

//go:noescape
func Yvfprintf(t *TLS, stream, format, ap uintptr) (_2 int32)

//go:noescape
func Yvprintf(t *TLS, s, ap uintptr) (_2 int32)

//go:noescape
func Yvsnprintf(t *TLS, str uintptr, size types.Size_t, format, va uintptr) (_4 int32)

//go:noescape
func Yvsprintf(t *TLS, str, format, va uintptr) (_2 int32)

//go:noescape
func Ywaitpid(t *TLS, pid types.Pid_t, wstatus uintptr, optname int32) (_4 types.Pid_t)

//go:noescape
func Ywcschr(t *TLS, wcs uintptr, wc wchar_t) (_3 uintptr)

//go:noescape
func Ywctomb(t *TLS, s uintptr, wc wchar_t) (_3 int32)

//go:noescape
func Ywcwidth(t *TLS, c wchar_t) (_2 int32)

//go:noescape
func Ywrite(t *TLS, fd int32, buf uintptr, count types.Size_t) (_4 types.Ssize_t)

//go:noescape
func Ywritev(t *TLS, fd int32, iov uintptr, iovcnt int32) (_4 types.Ssize_t)
