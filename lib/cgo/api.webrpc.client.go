package main

import (
	coin "github.com/skycoin/skycoin/src/coin"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_webrpc_Client_Do
func SKY_webrpc_Client_Do(_c *C.webrpc__Client, _obj interface{}, _method string, _params interface{}) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	method := _method
	____return_err := c.Do(obj, method, params)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_Client_GetUnspentOutputs
func SKY_webrpc_Client_GetUnspentOutputs(_c *C.webrpc__Client, _addrs *C.GoSlice_, _arg1 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1, ____return_err := c.GetUnspentOutputs(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofOutputsResult))
	}
	return
}

// export SKY_webrpc_Client_InjectTransactionString
func SKY_webrpc_Client_InjectTransactionString(_c *C.webrpc__Client, _rawtx string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	rawtx := _rawtx
	__arg1, ____return_err := c.InjectTransactionString(rawtx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_webrpc_Client_InjectTransaction
func SKY_webrpc_Client_InjectTransaction(_c *C.webrpc__Client, _tx *C.coin__Transaction, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	tx := (*coin.Transaction)(unsafe.Pointer(_tx))
	__arg1, ____return_err := c.InjectTransaction(tx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_webrpc_Client_GetStatus
func SKY_webrpc_Client_GetStatus(_c *C.webrpc__Client, _arg0 *C.webrpc__StatusResult) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.GetStatus()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofStatusResult))
	}
	return
}

// export SKY_webrpc_Client_GetTransactionByID
func SKY_webrpc_Client_GetTransactionByID(_c *C.webrpc__Client, _txid string, _arg1 *C.webrpc__TxnResult) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	txid := _txid
	__arg1, ____return_err := c.GetTransactionByID(txid)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTxnResult))
	}
	return
}

// export SKY_webrpc_Client_GetAddressUxOuts
func SKY_webrpc_Client_GetAddressUxOuts(_c *C.webrpc__Client, _addrs *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1, ____return_err := c.GetAddressUxOuts(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_webrpc_Client_GetBlocks
func SKY_webrpc_Client_GetBlocks(_c *C.webrpc__Client, _start, _end uint64, _arg1 *C.visor__ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	start := _start
	end := _end
	__arg1, ____return_err := c.GetBlocks(start, end)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_Client_GetBlocksBySeq
func SKY_webrpc_Client_GetBlocksBySeq(_c *C.webrpc__Client, _ss *C.GoSlice_, _arg1 *C.visor__ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	ss := *(*[]uint64)(unsafe.Pointer(_ss))
	__arg1, ____return_err := c.GetBlocksBySeq(ss)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_webrpc_Client_GetLastBlocks
func SKY_webrpc_Client_GetLastBlocks(_c *C.webrpc__Client, _n uint64, _arg1 *C.visor__ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	n := _n
	__arg1, ____return_err := c.GetLastBlocks(n)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
