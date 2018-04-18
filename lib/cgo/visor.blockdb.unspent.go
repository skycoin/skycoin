package main

import (
	blockdb "github.com/skycoin/skycoin/src/visor/blockdb"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_blockdb_NewErrUnspentNotExist
func SKY_blockdb_NewErrUnspentNotExist(_uxID string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uxID := _uxID
	____return_err := blockdb.NewErrUnspentNotExist(uxID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_blockdb_ErrUnspentNotExist_Error
func SKY_blockdb_ErrUnspentNotExist_Error(_e *C.blockdb__ErrUnspentNotExist, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := *(*blockdb.ErrUnspentNotExist)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}
