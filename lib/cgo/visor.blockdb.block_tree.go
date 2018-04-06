package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_blockdb_blockTree_AddBlock
func SKY_blockdb_blockTree_AddBlock(_bt blockTree, _b *C.Block) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bt := (*cipher.blockTree)(unsafe.Pointer(_bt))
	____return_err := bt.AddBlock(b)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_blockTree_AddBlockWithTx
func SKY_blockdb_blockTree_AddBlockWithTx(_bt blockTree, _tx *C.Tx, _b *C.Block) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bt := (*cipher.blockTree)(unsafe.Pointer(_bt))
	____return_err := bt.AddBlockWithTx(tx, b)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_blockTree_RemoveBlock
func SKY_blockdb_blockTree_RemoveBlock(_bt blockTree, _b *C.Block) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bt := (*cipher.blockTree)(unsafe.Pointer(_bt))
	____return_err := bt.RemoveBlock(b)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_blockTree_GetBlock
func SKY_blockdb_blockTree_GetBlock(_bt blockTree, _hash *C.SHA256, _arg1 *C.Block) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bt := (*cipher.blockTree)(unsafe.Pointer(_bt))
	__arg1 := bt.GetBlock(hash)
	return
}

// export SKY_blockdb_blockTree_GetBlockInDepth
func SKY_blockdb_blockTree_GetBlockInDepth(_bt blockTree, _depth uint64, _filter C.Handle, _arg2 *C.Block) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bt := (*cipher.blockTree)(unsafe.Pointer(_bt))
	depth := _depth
	__arg2 := bt.GetBlockInDepth(depth, filter)
	return
}
