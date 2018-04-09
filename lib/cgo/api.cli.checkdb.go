package main

import cli "github.com/skycoin/skycoin/src/api/cli"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_IntegrityCheck
func SKY_cli_IntegrityCheck(_db *C.DB, _genesisPubkey *C.PubKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	____return_err := cli.IntegrityCheck(db, genesisPubkey)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
