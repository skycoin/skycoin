package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	gui "github.com/skycoin/skycoin/src/gui"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gui_APIError_Error
func SKY_gui_APIError_Error(_e *C.APIError, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	e := *(*cipher.APIError)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_gui_NewClient
func SKY_gui_NewClient(_addr string, _arg1 *C.Client) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	addr := _addr
	__arg1 := gui.NewClient(addr)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofClient))
	return
}

// export SKY_gui_Client_Get
func SKY_gui_Client_Get(_c *C.Client, _endpoint string, _obj interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	endpoint := _endpoint
	____return_err := c.Get(endpoint, obj)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_PostForm
func SKY_gui_Client_PostForm(_c *C.Client, _endpoints string, _body *C.Reader, _obj interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	endpoints := _endpoints
	____return_err := c.PostForm(endpoints, body, obj)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_PostJSON
func SKY_gui_Client_PostJSON(_c *C.Client, _endpoints string, _body *C.Reader, _obj interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	endpoints := _endpoints
	____return_err := c.PostJSON(endpoints, body, obj)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_CSRF
func SKY_gui_Client_CSRF(_c *C.Client, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.CSRF()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg0, _arg0)
	}
	return
}

// export SKY_gui_Client_Version
func SKY_gui_Client_Version(_c *C.Client, _arg0 *C.BuildInfo) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.Version()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_Outputs
func SKY_gui_Client_Outputs(_c *C.Client, _arg0 *C.ReadableOutputSet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.Outputs()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_OutputsForAddresses
func SKY_gui_Client_OutputsForAddresses(_c *C.Client, _addrs *C.GoSlice_, _arg1 *C.ReadableOutputSet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1, ____return_err := c.OutputsForAddresses(addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_OutputsForHashes
func SKY_gui_Client_OutputsForHashes(_c *C.Client, _hashes *C.GoSlice_, _arg1 *C.ReadableOutputSet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	hashes := *(*[]string)(unsafe.Pointer(_hashes))
	__arg1, ____return_err := c.OutputsForHashes(hashes)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_CoinSupply
func SKY_gui_Client_CoinSupply(_c *C.Client, _arg0 *C.CoinSupply) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.CoinSupply()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofCoinSupply))
	}
	return
}

// export SKY_gui_Client_BlockByHash
func SKY_gui_Client_BlockByHash(_c *C.Client, _hash string, _arg1 *C.ReadableBlock) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	hash := _hash
	__arg1, ____return_err := c.BlockByHash(hash)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_BlockBySeq
func SKY_gui_Client_BlockBySeq(_c *C.Client, _seq uint64, _arg1 *C.ReadableBlock) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	seq := _seq
	__arg1, ____return_err := c.BlockBySeq(seq)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_Blocks
func SKY_gui_Client_Blocks(_c *C.Client, _start, _end int, _arg1 *C.ReadableBlocks) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	start := _start
	end := _end
	__arg1, ____return_err := c.Blocks(start, end)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_LastBlocks
func SKY_gui_Client_LastBlocks(_c *C.Client, _n int, _arg1 *C.ReadableBlocks) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	n := _n
	__arg1, ____return_err := c.LastBlocks(n)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_BlockchainMetadata
func SKY_gui_Client_BlockchainMetadata(_c *C.Client, _arg0 *C.BlockchainMetadata) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.BlockchainMetadata()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_BlockchainProgress
func SKY_gui_Client_BlockchainProgress(_c *C.Client, _arg0 *C.BlockchainProgress) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.BlockchainProgress()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_Balance
func SKY_gui_Client_Balance(_c *C.Client, _addrs *C.GoSlice_, _arg1 *C.BalancePair) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1, ____return_err := c.Balance(addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_UxOut
func SKY_gui_Client_UxOut(_c *C.Client, _uxID string, _arg1 *C.UxOutJSON) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	uxID := _uxID
	__arg1, ____return_err := c.UxOut(uxID)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_AddressUxOuts
func SKY_gui_Client_AddressUxOuts(_c *C.Client, _addr string, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	addr := _addr
	__arg1, ____return_err := c.AddressUxOuts(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_gui_Client_Wallet
func SKY_gui_Client_Wallet(_c *C.Client, _id string, _arg1 *C.Wallet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	id := _id
	__arg1, ____return_err := c.Wallet(id)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_Wallets
func SKY_gui_Client_Wallets(_c *C.Client, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.Wallets()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_gui_Client_CreateWallet
func SKY_gui_Client_CreateWallet(_c *C.Client, _seed, _label string, _scanN int, _arg2 *C.ReadableWallet) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	seed := _seed
	label := _label
	scanN := _scanN
	__arg2, ____return_err := c.CreateWallet(seed, label, scanN)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_NewWalletAddress
func SKY_gui_Client_NewWalletAddress(_c *C.Client, _id string, _n int, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	id := _id
	n := _n
	__arg2, ____return_err := c.NewWalletAddress(id, n)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	}
	return
}

// export SKY_gui_Client_WalletBalance
func SKY_gui_Client_WalletBalance(_c *C.Client, _id string, _arg1 *C.BalancePair) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	id := _id
	__arg1, ____return_err := c.WalletBalance(id)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_Spend
func SKY_gui_Client_Spend(_c *C.Client, _id, _dst string, _coins uint64, _arg2 *C.SpendResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	id := _id
	dst := _dst
	coins := _coins
	__arg2, ____return_err := c.Spend(id, dst, coins)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofSpendResult))
	}
	return
}

// export SKY_gui_Client_WalletTransactions
func SKY_gui_Client_WalletTransactions(_c *C.Client, _id string, _arg1 *C.UnconfirmedTxnsResponse) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	id := _id
	__arg1, ____return_err := c.WalletTransactions(id)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofUnconfirmedTxnsResponse))
	}
	return
}

// export SKY_gui_Client_UpdateWallet
func SKY_gui_Client_UpdateWallet(_c *C.Client, _id, _label string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	id := _id
	label := _label
	____return_err := c.UpdateWallet(id, label)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_WalletFolderName
func SKY_gui_Client_WalletFolderName(_c *C.Client, _arg0 *C.WalletFolder) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.WalletFolderName()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofWalletFolder))
	}
	return
}

// export SKY_gui_Client_NewSeed
func SKY_gui_Client_NewSeed(_c *C.Client, _entropy int, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	entropy := _entropy
	__arg1, ____return_err := c.NewSeed(entropy)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_gui_Client_NetworkConnection
func SKY_gui_Client_NetworkConnection(_c *C.Client, _addr string, _arg1 *C.Connection) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	addr := _addr
	__arg1, ____return_err := c.NetworkConnection(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_NetworkConnections
func SKY_gui_Client_NetworkConnections(_c *C.Client, _arg0 *C.Connections) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.NetworkConnections()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_NetworkDefaultConnections
func SKY_gui_Client_NetworkDefaultConnections(_c *C.Client, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.NetworkDefaultConnections()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_gui_Client_NetworkTrustedConnections
func SKY_gui_Client_NetworkTrustedConnections(_c *C.Client, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.NetworkTrustedConnections()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_gui_Client_NetworkExchangeableConnections
func SKY_gui_Client_NetworkExchangeableConnections(_c *C.Client, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.NetworkExchangeableConnections()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_gui_Client_PendingTransactions
func SKY_gui_Client_PendingTransactions(_c *C.Client, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.PendingTransactions()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_gui_Client_Transaction
func SKY_gui_Client_Transaction(_c *C.Client, _txid string, _arg1 *C.TransactionResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	txid := _txid
	__arg1, ____return_err := c.Transaction(txid)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_Transactions
func SKY_gui_Client_Transactions(_c *C.Client, _addrs *C.GoSlice_, _arg1 *[]C.TransactionResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1, ____return_err := c.Transactions(addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_gui_Client_ConfirmedTransactions
func SKY_gui_Client_ConfirmedTransactions(_c *C.Client, _addrs *C.GoSlice_, _arg1 *[]C.TransactionResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1, ____return_err := c.ConfirmedTransactions(addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_gui_Client_UnconfirmedTransactions
func SKY_gui_Client_UnconfirmedTransactions(_c *C.Client, _addrs *C.GoSlice_, _arg1 *[]C.TransactionResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	addrs := *(*[]string)(unsafe.Pointer(_addrs))
	__arg1, ____return_err := c.UnconfirmedTransactions(addrs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_gui_Client_InjectTransaction
func SKY_gui_Client_InjectTransaction(_c *C.Client, _rawTx string, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	rawTx := _rawTx
	__arg1, ____return_err := c.InjectTransaction(rawTx)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_gui_Client_ResendUnconfirmedTransactions
func SKY_gui_Client_ResendUnconfirmedTransactions(_c *C.Client, _arg0 *C.ResendResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.ResendUnconfirmedTransactions()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Client_RawTransaction
func SKY_gui_Client_RawTransaction(_c *C.Client, _txid string, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	txid := _txid
	__arg1, ____return_err := c.RawTransaction(txid)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_gui_Client_AddressTransactions
func SKY_gui_Client_AddressTransactions(_c *C.Client, _addr string, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	addr := _addr
	__arg1, ____return_err := c.AddressTransactions(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_gui_Client_Richlist
func SKY_gui_Client_Richlist(_c *C.Client, _params *C.RichlistParams, _arg1 *C.Richlist) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	params := (*cipher.RichlistParams)(unsafe.Pointer(_params))
	__arg1, ____return_err := c.Richlist(params)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofRichlist))
	}
	return
}

// export SKY_gui_Client_AddressCount
func SKY_gui_Client_AddressCount(_c *C.Client, _arg0 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.AddressCount()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

// export SKY_gui_Client_UnloadWallet
func SKY_gui_Client_UnloadWallet(_c *C.Client, _id string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := (*cipher.Client)(unsafe.Pointer(_c))
	id := _id
	____return_err := c.UnloadWallet(id)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
