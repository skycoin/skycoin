package main

import "unsafe"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_consensus_BlockBase_Init
func SKY_consensus_BlockBase_Init(_self *C.BlockBase, _sig *C.Sig, _hash *C.SHA256, _seqno uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockBase)(unsafe.Pointer(_self))
	seqno := _seqno
	self.Init(sig, hash, seqno)
	return
}

// export SKY_consensus_BlockBase_Print
func SKY_consensus_BlockBase_Print(_self *C.BlockBase) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockBase)(unsafe.Pointer(_self))
	self.Print()
	return
}

// export SKY_consensus_BlockBase_String
func SKY_consensus_BlockBase_String(_self *C.BlockBase, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockBase)(unsafe.Pointer(_self))
	__arg0 := self.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_consensus_BlockchainTail_Init
func SKY_consensus_BlockchainTail_Init(_self *C.BlockchainTail) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockchainTail)(unsafe.Pointer(_self))
	self.Init()
	return
}

// export SKY_consensus_BlockchainTail_GetNextSeqNo
func SKY_consensus_BlockchainTail_GetNextSeqNo(_self *C.BlockchainTail, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockchainTail)(unsafe.Pointer(_self))
	__arg0 := self.GetNextSeqNo()
	*_arg0 = __arg0
	return
}

// export SKY_consensus_BlockchainTail_Print
func SKY_consensus_BlockchainTail_Print(_self *C.BlockchainTail) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockchainTail)(unsafe.Pointer(_self))
	self.Print()
	return
}

// export SKY_consensus_HashCandidate_Init
func SKY_consensus_HashCandidate_Init(_self *C.HashCandidate) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*HashCandidate)(unsafe.Pointer(_self))
	self.Init()
	return
}

// export SKY_consensus_HashCandidate_ObserveSigAndPubkey
func SKY_consensus_HashCandidate_ObserveSigAndPubkey(_self *C.HashCandidate, _sig *C.Sig, _pubkey *C.PubKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*HashCandidate)(unsafe.Pointer(_self))
	self.ObserveSigAndPubkey(sig, pubkey)
	return
}

// export SKY_consensus_HashCandidate_Clear
func SKY_consensus_HashCandidate_Clear(_self *C.HashCandidate) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*HashCandidate)(unsafe.Pointer(_self))
	self.Clear()
	return
}
