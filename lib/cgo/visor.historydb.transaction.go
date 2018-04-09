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

// export SKY_historydb_Transaction_Hash
func SKY_historydb_Transaction_Hash(_tx *C.Transaction, _arg0 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tx := (*Transaction)(unsafe.Pointer(_tx))
	__arg0 := tx.Hash()
	return
}

// export SKY_historydb_transactions_Add
func SKY_historydb_transactions_Add(_txs transactions, _t *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txs := (*transactions)(unsafe.Pointer(_txs))
	t := (*Transaction)(unsafe.Pointer(_t))
	____return_err := txs.Add(t)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_transactions_Get
func SKY_historydb_transactions_Get(_txs transactions, _hash *C.SHA256, _arg1 *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txs := (*transactions)(unsafe.Pointer(_txs))
	__arg1, ____return_err := txs.Get(hash)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	}
	return
}

// export SKY_historydb_transactions_GetSlice
func SKY_historydb_transactions_GetSlice(_txs transactions, _hashes *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txs := (*transactions)(unsafe.Pointer(_txs))
	__arg1, ____return_err := txs.GetSlice(hashes)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_transactions_IsEmpty
func SKY_historydb_transactions_IsEmpty(_txs transactions, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txs := (*transactions)(unsafe.Pointer(_txs))
	__arg0 := txs.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_transactions_Reset
func SKY_historydb_transactions_Reset(_txs transactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txs := (*transactions)(unsafe.Pointer(_txs))
	____return_err := txs.Reset()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
