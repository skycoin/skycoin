package main

import (
	coin "github.com/skycoin/skycoin/src/coin"
	visor "github.com/skycoin/skycoin/src/visor"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_visor_DefaultWalker
func SKY_visor_DefaultWalker(_hps []C.coin__HashPair, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hps := *(*[]coin.HashPair)(unsafe.Pointer(&_hps))
	__arg1 := visor.DefaultWalker(hps)
	*_arg1 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_visor_Arbitrating
func SKY_visor_Arbitrating(_enable bool, _arg1 *C.visor__Option) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	enable := _enable
	__arg1 := visor.Arbitrating(enable)
	*_arg1 = *(*C.visor__Option)(unsafe.Pointer(&__arg1))
	return
}
