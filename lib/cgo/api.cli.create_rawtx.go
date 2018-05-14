package main

/*
#include <string.h>
#include <stdlib.h>

#include "skytypes.h"

*/
import "C"

import (
	"unsafe"

	"github.com/skycoin/skycoin/src/cli"
	//	"github.com/skycoin/skycoin/src/wallet"
)

/**
 * Functions in github.com/skycoin/skycoin/src/api/cli/transaction.go
 */

//export SKY_cli_CreateRawTxFromWallet
func SKY_cli_CreateRawTxFromWallet(_ctx C.Handle, _walletFile, _chgAddr string, _toAddrs []C.cli__SendAmount, _pr C.PasswordReader__Handle, _tx *C.coin__Transaction) uint32 {
	// TODO: Instantiate _ctx . Not used in cli function
	toAddrs := (*[]cli.SendAmount)(unsafe.Pointer(&_toAddrs))
	pr, isOk := lookupPasswordReaderHandle(_pr)
	if !isOk {
		return SKY_ERROR
	}
	tx, err := cli.CreateRawTxFromWallet(nil, _walletFile, _chgAddr, *toAddrs, pr)
	*_tx = *(*C.coin__Transaction)(unsafe.Pointer(&tx))
	if err != nil {
		return SKY_ERROR
	}
	return SKY_OK
}

//export SKY_cli_CreateRawTxFromAddress
func SKY_cli_CreateRawTxFromAddress(_ctx C.WebRpcClient__Handle, _addr, _walletFile, _chgAddr string, _toAddrs []C.cli__SendAmount, _tx *C.coin__Transaction) uint32 {
	// TODO: Implement
	return SKY_ERROR
}

//export SKY_cli_CreateRawTx
func SKY_cli_CreateRawTx(_ctx C.Handle, _wlt *C.wallet__Wallet, _inAddrs []string, _chgAddr string, _toAddrs []C.cli__SendAmount, _tx *C.coin__Transaction) {
	// TODO: Implement
}

//export SKY_cli_NewTransaction
func SKY_cli_NewTransaction(_utxos []C.wallet__UxBalance, _keys []C.cipher__SecKey, _outs []C.coin__TransactionOutput, _tx *C.coin__Transaction) {
	//	utxos := (*wallet.UxBalance)(unsafe.Pointer(_utxos))
	// TODO: Implement
}
