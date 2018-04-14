package main

import (
	pbkdf2 "github.com/skycoin/skycoin/src/cipher/pbkdf2"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_pbkdf2_Key
func SKY_pbkdf2_Key(_password, _salt *C.GoSlice_, _iter, _keyLen int, _h C.Handle, _arg3 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	password := *(*[]byte)(unsafe.Pointer(_password))
	salt := *(*[]byte)(unsafe.Pointer(_salt))
	iter := _iter
	keyLen := _keyLen
	h := copyToFunc(_h)
	__arg3 := pbkdf2.Key(password, salt, iter, keyLen, h)
	copyToGoSlice(reflect.ValueOf(__arg3), _arg3)
	return
}
