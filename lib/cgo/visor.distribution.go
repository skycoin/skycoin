package main

import (
	coin "github.com/skycoin/skycoin/src/coin"
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

//export SKY_visor_GetDistributionAddresses
func SKY_visor_GetDistributionAddresses(_arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := visor.GetDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_visor_GetUnlockedDistributionAddresses
func SKY_visor_GetUnlockedDistributionAddresses(_arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := visor.GetUnlockedDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_visor_GetLockedDistributionAddresses
func SKY_visor_GetLockedDistributionAddresses(_arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := visor.GetLockedDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_visor_TransactionIsLocked
func SKY_visor_TransactionIsLocked(_inUxs *C.coin__UxArray, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	inUxs := *(*coin.UxArray)(unsafe.Pointer(_inUxs))
	__arg1 := visor.TransactionIsLocked(inUxs)
	*_arg1 = __arg1
	return
}
