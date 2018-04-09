package main

import (
	chacha20poly1305 "github.com/skycoin/skycoin/src/cipher/chacha20poly1305"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_chacha20poly1305_New
func SKY_chacha20poly1305_New(_key *C.GoSlice_, _arg1 *C.AEAD) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	key := *(*[]byte)(unsafe.Pointer(_key))
	__arg1, ____return_err := chacha20poly1305.New(key)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_chacha20poly1305_chacha20poly1305_NonceSize
func SKY_chacha20poly1305_chacha20poly1305_NonceSize(_c chacha20poly1305, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := (*chacha20poly1305)(unsafe.Pointer(_c))
	__arg0 := c.NonceSize()
	*_arg0 = __arg0
	return
}

// export SKY_chacha20poly1305_chacha20poly1305_Overhead
func SKY_chacha20poly1305_chacha20poly1305_Overhead(_c chacha20poly1305, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := (*chacha20poly1305)(unsafe.Pointer(_c))
	__arg0 := c.Overhead()
	*_arg0 = __arg0
	return
}

// export SKY_chacha20poly1305_chacha20poly1305_Seal
func SKY_chacha20poly1305_chacha20poly1305_Seal(_c chacha20poly1305, _dst, _nonce, _plaintext, _additionalData *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := (*chacha20poly1305)(unsafe.Pointer(_c))
	dst := *(*[]byte)(unsafe.Pointer(_dst))
	nonce := *(*[]byte)(unsafe.Pointer(_nonce))
	plaintext := *(*[]byte)(unsafe.Pointer(_plaintext))
	additionalData := *(*[]byte)(unsafe.Pointer(_additionalData))
	__arg1 := c.Seal(dst, nonce, plaintext, additionalData)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_chacha20poly1305_chacha20poly1305_Open
func SKY_chacha20poly1305_chacha20poly1305_Open(_c chacha20poly1305, _dst, _nonce, _ciphertext, _additionalData *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := (*chacha20poly1305)(unsafe.Pointer(_c))
	dst := *(*[]byte)(unsafe.Pointer(_dst))
	nonce := *(*[]byte)(unsafe.Pointer(_nonce))
	ciphertext := *(*[]byte)(unsafe.Pointer(_ciphertext))
	additionalData := *(*[]byte)(unsafe.Pointer(_additionalData))
	__arg1, ____return_err := c.Open(dst, nonce, ciphertext, additionalData)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}
