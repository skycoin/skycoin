package main

import (
	base58 "github.com/skycoin/skycoin/src/base58"
	cipher "github.com/skycoin/skycoin/src/cipher"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_base58_Hex2Big
func SKY_base58_Hex2Big(_b *C.GoSlice_, _arg1 *C.Int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(_b))
	__arg1 := base58.Hex2Big(b)
	return
}

// export SKY_base58_String2Hex
func SKY_base58_String2Hex(_s string, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := _s
	__arg1 := base58.String2Hex(s)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_base58_Base58_ToBig
func SKY_base58_Base58_ToBig(_b *C.Base58, _arg0 *C.Int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Base58)(unsafe.Pointer(_b))
	__arg0, ____return_err := b.ToBig()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_base58_Base58_ToInt
func SKY_base58_Base58_ToInt(_b *C.Base58, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Base58)(unsafe.Pointer(_b))
	__arg0, ____return_err := b.ToInt()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

// export SKY_base58_Base58_ToHex
func SKY_base58_Base58_ToHex(_b *C.Base58, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Base58)(unsafe.Pointer(_b))
	__arg0, ____return_err := b.ToHex()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_base58_Base58_Base582Big
func SKY_base58_Base58_Base582Big(_b *C.Base58, _arg0 *C.Int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Base58)(unsafe.Pointer(_b))
	__arg0, ____return_err := b.Base582Big()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_base58_Base58_Base582Int
func SKY_base58_Base58_Base582Int(_b *C.Base58, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Base58)(unsafe.Pointer(_b))
	__arg0, ____return_err := b.Base582Int()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

// export SKY_base58_Base582Hex
func SKY_base58_Base582Hex(_b string, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := _b
	__arg1, ____return_err := base58.Base582Hex(b)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_base58_Base58_BitHex
func SKY_base58_Base58_BitHex(_b *C.Base58, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Base58)(unsafe.Pointer(_b))
	__arg0, ____return_err := b.BitHex()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_base58_Big2Base58
func SKY_base58_Big2Base58(_val *C.Int, _arg1 *C.Base58) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1 := base58.Big2Base58(val)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofBase58))
	return
}

// export SKY_base58_Int2Base58
func SKY_base58_Int2Base58(_val int, _arg1 *C.Base58) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	val := _val
	__arg1 := base58.Int2Base58(val)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofBase58))
	return
}

// export SKY_base58_Hex2Base58
func SKY_base58_Hex2Base58(_val *C.GoSlice_, _arg1 *C.Base58) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	val := *(*[]byte)(unsafe.Pointer(_val))
	__arg1 := base58.Hex2Base58(val)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofBase58))
	return
}

// export SKY_base58_Hex2Base58String
func SKY_base58_Hex2Base58String(_val *C.GoSlice_, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	val := *(*[]byte)(unsafe.Pointer(_val))
	__arg1 := base58.Hex2Base58String(val)
	copyString(__arg1, _arg1)
	return
}

// export SKY_base58_Hex2Base58Str
func SKY_base58_Hex2Base58Str(_val *C.GoSlice_, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	val := *(*[]byte)(unsafe.Pointer(_val))
	__arg1 := base58.Hex2Base58Str(val)
	copyString(__arg1, _arg1)
	return
}
