package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_coin_SignedBlock_VerifySignature
func SKY_coin_SignedBlock_VerifySignature(_b *C.SignedBlock, _pubkey *C.PubKey) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.SignedBlock)(unsafe.Pointer(_b))
	____return_err := b.VerifySignature(pubkey)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_coin_NewBlock
func SKY_coin_NewBlock(_prev *C.Block, _currentTime uint64, _uxHash *C.SHA256, _txns *C.Transactions, _calc *C.FeeCalculator, _arg5 *C.Block) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	prev := *(*cipher.Block)(unsafe.Pointer(_prev))
	currentTime := _currentTime
	txns := *(*cipher.Transactions)(unsafe.Pointer(_txns))
	calc := *(*cipher.FeeCalculator)(unsafe.Pointer(_calc))
	__arg5, ____return_err := coin.NewBlock(prev, currentTime, uxHash, txns, calc)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg5)[:]), unsafe.Pointer(_arg5), uint(SizeofBlock))
	}
	return
}

// export SKY_coin_NewGenesisBlock
func SKY_coin_NewGenesisBlock(_genesisAddr *C.Address, _genesisCoins, _timestamp uint64, _arg2 *C.Block) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	genesisCoins := _genesisCoins
	timestamp := _timestamp
	__arg2, ____return_err := coin.NewGenesisBlock(genesisAddr, genesisCoins, timestamp)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofBlock))
	}
	return
}

// export SKY_coin_Block_HashHeader
func SKY_coin_Block_HashHeader(_b *C.Block, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Block)(unsafe.Pointer(_b))
	__arg0 := b.HashHeader()
	return
}

// export SKY_coin_Block_PreHashHeader
func SKY_coin_Block_PreHashHeader(_b *C.Block, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Block)(unsafe.Pointer(_b))
	__arg0 := b.PreHashHeader()
	return
}

// export SKY_coin_Block_Time
func SKY_coin_Block_Time(_b *C.Block, _arg0 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Block)(unsafe.Pointer(_b))
	__arg0 := b.Time()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Block_Seq
func SKY_coin_Block_Seq(_b *C.Block, _arg0 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Block)(unsafe.Pointer(_b))
	__arg0 := b.Seq()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Block_HashBody
func SKY_coin_Block_HashBody(_b *C.Block, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Block)(unsafe.Pointer(_b))
	__arg0 := b.HashBody()
	return
}

// export SKY_coin_Block_Size
func SKY_coin_Block_Size(_b *C.Block, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Block)(unsafe.Pointer(_b))
	__arg0 := b.Size()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Block_String
func SKY_coin_Block_String(_b *C.Block, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Block)(unsafe.Pointer(_b))
	__arg0 := b.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_coin_Block_GetTransaction
func SKY_coin_Block_GetTransaction(_b *C.Block, _txHash *C.SHA256, _arg1 *C.Transaction, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Block)(unsafe.Pointer(_b))
	__arg1, __arg2 := b.GetTransaction(txHash)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	*_arg2 = __arg2
	return
}

// export SKY_coin_NewBlockHeader
func SKY_coin_NewBlockHeader(_prev *C.BlockHeader, _uxHash *C.SHA256, _currentTime, _fee uint64, _body *C.BlockBody, _arg4 *C.BlockHeader) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	prev := *(*cipher.BlockHeader)(unsafe.Pointer(_prev))
	currentTime := _currentTime
	fee := _fee
	body := *(*cipher.BlockBody)(unsafe.Pointer(_body))
	__arg4 := coin.NewBlockHeader(prev, uxHash, currentTime, fee, body)
	copyToBuffer(reflect.ValueOf(__arg4[:]), unsafe.Pointer(_arg4), uint(SizeofBlockHeader))
	return
}

// export SKY_coin_BlockHeader_Hash
func SKY_coin_BlockHeader_Hash(_bh *C.BlockHeader, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bh := *(*cipher.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.Hash()
	return
}

// export SKY_coin_BlockHeader_Bytes
func SKY_coin_BlockHeader_Bytes(_bh *C.BlockHeader, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bh := *(*cipher.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_BlockHeader_String
func SKY_coin_BlockHeader_String(_bh *C.BlockHeader, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bh := *(*cipher.BlockHeader)(unsafe.Pointer(_bh))
	__arg0 := bh.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_coin_BlockBody_Hash
func SKY_coin_BlockBody_Hash(_bb *C.BlockBody, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bb := *(*cipher.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Hash()
	return
}

// export SKY_coin_BlockBody_Size
func SKY_coin_BlockBody_Size(_bb *C.BlockBody, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bb := *(*cipher.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Size()
	*_arg0 = __arg0
	return
}

// export SKY_coin_BlockBody_Bytes
func SKY_coin_BlockBody_Bytes(_bb *C.BlockBody, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bb := *(*cipher.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_CreateUnspents
func SKY_coin_CreateUnspents(_bh *C.BlockHeader, _tx *C.Transaction, _arg2 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bh := *(*cipher.BlockHeader)(unsafe.Pointer(_bh))
	tx := *(*cipher.Transaction)(unsafe.Pointer(_tx))
	__arg2 := coin.CreateUnspents(bh, tx)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofUxArray))
	return
}

// export SKY_coin_CreateUnspent
func SKY_coin_CreateUnspent(_bh *C.BlockHeader, _tx *C.Transaction, _outIndex int, _arg3 *C.UxOut) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bh := *(*cipher.BlockHeader)(unsafe.Pointer(_bh))
	tx := *(*cipher.Transaction)(unsafe.Pointer(_tx))
	outIndex := _outIndex
	__arg3, ____return_err := coin.CreateUnspent(bh, tx, outIndex)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg3[:]), unsafe.Pointer(_arg3), uint(SizeofUxOut))
	}
	return
}
