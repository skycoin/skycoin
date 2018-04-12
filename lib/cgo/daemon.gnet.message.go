package main

import (
	gnet "github.com/skycoin/skycoin/src/daemon/gnet"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gnet_MessagePrefixFromString
func SKY_gnet_MessagePrefixFromString(_prefix string, _arg1 *C.gnet__MessagePrefix) (____error_code uint32) {
	____error_code = 0
	prefix := _prefix
	__arg1 := gnet.MessagePrefixFromString(prefix)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofMessagePrefix))
	return
}

// export SKY_gnet_RegisterMessage
func SKY_gnet_RegisterMessage(_prefix *C.gnet__MessagePrefix, _msg interface{}) (____error_code uint32) {
	____error_code = 0
	prefix := *(*gnet.MessagePrefix)(unsafe.Pointer(_prefix))
	gnet.RegisterMessage(prefix, msg)
	return
}

// export SKY_gnet_VerifyMessages
func SKY_gnet_VerifyMessages() (____error_code uint32) {
	____error_code = 0
	gnet.VerifyMessages()
	return
}

// export SKY_gnet_EraseMessages
func SKY_gnet_EraseMessages() (____error_code uint32) {
	____error_code = 0
	gnet.EraseMessages()
	return
}
