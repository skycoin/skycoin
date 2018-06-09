package main

import (
	"reflect"
	"unsafe"

	secp256k1go "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_secp256k1_SumSHA256
func SKY_secp256k1_SumSHA256(_b []byte, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1 := secp256k1go.SumSHA256(b)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_secp256k1_RandByte
func SKY_secp256k1_RandByte(_n int, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	n := _n
	__arg1 := secp256k1go.RandByte(n)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}
