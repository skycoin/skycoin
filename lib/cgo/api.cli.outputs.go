package main

import cli "github.com/skycoin/skycoin/src/api/cli"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_GetWalletOutputsFromFile
func SKY_cli_GetWalletOutputsFromFile(_c *C.Client, _walletFile string, _arg2 *C.OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	walletFile := _walletFile
	__arg2, ____return_err := cli.GetWalletOutputsFromFile(c, walletFile)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_cli_GetWalletOutputs
func SKY_cli_GetWalletOutputs(_c *C.Client, _wlt *C.Wallet, _arg2 *C.OutputsResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg2, ____return_err := cli.GetWalletOutputs(c, wlt)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
