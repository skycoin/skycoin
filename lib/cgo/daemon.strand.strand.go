package main

import (
	strand "github.com/skycoin/skycoin/src/daemon/strand"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_strand_Strand
func SKY_strand_Strand(_logger *C.Logger, _c C.GoChan_, _name string, _f C.Handle, _quit C.GoChan_, _quitErr error) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	name := _name
	quitErr := *(*error)(unsafe.Pointer(_quitErr))
	____return_err := strand.Strand(logger, c, name, f, quit, quitErr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
