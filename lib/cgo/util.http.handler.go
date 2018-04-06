package main

import httphelper "github.com/skycoin/skycoin/src/httphelper"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_httphelper_HostCheck
func SKY_httphelper_HostCheck(_logger *C.Logger, _host string, _handler *C.Handler, _arg3 *C.Handler) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	host := _host
	__arg3 := httphelper.HostCheck(logger, host, handler)
	return
}
