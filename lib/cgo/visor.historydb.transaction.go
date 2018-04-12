package main

import "unsafe"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_historydb_Transaction_Hash
func SKY_historydb_Transaction_Hash(_tx *C.historydb__Transaction, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	tx := (*historydb.Transaction)(unsafe.Pointer(_tx))
	__arg0 := tx.Hash()
	return
}
