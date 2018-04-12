package main

import utc "github.com/skycoin/skycoin/src/util/utc"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_utc_UnixNow
func SKY_utc_UnixNow(_arg0 *int64) (____error_code uint32) {
	____error_code = 0
	__arg0 := utc.UnixNow()
	*_arg0 = __arg0
	return
}
