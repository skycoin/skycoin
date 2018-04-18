package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	visor "github.com/skycoin/skycoin/src/visor"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_visor_MaxDropletDivisor
func SKY_visor_MaxDropletDivisor(_arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := visor.MaxDropletDivisor()
	*_arg0 = __arg0
	return
}

//export SKY_visor_DropletPrecisionCheck
func SKY_visor_DropletPrecisionCheck(_amount uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	amount := _amount
	____return_err := visor.DropletPrecisionCheck(amount)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_visor_AddrsFilter
func SKY_visor_AddrsFilter(_addrs []C.cipher__Address, _arg1 *C.visor__TxFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]cipher.Address)(unsafe.Pointer(&_addrs))
	__arg1 := visor.AddrsFilter(addrs)
	*_arg1 = *(*C.visor__TxFilter)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_visor_ConfirmedTxFilter
func SKY_visor_ConfirmedTxFilter(_isConfirmed bool, _arg1 *C.visor__TxFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	isConfirmed := _isConfirmed
	__arg1 := visor.ConfirmedTxFilter(isConfirmed)
	*_arg1 = *(*C.visor__TxFilter)(unsafe.Pointer(&__arg1))
	return
}
