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

// export SKY_visor_NewTransactionResult
func SKY_visor_NewTransactionResult(_tx *C.Transaction, _arg1 *C.TransactionResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	tx := (*cipher.Transaction)(unsafe.Pointer(_tx))
	__arg1, ____return_err := visor.NewTransactionResult(tx)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransactionResult))
	}
	return
}

// export SKY_visor_NewTransactionResults
func SKY_visor_NewTransactionResults(_txs *C.GoSlice_, _arg1 *C.TransactionResults) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	txs := *(*[]cipher.Transaction)(unsafe.Pointer(_txs))
	__arg1, ____return_err := visor.NewTransactionResults(txs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofTransactionResults))
	}
	return
}

// export SKY_visor_MakeRPC
func SKY_visor_MakeRPC(_v *C.Visor, _arg1 *C.RPC) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg1 := visor.MakeRPC(v)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofRPC))
	return
}

// export SKY_visor_RPC_GetBlockchainMetadata
func SKY_visor_RPC_GetBlockchainMetadata(_rpc *C.RPC, _v *C.Visor, _arg1 *C.BlockchainMetadata) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg1 := rpc.GetBlockchainMetadata(v)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofBlockchainMetadata))
	return
}

// export SKY_visor_RPC_GetUnspent
func SKY_visor_RPC_GetUnspent(_rpc *C.RPC, _v *C.Visor, _arg1 *C.UnspentPool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg1 := rpc.GetUnspent(v)
	return
}

// export SKY_visor_RPC_GetUnconfirmedSpends
func SKY_visor_RPC_GetUnconfirmedSpends(_rpc *C.RPC, _v *C.Visor, _addrs *C.GoSlice_, _arg2 *C.AddressUxOuts) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg2, ____return_err := rpc.GetUnconfirmedSpends(v, addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_GetUnconfirmedReceiving
func SKY_visor_RPC_GetUnconfirmedReceiving(_rpc *C.RPC, _v *C.Visor, _addrs *C.GoSlice_, _arg2 *C.AddressUxOuts) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg2, ____return_err := rpc.GetUnconfirmedReceiving(v, addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_GetUnconfirmedTxns
func SKY_visor_RPC_GetUnconfirmedTxns(_rpc *C.RPC, _v *C.Visor, _addresses *C.GoSlice_, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg2 := rpc.GetUnconfirmedTxns(v, addresses)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

// export SKY_visor_RPC_GetBlock
func SKY_visor_RPC_GetBlock(_rpc *C.RPC, _v *C.Visor, _seq uint64, _arg2 *C.SignedBlock) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	seq := _seq
	__arg2, ____return_err := rpc.GetBlock(v, seq)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_GetBlocks
func SKY_visor_RPC_GetBlocks(_rpc *C.RPC, _v *C.Visor, _start, _end uint64, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	start := _start
	end := _end
	__arg2 := rpc.GetBlocks(v, start, end)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

// export SKY_visor_RPC_GetLastBlocks
func SKY_visor_RPC_GetLastBlocks(_rpc *C.RPC, _v *C.Visor, _num uint64, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	num := _num
	__arg2 := rpc.GetLastBlocks(v, num)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

// export SKY_visor_RPC_GetBlockBySeq
func SKY_visor_RPC_GetBlockBySeq(_rpc *C.RPC, _v *C.Visor, _n uint64, _arg2 *C.SignedBlock) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	n := _n
	__arg2, ____return_err := rpc.GetBlockBySeq(v, n)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_GetTransaction
func SKY_visor_RPC_GetTransaction(_rpc *C.RPC, _v *C.Visor, _txHash *C.SHA256, _arg2 *C.Transaction) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg2, ____return_err := rpc.GetTransaction(v, txHash)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofTransaction))
	}
	return
}

// export SKY_visor_RPC_GetAddressTxns
func SKY_visor_RPC_GetAddressTxns(_rpc *C.RPC, _v *C.Visor, _addr *C.Address, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg2, ____return_err := rpc.GetAddressTxns(v, addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

// export SKY_visor_RPC_CreateWallet
func SKY_visor_RPC_CreateWallet(_rpc *C.RPC, _wltName string, _options *C.Options, _arg2 *C.Wallet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	wltName := _wltName
	__arg2, ____return_err := rpc.CreateWallet(wltName, options)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_NewAddresses
func SKY_visor_RPC_NewAddresses(_rpc *C.RPC, _wltName string, _num uint64, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	wltName := _wltName
	num := _num
	__arg2, ____return_err := rpc.NewAddresses(wltName, num)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

// export SKY_visor_RPC_GetWalletAddresses
func SKY_visor_RPC_GetWalletAddresses(_rpc *C.RPC, _wltID string, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	wltID := _wltID
	__arg1, ____return_err := rpc.GetWalletAddresses(wltID)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_visor_RPC_CreateAndSignTransaction
func SKY_visor_RPC_CreateAndSignTransaction(_rpc *C.RPC, _wltID string, _password *C.GoSlice_, _vld *C.Validator, _unspent *C.UnspentGetter, _headTime, _coins uint64, _dest *C.Address, _arg6 *C.Transaction) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	wltID := _wltID
	password := *(*[]byte)(unsafe.Pointer(_password))
	headTime := _headTime
	coins := _coins
	__arg6, ____return_err := rpc.CreateAndSignTransaction(wltID, password, vld, unspent, headTime, coins, dest)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_UpdateWalletLabel
func SKY_visor_RPC_UpdateWalletLabel(_rpc *C.RPC, _wltID, _label string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	wltID := _wltID
	label := _label
	____return_err := rpc.UpdateWalletLabel(wltID, label)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_GetWallet
func SKY_visor_RPC_GetWallet(_rpc *C.RPC, _wltID string, _arg1 *C.Wallet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	wltID := _wltID
	__arg1, ____return_err := rpc.GetWallet(wltID)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_GetWallets
func SKY_visor_RPC_GetWallets(_rpc *C.RPC, _arg0 *C.Wallets) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	__arg0 := rpc.GetWallets()
	return
}

// export SKY_visor_RPC_ReloadWallets
func SKY_visor_RPC_ReloadWallets(_rpc *C.RPC) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	____return_err := rpc.ReloadWallets()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_RPC_GetBuildInfo
func SKY_visor_RPC_GetBuildInfo(_rpc *C.RPC, _arg0 *C.BuildInfo) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	__arg0 := rpc.GetBuildInfo()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofBuildInfo))
	return
}

// export SKY_visor_RPC_UnloadWallet
func SKY_visor_RPC_UnloadWallet(_rpc *C.RPC, _id string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := (*cipher.RPC)(unsafe.Pointer(_rpc))
	id := _id
	rpc.UnloadWallet(id)
	return
}
