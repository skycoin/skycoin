package main

import (
	"unsafe"

	cli "github.com/skycoin/skycoin/src/cli"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cli_GetWalletOutputsFromFile
func SKY_cli_GetWalletOutputsFromFile(_c C.WebRpcClient__Handle, _walletFile string, _arg2 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	walletFile := _walletFile
	__arg2, ____return_err := cli.GetWalletOutputsFromFile(c, walletFile)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = *(*C.webrpc__OutputsResult)(unsafe.Pointer(__arg2))
	}
	return
}

//export SKY_cli_GetWalletOutputs
func SKY_cli_GetWalletOutputs(_c C.WebRpcClient__Handle, _wlt *C.Wallet__Handle, _arg2 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	wlt, okwlt := lookupWalletHandle(*_wlt)
	if !okwlt {
		____error_code = SKY_ERROR
		return
	}
	__arg2, ____return_err := cli.GetWalletOutputs(c, wlt)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = *(*C.webrpc__OutputsResult)(unsafe.Pointer(__arg2))
	}
	return
}
