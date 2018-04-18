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

//export SKY_cli_LoadConfig
func SKY_cli_LoadConfig(_config *C.Config__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	config, ____return_err := cli.LoadConfig()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_config = (registerConfigHandle(&config))
	}
	return
}

//export SKY_cli_Config_FullWalletPath
func SKY_cli_Config_FullWalletPath(_c *C.Config__Handle, _path *C.GoString_) (____error_code uint32) {
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	config, ok := lookupConfigHandle(*_c)
	____error_code = SKY_ERROR
	if ok {
		path := config.FullWalletPath()
		copyString(path, _path)
		____error_code = SKY_OK
	} 
	return
}

//export SKY_cli_Config_FullDBPath
func SKY_cli_Config_FullDBPath(_c *C.Config__Handle, _path *C.GoString_) (____error_code uint32) {
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	config, ok := lookupConfigHandle(*_c)
	____error_code = SKY_ERROR
	if ok {
		path := config.FullDBPath()
		copyString(path, _path)
		____error_code = SKY_OK
	} 
	return
}

//export SKY_cli_NewApp
func SKY_cli_NewApp(_cfg *C.Config__Handle, _app *C.App__Handle) (____error_code uint32) {
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	config, ok := lookupConfigHandle(*_cfg)
	____error_code = SKY_ERROR
	if ok {
		app := cli.NewApp(*config)
		*_app = registerAppHandle( app )
		____error_code = SKY_OK
	} 
	return
}

//export SKY_cli_App_Run
func SKY_cli_App_Run(_app *C.App__Handle, _args *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	____error_code = SKY_ERROR
	app, ok := lookupAppHandle(*_app)
	if ok {
		args := *(*[]string)(unsafe.Pointer(_args))
		____return_err := app.Run(args)
		____error_code = libErrorCode(____return_err)
	} 
	return
}

//export SKY_cli_RPCClientFromContext
func SKY_cli_RPCClientFromContext(_c *C.Context__Handle, _client *C.WebRpcClient__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	context, ok := lookupContextHandle(*_c)
	____error_code = SKY_ERROR
	if ok {
		client := cli.RPCClientFromContext(context)
		*_client = registerWebRpcClientHandle( client )
		____error_code = SKY_OK
	} 
	return
}

//export SKY_cli_ConfigFromContext
func SKY_cli_ConfigFromContext(_c *C.Context__Handle, _config *C.Config__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	context, ok := lookupContextHandle(*_c)
	____error_code = SKY_ERROR
	if ok {
		config := cli.ConfigFromContext(context)
		*_config = registerConfigHandle( &config )
		____error_code = SKY_OK
	} 
	return
}
