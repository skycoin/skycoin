package main

import (
	logging "github.com/skycoin/skycoin/src/util/logging"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_logging_MustGetLogger
func SKY_logging_MustGetLogger(_module string, _arg1 *C.Logger) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	module := _module
	__arg1 := logging.MustGetLogger(module)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofLogger))
	return
}

// export SKY_logging_Disable
func SKY_logging_Disable() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logging.Disable()
	return
}

// export SKY_logging_RedirectTo
func SKY_logging_RedirectTo(_w *C.Writer) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logging.RedirectTo(w)
	return
}
