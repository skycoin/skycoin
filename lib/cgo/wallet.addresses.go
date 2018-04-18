package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	wallet "github.com/skycoin/skycoin/src/wallet"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_wallet_CreateAddresses
func SKY_wallet_CreateAddresses(_coinType *C.wallet__CoinType, _seed string, _genCount int, _hideSecretKey bool, _arg4 *C.wallet__ReadableWallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	coinType := *(*wallet.CoinType)(unsafe.Pointer(_coinType))
	seed := _seed
	genCount := _genCount
	hideSecretKey := _hideSecretKey
	__arg4, ____return_err := wallet.CreateAddresses(coinType, seed, genCount, hideSecretKey)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg4 = *(*C.wallet__ReadableWallet)(unsafe.Pointer(__arg4))
	}
	return
}

//export SKY_wallet_GetSkycoinWalletEntry
func SKY_wallet_GetSkycoinWalletEntry(_pub *C.cipher__PubKey, _sec *C.cipher__SecKey, _arg2 *C.wallet__ReadableEntry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pub := *(*cipher.PubKey)(unsafe.Pointer(_pub))
	sec := *(*cipher.SecKey)(unsafe.Pointer(_sec))
	__arg2 := wallet.GetSkycoinWalletEntry(pub, sec)
	*_arg2 = *(*C.wallet__ReadableEntry)(unsafe.Pointer(&__arg2))
	return
}

//export SKY_wallet_GetBitcoinWalletEntry
func SKY_wallet_GetBitcoinWalletEntry(_pub *C.cipher__PubKey, _sec *C.cipher__SecKey, _arg2 *C.wallet__ReadableEntry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pub := *(*cipher.PubKey)(unsafe.Pointer(_pub))
	sec := *(*cipher.SecKey)(unsafe.Pointer(_sec))
	__arg2 := wallet.GetBitcoinWalletEntry(pub, sec)
	*_arg2 = *(*C.wallet__ReadableEntry)(unsafe.Pointer(&__arg2))
	return
}
