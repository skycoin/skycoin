package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	strand "github.com/skycoin/skycoin/src/strand"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_strand_Strand
func SKY_strand_Strand(_logger *C.Logger, _c C.GoChan_, _name string, _f C.Handle, _quit C.GoChan_, _quitErr error) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	name := _name
	quitErr := *(*cipher.error)(unsafe.Pointer(_quitErr))
	____return_err := strand.Strand(logger, c, name, f, quit, quitErr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
