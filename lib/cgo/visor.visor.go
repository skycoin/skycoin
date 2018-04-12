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

// export SKY_visor_MaxDropletDivisor
func SKY_visor_MaxDropletDivisor(_arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	__arg0 := visor.MaxDropletDivisor()
	*_arg0 = __arg0
	return
}

// export SKY_visor_DropletPrecisionCheck
func SKY_visor_DropletPrecisionCheck(_amount uint64) (____error_code uint32) {
	____error_code = 0
	amount := _amount
	____return_err := visor.DropletPrecisionCheck(amount)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_baseFilter_Match
func SKY_visor_baseFilter_Match(_f *C.visor__baseFilter, _tx *C.visor__Transaction, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	f := *(*visor.baseFilter)(unsafe.Pointer(_f))
	tx := (*visor.Transaction)(unsafe.Pointer(_tx))
	__arg1 := f.Match(tx)
	*_arg1 = __arg1
	return
}

// export SKY_visor_AddrsFilter
func SKY_visor_AddrsFilter(_addrs *C.GoSlice_, _arg1 *C.visor__TxFilter) (____error_code uint32) {
	____error_code = 0
	__arg1 := visor.AddrsFilter(addrs)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTxFilter))
	return
}

// export SKY_visor_addrsFilter_Match
func SKY_visor_addrsFilter_Match(_af *C.visor__addrsFilter, _tx *C.visor__Transaction, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	af := *(*visor.addrsFilter)(unsafe.Pointer(_af))
	tx := (*visor.Transaction)(unsafe.Pointer(_tx))
	__arg1 := af.Match(tx)
	*_arg1 = __arg1
	return
}

// export SKY_visor_ConfirmedTxFilter
func SKY_visor_ConfirmedTxFilter(_isConfirmed bool, _arg1 *C.visor__TxFilter) (____error_code uint32) {
	____error_code = 0
	isConfirmed := _isConfirmed
	__arg1 := visor.ConfirmedTxFilter(isConfirmed)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTxFilter))
	return
}

// export SKY_visor_ToAddresses
func SKY_visor_ToAddresses(_addresses *C.GoSlice_, _arg1 C.Handle) (____error_code uint32) {
	____error_code = 0
	__arg1 := visor.ToAddresses(addresses)
	return
}
