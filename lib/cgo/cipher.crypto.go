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

//export SKY_cipher_RandByte
func SKY_cipher_RandByte(_n int, _arg1 *C.GoSlice_) (____error_code uint32) {
	b := cipher.RandByte(_n)
	copyToGoSlice(reflect.ValueOf(b), _arg1)
	return
}

//export SKY_cipher_NewPubKey
func SKY_cipher_NewPubKey(_b []byte, _arg1 *C.cipher__PubKey) (____error_code uint32) {
	pubkey, err := cipher.NewPubKey(_b)
	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	}
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_PubKeyFromHex
func SKY_cipher_PubKeyFromHex(_s string, _arg1 *C.cipher__PubKey) (____error_code uint32) {
	pubkey, err := cipher.PubKeyFromHex(_s)
	____error_code = libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	}
	return
}

//export SKY_cipher_PubKeyFromSecKey
func SKY_cipher_PubKeyFromSecKey(_seckey *C.cipher__SecKey, _arg1 *C.cipher__PubKey) (____error_code uint32) {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))

	pubkey, err := cipher.PubKeyFromSecKey(*seckey)
	____error_code = libErrorCode(err)

	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	}

	return
}

//export SKY_cipher_PubKeyFromSig
func SKY_cipher_PubKeyFromSig(_sig *C.cipher__Sig, _hash *C.cipher__SHA256, _arg2 *C.cipher__PubKey) (____error_code uint32) {
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	pubkey, err := cipher.PubKeyFromSig(*sig, *hash)

	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg2), uint(SizeofPubKey))

	}
	____error_code = errcode
	return
}

//export SKY_cipher_PubKey_Verify
func SKY_cipher_PubKey_Verify(_pk *C.cipher__PubKey) (____error_code uint32) {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))

	err := pk.Verify()
	errcode := libErrorCode(err)
	____error_code = errcode
	return
}

//export SKY_cipher_PubKey_Hex
func SKY_cipher_PubKey_Hex(_pk *C.cipher__PubKey, _arg1 *C.GoString_) (____error_code uint32) {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))
	s := pk.Hex()
	copyString(s, _arg1)
	return SKY_OK
}

//export SKY_cipher_PubKeyRipemd160
func SKY_cipher_PubKeyRipemd160(_pk *C.cipher__PubKey, _arg0 *C.cipher__Ripemd160) (____error_code uint32) {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))
	h := cipher.PubKeyRipemd160(*pk)
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg0), uint(SizeofRipemd160))
	return
}

//export SKY_cipher_NewSecKey
func SKY_cipher_NewSecKey(_b []byte, _arg1 *C.cipher__SecKey) (____error_code uint32) {
	sk, err := cipher.NewSecKey(_b)
	if err == nil {
		copyToBuffer(reflect.ValueOf(sk[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	}

	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_SecKeyFromHex
func SKY_cipher_SecKeyFromHex(_s string, _arg1 *C.cipher__SecKey) (____error_code uint32) {
	sk, err := cipher.SecKeyFromHex(_s)
	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(sk[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	}
	____error_code = errcode
	return
}

//export SKY_cipher_SecKey_Verify
func SKY_cipher_SecKey_Verify(_sk *C.cipher__SecKey) (____error_code uint32) {
	sk := (*cipher.SecKey)(unsafe.Pointer(_sk))
	err := sk.Verify()
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_SecKey_Hex
func SKY_cipher_SecKey_Hex(_sk *C.cipher__SecKey, _arg1 *C.GoString_) (____error_code uint32) {
	sk := (*cipher.SecKey)(unsafe.Pointer(_sk))
	s := sk.Hex()
	copyString(s, _arg1)
	return
}

//export SKY_cipher_ECDH
func SKY_cipher_ECDH(_pub *C.cipher__PubKey, _sec *C.cipher__SecKey, _arg2 *C.GoSlice_) (____error_code uint32) {
	pub := (*cipher.PubKey)(unsafe.Pointer(_pub))
	sec := (*cipher.SecKey)(unsafe.Pointer(_sec))
	b, err := cipher.ECDH(*pub, *sec)
	____error_code = libErrorCode(err)
	if err == nil {
		copyToGoSlice(reflect.ValueOf(b), _arg2)
	}
	return
}

//export SKY_cipher_NewSig
func SKY_cipher_NewSig(_b []byte, _arg1 *C.cipher__Sig) (____error_code uint32) {
	s, err := cipher.NewSig(_b)
	if err == nil {
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSig))
	}
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_SigFromHex
func SKY_cipher_SigFromHex(_s string, _arg1 *C.cipher__Sig) (____error_code uint32) {
	s, err := cipher.SigFromHex(_s)
	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSig))
	}
	____error_code = errcode
	return
}

//export SKY_cipher_Sig_Hex
func SKY_cipher_Sig_Hex(_s *C.cipher__Sig, _arg1 *C.GoString_) (____error_code uint32) {
	s := (*cipher.Sig)(unsafe.Pointer(_s))
	copyString(s.Hex(), _arg1)
	return
}

//export SKY_cipher_SignHash
func SKY_cipher_SignHash(_hash *C.cipher__SHA256, _sec *C.cipher__SecKey, _arg2 *C.cipher__Sig) (____error_code uint32) {
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sec := (*cipher.SecKey)(unsafe.Pointer(_sec))
	s, err := cipher.SignHash(*hash, *sec)
	____error_code = libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg2), uint(SizeofSig))
	}
	return
}

//export SKY_cipher_VerifyAddressSignedHash
func SKY_cipher_VerifyAddressSignedHash(_address *C.cipher__Address, _sig *C.cipher__Sig, _hash *C.cipher__SHA256) (____error_code uint32) {
	address := inplaceAddress(_address)
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))

	err := cipher.VerifyAddressSignedHash(*address, *sig, *hash)
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_VerifySignedHash
func SKY_cipher_VerifySignedHash(_sig *C.cipher__Sig, _hash *C.cipher__SHA256) (____error_code uint32) {
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))

	err := cipher.VerifySignedHash(*sig, *hash)
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_VerifyPubKeySignedHash
func SKY_cipher_VerifyPubKeySignedHash(_pubkey *C.cipher__PubKey, _sig *C.cipher__Sig, _hash *C.cipher__SHA256) (____error_code uint32) {
	pubkey := (*cipher.PubKey)(unsafe.Pointer(_pubkey))
	sig := (*cipher.Sig)(unsafe.Pointer(_sig))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	err := cipher.VerifyPubKeySignedHash(*pubkey, *sig, *hash)
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_GenerateKeyPair
func SKY_cipher_GenerateKeyPair(_arg0 *C.cipher__PubKey, _arg1 *C.cipher__SecKey) (____error_code uint32) {
	p, s := cipher.GenerateKeyPair()
	copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg0), uint(SizeofPubKey))
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	return
}

//export SKY_cipher_GenerateDeterministicKeyPair
func SKY_cipher_GenerateDeterministicKeyPair(_seed []byte, _arg1 *C.cipher__PubKey, _arg2 *C.cipher__SecKey) (____error_code uint32) {
	p, s, err := cipher.GenerateDeterministicKeyPair(_seed)
	if err == nil {
		copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg2), uint(SizeofSecKey))
	}

	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_DeterministicKeyPairIterator
func SKY_cipher_DeterministicKeyPairIterator(_seed []byte, _arg1 *C.GoSlice_, _arg2 *C.cipher__PubKey, _arg3 *C.cipher__SecKey) (____error_code uint32) {
	h, p, s, err := cipher.DeterministicKeyPairIterator(_seed)
	____error_code = libErrorCode(err)

	if err == nil {
		copyToGoSlice(reflect.ValueOf(h), _arg1)
		copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg2), uint(SizeofPubKey))
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg3), uint(SizeofSecKey))
	}

	return
}

//export SKY_cipher_GenerateDeterministicKeyPairs
func SKY_cipher_GenerateDeterministicKeyPairs(_seed []byte, _n int, _arg2 *C.GoSlice_) (____error_code uint32) {
	sks, err := cipher.GenerateDeterministicKeyPairs(_seed, _n)
	____error_code = libErrorCode(err)
	if err == nil {
		copyToGoSlice(reflect.ValueOf(sks), _arg2)
	}

	return
}

//export SKY_cipher_GenerateDeterministicKeyPairsSeed
func SKY_cipher_GenerateDeterministicKeyPairsSeed(_seed []byte, _n int, _arg2 *C.GoSlice_, _arg3 *C.GoSlice_) (____error_code uint32) {
	h, sks, err := cipher.GenerateDeterministicKeyPairsSeed(_seed, _n)
	if err == nil {
		copyToGoSlice(reflect.ValueOf(h), _arg2)
		copyToGoSlice(reflect.ValueOf(sks), _arg3)
	}

	return
}

//export SKY_cipher_CheckSecKey
func SKY_cipher_CheckSecKey(_seckey *C.cipher__SecKey) (____error_code uint32) {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))

	err := cipher.CheckSecKey(*seckey)
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_CheckSecKeyHash
func SKY_cipher_CheckSecKeyHash(_seckey *C.cipher__SecKey, _hash *C.cipher__SHA256) (____error_code uint32) {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	err := cipher.CheckSecKeyHash(*seckey, *hash)
	____error_code = libErrorCode(err)
	return
}
