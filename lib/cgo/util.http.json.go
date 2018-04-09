package main

import http "github.com/skycoin/skycoin/src/util/http"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_httphelper_SendJSON
func SKY_httphelper_SendJSON(_w *C.ResponseWriter, _m interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	____return_err := http.SendJSON(w, m)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_httphelper_SendJSONOr500
func SKY_httphelper_SendJSONOr500(_log *C.Logger, _w *C.ResponseWriter, _m interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	http.SendJSONOr500(log, w, m)
	return
}
