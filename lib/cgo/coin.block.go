package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_coin_SignedBlock_VerifySignature
func SKY_coin_SignedBlock_VerifySignature(_b *C.coin__SignedBlock, _pubkey *C.cipher__PubKey) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.SignedBlock)(unsafe.Pointer(_b))
	pubkey := *(*cipher.PubKey)(unsafe.Pointer(_pubkey))
	____return_err := b.VerifySignature(pubkey)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_coin_NewBlock
func SKY_coin_NewBlock(_prev *C.coin__Block, _currentTime uint64, _uxHash *C.cipher__SHA256, _txns *C.coin__Transactions, _calc *C.coin__FeeCalculator, _arg5 *C.coin__Block) (____error_code uint32) {
	____error_code = 0
	prev := *(*coin.Block)(unsafe.Pointer(_prev))
	currentTime := _currentTime
	uxHash := *(*cipher.SHA256)(unsafe.Pointer(_uxHash))
	txns := *(*coin.Transactions)(unsafe.Pointer(_txns))
	calc := *(*coin.FeeCalculator)(unsafe.Pointer(_calc))
	__arg5, ____return_err := coin.NewBlock(prev, currentTime, uxHash, txns, calc)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg5)[:]), unsafe.Pointer(_arg5), uint(SizeofBlock))
	}
	return
}

// export SKY_coin_NewGenesisBlock
func SKY_coin_NewGenesisBlock(_genesisAddr *C.cipher__Address, _genesisCoins, _timestamp uint64, _arg2 *C.coin__Block) (____error_code uint32) {
	____error_code = 0
	genesisAddr := *(*cipher.Address)(unsafe.Pointer(_genesisAddr))
	genesisCoins := _genesisCoins
	timestamp := _timestamp
	__arg2, ____return_err := coin.NewGenesisBlock(genesisAddr, genesisCoins, timestamp)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofBlock))
	}
	return
}

// export SKY_coin_Block_HashHeader
func SKY_coin_Block_HashHeader(_b *C.coin__Block, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.HashHeader()
	return
}

// export SKY_coin_Block_PreHashHeader
func SKY_coin_Block_PreHashHeader(_b *C.coin__Block, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.PreHashHeader()
	return
}

// export SKY_coin_Block_Time
func SKY_coin_Block_Time(_b *C.coin__Block, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.Time()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Block_Seq
func SKY_coin_Block_Seq(_b *C.coin__Block, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.Seq()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Block_HashBody
func SKY_coin_Block_HashBody(_b *C.coin__Block, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.HashBody()
	return
}

// export SKY_coin_Block_Size
func SKY_coin_Block_Size(_b *C.coin__Block, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.Size()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Block_String
func SKY_coin_Block_String(_b *C.coin__Block, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_coin_Block_GetTransaction
func SKY_coin_Block_GetTransaction(_b *C.coin__Block, _txHash *C.cipher__SHA256, _arg1 *C.coin__Transaction, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	b := *(*coin.Block)(unsafe.Pointer(_b))
	txHash := *(*cipher.SHA256)(unsafe.Pointer(_txHash))
	__arg1, __arg2 := b.GetTransaction(txHash)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	*_arg2 = __arg2
	return
}

// export SKY_coin_NewBlockHeader
func SKY_coin_NewBlockHeader(_prev *C.coin__BlockHeader, _uxHash *C.cipher__SHA256, _currentTime, _fee uint64, _body *C.coin__BlockBody, _arg4 *C.coin__BlockHeader) (____error_code uint32) {
	____error_code = 0
	prev := *(*coin.BlockHeader)(unsafe.Pointer(_prev))
	uxHash := *(*cipher.SHA256)(unsafe.Pointer(_uxHash))
	currentTime := _currentTime
	fee := _fee
	body := *(*coin.BlockBody)(unsafe.Pointer(_body))
	__arg4 := coin.NewBlockHeader(prev, uxHash, currentTime, fee, body)
	copyToBuffer(reflect.ValueOf(__arg4[:]), unsafe.Pointer(_arg4), uint(SizeofBlockHeader))
	return
}

// export SKY_coin_BlockHeader_Hash
func SKY_coin_BlockHeader_Hash(_bh *C.coin__BlockHeader, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.Hash()
	return
}

// export SKY_coin_BlockHeader_Bytes
func SKY_coin_BlockHeader_Bytes(_bh *C.coin__BlockHeader, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_BlockHeader_String
func SKY_coin_BlockHeader_String(_bh *C.coin__BlockHeader, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_coin_BlockBody_Hash
func SKY_coin_BlockBody_Hash(_bb *C.coin__BlockBody, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	bb := *(*coin.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Hash()
	return
}

// export SKY_coin_BlockBody_Size
func SKY_coin_BlockBody_Size(_bb *C.coin__BlockBody, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	bb := *(*coin.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Size()
	*_arg0 = __arg0
	return
}

// export SKY_coin_BlockBody_Bytes
func SKY_coin_BlockBody_Bytes(_bb *C.coin__BlockBody, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	bb := *(*coin.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_CreateUnspents
func SKY_coin_CreateUnspents(_bh *C.coin__BlockHeader, _tx *C.coin__Transaction, _arg2 *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	tx := *(*coin.Transaction)(unsafe.Pointer(_tx))
	__arg2 := coin.CreateUnspents(bh, tx)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofUxArray))
	return
}

// export SKY_coin_CreateUnspent
func SKY_coin_CreateUnspent(_bh *C.coin__BlockHeader, _tx *C.coin__Transaction, _outIndex int, _arg3 *C.coin__UxOut) (____error_code uint32) {
	____error_code = 0
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	tx := *(*coin.Transaction)(unsafe.Pointer(_tx))
	outIndex := _outIndex
	__arg3, ____return_err := coin.CreateUnspent(bh, tx, outIndex)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg3[:]), unsafe.Pointer(_arg3), uint(SizeofUxOut))
	}
	return
}
