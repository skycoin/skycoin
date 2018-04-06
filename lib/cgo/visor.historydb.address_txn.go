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

// export SKY_historydb_addressTxns_Get
func SKY_historydb_addressTxns_Get(_atx addressTxns, _address *C.Address, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	atx := (*cipher.addressTxns)(unsafe.Pointer(_atx))
	__arg1, ____return_err := atx.Get(address)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_addressTxns_IsEmpty
func SKY_historydb_addressTxns_IsEmpty(_atx addressTxns, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	atx := (*cipher.addressTxns)(unsafe.Pointer(_atx))
	__arg0 := atx.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_addressTxns_Reset
func SKY_historydb_addressTxns_Reset(_atx addressTxns) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	atx := (*cipher.addressTxns)(unsafe.Pointer(_atx))
	____return_err := atx.Reset()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
