package main

import (
	cli "github.com/skycoin/skycoin/src/cli"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cli_LoadConfig
func SKY_cli_LoadConfig(_arg0 *C.Config__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0, ____return_err := cli.LoadConfig()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerConfigHandle(&__arg0)
	}
	return
}

// TODO: RpcClientFromConfig
