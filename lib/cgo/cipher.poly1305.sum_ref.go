package main

import (
	poly1305 "github.com/skycoin/skycoin/src/poly1305"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_poly1305_Sum
func SKY_poly1305_Sum(_out *[]byte, _msg *C.GoSlice_, _key *[]byte) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	out := *(*[]byte)(unsafe.Pointer(_out))
	msg := *(*[]byte)(unsafe.Pointer(_msg))
	key := *(*[]byte)(unsafe.Pointer(_key))
	poly1305.Sum(out, msg, key)
	return
}
