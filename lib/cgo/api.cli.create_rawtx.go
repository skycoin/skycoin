package main

/*
#include <string.h>
#include <stdlib.h>

#include "../../include/skytypes.h"

*/
import "C"

import (
	"unsafe"

	"github.com/skycoin/skycoin/src/api/cli"
	//	"github.com/skycoin/skycoin/src/wallet"
)

/**
 * Functions in github.com/skycoin/skycoin/src/api/cli/transaction.go
 */

//export SKY_cli_CreateRawTxFromWallet
func SKY_cli_CreateRawTxFromWallet(_ctx C.Handle, _walletFile, _chgAddr string, _toAddrs []C.SendAmount, _tx *C.Transaction) uint32 {
	// TODO: Instantiate _ctx . Not used in cli function
	toAddrs := (*[]cli.SendAmount)(unsafe.Pointer(&_toAddrs))
	tx, err := cli.CreateRawTxFromWallet(nil, _walletFile, _chgAddr, *toAddrs)
	*_tx = *(*C.Transaction)(unsafe.Pointer(&tx))
	if err != nil {
		return 1
	}
	return 0
}

//export SKY_cli_CreateRawTxFromAddress
func SKY_cli_CreateRawTxFromAddress(_ctx C.Handle, _addr, _walletFile, _chgAddr string, _toAddrs []C.SendAmount, _tx *C.Transaction) uint32 {
	// TODO: Implement
	return 0
}

//export SKY_cli_CreateRawTx
func SKY_cli_CreateRawTx(_ctx C.Handle, _wlt *C.Wallet, _inAddrs []string, _chgAddr string, _toAddrs []C.SendAmount, _tx *C.Transaction) {
	// TODO: Implement
}

//export SKY_cli_NewTransaction
func SKY_cli_NewTransaction(_utxos []C.UxBalance, _keys []C.SecKey, _outs []C.TransactionOutput, _tx *C.Transaction) {
	//	utxos := (*wallet.UxBalance)(unsafe.Pointer(_utxos))
	// TODO: Implement
}
