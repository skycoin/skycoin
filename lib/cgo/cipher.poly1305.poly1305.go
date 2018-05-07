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

//export SKY_poly1305_Verify
func SKY_poly1305_Verify(__mac *C.poly1305__Mac, __m []byte, __key *C.poly1305__Key, _arg3 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	_mac := (*[16]byte)(unsafe.Pointer(__mac))
	_key := (*[32]byte)(unsafe.Pointer(__key))
	__arg3 := poly1305.Verify(_mac, __m, _key)
	*_arg3 = __arg3
	return
}
