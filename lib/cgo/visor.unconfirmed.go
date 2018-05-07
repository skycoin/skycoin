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

//export SKY_visor_TxnUnspents_AllForAddress
func SKY_visor_TxnUnspents_AllForAddress(_tus *C.visor__TxnUnspents, _a *C.cipher__Address, _arg1 *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tus := *(*visor.TxnUnspents)(unsafe.Pointer(_tus))
	a := *(*cipher.Address)(unsafe.Pointer(_a))
	__arg1 := tus.AllForAddress(a)
	*_arg1 = *(*C.coin__UxArray)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_visor_UnconfirmedTxn_Hash
func SKY_visor_UnconfirmedTxn_Hash(_ut *C.visor__UnconfirmedTxn, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ut := (*visor.UnconfirmedTxn)(unsafe.Pointer(_ut))
	__arg0 := ut.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_visor_IsValid
func SKY_visor_IsValid(_tx *C.visor__UnconfirmedTxn, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tx := *(*visor.UnconfirmedTxn)(unsafe.Pointer(_tx))
	__arg1 := visor.IsValid(tx)
	*_arg1 = __arg1
	return
}

//export SKY_visor_All
func SKY_visor_All(_tx *C.visor__UnconfirmedTxn, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tx := *(*visor.UnconfirmedTxn)(unsafe.Pointer(_tx))
	__arg1 := visor.All(tx)
	*_arg1 = __arg1
	return
}
