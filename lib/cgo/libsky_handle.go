package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

import (
	"unsafe"

	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
	cli "github.com/skycoin/skycoin/src/cli"
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

func lookupHandle(handle Handle) (interface{}, bool) {
	obj, ok := handleMap[handle]
	return obj, ok
}

func registerWebRpcClientHandle(obj *webrpc.Client) C.WebRpcClient__Handle {
	return (C.WebRpcClient__Handle)(registerHandle(obj))
}

func lookupWebRpcClientHandle(handle C.WebRpcClient__Handle) (*webrpc.Client, bool) {
	obj, ok := lookupHandle(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*webrpc.Client); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerWalletHandle(obj *wallet.Wallet) C.Wallet__Handle {
	return (C.Wallet__Handle)(registerHandle(obj))
}

func lookupWalletHandle(handle C.Wallet__Handle) (*wallet.Wallet, bool) {
	obj, ok := lookupHandle(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*wallet.Wallet); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerOptionsHandle(obj *wallet.Options) C.Options__Handle {
	return (C.Options__Handle)(registerHandle(obj))
}

func lookupOptionsHandle(handle C.Options__Handle) (*wallet.Options, bool) {
	obj, ok := lookupHandle(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*wallet.Options); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerConfigHandle(obj *cli.Config) C.Config__Handle {
	return (C.Config__Handle)(registerHandle(obj))
}

func lookupConfigHandle(handle C.Config__Handle) (*cli.Config, bool) {
	obj, ok := lookupHandle(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*cli.Config); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerPasswordReaderHandle(obj cli.PasswordReader) C.PasswordReader__Handle {
	return (C.PasswordReader__Handle)(registerHandle(obj))
}

func lookupPasswordReaderHandle(handle C.PasswordReader__Handle) (cli.PasswordReader, bool) {
	obj, ok := lookupHandle(Handle(handle))
	if ok {
		if obj, isOK := (obj).(cli.PasswordReader); isOK {
			return obj, true
		}
	}
	return nil, false
}

func closeHandle(handle Handle) {
	delete(handleMap, handle)
}

//export SKY_handle_close
func SKY_handle_close(handle C.Handle) {
	closeHandle(Handle(handle))
}
