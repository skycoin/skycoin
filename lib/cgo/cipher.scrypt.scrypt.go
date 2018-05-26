package main

import (
	"reflect"
	"unsafe"

	scrypt "github.com/skycoin/skycoin/src/cipher/scrypt"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_scrypt_Key
func SKY_scrypt_Key(_password, _salt []byte, _N, _r, _p, _keyLen int, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	password := *(*[]byte)(unsafe.Pointer(&_password))
	salt := *(*[]byte)(unsafe.Pointer(&_salt))
	N := _N
	r := _r
	p := _p
	keyLen := _keyLen
	__arg2, ____return_err := scrypt.Key(password, salt, N, r, p, keyLen)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}
