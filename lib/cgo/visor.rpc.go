package main

import (
	visor "github.com/skycoin/skycoin/src/visor"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_visor_NewTransactionResult
func SKY_visor_NewTransactionResult(_tx *C.visor__Transaction, _arg1 *C.visor__TransactionResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tx := (*visor.Transaction)(unsafe.Pointer(_tx))
	__arg1, ____return_err := visor.NewTransactionResult(tx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.visor__TransactionResult)(unsafe.Pointer(__arg1))
	}
	return
}

//export SKY_visor_NewTransactionResults
func SKY_visor_NewTransactionResults(_txs []C.visor__Transaction, _arg1 *C.visor__TransactionResults) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txs := *(*[]visor.Transaction)(unsafe.Pointer(&_txs))
	__arg1, ____return_err := visor.NewTransactionResults(txs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.visor__TransactionResults)(unsafe.Pointer(__arg1))
	}
	return
}
