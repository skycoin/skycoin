package main

import (
	"unsafe"

	poly1305 "github.com/skycoin/skycoin/src/cipher/poly1305"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_poly1305_Verify
func SKY_poly1305_Verify(_mac *C.GoSlice_, _m []byte, _key *C.GoSlice_, _arg3 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	mac := (*[16]byte)(unsafe.Pointer(_mac))
	m := *(*[]byte)(unsafe.Pointer(&_m))
	key := (*[32]byte)(unsafe.Pointer(_key))
	__arg3 := poly1305.Verify(mac, m, key)
	*_arg3 = __arg3
	return
}
