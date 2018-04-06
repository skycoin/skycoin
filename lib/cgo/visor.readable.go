package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	visor "github.com/skycoin/skycoin/src/visor"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_NewBlockchainMetadata
func SKY_visor_NewBlockchainMetadata(_v *C.Visor, _arg1 *C.BlockchainMetadata) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg1 := visor.NewBlockchainMetadata(v)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofBlockchainMetadata))
	return
}

// export SKY_visor_NewUnconfirmedTransactionStatus
func SKY_visor_NewUnconfirmedTransactionStatus(_arg0 *C.TransactionStatus) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := visor.NewUnconfirmedTransactionStatus()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofTransactionStatus))
	return
}

// export SKY_visor_NewUnknownTransactionStatus
func SKY_visor_NewUnknownTransactionStatus(_arg0 *C.TransactionStatus) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := visor.NewUnknownTransactionStatus()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofTransactionStatus))
	return
}

// export SKY_visor_NewConfirmedTransactionStatus
func SKY_visor_NewConfirmedTransactionStatus(_height uint64, _blockSeq uint64, _arg2 *C.TransactionStatus) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	height := _height
	blockSeq := _blockSeq
	__arg2 := visor.NewConfirmedTransactionStatus(height, blockSeq)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofTransactionStatus))
	return
}

// export SKY_visor_NewReadableTransactionOutput
func SKY_visor_NewReadableTransactionOutput(_t *C.TransactionOutput, _txid *C.SHA256, _arg2 *C.ReadableTransactionOutput) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg2, ____return_err := visor.NewReadableTransactionOutput(t, txid)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofReadableTransactionOutput))
	}
	return
}

// export SKY_visor_NewReadableTransactionInput
func SKY_visor_NewReadableTransactionInput(_uxID, _ownerAddress string, _coins, _hours uint64, _arg2 *C.ReadableTransactionInput) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	uxID := _uxID
	ownerAddress := _ownerAddress
	coins := _coins
	hours := _hours
	__arg2, ____return_err := visor.NewReadableTransactionInput(uxID, ownerAddress, coins, hours)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofReadableTransactionInput))
	}
	return
}

// export SKY_visor_ReadableOutputs_Balance
func SKY_visor_ReadableOutputs_Balance(_ros *C.ReadableOutputs, _arg0 *C.Balance) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ros := *(*cipher.ReadableOutputs)(unsafe.Pointer(_ros))
	__arg0, ____return_err := ros.Balance()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_ReadableOutputs_ToUxArray
func SKY_visor_ReadableOutputs_ToUxArray(_ros *C.ReadableOutputs, _arg0 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ros := *(*cipher.ReadableOutputs)(unsafe.Pointer(_ros))
	__arg0, ____return_err := ros.ToUxArray()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_ReadableOutputSet_SpendableOutputs
func SKY_visor_ReadableOutputSet_SpendableOutputs(_os *C.ReadableOutputSet, _arg0 *C.ReadableOutputs) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	os := *(*cipher.ReadableOutputSet)(unsafe.Pointer(_os))
	__arg0 := os.SpendableOutputs()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofReadableOutputs))
	return
}

// export SKY_visor_ReadableOutputSet_ExpectedOutputs
func SKY_visor_ReadableOutputSet_ExpectedOutputs(_os *C.ReadableOutputSet, _arg0 *C.ReadableOutputs) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	os := *(*cipher.ReadableOutputSet)(unsafe.Pointer(_os))
	__arg0 := os.ExpectedOutputs()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofReadableOutputs))
	return
}

// export SKY_visor_ReadableOutputSet_AggregateUnspentOutputs
func SKY_visor_ReadableOutputSet_AggregateUnspentOutputs(_os *C.ReadableOutputSet, _arg0 map[string]uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	os := *(*cipher.ReadableOutputSet)(unsafe.Pointer(_os))
	__arg0, ____return_err := os.AggregateUnspentOutputs()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_NewReadableOutput
func SKY_visor_NewReadableOutput(_headTime uint64, _t *C.UxOut, _arg2 *C.ReadableOutput) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	headTime := _headTime
	__arg2, ____return_err := visor.NewReadableOutput(headTime, t)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofReadableOutput))
	}
	return
}

// export SKY_visor_NewReadableOutputs
func SKY_visor_NewReadableOutputs(_headTime uint64, _uxs *C.UxArray, _arg2 *C.ReadableOutputs) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	headTime := _headTime
	__arg2, ____return_err := visor.NewReadableOutputs(headTime, uxs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofReadableOutputs))
	}
	return
}

// export SKY_visor_ReadableOutputsToUxBalances
func SKY_visor_ReadableOutputsToUxBalances(_ros *C.ReadableOutputs, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ros := *(*cipher.ReadableOutputs)(unsafe.Pointer(_ros))
	__arg1, ____return_err := visor.ReadableOutputsToUxBalances(ros)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_visor_NewReadableUnconfirmedTxn
func SKY_visor_NewReadableUnconfirmedTxn(_unconfirmed *C.UnconfirmedTxn, _arg1 *C.ReadableUnconfirmedTxn) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	unconfirmed := (*cipher.UnconfirmedTxn)(unsafe.Pointer(_unconfirmed))
	__arg1, ____return_err := visor.NewReadableUnconfirmedTxn(unconfirmed)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableUnconfirmedTxn))
	}
	return
}

// export SKY_visor_NewReadableUnconfirmedTxns
func SKY_visor_NewReadableUnconfirmedTxns(_txs *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	txs := *(*[]cipher.UnconfirmedTxn)(unsafe.Pointer(_txs))
	__arg1, ____return_err := visor.NewReadableUnconfirmedTxns(txs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_visor_NewGenesisReadableTransaction
func SKY_visor_NewGenesisReadableTransaction(_t *C.Transaction, _arg1 *C.ReadableTransaction) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	t := (*cipher.Transaction)(unsafe.Pointer(_t))
	__arg1, ____return_err := visor.NewGenesisReadableTransaction(t)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableTransaction))
	}
	return
}

// export SKY_visor_NewReadableTransaction
func SKY_visor_NewReadableTransaction(_t *C.Transaction, _arg1 *C.ReadableTransaction) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	t := (*cipher.Transaction)(unsafe.Pointer(_t))
	__arg1, ____return_err := visor.NewReadableTransaction(t)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableTransaction))
	}
	return
}

// export SKY_visor_NewReadableBlockHeader
func SKY_visor_NewReadableBlockHeader(_b *C.BlockHeader, _arg1 *C.ReadableBlockHeader) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1 := visor.NewReadableBlockHeader(b)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableBlockHeader))
	return
}

// export SKY_visor_NewReadableBlockBody
func SKY_visor_NewReadableBlockBody(_b *C.Block, _arg1 *C.ReadableBlockBody) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1, ____return_err := visor.NewReadableBlockBody(b)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableBlockBody))
	}
	return
}

// export SKY_visor_NewReadableBlock
func SKY_visor_NewReadableBlock(_b *C.Block, _arg1 *C.ReadableBlock) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1, ____return_err := visor.NewReadableBlock(b)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableBlock))
	}
	return
}

// export SKY_visor_NewReadableBlocks
func SKY_visor_NewReadableBlocks(_blocks *C.GoSlice_, _arg1 *C.ReadableBlocks) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1, ____return_err := visor.NewReadableBlocks(blocks)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableBlocks))
	}
	return
}

// export SKY_visor_NewTxOutputJSON
func SKY_visor_NewTxOutputJSON(_ux *C.TransactionOutput, _srcTx *C.SHA256, _arg2 *C.TransactionOutputJSON) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg2, ____return_err := visor.NewTxOutputJSON(ux, srcTx)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofTransactionOutputJSON))
	}
	return
}

// export SKY_visor_TransactionToJSON
func SKY_visor_TransactionToJSON(_tx *C.Transaction, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1, ____return_err := visor.TransactionToJSON(tx)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}
