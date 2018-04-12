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

// export SKY_wallet_NewReadableEntry
func SKY_wallet_NewReadableEntry(_w *C.wallet__Entry, _arg1 *C.wallet__ReadableEntry) (____error_code uint32) {
	____error_code = 0
	w := *(*wallet.Entry)(unsafe.Pointer(_w))
	__arg1 := wallet.NewReadableEntry(w)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableEntry))
	return
}

// export SKY_wallet_LoadReadableEntry
func SKY_wallet_LoadReadableEntry(_filename string, _arg1 *C.wallet__ReadableEntry) (____error_code uint32) {
	____error_code = 0
	filename := _filename
	__arg1, ____return_err := wallet.LoadReadableEntry(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableEntry))
	}
	return
}

// export SKY_wallet_NewReadableEntryFromPubkey
func SKY_wallet_NewReadableEntryFromPubkey(_pub string, _arg1 *C.wallet__ReadableEntry) (____error_code uint32) {
	____error_code = 0
	pub := _pub
	__arg1 := wallet.NewReadableEntryFromPubkey(pub)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableEntry))
	return
}

// export SKY_wallet_ReadableEntry_Save
func SKY_wallet_ReadableEntry_Save(_re *C.wallet__ReadableEntry, _filename string) (____error_code uint32) {
	____error_code = 0
	re := (*wallet.ReadableEntry)(unsafe.Pointer(_re))
	filename := _filename
	____return_err := re.Save(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_NewReadableWallet
func SKY_wallet_NewReadableWallet(_w *C.wallet__Wallet, _arg1 *C.wallet__ReadableWallet) (____error_code uint32) {
	____error_code = 0
	w := (*wallet.Wallet)(unsafe.Pointer(_w))
	__arg1 := wallet.NewReadableWallet(w)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableWallet))
	return
}

// export SKY_wallet_LoadReadableWallet
func SKY_wallet_LoadReadableWallet(_filename string, _arg1 *C.wallet__ReadableWallet) (____error_code uint32) {
	____error_code = 0
	filename := _filename
	__arg1, ____return_err := wallet.LoadReadableWallet(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableWallet))
	}
	return
}

// export SKY_wallet_ReadableWallet_ToWallet
func SKY_wallet_ReadableWallet_ToWallet(_rw *C.wallet__ReadableWallet, _arg0 *C.wallet__Wallet) (____error_code uint32) {
	____error_code = 0
	rw := (*wallet.ReadableWallet)(unsafe.Pointer(_rw))
	__arg0, ____return_err := rw.ToWallet()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofWallet))
	}
	return
}

// export SKY_wallet_ReadableWallet_Save
func SKY_wallet_ReadableWallet_Save(_rw *C.wallet__ReadableWallet, _filename string) (____error_code uint32) {
	____error_code = 0
	rw := (*wallet.ReadableWallet)(unsafe.Pointer(_rw))
	filename := _filename
	____return_err := rw.Save(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_ReadableWallet_Load
func SKY_wallet_ReadableWallet_Load(_rw *C.wallet__ReadableWallet, _filename string) (____error_code uint32) {
	____error_code = 0
	rw := (*wallet.ReadableWallet)(unsafe.Pointer(_rw))
	filename := _filename
	____return_err := rw.Load(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
