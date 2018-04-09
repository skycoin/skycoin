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

// export SKY_historydb_addressTxns_Get
func SKY_historydb_addressTxns_Get(_atx addressTxns, _address *C.Address, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	atx := (*addressTxns)(unsafe.Pointer(_atx))
	__arg1, ____return_err := atx.Get(address)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_addressTxns_IsEmpty
func SKY_historydb_addressTxns_IsEmpty(_atx addressTxns, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	atx := (*addressTxns)(unsafe.Pointer(_atx))
	__arg0 := atx.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_addressTxns_Reset
func SKY_historydb_addressTxns_Reset(_atx addressTxns) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	atx := (*addressTxns)(unsafe.Pointer(_atx))
	____return_err := atx.Reset()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
