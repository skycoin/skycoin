package main

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	wallet "github.com/skycoin/skycoin/src/wallet"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_AddPrivateKey
func SKY_cli_AddPrivateKey(_wlt *C.wallet__Wallet, _key string) (____error_code uint32) {
	____error_code = 0
	wlt := (*wallet.Wallet)(unsafe.Pointer(_wlt))
	key := _key
	____return_err := cli.AddPrivateKey(wlt, key)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_cli_AddPrivateKeyToFile
func SKY_cli_AddPrivateKeyToFile(_walletFile, _key string) (____error_code uint32) {
	____error_code = 0
	walletFile := _walletFile
	key := _key
	____return_err := cli.AddPrivateKeyToFile(walletFile, key)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
