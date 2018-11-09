package main

import (
	"reflect"
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher"
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
	uo := (*coin.UxOut)(unsafe.Pointer(_uo))
	__arg0 := uo.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_UxOut_SnapshotHash
func SKY_coin_UxOut_SnapshotHash(_uo *C.coin__UxOut, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	uo := (*coin.UxOut)(unsafe.Pointer(_uo))
	__arg0 := uo.SnapshotHash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_UxBody_Hash
func SKY_coin_UxBody_Hash(_ub *C.coin__UxBody, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	ub := (*coin.UxBody)(unsafe.Pointer(_ub))
	__arg0 := ub.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_UxOut_CoinHours
func SKY_coin_UxOut_CoinHours(_uo *C.coin__UxOut, _t uint64, _arg1 *uint64) (____error_code uint32) {
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
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	__arg0 := ua.Hashes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_UxArray_HasDupes
func SKY_coin_UxArray_HasDupes(_ua *C.coin__UxArray, _arg0 *bool) (____error_code uint32) {
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	*_arg0 = ua.HasDupes()
	return
}

//export SKY_coin_UxArray_Sort
func SKY_coin_UxArray_Sort(_ua *C.coin__UxArray) (____error_code uint32) {
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	ua.Sort()
	return
}

//export SKY_coin_UxArray_Len
func SKY_coin_UxArray_Len(_ua *C.coin__UxArray, _arg0 *int) (____error_code uint32) {
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	*_arg0 = ua.Len()
	return
}

//export SKY_coin_UxArray_Less
func SKY_coin_UxArray_Less(_ua *C.coin__UxArray, _i, _j int, _arg0 *bool) (____error_code uint32) {
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	*_arg0 = ua.Less(_i, _j)
	return
}

//export SKY_coin_UxArray_Swap
func SKY_coin_UxArray_Swap(_ua *C.coin__UxArray, _i, _j int) (____error_code uint32) {
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	ua.Swap(_i, _j)
	return
}

//export SKY_coin_UxArray_Coins
func SKY_coin_UxArray_Coins(_ua *C.coin__UxArray, _arg0 *uint64) (____error_code uint32) {
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
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	other := *(*coin.UxArray)(unsafe.Pointer(_other))
	__arg1 := ua.Sub(other)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_coin_UxArray_Add
func SKY_coin_UxArray_Add(_ua *C.coin__UxArray, _other *C.coin__UxArray, _arg1 *C.coin__UxArray) (____error_code uint32) {
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	other := *(*coin.UxArray)(unsafe.Pointer(_other))
	__arg1 := ua.Add(other)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_coin_NewAddressUxOuts
func SKY_coin_NewAddressUxOuts(_ua *C.coin__UxArray, _address_outs *C.AddressUxOuts_Handle) (____error_code uint32) {
	ua := *(*coin.UxArray)(unsafe.Pointer(_ua))
	address_outs := coin.NewAddressUxOuts(ua)
	*_address_outs = registerAddressUxOutHandle(&address_outs)
	return
}

//export SKY_coin_AddressUxOuts_Keys
func SKY_coin_AddressUxOuts_Keys(_address_outs C.AddressUxOuts_Handle, _keys *C.GoSlice_) (____error_code uint32) {
	address_outs, ok := lookupAddressUxOutHandle(_address_outs)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	keys := (*address_outs).Keys()
	copyToGoSlice(reflect.ValueOf(keys), _keys)
	return
}

//export SKY_coin_AddressUxOuts_Flatten
func SKY_coin_AddressUxOuts_Flatten(_address_outs C.AddressUxOuts_Handle, _ua *C.coin__UxArray) (____error_code uint32) {
	address_outs, ok := lookupAddressUxOutHandle(_address_outs)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	ux := (*address_outs).Flatten()
	copyToGoSlice(reflect.ValueOf(ux), _ua)
	return
}

//export SKY_coin_AddressUxOuts_Sub
func SKY_coin_AddressUxOuts_Sub(_auo1 C.AddressUxOuts_Handle, _auo2 C.AddressUxOuts_Handle, _auo_result *C.AddressUxOuts_Handle) (____error_code uint32) {
	auo1, ok := lookupAddressUxOutHandle(_auo1)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	auo2, ok := lookupAddressUxOutHandle(_auo2)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	auo_result := (*auo1).Sub(*auo2)
	*_auo_result = registerAddressUxOutHandle(&auo_result)
	return
}

//export SKY_coin_AddressUxOuts_Add
func SKY_coin_AddressUxOuts_Add(_auo1 C.AddressUxOuts_Handle, _auo2 C.AddressUxOuts_Handle, _auo_result *C.AddressUxOuts_Handle) (____error_code uint32) {
	auo1, ok := lookupAddressUxOutHandle(_auo1)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	auo2, ok := lookupAddressUxOutHandle(_auo2)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	auo_result := (*auo1).Add(*auo2)
	*_auo_result = registerAddressUxOutHandle(&auo_result)
	return
}

//export SKY_coin_AddressUxOuts_Get
func SKY_coin_AddressUxOuts_Get(handle C.AddressUxOuts_Handle, _key *C.cipher__Address, _uxOuts *C.coin__UxArray) (____error_code uint32) {
	a, ok := lookupAddressUxOutHandle(handle)
	if ok {
		key := *(*cipher.Address)(unsafe.Pointer(_key))
		uxOuts, found := (*a)[key]
		if found {
			copyToGoSlice(reflect.ValueOf(uxOuts), _uxOuts)
			____error_code = SKY_OK
		}
	} else {
		____error_code = SKY_BAD_HANDLE
	}
	return
}

//export SKY_coin_AddressUxOuts_HasKey
func SKY_coin_AddressUxOuts_HasKey(handle C.AddressUxOuts_Handle, _key *C.cipher__Address, _hasKey *bool) (____error_code uint32) {
	a, ok := lookupAddressUxOutHandle(handle)
	if ok {
		key := *(*cipher.Address)(unsafe.Pointer(_key))
		_, found := (*a)[key]
		*_hasKey = found
		____error_code = SKY_OK
	} else {
		____error_code = SKY_BAD_HANDLE
	}
	return
}

//export SKY_coin_AddressUxOuts_GetOutputLength
func SKY_coin_AddressUxOuts_GetOutputLength(handle C.AddressUxOuts_Handle, _key *C.cipher__Address, _length *int) (____error_code uint32) {
	a, ok := lookupAddressUxOutHandle(handle)
	if ok {
		key := *(*cipher.Address)(unsafe.Pointer(_key))
		uxOuts, found := (*a)[key]
		if found {
			*_length = len(uxOuts)
			____error_code = SKY_OK
		}
	} else {
		____error_code = SKY_BAD_HANDLE
	}
	return
}

//export SKY_coin_AddressUxOuts_Length
func SKY_coin_AddressUxOuts_Length(handle C.AddressUxOuts_Handle, _length *int) (____error_code uint32) {
	a, ok := lookupAddressUxOutHandle(handle)
	if ok {
		*_length = len(*a)
		____error_code = SKY_OK
	} else {
		____error_code = SKY_BAD_HANDLE
	}
	return
}

//export SKY_coin_AddressUxOuts_Set
func SKY_coin_AddressUxOuts_Set(handle C.AddressUxOuts_Handle, _key *C.cipher__Address, _uxOuts *C.coin__UxArray) (____error_code uint32) {
	a, ok := lookupAddressUxOutHandle(handle)
	if ok {
		key := *(*cipher.Address)(unsafe.Pointer(_key))
		//Copy the slice because it is going to be kept
		//We can't hold memory allocated outside Go
		tempUxOuts := *(*coin.UxArray)(unsafe.Pointer(_uxOuts))
		uxOuts := make(coin.UxArray, 0, len(tempUxOuts))
		for _, ux := range tempUxOuts {
			uxOuts = append(uxOuts, ux)
		}
		(*a)[key] = uxOuts
		____error_code = SKY_OK
	} else {
		____error_code = SKY_BAD_HANDLE
	}
	return
}
