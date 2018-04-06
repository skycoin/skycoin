package main

import (
	bip39 "github.com/skycoin/skycoin/src/bip39"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_bip39_NewDefaultMnemomic
func SKY_bip39_NewDefaultMnemomic(_arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0, ____return_err := bip39.NewDefaultMnemomic()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg0, _arg0)
	}
	return
}

// export SKY_bip39_NewEntropy
func SKY_bip39_NewEntropy(_bitSize int, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bitSize := _bitSize
	__arg1, ____return_err := bip39.NewEntropy(bitSize)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_bip39_NewMnemonic
func SKY_bip39_NewMnemonic(_entropy *C.GoSlice_, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	entropy := *(*[]byte)(unsafe.Pointer(_entropy))
	__arg1, ____return_err := bip39.NewMnemonic(entropy)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_bip39_MnemonicToByteArray
func SKY_bip39_MnemonicToByteArray(_mnemonic string, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	mnemonic := _mnemonic
	__arg1, ____return_err := bip39.MnemonicToByteArray(mnemonic)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_bip39_IsMnemonicValid
func SKY_bip39_IsMnemonicValid(_mnemonic string, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	mnemonic := _mnemonic
	__arg1 := bip39.IsMnemonicValid(mnemonic)
	*_arg1 = __arg1
	return
}
