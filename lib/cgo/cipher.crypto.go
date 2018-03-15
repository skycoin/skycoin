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
	slice := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_slice), )
	return slice.Len()
}

//export SKY_cipher_PubKeySlice_Less
func SKY_cipher_PubKeySlice_Less(_slice *C.PubKeySlice, _i, _j int) bool {
	slice := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_slice))
	return slice.Less(_i, _j)
}

//export SKY_cipher_PubKeySlice_Swap
func SKY_cipher_PubKeySlice_Swap(_slice *C.PubKeySlice, _i, _j int) {
	slice := (*cipher.PubKeySlice)inplaceArrayObj(unsafe.Pointer(_slice))
	slice.Swap(_i, _j)
}

//export SKY_cipher_RandByte
func SKY_cipher_RandByte(_n int, _arg1 *C.GoSlice_) {
	if _n > _arg1.cap {
		_arg1.len = _arg1.cap - _n
		return
	}
	b := cipher.RandByte(_n)
	_arg1.len = len(b)
	C.memcpy(unsafe.Pointer(_arg1.data), unsafe.Pointer(&b), C.size_t(_arg1.len))
}

//export SKY_cipher_NewPubKey
func SKY_cipher_NewPubKey(_b []byte, _arg1 *C.PubKey) (retVal C.uint32) {
	defer func() {
		if err := recover(); err != nil {
			retVal = libErrorCode(err)
			retVal = 0
		}
	}()

	arg1 := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_arg1))
	pubkey := cipher.NewPubKey(b)
	copy(arg1[:], pubkey[:])
	return 1
}

//export SKY_cipher_PubKeyFromHex
func SKY_cipher_PubKeyFromHex(_s string, _arg1 *C.PubKey) C.uint32 {
	pubkey, err := cipher.PubKeyFromHex(s)
	if err == nil {
		return 0;
	}
	arg1 := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_arg1))
	copy(arg1[:], pubkey[:])
	return 1
}

//export SKY_cipher_PubKeyFromSecKey
func SKY_cipher_PubKeyFromSecKey(_seckey *C.SecKey, _arg1 *C.PubKey) (retaval C.uint32) {
	defer func() {
		if err := recover(); err != nil {
			retVal = 0
		}
	}()

	pubkey := cipher.PubKeyFromSecKey(seckey)
	arg1 := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_arg1))
	copy(arg1[:], pubkey[:])
	return 1
}

//export SKY_cipher_PubKeyFromSig
func SKY_cipher_PubKeyFromSig(_sig *C.Sig, _hash *C.SHA256, _arg2 *C.PubKey) C.uint32 {
	sig := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_sig))
	h := (*cipher.SHA256)inplaceArrayObj(unsafe.Pointer(_hash))
	arg2 := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_arg2))
	pubkey, err := cipher.PubKeyFromSig(sig, hash)
	if err == nil {
		return 0;
	}
	copy(arg2[:], pubkey[:])
	return 1;
}

//export SKY_cipher_PubKey_Verify
func SKY_cipher_PubKey_Verify(_pk *C.PubKey) C.uint32 {
	pk := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_pk))
	if pk.Verify() {
		return 1;
	} else {
		retur 0;
	}
}

//export SKY_cipher_PubKey_Hex
func SKY_cipher_PubKey_Hex(_pk *C.PubKey) string {
	pk := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_pk))
	return pk.Hex()
}

//export SKY_cipher_PubKey_ToAddressHash
func SKY_cipher_PubKey_ToAddressHash(_pk *C.PubKey, _arg0 *C.Ripemd160) {
	pk := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_pk))
	arg0 := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_arg0))
	h := pk.ToAddressHash()
	copy(arg0[:], h[:])
}

//export SKY_cipher_NewSecKey
func SKY_cipher_NewSecKey(_b []byte, _arg1 *C.SecKey) {
	arg1 := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_arg1))
	sk := cipher.NewSecKey(b)
	copy(arg1[:], sk[:])
}

//export SKY_cipher_SecKeyFromHex
func SKY_cipher_SecKeyFromHex(_s string, _arg1 *C.SecKey) C.uint32 {
	arg1 := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_arg1))
	sk, err := cipher.SecKeyFromHex(s)
	errcode := libErrorCode(err)
	if err != nil {
		copy(arg1[:], sk[:])
	}
	return errcode
}

//export SKY_cipher_SecKey_Verify
func SKY_cipher_SecKey_Verify(_sk *C.SecKey) C.uint32 {
	sk := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_sk))
	err := sk.Verify()
	return libErrorCode(err)
}

//export SKY_cipher_SecKey_Hex
func SKY_cipher_SecKey_Hex(_sk *C.SecKey) string {
	sk := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_sk))
	return sk.Hex()
}

//export SKY_cipher_ECDH
func SKY_cipher_ECDH(_pub *C.PubKey, _sec *C.SecKey, _arg2 *C.GoSlice_) {
	pub := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_pub))
	sec := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_sec))
	b := cipher.ECDH(pub, sec)
	blen := len(b)
	if _arg2.cap < blen {
		_arg2.len = _arg2.cap - blen
	}
	else {
		C.memcpy(unsafe.Pointer(_arg2.data), unsafe.Pointer(&b[:]), C.size_t(blen))
		_arg2.len = blen
	}
}

//export SKY_cipher_NewSig
func SKY_cipher_NewSig(_b []byte, _arg1 *C.Sig) {
	arg1 := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_arg1))
	s := cipher.NewSig(_b)
	copy(arg1[:], s[:])
}

//export SKY_cipher_SigFromHex
func SKY_cipher_SigFromHex(_s string, _arg1 *C.Sig) C.uint32 {
	arg1 := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_arg1))
	s, err := cipher.SigFromHex(_s)
	errcode := libErrorCode(err)
	if err == nil {
		copy(_arg1[:], s[:])
	}
}

//export SKY_cipher_Sig_Hex
func SKY_cipher_Sig_Hex(_s *C.Sig) string {
	s := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_s))
	return s.Hex()
}

//export SKY_cipher_SignHash
func SKY_cipher_SignHash(_hash *C.SHA256, _sec *C.SecKey, _arg2 *C.Sig) {
	hash := (*cipher.SHA256)inplaceArrayObj(unsafe.Pointer(_hash))
	sec := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_sec))
	arg2 := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_arg2))
	s := cipher.SignHash(hash, sec)
	copy(arg2[:], s[:])
}

//export SKY_cipher_ChkSig
func SKY_cipher_ChkSig(_address *C.Address, _hash *C.SHA256, _sig *C.Sig) C.uint32 {
	address := (*cipher.Address)inplaceArrayObj(unsafe.Pointer(_address))
	hash := (*cipher.SHA256)inplaceArrayObj(unsafe.Pointer(_hash))
	sig := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_sig))
	err := cipher.ChkSig(address, hash, sig)
	return libErrorCode(err)
}

//export SKY_cipher_VerifySignedHash
func SKY_cipher_VerifySignedHash(_sig *C.Sig, _hash *C.SHA256) C.uint32 {
	hash := (*cipher.SHA256)inplaceArrayObj(unsafe.Pointer(_hash))
	sig := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_sig))
	err := cipher.VerifySignedHash(sig, hash)
	return libErrorCode(err)
}

//export SKY_cipher_VerifySignature
func SKY_cipher_VerifySignature(_pubkey *C.PubKey, _sig *C.Sig, _hash *C.SHA256) C.uint32 {
	pubkey := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_pubkey))
	sig := (*cipher.Sig)inplaceArrayObj(unsafe.Pointer(_sig))
	hash := (*cipher.SHA256)inplaceArrayObj(unsafe.Pointer(_hash))
	err := cipher.VerifySignature(pubkey, sig, hash)
	return libErrorCode(err)
}

//export SKY_cipher_GenerateKeyPair
func SKY_cipher_GenerateKeyPair(_arg0 *C.PubKey, _arg1 *C.SecKey) {
	arg0 := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_arg0))
	arg1 := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_arg1))
	p, s := cipher.GenerateKeyPair()
	copy(arg0[:], p[:])
	copy(arg1[:], s[:])
}

//export SKY_cipher_GenerateDeterministicKeyPair
func SKY_cipher_GenerateDeterministicKeyPair(_seed []byte, _arg1 *C.PubKey, _arg2 *C.SecKey) {
	arg1 := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_arg1))
	arg2 := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_arg2))
	p, s := cipher.GenerateDeterministicKeyPair(_seed)
	copy(arg0[:], p[:])
	copy(arg1[:], s[:])
}

//export SKY_cipher_DeterministicKeyPairIterator
func SKY_cipher_DeterministicKeyPairIterator(_seed []byte, _arg1 *C.GoSlice_, _arg2 *C.PubKey, _arg3 *C.SecKey) {
	arg1 := (*[]byte)unsafe.Pointer(_arg1)
	arg2 := (*cipher.PubKey)inplaceArrayObj(unsafe.Pointer(_arg2))
	arg3 := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_arg3))
	h, p, s := cipher.DeterministicKeyPairIterator(seed)
	hlen := len(h)
	if hlen > _arg1.cap {
		_arg1.len = _arg1.cap - hlen
		return
	}
	copy(arg1[:], h[:])
	copy(arg2[:], p[:])
	copy(arg3[:], s[:])
}

//export SKY_cipher_GenerateDeterministicKeyPairs
func SKY_cipher_GenerateDeterministicKeyPairs(_seed []byte, _n int, _arg2 *C.GoSlice_) {
	arg2 := (*[]SecKey)unsafe.Pointer(_arg2)
	sks := cipher.GenerateDeterministicKeyPairs(_seed, _n)
	skslen := len(sks)
	if skslen > _arg1.cap {
		_arg1.len = _arg1.cap - skslen
		return
	}
	copy(arg2[:], sks[:])
}

//export SKY_cipher_GenerateDeterministicKeyPairsSeed
func SKY_cipher_GenerateDeterministicKeyPairsSeed(_seed []byte, _n int, _arg2 *C.GoSlice_, _arg3 *C.GoSlice_) {
	arg2 := (*[]byte)unsafe.Pointer(_arg2)
	arg3 := (*[]SecKey)unsafe.Pointer(_arg3)
	h, sks := cipher.GenerateDeterministicKeyPairsSeed(seed, n)
	skslen, hlen := len(sks), len(h)
	nospace := false
	if skslen > _arg3.cap {
		_arg3.len = _arg3.cap - skslen
		nospace = true
	}
	if hlen > _arg2.cap {
		_arg2.len = _arg2.cap - hlen
		nospace = true
	}
	if nospace {
		return
	}
	copy(arg3[:], sks[:])
	copy(arg2[:], h[:])
}

//export SKY_cipher_TestSecKey
func SKY_cipher_TestSecKey(_seckey *C.SecKey) C.uint32 {
	seckey := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_seckey))
	err := cipher.TestSecKey(seckey)
	return libErrorCode(err)
}

//export SKY_cipher_TestSecKeyHash
func SKY_cipher_TestSecKeyHash(_seckey *C.SecKey, _hash *C.SHA256) C.uint32 {
	seckey := (*cipher.SecKey)inplaceArrayObj(unsafe.Pointer(_seckey))
	hash := (*cipher.SHA256)inplaceArrayObj(unsafe.Pointer(_hash))
	err := cipher.TestSecKeyHash(seckey, hash)
	return libErrorCode(err)
}
