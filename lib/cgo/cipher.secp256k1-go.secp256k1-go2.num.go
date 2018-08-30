package main

import (
	"encoding/hex"

	secp256k1go2 "github.com/skycoin/skycoin/src/cipher/secp256k1-go/secp256k1-go2"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_secp256k1go_Number_Create
func SKY_secp256k1go_Number_Create(handle *C.Number_Handle) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	var num secp256k1go2.Number
	*handle = registerNumberHandle(&num)
	return
}

//export SKY_secp256k1go_Number_Print
func SKY_secp256k1go_Number_Print(handle C.Number_Handle, _label string) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	num, ok := lookupNumberHandle(handle)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	num.Print(_label)
	return
}

//export SKY_secp256k1go_Number_SetHex
func SKY_secp256k1go_Number_SetHex(handle C.Number_Handle, _s string) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	num, ok := lookupNumberHandle(handle)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	num.SetHex(_s)
	return
}

//export SKY_secp256k1go_Number_IsOdd
func SKY_secp256k1go_Number_IsOdd(handle C.Number_Handle, _arg0 *bool) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	num, ok := lookupNumberHandle(handle)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	__arg0 := num.IsOdd()
	*_arg0 = __arg0
	return
}

//export SKY_secp256k1go_Number_IsEqual
func SKY_secp256k1go_Number_IsEqual(handle1 C.Number_Handle, handle2 C.Number_Handle, result *bool) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	num1, ok := lookupNumberHandle(handle1)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	num2, ok := lookupNumberHandle(handle2)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	*result = hex.EncodeToString(num1.Bytes()) == hex.EncodeToString(num2.Bytes())
	return
}
