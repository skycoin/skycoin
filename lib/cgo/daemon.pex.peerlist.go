package main

import (
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_pex_Peers_ToAddrs
func SKY_pex_Peers_ToAddrs(_ps *C.pex__Peers, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	ps := *(*pex.Peers)(unsafe.Pointer(_ps))
	__arg0 := ps.ToAddrs()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
