package main

import (
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

// export SKY_wallet_CryptoTypeFromString
func SKY_wallet_CryptoTypeFromString(_s string, _arg1 *C.CryptoType) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := _s
	__arg1, ____return_err := wallet.CryptoTypeFromString(s)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofCryptoType))
	}
	return
}
