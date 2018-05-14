package main

import (
	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_webrpc_ClientError_Error
func SKY_webrpc_ClientError_Error(_e *C.webrpc__ClientError, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := *(*webrpc.ClientError)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}
