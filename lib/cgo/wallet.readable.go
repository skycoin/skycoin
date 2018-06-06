package main

import (
	"unsafe"

	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_wallet_NewReadableEntry
func SKY_wallet_NewReadableEntry(_w *C.wallet__Entry, _arg1 *C.ReadableEntry__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w := *(*wallet.Entry)(unsafe.Pointer(_w))
	__arg1 := wallet.NewReadableEntry(w)
	*_arg1 = registerReadableEntryHandle(&__arg1)
	return
}

//export SKY_wallet_LoadReadableEntry
func SKY_wallet_LoadReadableEntry(_filename string, _arg1 *C.ReadableEntry__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	filename := _filename
	__arg1, ____return_err := wallet.LoadReadableEntry(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerReadableEntryHandle(&__arg1)
	}
	return
}

//export SKY_wallet_NewReadableEntryFromPubkey
func SKY_wallet_NewReadableEntryFromPubkey(_pub string, _arg1 *C.ReadableEntry__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pub := _pub
	__arg1 := wallet.NewReadableEntryFromPubkey(pub)
	*_arg1 = registerReadableEntryHandle(&__arg1)
	return
}

//export SKY_wallet_ReadableEntry_Save
func SKY_wallet_ReadableEntry_Save(_re C.ReadableEntry__Handle, _filename string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	re, okre := lookupReadableEntryHandle(_re)
	if !okre {
		____error_code = SKY_ERROR
		return
	}
	filename := _filename
	____return_err := re.Save(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_LoadReadableWallet
func SKY_wallet_LoadReadableWallet(_filename string, _arg1 *C.ReadableWallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	filename := _filename
	__arg1, ____return_err := wallet.LoadReadableWallet(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerReadableWalletHandle(__arg1)
	}
	return
}

//export SKY_wallet_ReadableWallet_Save
func SKY_wallet_ReadableWallet_Save(_rw C.ReadableWallet__Handle, _filename string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rw, okrw := lookupReadableWalletHandle(_rw)
	if !okrw {
		____error_code = SKY_ERROR
		return
	}
	filename := _filename
	____return_err := rw.Save(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_ReadableWallet_Load
func SKY_wallet_ReadableWallet_Load(_rw C.ReadableWallet__Handle, _filename string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rw, okrw := lookupReadableWalletHandle(_rw)
	if !okrw {
		____error_code = SKY_ERROR
		return
	}
	filename := _filename
	____return_err := rw.Load(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_ReadableWallet_Erase
func SKY_wallet_ReadableWallet_Erase(_rw C.ReadableWallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rw, okrw := lookupReadableWalletHandle(_rw)
	if !okrw {
		____error_code = SKY_ERROR
		return
	}
	rw.Erase()
	return
}
