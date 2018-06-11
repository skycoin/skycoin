package main

import (
	"unsafe"

	api "github.com/skycoin/skycoin/src/api"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_api_NewClient
func SKY_api_NewClient(_addr string, _arg1 *C.Client__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := _addr
	__arg1 := api.NewClient(addr)
	*_arg1 = registerClientHandle(__arg1)
	return
}

//export SKY_api_Client_CSRF
func SKY_api_Client_CSRF(_c C.Client__Handle, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.CSRF()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg0, _arg0)
	}
	return
}

//export SKY_api_Client_Version
func SKY_api_Client_Version(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.Version()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_Outputs
func SKY_api_Client_Outputs(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.Outputs()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_OutputsForAddresses
func SKY_api_Client_OutputsForAddresses(_c C.Client__Handle, _addrs []string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1, ____return_err := c.OutputsForAddresses(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_OutputsForHashes
func SKY_api_Client_OutputsForHashes(_c C.Client__Handle, _hashes []string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	hashes := *(*[]string)(unsafe.Pointer(&_hashes))
	__arg1, ____return_err := c.OutputsForHashes(hashes)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_CoinSupply
func SKY_api_Client_CoinSupply(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.CoinSupply()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_BlockByHash
func SKY_api_Client_BlockByHash(_c C.Client__Handle, _hash string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	hash := _hash
	__arg1, ____return_err := c.BlockByHash(hash)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_BlockBySeq
func SKY_api_Client_BlockBySeq(_c C.Client__Handle, _seq uint64, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	seq := _seq
	__arg1, ____return_err := c.BlockBySeq(seq)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_Blocks
func SKY_api_Client_Blocks(_c C.Client__Handle, _start, _end int, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	start := _start
	end := _end
	__arg1, ____return_err := c.Blocks(start, end)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_LastBlocks
func SKY_api_Client_LastBlocks(_c C.Client__Handle, _n int, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	n := _n
	__arg1, ____return_err := c.LastBlocks(n)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_BlockchainMetadata
func SKY_api_Client_BlockchainMetadata(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.BlockchainMetadata()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_BlockchainProgress
func SKY_api_Client_BlockchainProgress(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.BlockchainProgress()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_Balance
func SKY_api_Client_Balance(_c C.Client__Handle, _addrs []string, _arg1 *C.wallet__BalancePair) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1, ____return_err := c.Balance(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.wallet__BalancePair)(unsafe.Pointer(__arg1))
	}
	return
}

//export SKY_api_Client_UxOut
func SKY_api_Client_UxOut(_c C.Client__Handle, _uxID string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	uxID := _uxID
	__arg1, ____return_err := c.UxOut(uxID)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_AddressUxOuts
func SKY_api_Client_AddressUxOuts(_c C.Client__Handle, _addr string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addr := _addr
	__arg1, ____return_err := c.AddressUxOuts(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_Wallet
func SKY_api_Client_Wallet(_c C.Client__Handle, _id string, _arg1 *C.WalletResponse__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	__arg1, ____return_err := c.Wallet(id)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerWalletResponseHandle(__arg1)
	}
	return
}

//export SKY_api_Client_Wallets
func SKY_api_Client_Wallets(_c C.Client__Handle, _arg0 *C.Wallets__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.Wallets()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerWalletsHandle(__arg0)
	}
	return
}

//export SKY_api_Client_CreateUnencryptedWallet
func SKY_api_Client_CreateUnencryptedWallet(_c C.Client__Handle, _seed, _label string, _scanN int, _arg2 *C.WalletResponse__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	seed := _seed
	label := _label
	scanN := _scanN
	__arg2, ____return_err := c.CreateUnencryptedWallet(seed, label, scanN)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerWalletResponseHandle(__arg2)
	}
	return
}

//export SKY_api_Client_CreateEncryptedWallet
func SKY_api_Client_CreateEncryptedWallet(_c C.Client__Handle, _seed, _label, _password string, _scanN int, _arg2 *C.WalletResponse__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	seed := _seed
	label := _label
	password := _password
	scanN := _scanN
	__arg2, ____return_err := c.CreateEncryptedWallet(seed, label, password, scanN)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerWalletResponseHandle(__arg2)
	}
	return
}

//export SKY_api_Client_NewWalletAddress
func SKY_api_Client_NewWalletAddress(_c C.Client__Handle, _id string, _n int, _password string, _arg3 *C.Strings__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	n := _n
	password := _password
	__arg3, ____return_err := c.NewWalletAddress(id, n, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg3 = (C.Strings__Handle)(registerHandle(__arg3))
	}
	return
}

//export SKY_api_Client_WalletBalance
func SKY_api_Client_WalletBalance(_c C.Client__Handle, _id string, _arg1 *C.wallet__BalancePair) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	__arg1, ____return_err := c.WalletBalance(id)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.wallet__BalancePair)(unsafe.Pointer(__arg1))
	}
	return
}

//export SKY_api_Client_Spend
func SKY_api_Client_Spend(_c C.Client__Handle, _id, _dst string, _coins uint64, _password string, _arg3 *C.api__SpendResult) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	dst := _dst
	coins := _coins
	password := _password
	__arg3, ____return_err := c.Spend(id, dst, coins, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg3 = *(*C.api__SpendResult)(unsafe.Pointer(__arg3))
	}
	return
}

//export SKY_api_Client_CreateTransaction
func SKY_api_Client_CreateTransaction(_c C.Client__Handle, _req *C.Handle, _arg1 *C.api__CreateTransactionResponse) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	req, okreq := lookupCreateTransactionRequestHandle(C.CreateTransactionRequest__Handle(*_req))
	if !okreq {
		____error_code = SKY_ERROR
		return
	}
	__arg1, ____return_err := c.CreateTransaction(*req)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.api__CreateTransactionResponse)(unsafe.Pointer(__arg1))
	}
	return
}

//export SKY_api_Client_WalletTransactions
func SKY_api_Client_WalletTransactions(_c C.Client__Handle, _id string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	__arg1, ____return_err := c.WalletTransactions(id)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_UpdateWallet
func SKY_api_Client_UpdateWallet(_c C.Client__Handle, _id, _label string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	label := _label
	____return_err := c.UpdateWallet(id, label)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_api_Client_WalletFolderName
func SKY_api_Client_WalletFolderName(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.WalletFolderName()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_NewSeed
func SKY_api_Client_NewSeed(_c C.Client__Handle, _entropy int, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	entropy := _entropy
	__arg1, ____return_err := c.NewSeed(entropy)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

//export SKY_api_Client_GetWalletSeed
func SKY_api_Client_GetWalletSeed(_c C.Client__Handle, _id string, _password string, _arg2 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	password := _password
	__arg2, ____return_err := c.GetWalletSeed(id, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg2, _arg2)
	}
	return
}

//export SKY_api_Client_NetworkConnection
func SKY_api_Client_NetworkConnection(_c C.Client__Handle, _addr string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addr := _addr
	__arg1, ____return_err := c.NetworkConnection(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_NetworkConnections
func SKY_api_Client_NetworkConnections(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.NetworkConnections()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_NetworkDefaultConnections
func SKY_api_Client_NetworkDefaultConnections(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.NetworkDefaultConnections()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_NetworkTrustedConnections
func SKY_api_Client_NetworkTrustedConnections(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.NetworkTrustedConnections()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_NetworkExchangeableConnections
func SKY_api_Client_NetworkExchangeableConnections(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.NetworkExchangeableConnections()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_PendingTransactions
func SKY_api_Client_PendingTransactions(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.PendingTransactions()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_Transaction
func SKY_api_Client_Transaction(_c C.Client__Handle, _txid string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	txid := _txid
	__arg1, ____return_err := c.Transaction(txid)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_Transactions
func SKY_api_Client_Transactions(_c C.Client__Handle, _addrs []string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1, ____return_err := c.Transactions(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_ConfirmedTransactions
func SKY_api_Client_ConfirmedTransactions(_c C.Client__Handle, _addrs []string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1, ____return_err := c.ConfirmedTransactions(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_UnconfirmedTransactions
func SKY_api_Client_UnconfirmedTransactions(_c C.Client__Handle, _addrs []string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addrs := *(*[]string)(unsafe.Pointer(&_addrs))
	__arg1, ____return_err := c.UnconfirmedTransactions(addrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_InjectTransaction
func SKY_api_Client_InjectTransaction(_c C.Client__Handle, _rawTx string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	rawTx := _rawTx
	__arg1, ____return_err := c.InjectTransaction(rawTx)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

//export SKY_api_Client_ResendUnconfirmedTransactions
func SKY_api_Client_ResendUnconfirmedTransactions(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.ResendUnconfirmedTransactions()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_RawTransaction
func SKY_api_Client_RawTransaction(_c C.Client__Handle, _txid string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	txid := _txid
	__arg1, ____return_err := c.RawTransaction(txid)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

//export SKY_api_Client_AddressTransactions
func SKY_api_Client_AddressTransactions(_c C.Client__Handle, _addr string, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addr := _addr
	__arg1, ____return_err := c.AddressTransactions(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_Richlist
func SKY_api_Client_Richlist(_c C.Client__Handle, _params *C.api__RichlistParams, _arg1 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	params := (*api.RichlistParams)(unsafe.Pointer(_params))
	__arg1, ____return_err := c.Richlist(params)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerHandle(__arg1)
	}
	return
}

//export SKY_api_Client_AddressCount
func SKY_api_Client_AddressCount(_c C.Client__Handle, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.AddressCount()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

//export SKY_api_Client_UnloadWallet
func SKY_api_Client_UnloadWallet(_c C.Client__Handle, _id string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	____return_err := c.UnloadWallet(id)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_api_Client_Health
func SKY_api_Client_Health(_c C.Client__Handle, _arg0 *C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	__arg0, ____return_err := c.Health()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = registerHandle(__arg0)
	}
	return
}

//export SKY_api_Client_EncryptWallet
func SKY_api_Client_EncryptWallet(_c C.Client__Handle, _id string, _password string, _arg2 *C.WalletResponse__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	password := _password
	__arg2, ____return_err := c.EncryptWallet(id, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerWalletResponseHandle(__arg2)
	}
	return
}

//export SKY_api_Client_DecryptWallet
func SKY_api_Client_DecryptWallet(_c C.Client__Handle, _id string, _password string, _arg2 *C.WalletResponse__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	id := _id
	password := _password
	__arg2, ____return_err := c.DecryptWallet(id, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerWalletResponseHandle(__arg2)
	}
	return
}
