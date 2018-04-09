package main

import http "github.com/skycoin/skycoin/src/util/http"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_httphelper_HTTPError
func SKY_httphelper_HTTPError(_w *C.ResponseWriter, _status int, _httpMsg string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	status := _status
	httpMsg := _httpMsg
	http.HTTPError(w, status, httpMsg)
	return
}

// export SKY_httphelper_Error400
func SKY_httphelper_Error400(_w *C.ResponseWriter, _msg string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := _msg
	http.Error400(w, msg)
	return
}

// export SKY_httphelper_Error403
func SKY_httphelper_Error403(_w *C.ResponseWriter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	http.Error403(w)
	return
}

// export SKY_httphelper_Error403Msg
func SKY_httphelper_Error403Msg(_w *C.ResponseWriter, _msg string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := _msg
	http.Error403Msg(w, msg)
	return
}

// export SKY_httphelper_Error404
func SKY_httphelper_Error404(_w *C.ResponseWriter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	http.Error404(w)
	return
}

// export SKY_httphelper_Error405
func SKY_httphelper_Error405(_w *C.ResponseWriter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	http.Error405(w)
	return
}

// export SKY_httphelper_Error501
func SKY_httphelper_Error501(_w *C.ResponseWriter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	http.Error501(w)
	return
}

// export SKY_httphelper_Error500
func SKY_httphelper_Error500(_w *C.ResponseWriter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	http.Error500(w)
	return
}

// export SKY_httphelper_Error500Msg
func SKY_httphelper_Error500Msg(_w *C.ResponseWriter, _msg string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := _msg
	http.Error500Msg(w, msg)
	return
}
