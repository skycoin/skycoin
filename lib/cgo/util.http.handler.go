package main

import http "github.com/skycoin/skycoin/src/util/http"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_httphelper_HostCheck
func SKY_httphelper_HostCheck(_logger *C.Logger, _host string, _handler *C.Handler, _arg3 *C.Handler) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	host := _host
	__arg3 := http.HostCheck(logger, host, handler)
	return
}
