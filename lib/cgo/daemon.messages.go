package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	daemon "github.com/skycoin/skycoin/src/daemon"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_daemon_NewMessageConfig
func SKY_daemon_NewMessageConfig(_prefix string, _m interface{}, _arg2 *C.MessageConfig) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	prefix := _prefix
	__arg2 := daemon.NewMessageConfig(prefix, m)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofMessageConfig))
	return
}

// export SKY_daemon_NewMessagesConfig
func SKY_daemon_NewMessagesConfig(_arg0 *C.MessagesConfig) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewMessagesConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofMessagesConfig))
	return
}

// export SKY_daemon_MessagesConfig_Register
func SKY_daemon_MessagesConfig_Register(_msc *C.MessagesConfig) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	msc := (*cipher.MessagesConfig)(unsafe.Pointer(_msc))
	msc.Register()
	return
}

// export SKY_daemon_NewMessages
func SKY_daemon_NewMessages(_c *C.MessagesConfig, _arg1 *C.Messages) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := *(*cipher.MessagesConfig)(unsafe.Pointer(_c))
	__arg1 := daemon.NewMessages(c)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofMessages))
	return
}

// export SKY_daemon_NewIPAddr
func SKY_daemon_NewIPAddr(_addr string, _arg1 *C.IPAddr) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	addr := _addr
	__arg1, ____return_err := daemon.NewIPAddr(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofIPAddr))
	}
	return
}

// export SKY_daemon_IPAddr_String
func SKY_daemon_IPAddr_String(_ipa *C.IPAddr, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ipa := *(*cipher.IPAddr)(unsafe.Pointer(_ipa))
	__arg0 := ipa.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_daemon_NewGetPeersMessage
func SKY_daemon_NewGetPeersMessage(_arg0 *C.GetPeersMessage) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewGetPeersMessage()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofGetPeersMessage))
	return
}

// export SKY_daemon_GetPeersMessage_Handle
func SKY_daemon_GetPeersMessage_Handle(_gpm *C.GetPeersMessage, _mc *C.MessageContext, _daemon interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	gpm := (*cipher.GetPeersMessage)(unsafe.Pointer(_gpm))
	____return_err := gpm.Handle(mc, daemon)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_GetPeersMessage_Process
func SKY_daemon_GetPeersMessage_Process(_gpm *C.GetPeersMessage, _d *C.Daemon) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	gpm := (*cipher.GetPeersMessage)(unsafe.Pointer(_gpm))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	gpm.Process(d)
	return
}

// export SKY_daemon_NewGivePeersMessage
func SKY_daemon_NewGivePeersMessage(_peers *C.GoSlice_, _arg1 *C.GivePeersMessage) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1 := daemon.NewGivePeersMessage(peers)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofGivePeersMessage))
	return
}

// export SKY_daemon_GivePeersMessage_GetPeers
func SKY_daemon_GivePeersMessage_GetPeers(_gpm *C.GivePeersMessage, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	gpm := (*cipher.GivePeersMessage)(unsafe.Pointer(_gpm))
	__arg0 := gpm.GetPeers()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_GivePeersMessage_Handle
func SKY_daemon_GivePeersMessage_Handle(_gpm *C.GivePeersMessage, _mc *C.MessageContext, _daemon interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	gpm := (*cipher.GivePeersMessage)(unsafe.Pointer(_gpm))
	____return_err := gpm.Handle(mc, daemon)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_GivePeersMessage_Process
func SKY_daemon_GivePeersMessage_Process(_gpm *C.GivePeersMessage, _d *C.Daemon) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	gpm := (*cipher.GivePeersMessage)(unsafe.Pointer(_gpm))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	gpm.Process(d)
	return
}

// export SKY_daemon_NewIntroductionMessage
func SKY_daemon_NewIntroductionMessage(_mirror uint32, _version int32, _port uint16, _arg3 *C.IntroductionMessage) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	mirror := _mirror
	version := _version
	port := _port
	__arg3 := daemon.NewIntroductionMessage(mirror, version, port)
	copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofIntroductionMessage))
	return
}

// export SKY_daemon_IntroductionMessage_Handle
func SKY_daemon_IntroductionMessage_Handle(_intro *C.IntroductionMessage, _mc *C.MessageContext, _daemon interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	intro := (*cipher.IntroductionMessage)(unsafe.Pointer(_intro))
	____return_err := intro.Handle(mc, daemon)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_IntroductionMessage_Process
func SKY_daemon_IntroductionMessage_Process(_intro *C.IntroductionMessage, _d *C.Daemon) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	intro := (*cipher.IntroductionMessage)(unsafe.Pointer(_intro))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	intro.Process(d)
	return
}

// export SKY_daemon_PingMessage_Handle
func SKY_daemon_PingMessage_Handle(_ping *C.PingMessage, _mc *C.MessageContext, _daemon interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ping := (*cipher.PingMessage)(unsafe.Pointer(_ping))
	____return_err := ping.Handle(mc, daemon)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_PingMessage_Process
func SKY_daemon_PingMessage_Process(_ping *C.PingMessage, _d *C.Daemon) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ping := (*cipher.PingMessage)(unsafe.Pointer(_ping))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	ping.Process(d)
	return
}

// export SKY_daemon_PongMessage_Handle
func SKY_daemon_PongMessage_Handle(_pong *C.PongMessage, _mc *C.MessageContext, _daemon interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pong := (*cipher.PongMessage)(unsafe.Pointer(_pong))
	____return_err := pong.Handle(mc, daemon)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
