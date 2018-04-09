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

// export SKY_visor_DefaultWalker
func SKY_visor_DefaultWalker(_hps *C.GoSlice_, _arg1 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := visor.DefaultWalker(hps)
	return
}

// export SKY_visor_NewBlockchain
func SKY_visor_NewBlockchain(_db *C.DB, _pubkey *C.PubKey, _ops ...*C.Option, _arg3 *C.Blockchain) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ops := _ops
	__arg3, ____return_err := visor.NewBlockchain(db, pubkey, ops)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofBlockchain))
	}
	return
}

// export SKY_visor_Arbitrating
func SKY_visor_Arbitrating(_enable bool, _arg1 *C.Option) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	enable := _enable
	__arg1 := visor.Arbitrating(enable)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofOption))
	return
}

// export SKY_visor_Blockchain_GetGenesisBlock
func SKY_visor_Blockchain_GetGenesisBlock(_bc *C.Blockchain, _arg0 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.GetGenesisBlock()
	return
}

// export SKY_visor_Blockchain_GetBlockByHash
func SKY_visor_Blockchain_GetBlockByHash(_bc *C.Blockchain, _hash *C.SHA256, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg1, ____return_err := bc.GetBlockByHash(hash)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Blockchain_GetBlockBySeq
func SKY_visor_Blockchain_GetBlockBySeq(_bc *C.Blockchain, _seq uint64, _arg1 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	seq := _seq
	__arg1, ____return_err := bc.GetBlockBySeq(seq)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Blockchain_Unspent
func SKY_visor_Blockchain_Unspent(_bc *C.Blockchain, _arg0 *C.UnspentPool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.Unspent()
	return
}

// export SKY_visor_Blockchain_Len
func SKY_visor_Blockchain_Len(_bc *C.Blockchain, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.Len()
	*_arg0 = __arg0
	return
}

// export SKY_visor_Blockchain_Head
func SKY_visor_Blockchain_Head(_bc *C.Blockchain, _arg0 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	__arg0, ____return_err := bc.Head()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Blockchain_HeadSeq
func SKY_visor_Blockchain_HeadSeq(_bc *C.Blockchain, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.HeadSeq()
	*_arg0 = __arg0
	return
}

// export SKY_visor_Blockchain_Time
func SKY_visor_Blockchain_Time(_bc *C.Blockchain, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.Time()
	*_arg0 = __arg0
	return
}

// export SKY_visor_Blockchain_NewBlock
func SKY_visor_Blockchain_NewBlock(_bc *C.Blockchain, _txns *C.Transactions, _currentTime uint64, _arg2 *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	currentTime := _currentTime
	__arg2, ____return_err := bc.NewBlock(txns, currentTime)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Blockchain_ExecuteBlockWithTx
func SKY_visor_Blockchain_ExecuteBlockWithTx(_bc *C.Blockchain, _tx *C.Tx, _sb *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	____return_err := bc.ExecuteBlockWithTx(tx, sb)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Blockchain_VerifyBlockTxnConstraints
func SKY_visor_Blockchain_VerifyBlockTxnConstraints(_bc *C.Blockchain, _tx *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	____return_err := bc.VerifyBlockTxnConstraints(tx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Blockchain_VerifySingleTxnHardConstraints
func SKY_visor_Blockchain_VerifySingleTxnHardConstraints(_bc *C.Blockchain, _tx *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	____return_err := bc.VerifySingleTxnHardConstraints(tx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Blockchain_VerifySingleTxnAllConstraints
func SKY_visor_Blockchain_VerifySingleTxnAllConstraints(_bc *C.Blockchain, _tx *C.Transaction, _maxSize int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	maxSize := _maxSize
	____return_err := bc.VerifySingleTxnAllConstraints(tx, maxSize)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_Blockchain_GetBlocks
func SKY_visor_Blockchain_GetBlocks(_bc *C.Blockchain, _start, _end uint64, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	start := _start
	end := _end
	__arg1 := bc.GetBlocks(start, end)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_visor_Blockchain_GetLastBlocks
func SKY_visor_Blockchain_GetLastBlocks(_bc *C.Blockchain, _num uint64, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	num := _num
	__arg1 := bc.GetLastBlocks(num)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_visor_Blockchain_TransactionFee
func SKY_visor_Blockchain_TransactionFee(_bc *C.Blockchain, _t *C.Transaction, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := *(*Blockchain)(unsafe.Pointer(_bc))
	__arg1, ____return_err := bc.TransactionFee(t)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_visor_Blockchain_BindListener
func SKY_visor_Blockchain_BindListener(_bc *C.Blockchain, _ls *C.BlockListener) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	ls := *(*BlockListener)(unsafe.Pointer(_ls))
	bc.BindListener(ls)
	return
}

// export SKY_visor_Blockchain_Notify
func SKY_visor_Blockchain_Notify(_bc *C.Blockchain, _b *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	bc.Notify(b)
	return
}

// export SKY_visor_Blockchain_UpdateDB
func SKY_visor_Blockchain_UpdateDB(_bc *C.Blockchain, _f C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	____return_err := bc.UpdateDB(f)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
