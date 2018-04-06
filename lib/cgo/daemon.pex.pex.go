package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	pex "github.com/skycoin/skycoin/src/pex"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_pex_NewPeer
func SKY_pex_NewPeer(_address string, _arg1 *C.Peer) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	address := _address
	__arg1 := pex.NewPeer(address)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofPeer))
	return
}

// export SKY_pex_Peer_Seen
func SKY_pex_Peer_Seen(_peer *C.Peer) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	peer := (*cipher.Peer)(unsafe.Pointer(_peer))
	peer.Seen()
	return
}

// export SKY_pex_Peer_IncreaseRetryTimes
func SKY_pex_Peer_IncreaseRetryTimes(_peer *C.Peer) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	peer := (*cipher.Peer)(unsafe.Pointer(_peer))
	peer.IncreaseRetryTimes()
	return
}

// export SKY_pex_Peer_ResetRetryTimes
func SKY_pex_Peer_ResetRetryTimes(_peer *C.Peer) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	peer := (*cipher.Peer)(unsafe.Pointer(_peer))
	peer.ResetRetryTimes()
	return
}

// export SKY_pex_Peer_CanTry
func SKY_pex_Peer_CanTry(_peer *C.Peer, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	peer := (*cipher.Peer)(unsafe.Pointer(_peer))
	__arg0 := peer.CanTry()
	*_arg0 = __arg0
	return
}

// export SKY_pex_Peer_String
func SKY_pex_Peer_String(_peer *C.Peer, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	peer := (*cipher.Peer)(unsafe.Pointer(_peer))
	__arg0 := peer.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_pex_NewConfig
func SKY_pex_NewConfig(_arg0 *C.Config) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := pex.NewConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofConfig))
	return
}

// export SKY_pex_New
func SKY_pex_New(_cfg *C.Config, _defaultConns *C.GoSlice_, _arg2 *C.Pex) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	cfg := *(*cipher.Config)(unsafe.Pointer(_cfg))
	defaultConns := *(*[]string)(unsafe.Pointer(_defaultConns))
	__arg2, ____return_err := pex.New(cfg, defaultConns)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofPex))
	}
	return
}

// export SKY_pex_Pex_Run
func SKY_pex_Pex_Run(_px *C.Pex) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	____return_err := px.Run()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_pex_Pex_Shutdown
func SKY_pex_Pex_Shutdown(_px *C.Pex) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	px.Shutdown()
	return
}

// export SKY_pex_Pex_AddPeer
func SKY_pex_Pex_AddPeer(_px *C.Pex, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addr := _addr
	____return_err := px.AddPeer(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_pex_Pex_AddPeers
func SKY_pex_Pex_AddPeers(_px *C.Pex, _addrs *C.GoSlice_, _arg1 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1 := px.AddPeers(addrs)
	*_arg1 = __arg1
	return
}

// export SKY_pex_Pex_SetPrivate
func SKY_pex_Pex_SetPrivate(_px *C.Pex, _addr string, _private bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addr := _addr
	private := _private
	____return_err := px.SetPrivate(addr, private)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_pex_Pex_SetTrusted
func SKY_pex_Pex_SetTrusted(_px *C.Pex, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addr := _addr
	____return_err := px.SetTrusted(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_pex_Pex_SetHasIncomingPort
func SKY_pex_Pex_SetHasIncomingPort(_px *C.Pex, _addr string, _hasPublicPort bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addr := _addr
	hasPublicPort := _hasPublicPort
	____return_err := px.SetHasIncomingPort(addr, hasPublicPort)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_pex_Pex_RemovePeer
func SKY_pex_Pex_RemovePeer(_px *C.Pex, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addr := _addr
	px.RemovePeer(addr)
	return
}

// export SKY_pex_Pex_GetPeerByAddr
func SKY_pex_Pex_GetPeerByAddr(_px *C.Pex, _addr string, _arg1 *C.Peer, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addr := _addr
	__arg1, __arg2 := px.GetPeerByAddr(addr)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofPeer))
	*_arg2 = __arg2
	return
}

// export SKY_pex_Pex_Trusted
func SKY_pex_Pex_Trusted(_px *C.Pex, _arg0 *C.Peers) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	__arg0 := px.Trusted()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofPeers))
	return
}

// export SKY_pex_Pex_Private
func SKY_pex_Pex_Private(_px *C.Pex, _arg0 *C.Peers) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	__arg0 := px.Private()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofPeers))
	return
}

// export SKY_pex_Pex_TrustedPublic
func SKY_pex_Pex_TrustedPublic(_px *C.Pex, _arg0 *C.Peers) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	__arg0 := px.TrustedPublic()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofPeers))
	return
}

// export SKY_pex_Pex_RandomPublic
func SKY_pex_Pex_RandomPublic(_px *C.Pex, _n int, _arg1 *C.Peers) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	n := _n
	__arg1 := px.RandomPublic(n)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofPeers))
	return
}

// export SKY_pex_Pex_RandomExchangeable
func SKY_pex_Pex_RandomExchangeable(_px *C.Pex, _n int, _arg1 *C.Peers) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	n := _n
	__arg1 := px.RandomExchangeable(n)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofPeers))
	return
}

// export SKY_pex_Pex_IncreaseRetryTimes
func SKY_pex_Pex_IncreaseRetryTimes(_px *C.Pex, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addr := _addr
	px.IncreaseRetryTimes(addr)
	return
}

// export SKY_pex_Pex_ResetRetryTimes
func SKY_pex_Pex_ResetRetryTimes(_px *C.Pex, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	addr := _addr
	px.ResetRetryTimes(addr)
	return
}

// export SKY_pex_Pex_ResetAllRetryTimes
func SKY_pex_Pex_ResetAllRetryTimes(_px *C.Pex) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	px.ResetAllRetryTimes()
	return
}

// export SKY_pex_Pex_IsFull
func SKY_pex_Pex_IsFull(_px *C.Pex, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	px := (*cipher.Pex)(unsafe.Pointer(_px))
	__arg0 := px.IsFull()
	*_arg0 = __arg0
	return
}
