package main

import (
	coin "github.com/skycoin/skycoin/src/coin"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_coin_Transaction_Verify
func SKY_coin_Transaction_Verify(_txn *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	____return_err := txn.Verify()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_coin_Transaction_VerifyInput
func SKY_coin_Transaction_VerifyInput(_txn *C.Transaction, _uxIn *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := *(*Transaction)(unsafe.Pointer(_txn))
	uxIn := *(*UxArray)(unsafe.Pointer(_uxIn))
	____return_err := txn.VerifyInput(uxIn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_coin_Transaction_PushInput
func SKY_coin_Transaction_PushInput(_txn *C.Transaction, _uxOut *C.SHA256, _arg1 *uint16) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg1 := txn.PushInput(uxOut)
	*_arg1 = __arg1
	return
}

// export SKY_coin_TransactionOutput_UxID
func SKY_coin_TransactionOutput_UxID(_txOut *C.TransactionOutput, _TxID *C.SHA256, _arg1 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txOut := *(*TransactionOutput)(unsafe.Pointer(_txOut))
	__arg1 := txOut.UxID(TxID)
	return
}

// export SKY_coin_Transaction_PushOutput
func SKY_coin_Transaction_PushOutput(_txn *C.Transaction, _dst *C.Address, _coins, _hours uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	coins := _coins
	hours := _hours
	txn.PushOutput(dst, coins, hours)
	return
}

// export SKY_coin_Transaction_SignInputs
func SKY_coin_Transaction_SignInputs(_txn *C.Transaction, _keys *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	txn.SignInputs(keys)
	return
}

// export SKY_coin_Transaction_Size
func SKY_coin_Transaction_Size(_txn *C.Transaction, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.Size()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Transaction_Hash
func SKY_coin_Transaction_Hash(_txn *C.Transaction, _arg0 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.Hash()
	return
}

// export SKY_coin_Transaction_SizeHash
func SKY_coin_Transaction_SizeHash(_txn *C.Transaction, _arg0 *int, _arg1 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg0, __arg1 := txn.SizeHash()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Transaction_TxID
func SKY_coin_Transaction_TxID(_txn *C.Transaction, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.TxID()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_Transaction_TxIDHex
func SKY_coin_Transaction_TxIDHex(_txn *C.Transaction, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.TxIDHex()
	copyString(__arg0, _arg0)
	return
}

// export SKY_coin_Transaction_UpdateHeader
func SKY_coin_Transaction_UpdateHeader(_txn *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	txn.UpdateHeader()
	return
}

// export SKY_coin_Transaction_HashInner
func SKY_coin_Transaction_HashInner(_txn *C.Transaction, _arg0 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.HashInner()
	return
}

// export SKY_coin_Transaction_Serialize
func SKY_coin_Transaction_Serialize(_txn *C.Transaction, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.Serialize()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_MustTransactionDeserialize
func SKY_coin_MustTransactionDeserialize(_b *C.GoSlice_, _arg1 *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(_b))
	__arg1 := coin.MustTransactionDeserialize(b)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	return
}

// export SKY_coin_TransactionDeserialize
func SKY_coin_TransactionDeserialize(_b *C.GoSlice_, _arg1 *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(_b))
	__arg1, ____return_err := coin.TransactionDeserialize(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTransaction))
	}
	return
}

// export SKY_coin_Transaction_OutputHours
func SKY_coin_Transaction_OutputHours(_txn *C.Transaction, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*Transaction)(unsafe.Pointer(_txn))
	__arg0, ____return_err := txn.OutputHours()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

// export SKY_coin_Transactions_Fees
func SKY_coin_Transactions_Fees(_txns *C.Transactions, _calc *C.FeeCalculator, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*Transactions)(unsafe.Pointer(_txns))
	calc := *(*FeeCalculator)(unsafe.Pointer(_calc))
	__arg1, ____return_err := txns.Fees(calc)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_coin_Transactions_Hashes
func SKY_coin_Transactions_Hashes(_txns *C.Transactions, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*Transactions)(unsafe.Pointer(_txns))
	__arg0 := txns.Hashes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_coin_Transactions_Size
func SKY_coin_Transactions_Size(_txns *C.Transactions, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*Transactions)(unsafe.Pointer(_txns))
	__arg0 := txns.Size()
	*_arg0 = __arg0
	return
}

// export SKY_coin_Transactions_TruncateBytesTo
func SKY_coin_Transactions_TruncateBytesTo(_txns *C.Transactions, _size int, _arg1 *C.Transactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*Transactions)(unsafe.Pointer(_txns))
	size := _size
	__arg1 := txns.TruncateBytesTo(size)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofTransactions))
	return
}

// export SKY_coin_SortTransactions
func SKY_coin_SortTransactions(_txns *C.Transactions, _feeCalc *C.FeeCalculator, _arg2 *C.Transactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*Transactions)(unsafe.Pointer(_txns))
	feeCalc := *(*FeeCalculator)(unsafe.Pointer(_feeCalc))
	__arg2 := coin.SortTransactions(txns, feeCalc)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofTransactions))
	return
}

// export SKY_coin_NewSortableTransactions
func SKY_coin_NewSortableTransactions(_txns *C.Transactions, _feeCalc *C.FeeCalculator, _arg2 *C.SortableTransactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*Transactions)(unsafe.Pointer(_txns))
	feeCalc := *(*FeeCalculator)(unsafe.Pointer(_feeCalc))
	__arg2 := coin.NewSortableTransactions(txns, feeCalc)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofSortableTransactions))
	return
}

// export SKY_coin_SortableTransactions_Sort
func SKY_coin_SortableTransactions_Sort(_txns *C.SortableTransactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*SortableTransactions)(unsafe.Pointer(_txns))
	txns.Sort()
	return
}

// export SKY_coin_SortableTransactions_Len
func SKY_coin_SortableTransactions_Len(_txns *C.SortableTransactions, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*SortableTransactions)(unsafe.Pointer(_txns))
	__arg0 := txns.Len()
	*_arg0 = __arg0
	return
}

// export SKY_coin_SortableTransactions_Less
func SKY_coin_SortableTransactions_Less(_txns *C.SortableTransactions, _i, _j int, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*SortableTransactions)(unsafe.Pointer(_txns))
	i := _i
	j := _j
	__arg1 := txns.Less(i, j)
	*_arg1 = __arg1
	return
}

// export SKY_coin_SortableTransactions_Swap
func SKY_coin_SortableTransactions_Swap(_txns *C.SortableTransactions, _i, _j int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*SortableTransactions)(unsafe.Pointer(_txns))
	i := _i
	j := _j
	txns.Swap(i, j)
	return
}

// export SKY_coin_VerifyTransactionCoinsSpending
func SKY_coin_VerifyTransactionCoinsSpending(_uxIn *C.UxArray, _uxOut *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	uxIn := *(*UxArray)(unsafe.Pointer(_uxIn))
	uxOut := *(*UxArray)(unsafe.Pointer(_uxOut))
	____return_err := coin.VerifyTransactionCoinsSpending(uxIn, uxOut)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_coin_VerifyTransactionHoursSpending
func SKY_coin_VerifyTransactionHoursSpending(_headTime uint64, _uxIn *C.UxArray, _uxOut *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	headTime := _headTime
	uxIn := *(*UxArray)(unsafe.Pointer(_uxIn))
	uxOut := *(*UxArray)(unsafe.Pointer(_uxOut))
	____return_err := coin.VerifyTransactionHoursSpending(headTime, uxIn, uxOut)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_coin_AddUint64
func SKY_coin_AddUint64(_a, _b uint64, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	a := _a
	b := _b
	__arg1, ____return_err := coin.AddUint64(a, b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}
