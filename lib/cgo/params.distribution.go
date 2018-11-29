package main

import (
	"reflect"

	params "github.com/skycoin/skycoin/src/params"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_params_GetDistributionAddresses
func SKY_params_GetDistributionAddresses(_arg0 *C.GoSlice_) {
	__arg0 := params.GetDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
}

//export SKY_params_GetUnlockedDistributionAddresses
func SKY_params_GetUnlockedDistributionAddresses(_arg0 *C.GoSlice_) {
	__arg0 := params.GetUnlockedDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
}

//export SKY_params_GetLockedDistributionAddresses
func SKY_params_GetLockedDistributionAddresses(_arg0 *C.GoSlice_) {
	__arg0 := params.GetLockedDistributionAddresses()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
}
