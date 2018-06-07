package main

import (
	"reflect"
	"unsafe"

	base58 "github.com/skycoin/skycoin/src/cipher/base58"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_base58_String2Hex
func SKY_base58_String2Hex(_s string, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	s := _s
	__arg1 := base58.String2Hex(s)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_base58_Base58_ToInt
func SKY_base58_Base58_ToInt(_b string, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := (base58.Base58)(_b)
	__arg0, ____return_err := b.ToInt()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

//export SKY_base58_Base58_ToHex
func SKY_base58_Base58_ToHex(_b string, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := (base58.Base58)(_b)
	__arg0, ____return_err := b.ToHex()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

//export SKY_base58_Base58_Base582Int
func SKY_base58_Base58_Base582Int(_b string, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := (base58.Base58)(_b)
	__arg0, ____return_err := b.Base582Int()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

//export SKY_base58_Base582Hex
func SKY_base58_Base582Hex(_b string, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := _b
	__arg1, ____return_err := base58.Base582Hex(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

//export SKY_base58_Base58_BitHex
func SKY_base58_Base58_BitHex(_b string, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := (base58.Base58)(_b)
	__arg0, ____return_err := b.BitHex()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

//export SKY_base58_Int2Base58
func SKY_base58_Int2Base58(_val int, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	val := _val
	__arg1 := base58.Int2Base58(val)
	copyString(string(__arg1), _arg1)
	return
}

//export SKY_base58_Hex2Base58
func SKY_base58_Hex2Base58(_val []byte, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	val := *(*[]byte)(unsafe.Pointer(&_val))
	__arg1 := base58.Hex2Base58(val)
	copyString(string(__arg1), _arg1)
	return
}

//export SKY_base58_Hex2Base58String
func SKY_base58_Hex2Base58String(_val []byte, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	val := *(*[]byte)(unsafe.Pointer(&_val))
	__arg1 := base58.Hex2Base58String(val)
	copyString(__arg1, _arg1)
	return
}

//export SKY_base58_Hex2Base58Str
func SKY_base58_Hex2Base58Str(_val []byte, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	val := *(*[]byte)(unsafe.Pointer(&_val))
	__arg1 := base58.Hex2Base58Str(val)
	copyString(__arg1, _arg1)
	return
}
