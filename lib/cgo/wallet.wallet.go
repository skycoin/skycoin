package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	blockdb "github.com/skycoin/skycoin/src/visor/blockdb"
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

//export SKY_wallet_NewWallet
func SKY_wallet_NewWallet(_wltName string, _opts *C.Options__Handle, _arg2 *C.Wallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	wltName := _wltName
	__opts, okopts := lookupOptionsHandle(*_opts)
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
func SKY_wallet_Wallet_Save(_w *C.Wallet__Handle, _dir string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
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
func SKY_wallet_Wallet_Validate(_w *C.Wallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
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
func SKY_wallet_Wallet_Type(_w *C.Wallet__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.Type()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_Wallet_Version
func SKY_wallet_Wallet_Version(_w *C.Wallet__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.Version()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_Wallet_Filename
func SKY_wallet_Wallet_Filename(_w *C.Wallet__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.Filename()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_Wallet_Label
func SKY_wallet_Wallet_Label(_w *C.Wallet__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.Label()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_Wallet_IsEncrypted
func SKY_wallet_Wallet_IsEncrypted(_w *C.Wallet__Handle, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.IsEncrypted()
	*_arg0 = __arg0
	return
}

//export SKY_wallet_Wallet_GenerateAddresses
func SKY_wallet_Wallet_GenerateAddresses(_w *C.Wallet__Handle, _num uint64, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
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

//export SKY_wallet_Wallet_ScanAddresses
func SKY_wallet_Wallet_ScanAddresses(_w *C.Wallet__Handle, _scanN uint64, _bg *C.wallet__BalanceGetter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	scanN := _scanN
	bg := *(*wallet.BalanceGetter)(unsafe.Pointer(_bg))
	____return_err := w.ScanAddresses(scanN, bg)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_Wallet_GetAddresses
func SKY_wallet_Wallet_GetAddresses(_w *C.Wallet__Handle, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := w.GetAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_wallet_Wallet_GetEntry
func SKY_wallet_Wallet_GetEntry(_w *C.Wallet__Handle, _a *C.cipher__Address, _arg1 *C.wallet__Entry, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
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
func SKY_wallet_Wallet_AddEntry(_w *C.Wallet__Handle, _entry *C.wallet__Entry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
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

//export SKY_wallet_Wallet_CreateAndSignTransaction
func SKY_wallet_Wallet_CreateAndSignTransaction(_w *C.Wallet__Handle, _vld *C.wallet__Validator, _unspent *C.blockdb__UnspentGetter, _headTime, _coins uint64, _dest *C.cipher__Address, _arg4 *C.coin__Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(*_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	vld := *(*wallet.Validator)(unsafe.Pointer(_vld))
	unspent := *(*blockdb.UnspentGetter)(unsafe.Pointer(_unspent))
	headTime := _headTime
	coins := _coins
	dest := *(*cipher.Address)(unsafe.Pointer(_dest))
	__arg4, ____return_err := w.CreateAndSignTransaction(vld, unspent, headTime, coins, dest)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg4 = *(*C.coin__Transaction)(unsafe.Pointer(__arg4))
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

//export SKY_wallet_ChooseSpendsMinimizeUxOuts
func SKY_wallet_ChooseSpendsMinimizeUxOuts(_uxa []C.wallet__UxBalance, _coins uint64, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uxa := *(*[]wallet.UxBalance)(unsafe.Pointer(&_uxa))
	coins := _coins
	__arg2, ____return_err := wallet.ChooseSpendsMinimizeUxOuts(uxa, coins)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

//export SKY_wallet_ChooseSpendsMaximizeUxOuts
func SKY_wallet_ChooseSpendsMaximizeUxOuts(_uxa []C.wallet__UxBalance, _coins uint64, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uxa := *(*[]wallet.UxBalance)(unsafe.Pointer(&_uxa))
	coins := _coins
	__arg2, ____return_err := wallet.ChooseSpendsMaximizeUxOuts(uxa, coins)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}
