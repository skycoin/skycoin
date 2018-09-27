package main

import (
	"reflect"
	"unsafe"

	secp256k1go2 "github.com/skycoin/skycoin/src/cipher/secp256k1-go/secp256k1-go2"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_secp256k1go_DecompressPoint
func SKY_secp256k1go_DecompressPoint(_X []byte, _off bool, _Y []byte) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	X := *(*[]byte)(unsafe.Pointer(&_X))
	off := _off
	Y := *(*[]byte)(unsafe.Pointer(&_Y))
	secp256k1go2.DecompressPoint(X, off, Y)
	return
}

//export SKY_secp256k1go_RecoverPublicKey
func SKY_secp256k1go_RecoverPublicKey(_sigByte []byte, _h []byte, _recid int, _arg3 *C.GoSlice_, _arg4 *int) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	sigByte := *(*[]byte)(unsafe.Pointer(&_sigByte))
	h := *(*[]byte)(unsafe.Pointer(&_h))
	recid := _recid
	__arg3, __arg4 := secp256k1go2.RecoverPublicKey(sigByte, h, recid)
	copyToGoSlice(reflect.ValueOf(__arg3), _arg3)
	*_arg4 = __arg4
	return
}

//export SKY_secp256k1go_Multiply
func SKY_secp256k1go_Multiply(_xy, _k []byte, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	xy := *(*[]byte)(unsafe.Pointer(&_xy))
	k := *(*[]byte)(unsafe.Pointer(&_k))
	__arg1 := secp256k1go2.Multiply(xy, k)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_secp256k1go_BaseMultiply
func SKY_secp256k1go_BaseMultiply(_k []byte, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	k := *(*[]byte)(unsafe.Pointer(&_k))
	__arg1 := secp256k1go2.BaseMultiply(k)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_secp256k1go_BaseMultiplyAdd
func SKY_secp256k1go_BaseMultiplyAdd(_xy, _k []byte, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	xy := *(*[]byte)(unsafe.Pointer(&_xy))
	k := *(*[]byte)(unsafe.Pointer(&_k))
	__arg1 := secp256k1go2.BaseMultiplyAdd(xy, k)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_secp256k1go_GeneratePublicKey
func SKY_secp256k1go_GeneratePublicKey(_k []byte, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	k := *(*[]byte)(unsafe.Pointer(&_k))
	__arg1 := secp256k1go2.GeneratePublicKey(k)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_secp256k1go_SeckeyIsValid
func SKY_secp256k1go_SeckeyIsValid(_seckey []byte, _arg1 *int) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	seckey := *(*[]byte)(unsafe.Pointer(&_seckey))
	__arg1 := secp256k1go2.SeckeyIsValid(seckey)
	*_arg1 = __arg1
	return
}

//export SKY_secp256k1go_PubkeyIsValid
func SKY_secp256k1go_PubkeyIsValid(_pubkey []byte, _arg1 *int) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	pubkey := *(*[]byte)(unsafe.Pointer(&_pubkey))
	__arg1 := secp256k1go2.PubkeyIsValid(pubkey)
	*_arg1 = __arg1
	return
}
