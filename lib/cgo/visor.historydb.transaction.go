package main

import (
	historydb "github.com/skycoin/skycoin/src/visor/historydb"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_historydb_Transaction_Hash
func SKY_historydb_Transaction_Hash(_tx *C.historydb__Transaction, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tx := (*historydb.Transaction)(unsafe.Pointer(_tx))
	__arg0 := tx.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}
