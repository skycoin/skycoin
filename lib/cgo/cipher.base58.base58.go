package main

import (
	"encoding/hex"
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
	s := _s
	__arg1, ____return_err := hex.DecodeString(s)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}

	return
}

//export SKY_base58_Base582Hex
func SKY_base58_Base582Hex(_b string, _arg1 *C.GoSlice_) (____error_code uint32) {
	b := _b
	__arg1, ____return_err := base58.Decode(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

//export SKY_base58_Hex2Base58
func SKY_base58_Hex2Base58(_val []byte, _arg1 *C.GoString_) (____error_code uint32) {
	val := *(*[]byte)(unsafe.Pointer(&_val))
	__arg1 := base58.Encode(val)
	copyString(string(__arg1), _arg1)
	return
}
