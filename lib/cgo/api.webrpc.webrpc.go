package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	webrpc "github.com/skycoin/skycoin/src/webrpc"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_webrpc_RPCError_Error
func SKY_webrpc_RPCError_Error(_e *C.RPCError, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	e := *(*cipher.RPCError)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_webrpc_NewRequest
func SKY_webrpc_NewRequest(_method string, _params interface{}, _id string, _arg3 *C.Request) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	method := _method
	id := _id
	__arg3, ____return_err := webrpc.NewRequest(method, params, id)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofRequest))
	}
	return
}

// export SKY_webrpc_Request_DecodeParams
func SKY_webrpc_Request_DecodeParams(_r *C.Request, _v interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	r := (*cipher.Request)(unsafe.Pointer(_r))
	____return_err := r.DecodeParams(v)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_New
func SKY_webrpc_New(_addr string, _c *C.Config, _gw *C.Gatewayer, _arg3 *C.WebRPC) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	addr := _addr
	c := *(*cipher.Config)(unsafe.Pointer(_c))
	gw := *(*cipher.Gatewayer)(unsafe.Pointer(_gw))
	__arg3, ____return_err := webrpc.New(addr, c, gw)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofWebRPC))
	}
	return
}

// export SKY_webrpc_WebRPC_Run
func SKY_webrpc_WebRPC_Run(_rpc *C.WebRPC) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.WebRPC)(unsafe.Pointer(_rpc))
	____return_err := rpc.Run()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_WebRPC_Shutdown
func SKY_webrpc_WebRPC_Shutdown(_rpc *C.WebRPC) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.WebRPC)(unsafe.Pointer(_rpc))
	____return_err := rpc.Shutdown()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_WebRPC_HandleFunc
func SKY_webrpc_WebRPC_HandleFunc(_rpc *C.WebRPC, _method string, _h *C.HandlerFunc) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.WebRPC)(unsafe.Pointer(_rpc))
	method := _method
	h := *(*cipher.HandlerFunc)(unsafe.Pointer(_h))
	____return_err := rpc.HandleFunc(method, h)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_WebRPC_Handler
func SKY_webrpc_WebRPC_Handler(_rpc *C.WebRPC, _w *C.ResponseWriter, _r *C.Request) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.WebRPC)(unsafe.Pointer(_rpc))
	rpc.Handler(w, r)
	return
}
