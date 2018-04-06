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

// export SKY_historydb_addressUx_Get
func SKY_historydb_addressUx_Get(_au addressUx, _address *C.Address, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	au := (*cipher.addressUx)(unsafe.Pointer(_au))
	__arg1, ____return_err := au.Get(address)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_addressUx_Add
func SKY_historydb_addressUx_Add(_au addressUx, _address *C.Address, _uxHash *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	au := (*cipher.addressUx)(unsafe.Pointer(_au))
	____return_err := au.Add(address, uxHash)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_addressUx_IsEmpty
func SKY_historydb_addressUx_IsEmpty(_au addressUx, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	au := (*cipher.addressUx)(unsafe.Pointer(_au))
	__arg0 := au.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_addressUx_Reset
func SKY_historydb_addressUx_Reset(_au addressUx) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	au := (*cipher.addressUx)(unsafe.Pointer(_au))
	____return_err := au.Reset()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
