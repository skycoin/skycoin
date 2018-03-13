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
func SKY_cipher_Ripemd160_Set(_rd *C.Ripemd160, _b []byte) (retVal C.uint) {
	defer func() {
		if r := recover(); r != nil {
			retVal = 0
		}
	}()

	rd := inplaceArrayObj(unsafe.Pointer(_rd), 20).(cipher.Ripemd160)
	rd.Set(_b)
	return 1
}

//export SKY_cipher_HashRipemd160
func SKY_cipher_HashRipemd160(_data []byte, _arg1 *C.Ripemd160) {
	rd := cipher.HashRipemd160(_data)
	rdout := inplaceArrayObj(unsafe.Pointer(_arg1), 20).(cipher.Ripemd160)
	copy(rdout[:], rd[:])
}

//export SKY_cipher_SHA256_Set
func SKY_cipher_SHA256_Set(_g *C.SHA256, _b []byte) (retVal C.uint) {
	defer func() {
		if r := recover(); r != nil {
			retVal = 0
		}
	}()

	g := inplaceArrayObj(unsafe.Pointer(_g), 32).(cipher.SHA256)
	g.Set(_b)
	return 1
}

//export SKY_cipher_SHA256_Hex
func SKY_cipher_SHA256_Hex(_g *C.SHA256) string {
	g := inplaceArrayObj(unsafe.Pointer(_g), 32).(cipher.SHA256)
	return g.Hex()
}

//export SKY_cipher_SHA256_Xor
func SKY_cipher_SHA256_Xor(_g *C.SHA256, _b *C.SHA256, _arg1 *C.SHA256) {
	g := inplaceArrayObj(unsafe.Pointer(_g), 32).(cipher.SHA256)
	b := inplaceArrayObj(unsafe.Pointer(_b), 32).(cipher.SHA256)
	r := inplaceArrayObj(unsafe.Pointer(_arg1), 32).(cipher.SHA256)
	x := g.Xor(b)
	copy(r[:], x[:])
}

//export SKY_cipher_SumSHA256
func SKY_cipher_SumSHA256(_b []byte, _arg1 *C.SHA256) {
	r := inplaceArrayObj(unsafe.Pointer(_arg1), 32).(cipher.SHA256)
	h := cipher.SumSHA256(_b)
	copy(r[:], h[:])
}

//export SKY_cipher_SHA256FromHex
func SKY_cipher_SHA256FromHex(_hs string, _arg1 *C.SHA256) C.uint {
	h, err := cipher.SHA256FromHex(_hs)
	if err != nil {
		return 0
	}
	r := inplaceArrayObj(unsafe.Pointer(_arg1), 32).(cipher.SHA256)
	copy(r[:], h[:])
	return 1
}

//export SKY_cipher_DoubleSHA256
func SKY_cipher_DoubleSHA256(_b []byte, _arg1 *C.SHA256) {
	r := inplaceArrayObj(unsafe.Pointer(_arg1), 32).(cipher.SHA256)
	h := cipher.DoubleSHA256(_b)
	copy(r[:], h[:])
}

//export SKY_cipher_AddSHA256
func SKY_cipher_AddSHA256(_a *C.SHA256, _b *C.SHA256, _arg2 *C.SHA256) {
	a := inplaceArrayObj(unsafe.Pointer(_a), 32).(cipher.SHA256)
	b := inplaceArrayObj(unsafe.Pointer(_b), 32).(cipher.SHA256)

	h := cipher.AddSHA256(a, b)
	r := inplaceArrayObj(unsafe.Pointer(_arg2), 32).(cipher.SHA256)
	copy(r[:], h[:])
}

//export SKY_cipher_Merkle
func SKY_cipher_Merkle(_h0 []C.SHA256, _arg1 *C.SHA256) {
	h0 := []cipher.SHA256{}
	for _, _h := range _h0 {
		h := inplaceArrayObj(unsafe.Pointer(&_h), 32).(cipher.SHA256)
		h0 = append(h0, h)
	}
	h := cipher.Merkle(h0)
	r := inplaceArrayObj(unsafe.Pointer(_arg1), 32).(cipher.SHA256)
	copy(r[:], h[:])
}
