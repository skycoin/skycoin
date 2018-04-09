package main

import (
	http "github.com/skycoin/skycoin/src/util/http"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_httphelper_ElapsedHandler
func SKY_httphelper_ElapsedHandler(_logger *C.FieldLogger, _handler *C.Handler, _arg2 *C.Handler) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg2 := http.ElapsedHandler(logger, handler)
	return
}

// export SKY_httphelper_wrappedResponseWriter_WriteHeader
func SKY_httphelper_wrappedResponseWriter_WriteHeader(_lrw wrappedResponseWriter, _code int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	lrw := (*wrappedResponseWriter)(unsafe.Pointer(_lrw))
	code := _code
	lrw.WriteHeader(code)
	return
}
