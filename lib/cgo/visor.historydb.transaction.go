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

// export SKY_historydb_Transaction_Hash
func SKY_historydb_Transaction_Hash(_tx *C.Transaction, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	tx := (*cipher.Transaction)(unsafe.Pointer(_tx))
	__arg0 := tx.Hash()
	return
}

// export SKY_historydb_transactions_Add
func SKY_historydb_transactions_Add(_txs transactions, _t *C.Transaction) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	txs := (*cipher.transactions)(unsafe.Pointer(_txs))
	t := (*cipher.Transaction)(unsafe.Pointer(_t))
	____return_err := txs.Add(t)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_transactions_Get
func SKY_historydb_transactions_Get(_txs transactions, _hash *C.SHA256, _arg1 *C.Transaction) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	txs := (*cipher.transactions)(unsafe.Pointer(_txs))
	__arg1, ____return_err := txs.Get(hash)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	}
	return
}

// export SKY_historydb_transactions_GetSlice
func SKY_historydb_transactions_GetSlice(_txs transactions, _hashes *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	txs := (*cipher.transactions)(unsafe.Pointer(_txs))
	__arg1, ____return_err := txs.GetSlice(hashes)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_transactions_IsEmpty
func SKY_historydb_transactions_IsEmpty(_txs transactions, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	txs := (*cipher.transactions)(unsafe.Pointer(_txs))
	__arg0 := txs.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_transactions_Reset
func SKY_historydb_transactions_Reset(_txs transactions) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	txs := (*cipher.transactions)(unsafe.Pointer(_txs))
	____return_err := txs.Reset()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
