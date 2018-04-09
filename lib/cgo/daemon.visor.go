package main

import (
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

// export SKY_daemon_NewVisorConfig
func SKY_daemon_NewVisorConfig(_arg0 *C.VisorConfig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := daemon.NewVisorConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofVisorConfig))
	return
}

// export SKY_daemon_NewVisor
func SKY_daemon_NewVisor(_c *C.VisorConfig, _db *C.DB, _arg2 *C.Visor) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*VisorConfig)(unsafe.Pointer(_c))
	__arg2, ____return_err := daemon.NewVisor(c, db)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofVisor))
	}
	return
}

// export SKY_daemon_Visor_Run
func SKY_daemon_Visor_Run(_vs *C.Visor) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	____return_err := vs.Run()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_Shutdown
func SKY_daemon_Visor_Shutdown(_vs *C.Visor) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	vs.Shutdown()
	return
}

// export SKY_daemon_Visor_RefreshUnconfirmed
func SKY_daemon_Visor_RefreshUnconfirmed(_vs *C.Visor, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0, ____return_err := vs.RefreshUnconfirmed()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_daemon_Visor_RemoveInvalidUnconfirmed
func SKY_daemon_Visor_RemoveInvalidUnconfirmed(_vs *C.Visor, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0, ____return_err := vs.RemoveInvalidUnconfirmed()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_daemon_Visor_RequestBlocks
func SKY_daemon_Visor_RequestBlocks(_vs *C.Visor, _pool *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	____return_err := vs.RequestBlocks(pool)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_AnnounceBlocks
func SKY_daemon_Visor_AnnounceBlocks(_vs *C.Visor, _pool *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	____return_err := vs.AnnounceBlocks(pool)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_AnnounceAllTxns
func SKY_daemon_Visor_AnnounceAllTxns(_vs *C.Visor, _pool *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	____return_err := vs.AnnounceAllTxns(pool)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_AnnounceTxns
func SKY_daemon_Visor_AnnounceTxns(_vs *C.Visor, _pool *C.Pool, _txns *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	____return_err := vs.AnnounceTxns(pool, txns)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_RequestBlocksFromAddr
func SKY_daemon_Visor_RequestBlocksFromAddr(_vs *C.Visor, _pool *C.Pool, _addr string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	addr := _addr
	____return_err := vs.RequestBlocksFromAddr(pool, addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_SetTxnsAnnounced
func SKY_daemon_Visor_SetTxnsAnnounced(_vs *C.Visor, _txns *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	vs.SetTxnsAnnounced(txns)
	return
}

// export SKY_daemon_Visor_InjectBroadcastTransaction
func SKY_daemon_Visor_InjectBroadcastTransaction(_vs *C.Visor, _txn *C.Transaction, _pool *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	____return_err := vs.InjectBroadcastTransaction(txn, pool)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_InjectTransaction
func SKY_daemon_Visor_InjectTransaction(_vs *C.Visor, _tx *C.Transaction, _arg1 *bool, _arg2 *C.ErrTxnViolatesSoftConstraint) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1, __arg2, ____return_err := vs.InjectTransaction(tx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_daemon_Visor_ResendTransaction
func SKY_daemon_Visor_ResendTransaction(_vs *C.Visor, _h *C.SHA256, _pool *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	____return_err := vs.ResendTransaction(h, pool)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_ResendUnconfirmedTxns
func SKY_daemon_Visor_ResendUnconfirmedTxns(_vs *C.Visor, _pool *C.Pool, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	__arg1 := vs.ResendUnconfirmedTxns(pool)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_daemon_Visor_CreateAndPublishBlock
func SKY_daemon_Visor_CreateAndPublishBlock(_vs *C.Visor, _pool *C.Pool, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	pool := (*Pool)(unsafe.Pointer(_pool))
	__arg1, ____return_err := vs.CreateAndPublishBlock(pool)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_RemoveConnection
func SKY_daemon_Visor_RemoveConnection(_vs *C.Visor, _addr string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	addr := _addr
	vs.RemoveConnection(addr)
	return
}

// export SKY_daemon_Visor_RecordBlockchainHeight
func SKY_daemon_Visor_RecordBlockchainHeight(_vs *C.Visor, _addr string, _bkLen uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	addr := _addr
	bkLen := _bkLen
	vs.RecordBlockchainHeight(addr, bkLen)
	return
}

// export SKY_daemon_Visor_EstimateBlockchainHeight
func SKY_daemon_Visor_EstimateBlockchainHeight(_vs *C.Visor, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0 := vs.EstimateBlockchainHeight()
	*_arg0 = __arg0
	return
}

// export SKY_daemon_Visor_ScanAheadWalletAddresses
func SKY_daemon_Visor_ScanAheadWalletAddresses(_vs *C.Visor, _wltName string, _password *C.GoSlice_, _scanN uint64, _arg3 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	wltName := _wltName
	password := *(*[]byte)(unsafe.Pointer(_password))
	scanN := _scanN
	__arg3, ____return_err := vs.ScanAheadWalletAddresses(wltName, password, scanN)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_GetPeerBlockchainHeights
func SKY_daemon_Visor_GetPeerBlockchainHeights(_vs *C.Visor, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0 := vs.GetPeerBlockchainHeights()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_Visor_HeadBkSeq
func SKY_daemon_Visor_HeadBkSeq(_vs *C.Visor, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0 := vs.HeadBkSeq()
	*_arg0 = __arg0
	return
}

// export SKY_daemon_Visor_ExecuteSignedBlock
func SKY_daemon_Visor_ExecuteSignedBlock(_vs *C.Visor, _b *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	____return_err := vs.ExecuteSignedBlock(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_GetSignedBlock
func SKY_daemon_Visor_GetSignedBlock(_vs *C.Visor, _seq uint64, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	seq := _seq
	__arg1, ____return_err := vs.GetSignedBlock(seq)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Visor_GetSignedBlocksSince
func SKY_daemon_Visor_GetSignedBlocksSince(_vs *C.Visor, _seq uint64, _ct uint64, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	seq := _seq
	ct := _ct
	__arg2, ____return_err := vs.GetSignedBlocksSince(seq, ct)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

// export SKY_daemon_Visor_UnConfirmFilterKnown
func SKY_daemon_Visor_UnConfirmFilterKnown(_vs *C.Visor, _txns *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1 := vs.UnConfirmFilterKnown(txns)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_daemon_Visor_UnConfirmKnow
func SKY_daemon_Visor_UnConfirmKnow(_vs *C.Visor, _hashes *C.GoSlice_, _arg1 *C.Transactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1 := vs.UnConfirmKnow(hashes)
	return
}

// export SKY_daemon_NewGetBlocksMessage
func SKY_daemon_NewGetBlocksMessage(_lastBlock uint64, _requestedBlocks uint64, _arg2 *C.GetBlocksMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	lastBlock := _lastBlock
	requestedBlocks := _requestedBlocks
	__arg2 := daemon.NewGetBlocksMessage(lastBlock, requestedBlocks)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofGetBlocksMessage))
	return
}

// export SKY_daemon_GetBlocksMessage_Handle
func SKY_daemon_GetBlocksMessage_Handle(_gbm *C.GetBlocksMessage, _mc *C.MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gbm := (*GetBlocksMessage)(unsafe.Pointer(_gbm))
	____return_err := gbm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_GetBlocksMessage_Process
func SKY_daemon_GetBlocksMessage_Process(_gbm *C.GetBlocksMessage, _d *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gbm := (*GetBlocksMessage)(unsafe.Pointer(_gbm))
	d := (*Daemon)(unsafe.Pointer(_d))
	gbm.Process(d)
	return
}

// export SKY_daemon_NewGiveBlocksMessage
func SKY_daemon_NewGiveBlocksMessage(_blocks *C.GoSlice_, _arg1 *C.GiveBlocksMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := daemon.NewGiveBlocksMessage(blocks)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofGiveBlocksMessage))
	return
}

// export SKY_daemon_GiveBlocksMessage_Handle
func SKY_daemon_GiveBlocksMessage_Handle(_gbm *C.GiveBlocksMessage, _mc *C.MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gbm := (*GiveBlocksMessage)(unsafe.Pointer(_gbm))
	____return_err := gbm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_GiveBlocksMessage_Process
func SKY_daemon_GiveBlocksMessage_Process(_gbm *C.GiveBlocksMessage, _d *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gbm := (*GiveBlocksMessage)(unsafe.Pointer(_gbm))
	d := (*Daemon)(unsafe.Pointer(_d))
	gbm.Process(d)
	return
}

// export SKY_daemon_NewAnnounceBlocksMessage
func SKY_daemon_NewAnnounceBlocksMessage(_seq uint64, _arg1 *C.AnnounceBlocksMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seq := _seq
	__arg1 := daemon.NewAnnounceBlocksMessage(seq)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofAnnounceBlocksMessage))
	return
}

// export SKY_daemon_AnnounceBlocksMessage_Handle
func SKY_daemon_AnnounceBlocksMessage_Handle(_abm *C.AnnounceBlocksMessage, _mc *C.MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	abm := (*AnnounceBlocksMessage)(unsafe.Pointer(_abm))
	____return_err := abm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_AnnounceBlocksMessage_Process
func SKY_daemon_AnnounceBlocksMessage_Process(_abm *C.AnnounceBlocksMessage, _d *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	abm := (*AnnounceBlocksMessage)(unsafe.Pointer(_abm))
	d := (*Daemon)(unsafe.Pointer(_d))
	abm.Process(d)
	return
}

// export SKY_daemon_NewAnnounceTxnsMessage
func SKY_daemon_NewAnnounceTxnsMessage(_txns *C.GoSlice_, _arg1 *C.AnnounceTxnsMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := daemon.NewAnnounceTxnsMessage(txns)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofAnnounceTxnsMessage))
	return
}

// export SKY_daemon_AnnounceTxnsMessage_GetTxns
func SKY_daemon_AnnounceTxnsMessage_GetTxns(_atm *C.AnnounceTxnsMessage, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	atm := (*AnnounceTxnsMessage)(unsafe.Pointer(_atm))
	__arg0 := atm.GetTxns()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_AnnounceTxnsMessage_Handle
func SKY_daemon_AnnounceTxnsMessage_Handle(_atm *C.AnnounceTxnsMessage, _mc *C.MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	atm := (*AnnounceTxnsMessage)(unsafe.Pointer(_atm))
	____return_err := atm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_AnnounceTxnsMessage_Process
func SKY_daemon_AnnounceTxnsMessage_Process(_atm *C.AnnounceTxnsMessage, _d *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	atm := (*AnnounceTxnsMessage)(unsafe.Pointer(_atm))
	d := (*Daemon)(unsafe.Pointer(_d))
	atm.Process(d)
	return
}

// export SKY_daemon_NewGetTxnsMessage
func SKY_daemon_NewGetTxnsMessage(_txns *C.GoSlice_, _arg1 *C.GetTxnsMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := daemon.NewGetTxnsMessage(txns)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofGetTxnsMessage))
	return
}

// export SKY_daemon_GetTxnsMessage_Handle
func SKY_daemon_GetTxnsMessage_Handle(_gtm *C.GetTxnsMessage, _mc *C.MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gtm := (*GetTxnsMessage)(unsafe.Pointer(_gtm))
	____return_err := gtm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_GetTxnsMessage_Process
func SKY_daemon_GetTxnsMessage_Process(_gtm *C.GetTxnsMessage, _d *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gtm := (*GetTxnsMessage)(unsafe.Pointer(_gtm))
	d := (*Daemon)(unsafe.Pointer(_d))
	gtm.Process(d)
	return
}

// export SKY_daemon_NewGiveTxnsMessage
func SKY_daemon_NewGiveTxnsMessage(_txns *C.Transactions, _arg1 *C.GiveTxnsMessage) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := daemon.NewGiveTxnsMessage(txns)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofGiveTxnsMessage))
	return
}

// export SKY_daemon_GiveTxnsMessage_GetTxns
func SKY_daemon_GiveTxnsMessage_GetTxns(_gtm *C.GiveTxnsMessage, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gtm := (*GiveTxnsMessage)(unsafe.Pointer(_gtm))
	__arg0 := gtm.GetTxns()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_GiveTxnsMessage_Handle
func SKY_daemon_GiveTxnsMessage_Handle(_gtm *C.GiveTxnsMessage, _mc *C.MessageContext, _daemon interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gtm := (*GiveTxnsMessage)(unsafe.Pointer(_gtm))
	____return_err := gtm.Handle(mc, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_GiveTxnsMessage_Process
func SKY_daemon_GiveTxnsMessage_Process(_gtm *C.GiveTxnsMessage, _d *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gtm := (*GiveTxnsMessage)(unsafe.Pointer(_gtm))
	d := (*Daemon)(unsafe.Pointer(_d))
	gtm.Process(d)
	return
}
