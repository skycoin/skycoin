package main

import (
	"reflect"
	"unsafe"

	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_coin_NewBlock
func SKY_coin_NewBlock(_b C.Block__Handle, _currentTime uint64, _hash *C.cipher__SHA256, _txns C.Transactions__Handle, _fee uint64, _arg2 *C.Block__Handle) (____error_code uint32) {
	feeCalculator := func(t *coin.Transaction) (uint64, error) {
		return _fee, nil
	}

	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	hash := *(*cipher.SHA256)(unsafe.Pointer(_hash))
	txns, ok := lookupTransactionsHandle(_txns)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg2, ____return_err := coin.NewBlock(*b, _currentTime, hash, *txns, feeCalculator)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerBlockHandle(__arg2)
	}
	return
}

//export SKY_coin_SignedBlock_VerifySignature
func SKY_coin_SignedBlock_VerifySignature(_b *C.coin__SignedBlock, _pubkey *C.cipher__PubKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.SignedBlock)(unsafe.Pointer(_b))
	pubkey := *(*cipher.PubKey)(unsafe.Pointer(_pubkey))
	____return_err := b.VerifySignature(pubkey)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_coin_NewGenesisBlock
func SKY_coin_NewGenesisBlock(_genesisAddr *C.cipher__Address, _genesisCoins, _timestamp uint64, _arg2 *C.Block__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	genesisAddr := *(*cipher.Address)(unsafe.Pointer(_genesisAddr))
	genesisCoins := _genesisCoins
	timestamp := _timestamp
	__arg2, ____return_err := coin.NewGenesisBlock(genesisAddr, genesisCoins, timestamp)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerBlockHandle(__arg2)
	}
	return
}

//export SKY_coin_Block_HashHeader
func SKY_coin_Block_HashHeader(_b C.Block__Handle, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := b.HashHeader()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Block_PreHashHeader
func SKY_coin_Block_PreHashHeader(_b C.Block__Handle, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := b.PreHashHeader()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Block_Time
func SKY_coin_Block_Time(_b C.Block__Handle, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := b.Time()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Block_Seq
func SKY_coin_Block_Seq(_b C.Block__Handle, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := b.Seq()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Block_HashBody
func SKY_coin_Block_HashBody(_b C.Block__Handle, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := b.HashBody()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Block_Size
func SKY_coin_Block_Size(_b C.Block__Handle, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := b.Size()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Block_String
func SKY_coin_Block_String(_b C.Block__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := b.String()
	copyString(__arg0, _arg0)
	return
}

//export SKY_coin_Block_GetTransaction
func SKY_coin_Block_GetTransaction(_b C.Block__Handle, _txHash *C.cipher__SHA256, _arg1 *C.Transaction__Handle, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	txHash := *(*cipher.SHA256)(unsafe.Pointer(_txHash))
	__arg1, __arg2 := b.GetTransaction(txHash)
	*_arg1 = registerTransactionHandle(&__arg1)
	*_arg2 = __arg2
	return
}

//export SKY_coin_NewBlockHeader
func SKY_coin_NewBlockHeader(_prev *C.coin__BlockHeader, _uxHash *C.cipher__SHA256, _currentTime, _fee uint64, _body C.BlockBody__Handle, _arg4 *C.coin__BlockHeader) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	prev := *(*coin.BlockHeader)(unsafe.Pointer(_prev))
	uxHash := *(*cipher.SHA256)(unsafe.Pointer(_uxHash))
	currentTime := _currentTime
	fee := _fee
	body, ok := lookupBlockBodyHandle(_body)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg4 := coin.NewBlockHeader(prev, uxHash, currentTime, fee, *body)
	*_arg4 = *(*C.coin__BlockHeader)(unsafe.Pointer(&__arg4))
	return
}

//export SKY_coin_BlockHeader_Hash
func SKY_coin_BlockHeader_Hash(_bh *C.coin__BlockHeader, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_BlockHeader_Bytes
func SKY_coin_BlockHeader_Bytes(_bh *C.coin__BlockHeader, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_BlockHeader_String
func SKY_coin_BlockHeader_String(_bh *C.coin__BlockHeader, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.String()
	copyString(__arg0, _arg0)
	return
}

//export SKY_coin_BlockBody_Hash
func SKY_coin_BlockBody_Hash(_body C.BlockBody__Handle, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	body, ok := lookupBlockBodyHandle(_body)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := body.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_BlockBody_Size
func SKY_coin_BlockBody_Size(_bb *C.BlockBody__Handle, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bb, ok := lookupBlockBodyHandle(*_bb)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := bb.Size()
	*_arg0 = __arg0
	return
}

//export SKY_coin_BlockBody_Bytes
func SKY_coin_BlockBody_Bytes(_bb C.BlockBody__Handle, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bb, ok := lookupBlockBodyHandle(_bb)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := bb.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_CreateUnspents
func SKY_coin_CreateUnspents(_bh *C.coin__BlockHeader, _tx C.Transaction__Handle, _arg2 *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	tx, ok := lookupTransactionHandle(_tx)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg2 := coin.CreateUnspents(bh, *tx)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

//export SKY_coin_CreateUnspent
func SKY_coin_CreateUnspent(_bh *C.coin__BlockHeader, _tx C.Transaction__Handle, _outIndex int, _arg3 *C.coin__UxOut) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	tx, ok := lookupTransactionHandle(_tx)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	outIndex := _outIndex
	__arg3, ____return_err := coin.CreateUnspent(bh, *tx, outIndex)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg3 = *(*C.coin__UxOut)(unsafe.Pointer(&__arg3))
	}
	return
}

//export SKY_coin_GetBlockObject
func SKY_coin_GetBlockObject(_b C.Block__Handle, _p **C.coin__Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
	} else {
		*_p = (*C.coin__Block)(unsafe.Pointer(b))
	}
	return
}

//export SKY_coin_GetBlockBody
func SKY_coin_GetBlockBody(_b C.Block__Handle, _p *C.BlockBody__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b, ok := lookupBlockHandle(_b)
	if !ok {
		____error_code = SKY_ERROR
	} else {
		*_p = registerBlockBodyHandle(&b.Body)
	}
	return
}

//export SKY_coin_NewEmptyBlock
func SKY_coin_NewEmptyBlock(_txns C.Transactions__Handle, handle *C.Block__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns, ok := lookupTransactionsHandle(_txns)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	body := coin.BlockBody{
		Transactions: *txns,
	}
	block := coin.Block{
		Body: body,
		Head: coin.BlockHeader{
			Version:  0x02,
			Time:     100,
			BkSeq:    0,
			Fee:      10,
			PrevHash: cipher.SHA256{},
			BodyHash: body.Hash(),
		}}
	*handle = registerBlockHandle(&block)
	return
}
