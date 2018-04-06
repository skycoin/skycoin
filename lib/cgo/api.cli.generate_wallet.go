package main

import cli "github.com/skycoin/skycoin/src/cli"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_GenerateWallet
func SKY_cli_GenerateWallet(_walletFile, _label, _seed string, _numAddrs uint64, _arg2 *C.Wallet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	walletFile := _walletFile
	label := _label
	seed := _seed
	numAddrs := _numAddrs
	__arg2, ____return_err := cli.GenerateWallet(walletFile, label, seed, numAddrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_cli_MakeAlphanumericSeed
func SKY_cli_MakeAlphanumericSeed(_arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := cli.MakeAlphanumericSeed()
	copyString(__arg0, _arg0)
	return
}
