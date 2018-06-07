package main

import (
	"reflect"
	"unsafe"

	coin "github.com/skycoin/skycoin/src/coin"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_coin_UxOut_Hash
func SKY_coin_UxOut_Hash(_uo *C.coin__UxOut, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uo := (*coin.UxOut)(unsafe.Pointer(_uo))
	__arg0 := uo.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_UxOut_SnapshotHash
func SKY_coin_UxOut_SnapshotHash(_uo *C.coin__UxOut, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uo := (*coin.UxOut)(unsafe.Pointer(_uo))
	__arg0 := uo.SnapshotHash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_UxBody_Hash
func SKY_coin_UxBody_Hash(_ub *C.coin__UxBody, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ub := (*coin.UxBody)(unsafe.Pointer(_ub))
	__arg0 := ub.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_UxOut_CoinHours
func SKY_coin_UxOut_CoinHours(_uo *C.coin__UxOut, _t uint64, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uo := (*coin.UxOut)(unsafe.Pointer(_uo))
	t := _t
	__arg1, ____return_err := uo.CoinHours(t)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

//export SKY_coin_UxArray_Hashes
func SKY_coin_UxArray_Hashes(_ua *C.coin__UxArray, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	__arg0 := ua.Hashes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_UxArray_HasDupes
func SKY_coin_UxArray_HasDupes(_ua *C.coin__UxArray, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	__arg0 := ua.HasDupes()
	*_arg0 = __arg0
	return
}

//export SKY_coin_UxArray_Sort
func SKY_coin_UxArray_Sort(_ua *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	ua.Sort()
	return
}

//export SKY_coin_UxArray_Len
func SKY_coin_UxArray_Len(_ua *C.coin__UxArray, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	__arg0 := ua.Len()
	*_arg0 = __arg0
	return
}

//export SKY_coin_UxArray_Less
func SKY_coin_UxArray_Less(_ua *C.coin__UxArray, _i, _j int, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	i := _i
	j := _j
	__arg1 := ua.Less(i, j)
	*_arg1 = __arg1
	return
}

//export SKY_coin_UxArray_Swap
func SKY_coin_UxArray_Swap(_ua *C.coin__UxArray, _i, _j int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	i := _i
	j := _j
	ua.Swap(i, j)
	return
}

//export SKY_coin_UxArray_Coins
func SKY_coin_UxArray_Coins(_ua *C.coin__UxArray, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	__arg0, ____return_err := ua.Coins()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

//export SKY_coin_UxArray_CoinHours
func SKY_coin_UxArray_CoinHours(_ua *C.coin__UxArray, _headTime uint64, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	headTime := _headTime
	__arg1, ____return_err := ua.CoinHours(headTime)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

//export SKY_coin_UxArray_Sub
func SKY_coin_UxArray_Sub(_ua *C.coin__UxArray, _other *C.coin__UxArray, _arg1 *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	other := *(*coin.UxArray)(unsafe.Pointer(_other))
	__arg1 := ua.Sub(other)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_coin_UxArray_Add
func SKY_coin_UxArray_Add(_ua *C.coin__UxArray, _other *C.coin__UxArray, _arg1 *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	other := *(*coin.UxArray)(unsafe.Pointer(_other))
	__arg1 := ua.Add(other)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}
