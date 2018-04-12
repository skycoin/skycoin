package main

import (
	coin "github.com/skycoin/skycoin/src/coin"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_BlockchainParser_FeedBlock
func SKY_visor_BlockchainParser_FeedBlock(_bcp *C.visor__BlockchainParser, _b *C.coin__Block) (____error_code uint32) {
	____error_code = 0
	bcp := (*visor.BlockchainParser)(unsafe.Pointer(_bcp))
	b := *(*coin.Block)(unsafe.Pointer(_b))
	bcp.FeedBlock(b)
	return
}

// export SKY_visor_BlockchainParser_Run
func SKY_visor_BlockchainParser_Run(_bcp *C.visor__BlockchainParser) (____error_code uint32) {
	____error_code = 0
	bcp := (*visor.BlockchainParser)(unsafe.Pointer(_bcp))
	____return_err := bcp.Run()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_BlockchainParser_Shutdown
func SKY_visor_BlockchainParser_Shutdown(_bcp *C.visor__BlockchainParser) (____error_code uint32) {
	____error_code = 0
	bcp := (*visor.BlockchainParser)(unsafe.Pointer(_bcp))
	bcp.Shutdown()
	return
}
