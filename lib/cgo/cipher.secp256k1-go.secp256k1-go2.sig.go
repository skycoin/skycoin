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

//export SKY_secp256k1go_Signature_Print
func SKY_secp256k1go_Signature_Print(_sig *C.Signature, _lab string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go2.Signature)(unsafe.Pointer(_sig))
	lab := _lab
	sig.Print(lab)
	return
}

//export SKY_secp256k1go_Signature_Verify
func SKY_secp256k1go_Signature_Verify(_sig *C.Signature, _pubkey *C.secp256k1go__XY, _message *C.Number, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go2.Signature)(unsafe.Pointer(_sig))
	pubkey := (*secp256k1go2.XY)(unsafe.Pointer(_pubkey))
	message := (*secp256k1go2.Number)(unsafe.Pointer(_message))
	__arg2 := sig.Verify(pubkey, message)
	*_arg2 = __arg2
	return
}

//export SKY_secp256k1go_Signature_Recover
func SKY_secp256k1go_Signature_Recover(_sig *C.Signature, _pubkey *C.secp256k1go__XY, _m *C.Number, _recid int, _arg3 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go2.Signature)(unsafe.Pointer(_sig))
	pubkey := (*secp256k1go2.XY)(unsafe.Pointer(_pubkey))
	m := (*secp256k1go2.Number)(unsafe.Pointer(_m))
	recid := _recid
	__arg3 := sig.Recover(pubkey, m, recid)
	*_arg3 = __arg3
	return
}

//export SKY_secp256k1go_Signature_Sign
func SKY_secp256k1go_Signature_Sign(_sig *C.Signature, _seckey, _message, _nonce *C.Number, _recid *int, _arg2 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go2.Signature)(unsafe.Pointer(_sig))
	seckey := (*secp256k1go2.Number)(unsafe.Pointer(_seckey))
	message := (*secp256k1go2.Number)(unsafe.Pointer(_message))
	nonce := (*secp256k1go2.Number)(unsafe.Pointer(_nonce))
	recid := _recid
	__arg2 := sig.Sign(seckey, message, nonce, recid)
	*_arg2 = __arg2
	return
}

//export SKY_secp256k1go_Signature_ParseBytes
func SKY_secp256k1go_Signature_ParseBytes(_sig *C.Signature, _v []byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go2.Signature)(unsafe.Pointer(_sig))
	v := *(*[]byte)(unsafe.Pointer(&_v))
	sig.ParseBytes(v)
	return
}

//export SKY_secp256k1go_Signature_Bytes
func SKY_secp256k1go_Signature_Bytes(_sig *C.Signature, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go2.Signature)(unsafe.Pointer(_sig))
	__arg0 := sig.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
