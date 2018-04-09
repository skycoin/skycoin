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
func SKY_main_MinimalConnectionManager_RegisterPublisher(_self *C.MinimalConnectionManager, _key *C.MinimalConnectionManager, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	key := (*MinimalConnectionManager)(unsafe.Pointer(_key))
	__arg1 := self.RegisterPublisher(key)
	*_arg1 = __arg1
	return
}

// export SKY_main_MinimalConnectionManager_SendBlockToAllMySubscriber
func SKY_main_MinimalConnectionManager_SendBlockToAllMySubscriber(_self *C.MinimalConnectionManager, _blockPtr *C.BlockBase) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	self.SendBlockToAllMySubscriber(blockPtr)
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
func SKY_main_MinimalConnectionManager_OnSubscriberConnectionRequest(_self *C.MinimalConnectionManager, _other *C.MinimalConnectionManager) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*MinimalConnectionManager)(unsafe.Pointer(_self))
	other := (*MinimalConnectionManager)(unsafe.Pointer(_other))
	self.OnSubscriberConnectionRequest(other)
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

// export SKY_main_Simulate_compare_node_StateQueue
func SKY_main_Simulate_compare_node_StateQueue(_X *C.GoSlice_, _global_seqno2h map[uint64]*C.SHA256, _global_seqno2h_alt map[uint64]*C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	consensus.Simulate_compare_node_StateQueue(X, global_seqno2h, global_seqno2h_alt)
	return
}
