package main

import (
	"reflect"
	"unsafe"

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

//export SKY_webrpc_Client_GetUnspentOutputs
func SKY_webrpc_Client_GetUnspentOutputs(_c C.WebRpcClient__Handle, _addrs []string, _arg1 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1, ____return_err := c.GetUnspentOutputs(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.webrpc__OutputsResult)(unsafe.Pointer(__arg1))
	}
	return
}

//export SKY_webrpc_Client_InjectTransactionString
func SKY_webrpc_Client_InjectTransactionString(_c C.WebRpcClient__Handle, _rawtx string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	rawtx := _rawtx
	__arg1, ____return_err := c.InjectTransactionString(rawtx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

//export SKY_webrpc_Client_InjectTransaction
func SKY_webrpc_Client_InjectTransaction(_c C.WebRpcClient__Handle, _tx C.Transaction__Handle, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	tx, ok := lookupTransactionHandle(_tx)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg1, ____return_err := c.InjectTransaction(tx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

//export SKY_webrpc_Client_GetStatus
func SKY_webrpc_Client_GetStatus(_c C.WebRpcClient__Handle, _arg0 *C.webrpc__StatusResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.GetStatus()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = *(*C.webrpc__StatusResult)(unsafe.Pointer(__arg0))
	}
	return
}

//export SKY_webrpc_Client_GetTransactionByID
func SKY_webrpc_Client_GetTransactionByID(_c C.WebRpcClient__Handle, _txid string, _arg1 *C.webrpc__TxnResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	txid := _txid
	__arg1, ____return_err := c.GetTransactionByID(txid)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.webrpc__TxnResult)(unsafe.Pointer(__arg1))
	}
	return
}

//export SKY_webrpc_Client_GetAddressUxOuts
func SKY_webrpc_Client_GetAddressUxOuts(_c C.WebRpcClient__Handle, _addrs []string, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1, ____return_err := c.GetAddressUxOuts(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

//export SKY_webrpc_Client_GetBlocks
func SKY_webrpc_Client_GetBlocks(_c C.WebRpcClient__Handle, _start, _end uint64, _arg1 *C.visor__ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	start := _start
	end := _end
	__arg1, ____return_err := c.GetBlocks(start, end)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.visor__ReadableBlocks)(unsafe.Pointer(__arg1))
	}
	return
}

//export SKY_webrpc_Client_GetBlocksBySeq
func SKY_webrpc_Client_GetBlocksBySeq(_c C.WebRpcClient__Handle, _ss []uint64, _arg1 *C.visor__ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	ss := *(*[]uint64)(unsafe.Pointer(&_ss))
	__arg1, ____return_err := c.GetBlocksBySeq(ss)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.visor__ReadableBlocks)(unsafe.Pointer(__arg1))
	}
	return
}

//export SKY_webrpc_Client_GetLastBlocks
func SKY_webrpc_Client_GetLastBlocks(_c C.WebRpcClient__Handle, _n uint64, _arg1 *C.visor__ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	n := _n
	__arg1, ____return_err := c.GetLastBlocks(n)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.visor__ReadableBlocks)(unsafe.Pointer(__arg1))
	}
	return
}
