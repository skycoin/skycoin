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
func SKY_cipher_Ripemd160_Set(_rd *C.cipher__Ripemd160, _b []byte) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()

	rd := (*cipher.Ripemd160)(unsafe.Pointer(_rd))

	rd.Set(_b)
	return libErrorCode(nil)
}

//export SKY_cipher_HashRipemd160
func SKY_cipher_HashRipemd160(_data []byte, _arg1 *C.cipher__Ripemd160) uint32 {
	rd := cipher.HashRipemd160(_data)

	copyToBuffer(reflect.ValueOf(rd[:]), unsafe.Pointer(_arg1), uint(SizeofRipemd160))
	return SKY_OK
}

//export SKY_cipher_SHA256_Set
func SKY_cipher_SHA256_Set(_g *C.cipher__SHA256, _b []byte) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()

	g := (*cipher.SHA256)(unsafe.Pointer(_g))

	g.Set(_b)
	return libErrorCode(nil)
}

//export SKY_cipher_SHA256_Hex
func SKY_cipher_SHA256_Hex(_g *C.cipher__SHA256, _arg1 *C.GoString_) uint32 {
	g := (*cipher.SHA256)(unsafe.Pointer(_g))
	copyString(g.Hex(), _arg1)
	return SKY_OK
}

//export SKY_cipher_SHA256_Xor
func SKY_cipher_SHA256_Xor(_g *C.cipher__SHA256, _b *C.cipher__SHA256, _arg1 *C.cipher__SHA256) uint32 {
	g := (*cipher.SHA256)(unsafe.Pointer(_g))
	b := (*cipher.SHA256)(unsafe.Pointer(_b))

	x := g.Xor(*b)
	copyToBuffer(reflect.ValueOf(x[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	return SKY_OK
}

//export SKY_cipher_SumSHA256
func SKY_cipher_SumSHA256(_b []byte, _arg1 *C.cipher__SHA256) (errcode uint32) {
	defer func() {
		errcode = catchApiPanic(errcode, recover())
	}()

	h := cipher.SumSHA256(_b)

	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	return libErrorCode(nil)
}

//export SKY_cipher_SHA256FromHex
func SKY_cipher_SHA256FromHex(_hs string, _arg1 *C.cipher__SHA256) uint32 {
	h, err := cipher.SHA256FromHex(_hs)
	errcode := libErrorCode(err)
	if err == nil {
		copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	}
	return errcode
}

//export SKY_cipher_DoubleSHA256
func SKY_cipher_DoubleSHA256(_b []byte, _arg1 *C.cipher__SHA256) uint32 {
	h := cipher.DoubleSHA256(_b)
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	return SKY_OK
}

//export SKY_cipher_AddSHA256
func SKY_cipher_AddSHA256(_a *C.cipher__SHA256, _b *C.cipher__SHA256, _arg2 *C.cipher__SHA256) uint32 {
	a := (*cipher.SHA256)(unsafe.Pointer(_a))
	b := (*cipher.SHA256)(unsafe.Pointer(_b))

	h := cipher.AddSHA256(*a, *b)
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg2), uint(SizeofSHA256))
	return SKY_OK
}

//export SKY_cipher_Merkle
func SKY_cipher_Merkle(_h0 *[]C.cipher__SHA256, _arg1 *C.cipher__SHA256) uint32 {
	h0 := (*[]cipher.SHA256)(unsafe.Pointer(_h0))
	h := cipher.Merkle(*h0)
	copyToBuffer(reflect.ValueOf(h[:]), unsafe.Pointer(_arg1), uint(SizeofSHA256))
	return SKY_OK
}

//export SKY_cipher_MustSumSHA256
func SKY_cipher_MustSumSHA256(_b []byte, _n int, _arg2 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	n := _n
	__arg2 := cipher.MustSumSHA256(b, n)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofSHA256))
	return
}
