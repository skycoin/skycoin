package main

import (
	daemon "github.com/skycoin/skycoin/src/daemon"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_daemon_NewGatewayConfig
func SKY_daemon_NewGatewayConfig(_arg0 *C.GatewayConfig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := daemon.NewGatewayConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofGatewayConfig))
	return
}

// export SKY_daemon_NewGateway
func SKY_daemon_NewGateway(_c *C.GatewayConfig, _d *C.Daemon, _arg2 *C.Gateway) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*GatewayConfig)(unsafe.Pointer(_c))
	d := (*Daemon)(unsafe.Pointer(_d))
	__arg2 := daemon.NewGateway(c, d)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofGateway))
	return
}

// export SKY_daemon_Gateway_Shutdown
func SKY_daemon_Gateway_Shutdown(_gw *C.Gateway) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	gw.Shutdown()
	return
}

// export SKY_daemon_Gateway_GetConnections
func SKY_daemon_Gateway_GetConnections(_gw *C.Gateway, _arg0 *C.Connections) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetConnections()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofConnections))
	return
}

// export SKY_daemon_Gateway_GetDefaultConnections
func SKY_daemon_Gateway_GetDefaultConnections(_gw *C.Gateway, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetDefaultConnections()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_Gateway_GetConnection
func SKY_daemon_Gateway_GetConnection(_gw *C.Gateway, _addr string, _arg1 *C.Connection) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	addr := _addr
	__arg1 := gw.GetConnection(addr)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofConnection))
	return
}

// export SKY_daemon_Gateway_GetTrustConnections
func SKY_daemon_Gateway_GetTrustConnections(_gw *C.Gateway, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetTrustConnections()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_Gateway_GetExchgConnection
func SKY_daemon_Gateway_GetExchgConnection(_gw *C.Gateway, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetExchgConnection()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_Gateway_GetBlockchainProgress
func SKY_daemon_Gateway_GetBlockchainProgress(_gw *C.Gateway, _arg0 *C.BlockchainProgress) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetBlockchainProgress()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofBlockchainProgress))
	return
}

// export SKY_daemon_Gateway_ResendTransaction
func SKY_daemon_Gateway_ResendTransaction(_gw *C.Gateway, _txn *C.SHA256, _arg1 *C.ResendResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1 := gw.ResendTransaction(txn)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofResendResult))
	return
}

// export SKY_daemon_Gateway_ResendUnconfirmedTxns
func SKY_daemon_Gateway_ResendUnconfirmedTxns(_gw *C.Gateway, _arg0 *C.ResendResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.ResendUnconfirmedTxns()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofResendResult))
	return
}

// export SKY_daemon_Gateway_GetBlockchainMetadata
func SKY_daemon_Gateway_GetBlockchainMetadata(_gw *C.Gateway, _arg0 *C.BlockchainMetadata) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetBlockchainMetadata()
	return
}

// export SKY_daemon_Gateway_GetBlockByHash
func SKY_daemon_Gateway_GetBlockByHash(_gw *C.Gateway, _hash *C.SHA256, _arg1 *C.SignedBlock, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1, __arg2 := gw.GetBlockByHash(hash)
	*_arg2 = __arg2
	return
}

// export SKY_daemon_Gateway_GetBlockBySeq
func SKY_daemon_Gateway_GetBlockBySeq(_gw *C.Gateway, _seq uint64, _arg1 *C.SignedBlock, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	seq := _seq
	__arg1, __arg2 := gw.GetBlockBySeq(seq)
	*_arg2 = __arg2
	return
}

// export SKY_daemon_Gateway_GetBlocks
func SKY_daemon_Gateway_GetBlocks(_gw *C.Gateway, _start, _end uint64, _arg1 *C.ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	start := _start
	end := _end
	__arg1, ____return_err := gw.GetBlocks(start, end)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetBlocksInDepth
func SKY_daemon_Gateway_GetBlocksInDepth(_gw *C.Gateway, _vs *C.GoSlice_, _arg1 *C.ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	vs := *(*[]uint64)(unsafe.Pointer(_vs))
	__arg1, ____return_err := gw.GetBlocksInDepth(vs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetLastBlocks
func SKY_daemon_Gateway_GetLastBlocks(_gw *C.Gateway, _num uint64, _arg1 *C.ReadableBlocks) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	num := _num
	__arg1, ____return_err := gw.GetLastBlocks(num)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetUnspentOutputs
func SKY_daemon_Gateway_GetUnspentOutputs(_gw *C.Gateway, _filters ...*C.OutputsFilter, _arg1 *C.ReadableOutputSet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	filters := _filters
	__arg1, ____return_err := gw.GetUnspentOutputs(filters)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_FbyAddressesNotIncluded
func SKY_daemon_FbyAddressesNotIncluded(_addrs *C.GoSlice_, _arg1 *C.OutputsFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1 := daemon.FbyAddressesNotIncluded(addrs)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofOutputsFilter))
	return
}

// export SKY_daemon_FbyAddresses
func SKY_daemon_FbyAddresses(_addrs *C.GoSlice_, _arg1 *C.OutputsFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1 := daemon.FbyAddresses(addrs)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofOutputsFilter))
	return
}

// export SKY_daemon_FbyHashes
func SKY_daemon_FbyHashes(_hashes *C.GoSlice_, _arg1 *C.OutputsFilter) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hashes := *(*[]string)(unsafe.Pointer(_hashes))
	__arg1 := daemon.FbyHashes(hashes)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofOutputsFilter))
	return
}

// export SKY_daemon_MakeSearchMap
func SKY_daemon_MakeSearchMap(_addrs *C.GoSlice_, _arg1 map[string]struct{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1 := daemon.MakeSearchMap(addrs)
	return
}

// export SKY_daemon_Gateway_GetTransaction
func SKY_daemon_Gateway_GetTransaction(_gw *C.Gateway, _txid *C.SHA256, _arg1 *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1, ____return_err := gw.GetTransaction(txid)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetTransactionResult
func SKY_daemon_Gateway_GetTransactionResult(_gw *C.Gateway, _txid *C.SHA256, _arg1 *C.TransactionResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1, ____return_err := gw.GetTransactionResult(txid)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_InjectBroadcastTransaction
func SKY_daemon_Gateway_InjectBroadcastTransaction(_gw *C.Gateway, _txn *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	____return_err := gw.InjectBroadcastTransaction(txn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetAddressTxns
func SKY_daemon_Gateway_GetAddressTxns(_gw *C.Gateway, _a *C.Address, _arg1 *C.TransactionResults) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1, ____return_err := gw.GetAddressTxns(a)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetTransactions
func SKY_daemon_Gateway_GetTransactions(_gw *C.Gateway, _flts ...*C.TxFilter, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	flts := _flts
	__arg1, ____return_err := gw.GetTransactions(flts)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_daemon_Gateway_GetUxOutByID
func SKY_daemon_Gateway_GetUxOutByID(_gw *C.Gateway, _id *C.SHA256, _arg1 *C.UxOut) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1, ____return_err := gw.GetUxOutByID(id)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetAddrUxOuts
func SKY_daemon_Gateway_GetAddrUxOuts(_gw *C.Gateway, _addr *C.Address, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1, ____return_err := gw.GetAddrUxOuts(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_daemon_Gateway_GetTimeNow
func SKY_daemon_Gateway_GetTimeNow(_gw *C.Gateway, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetTimeNow()
	*_arg0 = __arg0
	return
}

// export SKY_daemon_Gateway_GetAllUnconfirmedTxns
func SKY_daemon_Gateway_GetAllUnconfirmedTxns(_gw *C.Gateway, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetAllUnconfirmedTxns()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_daemon_Gateway_GetUnconfirmedTxns
func SKY_daemon_Gateway_GetUnconfirmedTxns(_gw *C.Gateway, _addrs *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1 := gw.GetUnconfirmedTxns(addrs)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_daemon_Gateway_GetUnspent
func SKY_daemon_Gateway_GetUnspent(_gw *C.Gateway, _arg0 *C.UnspentPool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetUnspent()
	return
}

// export SKY_daemon_spendValidator_HasUnconfirmedSpendTx
func SKY_daemon_spendValidator_HasUnconfirmedSpendTx(_sv spendValidator, _addr *C.GoSlice_, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sv := *(*spendValidator)(unsafe.Pointer(_sv))
	__arg1, ____return_err := sv.HasUnconfirmedSpendTx(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_daemon_Gateway_Spend
func SKY_daemon_Gateway_Spend(_gw *C.Gateway, _wltID string, _password *C.GoSlice_, _coins uint64, _dest *C.Address, _arg4 *C.Transaction) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltID := _wltID
	password := *(*[]byte)(unsafe.Pointer(_password))
	coins := _coins
	__arg4, ____return_err := gw.Spend(wltID, password, coins, dest)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_CreateWallet
func SKY_daemon_Gateway_CreateWallet(_gw *C.Gateway, _wltName string, _options *C.Options, _arg2 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltName := _wltName
	__arg2, ____return_err := gw.CreateWallet(wltName, options)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_ScanAheadWalletAddresses
func SKY_daemon_Gateway_ScanAheadWalletAddresses(_gw *C.Gateway, _wltName string, _password *C.GoSlice_, _scanN uint64, _arg3 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltName := _wltName
	password := *(*[]byte)(unsafe.Pointer(_password))
	scanN := _scanN
	__arg3, ____return_err := gw.ScanAheadWalletAddresses(wltName, password, scanN)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_EncryptWallet
func SKY_daemon_Gateway_EncryptWallet(_gw *C.Gateway, _wltName string, _password *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltName := _wltName
	password := *(*[]byte)(unsafe.Pointer(_password))
	____return_err := gw.EncryptWallet(wltName, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetWalletBalance
func SKY_daemon_Gateway_GetWalletBalance(_gw *C.Gateway, _wltID string, _arg1 *C.BalancePair) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltID := _wltID
	__arg1, ____return_err := gw.GetWalletBalance(wltID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetBalanceOfAddrs
func SKY_daemon_Gateway_GetBalanceOfAddrs(_gw *C.Gateway, _addrs *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg1, ____return_err := gw.GetBalanceOfAddrs(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_daemon_Gateway_GetWalletDir
func SKY_daemon_Gateway_GetWalletDir(_gw *C.Gateway, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0, ____return_err := gw.GetWalletDir()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg0, _arg0)
	}
	return
}

// export SKY_daemon_Gateway_NewAddresses
func SKY_daemon_Gateway_NewAddresses(_gw *C.Gateway, _wltID string, _password *C.GoSlice_, _n uint64, _arg3 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltID := _wltID
	password := *(*[]byte)(unsafe.Pointer(_password))
	n := _n
	__arg3, ____return_err := gw.NewAddresses(wltID, password, n)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg3), _arg3)
	}
	return
}

// export SKY_daemon_Gateway_UpdateWalletLabel
func SKY_daemon_Gateway_UpdateWalletLabel(_gw *C.Gateway, _wltID, _label string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltID := _wltID
	label := _label
	____return_err := gw.UpdateWalletLabel(wltID, label)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetWallet
func SKY_daemon_Gateway_GetWallet(_gw *C.Gateway, _wltID string, _arg1 *C.Wallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltID := _wltID
	__arg1, ____return_err := gw.GetWallet(wltID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetWallets
func SKY_daemon_Gateway_GetWallets(_gw *C.Gateway, _arg0 *C.Wallets) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0, ____return_err := gw.GetWallets()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetWalletUnconfirmedTxns
func SKY_daemon_Gateway_GetWalletUnconfirmedTxns(_gw *C.Gateway, _wltID string, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	wltID := _wltID
	__arg1, ____return_err := gw.GetWalletUnconfirmedTxns(wltID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_daemon_Gateway_ReloadWallets
func SKY_daemon_Gateway_ReloadWallets(_gw *C.Gateway) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	____return_err := gw.ReloadWallets()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_UnloadWallet
func SKY_daemon_Gateway_UnloadWallet(_gw *C.Gateway, _id string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	id := _id
	____return_err := gw.UnloadWallet(id)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_IsWalletAPIDisabled
func SKY_daemon_Gateway_IsWalletAPIDisabled(_gw *C.Gateway, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.IsWalletAPIDisabled()
	*_arg0 = __arg0
	return
}

// export SKY_daemon_Gateway_GetBuildInfo
func SKY_daemon_Gateway_GetBuildInfo(_gw *C.Gateway, _arg0 *C.BuildInfo) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0 := gw.GetBuildInfo()
	return
}

// export SKY_daemon_Gateway_GetRichlist
func SKY_daemon_Gateway_GetRichlist(_gw *C.Gateway, _includeDistribution bool, _arg1 *C.Richlist) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	includeDistribution := _includeDistribution
	__arg1, ____return_err := gw.GetRichlist(includeDistribution)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Gateway_GetAddressCount
func SKY_daemon_Gateway_GetAddressCount(_gw *C.Gateway, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	gw := (*Gateway)(unsafe.Pointer(_gw))
	__arg0, ____return_err := gw.GetAddressCount()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}
