package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
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
func SKY_consensus_ConsensusParticipant_GetConnectionManager(_self *C.consensus__ConsensusParticipant, _arg0 *C.consensus__ConnectionManagerInterface) (____error_code uint32) {
	____error_code = 0
	self := (*consensus.ConsensusParticipant)(unsafe.Pointer(_self))
	__arg0 := self.GetConnectionManager()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofConnectionManagerInterface))
	return
}

// export SKY_consensus_ConsensusParticipant_GetNextBlockSeqNo
func SKY_consensus_ConsensusParticipant_GetNextBlockSeqNo(_self *C.consensus__ConsensusParticipant, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	self := (*consensus.ConsensusParticipant)(unsafe.Pointer(_self))
	__arg0 := self.GetNextBlockSeqNo()
	*_arg0 = __arg0
	return
}

// export SKY_consensus_ConsensusParticipant_SetPubkeySeckey
func SKY_consensus_ConsensusParticipant_SetPubkeySeckey(_self *C.consensus__ConsensusParticipant, _pubkey *C.cipher__PubKey, _seckey *C.cipher__SecKey) (____error_code uint32) {
	____error_code = 0
	self := (*consensus.ConsensusParticipant)(unsafe.Pointer(_self))
	pubkey := *(*cipher.PubKey)(unsafe.Pointer(_pubkey))
	seckey := *(*cipher.SecKey)(unsafe.Pointer(_seckey))
	self.SetPubkeySeckey(pubkey, seckey)
	return
}

// export SKY_consensus_ConsensusParticipant_Print
func SKY_consensus_ConsensusParticipant_Print(_self *C.consensus__ConsensusParticipant) (____error_code uint32) {
	____error_code = 0
	self := (*consensus.ConsensusParticipant)(unsafe.Pointer(_self))
	self.Print()
	return
}

// export SKY_consensus_NewConsensusParticipantPtr
func SKY_consensus_NewConsensusParticipantPtr(_pMan *C.consensus__ConnectionManagerInterface, _arg1 *C.consensus__ConsensusParticipant) (____error_code uint32) {
	____error_code = 0
	pMan := *(*consensus.ConnectionManagerInterface)(unsafe.Pointer(_pMan))
	__arg1 := consensus.NewConsensusParticipantPtr(pMan)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofConsensusParticipant))
	return
}

// export SKY_consensus_ConsensusParticipant_SignatureOf
func SKY_consensus_ConsensusParticipant_SignatureOf(_self *C.consensus__ConsensusParticipant, _hash *C.cipher__SHA256, _arg1 *C.cipher__Sig) (____error_code uint32) {
	____error_code = 0
	self := (*consensus.ConsensusParticipant)(unsafe.Pointer(_self))
	hash := *(*cipher.SHA256)(unsafe.Pointer(_hash))
	__arg1 := self.SignatureOf(hash)
	return
}

// export SKY_consensus_ConsensusParticipant_Get_block_stat_queue_Len
func SKY_consensus_ConsensusParticipant_Get_block_stat_queue_Len(_self *C.consensus__ConsensusParticipant, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	self := (*consensus.ConsensusParticipant)(unsafe.Pointer(_self))
	__arg0 := self.Get_block_stat_queue_Len()
	*_arg0 = __arg0
	return
}

// export SKY_consensus_ConsensusParticipant_Get_block_stat_queue_element_at
func SKY_consensus_ConsensusParticipant_Get_block_stat_queue_element_at(_self *C.consensus__ConsensusParticipant, _j int, _arg1 *C.consensus__BlockStat) (____error_code uint32) {
	____error_code = 0
	self := (*consensus.ConsensusParticipant)(unsafe.Pointer(_self))
	j := _j
	__arg1 := self.Get_block_stat_queue_element_at(j)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofBlockStat))
	return
}

// export SKY_consensus_ConsensusParticipant_OnBlockHeaderArrived
func SKY_consensus_ConsensusParticipant_OnBlockHeaderArrived(_self *C.consensus__ConsensusParticipant, _blockPtr *C.consensus__BlockBase) (____error_code uint32) {
	____error_code = 0
	self := (*consensus.ConsensusParticipant)(unsafe.Pointer(_self))
	blockPtr := (*consensus.BlockBase)(unsafe.Pointer(_blockPtr))
	self.OnBlockHeaderArrived(blockPtr)
	return
}
