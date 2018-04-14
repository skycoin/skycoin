package main

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//TODO: Use WebRpc.Client specific handle type
//export SKY_cli_GetWalletOutputsFromFile
func SKY_cli_GetWalletOutputsFromFile(_c *C.WebrpcClient__Handle, _walletFile string, _arg2 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	client, ok := lookupWebRpcClientHandle(*_c)
	____error_code = SKY_ERROR
	if ok {
		__arg2, ____return_err := cli.GetWalletOutputsFromFile(client, _walletFile)
		____error_code = libErrorCode(____return_err)
		if ____return_err == nil {
			*_arg2 = *(*C.webrpc__OutputsResult)(unsafe.Pointer(__arg2))
		}
		____error_code = libErrorCode(____return_err)
	}
	return
}

//TODO: Use WebRpc.Client specific handle type
//export SKY_cli_GetWalletOutputs
func SKY_cli_GetWalletOutputs(_c *C.WebrpcClient__Handle, _wlt *C.Wallet__Handle, _arg2 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	client, isClient := lookupWebRpcClientHandle(*_c)
	wlt, isWallet := lookupWalletHandle(*_wlt)
	____error_code = SKY_ERROR
	if isClient && isWallet {
		__arg2, ____return_err := cli.GetWalletOutputs(client, wlt)
		____error_code = libErrorCode(____return_err)
		if ____return_err == nil {
			*_arg2 = *(*C.webrpc__OutputsResult)(unsafe.Pointer(__arg2))
		}
		____error_code = libErrorCode(____return_err)
	}
	return
}
