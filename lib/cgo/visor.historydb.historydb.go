package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	historydb "github.com/skycoin/skycoin/src/historydb"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_historydb_New
func SKY_historydb_New(_db *C.DB, _arg1 *C.HistoryDB) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1, ____return_err := historydb.New(db)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofHistoryDB))
	}
	return
}

// export SKY_historydb_HistoryDB_ResetIfNeed
func SKY_historydb_HistoryDB_ResetIfNeed(_hd *C.HistoryDB) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hd := (*cipher.HistoryDB)(unsafe.Pointer(_hd))
	____return_err := hd.ResetIfNeed()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_HistoryDB_GetUxout
func SKY_historydb_HistoryDB_GetUxout(_hd *C.HistoryDB, _uxID *C.SHA256, _arg1 *C.UxOut) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hd := (*cipher.HistoryDB)(unsafe.Pointer(_hd))
	__arg1, ____return_err := hd.GetUxout(uxID)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUxOut))
	}
	return
}

// export SKY_historydb_HistoryDB_ParseBlock
func SKY_historydb_HistoryDB_ParseBlock(_hd *C.HistoryDB, _b *C.Block) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hd := (*cipher.HistoryDB)(unsafe.Pointer(_hd))
	____return_err := hd.ParseBlock(b)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_HistoryDB_GetTransaction
func SKY_historydb_HistoryDB_GetTransaction(_hd *C.HistoryDB, _hash *C.SHA256, _arg1 *C.Transaction) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hd := *(*cipher.HistoryDB)(unsafe.Pointer(_hd))
	__arg1, ____return_err := hd.GetTransaction(hash)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	}
	return
}

// export SKY_historydb_HistoryDB_GetAddrUxOuts
func SKY_historydb_HistoryDB_GetAddrUxOuts(_hd *C.HistoryDB, _address *C.Address, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hd := *(*cipher.HistoryDB)(unsafe.Pointer(_hd))
	__arg1, ____return_err := hd.GetAddrUxOuts(address)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_HistoryDB_GetAddrTxns
func SKY_historydb_HistoryDB_GetAddrTxns(_hd *C.HistoryDB, _address *C.Address, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hd := *(*cipher.HistoryDB)(unsafe.Pointer(_hd))
	__arg1, ____return_err := hd.GetAddrTxns(address)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_HistoryDB_ForEach
func SKY_historydb_HistoryDB_ForEach(_hd *C.HistoryDB, _f C.Handle) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hd := *(*cipher.HistoryDB)(unsafe.Pointer(_hd))
	____return_err := hd.ForEach(f)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
