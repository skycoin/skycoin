package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_historydb_historyMeta_ParsedHeight
func SKY_historydb_historyMeta_ParsedHeight(_hm historyMeta, _arg0 *int64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hm := (*cipher.historyMeta)(unsafe.Pointer(_hm))
	__arg0 := hm.ParsedHeight()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_historyMeta_SetParsedHeight
func SKY_historydb_historyMeta_SetParsedHeight(_hm historyMeta, _h uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hm := (*cipher.historyMeta)(unsafe.Pointer(_hm))
	h := _h
	____return_err := hm.SetParsedHeight(h)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_historyMeta_SetParsedHeightWithTx
func SKY_historydb_historyMeta_SetParsedHeightWithTx(_hm historyMeta, _tx *C.Tx, _h uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hm := (*cipher.historyMeta)(unsafe.Pointer(_hm))
	h := _h
	____return_err := hm.SetParsedHeightWithTx(tx, h)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_historydb_historyMeta_IsEmpty
func SKY_historydb_historyMeta_IsEmpty(_hm historyMeta, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hm := (*cipher.historyMeta)(unsafe.Pointer(_hm))
	__arg0 := hm.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_historydb_historyMeta_Reset
func SKY_historydb_historyMeta_Reset(_hm historyMeta) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hm := (*cipher.historyMeta)(unsafe.Pointer(_hm))
	____return_err := hm.Reset()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
