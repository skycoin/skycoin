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

//export SKY_cipher_Ripemd160_Set
func SKY_cipher_Ripemd160_Set(_rd *C.Ripemd160, _b []byte) (retVal uint32) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: Fix to be like retVal = libErrorCode(err)
			retVal = SKY_ERROR
		}
	}()

	rd := inplaceRipemd160(_rd)
	rd.Set(_b)
	return libErrorCode(nil)
}

//export SKY_cipher_HashRipemd160
func SKY_cipher_HashRipemd160(_data []byte, _arg1 *C.Ripemd160) {
	rd := cipher.HashRipemd160(_data)
	arg1 := inplaceRipemd160(_arg1)
	copy(arg1[:], rd[:])
}

//export SKY_cipher_SHA256_Set
func SKY_cipher_SHA256_Set(_g *C.SHA256, _b []byte) (retVal uint32) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: Fix to be like retVal = libErrorCode(err)
			retVal = SKY_ERROR
		}
	}()

	g := inplaceSHA256(_g)
	g.Set(_b)
	return libErrorCode(nil)
}

//export SKY_cipher_SHA256_Hex
func SKY_cipher_SHA256_Hex(_g *C.SHA256) string {
	g := inplaceSHA256(_g)
	return g.Hex()
}

//export SKY_cipher_SHA256_Xor
func SKY_cipher_SHA256_Xor(_g *C.SHA256, _b *C.SHA256, _arg1 *C.SHA256) {
	g := inplaceSHA256(_g)
	b := inplaceSHA256(_b)
	arg1 := inplaceSHA256(_arg1)
	x := g.Xor(*b)
	copy(arg1[:], x[:])
}

//export SKY_cipher_SumSHA256
func SKY_cipher_SumSHA256(_b []byte, _arg1 *C.SHA256) {
	arg1 := inplaceSHA256(_arg1)
	h := cipher.SumSHA256(_b)
	copy(arg1[:], h[:])
}

//export SKY_cipher_SHA256FromHex
func SKY_cipher_SHA256FromHex(_hs string, _arg1 *C.SHA256) uint32 {
	h, err := cipher.SHA256FromHex(_hs)
	errcode := libErrorCode(err)
	if err == nil {
		arg1 := inplaceSHA256(_arg1)
		copy(arg1[:], h[:])
	}
	return errcode
}

//export SKY_cipher_DoubleSHA256
func SKY_cipher_DoubleSHA256(_b []byte, _arg1 *C.SHA256) {
	arg1 := inplaceSHA256(_arg1)
	h := cipher.DoubleSHA256(_b)
	copy(arg1[:], h[:])
}

//export SKY_cipher_AddSHA256
func SKY_cipher_AddSHA256(_a *C.SHA256, _b *C.SHA256, _arg2 *C.SHA256) {
	a := inplaceSHA256(_a)
	b := inplaceSHA256(_b)

	h := cipher.AddSHA256(*a, *b)
	arg2 := inplaceSHA256(_arg2)
	copy(arg2[:], h[:])
}

//export SKY_cipher_Merkle
func SKY_cipher_Merkle(_h0 *C.GoSlice_, _arg1 *C.SHA256) {
	h0 := (*[]cipher.SHA256)(unsafe.Pointer(_h0))
	h := cipher.Merkle(*h0)
	arg1 := inplaceSHA256(_arg1)
	copy(arg1[:], h[:])
}
