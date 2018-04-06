package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	secp256k1 "github.com/skycoin/skycoin/src/secp256k1"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_secp256k1_SumSHA256
func SKY_secp256k1_SumSHA256(_b *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(_b))
	__arg1 := secp256k1.SumSHA256(b)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_secp256k1_EntropyPool_Mix256
func SKY_secp256k1_EntropyPool_Mix256(_ep *C.EntropyPool, _in *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ep := (*cipher.EntropyPool)(unsafe.Pointer(_ep))
	in := *(*[]byte)(unsafe.Pointer(_in))
	__arg1 := ep.Mix256(in)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_secp256k1_EntropyPool_Mix
func SKY_secp256k1_EntropyPool_Mix(_ep *C.EntropyPool, _in *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ep := (*cipher.EntropyPool)(unsafe.Pointer(_ep))
	in := *(*[]byte)(unsafe.Pointer(_in))
	__arg1 := ep.Mix(in)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_secp256k1_RandByte
func SKY_secp256k1_RandByte(_n int, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	n := _n
	__arg1 := secp256k1.RandByte(n)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}
