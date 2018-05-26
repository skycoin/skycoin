package main

import (
	"unsafe"

	secp256k1go2 "github.com/skycoin/skycoin/src/cipher/secp256k1-go/secp256k1-go2"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_secp256k1go_Number_Print
func SKY_secp256k1go_Number_Print(_num *C.Number, _label string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	num := (*secp256k1go2.Number)(unsafe.Pointer(_num))
	label := _label
	num.Print(label)
	return
}

//export SKY_secp256k1go_Number_SetHex
func SKY_secp256k1go_Number_SetHex(_num *C.Number, _s string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	num := (*secp256k1go2.Number)(unsafe.Pointer(_num))
	s := _s
	num.SetHex(s)
	return
}

//export SKY_secp256k1go_Number_IsOdd
func SKY_secp256k1go_Number_IsOdd(_num *C.Number, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	num := (*secp256k1go2.Number)(unsafe.Pointer(_num))
	__arg0 := num.IsOdd()
	*_arg0 = __arg0
	return
}
