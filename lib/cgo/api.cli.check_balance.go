package main

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_CheckWalletBalance
func SKY_cli_CheckWalletBalance(_c *C.Handle, _walletFile string, _arg2 *C.cli__BalanceResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_c))
	____error_code = SKY_ERROR
	if ok {
		if client, isClient := (obj).(*webrpc.Client); isClient {
			balance, ____return_err := cli.CheckWalletBalance(client, _walletFile)
			if ____return_err == nil {
				*_arg2 = *(*C.cli__BalanceResult)(unsafe.Pointer(balance))
			}
			____error_code = libErrorCode(____return_err)
		}
	}
	return
}

// export SKY_cli_GetBalanceOfAddresses
func SKY_cli_GetBalanceOfAddresses(_c *C.Handle, _addrs *C.GoSlice_, _arg2 *C.cli__BalanceResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_c))
	____error_code = SKY_ERROR
	if ok {
		if client, isClient := (obj).(*webrpc.Client); isClient {
			addrs := *(*[]string)(unsafe.Pointer(_addrs))
			__arg2, ____return_err := cli.GetBalanceOfAddresses(client, addrs)
			____error_code = libErrorCode(____return_err)
			if ____return_err == nil {
				*_arg2 = *(*C.cli__BalanceResult)(unsafe.Pointer(__arg2))
			}
			____error_code = libErrorCode(____return_err)
		}
	}
	return
}
