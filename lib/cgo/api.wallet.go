package main

import api "github.com/skycoin/skycoin/src/api"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_api_NewWalletResponse
func SKY_api_NewWalletResponse(_w C.Wallet__Handle, _arg1 *C.WalletResponse__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, okw := lookupWalletHandle(_w)
	if !okw {
		____error_code = SKY_ERROR
		return
	}
	__arg1, ____return_err := api.NewWalletResponse(w)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerWalletResponseHandle(__arg1)
	}
	return
}
