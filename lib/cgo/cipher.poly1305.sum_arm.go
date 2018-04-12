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

// export SKY_poly1305_Sum
func SKY_poly1305_Sum(_out *[]byte, _m *C.GoSlice_, _key *[]byte) (____error_code uint32) {
	____error_code = 0
	out := *(*[]byte)(unsafe.Pointer(_out))
	m := *(*[]byte)(unsafe.Pointer(_m))
	key := *(*[]byte)(unsafe.Pointer(_key))
	poly1305.Sum(out, m, key)
	return
}
