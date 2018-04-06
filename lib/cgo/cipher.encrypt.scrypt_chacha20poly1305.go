package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_encrypt_ScryptChacha20poly1305_Encrypt
func SKY_encrypt_ScryptChacha20poly1305_Encrypt(_s *C.ScryptChacha20poly1305, _data, _password *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := *(*cipher.ScryptChacha20poly1305)(unsafe.Pointer(_s))
	data := *(*[]byte)(unsafe.Pointer(_data))
	password := *(*[]byte)(unsafe.Pointer(_password))
	__arg1, ____return_err := s.Encrypt(data, password)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_encrypt_ScryptChacha20poly1305_Decrypt
func SKY_encrypt_ScryptChacha20poly1305_Decrypt(_s *C.ScryptChacha20poly1305, _data, _password *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := *(*cipher.ScryptChacha20poly1305)(unsafe.Pointer(_s))
	data := *(*[]byte)(unsafe.Pointer(_data))
	password := *(*[]byte)(unsafe.Pointer(_password))
	__arg1, ____return_err := s.Decrypt(data, password)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}
