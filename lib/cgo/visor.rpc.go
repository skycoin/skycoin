package main

import (
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

// export SKY_visor_NewTransactionResult
func SKY_visor_NewTransactionResult(_tx *C.visor__Transaction, _arg1 *C.visor__TransactionResult) (____error_code uint32) {
	____error_code = 0
	tx := (*visor.Transaction)(unsafe.Pointer(_tx))
	__arg1, ____return_err := visor.NewTransactionResult(tx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransactionResult))
	}
	return
}

// export SKY_visor_NewTransactionResults
func SKY_visor_NewTransactionResults(_txs *C.GoSlice_, _arg1 *C.visor__TransactionResults) (____error_code uint32) {
	____error_code = 0
	txs := *(*[]Transaction)(unsafe.Pointer(_txs))
	__arg1, ____return_err := visor.NewTransactionResults(txs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransactionResults))
	}
	return
}
