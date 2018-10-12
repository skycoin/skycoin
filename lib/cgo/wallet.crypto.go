package main

import (
	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_wallet_CryptoTypeFromString
func SKY_wallet_CryptoTypeFromString(_s string, _arg1 *C.GoString_) (____error_code uint32) {
	s := _s
	__arg1, ____return_err := wallet.CryptoTypeFromString(s)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(string(__arg1), _arg1)
	}
	return
}
