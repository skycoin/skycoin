package main

import cli "github.com/skycoin/skycoin/src/cli"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_IntegrityCheck
func SKY_cli_IntegrityCheck(_db *C.DB, _genesisPubkey *C.PubKey) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	____return_err := cli.IntegrityCheck(db, genesisPubkey)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
