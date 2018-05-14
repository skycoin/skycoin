package main

import (
	api "github.com/skycoin/skycoin/src/api"
	daemon "github.com/skycoin/skycoin/src/daemon"
	visor "github.com/skycoin/skycoin/src/visor"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_api_NewReadableTransaction
func SKY_api_NewReadableTransaction(_t *C.daemon__TransactionResult, _inputs []C.visor__ReadableTransactionInput, _arg2 *C.api__ReadableTransaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	t := *(*daemon.TransactionResult)(unsafe.Pointer(_t))
	inputs := *(*[]visor.ReadableTransactionInput)(unsafe.Pointer(&_inputs))
	__arg2 := api.NewReadableTransaction(t, inputs)
	*_arg2 = *(*C.api__ReadableTransaction)(unsafe.Pointer(&__arg2))
	return
}
