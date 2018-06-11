package main

import (
	"unsafe"

	api "github.com/skycoin/skycoin/src/api"
	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_api_NewCreateTransactionResponse
func SKY_api_NewCreateTransactionResponse(_txn C.Transaction__Handle, _inputs []C.wallet__UxBalance, _arg2 *C.CreateTransactionResponse__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(_txn)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	inputs := *(*[]wallet.UxBalance)(unsafe.Pointer(&_inputs))
	__arg2, ____return_err := api.NewCreateTransactionResponse(txn, inputs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerCreateTransactionResponseHandle(__arg2)
	}
	return
}

//export SKY_api_NewCreatedTransaction
func SKY_api_NewCreatedTransaction(_txn C.Transaction__Handle, _inputs []C.wallet__UxBalance, _arg2 *C.CreatedTransaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(_txn)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	inputs := *(*[]wallet.UxBalance)(unsafe.Pointer(&_inputs))
	__arg2, ____return_err := api.NewCreatedTransaction(txn, inputs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerCreatedTransactionHandle(__arg2)
	}
	return
}

//export SKY_api_CreatedTransaction_ToTransaction
func SKY_api_CreatedTransaction_ToTransaction(_r C.CreatedTransaction__Handle, _arg0 *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	r, ok := lookupCreatedTransactionHandle(_r)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := r.ToTransaction()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerTransactionHandle(__arg0)
	}
	return
}

//export SKY_api_NewCreatedTransactionOutput
func SKY_api_NewCreatedTransactionOutput(_out *C.coin__TransactionOutput, _txid *C.cipher__SHA256, _arg2 *C.api__CreatedTransactionOutput) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	out := *(*coin.TransactionOutput)(unsafe.Pointer(_out))
	txid := *(*cipher.SHA256)(unsafe.Pointer(_txid))
	__arg2, ____return_err := api.NewCreatedTransactionOutput(out, txid)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = *(*C.api__CreatedTransactionOutput)(unsafe.Pointer(__arg2))
	}
	return
}

//export SKY_api_NewCreatedTransactionInput
func SKY_api_NewCreatedTransactionInput(_out *C.wallet__UxBalance, _arg1 *C.CreatedTransactionInput__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	out := *(*wallet.UxBalance)(unsafe.Pointer(_out))
	__arg1, ____return_err := api.NewCreatedTransactionInput(out)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerCreatedTransactionInputHandle(__arg1)
	}
	return
}
