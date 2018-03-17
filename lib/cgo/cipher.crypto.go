package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"

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
	if _n > int(_arg1.cap) {
		_arg1.len = _arg1.cap - C.GoInt_(_n)
		return
	}
	b := cipher.RandByte(_n)
	_arg1.len = C.GoInt_(len(b))
	arg1 := (*[]byte)(unsafe.Pointer(_arg1))
	copy(*arg1, b)
}

//export SKY_cipher_NewPubKey
func SKY_cipher_NewPubKey(_b []byte, _arg1 *C.PubKey) (retVal uint32) {
	defer func() {
		if err := recover(); err != nil {
			// TODO: Fix to be like retVal = libErrorCode(err)
			retVal = ERR_UNKNOWN
		}
	}()

	arg1 := inplacePubKey(_arg1)
	pubkey := cipher.NewPubKey(_b)
	copy(arg1[:], pubkey[:])
	return libErrorCode(nil)
}

//export SKY_cipher_PubKeyFromHex
func SKY_cipher_PubKeyFromHex(_s string, _arg1 *C.PubKey) uint32 {
	pubkey, err := cipher.PubKeyFromHex(_s)
	errcode := libErrorCode(err)
	if err != nil {
		arg1 := inplacePubKey(_arg1)
		copy(arg1[:], pubkey[:])
	}
	return errcode
}

//export SKY_cipher_PubKeyFromSecKey
func SKY_cipher_PubKeyFromSecKey(_seckey *C.SecKey, _arg1 *C.PubKey) (retVal uint32) {
	defer func() {
		if err := recover(); err != nil {
			// TODO: Fix to be like retVal = libErrorCode(err)
			retVal = ERR_UNKNOWN
		}
	}()

	seckey := inplaceSecKey(_seckey)
	pubkey := cipher.PubKeyFromSecKey(*seckey)
	arg1 := inplacePubKey(_arg1)
	copy(arg1[:], pubkey[:])
	return libErrorCode(nil)
}

//export SKY_cipher_PubKeyFromSig
func SKY_cipher_PubKeyFromSig(_sig *C.Sig, _hash *C.SHA256, _arg2 *C.PubKey) uint32 {
	sig := inplaceSig(_sig)
	hash := inplaceSHA256(_hash)
	arg2 := inplacePubKey(_arg2)
	pubkey, err := cipher.PubKeyFromSig(*sig, *hash)
	errcode := libErrorCode(err)
	if err != nil {
		copy(arg2[:], pubkey[:])
	}
	return errcode;
}

//export SKY_cipher_PubKey_Verify
func SKY_cipher_PubKey_Verify(_pk *C.PubKey) uint32 {
	pk := inplacePubKey(_pk)
	err := pk.Verify()
	errcode := libErrorCode(err)
	return errcode
}

//export SKY_cipher_PubKey_Hex
func SKY_cipher_PubKey_Hex(_pk *C.PubKey) string {
	pk := inplacePubKey(_pk)
	return pk.Hex()
}

//export SKY_cipher_PubKey_ToAddressHash
func SKY_cipher_PubKey_ToAddressHash(_pk *C.PubKey, _arg0 *C.Ripemd160) {
	pk := inplacePubKey(_pk)
	arg0 := inplaceRipemd160(_arg0)
	h := pk.ToAddressHash()
	copy(arg0[:], h[:])
}

//export SKY_cipher_NewSecKey
func SKY_cipher_NewSecKey(_b []byte, _arg1 *C.SecKey) {
	arg1 := inplaceSecKey(_arg1)
	sk := cipher.NewSecKey(_b)
	copy(arg1[:], sk[:])
}

//export SKY_cipher_SecKeyFromHex
func SKY_cipher_SecKeyFromHex(_s string, _arg1 *C.SecKey) uint32 {
	arg1 := inplaceSecKey(_arg1)
	sk, err := cipher.SecKeyFromHex(_s)
	errcode := libErrorCode(err)
	if err != nil {
		copy(arg1[:], sk[:])
	}
	return errcode
}

//export SKY_cipher_SecKey_Verify
func SKY_cipher_SecKey_Verify(_sk *C.SecKey) uint32 {
	sk := inplaceSecKey(_sk)
	err := sk.Verify()
	return libErrorCode(err)
}

//export SKY_cipher_SecKey_Hex
func SKY_cipher_SecKey_Hex(_sk *C.SecKey) string {
	sk := inplaceSecKey(_sk)
	return sk.Hex()
}

//export SKY_cipher_ECDH
func SKY_cipher_ECDH(_pub *C.PubKey, _sec *C.SecKey, _arg2 *C.GoSlice_) {
	pub := inplacePubKey(_pub)
	sec := inplaceSecKey(_sec)
	b := cipher.ECDH(*pub, *sec)
	blen := len(b)
	if int(_arg2.cap) < blen {
		_arg2.len = _arg2.cap - C.GoInt_(blen)
		return
	}
	_arg2.len = C.GoInt_(blen)
	arg2 := inplaceByteArray(unsafe.Pointer(_arg2.data), int(_arg2.len))
	copy(*arg2, b)
}

//export SKY_cipher_NewSig
func SKY_cipher_NewSig(_b []byte, _arg1 *C.Sig) {
	arg1 := inplaceSig(_arg1)
	s := cipher.NewSig(_b)
	copy(arg1[:], s[:])
}

//export SKY_cipher_SigFromHex
func SKY_cipher_SigFromHex(_s string, _arg1 *C.Sig) uint32 {
	arg1 := inplaceSig(_arg1)
	s, err := cipher.SigFromHex(_s)
	errcode := libErrorCode(err)
	if err != nil {
		copy(arg1[:], s[:])
	}
	return errcode
}

//export SKY_cipher_Sig_Hex
func SKY_cipher_Sig_Hex(_s *C.Sig) string {
	s := inplaceSig(_s)
	return s.Hex()
}

//export SKY_cipher_SignHash
func SKY_cipher_SignHash(_hash *C.SHA256, _sec *C.SecKey, _arg2 *C.Sig) {
	hash := inplaceSHA256(_hash)
	sec := inplaceSecKey(_sec)
	arg2 := inplaceSig(_arg2)
	s := cipher.SignHash(*hash, *sec)
	copy(arg2[:], s[:])
}

//export SKY_cipher_ChkSig
func SKY_cipher_ChkSig(_address *C.Address, _hash *C.SHA256, _sig *C.Sig) uint32 {
	address := inplaceAddress(_address)
	hash := inplaceSHA256(_hash)
	sig := inplaceSig(_sig)
	err := cipher.ChkSig(*address, *hash, *sig)
	return libErrorCode(err)
}

//export SKY_cipher_VerifySignedHash
func SKY_cipher_VerifySignedHash(_sig *C.Sig, _hash *C.SHA256) uint32 {
	hash := inplaceSHA256(_hash)
	sig := inplaceSig(_sig)
	err := cipher.VerifySignedHash(*sig, *hash)
	return libErrorCode(err)
}

//export SKY_cipher_VerifySignature
func SKY_cipher_VerifySignature(_pubkey *C.PubKey, _sig *C.Sig, _hash *C.SHA256) uint32 {
	pubkey := inplacePubKey(_pubkey)
	sig := inplaceSig(_sig)
	hash := inplaceSHA256(_hash)
	err := cipher.VerifySignature(*pubkey, *sig, *hash)
	return libErrorCode(err)
}

//export SKY_cipher_GenerateKeyPair
func SKY_cipher_GenerateKeyPair(_arg0 *C.PubKey, _arg1 *C.SecKey) {
	arg0 := inplacePubKey(_arg0)
	arg1 := inplaceSecKey(_arg1)
	p, s := cipher.GenerateKeyPair()
	copy(arg0[:], p[:])
	copy(arg1[:], s[:])
}

//export SKY_cipher_GenerateDeterministicKeyPair
func SKY_cipher_GenerateDeterministicKeyPair(_seed []byte, _arg1 *C.PubKey, _arg2 *C.SecKey) {
	arg1 := inplacePubKey(_arg1)
	arg2 := inplaceSecKey(_arg2)
	p, s := cipher.GenerateDeterministicKeyPair(_seed)
	copy(arg1[:], p[:])
	copy(arg2[:], s[:])
}

//export SKY_cipher_DeterministicKeyPairIterator
func SKY_cipher_DeterministicKeyPairIterator(_seed []byte, _arg1 *C.GoSlice_, _arg2 *C.PubKey, _arg3 *C.SecKey) {
	arg1 := (*[]byte)(unsafe.Pointer(_arg1))
	arg2 := inplacePubKey(_arg2)
	arg3 := inplaceSecKey(_arg3)
	h, p, s := cipher.DeterministicKeyPairIterator(_seed)
	hlen := len(h)
	if hlen > int(_arg1.cap) {
		_arg1.len = _arg1.cap - C.GoInt_(hlen)
		return
	}
	copy(*arg1, h[:])
	copy(arg2[:], p[:])
	copy(arg3[:], s[:])
}

//export SKY_cipher_GenerateDeterministicKeyPairs
func SKY_cipher_GenerateDeterministicKeyPairs(_seed []byte, _n int, _arg2 *C.GoSlice_) {
	sks := cipher.GenerateDeterministicKeyPairs(_seed, _n)
	skslen := len(sks)
	if skslen > int(_arg2.cap) {
		_arg2.len = _arg2.cap - C.GoInt_(skslen)
		return
	}
	_arg2.len = C.GoInt_(skslen)
	arg2 := (*[]cipher.SecKey)(unsafe.Pointer(_arg2))
	copy(*arg2, sks)
}

//export SKY_cipher_GenerateDeterministicKeyPairsSeed
func SKY_cipher_GenerateDeterministicKeyPairsSeed(_seed []byte, _n int, _arg2 *C.GoSlice_, _arg3 *C.GoSlice_) {
	h, sks := cipher.GenerateDeterministicKeyPairsSeed(_seed, _n)
	skslen, hlen := len(sks), len(h)
	nospace := false
	if skslen > int(_arg3.cap) {
		_arg3.len = _arg3.cap - C.GoInt_(skslen)
		nospace = true
	}
	if hlen > int(_arg2.cap) {
		_arg2.len = _arg2.cap - C.GoInt_(hlen)
		nospace = true
	}
	if nospace {
		return
	}
	_arg2.len = C.GoInt_(hlen)
	_arg3.len = C.GoInt_(skslen)
	arg2 := (*[]byte)(unsafe.Pointer(_arg2))
	arg3 := (*[]cipher.SecKey)(unsafe.Pointer(_arg3))
	copy(*arg2, h)
	copy(*arg3, sks)
}

//export SKY_cipher_TestSecKey
func SKY_cipher_TestSecKey(_seckey *C.SecKey) uint32 {
	seckey := inplaceSecKey(_seckey)
	err := cipher.TestSecKey(*seckey)
	return libErrorCode(err)
}

//export SKY_cipher_TestSecKeyHash
func SKY_cipher_TestSecKeyHash(_seckey *C.SecKey, _hash *C.SHA256) uint32 {
	seckey := inplaceSecKey(_seckey)
	hash := inplaceSHA256(_hash)
	err := cipher.TestSecKeyHash(*seckey, *hash)
	return libErrorCode(err)
}
