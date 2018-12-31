package main

import (
	"os"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cli_CLI_Run
func SKY_cli_CLI_Run(_app C.CLI__Handle) (____error_code uint32) {
	cli, okapp := lookupCLIHandle(_app)
	if !okapp {
		____error_code = SKY_BAD_HANDLE
		return
	}

	____return_err := cli.Execute()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_cli_Config_GetCoin
func SKY_cli_Config_GetCoin(_c C.Config__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	__c, okc := lookupConfigHandle(_c)
	if !okc {
		____error_code = SKY_BAD_HANDLE
		return
	}
	c := *__c
	__arg0 := c.Coin
	copyString(__arg0, _arg0)
	return
}

//export SKY_cli_Config_GetRPCAddress
func SKY_cli_Config_GetRPCAddress(_c C.Config__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	__c, okc := lookupConfigHandle(_c)
	if !okc {
		____error_code = SKY_BAD_HANDLE
		return
	}
	c := *__c
	__arg0 := c.RPCAddress
	copyString(__arg0, _arg0)
	return
}

//export SKY_cli_Getenv
func SKY_cli_Getenv(varname string, _arg0 *C.GoString_) (____error_code uint32) {
	__arg0 := os.Getenv(varname)
	copyString(__arg0, _arg0)
	return
}

//export SKY_cli_Setenv
func SKY_cli_Setenv(varname string, value string) (____error_code uint32) {
	____return_err := os.Setenv(varname, value)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
