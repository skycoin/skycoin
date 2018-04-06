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

// export SKY_encrypt_Sha256Xor_Encrypt
func SKY_encrypt_Sha256Xor_Encrypt(_s *C.Sha256Xor, _data *C.GoSlice_, _password *C.GoSlice_, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := *(*cipher.Sha256Xor)(unsafe.Pointer(_s))
	data := *(*[]byte)(unsafe.Pointer(_data))
	password := *(*[]byte)(unsafe.Pointer(_password))
	__arg2, ____return_err := s.Encrypt(data, password)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

// export SKY_encrypt_Sha256Xor_Decrypt
func SKY_encrypt_Sha256Xor_Decrypt(_s *C.Sha256Xor, _data *C.GoSlice_, _password *C.GoSlice_, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := *(*cipher.Sha256Xor)(unsafe.Pointer(_s))
	data := *(*[]byte)(unsafe.Pointer(_data))
	password := *(*[]byte)(unsafe.Pointer(_password))
	__arg2, ____return_err := s.Decrypt(data, password)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}
