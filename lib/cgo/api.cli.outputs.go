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
func SKY_cli_GetWalletOutputsFromFile(_c *C.webrpc__Client, _walletFile string, _arg2 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	walletFile := _walletFile
	__arg2, ____return_err := cli.GetWalletOutputsFromFile(c, walletFile)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_cli_GetWalletOutputs
func SKY_cli_GetWalletOutputs(_c *C.webrpc__Client, _wlt *C.wallet__Wallet, _arg2 *C.webrpc__OutputsResult) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	wlt := (*wallet.Wallet)(unsafe.Pointer(_wlt))
	__arg2, ____return_err := cli.GetWalletOutputs(c, wlt)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
