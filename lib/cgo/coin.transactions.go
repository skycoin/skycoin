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

//export SKY_coin_Transaction_Verify
func SKY_coin_Transaction_Verify(_txn *C.coin__Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	____return_err := txn.Verify()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_coin_Transaction_VerifyInput
func SKY_coin_Transaction_VerifyInput(_txn *C.coin__Transaction, _uxIn *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := *(*coin.Transaction)(unsafe.Pointer(_txn))
	uxIn := *(*coin.UxArray)(unsafe.Pointer(_uxIn))
	____return_err := txn.VerifyInput(uxIn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_coin_Transaction_PushInput
func SKY_coin_Transaction_PushInput(_txn *C.coin__Transaction, _uxOut *C.cipher__SHA256, _arg1 *uint16) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
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
func SKY_coin_Transaction_PushOutput(_txn *C.coin__Transaction, _dst *C.cipher__Address, _coins, _hours uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	dst := *(*cipher.Address)(unsafe.Pointer(_dst))
	coins := _coins
	hours := _hours
	txn.PushOutput(dst, coins, hours)
	return
}

//export SKY_coin_Transaction_SignInputs
func SKY_coin_Transaction_SignInputs(_txn *C.coin__Transaction, _keys []C.cipher__SecKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	keys := *(*[]cipher.SecKey)(unsafe.Pointer(&_keys))
	txn.SignInputs(keys)
	return
}

//export SKY_coin_Transaction_Size
func SKY_coin_Transaction_Size(_txn *C.coin__Transaction, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.Size()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Transaction_Hash
func SKY_coin_Transaction_Hash(_txn *C.coin__Transaction, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.Hash()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Transaction_SizeHash
func SKY_coin_Transaction_SizeHash(_txn *C.coin__Transaction, _arg0 *int, _arg1 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	__arg0, __arg1 := txn.SizeHash()
	*_arg0 = __arg0
	*_arg1 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_coin_Transaction_TxID
func SKY_coin_Transaction_TxID(_txn *C.coin__Transaction, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.TxID()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_Transaction_TxIDHex
func SKY_coin_Transaction_TxIDHex(_txn *C.coin__Transaction, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.TxIDHex()
	copyString(__arg0, _arg0)
	return
}

//export SKY_coin_Transaction_UpdateHeader
func SKY_coin_Transaction_UpdateHeader(_txn *C.coin__Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	txn.UpdateHeader()
	return
}

//export SKY_coin_Transaction_HashInner
func SKY_coin_Transaction_HashInner(_txn *C.coin__Transaction, _arg0 *C.cipher__SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.HashInner()
	*_arg0 = *(*C.cipher__SHA256)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_coin_Transaction_Serialize
func SKY_coin_Transaction_Serialize(_txn *C.coin__Transaction, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	__arg0 := txn.Serialize()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_MustTransactionDeserialize
func SKY_coin_MustTransactionDeserialize(_b []byte, _arg1 *C.coin__Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1 := coin.MustTransactionDeserialize(b)
	*_arg1 = *(*C.coin__Transaction)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_coin_TransactionDeserialize
func SKY_coin_TransactionDeserialize(_b []byte, _arg1 *C.coin__Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1, ____return_err := coin.TransactionDeserialize(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.coin__Transaction)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_coin_Transaction_OutputHours
func SKY_coin_Transaction_OutputHours(_txn *C.coin__Transaction, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txn := (*coin.Transaction)(unsafe.Pointer(_txn))
	__arg0, ____return_err := txn.OutputHours()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}



//export SKY_coin_Transactions_Hashes
func SKY_coin_Transactions_Hashes(_txns *C.coin__Transactions, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*coin.Transactions)(unsafe.Pointer(_txns))
	__arg0 := txns.Hashes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_coin_Transactions_Size
func SKY_coin_Transactions_Size(_txns *C.coin__Transactions, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*coin.Transactions)(unsafe.Pointer(_txns))
	__arg0 := txns.Size()
	*_arg0 = __arg0
	return
}

//export SKY_coin_Transactions_TruncateBytesTo
func SKY_coin_Transactions_TruncateBytesTo(_txns *C.coin__Transactions, _size int, _arg1 *C.coin__Transactions) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	txns := *(*coin.Transactions)(unsafe.Pointer(_txns))
	size := _size
	__arg1 := txns.TruncateBytesTo(size)
	*_arg1 = *(*C.coin__Transactions)(unsafe.Pointer(&__arg1))
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
