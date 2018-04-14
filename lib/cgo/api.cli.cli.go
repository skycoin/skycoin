package main

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	gcli "github.com/urfave/cli"
	"unsafe"
) 

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//TODO: Create specific handle type for config
//export SKY_cli_LoadConfig
func SKY_cli_LoadConfig(_config *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	config, ____return_err := cli.LoadConfig()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_config = (C.Handle)(registerHandle(config))
	}
	return
}

//TODO: Create specific handle type for config
//export SKY_cli_Config_FullWalletPath
func SKY_cli_Config_FullWalletPath(_c *C.Handle, _path *C.GoString_) (____error_code uint32) {
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_c))
	____error_code = SKY_ERROR
	if ok {
		if config, isConfig := (obj).(cli.Config); isConfig {
			path := config.FullWalletPath()
			copyString(path, _path)
			____error_code = SKY_OK
		}
	} 
	return
}

//TODO: Create specific handle type for config
//export SKY_cli_Config_FullDBPath
func SKY_cli_Config_FullDBPath(_c *C.Handle, _path *C.GoString_) (____error_code uint32) {
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_c))
	____error_code = SKY_ERROR
	if ok {
		if config, isConfig := (obj).(cli.Config); isConfig {
			path := config.FullDBPath()
			copyString(path, _path)
			____error_code = SKY_OK
		}
	} 
	return
}

//TODO: Create specific handle type for config
//export SKY_cli_NewApp
func SKY_cli_NewApp(_cfg *C.Handle, _app *C.Handle) (____error_code uint32) {
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_cfg))
	____error_code = SKY_ERROR
	if ok {
		if config, isConfig := (obj).(cli.Config); isConfig {
			app := cli.NewApp(config)
			*_app = (C.Handle)(registerHandle( app ))
			____error_code = SKY_OK
		}
	} 
	return
}

//TODO: Create specific handle type for App
//export SKY_cli_App_Run
func SKY_cli_App_Run(_app *C.Handle, _args *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_app))
	if ok {
		if app, isApp := (obj).(*cli.App); isApp {
			//TODO: stdevEclipse Test this typecast
			args := *(*[]string)(unsafe.Pointer(_args))
			____return_err := app.Run(args)
			____error_code = libErrorCode(____return_err)
		} else {
			____error_code = SKY_ERROR
		}
	} else {
		____error_code = SKY_ERROR
	}
	return
}

//TODO: Create specific handle type for App
//export SKY_cli_RPCClientFromContext
func SKY_cli_RPCClientFromContext(_c *C.Handle, _client *C.WebrpcClient__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_c))
	____error_code = SKY_ERROR
	if ok {
		if context, isContext := (obj).(*gcli.Context); isContext {
			client := cli.RPCClientFromContext(context)
			*_client = registerWebRpcClientHandle( client )
			____error_code = SKY_OK
		}
	} 
	return
}

//TODO: Create specific handle type for config
//export SKY_cli_ConfigFromContext
func SKY_cli_ConfigFromContext(_c *C.Handle, _config *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_c))
	____error_code = SKY_ERROR
	if ok {
		if context, isContext := (obj).(*gcli.Context); isContext {
			config := cli.ConfigFromContext(context)
			*_config = (C.Handle)(registerHandle( config ))
			____error_code = SKY_OK
		}
	} 
	return
}
