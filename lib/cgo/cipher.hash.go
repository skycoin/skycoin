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

//export SKY_cipher_Ripemd160_Set
func SKY_cipher_Ripemd160_Set(_rd *C.cipher__Ripemd160, _b []byte) (____error_code uint32) {
	rd := (*cipher.Ripemd160)(unsafe.Pointer(_rd))

	err := rd.Set(_b)
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_HashRipemd160
func SKY_cipher_HashRipemd160(_data []byte, _arg1 *C.cipher__Ripemd160) (____error_code uint32) {
	rd := cipher.HashRipemd160(_data)

	copyToBuffer(reflect.ValueOf(rd[:]), unsafe.Pointer(_arg1), uint(SizeofRipemd160))
	return
}

//export SKY_cipher_SHA256_Set
func SKY_cipher_SHA256_Set(_g *C.cipher__SHA256, _b []byte) (____error_code uint32) {
	g := (*cipher.SHA256)(unsafe.Pointer(_g))

	err := g.Set(_b)
	____error_code = libErrorCode(err)
	return
}

//export SKY_cipher_SHA256_Hex
func SKY_cipher_SHA256_Hex(_g *C.cipher__SHA256, _arg1 *C.GoString_) (____error_code uint32) {
	g := (*cipher.SHA256)(unsafe.Pointer(_g))
	copyString(g.Hex(), _arg1)
	return
}

//export SKY_cipher_SHA256_Xor
func SKY_cipher_SHA256_Xor(_g *C.cipher__SHA256, _b *C.cipher__SHA256, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	g := (*cipher.SHA256)(unsafe.Pointer(_g))
	b := (*cipher.SHA256)(unsafe.Pointer(_b))

	x := g.Xor(*b)
	copyToBuffer(reflect.ValueOf(x[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	return
}

//export SKY_cipher_SumSHA256
func SKY_cipher_SumSHA256(_b []byte, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	h := cipher.SumSHA256(_b)

	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	return
}

//export SKY_cipher_SHA256FromHex
func SKY_cipher_SHA256FromHex(_hs string, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	h, err := cipher.SHA256FromHex(_hs)
	____error_code = libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	}
	return
}

//export SKY_cipher_DoubleSHA256
func SKY_cipher_DoubleSHA256(_b []byte, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	h := cipher.DoubleSHA256(_b)
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	return
}

//export SKY_cipher_AddSHA256
func SKY_cipher_AddSHA256(_a *C.cipher__SHA256, _b *C.cipher__SHA256, _arg2 *C.cipher__SHA256) (____error_code uint32) {
	a := (*cipher.SHA256)(unsafe.Pointer(_a))
	b := (*cipher.SHA256)(unsafe.Pointer(_b))

	h := cipher.AddSHA256(*a, *b)
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg2), uint(SizeofSHA256))
	return
}

//export SKY_cipher_Merkle
func SKY_cipher_Merkle(_h0 *[]C.cipher__SHA256, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	h0 := (*[]cipher.SHA256)(unsafe.Pointer(_h0))
	h := cipher.Merkle(*h0)
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	return
}

//export SKY_cipher_SHA256_Null
func SKY_cipher_SHA256_Null(_g *C.cipher__SHA256, _arg0 *bool) (____error_code uint32) {
	g := (*cipher.SHA256)(unsafe.Pointer(_g))
	*_arg0 = g.Null()
	return
}
