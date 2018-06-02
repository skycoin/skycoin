package main

import (
	"os"

	"github.com/skycoin/skycoin/src/api/webrpc"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cli_App_Run
func SKY_cli_App_Run(_app C.App__Handle, _args string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	app, okapp := lookupAppHandle(_app)
	if !okapp {
		____error_code = SKY_ERROR
		return
	}
	args := splitCliArgs(_args)
	____return_err := app.Run(args)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_cli_Config_GetCoin
func SKY_cli_Config_GetCoin(_c C.Config__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__c, okc := lookupConfigHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	c := *__c
	__arg0 := c.Coin
	copyString(__arg0, _arg0)
	return
}

//export SKY_cli_Config_GetRPCAddress
func SKY_cli_Config_GetRPCAddress(_c C.Config__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__c, okc := lookupConfigHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	c := *__c
	__arg0 := c.RPCAddress
	copyString(__arg0, _arg0)
	return
}

//export SKY_cli_RPCClientFromApp
func SKY_cli_RPCClientFromApp(_app C.App__Handle, _arg1 *C.WebRpcClient__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	app, okapp := lookupAppHandle(_app)
	if !okapp {
		____error_code = SKY_ERROR
		return
	}
	__arg1 := app.App.Metadata["rpc"].(*webrpc.Client)
	*_arg1 = registerWebRpcClientHandle(__arg1)
	return
}

//export SKY_cli_Getenv
func SKY_cli_Getenv(varname string, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := os.Getenv(varname)
	copyString(__arg0, _arg0)
	return
}

//export SKY_cli_Setenv
func SKY_cli_Setenv(varname string, value string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	os.Setenv(varname, value)
	return
}
