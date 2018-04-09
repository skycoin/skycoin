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

// export SKY_historydb_NewUxOutJSON
func SKY_historydb_NewUxOutJSON(_out *C.UxOut, _arg1 *C.UxOutJSON) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	out := (*UxOut)(unsafe.Pointer(_out))
	__arg1 := historydb.NewUxOutJSON(out)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUxOutJSON))
	return
}

// export SKY_historydb_UxOut_Hash
func SKY_historydb_UxOut_Hash(_o *C.UxOut, _arg0 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	o := *(*UxOut)(unsafe.Pointer(_o))
	__arg0 := o.Hash()
	return
}

// export SKY_historydb_UxOuts_Set
func SKY_historydb_UxOuts_Set(_ux *C.UxOuts, _out *C.UxOut) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ux := (*UxOuts)(unsafe.Pointer(_ux))
	out := *(*UxOut)(unsafe.Pointer(_out))
	____return_err := ux.Set(out)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_UxOuts_Get
func SKY_historydb_UxOuts_Get(_ux *C.UxOuts, _uxID *C.SHA256, _arg1 *C.UxOut) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ux := (*UxOuts)(unsafe.Pointer(_ux))
	__arg1, ____return_err := ux.Get(uxID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUxOut))
	}
	return
}

// export SKY_historydb_UxOuts_IsEmpty
func SKY_historydb_UxOuts_IsEmpty(_ux *C.UxOuts, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ux := (*UxOuts)(unsafe.Pointer(_ux))
	__arg0 := ux.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_UxOuts_Reset
func SKY_historydb_UxOuts_Reset(_ux *C.UxOuts) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ux := (*UxOuts)(unsafe.Pointer(_ux))
	____return_err := ux.Reset()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
