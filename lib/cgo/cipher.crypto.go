package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"

	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cipher_RandByte
func SKY_cipher_RandByte(_n int, _arg1 *C.GoSlice_) {
	b := cipher.RandByte(_n)
	copyToGoSlice(reflect.ValueOf(b), _arg1)
}

//export SKY_cipher_NewPubKey
func SKY_cipher_NewPubKey(_b []byte, _arg1 *C.cipher__PubKey) uint32 {
	pubkey, err := cipher.NewPubKey(_b)
	errcode := libErrorCode(err)

	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	}

	return errcode
}

//export SKY_cipher_PubKeyFromHex
func SKY_cipher_PubKeyFromHex(_s string, _arg1 *C.cipher__PubKey) uint32 {
	pubkey, err := cipher.PubKeyFromHex(_s)
	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	}
	return errcode
}

//export SKY_cipher_PubKeyFromSecKey
func SKY_cipher_PubKeyFromSecKey(_seckey *C.cipher__SecKey, _arg1 *C.cipher__PubKey) uint32 {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))

	pubkey, err := cipher.PubKeyFromSecKey(*seckey)
	errcode := libErrorCode(err)

	if err == nil {
		copyToBuffer(reflect.ValueOf(pubkey[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
	}

	return errcode
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
func SKY_cipher_PubKey_Hex(_pk *C.cipher__PubKey, _arg1 *C.GoString_) {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))
	s := pk.Hex()
	copyString(s, _arg1)
}

//export SKY_cipher_PubKeyRipemd160
func SKY_cipher_PubKeyRipemd160(_pk *C.cipher__PubKey, _arg0 *C.cipher__Ripemd160) {
	pk := (*cipher.PubKey)(unsafe.Pointer(_pk))
	h := cipher.PubKeyRipemd160(*pk)
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg0), uint(SizeofRipemd160))
}

//export SKY_cipher_NewSecKey
func SKY_cipher_NewSecKey(_b []byte, _arg1 *C.cipher__SecKey) uint32 {
	sk, err := cipher.NewSecKey(_b)
	errcode := libErrorCode(err)

	if err == nil {
		copyToBuffer(reflect.ValueOf(sk[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	}

	return errcode
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
func SKY_cipher_SecKey_Hex(_sk *C.cipher__SecKey, _arg1 *C.GoString_) {
	sk := (*cipher.SecKey)(unsafe.Pointer(_sk))
	s := sk.Hex()
	copyString(s, _arg1)
}

//export SKY_cipher_ECDH
func SKY_cipher_ECDH(_pub *C.cipher__PubKey, _sec *C.cipher__SecKey, _arg2 *C.GoSlice_) uint32 {
	pub := (*cipher.PubKey)(unsafe.Pointer(_pub))
	sec := (*cipher.SecKey)(unsafe.Pointer(_sec))
	b, err := cipher.ECDH(*pub, *sec)
	errcode := libErrorCode(err)
	if err == nil {
		copyToGoSlice(reflect.ValueOf(b), _arg2)
	}
	return errcode
}

//export SKY_cipher_NewSig
func SKY_cipher_NewSig(_b []byte, _arg1 *C.cipher__Sig) uint32 {
	s, err := cipher.NewSig(_b)
	errcode := libErrorCode(err)

	if err == nil {
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSig))
	}

	return errcode
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
func SKY_cipher_Sig_Hex(_s *C.cipher__Sig, _arg1 *C.GoString_) {
	s := (*cipher.Sig)(unsafe.Pointer(_s))
	copyString(s.Hex(), _arg1)
}

//export SKY_cipher_MustSignHash
func SKY_cipher_MustSignHash(_hash *C.cipher__SHA256, _sec *C.cipher__SecKey, _arg2 *C.cipher__Sig) {
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))
	sec := (*cipher.SecKey)(unsafe.Pointer(_sec))
	s := cipher.MustSignHash(*hash, *sec)
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg2), uint(SizeofSig))
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
func SKY_cipher_GenerateKeyPair(_arg0 *C.cipher__PubKey, _arg1 *C.cipher__SecKey) {
	p, s := cipher.GenerateKeyPair()
	copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg0), uint(SizeofPubKey))
	copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
}

//export SKY_cipher_GenerateDeterministicKeyPair
func SKY_cipher_GenerateDeterministicKeyPair(_seed []byte, _arg1 *C.cipher__PubKey, _arg2 *C.cipher__SecKey) uint32 {
	p, s, err := cipher.GenerateDeterministicKeyPair(_seed)
	errcode := libErrorCode(err)

	if err == nil {
		copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg1), uint(SizeofPubKey))
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg2), uint(SizeofSecKey))
	}

	return errcode
}

//export SKY_cipher_DeterministicKeyPairIterator
func SKY_cipher_DeterministicKeyPairIterator(_seed []byte, _arg1 *C.GoSlice_, _arg2 *C.cipher__PubKey, _arg3 *C.cipher__SecKey) uint32 {
	h, p, s, err := cipher.DeterministicKeyPairIterator(_seed)
	errcode := libErrorCode(err)

	if err == nil {
		copyToGoSlice(reflect.ValueOf(h), _arg1)
		copyToBuffer(reflect.ValueOf(p[:]), unsafe.Pointer(_arg2), uint(SizeofPubKey))
		copyToBuffer(reflect.ValueOf(s[:]), unsafe.Pointer(_arg3), uint(SizeofSecKey))
	}

	return errcode
}

//export SKY_cipher_GenerateDeterministicKeyPairs
func SKY_cipher_GenerateDeterministicKeyPairs(_seed []byte, _n int, _arg2 *C.GoSlice_) uint32 {
	sks, err := cipher.GenerateDeterministicKeyPairs(_seed, _n)
	errcode := libErrorCode(err)

	if err == nil {
		copyToGoSlice(reflect.ValueOf(sks), _arg2)
	}

	return errcode
}

//export SKY_cipher_GenerateDeterministicKeyPairsSeed
func SKY_cipher_GenerateDeterministicKeyPairsSeed(_seed []byte, _n int, _arg2 *C.GoSlice_, _arg3 *C.GoSlice_) uint32 {
	h, sks, err := cipher.GenerateDeterministicKeyPairsSeed(_seed, _n)
	errcode := libErrorCode(err)

	if err == nil {
		copyToGoSlice(reflect.ValueOf(h), _arg2)
		copyToGoSlice(reflect.ValueOf(sks), _arg3)
	}

	return errcode
}

//export SKY_cipher_CheckSecKey
func SKY_cipher_CheckSecKey(_seckey *C.cipher__SecKey) uint32 {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))

	err := cipher.CheckSecKey(*seckey)
	return libErrorCode(err)
}

//export SKY_cipher_CheckSecKeyHash
func SKY_cipher_CheckSecKeyHash(_seckey *C.cipher__SecKey, _hash *C.cipher__SHA256) uint32 {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))
	hash := (*cipher.SHA256)(unsafe.Pointer(_hash))

	err := cipher.CheckSecKeyHash(*seckey, *hash)
	return libErrorCode(err)
}
