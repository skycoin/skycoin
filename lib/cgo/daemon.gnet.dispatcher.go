package main

import (
	gnet "github.com/skycoin/skycoin/src/daemon/gnet"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gnet_EncodeMessage
func SKY_gnet_EncodeMessage(_msg *C.Message, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msg := *(*Message)(unsafe.Pointer(_msg))
	__arg1 := gnet.EncodeMessage(msg)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}
