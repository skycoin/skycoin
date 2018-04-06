package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	gui "github.com/skycoin/skycoin/src/gui"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gui_CSRFToken_String
func SKY_gui_CSRFToken_String(_c *C.CSRFToken, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.CSRFToken)(unsafe.Pointer(_c))
	__arg0 := c.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_gui_CSRFCheck
func SKY_gui_CSRFCheck(_store *C.CSRFStore, _handler *C.Handler, _arg2 *C.Handler) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	store := (*cipher.CSRFStore)(unsafe.Pointer(_store))
	__arg2 := gui.CSRFCheck(store, handler)
	return
}

// export SKY_gui_OriginRefererCheck
func SKY_gui_OriginRefererCheck(_host string, _handler *C.Handler, _arg2 *C.Handler) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	host := _host
	__arg2 := gui.OriginRefererCheck(host, handler)
	return
}
