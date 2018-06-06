package main

import (
	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
)

/*
  #include <string.h>
  #include <stdlib.h>

	#include "skytypes.h"
*/
import "C"

//export SKY_webrpc_NewClient
func SKY_webrpc_NewClient(_addr string, _arg1 *C.WebRpcClient__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := _addr
	__arg1, ____return_err := webrpc.NewClient(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerWebRpcClientHandle(__arg1)
	}
	return
}

//export SKY_webrpc_Client_CSRF
func SKY_webrpc_Client_CSRF(_c C.WebRpcClient__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.CSRF()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg0, _arg0)
	}
	return
}
