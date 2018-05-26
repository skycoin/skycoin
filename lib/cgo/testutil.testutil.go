package main

import (
	"unsafe"

	testutil "github.com/skycoin/skycoin/src/testutil"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_testutil_MakeAddress
func SKY_testutil_MakeAddress(_arg0 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := testutil.MakeAddress()
	*_arg0 = *(*C.cipher__Address)(unsafe.Pointer(&__arg0))
	return
}
