package main

import (
	daemon "github.com/skycoin/skycoin/src/daemon"
	pex "github.com/skycoin/skycoin/src/daemon/pex"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_daemon_NewMessagesConfig
func SKY_daemon_NewMessagesConfig(_arg0 *C.daemon__MessagesConfig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := daemon.NewMessagesConfig()
	*_arg0 = *(*C.daemon__MessagesConfig)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_daemon_MessagesConfig_Register
func SKY_daemon_MessagesConfig_Register(_msc *C.daemon__MessagesConfig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msc := (*daemon.MessagesConfig)(unsafe.Pointer(_msc))
	msc.Register()
	return
}

//export SKY_daemon_NewMessages
func SKY_daemon_NewMessages(_c *C.daemon__MessagesConfig, _arg1 *C.daemon__Messages) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*daemon.MessagesConfig)(unsafe.Pointer(_c))
	__arg1 := daemon.NewMessages(c)
	*_arg1 = *(*C.daemon__Messages)(unsafe.Pointer(__arg1))
	return
}

//export SKY_daemon_NewIPAddr
func SKY_daemon_NewIPAddr(_addr string, _arg1 *C.daemon__IPAddr) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := _addr
	__arg1, ____return_err := daemon.NewIPAddr(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.daemon__IPAddr)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_daemon_IPAddr_String
func SKY_daemon_IPAddr_String(_ipa *C.daemon__IPAddr, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ipa := *(*daemon.IPAddr)(unsafe.Pointer(_ipa))
	__arg0 := ipa.String()
	copyString(__arg0, _arg0)
	return
}

//export SKY_daemon_NewGetPeersMessage
func SKY_daemon_NewGetPeersMessage(_arg0 *C.daemon__GetPeersMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := daemon.NewGetPeersMessage()
	*_arg0 = *(*C.daemon__GetPeersMessage)(unsafe.Pointer(__arg0))
	return
}

//export SKY_daemon_NewGivePeersMessage
func SKY_daemon_NewGivePeersMessage(_peers []C.pex__Peer, _arg1 *C.daemon__GivePeersMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	peers := *(*[]pex.Peer)(unsafe.Pointer(&_peers))
	__arg1 := daemon.NewGivePeersMessage(peers)
	*_arg1 = *(*C.daemon__GivePeersMessage)(unsafe.Pointer(__arg1))
	return
}

//export SKY_daemon_GivePeersMessage_GetPeers
func SKY_daemon_GivePeersMessage_GetPeers(_gpm *C.daemon__GivePeersMessage, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gpm := (*daemon.GivePeersMessage)(unsafe.Pointer(_gpm))
	__arg0 := gpm.GetPeers()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_daemon_NewIntroductionMessage
func SKY_daemon_NewIntroductionMessage(_mirror uint32, _version int32, _port uint16, _arg3 *C.daemon__IntroductionMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	mirror := _mirror
	version := _version
	port := _port
	__arg3 := daemon.NewIntroductionMessage(mirror, version, port)
	*_arg3 = *(*C.daemon__IntroductionMessage)(unsafe.Pointer(__arg3))
	return
}
