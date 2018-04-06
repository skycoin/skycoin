package main

import utc "github.com/skycoin/skycoin/src/utc"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_utc_Now
func SKY_utc_Now(_arg0 *C.Time) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := utc.Now()
	return
}

// export SKY_utc_UnixNow
func SKY_utc_UnixNow(_arg0 *int64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := utc.UnixNow()
	*_arg0 = __arg0
	return
}
