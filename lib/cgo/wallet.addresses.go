package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

// FIXME: Removed from src/ . What's the way to do it now?
/*
//export SKY_wallet_CreateAddresses
func SKY_wallet_CreateAddresses(_coinType string, _seed string, _genCount int, _hideSecretKey bool, _arg4 *C.ReadableWallet__Handle) (____error_code uint32) {
	coinType := _coinType
	seed := _seed
	genCount := _genCount
	hideSecretKey := _hideSecretKey
	__arg4, ____return_err := wallet.CreateAddresses(wallet.CoinType(coinType), seed, genCount, hideSecretKey)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg4 = registerReadableWalletHandle(__arg4)
	}
	return
}
*/
