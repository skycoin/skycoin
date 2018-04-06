package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	visor "github.com/skycoin/skycoin/src/visor"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_TxnUnspents_AllForAddress
func SKY_visor_TxnUnspents_AllForAddress(_tus *C.TxnUnspents, _a *C.Address, _arg1 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	tus := *(*cipher.TxnUnspents)(unsafe.Pointer(_tus))
	__arg1 := tus.AllForAddress(a)
	return
}

// export SKY_visor_UnconfirmedTxn_Hash
func SKY_visor_UnconfirmedTxn_Hash(_ut *C.UnconfirmedTxn, _arg0 *C.SHA256) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ut := (*cipher.UnconfirmedTxn)(unsafe.Pointer(_ut))
	__arg0 := ut.Hash()
	return
}

// export SKY_visor_NewUnconfirmedTxnPool
func SKY_visor_NewUnconfirmedTxnPool(_db *C.DB, _arg1 *C.UnconfirmedTxnPool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1 := visor.NewUnconfirmedTxnPool(db)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUnconfirmedTxnPool))
	return
}

// export SKY_visor_UnconfirmedTxnPool_SetAnnounced
func SKY_visor_UnconfirmedTxnPool_SetAnnounced(_utp *C.UnconfirmedTxnPool, _h *C.SHA256, _t *C.Time) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	____return_err := utp.SetAnnounced(h, t)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_InjectTransaction
func SKY_visor_UnconfirmedTxnPool_InjectTransaction(_utp *C.UnconfirmedTxnPool, _bc *C.Blockchainer, _t *C.Transaction, _maxSize int, _arg3 *bool, _arg4 *C.ErrTxnViolatesSoftConstraint) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	bc := *(*cipher.Blockchainer)(unsafe.Pointer(_bc))
	maxSize := _maxSize
	__arg3, __arg4, ____return_err := utp.InjectTransaction(bc, t, maxSize)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg3 = __arg3
		copyToBuffer(reflect.ValueOf((*__arg4)[:]), unsafe.Pointer(_arg4), uint(SizeofErrTxnViolatesSoftConstraint))
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_RawTxns
func SKY_visor_UnconfirmedTxnPool_RawTxns(_utp *C.UnconfirmedTxnPool, _arg0 *C.Transactions) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg0 := utp.RawTxns()
	return
}

// export SKY_visor_UnconfirmedTxnPool_RemoveTransactions
func SKY_visor_UnconfirmedTxnPool_RemoveTransactions(_utp *C.UnconfirmedTxnPool, _txns *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	____return_err := utp.RemoveTransactions(txns)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_RemoveTransactionsWithTx
func SKY_visor_UnconfirmedTxnPool_RemoveTransactionsWithTx(_utp *C.UnconfirmedTxnPool, _tx *C.Tx, _txns *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	utp.RemoveTransactionsWithTx(tx, txns)
	return
}

// export SKY_visor_UnconfirmedTxnPool_Refresh
func SKY_visor_UnconfirmedTxnPool_Refresh(_utp *C.UnconfirmedTxnPool, _bc *C.Blockchainer, _maxBlockSize int, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	bc := *(*cipher.Blockchainer)(unsafe.Pointer(_bc))
	maxBlockSize := _maxBlockSize
	__arg2, ____return_err := utp.Refresh(bc, maxBlockSize)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_RemoveInvalid
func SKY_visor_UnconfirmedTxnPool_RemoveInvalid(_utp *C.UnconfirmedTxnPool, _bc *C.Blockchainer, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	bc := *(*cipher.Blockchainer)(unsafe.Pointer(_bc))
	__arg1, ____return_err := utp.RemoveInvalid(bc)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_FilterKnown
func SKY_visor_UnconfirmedTxnPool_FilterKnown(_utp *C.UnconfirmedTxnPool, _txns *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg1 := utp.FilterKnown(txns)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_visor_UnconfirmedTxnPool_GetKnown
func SKY_visor_UnconfirmedTxnPool_GetKnown(_utp *C.UnconfirmedTxnPool, _txns *C.GoSlice_, _arg1 *C.Transactions) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg1 := utp.GetKnown(txns)
	return
}

// export SKY_visor_UnconfirmedTxnPool_RecvOfAddresses
func SKY_visor_UnconfirmedTxnPool_RecvOfAddresses(_utp *C.UnconfirmedTxnPool, _bh *C.BlockHeader, _addrs *C.GoSlice_, _arg2 *C.AddressUxOuts) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg2, ____return_err := utp.RecvOfAddresses(bh, addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_SpendsOfAddresses
func SKY_visor_UnconfirmedTxnPool_SpendsOfAddresses(_utp *C.UnconfirmedTxnPool, _addrs *C.GoSlice_, _unspent *C.UnspentGetter, _arg2 *C.AddressUxOuts) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg2, ____return_err := utp.SpendsOfAddresses(addrs, unspent)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_GetSpendingOutputs
func SKY_visor_UnconfirmedTxnPool_GetSpendingOutputs(_utp *C.UnconfirmedTxnPool, _bcUnspent *C.UnspentPool, _arg1 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg1, ____return_err := utp.GetSpendingOutputs(bcUnspent)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_GetIncomingOutputs
func SKY_visor_UnconfirmedTxnPool_GetIncomingOutputs(_utp *C.UnconfirmedTxnPool, _bh *C.BlockHeader, _arg1 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg1 := utp.GetIncomingOutputs(bh)
	return
}

// export SKY_visor_UnconfirmedTxnPool_Get
func SKY_visor_UnconfirmedTxnPool_Get(_utp *C.UnconfirmedTxnPool, _key *C.SHA256, _arg1 *C.UnconfirmedTxn, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg1, __arg2 := utp.Get(key)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUnconfirmedTxn))
	*_arg2 = __arg2
	return
}

// export SKY_visor_UnconfirmedTxnPool_GetTxns
func SKY_visor_UnconfirmedTxnPool_GetTxns(_utp *C.UnconfirmedTxnPool, _filter C.Handle, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg1 := utp.GetTxns(filter)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_visor_UnconfirmedTxnPool_GetTxHashes
func SKY_visor_UnconfirmedTxnPool_GetTxHashes(_utp *C.UnconfirmedTxnPool, _filter C.Handle, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg1 := utp.GetTxHashes(filter)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_visor_UnconfirmedTxnPool_ForEach
func SKY_visor_UnconfirmedTxnPool_ForEach(_utp *C.UnconfirmedTxnPool, _f C.Handle) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	____return_err := utp.ForEach(f)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_UnconfirmedTxnPool_GetUnspentsOfAddr
func SKY_visor_UnconfirmedTxnPool_GetUnspentsOfAddr(_utp *C.UnconfirmedTxnPool, _addr *C.Address, _arg1 *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg1 := utp.GetUnspentsOfAddr(addr)
	return
}

// export SKY_visor_IsValid
func SKY_visor_IsValid(_tx *C.UnconfirmedTxn, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	tx := *(*cipher.UnconfirmedTxn)(unsafe.Pointer(_tx))
	__arg1 := visor.IsValid(tx)
	*_arg1 = __arg1
	return
}

// export SKY_visor_All
func SKY_visor_All(_tx *C.UnconfirmedTxn, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	tx := *(*cipher.UnconfirmedTxn)(unsafe.Pointer(_tx))
	__arg1 := visor.All(tx)
	*_arg1 = __arg1
	return
}

// export SKY_visor_UnconfirmedTxnPool_Len
func SKY_visor_UnconfirmedTxnPool_Len(_utp *C.UnconfirmedTxnPool, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	utp := (*cipher.UnconfirmedTxnPool)(unsafe.Pointer(_utp))
	__arg0 := utp.Len()
	*_arg0 = __arg0
	return
}
