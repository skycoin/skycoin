package main

import visor "github.com/skycoin/skycoin/src/visor"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_OpenDB
func SKY_visor_OpenDB(_dbFile string, _readOnly bool, _arg2 *C.DB) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dbFile := _dbFile
	readOnly := _readOnly
	__arg2, ____return_err := visor.OpenDB(dbFile, readOnly)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
