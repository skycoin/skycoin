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
import "unsafe"

//export SKY_poly1305_Sum_amd64
func SKY_poly1305_Sum_amd64(__out *C.poly1305__Mac, __m []byte, __key *C.poly1305__Key) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	//TODO: stdevEclipse Check Pointer casting
	_out := (*[16]byte)(unsafe.Pointer(__out))
	_key := (*[32]byte)(unsafe.Pointer(__key))
	poly1305.Sum(_out, __m, _key)
	return
}
