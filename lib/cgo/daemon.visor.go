package main

import (
	coin "github.com/skycoin/skycoin/src/coin"
	daemon "github.com/skycoin/skycoin/src/daemon"
	gnet "github.com/skycoin/skycoin/src/daemon/gnet"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_daemon_NewGetBlocksMessage
func SKY_daemon_NewGetBlocksMessage(_lastBlock uint64, _requestedBlocks uint64, _arg2 *C.daemon__GetBlocksMessage) (____error_code uint32) {
	____error_code = 0
	lastBlock := _lastBlock
	requestedBlocks := _requestedBlocks
	__arg2 := daemon.NewGetBlocksMessage(lastBlock, requestedBlocks)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofGetBlocksMessage))
	return
}

// export SKY_daemon_GetBlocksMessage_Handle
func SKY_daemon_GetBlocksMessage_Handle(_gbm *C.daemon__GetBlocksMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	gbm := (*daemon.GetBlocksMessage)(unsafe.Pointer(_gbm))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := gbm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_NewGiveBlocksMessage
func SKY_daemon_NewGiveBlocksMessage(_blocks *C.GoSlice_, _arg1 *C.daemon__GiveBlocksMessage) (____error_code uint32) {
	____error_code = 0
	__arg1 := daemon.NewGiveBlocksMessage(blocks)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofGiveBlocksMessage))
	return
}

// export SKY_daemon_GiveBlocksMessage_Handle
func SKY_daemon_GiveBlocksMessage_Handle(_gbm *C.daemon__GiveBlocksMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	gbm := (*daemon.GiveBlocksMessage)(unsafe.Pointer(_gbm))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := gbm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_NewAnnounceBlocksMessage
func SKY_daemon_NewAnnounceBlocksMessage(_seq uint64, _arg1 *C.daemon__AnnounceBlocksMessage) (____error_code uint32) {
	____error_code = 0
	seq := _seq
	__arg1 := daemon.NewAnnounceBlocksMessage(seq)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofAnnounceBlocksMessage))
	return
}

// export SKY_daemon_AnnounceBlocksMessage_Handle
func SKY_daemon_AnnounceBlocksMessage_Handle(_abm *C.daemon__AnnounceBlocksMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	abm := (*daemon.AnnounceBlocksMessage)(unsafe.Pointer(_abm))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := abm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_NewAnnounceTxnsMessage
func SKY_daemon_NewAnnounceTxnsMessage(_txns *C.GoSlice_, _arg1 *C.daemon__AnnounceTxnsMessage) (____error_code uint32) {
	____error_code = 0
	__arg1 := daemon.NewAnnounceTxnsMessage(txns)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofAnnounceTxnsMessage))
	return
}

// export SKY_daemon_AnnounceTxnsMessage_GetTxns
func SKY_daemon_AnnounceTxnsMessage_GetTxns(_atm *C.daemon__AnnounceTxnsMessage, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	atm := (*daemon.AnnounceTxnsMessage)(unsafe.Pointer(_atm))
	__arg0 := atm.GetTxns()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_AnnounceTxnsMessage_Handle
func SKY_daemon_AnnounceTxnsMessage_Handle(_atm *C.daemon__AnnounceTxnsMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	atm := (*daemon.AnnounceTxnsMessage)(unsafe.Pointer(_atm))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := atm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_NewGetTxnsMessage
func SKY_daemon_NewGetTxnsMessage(_txns *C.GoSlice_, _arg1 *C.daemon__GetTxnsMessage) (____error_code uint32) {
	____error_code = 0
	__arg1 := daemon.NewGetTxnsMessage(txns)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofGetTxnsMessage))
	return
}

// export SKY_daemon_GetTxnsMessage_Handle
func SKY_daemon_GetTxnsMessage_Handle(_gtm *C.daemon__GetTxnsMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	gtm := (*daemon.GetTxnsMessage)(unsafe.Pointer(_gtm))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := gtm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_NewGiveTxnsMessage
func SKY_daemon_NewGiveTxnsMessage(_txns *C.coin__Transactions, _arg1 *C.daemon__GiveTxnsMessage) (____error_code uint32) {
	____error_code = 0
	txns := *(*coin.Transactions)(unsafe.Pointer(_txns))
	__arg1 := daemon.NewGiveTxnsMessage(txns)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofGiveTxnsMessage))
	return
}

// export SKY_daemon_GiveTxnsMessage_GetTxns
func SKY_daemon_GiveTxnsMessage_GetTxns(_gtm *C.daemon__GiveTxnsMessage, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	gtm := (*daemon.GiveTxnsMessage)(unsafe.Pointer(_gtm))
	__arg0 := gtm.GetTxns()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_GiveTxnsMessage_Handle
func SKY_daemon_GiveTxnsMessage_Handle(_gtm *C.daemon__GiveTxnsMessage, _mc *C.gnet__MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	gtm := (*daemon.GiveTxnsMessage)(unsafe.Pointer(_gtm))
	mc := (*gnet.MessageContext)(unsafe.Pointer(_mc))
	____return_err := gtm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
