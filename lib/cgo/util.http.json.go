package main

import httphelper "github.com/skycoin/skycoin/src/httphelper"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_httphelper_SendJSON
func SKY_httphelper_SendJSON(_w *C.ResponseWriter, _m interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	____return_err := httphelper.SendJSON(w, m)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_httphelper_SendJSONOr500
func SKY_httphelper_SendJSONOr500(_log *C.Logger, _w *C.ResponseWriter, _m interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	httphelper.SendJSONOr500(log, w, m)
	return
}
