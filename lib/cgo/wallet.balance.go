package main

import (
	"unsafe"

	coin "github.com/skycoin/skycoin/src/coin"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_wallet_NewBalance
func SKY_wallet_NewBalance(_coins, _hours uint64, _arg1 *C.wallet__Balance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	coins := _coins
	hours := _hours
	__arg1 := wallet.NewBalance(coins, hours)
	*_arg1 = *(*C.wallet__Balance)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_wallet_NewBalanceFromUxOut
func SKY_wallet_NewBalanceFromUxOut(_headTime uint64, _ux *C.coin__UxOut, _arg2 *C.wallet__Balance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	headTime := _headTime
	ux := (*coin.UxOut)(unsafe.Pointer(_ux))
	__arg2, ____return_err := wallet.NewBalanceFromUxOut(headTime, ux)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = *(*C.wallet__Balance)(unsafe.Pointer(&__arg2))
	}
	return
}

//export SKY_wallet_Balance_Add
func SKY_wallet_Balance_Add(_bal *C.wallet__Balance, _other *C.wallet__Balance, _arg1 *C.wallet__Balance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bal := *(*wallet.Balance)(unsafe.Pointer(_bal))
	other := *(*wallet.Balance)(unsafe.Pointer(_other))
	__arg1, ____return_err := bal.Add(other)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.wallet__Balance)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_wallet_Balance_Sub
func SKY_wallet_Balance_Sub(_bal *C.wallet__Balance, _other *C.wallet__Balance, _arg1 *C.wallet__Balance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bal := *(*wallet.Balance)(unsafe.Pointer(_bal))
	other := *(*wallet.Balance)(unsafe.Pointer(_other))
	__arg1 := bal.Sub(other)
	*_arg1 = *(*C.wallet__Balance)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_wallet_Balance_Equals
func SKY_wallet_Balance_Equals(_bal *C.wallet__Balance, _other *C.wallet__Balance, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bal := *(*wallet.Balance)(unsafe.Pointer(_bal))
	other := *(*wallet.Balance)(unsafe.Pointer(_other))
	__arg1 := bal.Equals(other)
	*_arg1 = __arg1
	return
}

//export SKY_wallet_Balance_IsZero
func SKY_wallet_Balance_IsZero(_bal *C.wallet__Balance, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bal := *(*wallet.Balance)(unsafe.Pointer(_bal))
	__arg0 := bal.IsZero()
	*_arg0 = __arg0
	return
}
