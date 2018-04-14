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
func SKY_poly1305_Sum(_out *C.poly1305__Mac, _msg []byte, _key *C.poly1305__Key) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	//TODO: stdevEclipse Check Pointer Casting
	out := (*[16]byte)(unsafe.Pointer(_out))
	key := (*[32]byte)(unsafe.Pointer(_key))
	poly1305.Sum(out, _msg, key)
	return
}
