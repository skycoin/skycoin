package main

import (
	consensus "github.com/skycoin/skycoin/src/consensus"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_main_BlockBaseWrapper_String
func SKY_main_BlockBaseWrapper_String(_self *C.BlockBaseWrapper, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockBaseWrapper)(unsafe.Pointer(_self))
	__arg0 := self.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_main_BlockBaseWrapper_Handle
func SKY_main_BlockBaseWrapper_Handle(_self *C.BlockBaseWrapper, _context *C.MessageContext, _closure interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockBaseWrapper)(unsafe.Pointer(_self))
	____return_err := self.Handle(context, closure)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_main_PoolOwner_Shutdown
func SKY_main_PoolOwner_Shutdown(_self *C.PoolOwner) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	self.Shutdown()
	return
}

// export SKY_main_PoolOwner_Init
func SKY_main_PoolOwner_Init(_self *C.PoolOwner, _pCMan *C.MinimalConnectionManager, _listen_port uint16, _num_id int, _nickname string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	pCMan := (*MinimalConnectionManager)(unsafe.Pointer(_pCMan))
	listen_port := _listen_port
	num_id := _num_id
	nickname := _nickname
	self.Init(pCMan, listen_port, num_id, nickname)
	return
}

// export SKY_main_PoolOwner_Run
func SKY_main_PoolOwner_Run(_self *C.PoolOwner) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	self.Run()
	return
}

// export SKY_main_PoolOwner_DataCallback
func SKY_main_PoolOwner_DataCallback(_self *C.PoolOwner, _context *C.MessageContext, _xxx *C.BlockBaseWrapper) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	xxx := (*BlockBaseWrapper)(unsafe.Pointer(_xxx))
	self.DataCallback(context, xxx)
	return
}

// export SKY_main_PoolOwner_OnConnect
func SKY_main_PoolOwner_OnConnect(_self *C.PoolOwner, _c *C.Connection, _is_solicited bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	is_solicited := _is_solicited
	self.OnConnect(c, is_solicited)
	return
}

// export SKY_main_PoolOwner_OnDisconnect
func SKY_main_PoolOwner_OnDisconnect(_self *C.PoolOwner, _c *C.Connection, _reason *C.DisconnectReason) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	self.OnDisconnect(c, reason)
	return
}

// export SKY_main_PoolOwner_RequestConnectionToKeys
func SKY_main_PoolOwner_RequestConnectionToKeys(_self *C.PoolOwner) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	self.RequestConnectionToKeys()
	return
}

// export SKY_main_PoolOwner_RegisterKey
func SKY_main_PoolOwner_RegisterKey(_self *C.PoolOwner, _key string, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	key := _key
	__arg1 := self.RegisterKey(key)
	*_arg1 = __arg1
	return
}

// export SKY_main_PoolOwner_GetListenPort
func SKY_main_PoolOwner_GetListenPort(_self *C.PoolOwner, _arg0 *uint16) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	__arg0 := self.GetListenPort()
	*_arg0 = __arg0
	return
}

// export SKY_main_PoolOwner_BlockingConnectTo
func SKY_main_PoolOwner_BlockingConnectTo(_self *C.PoolOwner, _IPAddress string, _port uint16, _arg2 *C.Connection) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	IPAddress := _IPAddress
	port := _port
	__arg2, ____return_err := self.BlockingConnectTo(IPAddress, port)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_main_PoolOwner_BlockingConnectToUrl
func SKY_main_PoolOwner_BlockingConnectToUrl(_self *C.PoolOwner, _url string, _arg1 *C.Connection) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	url := _url
	__arg1, ____return_err := self.BlockingConnectToUrl(url)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_main_PoolOwner_BroadcastMessage
func SKY_main_PoolOwner_BroadcastMessage(_self *C.PoolOwner, _msg *C.Message) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	____return_err := self.BroadcastMessage(msg)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_main_PoolOwner_Print
func SKY_main_PoolOwner_Print(_self *C.PoolOwner) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*PoolOwner)(unsafe.Pointer(_self))
	self.Print()
	return
}

// export SKY_main_Simulate_compare_node_StateQueue
func SKY_main_Simulate_compare_node_StateQueue(_X *C.GoSlice_, _global_seqno2h map[uint64]*C.SHA256, _global_seqno2h_alt map[uint64]*C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	consensus.Simulate_compare_node_StateQueue(X, global_seqno2h, global_seqno2h_alt)
	return
}

// export SKY_main_MinimalConnectionManager_Init
func SKY_main_MinimalConnectionManager_Init(_self *C.MinimalConnectionManager, _listen_port uint16, _num_id int, _nickname string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	listen_port := _listen_port
	num_id := _num_id
	nickname := _nickname
	self.Init(listen_port, num_id, nickname)
	return
}

// export SKY_main_MinimalConnectionManager_Run
func SKY_main_MinimalConnectionManager_Run(_self *C.MinimalConnectionManager) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	self.Run()
	return
}

// export SKY_main_MinimalConnectionManager_ShutdownPublishing
func SKY_main_MinimalConnectionManager_ShutdownPublishing(_self *C.MinimalConnectionManager) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	self.ShutdownPublishing()
	return
}

// export SKY_main_MinimalConnectionManager_ShutdownSubscribing
func SKY_main_MinimalConnectionManager_ShutdownSubscribing(_self *C.MinimalConnectionManager) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	self.ShutdownSubscribing()
	return
}

// export SKY_main_MinimalConnectionManager_GetListenPort
func SKY_main_MinimalConnectionManager_GetListenPort(_self *C.MinimalConnectionManager, _arg0 *uint16) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	__arg0 := self.GetListenPort()
	*_arg0 = __arg0
	return
}

// export SKY_main_MinimalConnectionManager_GetNode
func SKY_main_MinimalConnectionManager_GetNode(_self *C.MinimalConnectionManager, _arg0 *C.ConsensusParticipant) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	__arg0 := self.GetNode()
	return
}

// export SKY_main_MinimalConnectionManager_RegisterPublisher
func SKY_main_MinimalConnectionManager_RegisterPublisher(_self *C.MinimalConnectionManager, _key string, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	key := _key
	__arg1 := self.RegisterPublisher(key)
	*_arg1 = __arg1
	return
}

// export SKY_main_MinimalConnectionManager_SendBlockToAllMySubscriber
func SKY_main_MinimalConnectionManager_SendBlockToAllMySubscriber(_self *C.MinimalConnectionManager, _xxx *C.BlockBase) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	self.SendBlockToAllMySubscriber(xxx)
	return
}

// export SKY_main_MinimalConnectionManager_RequestConnectionToAllMyPublisher
func SKY_main_MinimalConnectionManager_RequestConnectionToAllMyPublisher(_self *C.MinimalConnectionManager) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	self.RequestConnectionToAllMyPublisher()
	return
}

// export SKY_main_MinimalConnectionManager_OnSubscriberConnectionRequest
func SKY_main_MinimalConnectionManager_OnSubscriberConnectionRequest(_self *C.MinimalConnectionManager, _key string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	key := _key
	self.OnSubscriberConnectionRequest(key)
	return
}

// export SKY_main_MinimalConnectionManager_Print
func SKY_main_MinimalConnectionManager_Print(_self *C.MinimalConnectionManager) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	self.Print()
	return
}
