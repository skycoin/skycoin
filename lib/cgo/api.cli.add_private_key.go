package main

import cli "github.com/skycoin/skycoin/src/api/cli"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_cli_AddPrivateKey
func SKY_cli_AddPrivateKey(_wlt *C.Wallet__Handle, _key string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	wlt, okwlt := lookupWalletHandle(*_wlt)
	if !okwlt {
		____error_code = SKY_ERROR
		return
	}
	key := _key
	____return_err := cli.AddPrivateKey(wlt, key)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_cli_AddPrivateKeyToFile
func SKY_cli_AddPrivateKeyToFile(_walletFile, _key string, pwd string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	walletFile := _walletFile
	key := _key
	pr := cli.NewPasswordReader([]byte(pwd))
	____return_err := cli.AddPrivateKeyToFile(walletFile, key, pr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
