package main

import (
	gnet "github.com/skycoin/skycoin/src/daemon/gnet"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_gnet_MessagePrefixFromString
func SKY_gnet_MessagePrefixFromString(_prefix string, _arg1 *C.gnet__MessagePrefix) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	prefix := _prefix
	__arg1 := gnet.MessagePrefixFromString(prefix)
	*_arg1 = *(*C.gnet__MessagePrefix)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_gnet_VerifyMessages
func SKY_gnet_VerifyMessages() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gnet.VerifyMessages()
	return
}

//export SKY_gnet_EraseMessages
func SKY_gnet_EraseMessages() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gnet.EraseMessages()
	return
}
