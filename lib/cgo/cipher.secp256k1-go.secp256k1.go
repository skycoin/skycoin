package main

import (
	secp256k1go "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_secp256k1_GenerateKeyPair
func SKY_secp256k1_GenerateKeyPair(_arg0 *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0, __arg1 := secp256k1go.GenerateKeyPair()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_secp256k1_PubkeyFromSeckey
func SKY_secp256k1_PubkeyFromSeckey(_seckey *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seckey := *(*[]byte)(unsafe.Pointer(_seckey))
	__arg1 := secp256k1go.PubkeyFromSeckey(seckey)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_secp256k1_UncompressPubkey
func SKY_secp256k1_UncompressPubkey(_pubkey *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pubkey := *(*[]byte)(unsafe.Pointer(_pubkey))
	__arg1 := secp256k1go.UncompressPubkey(pubkey)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_secp256k1_UncompressedPubkeyFromSeckey
func SKY_secp256k1_UncompressedPubkeyFromSeckey(_seckey *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seckey := *(*[]byte)(unsafe.Pointer(_seckey))
	__arg1 := secp256k1go.UncompressedPubkeyFromSeckey(seckey)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_secp256k1_Secp256k1Hash
func SKY_secp256k1_Secp256k1Hash(_hash *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hash := *(*[]byte)(unsafe.Pointer(_hash))
	__arg1 := secp256k1go.Secp256k1Hash(hash)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_secp256k1_GenerateDeterministicKeyPair
func SKY_secp256k1_GenerateDeterministicKeyPair(_seed *C.GoSlice_, _arg1 *C.GoSlice_, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seed := *(*[]byte)(unsafe.Pointer(_seed))
	__arg1, __arg2 := secp256k1go.GenerateDeterministicKeyPair(seed)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

// export SKY_secp256k1_DeterministicKeyPairIterator
func SKY_secp256k1_DeterministicKeyPairIterator(_seedIn *C.GoSlice_, _arg1 *C.GoSlice_, _arg2 *C.GoSlice_, _arg3 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seedIn := *(*[]byte)(unsafe.Pointer(_seedIn))
	__arg1, __arg2, __arg3 := secp256k1go.DeterministicKeyPairIterator(seedIn)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	copyToGoSlice(reflect.ValueOf(__arg3), _arg3)
	return
}

// export SKY_secp256k1_Sign
func SKY_secp256k1_Sign(_msg *C.GoSlice_, _seckey *C.GoSlice_, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := *(*[]byte)(unsafe.Pointer(_msg))
	seckey := *(*[]byte)(unsafe.Pointer(_seckey))
	__arg2 := secp256k1go.Sign(msg, seckey)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

// export SKY_secp256k1_SignDeterministic
func SKY_secp256k1_SignDeterministic(_msg *C.GoSlice_, _seckey *C.GoSlice_, _nonceSeed *C.GoSlice_, _arg3 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := *(*[]byte)(unsafe.Pointer(_msg))
	seckey := *(*[]byte)(unsafe.Pointer(_seckey))
	nonceSeed := *(*[]byte)(unsafe.Pointer(_nonceSeed))
	__arg3 := secp256k1go.SignDeterministic(msg, seckey, nonceSeed)
	copyToGoSlice(reflect.ValueOf(__arg3), _arg3)
	return
}

// export SKY_secp256k1_VerifySeckey
func SKY_secp256k1_VerifySeckey(_seckey *C.GoSlice_, _arg1 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seckey := *(*[]byte)(unsafe.Pointer(_seckey))
	__arg1 := secp256k1go.VerifySeckey(seckey)
	*_arg1 = __arg1
	return
}

// export SKY_secp256k1_VerifyPubkey
func SKY_secp256k1_VerifyPubkey(_pubkey *C.GoSlice_, _arg1 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pubkey := *(*[]byte)(unsafe.Pointer(_pubkey))
	__arg1 := secp256k1go.VerifyPubkey(pubkey)
	*_arg1 = __arg1
	return
}

// export SKY_secp256k1_VerifySignatureValidity
func SKY_secp256k1_VerifySignatureValidity(_sig *C.GoSlice_, _arg1 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := *(*[]byte)(unsafe.Pointer(_sig))
	__arg1 := secp256k1go.VerifySignatureValidity(sig)
	*_arg1 = __arg1
	return
}

// export SKY_secp256k1_VerifySignature
func SKY_secp256k1_VerifySignature(_msg *C.GoSlice_, _sig *C.GoSlice_, _pubkey1 *C.GoSlice_, _arg3 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := *(*[]byte)(unsafe.Pointer(_msg))
	sig := *(*[]byte)(unsafe.Pointer(_sig))
	pubkey1 := *(*[]byte)(unsafe.Pointer(_pubkey1))
	__arg3 := secp256k1go.VerifySignature(msg, sig, pubkey1)
	*_arg3 = __arg3
	return
}

// export SKY_secp256k1_SignatureErrorString
func SKY_secp256k1_SignatureErrorString(_msg *C.GoSlice_, _sig *C.GoSlice_, _pubkey1 *C.GoSlice_, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := *(*[]byte)(unsafe.Pointer(_msg))
	sig := *(*[]byte)(unsafe.Pointer(_sig))
	pubkey1 := *(*[]byte)(unsafe.Pointer(_pubkey1))
	__arg3 := secp256k1go.SignatureErrorString(msg, sig, pubkey1)
	copyString(__arg3, _arg3)
	return
}

// export SKY_secp256k1_RecoverPubkey
func SKY_secp256k1_RecoverPubkey(_msg *C.GoSlice_, _sig *C.GoSlice_, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := *(*[]byte)(unsafe.Pointer(_msg))
	sig := *(*[]byte)(unsafe.Pointer(_sig))
	__arg2 := secp256k1go.RecoverPubkey(msg, sig)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

// export SKY_secp256k1_ECDH
func SKY_secp256k1_ECDH(_pub *C.GoSlice_, _sec *C.GoSlice_, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pub := *(*[]byte)(unsafe.Pointer(_pub))
	sec := *(*[]byte)(unsafe.Pointer(_sec))
	__arg2 := secp256k1go.ECDH(pub, sec)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}
