package main

import (
	consensus "github.com/skycoin/skycoin/src/consensus"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_consensus_PriorityQueue_Len
func SKY_consensus_PriorityQueue_Len(_pq *C.consensus__PriorityQueue, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pq := *(*consensus.PriorityQueue)(unsafe.Pointer(_pq))
	__arg0 := pq.Len()
	*_arg0 = __arg0
	return
}

//export SKY_consensus_PriorityQueue_Less
func SKY_consensus_PriorityQueue_Less(_pq *C.consensus__PriorityQueue, _i int, _j int, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pq := *(*consensus.PriorityQueue)(unsafe.Pointer(_pq))
	i := _i
	j := _j
	__arg2 := pq.Less(i, j)
	*_arg2 = __arg2
	return
}

//export SKY_consensus_PriorityQueue_Swap
func SKY_consensus_PriorityQueue_Swap(_pq *C.consensus__PriorityQueue, _i int, _j int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pq := *(*consensus.PriorityQueue)(unsafe.Pointer(_pq))
	i := _i
	j := _j
	pq.Swap(i, j)
	return
}

//export SKY_consensus_BlockStatQueue_Len
func SKY_consensus_BlockStatQueue_Len(_self *C.consensus__BlockStatQueue, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.BlockStatQueue)(unsafe.Pointer(_self))
	__arg0 := self.Len()
	*_arg0 = __arg0
	return
}

//export SKY_consensus_BlockStatQueue_Print
func SKY_consensus_BlockStatQueue_Print(_self *C.consensus__BlockStatQueue) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*consensus.BlockStatQueue)(unsafe.Pointer(_self))
	self.Print()
	return
}
