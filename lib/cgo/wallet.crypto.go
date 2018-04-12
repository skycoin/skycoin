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

// export SKY_wallet_CryptoTypeFromString
func SKY_wallet_CryptoTypeFromString(_s string, _arg1 *C.wallet__CryptoType) (____error_code uint32) {
	____error_code = 0
	s := _s
	__arg1, ____return_err := wallet.CryptoTypeFromString(s)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofCryptoType))
	}
	return
}
