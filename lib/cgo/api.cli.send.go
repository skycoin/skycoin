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

//export SKY_cli_SendFromWallet
func SKY_cli_SendFromWallet(_c *C.WebRpcClient__Handle, _walletFile, _chgAddr string, _toAddrs []C.cli__SendAmount, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(*_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	walletFile := _walletFile
	chgAddr := _chgAddr
	toAddrs := *(*[]cli.SendAmount)(unsafe.Pointer(&_toAddrs))
	__arg3, ____return_err := cli.SendFromWallet(c, walletFile, chgAddr, toAddrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}

//export SKY_cli_SendFromAddress
func SKY_cli_SendFromAddress(_c *C.WebRpcClient__Handle, _addr, _walletFile, _chgAddr string, _toAddrs []C.cli__SendAmount, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(*_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addr := _addr
	walletFile := _walletFile
	chgAddr := _chgAddr
	toAddrs := *(*[]cli.SendAmount)(unsafe.Pointer(&_toAddrs))
	__arg3, ____return_err := cli.SendFromAddress(c, addr, walletFile, chgAddr, toAddrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}
