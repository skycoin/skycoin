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

// export SKY_poly1305_Sum_arm
func SKY_poly1305_Sum_arm(_out *[16]byte, _m []byte, _key *[32]byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	//TODO: stdevEclipse Check Pointer Casting
	poly1305.Sum(_out, _m, _key)
	return
}
