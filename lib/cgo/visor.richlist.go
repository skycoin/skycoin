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

// export SKY_visor_NewRichlist
func SKY_visor_NewRichlist(_allAccounts map[string]uint64, _lockedAddrs map[string]struct{}, _arg2 *C.visor__Richlist) (____error_code uint32) {
	____error_code = 0
	__arg2, ____return_err := visor.NewRichlist(allAccounts, lockedAddrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofRichlist))
	}
	return
}

// export SKY_visor_Richlist_FilterAddresses
func SKY_visor_Richlist_FilterAddresses(_r *C.visor__Richlist, _addrs map[string]struct{}, _arg1 *C.visor__Richlist) (____error_code uint32) {
	____error_code = 0
	r := *(*visor.Richlist)(unsafe.Pointer(_r))
	__arg1 := r.FilterAddresses(addrs)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofRichlist))
	return
}
