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

	__rd := (*[1 << 30]byte)(
		unsafe.Pointer(_rd))[:SizeofSecKey:SizeofSecKey]
	rd := (*cipher.Ripemd160)(unsafe.Pointer(&__rd))
	rd.Set(_b)
	return libErrorCode(nil)
}

//export SKY_cipher_HashRipemd160
func SKY_cipher_HashRipemd160(_data []byte, _arg1 *C.Ripemd160) {
	rd := cipher.HashRipemd160(_data)
	__arg1 := (*[1 << 30]byte)(
		unsafe.Pointer(_arg1))[:SizeofSecKey:SizeofSecKey]
	arg1 := (*cipher.Ripemd160)(unsafe.Pointer(&__arg1))
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

	__g := (*[1 << 30]byte)(
		unsafe.Pointer(_g))[:SizeofSecKey:SizeofSecKey]
	g := (*cipher.SHA256)(unsafe.Pointer(&__g))
	g.Set(_b)
	return libErrorCode(nil)
}

//export SKY_cipher_SHA256_Hex
func SKY_cipher_SHA256_Hex(_g *C.SHA256) string {
	__g := (*[1 << 30]byte)(
		unsafe.Pointer(_g))[:SizeofSecKey:SizeofSecKey]
	g := (*cipher.SHA256)(unsafe.Pointer(&__g))
	return g.Hex()
}

//export SKY_cipher_SHA256_Xor
func SKY_cipher_SHA256_Xor(_g *C.SHA256, _b *C.SHA256, _arg1 *C.SHA256) {
	__g := (*[1 << 30]byte)(
		unsafe.Pointer(_g))[:SizeofSecKey:SizeofSecKey]
	g := (*cipher.SHA256)(unsafe.Pointer(&__g))
	__b := (*[1 << 30]byte)(
		unsafe.Pointer(_b))[:SizeofSecKey:SizeofSecKey]
	b := (*cipher.SHA256)(unsafe.Pointer(&__b))
	__arg1 := (*[1 << 30]byte)(
		unsafe.Pointer(_arg1))[:SizeofSecKey:SizeofSecKey]
	arg1 := (*cipher.SHA256)(unsafe.Pointer(&__arg1))
	x := g.Xor(*b)
	copy(arg1[:], x[:])
}

//export SKY_cipher_SumSHA256
func SKY_cipher_SumSHA256(_b []byte, _arg1 *C.SHA256) {
	__arg1 := (*[1 << 30]byte)(
		unsafe.Pointer(_arg1))[:SizeofSecKey:SizeofSecKey]
	arg1 := (*cipher.SHA256)(unsafe.Pointer(&__arg1))
	h := cipher.SumSHA256(_b)
	copy(arg1[:], h[:])
}

//export SKY_cipher_SHA256FromHex
func SKY_cipher_SHA256FromHex(_hs string, _arg1 *C.SHA256) uint32 {
	h, err := cipher.SHA256FromHex(_hs)
	errcode := libErrorCode(err)
	if err == nil {
		__arg1 := (*[1 << 30]byte)(
			unsafe.Pointer(_arg1))[:SizeofSecKey:SizeofSecKey]
		arg1 := (*cipher.SHA256)(unsafe.Pointer(&__arg1))
		copy(arg1[:], h[:])
	}
	return errcode
}

//export SKY_cipher_DoubleSHA256
func SKY_cipher_DoubleSHA256(_b []byte, _arg1 *C.SHA256) {
	__arg1 := (*[1 << 30]byte)(
		unsafe.Pointer(_arg1))[:SizeofSecKey:SizeofSecKey]
	arg1 := (*cipher.SHA256)(unsafe.Pointer(&__arg1))
	h := cipher.DoubleSHA256(_b)
	copy(arg1[:], h[:])
}

//export SKY_cipher_AddSHA256
func SKY_cipher_AddSHA256(_a *C.SHA256, _b *C.SHA256, _arg2 *C.SHA256) {
	__a := (*[1 << 30]byte)(
		unsafe.Pointer(_a))[:SizeofSecKey:SizeofSecKey]
	a := (*cipher.SHA256)(unsafe.Pointer(&__a))
	__b := (*[1 << 30]byte)(
		unsafe.Pointer(_b))[:SizeofSecKey:SizeofSecKey]
	b := (*cipher.SHA256)(unsafe.Pointer(&__b))

	h := cipher.AddSHA256(*a, *b)
	__arg2 := (*[1 << 30]byte)(
		unsafe.Pointer(_arg2))[:SizeofSecKey:SizeofSecKey]
	arg2 := (*cipher.SHA256)(unsafe.Pointer(&__arg2))
	copy(arg2[:], h[:])
}

//export SKY_cipher_Merkle
func SKY_cipher_Merkle(_h0 *C.GoSlice_, _arg1 *C.SHA256) {
	h0 := (*[]cipher.SHA256)(unsafe.Pointer(_h0))
	h := cipher.Merkle(*h0)
	__arg1 := (*[1 << 30]byte)(
		unsafe.Pointer(_arg1))[:SizeofSecKey:SizeofSecKey]
	arg1 := (*cipher.SHA256)(unsafe.Pointer(&__arg1))
	copy(arg1[:], h[:])
}
