package main

import testutil "github.com/skycoin/skycoin/src/testutil"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_testutil_MakeAddress
func SKY_testutil_MakeAddress(_arg0 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	__arg0 := testutil.MakeAddress()
	return
}
