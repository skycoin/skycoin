package main

import cli "github.com/skycoin/skycoin/src/cli"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_AddPrivateKey
func SKY_cli_AddPrivateKey(_wlt *C.Wallet, _key string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	key := _key
	____return_err := cli.AddPrivateKey(wlt, key)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_cli_AddPrivateKeyToFile
func SKY_cli_AddPrivateKeyToFile(_walletFile, _key string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	walletFile := _walletFile
	key := _key
	____return_err := cli.AddPrivateKeyToFile(walletFile, key)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
