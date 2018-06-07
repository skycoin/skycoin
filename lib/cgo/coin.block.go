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
func SKY_coin_NewBlock(_b *C.coin__Block, _currentTime uint64, _hash *C.cipher__SHA256, _txns *C.coin__Transactions, _fee uint64, _arg2 *C.coin__Block) (____error_code uint32) {
	feeCalculator := func(t *coin.Transaction) (uint64, error) {
		return _fee, nil
	}

	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	hash := *(*cipher.SHA256)(unsafe.Pointer(_hash))
	txns := *(*coin.Transactions)(unsafe.Pointer(_txns))
	__arg2, ____return_err := coin.NewBlock(b, _currentTime, hash, txns, feeCalculator)
	if ____return_err == nil {
		*_arg2 = *(*C.coin__Block)(unsafe.Pointer(__arg2))
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
func SKY_coin_NewGenesisBlock(_genesisAddr *C.cipher__Address, _genesisCoins, _timestamp uint64, _arg2 *C.coin__Block) (____error_code uint32) {
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
		*_arg2 = *(*C.coin__Block)(unsafe.Pointer(__arg2))
	}
	return
}

//export SKY_coin_Block_HashHeader
func SKY_coin_Block_HashHeader(_b *C.coin__Block, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.HashHeader()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Block_PreHashHeader
func SKY_coin_Block_PreHashHeader(_b *C.coin__Block, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.PreHashHeader()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Block_Time
func SKY_coin_Block_Time(_b *C.coin__Block, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.Time()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Block_Seq
func SKY_coin_Block_Seq(_b *C.coin__Block, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.Seq()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Block_HashBody
func SKY_coin_Block_HashBody(_b *C.coin__Block, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.HashBody()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Block_Size
func SKY_coin_Block_Size(_b *C.coin__Block, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.Size()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Block_String
func SKY_coin_Block_String(_b *C.coin__Block, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	__arg0 := b.String()
	copyString(__arg0, _arg0)
	return
}

//export SKY_coin_Block_GetTransaction
func SKY_coin_Block_GetTransaction(_b *C.coin__Block, _txHash *C.cipher__SHA256, _arg1 *C.coin__Transaction, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*coin.Block)(unsafe.Pointer(_b))
	txHash := *(*cipher.SHA256)(unsafe.Pointer(_txHash))
	__arg1, __arg2 := b.GetTransaction(txHash)
	*_arg1 = *(*C.coin__Transaction)(unsafe.Pointer(&__arg1))
	*_arg2 = __arg2
	return
}

//export SKY_coin_NewBlockHeader
func SKY_coin_NewBlockHeader(_prev *C.coin__BlockHeader, _uxHash *C.cipher__SHA256, _currentTime, _fee uint64, _body *C.coin__BlockBody, _arg4 *C.coin__BlockHeader) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	prev := *(*coin.BlockHeader)(unsafe.Pointer(_prev))
	uxHash := *(*cipher.SHA256)(unsafe.Pointer(_uxHash))
	currentTime := _currentTime
	fee := _fee
	body := *(*coin.BlockBody)(unsafe.Pointer(_body))
	__arg4 := coin.NewBlockHeader(prev, uxHash, currentTime, fee, body)
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
func SKY_coin_BlockBody_Hash(_bb *C.coin__BlockBody, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bb := *(*coin.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_BlockBody_Size
func SKY_coin_BlockBody_Size(_bb *C.coin__BlockBody, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bb := *(*coin.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Size()
	*_arg0 = __arg0
	return
}

//export SKY_coin_BlockBody_Bytes
func SKY_coin_BlockBody_Bytes(_bb *C.coin__BlockBody, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bb := *(*coin.BlockBody)(unsafe.Pointer(_bb))
	__arg0 := bb.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_CreateUnspents
func SKY_coin_CreateUnspents(_bh *C.coin__BlockHeader, _tx *C.coin__Transaction, _arg2 *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	tx := *(*coin.Transaction)(unsafe.Pointer(_tx))
	__arg2 := coin.CreateUnspents(bh, tx)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

//export SKY_coin_CreateUnspent
func SKY_coin_CreateUnspent(_bh *C.coin__BlockHeader, _tx *C.coin__Transaction, _outIndex int, _arg3 *C.coin__UxOut) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bh := *(*coin.BlockHeader)(unsafe.Pointer(_bh))
	tx := *(*coin.Transaction)(unsafe.Pointer(_tx))
	outIndex := _outIndex
	__arg3, ____return_err := coin.CreateUnspent(bh, tx, outIndex)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg3 = *(*C.coin__UxOut)(unsafe.Pointer(&__arg3))
	}
	return
}
