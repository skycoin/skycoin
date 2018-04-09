package main

import (
	blockdb "github.com/skycoin/skycoin/src/visor/blockdb"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_blockdb_ErrMissingSignature_Error
func SKY_blockdb_ErrMissingSignature_Error(_e *C.ErrMissingSignature, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := *(*ErrMissingSignature)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_blockdb_NewBlockchain
func SKY_blockdb_NewBlockchain(_db *C.DB, _walker *C.Walker, _arg2 *C.Blockchain) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	walker := *(*Walker)(unsafe.Pointer(_walker))
	__arg2, ____return_err := blockdb.NewBlockchain(db, walker)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofBlockchain))
	}
	return
}

// export SKY_blockdb_Blockchain_AddBlockWithTx
func SKY_blockdb_Blockchain_AddBlockWithTx(_bc *C.Blockchain, _tx *C.Tx, _sb *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	____return_err := bc.AddBlockWithTx(tx, sb)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_Blockchain_Head
func SKY_blockdb_Blockchain_Head(_bc *C.Blockchain, _arg0 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0, ____return_err := bc.Head()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_Blockchain_HeadSeq
func SKY_blockdb_Blockchain_HeadSeq(_bc *C.Blockchain, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.HeadSeq()
	*_arg0 = __arg0
	return
}

// export SKY_blockdb_Blockchain_UnspentPool
func SKY_blockdb_Blockchain_UnspentPool(_bc *C.Blockchain, _arg0 *C.UnspentPool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.UnspentPool()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofUnspentPool))
	return
}

// export SKY_blockdb_Blockchain_Len
func SKY_blockdb_Blockchain_Len(_bc *C.Blockchain, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.Len()
	*_arg0 = __arg0
	return
}

// export SKY_blockdb_Blockchain_GetBlockByHash
func SKY_blockdb_Blockchain_GetBlockByHash(_bc *C.Blockchain, _hash *C.SHA256, _arg1 *C.SignedBlock) (____error_code uint32) {
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

// export SKY_blockdb_Blockchain_GetBlockBySeq
func SKY_blockdb_Blockchain_GetBlockBySeq(_bc *C.Blockchain, _seq uint64, _arg1 *C.SignedBlock) (____error_code uint32) {
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

// export SKY_blockdb_Blockchain_GetGenesisBlock
func SKY_blockdb_Blockchain_GetGenesisBlock(_bc *C.Blockchain, _arg0 *C.SignedBlock) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	__arg0 := bc.GetGenesisBlock()
	return
}
