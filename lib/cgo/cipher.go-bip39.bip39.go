package main

import (
	gobip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_bip39_NewDefaultMnemomic
func SKY_bip39_NewDefaultMnemomic(_arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	__arg0, ____return_err := gobip39.NewDefaultMnemomic()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg0, _arg0)
	}
	return
}

// export SKY_bip39_NewEntropy
func SKY_bip39_NewEntropy(_bitSize int, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	bitSize := _bitSize
	__arg1, ____return_err := gobip39.NewEntropy(bitSize)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_bip39_NewMnemonic
func SKY_bip39_NewMnemonic(_entropy *C.GoSlice_, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	entropy := *(*[]byte)(unsafe.Pointer(_entropy))
	__arg1, ____return_err := gobip39.NewMnemonic(entropy)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_bip39_MnemonicToByteArray
func SKY_bip39_MnemonicToByteArray(_mnemonic string, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	mnemonic := _mnemonic
	__arg1, ____return_err := gobip39.MnemonicToByteArray(mnemonic)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_bip39_IsMnemonicValid
func SKY_bip39_IsMnemonicValid(_mnemonic string, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	mnemonic := _mnemonic
	__arg1 := gobip39.IsMnemonicValid(mnemonic)
	*_arg1 = __arg1
	return
}
