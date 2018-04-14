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

//export SKY_cli_GenerateWallet
func SKY_cli_GenerateWallet(_walletFile, _label, _seed string, _numAddrs uint64, _arg2 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	walletFile := _walletFile
	label := _label
	seed := _seed
	numAddrs := _numAddrs
	__arg2, ____return_err := cli.GenerateWallet(walletFile, label, seed, numAddrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = *(*C.Wallet)(unsafe.Pointer(__arg2))
	}
	return
}

//export SKY_cli_MakeAlphanumericSeed
func SKY_cli_MakeAlphanumericSeed(_arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := cli.MakeAlphanumericSeed()
	copyString(__arg0, _arg0)
	return
}
