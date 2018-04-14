package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

import (
	"unsafe"
	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

type Handle uint64

var (
	handleMap = make(map[Handle]interface{})
)

func registerHandle(obj interface{}) Handle {
	ptr := &obj
	handle := *(*Handle)(unsafe.Pointer(&ptr))
	handleMap[handle] = obj
	return handle
}

func lookupHandleObj(handle Handle) (interface{}, bool) {
	obj, ok := handleMap[handle]
	return obj, ok
}

func registerWebRpcClientHandle(obj *webrpc.Client) C.WebrpcClient__Handle{
	return (C.WebrpcClient__Handle)(registerHandle(obj))
}

func lookupWebRpcClientHandle(handle C.WebrpcClient__Handle) (*webrpc.Client, bool){
	obj, ok := lookupHandleObj(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*webrpc.Client); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerWalletHandle(obj *wallet.Wallet) C.Wallet__Handle{
	return (C.Wallet__Handle)(registerHandle(obj))
}

func lookupWalletHandle(handle C.Wallet__Handle) (*wallet.Wallet, bool){
	obj, ok := lookupHandleObj(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*wallet.Wallet); isOK {
			return obj, true
		}
	}
	return nil, false
}

func closeHandle(handle Handle) {
	delete(handleMap, handle)
}

//export SKY_handle_close
func SKY_handle_close(handle C.Handle){
	closeHandle(Handle(handle))
}