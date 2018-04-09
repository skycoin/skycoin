package main

import (
	historydb "github.com/skycoin/skycoin/src/visor/historydb"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_historydb_New
func SKY_historydb_New(_db *C.DB, _arg1 *C.HistoryDB) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1, ____return_err := historydb.New(db)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofHistoryDB))
	}
	return
}

// export SKY_historydb_HistoryDB_ResetIfNeed
func SKY_historydb_HistoryDB_ResetIfNeed(_hd *C.HistoryDB) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hd := (*HistoryDB)(unsafe.Pointer(_hd))
	____return_err := hd.ResetIfNeed()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_HistoryDB_GetUxout
func SKY_historydb_HistoryDB_GetUxout(_hd *C.HistoryDB, _uxID *C.SHA256, _arg1 *C.UxOut) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hd := (*HistoryDB)(unsafe.Pointer(_hd))
	__arg1, ____return_err := hd.GetUxout(uxID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUxOut))
	}
	return
}

// export SKY_historydb_HistoryDB_ParseBlock
func SKY_historydb_HistoryDB_ParseBlock(_hd *C.HistoryDB, _b *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hd := (*HistoryDB)(unsafe.Pointer(_hd))
	____return_err := hd.ParseBlock(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_HistoryDB_GetTransaction
func SKY_historydb_HistoryDB_GetTransaction(_hd *C.HistoryDB, _hash *C.SHA256, _arg1 *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hd := *(*HistoryDB)(unsafe.Pointer(_hd))
	__arg1, ____return_err := hd.GetTransaction(hash)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	}
	return
}

// export SKY_historydb_HistoryDB_GetAddrUxOuts
func SKY_historydb_HistoryDB_GetAddrUxOuts(_hd *C.HistoryDB, _address *C.Address, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hd := *(*HistoryDB)(unsafe.Pointer(_hd))
	__arg1, ____return_err := hd.GetAddrUxOuts(address)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_HistoryDB_GetAddrTxns
func SKY_historydb_HistoryDB_GetAddrTxns(_hd *C.HistoryDB, _address *C.Address, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hd := *(*HistoryDB)(unsafe.Pointer(_hd))
	__arg1, ____return_err := hd.GetAddrTxns(address)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_historydb_HistoryDB_ForEach
func SKY_historydb_HistoryDB_ForEach(_hd *C.HistoryDB, _f C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hd := *(*HistoryDB)(unsafe.Pointer(_hd))
	____return_err := hd.ForEach(f)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
