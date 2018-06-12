package main

import (
	"unsafe"

	cipher "github.com/skycoin/skycoin/src/cipher"
	cli "github.com/skycoin/skycoin/src/cli"
	coin "github.com/skycoin/skycoin/src/coin"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cli_CreateRawTxFromWallet
func SKY_cli_CreateRawTxFromWallet(_c C.WebRpcClient__Handle, _walletFile, _chgAddr string, _toAddrs []C.cli__SendAmount, pwd C.PasswordReader__Handle, _arg4 *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}

	walletFile := _walletFile
	chgAddr := _chgAddr
	toAddrs := *(*[]cli.SendAmount)(unsafe.Pointer(&_toAddrs))
	pr, okp := lookupPasswordReaderHandle(pwd)
	if !okp {
		____error_code = SKY_ERROR
		return
	}
	__arg4, ____return_err := cli.CreateRawTxFromWallet(c, walletFile, chgAddr, toAddrs, *pr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg4 = registerTransactionHandle(__arg4)
	}
	return
}

//export SKY_cli_CreateRawTxFromAddress
func SKY_cli_CreateRawTxFromAddress(_c C.WebRpcClient__Handle, _addr, _walletFile, _chgAddr string, _toAddrs []C.cli__SendAmount, pwd C.PasswordReader__Handle, _arg4 *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	addr := _addr
	walletFile := _walletFile
	chgAddr := _chgAddr
	toAddrs := *(*[]cli.SendAmount)(unsafe.Pointer(&_toAddrs))
	pr, okp := lookupPasswordReaderHandle(pwd)
	if !okp {
		____error_code = SKY_ERROR
		return
	}
	__arg4, ____return_err := cli.CreateRawTxFromAddress(c, addr, walletFile, chgAddr, toAddrs, *pr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg4 = registerTransactionHandle(__arg4)
	}
	return
}

//export SKY_cli_CreateRawTx
func SKY_cli_CreateRawTx(_c C.WebRpcClient__Handle, _wlt C.Wallet__Handle, _inAddrs []string, _chgAddr string, _toAddrs []C.cli__SendAmount, _password []byte, _arg6 *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_ERROR
		return
	}
	wlt, okwlt := lookupWalletHandle(_wlt)
	if !okwlt {
		____error_code = SKY_ERROR
		return
	}
	inAddrs := *(*[]string)(unsafe.Pointer(&_inAddrs))
	chgAddr := _chgAddr
	toAddrs := *(*[]cli.SendAmount)(unsafe.Pointer(&_toAddrs))
	password := *(*[]byte)(unsafe.Pointer(&_password))
	__arg6, ____return_err := cli.CreateRawTx(c, wlt, inAddrs, chgAddr, toAddrs, password)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg6 = registerTransactionHandle(__arg6)
	}
	return
}

//export SKY_cli_NewTransaction
func SKY_cli_NewTransaction(_utxos []C.wallet__UxBalance, _keys []C.cipher__SecKey, _outs []C.coin__TransactionOutput, _arg3 *C.Transaction__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	utxos := *(*[]wallet.UxBalance)(unsafe.Pointer(&_utxos))
	keys := *(*[]cipher.SecKey)(unsafe.Pointer(&_keys))
	outs := *(*[]coin.TransactionOutput)(unsafe.Pointer(&_outs))
	__arg3 := cli.NewTransaction(utxos, keys, outs)
	*_arg3 = registerTransactionHandle(__arg3)
	return
}
