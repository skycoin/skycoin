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

// export SKY_poly1305_Verify
func SKY_poly1305_Verify(_mac *[]byte, _m *C.GoSlice_, _key *[]byte, _arg3 *bool) (____error_code uint32) {
	____error_code = 0
	mac := *(*[]byte)(unsafe.Pointer(_mac))
	m := *(*[]byte)(unsafe.Pointer(_m))
	key := *(*[]byte)(unsafe.Pointer(_key))
	__arg3 := poly1305.Verify(mac, m, key)
	*_arg3 = __arg3
	return
}
