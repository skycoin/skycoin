package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"

	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_cipher_PubKeySlice_Len
func SKY_cipher_PubKeySlice_Len(_slice *C.PubKeySlice) int {
	slice := inplacePubKeySlice(_slice)
	return slice.Len()
}

//export SKY_cipher_PubKeySlice_Less
func SKY_cipher_PubKeySlice_Less(_slice *C.PubKeySlice, _i, _j int) bool {
	slice := inplacePubKeySlice(_slice)
	return slice.Less(_i, _j)
}

//export SKY_cipher_PubKeySlice_Swap
func SKY_cipher_PubKeySlice_Swap(_slice *C.PubKeySlice, _i, _j int) {
	slice := inplacePubKeySlice(_slice)
	slice.Swap(_i, _j)
}

//export SKY_cipher_RandByte
func SKY_cipher_RandByte(_n int, _arg1 *C.GoSlice_) {
	b := cipher.RandByte(_n)
	copyToGoSlice(reflect.ValueOf(b), _arg1)
}

//export SKY_cipher_NewPubKey
func SKY_cipher_NewPubKey(_b []byte, _arg1 *C.PubKey) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()

	pubkey := cipher.NewPubKey(_b)
	copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	return libErrorCode(nil)
}

//export SKY_cipher_PubKeyFromHex
func SKY_cipher_PubKeyFromHex(_s string, _arg1 *C.PubKey) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()

	pubkey, err := cipher.PubKeyFromHex(_s)
	errcode = libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	}
	return errcode
}

//export SKY_cipher_PubKeyFromSecKey
func SKY_cipher_PubKeyFromSecKey(_seckey *C.SecKey, _arg1 *C.PubKey) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()

	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))

	pubkey := cipher.PubKeyFromSecKey(*seckey)

	copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	return libErrorCode(nil)
}

//export SKY_cipher_PubKeyFromSig
func SKY_cipher_PubKeyFromSig(_sig *C.Sig, _hash *C.SHA256, _arg2 *C.PubKey) uint32 {
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
func SKY_cipher_PubKey_Verify(_pk *C.PubKey) uint32 {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))

	err := pk.Verify()
	errcode := libErrorCode(err)
	return errcode
}

//export SKY_cipher_PubKey_Hex
func SKY_cipher_PubKey_Hex(_pk *C.PubKey, _arg1 *C.GoString_) {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))
	s := pk.Hex()
	copyString(s, _arg1)
}

//export SKY_cipher_PubKey_ToAddressHash
func SKY_cipher_PubKey_ToAddressHash(_pk *C.PubKey, _arg0 *C.Ripemd160) {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))
	h := pk.ToAddressHash()
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg0), uint(SizeofRipemd160))
}

//export SKY_cipher_NewSecKey
func SKY_cipher_NewSecKey(_b []byte, _arg1 *C.SecKey) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()

	sk := cipher.NewSecKey(_b)
	copyToBuffer(reflect.ValueOf(sk[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	return SKY_OK
}

//export SKY_cipher_SecKeyFromHex
func SKY_cipher_SecKeyFromHex(_s string, _arg1 *C.SecKey) uint32 {
	sk, err := cipher.SecKeyFromHex(_s)
	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(sk[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	}
	return errcode
}

//export SKY_cipher_SecKey_Verify
func SKY_cipher_SecKey_Verify(_sk *C.SecKey) uint32 {
	sk := (*cipher.SecKey)(unsafe.Pointer(_sk))
	err := sk.Verify()
	return libErrorCode(err)
}

//export SKY_cipher_SecKey_Hex
func SKY_cipher_SecKey_Hex(_sk *C.SecKey, _arg1 *C.GoString_) {
	sk := (*cipher.SecKey)(unsafe.Pointer(_sk))
	s := sk.Hex()
	copyString(s, _arg1)
}

//export SKY_cipher_ECDH
func SKY_cipher_ECDH(_pub *C.PubKey, _sec *C.SecKey, _arg2 *C.GoSlice_) {
	pub := (*cipher.PubKey)(unsafe.Pointer(_pub))
	sec := (*cipher.SecKey)(unsafe.Pointer(_sec))
	b := cipher.ECDH(*pub, *sec)
	copyToGoSlice(reflect.ValueOf(b), _arg2)
}

//export SKY_cipher_NewSig
func SKY_cipher_NewSig(_b []byte, _arg1 *C.Sig) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()

	s := cipher.NewSig(_b)
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSig))

	return SKY_OK
}

//export SKY_cipher_SigFromHex
func SKY_cipher_SigFromHex(_s string, _arg1 *C.Sig) uint32 {
	s, err := cipher.SigFromHex(_s)
	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSig))
	}
	return errcode
}

//export SKY_cipher_Sig_Hex
func SKY_cipher_Sig_Hex(_s *C.Sig, _arg1 *C.GoString_) {
	s := (*cipher.Sig)(unsafe.Pointer(_s))
	copyString(s.Hex(), _arg1)
}

//export SKY_cipher_SignHash
func SKY_cipher_SignHash(_hash *C.SHA256, _sec *C.SecKey, _arg2 *C.Sig) {
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sec := (*cipher.SecKey)(unsafe.Pointer(_sec))
	s := cipher.SignHash(*hash, *sec)
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg2), uint(SizeofSig))
}

//export SKY_cipher_ChkSig
func SKY_cipher_ChkSig(_address *C.Address, _hash *C.SHA256, _sig *C.Sig) uint32 {
	address := inplaceAddress(_address)
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))

	err := cipher.ChkSig(*address, *hash, *sig)
	return libErrorCode(err)
}

//export SKY_cipher_VerifySignedHash
func SKY_cipher_VerifySignedHash(_sig *C.Sig, _hash *C.SHA256) uint32 {
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))

	err := cipher.VerifySignedHash(*sig, *hash)
	return libErrorCode(err)
}

//export SKY_cipher_VerifySignature
func SKY_cipher_VerifySignature(_pubkey *C.PubKey, _sig *C.Sig, _hash *C.SHA256) uint32 {
	pubkey := (*cipher.PubKey)(unsafe.Pointer(_pubkey))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	err := cipher.VerifySignature(*pubkey, *sig, *hash)
	return libErrorCode(err)
}

//export SKY_cipher_GenerateKeyPair
func SKY_cipher_GenerateKeyPair(_arg0 *C.PubKey, _arg1 *C.SecKey) {
	p, s := cipher.GenerateKeyPair()
	copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg0), uint(SizeofPubKey))
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
}

//export SKY_cipher_GenerateDeterministicKeyPair
func SKY_cipher_GenerateDeterministicKeyPair(_seed []byte, _arg1 *C.PubKey, _arg2 *C.SecKey) {
	p, s := cipher.GenerateDeterministicKeyPair(_seed)
	copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg2), uint(SizeofSecKey))
}

//export SKY_cipher_DeterministicKeyPairIterator
func SKY_cipher_DeterministicKeyPairIterator(_seed []byte, _arg1 *C.GoSlice_, _arg2 *C.PubKey, _arg3 *C.SecKey) {
	h, p, s := cipher.DeterministicKeyPairIterator(_seed)

	copyToGoSlice(reflect.ValueOf(h), _arg1)
	copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg2), uint(SizeofPubKey))
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg3), uint(SizeofSecKey))
}

//export SKY_cipher_GenerateDeterministicKeyPairs
func SKY_cipher_GenerateDeterministicKeyPairs(_seed []byte, _n int, _arg2 *C.GoSlice_) {
	sks := cipher.GenerateDeterministicKeyPairs(_seed, _n)
	copyToGoSlice(reflect.ValueOf(sks), _arg2)
}

//export SKY_cipher_GenerateDeterministicKeyPairsSeed
func SKY_cipher_GenerateDeterministicKeyPairsSeed(_seed []byte, _n int, _arg2 *C.GoSlice_, _arg3 *C.GoSlice_) {
	h, sks := cipher.GenerateDeterministicKeyPairsSeed(_seed, _n)
	copyToGoSlice(reflect.ValueOf(h), _arg2)
	copyToGoSlice(reflect.ValueOf(sks), _arg3)
}

//export SKY_cipher_TestSecKey
func SKY_cipher_TestSecKey(_seckey *C.SecKey) uint32 {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))

	err := cipher.TestSecKey(*seckey)
	return libErrorCode(err)
}

//export SKY_cipher_TestSecKeyHash
func SKY_cipher_TestSecKeyHash(_seckey *C.SecKey, _hash *C.SHA256) uint32 {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	err := cipher.TestSecKeyHash(*seckey, *hash)
	return libErrorCode(err)
}
