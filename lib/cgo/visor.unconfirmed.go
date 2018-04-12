package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	visor "github.com/skycoin/skycoin/src/visor"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_TxnUnspents_AllForAddress
func SKY_visor_TxnUnspents_AllForAddress(_tus *C.visor__TxnUnspents, _a *C.cipher__Address, _arg1 *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	tus := *(*visor.TxnUnspents)(unsafe.Pointer(_tus))
	a := *(*cipher.Address)(unsafe.Pointer(_a))
	__arg1 := tus.AllForAddress(a)
	return
}

// export SKY_visor_UnconfirmedTxn_Hash
func SKY_visor_UnconfirmedTxn_Hash(_ut *C.visor__UnconfirmedTxn, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	ut := (*visor.UnconfirmedTxn)(unsafe.Pointer(_ut))
	__arg0 := ut.Hash()
	return
}

// export SKY_visor_IsValid
func SKY_visor_IsValid(_tx *C.visor__UnconfirmedTxn, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	tx := *(*visor.UnconfirmedTxn)(unsafe.Pointer(_tx))
	__arg1 := visor.IsValid(tx)
	*_arg1 = __arg1
	return
}

// export SKY_visor_All
func SKY_visor_All(_tx *C.visor__UnconfirmedTxn, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	tx := *(*visor.UnconfirmedTxn)(unsafe.Pointer(_tx))
	__arg1 := visor.All(tx)
	*_arg1 = __arg1
	return
}
