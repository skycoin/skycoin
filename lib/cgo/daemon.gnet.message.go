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
func SKY_gnet_MessagePrefixFromString(_prefix string, _arg1 *C.MessagePrefix) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	prefix := _prefix
	__arg1 := gnet.MessagePrefixFromString(prefix)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofMessagePrefix))
	return
}

// export SKY_gnet_NewMessageContext
func SKY_gnet_NewMessageContext(_conn *C.Connection, _arg1 *C.MessageContext) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	conn := (*Connection)(unsafe.Pointer(_conn))
	__arg1 := gnet.NewMessageContext(conn)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofMessageContext))
	return
}

// export SKY_gnet_RegisterMessage
func SKY_gnet_RegisterMessage(_prefix *C.MessagePrefix, _msg interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	prefix := *(*MessagePrefix)(unsafe.Pointer(_prefix))
	gnet.RegisterMessage(prefix, msg)
	return
}

// export SKY_gnet_VerifyMessages
func SKY_gnet_VerifyMessages() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gnet.VerifyMessages()
	return
}

// export SKY_gnet_EraseMessages
func SKY_gnet_EraseMessages() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gnet.EraseMessages()
	return
}
