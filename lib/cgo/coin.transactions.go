package main

import (
	"reflect"
	"unsafe"

	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_coin_Create_Transaction
func SKY_coin_Create_Transaction(handle *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tx := coin.Transaction{}
	*handle = registerTransactionHandle(&tx)
	return
}

//export SKY_coin_Transaction_Copy
func SKY_coin_Transaction_Copy(handle C.Transaction__Handle, handle2 *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tx, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	ntx := coin.Transaction{}
	ntx.Length = tx.Length
	ntx.Type = tx.Type
	ntx.InnerHash = tx.InnerHash
	ntx.Sigs = make([]cipher.Sig, 0)
	ntx.Sigs = append(ntx.Sigs, tx.Sigs...)
	ntx.In = make([]cipher.SHA256, 0)
	ntx.In = append(ntx.In, tx.In...)
	ntx.Out = make([]coin.TransactionOutput, 0)
	ntx.Out = append(ntx.Out, tx.Out...)
	*handle2 = registerTransactionHandle(&ntx)
	return
}

//export SKY_coin_Get_Transaction_Object
func SKY_coin_Get_Transaction_Object(handle C.Transaction__Handle, _pptx **C.coin__Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ptx, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
	} else {
		*_pptx = (*C.coin__Transaction)(unsafe.Pointer(ptx))
	}
	return
}

//export SKY_coin_Transaction_ResetInputs
func SKY_coin_Transaction_ResetInputs(handle C.Transaction__Handle, count int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	txn.In = make([]cipher.SHA256, count)
	return
}

//export SKY_coin_Transaction_ResetOutputs
func SKY_coin_Transaction_ResetOutputs(handle C.Transaction__Handle, count int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	txn.Out = make([]coin.TransactionOutput, count)
	return
}

//export SKY_coin_Transaction_ResetSignatures
func SKY_coin_Transaction_ResetSignatures(handle C.Transaction__Handle, count int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	txn.Sigs = make([]cipher.Sig, count)
	return
}

//export SKY_coin_Transaction_Verify
func SKY_coin_Transaction_Verify(handle C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	____return_err := txn.Verify()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_coin_Transaction_VerifyInput
func SKY_coin_Transaction_VerifyInput(handle C.Transaction__Handle, _uxIn *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	uxIn := *(*coin.UxArray)(unsafe.Pointer(_uxIn))
	____return_err := txn.VerifyInput(uxIn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_coin_Transaction_PushInput
func SKY_coin_Transaction_PushInput(handle C.Transaction__Handle, _uxOut *C.cipher__SHA256, _arg1 *uint16) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	uxOut := *(*cipher.SHA256)(unsafe.Pointer(_uxOut))
	__arg1 := txn.PushInput(uxOut)
	*_arg1 = __arg1
	return
}

//export SKY_coin_TransactionOutput_UxID
func SKY_coin_TransactionOutput_UxID(_txOut *C.coin__TransactionOutput, _txID *C.cipher__SHA256, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txOut := *(*coin.TransactionOutput)(unsafe.Pointer(_txOut))
	txID := *(*cipher.SHA256)(unsafe.Pointer(_txID))
	__arg1 := txOut.UxID(txID)
	*_arg1 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_coin_Transaction_PushOutput
func SKY_coin_Transaction_PushOutput(handle C.Transaction__Handle, _dst *C.cipher__Address, _coins, _hours uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	dst := *(*cipher.Address)(unsafe.Pointer(_dst))
	coins := _coins
	hours := _hours
	txn.PushOutput(dst, coins, hours)
	return
}

//export SKY_coin_Transaction_SignInputs
func SKY_coin_Transaction_SignInputs(handle C.Transaction__Handle, _keys []C.cipher__SecKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	keys := *(*[]cipher.SecKey)(unsafe.Pointer(&_keys))
	txn.SignInputs(keys)
	return
}

//export SKY_coin_Transaction_Size
func SKY_coin_Transaction_Size(handle C.Transaction__Handle, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := txn.Size()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Transaction_Hash
func SKY_coin_Transaction_Hash(handle C.Transaction__Handle, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := txn.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Transaction_SizeHash
func SKY_coin_Transaction_SizeHash(handle C.Transaction__Handle, _arg0 *int, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0, __arg1 := txn.SizeHash()
	*_arg0 = __arg0
	*_arg1 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_coin_Transaction_TxID
func SKY_coin_Transaction_TxID(handle C.Transaction__Handle, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := txn.TxID()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_Transaction_TxIDHex
func SKY_coin_Transaction_TxIDHex(handle C.Transaction__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := txn.TxIDHex()
	copyString(__arg0, _arg0)
	return
}

//export SKY_coin_Transaction_UpdateHeader
func SKY_coin_Transaction_UpdateHeader(handle C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	txn.UpdateHeader()
	return
}

//export SKY_coin_Transaction_HashInner
func SKY_coin_Transaction_HashInner(handle C.Transaction__Handle, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := txn.HashInner()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Transaction_Serialize
func SKY_coin_Transaction_Serialize(handle C.Transaction__Handle, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := txn.Serialize()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_MustTransactionDeserialize
func SKY_coin_MustTransactionDeserialize(_b []byte, _arg1 *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1 := coin.MustTransactionDeserialize(b)
	*_arg1 = registerTransactionHandle(&__arg1)
	return
}

//export SKY_coin_TransactionDeserialize
func SKY_coin_TransactionDeserialize(_b []byte, _arg1 *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1, ____return_err := coin.TransactionDeserialize(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerTransactionHandle(&__arg1)
	}
	return
}

//export SKY_coin_Transaction_OutputHours
func SKY_coin_Transaction_OutputHours(handle C.Transaction__Handle, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn, ok := lookupTransactionHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := txn.OutputHours()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

//export SKY_coin_Create_Transactions
func SKY_coin_Create_Transactions(handle *C.Transactions__Handle) (____error_code uint32) {
	txs := make(coin.Transactions, 0, 0)
	*handle = registerTransactionsHandle(&txs)
	return SKY_OK
}

//export SKY_coin_Get_Transactions_Object
func SKY_coin_Get_Transactions_Object(handle C.Transactions__Handle, _pptx **C.coin__Transactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	ptx, ok := lookupTransactionsHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
	} else {
		*_pptx = (*C.coin__Transactions)(unsafe.Pointer(ptx))
	}
	return
}

//export SKY_coin_Transactions_Length
func SKY_coin_Transactions_Length(handle C.Transactions__Handle, _length *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns, ok := lookupTransactionsHandle(handle)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	*_length = len(*txns)
	return
}

//export SKY_coin_Transactions_Add
func SKY_coin_Transactions_Add(tsh C.Transactions__Handle, th C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns, ok := lookupTransactionsHandle(tsh)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	tx, okt := lookupTransactionHandle(th)
	if !okt {
		____error_code = SKY_ERROR
		return
	}
	*txns = append(*txns, *tx)
	result := overwriteHandle(tsh, txns)
	if !result {
		____error_code = SKY_ERROR
	}
	return
}

//export SKY_coin_Transactions_GetAt
func SKY_coin_Transactions_GetAt(tsh C.Transactions__Handle, n int, th *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns, ok := lookupTransactionsHandle(tsh)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	if n >= len(*txns) {
		____error_code = SKY_ERROR
		return
	}
	tx := (*txns)[n]
	*th = registerTransactionHandle(&tx)
	return
}

//export SKY_coin_Transactions_Hashes
func SKY_coin_Transactions_Hashes(tsh C.Transactions__Handle, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns, ok := lookupTransactionsHandle(tsh)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := txns.Hashes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_Transactions_Size
func SKY_coin_Transactions_Size(tsh C.Transactions__Handle, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns, ok := lookupTransactionsHandle(tsh)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	__arg0 := txns.Size()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Transactions_TruncateBytesTo
func SKY_coin_Transactions_TruncateBytesTo(tsh C.Transactions__Handle, _size int, _arg1 *C.Transactions__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns, ok := lookupTransactionsHandle(tsh)
	if !ok {
		____error_code = SKY_ERROR
		return
	}
	size := _size
	__arg1 := txns.TruncateBytesTo(size)
	*_arg1 = registerTransactionsHandle(&__arg1)
	return
}

//export SKY_coin_SortableTransactions_Sort
func SKY_coin_SortableTransactions_Sort(_txns *C.coin__SortableTransactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*coin.SortableTransactions)(unsafe.Pointer(_txns))
	txns.Sort()
	return
}

//export SKY_coin_SortableTransactions_Len
func SKY_coin_SortableTransactions_Len(_txns *C.coin__SortableTransactions, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*coin.SortableTransactions)(unsafe.Pointer(_txns))
	__arg0 := txns.Len()
	*_arg0 = __arg0
	return
}

//export SKY_coin_SortableTransactions_Less
func SKY_coin_SortableTransactions_Less(_txns *C.coin__SortableTransactions, _i, _j int, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*coin.SortableTransactions)(unsafe.Pointer(_txns))
	i := _i
	j := _j
	__arg1 := txns.Less(i, j)
	*_arg1 = __arg1
	return
}

//export SKY_coin_SortableTransactions_Swap
func SKY_coin_SortableTransactions_Swap(_txns *C.coin__SortableTransactions, _i, _j int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*coin.SortableTransactions)(unsafe.Pointer(_txns))
	i := _i
	j := _j
	txns.Swap(i, j)
	return
}

//export SKY_coin_VerifyTransactionCoinsSpending
func SKY_coin_VerifyTransactionCoinsSpending(_uxIn *C.coin__UxArray, _uxOut *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uxIn := *(*coin.UxArray)(unsafe.Pointer(_uxIn))
	uxOut := *(*coin.UxArray)(unsafe.Pointer(_uxOut))
	____return_err := coin.VerifyTransactionCoinsSpending(uxIn, uxOut)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_coin_VerifyTransactionHoursSpending
func SKY_coin_VerifyTransactionHoursSpending(_headTime uint64, _uxIn *C.coin__UxArray, _uxOut *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	headTime := _headTime
	uxIn := *(*coin.UxArray)(unsafe.Pointer(_uxIn))
	uxOut := *(*coin.UxArray)(unsafe.Pointer(_uxOut))
	____return_err := coin.VerifyTransactionHoursSpending(headTime, uxIn, uxOut)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
