package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_secp256k1go_Number_Print
func SKY_secp256k1go_Number_Print(_num *C.Number, _label string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	num := (*cipher.Number)(unsafe.Pointer(_num))
	label := _label
	num.Print(label)
	return
}

// export SKY_secp256k1go_Number_SetHex
func SKY_secp256k1go_Number_SetHex(_num *C.Number, _s string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	num := (*cipher.Number)(unsafe.Pointer(_num))
	s := _s
	num.SetHex(s)
	return
}

// export SKY_secp256k1go_Number_IsOdd
func SKY_secp256k1go_Number_IsOdd(_num *C.Number, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	num := (*cipher.Number)(unsafe.Pointer(_num))
	__arg0 := num.IsOdd()
	*_arg0 = __arg0
	return
}
