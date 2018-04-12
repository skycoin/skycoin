package main

import (
	poly1305 "github.com/skycoin/skycoin/src/cipher/poly1305"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_poly1305_Verify
func SKY_poly1305_Verify(_mac *[16]byte, _m []byte, _key *[32]byte, _arg3 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	//TODO: stdevEclipse Check Pointer casting
	__arg3 := poly1305.Verify(_mac, _m, _key)
	*_arg3 = __arg3
	return
}
