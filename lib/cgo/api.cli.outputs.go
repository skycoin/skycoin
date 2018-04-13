package main

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
	wallet "github.com/skycoin/skycoin/src/wallet"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_GetWalletOutputsFromFile
func SKY_cli_GetWalletOutputsFromFile(_c *C.Handle, _walletFile string, _arg2 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_c))
	____error_code = SKY_ERROR
	if ok {
		if client, isClient := (obj).(*webrpc.Client); isClient {
			__arg2, ____return_err := cli.GetWalletOutputsFromFile(client, _walletFile)
			____error_code = libErrorCode(____return_err)
			if ____return_err == nil {
				*_arg2 = *(*C.webrpc__OutputsResult)(unsafe.Pointer(__arg2))
			}
			____error_code = libErrorCode(____return_err)
		}
	}
	return
}

// export SKY_cli_GetWalletOutputs
func SKY_cli_GetWalletOutputs(_c *C.Handle, _wlt *C.wallet__Wallet, _arg2 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_c))
	____error_code = SKY_ERROR
	if ok {
		if client, isClient := (obj).(*webrpc.Client); isClient {
			wlt := (*wallet.Wallet)(unsafe.Pointer(_wlt))
			__arg2, ____return_err := cli.GetWalletOutputs(client, wlt)
			____error_code = libErrorCode(____return_err)
			if ____return_err == nil {
				*_arg2 = *(*C.webrpc__OutputsResult)(unsafe.Pointer(__arg2))
			}
			____error_code = libErrorCode(____return_err)
		}
	}
	return
}
