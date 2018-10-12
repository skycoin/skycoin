package main

import (
	"reflect"
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher/ripemd160"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_ripemd160_New
func SKY_ripemd160_New(handle *C.Hash_Handle) (____error_code uint32) {
	hash := ripemd160.New()
	*handle = registerHashHandle(&hash)
	return
}

//export SKY_ripemd160_Write
func SKY_ripemd160_Write(handle C.Hash_Handle, _p []byte, _nn *int) (____error_code uint32) {
	h, ok := lookupHashHandle(handle)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	p := *(*[]byte)(unsafe.Pointer(&_p))
	nn, error := (*h).Write(p)
	if error != nil {
		____error_code = SKY_ERROR
		return
	}
	*_nn = nn
	return
}

//export SKY_ripemd160_Sum
func SKY_ripemd160_Sum(handle C.Hash_Handle, _p []byte, _arg1 *C.GoSlice_) (____error_code uint32) {
	h, ok := lookupHashHandle(handle)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	p := *(*[]byte)(unsafe.Pointer(&_p))
	__arg1 := (*h).Sum(p)
	copyToGoSlice(reflect.ValueOf(__arg1[:]), _arg1)
	return
}
