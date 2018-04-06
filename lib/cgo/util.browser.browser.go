package main

import browser "github.com/skycoin/skycoin/src/browser"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_browser_Open
func SKY_browser_Open(_url string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	url := _url
	____return_err := browser.Open(url)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
