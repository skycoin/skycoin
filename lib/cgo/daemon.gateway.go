package main

import (
	daemon "github.com/skycoin/skycoin/src/daemon"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_daemon_NewGatewayConfig
func SKY_daemon_NewGatewayConfig(_arg0 *C.daemon__GatewayConfig) (____error_code uint32) {
	____error_code = 0
	__arg0 := daemon.NewGatewayConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofGatewayConfig))
	return
}

// export SKY_daemon_FbyAddressesNotIncluded
func SKY_daemon_FbyAddressesNotIncluded(_addrs *C.GoSlice_, _arg1 *C.daemon__OutputsFilter) (____error_code uint32) {
	____error_code = 0
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1 := daemon.FbyAddressesNotIncluded(addrs)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofOutputsFilter))
	return
}

// export SKY_daemon_FbyAddresses
func SKY_daemon_FbyAddresses(_addrs *C.GoSlice_, _arg1 *C.daemon__OutputsFilter) (____error_code uint32) {
	____error_code = 0
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1 := daemon.FbyAddresses(addrs)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofOutputsFilter))
	return
}

// export SKY_daemon_FbyHashes
func SKY_daemon_FbyHashes(_hashes *C.GoSlice_, _arg1 *C.daemon__OutputsFilter) (____error_code uint32) {
	____error_code = 0
	hashes := *(*[]string)(unsafe.Pointer(_hashes))
	__arg1 := daemon.FbyHashes(hashes)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofOutputsFilter))
	return
}

// export SKY_daemon_MakeSearchMap
func SKY_daemon_MakeSearchMap(_addrs *C.GoSlice_, _arg1 map[string]struct{}) (____error_code uint32) {
	____error_code = 0
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1 := daemon.MakeSearchMap(addrs)
	return
}

// export SKY_daemon_spendValidator_HasUnconfirmedSpendTx
func SKY_daemon_spendValidator_HasUnconfirmedSpendTx(_sv *C.daemon__spendValidator, _addr *C.GoSlice_, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	sv := *(*daemon.spendValidator)(unsafe.Pointer(_sv))
	__arg1, ____return_err := sv.HasUnconfirmedSpendTx(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}
