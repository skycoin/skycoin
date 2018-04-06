package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	cli "github.com/skycoin/skycoin/src/cli"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_SendFromWallet
func SKY_cli_SendFromWallet(_c *C.Client, _walletFile, _chgAddr string, _toAddrs *C.GoSlice_, _arg3 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	walletFile := _walletFile
	chgAddr := _chgAddr
	toAddrs := *(*[]cipher.SendAmount)(unsafe.Pointer(_toAddrs))
	__arg3, ____return_err := cli.SendFromWallet(c, walletFile, chgAddr, toAddrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}

// export SKY_cli_SendFromAddress
func SKY_cli_SendFromAddress(_c *C.Client, _addr, _walletFile, _chgAddr string, _toAddrs *C.GoSlice_, _arg3 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	addr := _addr
	walletFile := _walletFile
	chgAddr := _chgAddr
	toAddrs := *(*[]cipher.SendAmount)(unsafe.Pointer(_toAddrs))
	__arg3, ____return_err := cli.SendFromAddress(c, addr, walletFile, chgAddr, toAddrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}
