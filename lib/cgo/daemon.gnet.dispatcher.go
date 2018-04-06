package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	gnet "github.com/skycoin/skycoin/src/gnet"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gnet_EncodeMessage
func SKY_gnet_EncodeMessage(_msg *C.Message, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	msg := *(*cipher.Message)(unsafe.Pointer(_msg))
	__arg1 := gnet.EncodeMessage(msg)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}
