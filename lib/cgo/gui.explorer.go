package main

import (
	gui "github.com/skycoin/skycoin/src/gui"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gui_NewReadableTransaction
func SKY_gui_NewReadableTransaction(_t *C.TransactionResult, _inputs *C.GoSlice_, _arg2 *C.ReadableTransaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg2 := gui.NewReadableTransaction(t, inputs)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofReadableTransaction))
	return
}
