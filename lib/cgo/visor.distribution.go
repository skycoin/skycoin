package main

import (
	visor "github.com/skycoin/skycoin/src/visor"
	reflect "reflect"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_GetDistributionAddresses
func SKY_visor_GetDistributionAddresses(_arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := visor.GetDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_visor_GetUnlockedDistributionAddresses
func SKY_visor_GetUnlockedDistributionAddresses(_arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := visor.GetUnlockedDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_visor_GetLockedDistributionAddresses
func SKY_visor_GetLockedDistributionAddresses(_arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := visor.GetLockedDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_visor_TransactionIsLocked
func SKY_visor_TransactionIsLocked(_inUxs *C.UxArray, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1 := visor.TransactionIsLocked(inUxs)
	*_arg1 = __arg1
	return
}
