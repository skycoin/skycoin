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

//export SKY_cli_CheckWalletBalance
func SKY_cli_CheckWalletBalance(_c C.WebRpcClient__Handle, _walletFile string, _arg2 *C.cli__BalanceResult) (____error_code uint32) {
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
	__arg2, ____return_err := cli.CheckWalletBalance(c, walletFile)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = *(*C.cli__BalanceResult)(unsafe.Pointer(__arg2))
	}
	return
}

//export SKY_cli_GetBalanceOfAddresses
func SKY_cli_GetBalanceOfAddresses(_c C.WebRpcClient__Handle, _addrs []string, _arg2 *C.cli__BalanceResult) (____error_code uint32) {
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
	__arg2, ____return_err := cli.GetBalanceOfAddresses(c, addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = *(*C.cli__BalanceResult)(unsafe.Pointer(__arg2))
	}
	return
}
