package main

import (
	"reflect"
	"unsafe"

	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_wallet_NewError
func SKY_wallet_NewError(_err error) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	err := _err
	____return_err := wallet.NewError(err)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_NewWallet
func SKY_wallet_NewWallet(_wltName string, _opts C.Options__Handle, _arg2 *C.Wallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	wltName := _wltName
	__opts, okopts := lookupOptionsHandle(_opts)
	if !okopts {
		____error_code = SKY_ERROR
		return
	}
	opts := *__opts
	__arg2, ____return_err := wallet.NewWallet(wltName, opts)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerWalletHandle(__arg2)
	}
	return
}

//export SKY_wallet_Wallet_Lock
func SKY_wallet_Wallet_Lock(_w C.Wallet__Handle, _password []byte, _cryptoType string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	password := *(*[]byte)(unsafe.Pointer(&_password))
	cryptoType := wallet.CryptoType(_cryptoType)
	____return_err := w.Lock(password, cryptoType)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_Wallet_Unlock
func SKY_wallet_Wallet_Unlock(_w C.Wallet__Handle, _password []byte, _arg1 *C.Wallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	password := *(*[]byte)(unsafe.Pointer(&_password))
	__arg1, ____return_err := w.Unlock(password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerWalletHandle(__arg1)
	}
	return
}

//export SKY_wallet_Load
func SKY_wallet_Load(_wltFile string, _arg1 *C.Wallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	wltFile := _wltFile
	__arg1, ____return_err := wallet.Load(wltFile)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerWalletHandle(__arg1)
	}
	return
}

//export SKY_wallet_Wallet_Save
func SKY_wallet_Wallet_Save(_w C.Wallet__Handle, _dir string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	dir := _dir
	____return_err := w.Save(dir)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_Wallet_Validate
func SKY_wallet_Wallet_Validate(_w C.Wallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	____return_err := w.Validate()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_Wallet_Type
func SKY_wallet_Wallet_Type(_w C.Wallet__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.Type()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_Wallet_Version
func SKY_wallet_Wallet_Version(_w C.Wallet__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.Version()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_Wallet_Filename
func SKY_wallet_Wallet_Filename(_w C.Wallet__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.Filename()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_Wallet_Label
func SKY_wallet_Wallet_Label(_w C.Wallet__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.Label()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_Wallet_IsEncrypted
func SKY_wallet_Wallet_IsEncrypted(_w C.Wallet__Handle, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.IsEncrypted()
	*_arg0 = __arg0
	return
}

//export SKY_wallet_Wallet_GenerateAddresses
func SKY_wallet_Wallet_GenerateAddresses(_w C.Wallet__Handle, _num uint64, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	num := _num
	__arg1, ____return_err := w.GenerateAddresses(num)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

//export SKY_wallet_Wallet_GetAddresses
func SKY_wallet_Wallet_GetAddresses(_w C.Wallet__Handle, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.GetAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_wallet_Wallet_GetEntry
func SKY_wallet_Wallet_GetEntry(_w C.Wallet__Handle, _a *C.cipher__Address, _arg1 *C.wallet__Entry, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	a := *(*cipher.Address)(unsafe.Pointer(_a))
	__arg1, __arg2 := w.GetEntry(a)
	*_arg1 = *(*C.wallet__Entry)(unsafe.Pointer(&__arg1))
	*_arg2 = __arg2
	return
}

//export SKY_wallet_Wallet_AddEntry
func SKY_wallet_Wallet_AddEntry(_w C.Wallet__Handle, _entry *C.wallet__Entry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	entry := *(*wallet.Entry)(unsafe.Pointer(_entry))
	____return_err := w.AddEntry(entry)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_DistributeSpendHours
func SKY_wallet_DistributeSpendHours(_inputHours, _nAddrs uint64, _haveChange bool, _arg2 *uint64, _arg3 *C.GoSlice_, _arg4 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	inputHours := _inputHours
	nAddrs := _nAddrs
	haveChange := _haveChange
	__arg2, __arg3, __arg4 := wallet.DistributeSpendHours(inputHours, nAddrs, haveChange)
	*_arg2 = __arg2
	copyToGoSlice(reflect.ValueOf(__arg3), _arg3)
	*_arg4 = __arg4
	return
}

//export SKY_wallet_DistributeCoinHoursProportional
func SKY_wallet_DistributeCoinHoursProportional(_coins []uint64, _hours uint64, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	coins := *(*[]uint64)(unsafe.Pointer(&_coins))
	hours := _hours
	__arg2, ____return_err := wallet.DistributeCoinHoursProportional(coins, hours)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

//export SKY_wallet_NewUxBalances
func SKY_wallet_NewUxBalances(_headTime uint64, _uxa *C.coin__UxArray, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	headTime := _headTime
	uxa := *(*coin.UxArray)(unsafe.Pointer(_uxa))
	__arg2, ____return_err := wallet.NewUxBalances(headTime, uxa)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

//export SKY_wallet_NewUxBalance
func SKY_wallet_NewUxBalance(_headTime uint64, _ux *C.coin__UxOut, _arg2 *C.wallet__UxBalance) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	headTime := _headTime
	ux := *(*coin.UxOut)(unsafe.Pointer(_ux))
	__arg2, ____return_err := wallet.NewUxBalance(headTime, ux)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = *(*C.wallet__UxBalance)(unsafe.Pointer(&__arg2))
	}
	return
}

//export SKY_wallet_ChooseSpendsMinimizeUxOuts
func SKY_wallet_ChooseSpendsMinimizeUxOuts(_uxa []C.wallet__UxBalance, _coins, _hours uint64, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uxa := *(*[]wallet.UxBalance)(unsafe.Pointer(&_uxa))
	coins := _coins
	hours := _hours
	__arg2, ____return_err := wallet.ChooseSpendsMinimizeUxOuts(uxa, coins, hours)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

//export SKY_wallet_ChooseSpendsMaximizeUxOuts
func SKY_wallet_ChooseSpendsMaximizeUxOuts(_uxa []C.wallet__UxBalance, _coins, _hours uint64, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uxa := *(*[]wallet.UxBalance)(unsafe.Pointer(&_uxa))
	coins := _coins
	hours := _hours
	__arg2, ____return_err := wallet.ChooseSpendsMaximizeUxOuts(uxa, coins, hours)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}
