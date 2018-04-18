package main

import (
	wallet "github.com/skycoin/skycoin/src/wallet"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_wallet_Entry_Verify
func SKY_wallet_Entry_Verify(_we *C.wallet__Entry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	we := (*wallet.Entry)(unsafe.Pointer(_we))
	____return_err := we.Verify()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_Entry_VerifyPublic
func SKY_wallet_Entry_VerifyPublic(_we *C.wallet__Entry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	we := (*wallet.Entry)(unsafe.Pointer(_we))
	____return_err := we.VerifyPublic()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
