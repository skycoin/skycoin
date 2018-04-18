package main

import (
	poly1305 "github.com/skycoin/skycoin/src/cipher/poly1305"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_poly1305_Sum
func SKY_poly1305_Sum(_out *C.GoSlice_, _msg []byte, _key *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	out := (*[16]byte)(unsafe.Pointer(_out))
	msg := *(*[]byte)(unsafe.Pointer(&_msg))
	key := (*[32]byte)(unsafe.Pointer(_key))
	poly1305.Sum(out, msg, key)
	return
}
