package main

import (
	consensus "github.com/skycoin/skycoin/src/consensus"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_consensus_ConsensusParticipant_GetConnectionManager
func SKY_consensus_ConsensusParticipant_GetConnectionManager(_self *C.ConsensusParticipant, _arg0 *C.ConnectionManagerInterface) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*ConsensusParticipant)(unsafe.Pointer(_self))
	__arg0 := self.GetConnectionManager()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofConnectionManagerInterface))
	return
}

// export SKY_consensus_ConsensusParticipant_GetNextBlockSeqNo
func SKY_consensus_ConsensusParticipant_GetNextBlockSeqNo(_self *C.ConsensusParticipant, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*ConsensusParticipant)(unsafe.Pointer(_self))
	__arg0 := self.GetNextBlockSeqNo()
	*_arg0 = __arg0
	return
}

// export SKY_consensus_ConsensusParticipant_SetPubkeySeckey
func SKY_consensus_ConsensusParticipant_SetPubkeySeckey(_self *C.ConsensusParticipant, _pubkey *C.PubKey, _seckey *C.SecKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*ConsensusParticipant)(unsafe.Pointer(_self))
	self.SetPubkeySeckey(pubkey, seckey)
	return
}

// export SKY_consensus_ConsensusParticipant_Print
func SKY_consensus_ConsensusParticipant_Print(_self *C.ConsensusParticipant) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*ConsensusParticipant)(unsafe.Pointer(_self))
	self.Print()
	return
}

// export SKY_consensus_NewConsensusParticipantPtr
func SKY_consensus_NewConsensusParticipantPtr(_pMan *C.ConnectionManagerInterface, _arg1 *C.ConsensusParticipant) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pMan := *(*ConnectionManagerInterface)(unsafe.Pointer(_pMan))
	__arg1 := consensus.NewConsensusParticipantPtr(pMan)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofConsensusParticipant))
	return
}

// export SKY_consensus_ConsensusParticipant_SignatureOf
func SKY_consensus_ConsensusParticipant_SignatureOf(_self *C.ConsensusParticipant, _hash *C.SHA256, _arg1 *C.Sig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*ConsensusParticipant)(unsafe.Pointer(_self))
	__arg1 := self.SignatureOf(hash)
	return
}

// export SKY_consensus_ConsensusParticipant_Get_block_stat_queue_Len
func SKY_consensus_ConsensusParticipant_Get_block_stat_queue_Len(_self *C.ConsensusParticipant, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*ConsensusParticipant)(unsafe.Pointer(_self))
	__arg0 := self.Get_block_stat_queue_Len()
	*_arg0 = __arg0
	return
}

// export SKY_consensus_ConsensusParticipant_Get_block_stat_queue_element_at
func SKY_consensus_ConsensusParticipant_Get_block_stat_queue_element_at(_self *C.ConsensusParticipant, _j int, _arg1 *C.BlockStat) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*ConsensusParticipant)(unsafe.Pointer(_self))
	j := _j
	__arg1 := self.Get_block_stat_queue_element_at(j)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofBlockStat))
	return
}

// export SKY_consensus_ConsensusParticipant_OnBlockHeaderArrived
func SKY_consensus_ConsensusParticipant_OnBlockHeaderArrived(_self *C.ConsensusParticipant, _blockPtr *C.BlockBase) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*ConsensusParticipant)(unsafe.Pointer(_self))
	blockPtr := (*BlockBase)(unsafe.Pointer(_blockPtr))
	self.OnBlockHeaderArrived(blockPtr)
	return
}
