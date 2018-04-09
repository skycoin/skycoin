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

// export SKY_wallet_NewService
func SKY_wallet_NewService(_c *C.Config, _arg1 *C.Service) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*Config)(unsafe.Pointer(_c))
	__arg1, ____return_err := wallet.NewService(c)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofService))
	}
	return
}

// export SKY_wallet_Service_CreateWallet
func SKY_wallet_Service_CreateWallet(_serv *C.Service, _wltName string, _options *C.Options, _arg2 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltName := _wltName
	options := *(*Options)(unsafe.Pointer(_options))
	__arg2, ____return_err := serv.CreateWallet(wltName, options)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofWallet))
	}
	return
}

// export SKY_wallet_Service_ScanAheadWalletAddresses
func SKY_wallet_Service_ScanAheadWalletAddresses(_serv *C.Service, _wltName string, _password *C.GoSlice_, _scanN uint64, _bg *C.BalanceGetter, _arg4 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltName := _wltName
	password := *(*[]byte)(unsafe.Pointer(_password))
	scanN := _scanN
	bg := *(*BalanceGetter)(unsafe.Pointer(_bg))
	__arg4, ____return_err := serv.ScanAheadWalletAddresses(wltName, password, scanN, bg)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg4)[:]), unsafe.Pointer(_arg4), uint(SizeofWallet))
	}
	return
}

// export SKY_wallet_Service_EncryptWallet
func SKY_wallet_Service_EncryptWallet(_serv *C.Service, _wltID string, _password *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltID := _wltID
	password := *(*[]byte)(unsafe.Pointer(_password))
	____return_err := serv.EncryptWallet(wltID, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_Service_DecryptWallet
func SKY_wallet_Service_DecryptWallet(_serv *C.Service, _wltID string, _password *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltID := _wltID
	password := *(*[]byte)(unsafe.Pointer(_password))
	____return_err := serv.DecryptWallet(wltID, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_Service_NewAddresses
func SKY_wallet_Service_NewAddresses(_serv *C.Service, _wltID string, _password *C.GoSlice_, _num uint64, _arg3 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltID := _wltID
	password := *(*[]byte)(unsafe.Pointer(_password))
	num := _num
	__arg3, ____return_err := serv.NewAddresses(wltID, password, num)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg3), _arg3)
	}
	return
}

// export SKY_wallet_Service_GetAddresses
func SKY_wallet_Service_GetAddresses(_serv *C.Service, _wltID string, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltID := _wltID
	__arg1, ____return_err := serv.GetAddresses(wltID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_wallet_Service_GetWallet
func SKY_wallet_Service_GetWallet(_serv *C.Service, _wltID string, _arg1 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltID := _wltID
	__arg1, ____return_err := serv.GetWallet(wltID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofWallet))
	}
	return
}

// export SKY_wallet_Service_GetWallets
func SKY_wallet_Service_GetWallets(_serv *C.Service, _arg0 *C.Wallets) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	__arg0 := serv.GetWallets()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofWallets))
	return
}

// export SKY_wallet_Service_ReloadWallets
func SKY_wallet_Service_ReloadWallets(_serv *C.Service) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	____return_err := serv.ReloadWallets()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_Service_CreateAndSignTransaction
func SKY_wallet_Service_CreateAndSignTransaction(_serv *C.Service, _wltID string, _password *C.GoSlice_, _vld *C.Validator, _unspent *C.UnspentGetter, _headTime, _coins uint64, _dest *C.Address, _arg6 *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltID := _wltID
	password := *(*[]byte)(unsafe.Pointer(_password))
	vld := *(*Validator)(unsafe.Pointer(_vld))
	headTime := _headTime
	coins := _coins
	__arg6, ____return_err := serv.CreateAndSignTransaction(wltID, password, vld, unspent, headTime, coins, dest)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_Service_UpdateWalletLabel
func SKY_wallet_Service_UpdateWalletLabel(_serv *C.Service, _wltID, _label string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltID := _wltID
	label := _label
	____return_err := serv.UpdateWalletLabel(wltID, label)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_Service_Remove
func SKY_wallet_Service_Remove(_serv *C.Service, _wltID string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	serv := (*Service)(unsafe.Pointer(_serv))
	wltID := _wltID
	serv.Remove(wltID)
	return
}
