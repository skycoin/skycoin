package main

/*

  #include <string.h>
  #include <stdlib.h>


  #include "skytypes.h"
*/
import "C"

import (
	api "github.com/skycoin/skycoin/src/api"
	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
	cli "github.com/skycoin/skycoin/src/cli"
	wallet "github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

type Handle uint64

var (
	handlesCounter uint64 = 0
	handleMap             = make(map[Handle]interface{})
)

func registerHandle(obj interface{}) C.Handle {
	handlesCounter++
	handle := handlesCounter
	//handle := *(*Handle)(unsafe.Pointer(&obj))
	handleMap[Handle(handle)] = obj
	return (C.Handle)(handle)
}

func lookupHandle(handle C.Handle) (interface{}, bool) {
	obj, ok := handleMap[Handle(handle)]
	return obj, ok
}

func registerWebRpcClientHandle(obj *webrpc.Client) C.WebRpcClient__Handle {
	return (C.WebRpcClient__Handle)(registerHandle(obj))
}

func lookupWebRpcClientHandle(handle C.WebRpcClient__Handle) (*webrpc.Client, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
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
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*wallet.Wallet); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerReadableWalletHandle(obj *wallet.ReadableWallet) C.ReadableWallet__Handle {
	return (C.ReadableWallet__Handle)(registerHandle(obj))
}

func lookupReadableWalletHandle(handle C.ReadableWallet__Handle) (*wallet.ReadableWallet, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*wallet.ReadableWallet); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerReadableEntryHandle(obj *wallet.ReadableEntry) C.ReadableEntry__Handle {
	return (C.ReadableEntry__Handle)(registerHandle(obj))
}

func lookupReadableEntryHandle(handle C.ReadableEntry__Handle) (*wallet.ReadableEntry, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*wallet.ReadableEntry); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerOptionsHandle(obj *wallet.Options) C.Options__Handle {
	return (C.Options__Handle)(registerHandle(obj))
}

func lookupOptionsHandle(handle C.Options__Handle) (*wallet.Options, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
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
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*cli.Config); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerAppHandle(obj *cli.App) C.App__Handle {
	return (C.App__Handle)(registerHandle(obj))
}

func lookupAppHandle(handle C.App__Handle) (*cli.App, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*cli.App); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerContextHandle(obj *gcli.Context) C.Context__Handle {
	return (C.Context__Handle)(registerHandle(obj))
}

func lookupContextHandle(handle C.Context__Handle) (*gcli.Context, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*gcli.Context); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerClientHandle(obj *api.Client) C.Client__Handle {
	return (C.Client__Handle)(registerHandle(obj))
}

func lookupClientHandle(handle C.Client__Handle) (*api.Client, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*api.Client); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerWalletsHandle(obj []*api.WalletResponse) C.Wallets__Handle {
	return (C.Wallets__Handle)(registerHandle(obj))
}

func lookupWalletsHandle(handle C.Wallets__Handle) ([]*api.WalletResponse, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).([]*api.WalletResponse); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerWalletResponseHandle(obj *api.WalletResponse) C.WalletResponse__Handle {
	return (C.WalletResponse__Handle)(registerHandle(obj))
}

func lookupWalletResponseHandle(handle C.WalletResponse__Handle) (*api.WalletResponse, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*api.WalletResponse); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerCreateTransactionRequestHandle(obj *api.CreateTransactionRequest) C.CreateTransactionRequest__Handle {
	return (C.CreateTransactionRequest__Handle)(registerHandle(obj))
}

func lookupCreateTransactionRequestHandle(handle C.CreateTransactionRequest__Handle) (*api.CreateTransactionRequest, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*api.CreateTransactionRequest); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerPasswordReaderHandle(obj *cli.PasswordReader) C.PasswordReader__Handle {
	return (C.PasswordReader__Handle)(registerHandle(obj))
}

func lookupPasswordReaderHandle(handle C.PasswordReader__Handle) (*cli.PasswordReader, bool) {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*cli.PasswordReader); isOK {
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
