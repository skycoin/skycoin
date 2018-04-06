package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	gnet "github.com/skycoin/skycoin/src/gnet"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gnet_MessagePrefixFromString
func SKY_gnet_MessagePrefixFromString(_prefix string, _arg1 *C.MessagePrefix) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	prefix := _prefix
	__arg1 := gnet.MessagePrefixFromString(prefix)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofMessagePrefix))
	return
}

// export SKY_gnet_NewMessageContext
func SKY_gnet_NewMessageContext(_conn *C.Connection, _arg1 *C.MessageContext) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	conn := (*cipher.Connection)(unsafe.Pointer(_conn))
	__arg1 := gnet.NewMessageContext(conn)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofMessageContext))
	return
}

// export SKY_gnet_RegisterMessage
func SKY_gnet_RegisterMessage(_prefix *C.MessagePrefix, _msg interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	prefix := *(*cipher.MessagePrefix)(unsafe.Pointer(_prefix))
	gnet.RegisterMessage(prefix, msg)
	return
}

// export SKY_gnet_VerifyMessages
func SKY_gnet_VerifyMessages() (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	gnet.VerifyMessages()
	return
}

// export SKY_gnet_EraseMessages
func SKY_gnet_EraseMessages() (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	gnet.EraseMessages()
	return
}
