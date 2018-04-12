package main

import (
	gui "github.com/skycoin/skycoin/src/gui"
	visor "github.com/skycoin/skycoin/src/visor"
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
func SKY_gui_NewReadableTransaction(_t *C.visor__TransactionResult, _inputs *C.GoSlice_, _arg2 *C.gui__ReadableTransaction) (____error_code uint32) {
	____error_code = 0
	t := *(*visor.TransactionResult)(unsafe.Pointer(_t))
	__arg2 := gui.NewReadableTransaction(t, inputs)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofReadableTransaction))
	return
}
