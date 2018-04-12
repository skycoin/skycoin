package main

import (
	daemon "github.com/skycoin/skycoin/src/daemon"
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

// export SKY_daemon_NewMessageConfig
func SKY_daemon_NewMessageConfig(_prefix string, _m interface{}, _arg2 *C.daemon__MessageConfig) (____error_code uint32) {
	____error_code = 0
	prefix := _prefix
	__arg2 := daemon.NewMessageConfig(prefix, m)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofMessageConfig))
	return
}

// export SKY_daemon_NewMessagesConfig
func SKY_daemon_NewMessagesConfig(_arg0 *C.daemon__MessagesConfig) (____error_code uint32) {
	____error_code = 0
	__arg0 := daemon.NewMessagesConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofMessagesConfig))
	return
}

// export SKY_daemon_MessagesConfig_Register
func SKY_daemon_MessagesConfig_Register(_msc *C.daemon__MessagesConfig) (____error_code uint32) {
	____error_code = 0
	msc := (*daemon.MessagesConfig)(unsafe.Pointer(_msc))
	msc.Register()
	return
}

// export SKY_daemon_NewMessages
func SKY_daemon_NewMessages(_c *C.daemon__MessagesConfig, _arg1 *C.daemon__Messages) (____error_code uint32) {
	____error_code = 0
	c := *(*daemon.MessagesConfig)(unsafe.Pointer(_c))
	__arg1 := daemon.NewMessages(c)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofMessages))
	return
}

// export SKY_daemon_NewIPAddr
func SKY_daemon_NewIPAddr(_addr string, _arg1 *C.daemon__IPAddr) (____error_code uint32) {
	____error_code = 0
	addr := _addr
	__arg1, ____return_err := daemon.NewIPAddr(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofIPAddr))
	}
	return
}

// export SKY_daemon_IPAddr_String
func SKY_daemon_IPAddr_String(_ipa *C.daemon__IPAddr, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	ipa := *(*daemon.IPAddr)(unsafe.Pointer(_ipa))
	__arg0 := ipa.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_daemon_NewGetPeersMessage
func SKY_daemon_NewGetPeersMessage(_arg0 *C.daemon__GetPeersMessage) (____error_code uint32) {
	____error_code = 0
	__arg0 := daemon.NewGetPeersMessage()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofGetPeersMessage))
	return
}

// export SKY_daemon_GetPeersMessage_Handle
func SKY_daemon_GetPeersMessage_Handle(_gpm *C.daemon__GetPeersMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	gpm := (*daemon.GetPeersMessage)(unsafe.Pointer(_gpm))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := gpm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_NewGivePeersMessage
func SKY_daemon_NewGivePeersMessage(_peers *C.GoSlice_, _arg1 *C.daemon__GivePeersMessage) (____error_code uint32) {
	____error_code = 0
	__arg1 := daemon.NewGivePeersMessage(peers)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofGivePeersMessage))
	return
}

// export SKY_daemon_GivePeersMessage_GetPeers
func SKY_daemon_GivePeersMessage_GetPeers(_gpm *C.daemon__GivePeersMessage, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	gpm := (*daemon.GivePeersMessage)(unsafe.Pointer(_gpm))
	__arg0 := gpm.GetPeers()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_GivePeersMessage_Handle
func SKY_daemon_GivePeersMessage_Handle(_gpm *C.daemon__GivePeersMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	gpm := (*daemon.GivePeersMessage)(unsafe.Pointer(_gpm))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := gpm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_NewIntroductionMessage
func SKY_daemon_NewIntroductionMessage(_mirror uint32, _version int32, _port uint16, _arg3 *C.daemon__IntroductionMessage) (____error_code uint32) {
	____error_code = 0
	mirror := _mirror
	version := _version
	port := _port
	__arg3 := daemon.NewIntroductionMessage(mirror, version, port)
	copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofIntroductionMessage))
	return
}

// export SKY_daemon_IntroductionMessage_Handle
func SKY_daemon_IntroductionMessage_Handle(_intro *C.daemon__IntroductionMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	intro := (*daemon.IntroductionMessage)(unsafe.Pointer(_intro))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := intro.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_PingMessage_Handle
func SKY_daemon_PingMessage_Handle(_ping *C.daemon__PingMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	ping := (*daemon.PingMessage)(unsafe.Pointer(_ping))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := ping.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_PongMessage_Handle
func SKY_daemon_PongMessage_Handle(_pong *C.daemon__PongMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	pong := (*daemon.PongMessage)(unsafe.Pointer(_pong))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := pong.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
