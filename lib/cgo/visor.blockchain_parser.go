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

// export SKY_visor_NewBlockchainParser
func SKY_visor_NewBlockchainParser(_hisDB *C.HistoryDB, _bc *C.Blockchain, _ops ...*C.ParserOption, _arg3 *C.BlockchainParser) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bc := (*Blockchain)(unsafe.Pointer(_bc))
	ops := _ops
	__arg3 := visor.NewBlockchainParser(hisDB, bc, ops)
	copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofBlockchainParser))
	return
}

// export SKY_visor_BlockchainParser_FeedBlock
func SKY_visor_BlockchainParser_FeedBlock(_bcp *C.BlockchainParser, _b *C.Block) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bcp := (*BlockchainParser)(unsafe.Pointer(_bcp))
	bcp.FeedBlock(b)
	return
}

// export SKY_visor_BlockchainParser_Run
func SKY_visor_BlockchainParser_Run(_bcp *C.BlockchainParser) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bcp := (*BlockchainParser)(unsafe.Pointer(_bcp))
	____return_err := bcp.Run()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_BlockchainParser_Shutdown
func SKY_visor_BlockchainParser_Shutdown(_bcp *C.BlockchainParser) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	bcp := (*BlockchainParser)(unsafe.Pointer(_bcp))
	bcp.Shutdown()
	return
}
