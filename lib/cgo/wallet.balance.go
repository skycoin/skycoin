package main

import (
	wallet "github.com/skycoin/skycoin/src/wallet"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_wallet_NewBalance
func SKY_wallet_NewBalance(_coins, _hours uint64, _arg1 *C.Balance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	coins := _coins
	hours := _hours
	__arg1 := wallet.NewBalance(coins, hours)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofBalance))
	return
}

// export SKY_wallet_NewBalanceFromUxOut
func SKY_wallet_NewBalanceFromUxOut(_headTime uint64, _ux *C.UxOut, _arg2 *C.Balance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	headTime := _headTime
	__arg2, ____return_err := wallet.NewBalanceFromUxOut(headTime, ux)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofBalance))
	}
	return
}

// export SKY_wallet_Balance_Add
func SKY_wallet_Balance_Add(_bal *C.Balance, _other *C.Balance, _arg1 *C.Balance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bal := *(*Balance)(unsafe.Pointer(_bal))
	other := *(*Balance)(unsafe.Pointer(_other))
	__arg1, ____return_err := bal.Add(other)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofBalance))
	}
	return
}

// export SKY_wallet_Balance_Sub
func SKY_wallet_Balance_Sub(_bal *C.Balance, _other *C.Balance, _arg1 *C.Balance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bal := *(*Balance)(unsafe.Pointer(_bal))
	other := *(*Balance)(unsafe.Pointer(_other))
	__arg1 := bal.Sub(other)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofBalance))
	return
}

// export SKY_wallet_Balance_Equals
func SKY_wallet_Balance_Equals(_bal *C.Balance, _other *C.Balance, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bal := *(*Balance)(unsafe.Pointer(_bal))
	other := *(*Balance)(unsafe.Pointer(_other))
	__arg1 := bal.Equals(other)
	*_arg1 = __arg1
	return
}

// export SKY_wallet_Balance_IsZero
func SKY_wallet_Balance_IsZero(_bal *C.Balance, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bal := *(*Balance)(unsafe.Pointer(_bal))
	__arg0 := bal.IsZero()
	*_arg0 = __arg0
	return
}
