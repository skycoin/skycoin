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

//export SKY_secp256k1go_Signature_Create
func SKY_secp256k1go_Signature_Create(handle *C.Signature_Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	var sig secp256k1go2.Signature
	*handle = registerSignatureHandle(&sig)
	return
}

//export SKY_secp256k1go_Signature_Get_R
func SKY_secp256k1go_Signature_Get_R(handle C.Signature_Handle, r *C.Number_Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig, ok := lookupSignatureHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	*r = registerNumberHandle(&sig.R)
	return
}

//export SKY_secp256k1go_Signature_Get_S
func SKY_secp256k1go_Signature_Get_S(handle C.Signature_Handle, s *C.Number_Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig, ok := lookupSignatureHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	*s = registerNumberHandle(&sig.S)
	return
}

//export SKY_secp256k1go_Signature_Print
func SKY_secp256k1go_Signature_Print(handle C.Signature_Handle, _lab string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig, ok := lookupSignatureHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	sig.Print(_lab)
	return
}

//export SKY_secp256k1go_Signature_Verify
func SKY_secp256k1go_Signature_Verify(handle C.Signature_Handle, _pubkey *C.secp256k1go__XY, _message C.Number_Handle, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig, ok := lookupSignatureHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	pubkey := (*secp256k1go2.XY)(unsafe.Pointer(_pubkey))
	message, okm := lookupNumberHandle(_message)
	if !okm {
		____error_code = SKY_ERROR
		return
	}
	__arg2 := sig.Verify(pubkey, message)
	*_arg2 = __arg2
	return
}

//export SKY_secp256k1go_Signature_Recover
func SKY_secp256k1go_Signature_Recover(handle C.Signature_Handle, _pubkey *C.secp256k1go__XY, _message C.Number_Handle, _recid int, _arg3 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig, ok := lookupSignatureHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	pubkey := (*secp256k1go2.XY)(unsafe.Pointer(_pubkey))
	m, okm := lookupNumberHandle(_message)
	if !okm {
		____error_code = SKY_ERROR
		return
	}
	recid := _recid
	__arg3 := sig.Recover(pubkey, m, recid)
	*_arg3 = __arg3
	return
}

//export SKY_secp256k1go_Signature_Sign
func SKY_secp256k1go_Signature_Sign(handle C.Signature_Handle, _seckey, _message, _nonce C.Number_Handle, _recid *int, _arg2 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig, ok := lookupSignatureHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	seckey, oks := lookupNumberHandle(_seckey)
	if !oks {
		____error_code = SKY_ERROR
		return
	}
	message, okm := lookupNumberHandle(_message)
	if !okm {
		____error_code = SKY_ERROR
		return
	}
	nonce, okn := lookupNumberHandle(_nonce)
	if !okn {
		____error_code = SKY_ERROR
		return
	}
	recid := _recid
	__arg2 := sig.Sign(seckey, message, nonce, recid)
	*_arg2 = __arg2
	return
}

//export SKY_secp256k1go_Signature_ParseBytes
func SKY_secp256k1go_Signature_ParseBytes(handle C.Signature_Handle, _v []byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig, ok := lookupSignatureHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	v := *(*[]byte)(unsafe.Pointer(&_v))
	sig.ParseBytes(v)
	return
}

//export SKY_secp256k1go_Signature_Bytes
func SKY_secp256k1go_Signature_Bytes(handle C.Signature_Handle, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig, ok := lookupSignatureHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := sig.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
