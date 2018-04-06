package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	httphelper "github.com/skycoin/skycoin/src/httphelper"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_httphelper_ElapsedHandler
func SKY_httphelper_ElapsedHandler(_logger *C.FieldLogger, _handler *C.Handler, _arg2 *C.Handler) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg2 := httphelper.ElapsedHandler(logger, handler)
	return
}

// export SKY_httphelper_wrappedResponseWriter_WriteHeader
func SKY_httphelper_wrappedResponseWriter_WriteHeader(_lrw wrappedResponseWriter, _code int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	lrw := (*cipher.wrappedResponseWriter)(unsafe.Pointer(_lrw))
	code := _code
	lrw.WriteHeader(code)
	return
}
