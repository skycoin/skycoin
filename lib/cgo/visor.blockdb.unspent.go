package main

import (
	blockdb "github.com/skycoin/skycoin/src/visor/blockdb"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_blockdb_NewErrUnspentNotExist
func SKY_blockdb_NewErrUnspentNotExist(_uxID string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uxID := _uxID
	____return_err := blockdb.NewErrUnspentNotExist(uxID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_ErrUnspentNotExist_Error
func SKY_blockdb_ErrUnspentNotExist_Error(_e *C.ErrUnspentNotExist, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := *(*ErrUnspentNotExist)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_blockdb_NewUnspentPool
func SKY_blockdb_NewUnspentPool(_db *C.DB, _arg1 *C.Unspents) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1, ____return_err := blockdb.NewUnspentPool(db)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUnspents))
	}
	return
}

// export SKY_blockdb_Unspents_ProcessBlock
func SKY_blockdb_Unspents_ProcessBlock(_up *C.Unspents, _b *C.SignedBlock, _arg1 *C.TxHandler) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	up := (*Unspents)(unsafe.Pointer(_up))
	__arg1 := up.ProcessBlock(b)
	return
}

// export SKY_blockdb_Unspents_GetArray
func SKY_blockdb_Unspents_GetArray(_up *C.Unspents, _hashes *C.GoSlice_, _arg1 *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	up := (*Unspents)(unsafe.Pointer(_up))
	__arg1, ____return_err := up.GetArray(hashes)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_Unspents_Get
func SKY_blockdb_Unspents_Get(_up *C.Unspents, _h *C.SHA256, _arg1 *C.UxOut, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	up := (*Unspents)(unsafe.Pointer(_up))
	__arg1, __arg2 := up.Get(h)
	*_arg2 = __arg2
	return
}

// export SKY_blockdb_Unspents_GetAll
func SKY_blockdb_Unspents_GetAll(_up *C.Unspents, _arg0 *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	up := (*Unspents)(unsafe.Pointer(_up))
	__arg0, ____return_err := up.GetAll()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_Unspents_Len
func SKY_blockdb_Unspents_Len(_up *C.Unspents, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	up := (*Unspents)(unsafe.Pointer(_up))
	__arg0 := up.Len()
	*_arg0 = __arg0
	return
}

// export SKY_blockdb_Unspents_Contains
func SKY_blockdb_Unspents_Contains(_up *C.Unspents, _h *C.SHA256, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	up := (*Unspents)(unsafe.Pointer(_up))
	__arg1 := up.Contains(h)
	*_arg1 = __arg1
	return
}

// export SKY_blockdb_Unspents_GetUnspentsOfAddrs
func SKY_blockdb_Unspents_GetUnspentsOfAddrs(_up *C.Unspents, _addrs *C.GoSlice_, _arg1 *C.AddressUxOuts) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	up := (*Unspents)(unsafe.Pointer(_up))
	__arg1 := up.GetUnspentsOfAddrs(addrs)
	return
}

// export SKY_blockdb_Unspents_GetUxHash
func SKY_blockdb_Unspents_GetUxHash(_up *C.Unspents, _arg0 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	up := (*Unspents)(unsafe.Pointer(_up))
	__arg0 := up.GetUxHash()
	return
}
