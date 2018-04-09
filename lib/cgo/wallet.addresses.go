package main

import (
	wallet "github.com/skycoin/skycoin/src/wallet"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_wallet_CreateAddresses
func SKY_wallet_CreateAddresses(_coinType *C.CoinType, _seed string, _genCount int, _hideSecretKey bool, _arg4 *C.ReadableWallet) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	coinType := *(*CoinType)(unsafe.Pointer(_coinType))
	seed := _seed
	genCount := _genCount
	hideSecretKey := _hideSecretKey
	__arg4, ____return_err := wallet.CreateAddresses(coinType, seed, genCount, hideSecretKey)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg4)[:]), unsafe.Pointer(_arg4), uint(SizeofReadableWallet))
	}
	return
}

// export SKY_wallet_GetSkycoinWalletEntry
func SKY_wallet_GetSkycoinWalletEntry(_pub *C.PubKey, _sec *C.SecKey, _arg2 *C.ReadableEntry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg2 := wallet.GetSkycoinWalletEntry(pub, sec)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofReadableEntry))
	return
}

// export SKY_wallet_GetBitcoinWalletEntry
func SKY_wallet_GetBitcoinWalletEntry(_pub *C.PubKey, _sec *C.SecKey, _arg2 *C.ReadableEntry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg2 := wallet.GetBitcoinWalletEntry(pub, sec)
	copyToBuffer(reflect.ValueOf(__arg2[:]), unsafe.Pointer(_arg2), uint(SizeofReadableEntry))
	return
}
