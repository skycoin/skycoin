package main

import (
	cli "github.com/skycoin/skycoin/src/cli"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_CheckWalletBalance
func SKY_cli_CheckWalletBalance(_c *C.Client, _walletFile string, _arg2 *C.BalanceResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	walletFile := _walletFile
	__arg2, ____return_err := cli.CheckWalletBalance(c, walletFile)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofBalanceResult))
	}
	return
}

// export SKY_cli_GetBalanceOfAddresses
func SKY_cli_GetBalanceOfAddresses(_c *C.Client, _addrs *C.GoSlice_, _arg2 *C.BalanceResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg2, ____return_err := cli.GetBalanceOfAddresses(c, addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofBalanceResult))
	}
	return
}
