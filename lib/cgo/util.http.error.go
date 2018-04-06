package main

import httphelper "github.com/skycoin/skycoin/src/httphelper"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_httphelper_HTTPError
func SKY_httphelper_HTTPError(_w *C.ResponseWriter, _status int, _httpMsg string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	status := _status
	httpMsg := _httpMsg
	httphelper.HTTPError(w, status, httpMsg)
	return
}

// export SKY_httphelper_Error400
func SKY_httphelper_Error400(_w *C.ResponseWriter, _msg string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	msg := _msg
	httphelper.Error400(w, msg)
	return
}

// export SKY_httphelper_Error403
func SKY_httphelper_Error403(_w *C.ResponseWriter) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	httphelper.Error403(w)
	return
}

// export SKY_httphelper_Error403Msg
func SKY_httphelper_Error403Msg(_w *C.ResponseWriter, _msg string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	msg := _msg
	httphelper.Error403Msg(w, msg)
	return
}

// export SKY_httphelper_Error404
func SKY_httphelper_Error404(_w *C.ResponseWriter) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	httphelper.Error404(w)
	return
}

// export SKY_httphelper_Error405
func SKY_httphelper_Error405(_w *C.ResponseWriter) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	httphelper.Error405(w)
	return
}

// export SKY_httphelper_Error501
func SKY_httphelper_Error501(_w *C.ResponseWriter) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	httphelper.Error501(w)
	return
}

// export SKY_httphelper_Error500
func SKY_httphelper_Error500(_w *C.ResponseWriter) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	httphelper.Error500(w)
	return
}

// export SKY_httphelper_Error500Msg
func SKY_httphelper_Error500Msg(_w *C.ResponseWriter, _msg string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	msg := _msg
	httphelper.Error500Msg(w, msg)
	return
}
