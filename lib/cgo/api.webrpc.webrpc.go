package main

import (
	"unsafe"

	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_webrpc_RPCError_Error
func SKY_webrpc_RPCError_Error(_e *C.webrpc__RPCError, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := *(*webrpc.RPCError)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}
