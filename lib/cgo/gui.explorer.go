package main

import (
	gui "github.com/skycoin/skycoin/src/gui"
	visor "github.com/skycoin/skycoin/src/visor"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_gui_NewReadableTransaction
func SKY_gui_NewReadableTransaction(_t *C.visor__TransactionResult, _inputs []C.visor__ReadableTransactionInput, _arg2 *C.gui__ReadableTransaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	t := *(*visor.TransactionResult)(unsafe.Pointer(_t))
	inputs := *(*[]visor.ReadableTransactionInput)(unsafe.Pointer(&_inputs))
	__arg2 := gui.NewReadableTransaction(t, inputs)
	*_arg2 = *(*C.gui__ReadableTransaction)(unsafe.Pointer(&__arg2))
	return
}
