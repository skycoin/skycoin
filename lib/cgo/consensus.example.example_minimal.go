package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_main_MinimalConnectionManager_GetNode
func SKY_main_MinimalConnectionManager_GetNode(_self *C.MinimalConnectionManager, _arg0 *C.ConsensusParticipant) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	self := (*cipher.MinimalConnectionManager)(unsafe.Pointer(_self))
	__arg0 := self.GetNode()
	return
}

// export SKY_main_MinimalConnectionManager_RegisterPublisher
func SKY_main_MinimalConnectionManager_RegisterPublisher(_self *C.MinimalConnectionManager, _key *C.MinimalConnectionManager, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	self := (*cipher.MinimalConnectionManager)(unsafe.Pointer(_self))
	key := (*cipher.MinimalConnectionManager)(unsafe.Pointer(_key))
	__arg1 := self.RegisterPublisher(key)
	*_arg1 = __arg1
	return
}

// export SKY_main_MinimalConnectionManager_SendBlockToAllMySubscriber
func SKY_main_MinimalConnectionManager_SendBlockToAllMySubscriber(_self *C.MinimalConnectionManager, _blockPtr *C.BlockBase) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	self := (*cipher.MinimalConnectionManager)(unsafe.Pointer(_self))
	self.SendBlockToAllMySubscriber(blockPtr)
	return
}

// export SKY_main_MinimalConnectionManager_RequestConnectionToAllMyPublisher
func SKY_main_MinimalConnectionManager_RequestConnectionToAllMyPublisher(_self *C.MinimalConnectionManager) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	self := (*cipher.MinimalConnectionManager)(unsafe.Pointer(_self))
	self.RequestConnectionToAllMyPublisher()
	return
}

// export SKY_main_MinimalConnectionManager_OnSubscriberConnectionRequest
func SKY_main_MinimalConnectionManager_OnSubscriberConnectionRequest(_self *C.MinimalConnectionManager, _other *C.MinimalConnectionManager) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	self := (*cipher.MinimalConnectionManager)(unsafe.Pointer(_self))
	other := (*cipher.MinimalConnectionManager)(unsafe.Pointer(_other))
	self.OnSubscriberConnectionRequest(other)
	return
}

// export SKY_main_MinimalConnectionManager_Print
func SKY_main_MinimalConnectionManager_Print(_self *C.MinimalConnectionManager) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	self := (*cipher.MinimalConnectionManager)(unsafe.Pointer(_self))
	self.Print()
	return
}
