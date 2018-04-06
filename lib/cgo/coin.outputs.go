package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_coin_UxOut_Hash
func SKY_coin_UxOut_Hash(_uo *C.UxOut, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	uo := (*cipher.UxOut)(unsafe.Pointer(_uo))
	__arg0 := uo.Hash()
	return
}

// export SKY_coin_UxOut_SnapshotHash
func SKY_coin_UxOut_SnapshotHash(_uo *C.UxOut, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	uo := (*cipher.UxOut)(unsafe.Pointer(_uo))
	__arg0 := uo.SnapshotHash()
	return
}

// export SKY_coin_UxBody_Hash
func SKY_coin_UxBody_Hash(_ub *C.UxBody, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ub := (*cipher.UxBody)(unsafe.Pointer(_ub))
	__arg0 := ub.Hash()
	return
}

// export SKY_coin_UxOut_CoinHours
func SKY_coin_UxOut_CoinHours(_uo *C.UxOut, _t uint64, _arg1 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	uo := (*cipher.UxOut)(unsafe.Pointer(_uo))
	t := _t
	__arg1, ____return_err := uo.CoinHours(t)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_coin_UxArray_Hashes
func SKY_coin_UxArray_Hashes(_ua *C.UxArray, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	__arg0 := ua.Hashes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_UxArray_HasDupes
func SKY_coin_UxArray_HasDupes(_ua *C.UxArray, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	__arg0 := ua.HasDupes()
	*_arg0 = __arg0
	return
}

// export SKY_coin_UxArray_Set
func SKY_coin_UxArray_Set(_ua *C.UxArray, _arg0 *C.UxHashSet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	__arg0 := ua.Set()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofUxHashSet))
	return
}

// export SKY_coin_UxArray_Sort
func SKY_coin_UxArray_Sort(_ua *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	ua.Sort()
	return
}

// export SKY_coin_UxArray_Len
func SKY_coin_UxArray_Len(_ua *C.UxArray, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	__arg0 := ua.Len()
	*_arg0 = __arg0
	return
}

// export SKY_coin_UxArray_Less
func SKY_coin_UxArray_Less(_ua *C.UxArray, _i, _j int, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	i := _i
	j := _j
	__arg1 := ua.Less(i, j)
	*_arg1 = __arg1
	return
}

// export SKY_coin_UxArray_Swap
func SKY_coin_UxArray_Swap(_ua *C.UxArray, _i, _j int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	i := _i
	j := _j
	ua.Swap(i, j)
	return
}

// export SKY_coin_UxArray_Coins
func SKY_coin_UxArray_Coins(_ua *C.UxArray, _arg0 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	__arg0, ____return_err := ua.Coins()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

// export SKY_coin_UxArray_CoinHours
func SKY_coin_UxArray_CoinHours(_ua *C.UxArray, _headTime uint64, _arg1 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	headTime := _headTime
	__arg1, ____return_err := ua.CoinHours(headTime)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_coin_NewAddressUxOuts
func SKY_coin_NewAddressUxOuts(_uxs *C.UxArray, _arg1 *C.AddressUxOuts) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	uxs := *(*cipher.UxArray)(unsafe.Pointer(_uxs))
	__arg1 := coin.NewAddressUxOuts(uxs)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofAddressUxOuts))
	return
}

// export SKY_coin_AddressUxOuts_Keys
func SKY_coin_AddressUxOuts_Keys(_auo *C.AddressUxOuts, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	auo := *(*cipher.AddressUxOuts)(unsafe.Pointer(_auo))
	__arg0 := auo.Keys()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_AddressUxOuts_Flatten
func SKY_coin_AddressUxOuts_Flatten(_auo *C.AddressUxOuts, _arg0 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	auo := *(*cipher.AddressUxOuts)(unsafe.Pointer(_auo))
	__arg0 := auo.Flatten()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofUxArray))
	return
}

// export SKY_coin_AddressUxOuts_Sub
func SKY_coin_AddressUxOuts_Sub(_auo *C.AddressUxOuts, _other *C.AddressUxOuts, _arg1 *C.AddressUxOuts) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	auo := *(*cipher.AddressUxOuts)(unsafe.Pointer(_auo))
	other := *(*cipher.AddressUxOuts)(unsafe.Pointer(_other))
	__arg1 := auo.Sub(other)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofAddressUxOuts))
	return
}

// export SKY_coin_AddressUxOuts_Add
func SKY_coin_AddressUxOuts_Add(_auo *C.AddressUxOuts, _other *C.AddressUxOuts, _arg1 *C.AddressUxOuts) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	auo := *(*cipher.AddressUxOuts)(unsafe.Pointer(_auo))
	other := *(*cipher.AddressUxOuts)(unsafe.Pointer(_other))
	__arg1 := auo.Add(other)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofAddressUxOuts))
	return
}

// export SKY_coin_UxArray_Sub
func SKY_coin_UxArray_Sub(_ua *C.UxArray, _other *C.UxArray, _arg1 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	other := *(*cipher.UxArray)(unsafe.Pointer(_other))
	__arg1 := ua.Sub(other)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofUxArray))
	return
}

// export SKY_coin_UxArray_Add
func SKY_coin_UxArray_Add(_ua *C.UxArray, _other *C.UxArray, _arg1 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ua := *(*cipher.UxArray)(unsafe.Pointer(_ua))
	other := *(*cipher.UxArray)(unsafe.Pointer(_other))
	__arg1 := ua.Add(other)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofUxArray))
	return
}
