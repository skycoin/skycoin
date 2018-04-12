package main

import "unsafe"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_webrpc_RPCError_Error
func SKY_webrpc_RPCError_Error(_e *C.webrpc__RPCError, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	e := *(*webrpc.RPCError)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}
