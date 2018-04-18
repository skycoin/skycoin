package main

import (
	daemon "github.com/skycoin/skycoin/src/daemon"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_daemon_NewGatewayConfig
func SKY_daemon_NewGatewayConfig(_arg0 *C.daemon__GatewayConfig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := daemon.NewGatewayConfig()
	*_arg0 = *(*C.daemon__GatewayConfig)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_daemon_FbyAddressesNotIncluded
func SKY_daemon_FbyAddressesNotIncluded(_addrs []string, _arg1 *C.daemon__OutputsFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1 := daemon.FbyAddressesNotIncluded(addrs)
	*_arg1 = *(*C.daemon__OutputsFilter)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_daemon_FbyAddresses
func SKY_daemon_FbyAddresses(_addrs []string, _arg1 *C.daemon__OutputsFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1 := daemon.FbyAddresses(addrs)
	*_arg1 = *(*C.daemon__OutputsFilter)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_daemon_FbyHashes
func SKY_daemon_FbyHashes(_hashes []string, _arg1 *C.daemon__OutputsFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hashes := *(*[]string)(unsafe.Pointer(&_hashes))
	__arg1 := daemon.FbyHashes(hashes)
	*_arg1 = *(*C.daemon__OutputsFilter)(unsafe.Pointer(&__arg1))
	return
}
