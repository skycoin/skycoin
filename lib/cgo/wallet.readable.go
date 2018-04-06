package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	wallet "github.com/skycoin/skycoin/src/wallet"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_wallet_NewReadableEntry
func SKY_wallet_NewReadableEntry(_w *C.Entry, _arg1 *C.ReadableEntry) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	w := *(*cipher.Entry)(unsafe.Pointer(_w))
	__arg1 := wallet.NewReadableEntry(w)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableEntry))
	return
}

// export SKY_wallet_LoadReadableEntry
func SKY_wallet_LoadReadableEntry(_filename string, _arg1 *C.ReadableEntry) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	filename := _filename
	__arg1, ____return_err := wallet.LoadReadableEntry(filename)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableEntry))
	}
	return
}

// export SKY_wallet_NewReadableEntryFromPubkey
func SKY_wallet_NewReadableEntryFromPubkey(_pub string, _arg1 *C.ReadableEntry) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pub := _pub
	__arg1 := wallet.NewReadableEntryFromPubkey(pub)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableEntry))
	return
}

// export SKY_wallet_ReadableEntry_Save
func SKY_wallet_ReadableEntry_Save(_re *C.ReadableEntry, _filename string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	re := (*cipher.ReadableEntry)(unsafe.Pointer(_re))
	filename := _filename
	____return_err := re.Save(filename)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_NewReadableWallet
func SKY_wallet_NewReadableWallet(_w *C.Wallet, _arg1 *C.ReadableWallet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	w := (*cipher.Wallet)(unsafe.Pointer(_w))
	__arg1 := wallet.NewReadableWallet(w)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableWallet))
	return
}

// export SKY_wallet_LoadReadableWallet
func SKY_wallet_LoadReadableWallet(_filename string, _arg1 *C.ReadableWallet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	filename := _filename
	__arg1, ____return_err := wallet.LoadReadableWallet(filename)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableWallet))
	}
	return
}

// export SKY_wallet_ReadableWallet_ToWallet
func SKY_wallet_ReadableWallet_ToWallet(_rw *C.ReadableWallet, _arg0 *C.Wallet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rw := (*cipher.ReadableWallet)(unsafe.Pointer(_rw))
	__arg0, ____return_err := rw.ToWallet()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofWallet))
	}
	return
}

// export SKY_wallet_ReadableWallet_Save
func SKY_wallet_ReadableWallet_Save(_rw *C.ReadableWallet, _filename string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rw := (*cipher.ReadableWallet)(unsafe.Pointer(_rw))
	filename := _filename
	____return_err := rw.Save(filename)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_ReadableWallet_Load
func SKY_wallet_ReadableWallet_Load(_rw *C.ReadableWallet, _filename string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rw := (*cipher.ReadableWallet)(unsafe.Pointer(_rw))
	filename := _filename
	____return_err := rw.Load(filename)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
