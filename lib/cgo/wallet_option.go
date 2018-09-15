package main

import (
	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*
  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_wallet_CreateOptionsHandle
func SKY_wallet_CreateOptionsHandle(coin string, label string, seed string, encrypt bool, pwd string, cryptoType string, scanN uint64, _opts *C.Options__Handle) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()

	var walletOptions wallet.Options
	walletOptions.Coin = (wallet.CoinType)(coin)
	walletOptions.Label = label
	walletOptions.Seed = seed
	walletOptions.Encrypt = encrypt
	walletOptions.Password = []byte(pwd)
	walletOptions.CryptoType = (wallet.CryptoType)(cryptoType)
	walletOptions.ScanN = scanN
	*_opts = registerOptionsHandle(&walletOptions)
	return
}
