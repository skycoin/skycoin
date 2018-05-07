package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	daemon "github.com/skycoin/skycoin/src/daemon"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_daemon_NewGetBlocksMessage
func SKY_daemon_NewGetBlocksMessage(_lastBlock uint64, _requestedBlocks uint64, _arg2 *C.daemon__GetBlocksMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	lastBlock := _lastBlock
	requestedBlocks := _requestedBlocks
	__arg2 := daemon.NewGetBlocksMessage(lastBlock, requestedBlocks)
	*_arg2 = *(*C.daemon__GetBlocksMessage)(unsafe.Pointer(__arg2))
	return
}

//export SKY_daemon_NewGiveBlocksMessage
func SKY_daemon_NewGiveBlocksMessage(_blocks []C.coin__SignedBlock, _arg1 *C.daemon__GiveBlocksMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	blocks := *(*[]coin.SignedBlock)(unsafe.Pointer(&_blocks))
	__arg1 := daemon.NewGiveBlocksMessage(blocks)
	*_arg1 = *(*C.daemon__GiveBlocksMessage)(unsafe.Pointer(__arg1))
	return
}

//export SKY_daemon_NewAnnounceBlocksMessage
func SKY_daemon_NewAnnounceBlocksMessage(_seq uint64, _arg1 *C.daemon__AnnounceBlocksMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seq := _seq
	__arg1 := daemon.NewAnnounceBlocksMessage(seq)
	*_arg1 = *(*C.daemon__AnnounceBlocksMessage)(unsafe.Pointer(__arg1))
	return
}

//export SKY_daemon_NewAnnounceTxnsMessage
func SKY_daemon_NewAnnounceTxnsMessage(_txns []C.cipher__SHA256, _arg1 *C.daemon__AnnounceTxnsMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*[]cipher.SHA256)(unsafe.Pointer(&_txns))
	__arg1 := daemon.NewAnnounceTxnsMessage(txns)
	*_arg1 = *(*C.daemon__AnnounceTxnsMessage)(unsafe.Pointer(__arg1))
	return
}

//export SKY_daemon_AnnounceTxnsMessage_GetTxns
func SKY_daemon_AnnounceTxnsMessage_GetTxns(_atm *C.daemon__AnnounceTxnsMessage, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	atm := (*daemon.AnnounceTxnsMessage)(unsafe.Pointer(_atm))
	__arg0 := atm.GetTxns()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_daemon_NewGetTxnsMessage
func SKY_daemon_NewGetTxnsMessage(_txns []C.cipher__SHA256, _arg1 *C.daemon__GetTxnsMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*[]cipher.SHA256)(unsafe.Pointer(&_txns))
	__arg1 := daemon.NewGetTxnsMessage(txns)
	*_arg1 = *(*C.daemon__GetTxnsMessage)(unsafe.Pointer(__arg1))
	return
}

//export SKY_daemon_NewGiveTxnsMessage
func SKY_daemon_NewGiveTxnsMessage(_txns *C.coin__Transactions, _arg1 *C.daemon__GiveTxnsMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*coin.Transactions)(unsafe.Pointer(_txns))
	__arg1 := daemon.NewGiveTxnsMessage(txns)
	*_arg1 = *(*C.daemon__GiveTxnsMessage)(unsafe.Pointer(__arg1))
	return
}

//export SKY_daemon_GiveTxnsMessage_GetTxns
func SKY_daemon_GiveTxnsMessage_GetTxns(_gtm *C.daemon__GiveTxnsMessage, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gtm := (*daemon.GiveTxnsMessage)(unsafe.Pointer(_gtm))
	__arg0 := gtm.GetTxns()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
