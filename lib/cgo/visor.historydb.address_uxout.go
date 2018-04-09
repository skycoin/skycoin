package main

import (
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_historydb_addressUx_Get
func SKY_historydb_addressUx_Get(_au addressUx, _address *C.Address, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	au := (*addressUx)(unsafe.Pointer(_au))
	__arg1, ____return_err := au.Get(address)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_addressUx_Add
func SKY_historydb_addressUx_Add(_au addressUx, _address *C.Address, _uxHash *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	au := (*addressUx)(unsafe.Pointer(_au))
	____return_err := au.Add(address, uxHash)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_addressUx_IsEmpty
func SKY_historydb_addressUx_IsEmpty(_au addressUx, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	au := (*addressUx)(unsafe.Pointer(_au))
	__arg0 := au.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_addressUx_Reset
func SKY_historydb_addressUx_Reset(_au addressUx) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	au := (*addressUx)(unsafe.Pointer(_au))
	____return_err := au.Reset()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
