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

//export SKY_historydb_NewUxOutJSON
func SKY_historydb_NewUxOutJSON(_out *C.historydb__UxOut, _arg1 *C.historydb__UxOutJSON) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	out := (*historydb.UxOut)(unsafe.Pointer(_out))
	__arg1 := historydb.NewUxOutJSON(out)
	*_arg1 = *(*C.historydb__UxOutJSON)(unsafe.Pointer(__arg1))
	return
}

//export SKY_historydb_UxOut_Hash
func SKY_historydb_UxOut_Hash(_o *C.historydb__UxOut, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	o := *(*historydb.UxOut)(unsafe.Pointer(_o))
	__arg0 := o.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}
