package main

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_LoadConfig
func SKY_cli_LoadConfig(_arg0 *C.Config) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0, ____return_err := cli.LoadConfig()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofConfig))
	}
	return
}

// export SKY_cli_Config_FullWalletPath
func SKY_cli_Config_FullWalletPath(_c *C.Config, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*Config)(unsafe.Pointer(_c))
	__arg0 := c.FullWalletPath()
	copyString(__arg0, _arg0)
	return
}

// export SKY_cli_Config_FullDBPath
func SKY_cli_Config_FullDBPath(_c *C.Config, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*Config)(unsafe.Pointer(_c))
	__arg0 := c.FullDBPath()
	copyString(__arg0, _arg0)
	return
}

// export SKY_cli_NewApp
func SKY_cli_NewApp(_cfg *C.Config, _arg1 *C.App) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	cfg := *(*Config)(unsafe.Pointer(_cfg))
	__arg1 := cli.NewApp(cfg)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofApp))
	return
}

// export SKY_cli_App_Run
func SKY_cli_App_Run(_app *C.App, _args *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	app := (*App)(unsafe.Pointer(_app))
	args := *(*[]string)(unsafe.Pointer(_args))
	____return_err := app.Run(args)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_cli_RPCClientFromContext
func SKY_cli_RPCClientFromContext(_c *C.Context, _arg1 *C.Client) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := cli.RPCClientFromContext(c)
	return
}

// export SKY_cli_ConfigFromContext
func SKY_cli_ConfigFromContext(_c *C.Context, _arg1 *C.Config) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := cli.ConfigFromContext(c)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofConfig))
	return
}
