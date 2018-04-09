package main

import (
	visor "github.com/skycoin/skycoin/src/visor"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_MaxDropletDivisor
func SKY_visor_MaxDropletDivisor(_arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := visor.MaxDropletDivisor()
	*_arg0 = __arg0
	return
}

// export SKY_visor_DropletPrecisionCheck
func SKY_visor_DropletPrecisionCheck(_amount uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	amount := _amount
	____return_err := visor.DropletPrecisionCheck(amount)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_NewVisorConfig
func SKY_visor_NewVisorConfig(_arg0 *C.Config) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := visor.NewVisorConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofConfig))
	return
}

// export SKY_visor_Config_Verify
func SKY_visor_Config_Verify(_c *C.Config) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*Config)(unsafe.Pointer(_c))
	____return_err := c.Verify()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_NewVisor
func SKY_visor_NewVisor(_c *C.Config, _db *C.DB, _arg2 *C.Visor) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*Config)(unsafe.Pointer(_c))
	__arg2, ____return_err := visor.NewVisor(c, db)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofVisor))
	}
	return
}

// export SKY_visor_Visor_Run
func SKY_visor_Visor_Run(_vs *C.Visor) (____error_code uint32) {
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

// export SKY_visor_Visor_Shutdown
func SKY_visor_Visor_Shutdown(_vs *C.Visor) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	vs.Shutdown()
	return
}

// export SKY_visor_Visor_GenesisPreconditions
func SKY_visor_Visor_GenesisPreconditions(_vs *C.Visor) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	vs.GenesisPreconditions()
	return
}

// export SKY_visor_Visor_RefreshUnconfirmed
func SKY_visor_Visor_RefreshUnconfirmed(_vs *C.Visor, _arg0 *C.GoSlice_) (____error_code uint32) {
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

// export SKY_visor_Visor_RemoveInvalidUnconfirmed
func SKY_visor_Visor_RemoveInvalidUnconfirmed(_vs *C.Visor, _arg0 *C.GoSlice_) (____error_code uint32) {
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

// export SKY_visor_Visor_CreateBlock
func SKY_visor_Visor_CreateBlock(_vs *C.Visor, _when uint64, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	when := _when
	__arg1, ____return_err := vs.CreateBlock(when)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_CreateAndExecuteBlock
func SKY_visor_Visor_CreateAndExecuteBlock(_vs *C.Visor, _arg0 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0, ____return_err := vs.CreateAndExecuteBlock()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_ExecuteSignedBlock
func SKY_visor_Visor_ExecuteSignedBlock(_vs *C.Visor, _b *C.SignedBlock) (____error_code uint32) {
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

// export SKY_visor_Visor_SignBlock
func SKY_visor_Visor_SignBlock(_vs *C.Visor, _b *C.Block, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1 := vs.SignBlock(b)
	return
}

// export SKY_visor_Visor_GetUnspentOutputs
func SKY_visor_Visor_GetUnspentOutputs(_vs *C.Visor, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0, ____return_err := vs.GetUnspentOutputs()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_visor_Visor_UnconfirmedSpendingOutputs
func SKY_visor_Visor_UnconfirmedSpendingOutputs(_vs *C.Visor, _arg0 *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0, ____return_err := vs.UnconfirmedSpendingOutputs()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_UnconfirmedIncomingOutputs
func SKY_visor_Visor_UnconfirmedIncomingOutputs(_vs *C.Visor, _arg0 *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0, ____return_err := vs.UnconfirmedIncomingOutputs()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_GetSignedBlocksSince
func SKY_visor_Visor_GetSignedBlocksSince(_vs *C.Visor, _seq, _ct uint64, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	seq := _seq
	ct := _ct
	__arg1, ____return_err := vs.GetSignedBlocksSince(seq, ct)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_visor_Visor_HeadBkSeq
func SKY_visor_Visor_HeadBkSeq(_vs *C.Visor, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0 := vs.HeadBkSeq()
	*_arg0 = __arg0
	return
}

// export SKY_visor_Visor_GetBlockchainMetadata
func SKY_visor_Visor_GetBlockchainMetadata(_vs *C.Visor, _arg0 *C.BlockchainMetadata) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0 := vs.GetBlockchainMetadata()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofBlockchainMetadata))
	return
}

// export SKY_visor_Visor_GetBlock
func SKY_visor_Visor_GetBlock(_vs *C.Visor, _seq uint64, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	seq := _seq
	__arg1, ____return_err := vs.GetBlock(seq)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_GetBlocks
func SKY_visor_Visor_GetBlocks(_vs *C.Visor, _start, _end uint64, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	start := _start
	end := _end
	__arg1 := vs.GetBlocks(start, end)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_visor_Visor_InjectTransaction
func SKY_visor_Visor_InjectTransaction(_vs *C.Visor, _txn *C.Transaction, _arg1 *bool, _arg2 *C.ErrTxnViolatesSoftConstraint) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1, __arg2, ____return_err := vs.InjectTransaction(txn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofErrTxnViolatesSoftConstraint))
	}
	return
}

// export SKY_visor_Visor_InjectTransactionStrict
func SKY_visor_Visor_InjectTransactionStrict(_vs *C.Visor, _txn *C.Transaction, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1, ____return_err := vs.InjectTransactionStrict(txn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_visor_Visor_GetAddressTxns
func SKY_visor_Visor_GetAddressTxns(_vs *C.Visor, _a *C.Address, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1, ____return_err := vs.GetAddressTxns(a)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_visor_Visor_GetTransaction
func SKY_visor_Visor_GetTransaction(_vs *C.Visor, _txHash *C.SHA256, _arg1 *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1, ____return_err := vs.GetTransaction(txHash)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	}
	return
}

// export SKY_visor_baseFilter_Match
func SKY_visor_baseFilter_Match(_f baseFilter, _tx *C.Transaction, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	f := *(*baseFilter)(unsafe.Pointer(_f))
	tx := (*Transaction)(unsafe.Pointer(_tx))
	__arg1 := f.Match(tx)
	*_arg1 = __arg1
	return
}

// export SKY_visor_AddrsFilter
func SKY_visor_AddrsFilter(_addrs *C.GoSlice_, _arg1 *C.TxFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := visor.AddrsFilter(addrs)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTxFilter))
	return
}

// export SKY_visor_addrsFilter_Match
func SKY_visor_addrsFilter_Match(_af addrsFilter, _tx *C.Transaction, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	af := *(*addrsFilter)(unsafe.Pointer(_af))
	tx := (*Transaction)(unsafe.Pointer(_tx))
	__arg1 := af.Match(tx)
	*_arg1 = __arg1
	return
}

// export SKY_visor_ConfirmedTxFilter
func SKY_visor_ConfirmedTxFilter(_isConfirmed bool, _arg1 *C.TxFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	isConfirmed := _isConfirmed
	__arg1 := visor.ConfirmedTxFilter(isConfirmed)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTxFilter))
	return
}

// export SKY_visor_Visor_GetTransactions
func SKY_visor_Visor_GetTransactions(_vs *C.Visor, _flts *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1, ____return_err := vs.GetTransactions(flts)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_visor_Visor_AddressBalance
func SKY_visor_Visor_AddressBalance(_vs *C.Visor, _auxs *C.AddressUxOuts, _arg1 *uint64, _arg2 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1, __arg2, ____return_err := vs.AddressBalance(auxs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
		*_arg2 = __arg2
	}
	return
}

// export SKY_visor_Visor_GetUnconfirmedTxns
func SKY_visor_Visor_GetUnconfirmedTxns(_vs *C.Visor, _filter C.Handle, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1 := vs.GetUnconfirmedTxns(filter)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_visor_ToAddresses
func SKY_visor_ToAddresses(_addresses *C.GoSlice_, _arg1 C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := visor.ToAddresses(addresses)
	return
}

// export SKY_visor_Visor_GetAllUnconfirmedTxns
func SKY_visor_Visor_GetAllUnconfirmedTxns(_vs *C.Visor, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0 := vs.GetAllUnconfirmedTxns()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_visor_Visor_GetAllValidUnconfirmedTxHashes
func SKY_visor_Visor_GetAllValidUnconfirmedTxHashes(_vs *C.Visor, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg0 := vs.GetAllValidUnconfirmedTxHashes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_visor_Visor_GetBlockByHash
func SKY_visor_Visor_GetBlockByHash(_vs *C.Visor, _hash *C.SHA256, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	__arg1, ____return_err := vs.GetBlockByHash(hash)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_GetBlockBySeq
func SKY_visor_Visor_GetBlockBySeq(_vs *C.Visor, _seq uint64, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	seq := _seq
	__arg1, ____return_err := vs.GetBlockBySeq(seq)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_GetLastBlocks
func SKY_visor_Visor_GetLastBlocks(_vs *C.Visor, _num uint64, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := (*Visor)(unsafe.Pointer(_vs))
	num := _num
	__arg1 := vs.GetLastBlocks(num)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_visor_Visor_GetHeadBlock
func SKY_visor_Visor_GetHeadBlock(_vs *C.Visor, _arg0 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := *(*Visor)(unsafe.Pointer(_vs))
	__arg0, ____return_err := vs.GetHeadBlock()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_GetUxOutByID
func SKY_visor_Visor_GetUxOutByID(_vs *C.Visor, _id *C.SHA256, _arg1 *C.UxOut) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := *(*Visor)(unsafe.Pointer(_vs))
	__arg1, ____return_err := vs.GetUxOutByID(id)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_GetAddrUxOuts
func SKY_visor_Visor_GetAddrUxOuts(_vs *C.Visor, _address *C.Address, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := *(*Visor)(unsafe.Pointer(_vs))
	__arg1, ____return_err := vs.GetAddrUxOuts(address)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_visor_Visor_ScanAheadWalletAddresses
func SKY_visor_Visor_ScanAheadWalletAddresses(_vs *C.Visor, _wltName string, _password *C.GoSlice_, _scanN uint64, _arg3 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := *(*Visor)(unsafe.Pointer(_vs))
	wltName := _wltName
	password := *(*[]byte)(unsafe.Pointer(_password))
	scanN := _scanN
	__arg3, ____return_err := vs.ScanAheadWalletAddresses(wltName, password, scanN)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Visor_GetBalanceOfAddrs
func SKY_visor_Visor_GetBalanceOfAddrs(_vs *C.Visor, _addrs *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	vs := *(*Visor)(unsafe.Pointer(_vs))
	__arg1, ____return_err := vs.GetBalanceOfAddrs(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}
