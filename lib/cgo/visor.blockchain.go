package main

import (
	visor "github.com/skycoin/skycoin/src/visor"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_DefaultWalker
func SKY_visor_DefaultWalker(_hps *C.GoSlice_, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	__arg1 := visor.DefaultWalker(hps)
	return
}

// export SKY_visor_Arbitrating
func SKY_visor_Arbitrating(_enable bool, _arg1 *C.visor__Option) (____error_code uint32) {
	____error_code = 0
	enable := _enable
	__arg1 := visor.Arbitrating(enable)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofOption))
	return
}
