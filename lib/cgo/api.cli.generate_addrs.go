package main

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	cipher "github.com/skycoin/skycoin/src/cipher"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_cli_GenerateAddressesInFile
func SKY_cli_GenerateAddressesInFile(_walletFile string, _num uint64, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	walletFile := _walletFile
	num := _num
	__arg2, ____return_err := cli.GenerateAddressesInFile(walletFile, num)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

//export SKY_cli_FormatAddressesAsJSON
func SKY_cli_FormatAddressesAsJSON(_addrs *C.GoSlice_, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]cipher.Address)(unsafe.Pointer(_addrs))
	__arg1, ____return_err := cli.FormatAddressesAsJSON(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

//export SKY_cli_FormatAddressesAsJoinedArray
func SKY_cli_FormatAddressesAsJoinedArray(_addrs *C.GoSlice_, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]cipher.Address)(unsafe.Pointer(_addrs))
	__arg1 := cli.FormatAddressesAsJoinedArray(addrs)
	copyString(__arg1, _arg1)
	return
}

//export SKY_cli_AddressesToStrings
func SKY_cli_AddressesToStrings(_addrs *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]cipher.Address)(unsafe.Pointer(_addrs))
	__arg1 := cli.AddressesToStrings(addrs)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}
