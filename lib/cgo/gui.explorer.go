package main

import (
	gui "github.com/skycoin/skycoin/src/gui"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gui_NewReadableTransaction
func SKY_gui_NewReadableTransaction(_t *C.TransactionResult, _inputs *C.GoSlice_, _arg2 *C.ReadableTransaction) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg2 := gui.NewReadableTransaction(t, inputs)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofReadableTransaction))
	return
}
