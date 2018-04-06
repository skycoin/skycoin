package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_wallet_Entry_Verify
func SKY_wallet_Entry_Verify(_we *C.Entry) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	we := (*cipher.Entry)(unsafe.Pointer(_we))
	____return_err := we.Verify()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_Entry_VerifyPublic
func SKY_wallet_Entry_VerifyPublic(_we *C.Entry) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	we := (*cipher.Entry)(unsafe.Pointer(_we))
	____return_err := we.VerifyPublic()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
