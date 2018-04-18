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
func SKY_cli_SendFromWallet(_c *C.WebRpcClient__Handle, _walletFile, _chgAddr string, _toAddrs *C.GoSlice_, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	client, ok := lookupWebRpcClientHandle(*_c)
	____error_code = SKY_ERROR
	if ok {
		toAddrs := *(*[]cli.SendAmount)(unsafe.Pointer(_toAddrs))
		__arg3, ____return_err := cli.SendFromWallet(client, _walletFile, _chgAddr, toAddrs)
		____error_code = libErrorCode(____return_err)
		if ____return_err == nil {
			copyString(__arg3, _arg3)
		}
		____error_code = libErrorCode(____return_err)
	}
	return
}

//export SKY_cli_SendFromAddress
func SKY_cli_SendFromAddress(_c *C.WebRpcClient__Handle, _addr, _walletFile, _chgAddr string, _toAddrs *C.GoSlice_, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	client, ok := lookupWebRpcClientHandle(*_c)
	____error_code = SKY_ERROR
	if ok {
		toAddrs := *(*[]cli.SendAmount)(unsafe.Pointer(_toAddrs))
		__arg3, ____return_err := cli.SendFromAddress(client, _addr, _walletFile, _chgAddr, toAddrs)
		____error_code = libErrorCode(____return_err)
		if ____return_err == nil {
			copyString(__arg3, _arg3)
		}
		____error_code = libErrorCode(____return_err)
	}
	return
}
