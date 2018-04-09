package main

import (
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_secp256k1go_Signature_Print
func SKY_secp256k1go_Signature_Print(_sig *C.Signature, _lab string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*Signature)(unsafe.Pointer(_sig))
	lab := _lab
	sig.Print(lab)
	return
}

// export SKY_secp256k1go_Signature_Verify
func SKY_secp256k1go_Signature_Verify(_sig *C.Signature, _pubkey *C.XY, _message *C.Number, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*Signature)(unsafe.Pointer(_sig))
	pubkey := (*XY)(unsafe.Pointer(_pubkey))
	message := (*Number)(unsafe.Pointer(_message))
	__arg2 := sig.Verify(pubkey, message)
	*_arg2 = __arg2
	return
}

// export SKY_secp256k1go_Signature_Recover
func SKY_secp256k1go_Signature_Recover(_sig *C.Signature, _pubkey *C.XY, _m *C.Number, _recid int, _arg3 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*Signature)(unsafe.Pointer(_sig))
	pubkey := (*XY)(unsafe.Pointer(_pubkey))
	m := (*Number)(unsafe.Pointer(_m))
	recid := _recid
	__arg3 := sig.Recover(pubkey, m, recid)
	*_arg3 = __arg3
	return
}

// export SKY_secp256k1go_Signature_Sign
func SKY_secp256k1go_Signature_Sign(_sig *C.Signature, _seckey, _message, _nonce *C.Number, _recid *int, _arg2 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*Signature)(unsafe.Pointer(_sig))
	seckey := (*Number)(unsafe.Pointer(_seckey))
	message := (*Number)(unsafe.Pointer(_message))
	nonce := (*Number)(unsafe.Pointer(_nonce))
	recid := _recid
	__arg2 := sig.Sign(seckey, message, nonce, recid)
	*_arg2 = __arg2
	return
}

// export SKY_secp256k1go_Signature_ParseBytes
func SKY_secp256k1go_Signature_ParseBytes(_sig *C.Signature, _v *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*Signature)(unsafe.Pointer(_sig))
	v := *(*[]byte)(unsafe.Pointer(_v))
	sig.ParseBytes(v)
	return
}

// export SKY_secp256k1go_Signature_Bytes
func SKY_secp256k1go_Signature_Bytes(_sig *C.Signature, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*Signature)(unsafe.Pointer(_sig))
	__arg0 := sig.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
