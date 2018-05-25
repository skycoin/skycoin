package main

import browser "github.com/skycoin/skycoin/src/util/browser"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_browser_Open
func SKY_browser_Open(_url string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	url := _url
	____return_err := browser.Open(url)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
