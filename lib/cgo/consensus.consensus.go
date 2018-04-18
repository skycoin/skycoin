package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	consensus "github.com/skycoin/skycoin/src/consensus"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_consensus_BlockBase_Init
func SKY_consensus_BlockBase_Init(_self *C.consensus__BlockBase, _sig *C.cipher__Sig, _hash *C.cipher__SHA256, _seqno uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.BlockBase)(unsafe.Pointer(_self))
	sig := *(*cipher.Sig)(unsafe.Pointer(_sig))
	hash := *(*cipher.SHA256)(unsafe.Pointer(_hash))
	seqno := _seqno
	self.Init(sig, hash, seqno)
	return
}

//export SKY_consensus_BlockBase_Print
func SKY_consensus_BlockBase_Print(_self *C.consensus__BlockBase) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.BlockBase)(unsafe.Pointer(_self))
	self.Print()
	return
}

//export SKY_consensus_BlockBase_String
func SKY_consensus_BlockBase_String(_self *C.consensus__BlockBase, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.BlockBase)(unsafe.Pointer(_self))
	__arg0 := self.String()
	copyString(__arg0, _arg0)
	return
}

//export SKY_consensus_BlockchainTail_Init
func SKY_consensus_BlockchainTail_Init(_self *C.consensus__BlockchainTail) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.BlockchainTail)(unsafe.Pointer(_self))
	self.Init()
	return
}

//export SKY_consensus_BlockchainTail_GetNextSeqNo
func SKY_consensus_BlockchainTail_GetNextSeqNo(_self *C.consensus__BlockchainTail, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.BlockchainTail)(unsafe.Pointer(_self))
	__arg0 := self.GetNextSeqNo()
	*_arg0 = __arg0
	return
}

//export SKY_consensus_BlockchainTail_Print
func SKY_consensus_BlockchainTail_Print(_self *C.consensus__BlockchainTail) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.BlockchainTail)(unsafe.Pointer(_self))
	self.Print()
	return
}

//export SKY_consensus_HashCandidate_Init
func SKY_consensus_HashCandidate_Init(_self *C.consensus__HashCandidate) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.HashCandidate)(unsafe.Pointer(_self))
	self.Init()
	return
}

//export SKY_consensus_HashCandidate_ObserveSigAndPubkey
func SKY_consensus_HashCandidate_ObserveSigAndPubkey(_self *C.consensus__HashCandidate, _sig *C.cipher__Sig, _pubkey *C.cipher__PubKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.HashCandidate)(unsafe.Pointer(_self))
	sig := *(*cipher.Sig)(unsafe.Pointer(_sig))
	pubkey := *(*cipher.PubKey)(unsafe.Pointer(_pubkey))
	self.ObserveSigAndPubkey(sig, pubkey)
	return
}

//export SKY_consensus_HashCandidate_Clear
func SKY_consensus_HashCandidate_Clear(_self *C.consensus__HashCandidate) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.HashCandidate)(unsafe.Pointer(_self))
	self.Clear()
	return
}
