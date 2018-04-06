package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_pex_Peers_ToAddrs
func SKY_pex_Peers_ToAddrs(_ps *C.Peers, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ps := *(*cipher.Peers)(unsafe.Pointer(_ps))
	__arg0 := ps.ToAddrs()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
