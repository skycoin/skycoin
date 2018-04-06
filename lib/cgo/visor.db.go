package main

import visor "github.com/skycoin/skycoin/src/visor"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_OpenDB
func SKY_visor_OpenDB(_dbFile string, _readOnly bool, _arg2 *C.DB) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dbFile := _dbFile
	readOnly := _readOnly
	__arg2, ____return_err := visor.OpenDB(dbFile, readOnly)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
