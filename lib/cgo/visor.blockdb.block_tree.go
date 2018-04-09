package main

import "unsafe"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_blockdb_blockTree_AddBlock
func SKY_blockdb_blockTree_AddBlock(_bt blockTree, _b *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bt := (*blockTree)(unsafe.Pointer(_bt))
	____return_err := bt.AddBlock(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_blockTree_AddBlockWithTx
func SKY_blockdb_blockTree_AddBlockWithTx(_bt blockTree, _tx *C.Tx, _b *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bt := (*blockTree)(unsafe.Pointer(_bt))
	____return_err := bt.AddBlockWithTx(tx, b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_blockTree_RemoveBlock
func SKY_blockdb_blockTree_RemoveBlock(_bt blockTree, _b *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bt := (*blockTree)(unsafe.Pointer(_bt))
	____return_err := bt.RemoveBlock(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_blockdb_blockTree_GetBlock
func SKY_blockdb_blockTree_GetBlock(_bt blockTree, _hash *C.SHA256, _arg1 *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bt := (*blockTree)(unsafe.Pointer(_bt))
	__arg1 := bt.GetBlock(hash)
	return
}

// export SKY_blockdb_blockTree_GetBlockInDepth
func SKY_blockdb_blockTree_GetBlockInDepth(_bt blockTree, _depth uint64, _filter C.Handle, _arg2 *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bt := (*blockTree)(unsafe.Pointer(_bt))
	depth := _depth
	__arg2 := bt.GetBlockInDepth(depth, filter)
	return
}
