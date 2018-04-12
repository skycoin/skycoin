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

// export SKY_cli_SendFromWallet
func SKY_cli_SendFromWallet(_c *C.webrpc__Client, _walletFile, _chgAddr string, _toAddrs *C.GoSlice_, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	walletFile := _walletFile
	chgAddr := _chgAddr
	toAddrs := *(*[]SendAmount)(unsafe.Pointer(_toAddrs))
	__arg3, ____return_err := cli.SendFromWallet(c, walletFile, chgAddr, toAddrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}

// export SKY_cli_SendFromAddress
func SKY_cli_SendFromAddress(_c *C.webrpc__Client, _addr, _walletFile, _chgAddr string, _toAddrs *C.GoSlice_, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	c := (*webrpc.Client)(unsafe.Pointer(_c))
	addr := _addr
	walletFile := _walletFile
	chgAddr := _chgAddr
	toAddrs := *(*[]SendAmount)(unsafe.Pointer(_toAddrs))
	__arg3, ____return_err := cli.SendFromAddress(c, addr, walletFile, chgAddr, toAddrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}
