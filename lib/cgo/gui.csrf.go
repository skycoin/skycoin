package main

import (
	gui "github.com/skycoin/skycoin/src/gui"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gui_CSRFToken_String
func SKY_gui_CSRFToken_String(_c *C.CSRFToken, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := (*CSRFToken)(unsafe.Pointer(_c))
	__arg0 := c.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_gui_CSRFCheck
func SKY_gui_CSRFCheck(_store *C.CSRFStore, _handler *C.Handler, _arg2 *C.Handler) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	store := (*CSRFStore)(unsafe.Pointer(_store))
	__arg2 := gui.CSRFCheck(store, handler)
	return
}

// export SKY_gui_OriginRefererCheck
func SKY_gui_OriginRefererCheck(_host string, _handler *C.Handler, _arg2 *C.Handler) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	host := _host
	__arg2 := gui.OriginRefererCheck(host, handler)
	return
}
