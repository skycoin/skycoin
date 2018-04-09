package main

import (
	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_webrpc_RPCError_Error
func SKY_webrpc_RPCError_Error(_e *C.RPCError, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := *(*RPCError)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_webrpc_NewRequest
func SKY_webrpc_NewRequest(_method string, _params interface{}, _id string, _arg3 *C.Request) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	method := _method
	id := _id
	__arg3, ____return_err := webrpc.NewRequest(method, params, id)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofRequest))
	}
	return
}

// export SKY_webrpc_Request_DecodeParams
func SKY_webrpc_Request_DecodeParams(_r *C.Request, _v interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	r := (*Request)(unsafe.Pointer(_r))
	____return_err := r.DecodeParams(v)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_New
func SKY_webrpc_New(_addr string, _c *C.Config, _gw *C.Gatewayer, _arg3 *C.WebRPC) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := _addr
	c := *(*Config)(unsafe.Pointer(_c))
	gw := *(*Gatewayer)(unsafe.Pointer(_gw))
	__arg3, ____return_err := webrpc.New(addr, c, gw)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofWebRPC))
	}
	return
}

// export SKY_webrpc_WebRPC_Run
func SKY_webrpc_WebRPC_Run(_rpc *C.WebRPC) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rpc := (*WebRPC)(unsafe.Pointer(_rpc))
	____return_err := rpc.Run()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_WebRPC_Shutdown
func SKY_webrpc_WebRPC_Shutdown(_rpc *C.WebRPC) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rpc := (*WebRPC)(unsafe.Pointer(_rpc))
	____return_err := rpc.Shutdown()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_WebRPC_HandleFunc
func SKY_webrpc_WebRPC_HandleFunc(_rpc *C.WebRPC, _method string, _h *C.HandlerFunc) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rpc := (*WebRPC)(unsafe.Pointer(_rpc))
	method := _method
	h := *(*HandlerFunc)(unsafe.Pointer(_h))
	____return_err := rpc.HandleFunc(method, h)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_WebRPC_Handler
func SKY_webrpc_WebRPC_Handler(_rpc *C.WebRPC, _w *C.ResponseWriter, _r *C.Request) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rpc := (*WebRPC)(unsafe.Pointer(_rpc))
	rpc.Handler(w, r)
	return
}
