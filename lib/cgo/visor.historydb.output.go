package main

import (
	historydb "github.com/skycoin/skycoin/src/visor/historydb"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_historydb_NewUxOutJSON
func SKY_historydb_NewUxOutJSON(_out *C.historydb__UxOut, _arg1 *C.historydb__UxOutJSON) (____error_code uint32) {
	____error_code = 0
	out := (*historydb.UxOut)(unsafe.Pointer(_out))
	__arg1 := historydb.NewUxOutJSON(out)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUxOutJSON))
	return
}

// export SKY_historydb_UxOut_Hash
func SKY_historydb_UxOut_Hash(_o *C.historydb__UxOut, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	o := *(*historydb.UxOut)(unsafe.Pointer(_o))
	__arg0 := o.Hash()
	return
}
