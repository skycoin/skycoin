package main

import (
	"reflect"
	"unsafe"

	cipher "github.com/skycoin/skycoin/src/cipher"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cipher_PubKeySlice_Len
func SKY_cipher_PubKeySlice_Len(_slice *C.cipher__PubKeySlice) int {
	slice := inplacePubKeySlice(_slice)
	return slice.Len()
}

//export SKY_cipher_PubKeySlice_Less
func SKY_cipher_PubKeySlice_Less(_slice *C.cipher__PubKeySlice, _i, _j int) bool {
	slice := inplacePubKeySlice(_slice)
	return slice.Less(_i, _j)
}

//export SKY_cipher_PubKeySlice_Swap
func SKY_cipher_PubKeySlice_Swap(_slice *C.cipher__PubKeySlice, _i, _j int) uint32 {
	slice := inplacePubKeySlice(_slice)
	slice.Swap(_i, _j)
	return SKY_OK
}

//export SKY_cipher_RandByte
func SKY_cipher_RandByte(_n int, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	n := _n
	__arg1 := cipher.RandByte(n)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_cipher_NewPubKey
func SKY_cipher_NewPubKey(_b []byte, _arg1 *C.cipher__PubKey) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1 := cipher.NewPubKey(b)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	return
}

//export SKY_cipher_MustPubKeyFromHex
func SKY_cipher_MustPubKeyFromHex(_s string, _arg1 *C.cipher__PubKey) (errcode uint32) {
	errcode = 0
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()
	s := _s
	__arg1 := cipher.MustPubKeyFromHex(s)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	return
}

//export SKY_cipher_PubKeyFromHex
func SKY_cipher_PubKeyFromHex(_s string, _arg1 *C.cipher__PubKey) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()
	s := _s
	__arg1, ____return_err := cipher.PubKeyFromHex(s)
	errcode = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	}
	return
}

//export SKY_cipher_PubKeyFromSecKey
func SKY_cipher_PubKeyFromSecKey(_seckey *C.cipher__SecKey, _arg1 *C.cipher__PubKey) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()
	seckey := *(*cipher.SecKey)(unsafe.Pointer(_seckey))
	__arg1 := cipher.PubKeyFromSecKey(seckey)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	return
}

//export SKY_cipher_PubKeyFromSig
func SKY_cipher_PubKeyFromSig(_sig *C.cipher__Sig, _hash *C.cipher__SHA256, _arg2 *C.cipher__PubKey) uint32 {
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	pubkey, err := cipher.PubKeyFromSig(*sig, *hash)

	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg2), uint(SizeofPubKey))

	}
	return errcode
}

//export SKY_cipher_PubKey_Verify
func SKY_cipher_PubKey_Verify(_pk *C.cipher__PubKey) uint32 {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))

	err := pk.Verify()
	errcode := libErrorCode(err)
	return errcode
}

//export SKY_cipher_PubKey_Hex
func SKY_cipher_PubKey_Hex(_pk *C.cipher__PubKey, _arg1 *C.GoString_) uint32 {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))
	s := pk.Hex()
	copyString(s, _arg1)
	return SKY_OK
}

//export SKY_cipher_PubKey_ToAddressHash
func SKY_cipher_PubKey_ToAddressHash(_pk *C.cipher__PubKey, _arg0 *C.cipher__Ripemd160) uint32 {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))
	__arg0 := pk.ToAddressHash()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofRipemd160))
	return SKY_OK
}

//export SKY_cipher_NewSecKey
func SKY_cipher_NewSecKey(_b []byte, _arg1 *C.cipher__SecKey) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1 := cipher.NewSecKey(b)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	return
}

//export SKY_cipher_MustSecKeyFromHex
func SKY_cipher_MustSecKeyFromHex(_s string, _arg1 *C.cipher__SecKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	s := _s
	__arg1 := cipher.MustSecKeyFromHex(s)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	return
}

//export SKY_cipher_SecKeyFromHex
func SKY_cipher_SecKeyFromHex(_s string, _arg1 *C.cipher__SecKey) uint32 {
	sk, err := cipher.SecKeyFromHex(_s)
	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(sk[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	}
	return errcode
}

//export SKY_cipher_SecKey_Verify
func SKY_cipher_SecKey_Verify(_sk *C.cipher__SecKey) uint32 {
	sk := (*cipher.SecKey)(unsafe.Pointer(_sk))
	err := sk.Verify()
	return libErrorCode(err)
}

//export SKY_cipher_SecKey_Hex
func SKY_cipher_SecKey_Hex(_sk *C.cipher__SecKey, _arg1 *C.GoString_) uint32 {
	sk := (*cipher.SecKey)(unsafe.Pointer(_sk))
	s := sk.Hex()
	copyString(s, _arg1)
	return SKY_OK
}

//export SKY_cipher_ECDH
func SKY_cipher_ECDH(_pub *C.cipher__PubKey, _sec *C.cipher__SecKey, _arg2 *C.GoSlice_) uint32 {
	pub := (*cipher.PubKey)(unsafe.Pointer(_pub))
	sec := (*cipher.SecKey)(unsafe.Pointer(_sec))
	b := cipher.ECDH(*pub, *sec)
	copyToGoSlice(reflect.ValueOf(b), _arg2)
	return SKY_OK
}

//export SKY_cipher_NewSig
func SKY_cipher_NewSig(_b []byte, _arg1 *C.cipher__Sig) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1 := cipher.NewSig(b)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofSig))
	return
}

//export SKY_cipher_MustSigFromHex
func SKY_cipher_MustSigFromHex(_s string, _arg1 *C.cipher__Sig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	s := _s
	__arg1 := cipher.MustSigFromHex(s)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofSig))
	return
}

//export SKY_cipher_SigFromHex
func SKY_cipher_SigFromHex(_s string, _arg1 *C.cipher__Sig) uint32 {
	s, err := cipher.SigFromHex(_s)
	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSig))
	}
	return errcode
}

//export SKY_cipher_Sig_Hex
func SKY_cipher_Sig_Hex(_s *C.cipher__Sig, _arg1 *C.GoString_) uint32 {
	s := (*cipher.Sig)(unsafe.Pointer(_s))
	copyString(s.Hex(), _arg1)
	return SKY_OK
}

//export SKY_cipher_SignHash
func SKY_cipher_SignHash(_hash *C.cipher__SHA256, _sec *C.cipher__SecKey, _arg2 *C.cipher__Sig) uint32 {
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sec := (*cipher.SecKey)(unsafe.Pointer(_sec))
	s := cipher.SignHash(*hash, *sec)
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg2), uint(SizeofSig))
	return SKY_OK
}

//export SKY_cipher_ChkSig
func SKY_cipher_ChkSig(_address *C.cipher__Address, _hash *C.cipher__SHA256, _sig *C.cipher__Sig) uint32 {
	address := inplaceAddress(_address)
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))

	err := cipher.ChkSig(*address, *hash, *sig)
	return libErrorCode(err)
}

//export SKY_cipher_VerifySignedHash
func SKY_cipher_VerifySignedHash(_sig *C.cipher__Sig, _hash *C.cipher__SHA256) uint32 {
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))

	err := cipher.VerifySignedHash(*sig, *hash)
	return libErrorCode(err)
}

//export SKY_cipher_VerifySignature
func SKY_cipher_VerifySignature(_pubkey *C.cipher__PubKey, _sig *C.cipher__Sig, _hash *C.cipher__SHA256) uint32 {
	pubkey := (*cipher.PubKey)(unsafe.Pointer(_pubkey))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	err := cipher.VerifySignature(*pubkey, *sig, *hash)
	return libErrorCode(err)
}

//export SKY_cipher_GenerateKeyPair
func SKY_cipher_GenerateKeyPair(_arg0 *C.cipher__PubKey, _arg1 *C.cipher__SecKey) uint32 {
	p, s := cipher.GenerateKeyPair()
	copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg0), uint(SizeofPubKey))
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	return SKY_OK
}

//export SKY_cipher_GenerateDeterministicKeyPair
func SKY_cipher_GenerateDeterministicKeyPair(_seed []byte, _arg1 *C.cipher__PubKey, _arg2 *C.cipher__SecKey) uint32 {
	p, s := cipher.GenerateDeterministicKeyPair(_seed)
	copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg2), uint(SizeofSecKey))
	return SKY_OK
}

//export SKY_cipher_DeterministicKeyPairIterator
func SKY_cipher_DeterministicKeyPairIterator(_seed []byte, _arg1 *C.GoSlice_, _arg2 *C.cipher__PubKey, _arg3 *C.cipher__SecKey) uint32 {
	h, p, s := cipher.DeterministicKeyPairIterator(_seed)

	copyToGoSlice(reflect.ValueOf(h), _arg1)
	copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg2), uint(SizeofPubKey))
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg3), uint(SizeofSecKey))
	return SKY_OK
}

//export SKY_cipher_GenerateDeterministicKeyPairs
func SKY_cipher_GenerateDeterministicKeyPairs(_seed []byte, _n int, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seed := *(*[]byte)(unsafe.Pointer(&_seed))
	n := _n
	__arg2 := cipher.GenerateDeterministicKeyPairs(seed, n)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

//export SKY_cipher_GenerateDeterministicKeyPairsSeed
func SKY_cipher_GenerateDeterministicKeyPairsSeed(_seed []byte, _n int, _arg2 *C.GoSlice_, _arg3 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seed := *(*[]byte)(unsafe.Pointer(&_seed))
	n := _n
	__arg2, __arg3 := cipher.GenerateDeterministicKeyPairsSeed(seed, n)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	copyToGoSlice(reflect.ValueOf(__arg3), _arg3)
	return
}

//export SKY_cipher_TestSecKey
func SKY_cipher_TestSecKey(_seckey *C.cipher__SecKey) uint32 {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))

	err := cipher.TestSecKey(*seckey)
	return libErrorCode(err)
}

//export SKY_cipher_TestSecKeyHash
func SKY_cipher_TestSecKeyHash(_seckey *C.cipher__SecKey, _hash *C.cipher__SHA256) uint32 {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	err := cipher.TestSecKeyHash(*seckey, *hash)
	return libErrorCode(err)
}
