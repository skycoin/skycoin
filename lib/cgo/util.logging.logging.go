package main

import (
	logging "github.com/skycoin/skycoin/src/logging"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_logging_MustGetLogger
func SKY_logging_MustGetLogger(_module string, _arg1 *C.Logger) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	module := _module
	__arg1 := logging.MustGetLogger(module)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofLogger))
	return
}

// export SKY_logging_Disable
func SKY_logging_Disable() (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	logging.Disable()
	return
}

// export SKY_logging_RedirectTo
func SKY_logging_RedirectTo(_w *C.Writer) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	logging.RedirectTo(w)
	return
}
