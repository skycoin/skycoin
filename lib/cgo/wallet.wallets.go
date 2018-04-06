package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	wallet "github.com/skycoin/skycoin/src/wallet"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_wallet_LoadWallets
func SKY_wallet_LoadWallets(_dir string, _arg1 *C.Wallets) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dir := _dir
	__arg1, ____return_err := wallet.LoadWallets(dir)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofWallets))
	}
	return
}

// export SKY_wallet_Wallets_ToReadable
func SKY_wallet_Wallets_ToReadable(_wlts *C.Wallets, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	wlts := *(*cipher.Wallets)(unsafe.Pointer(_wlts))
	__arg0 := wlts.ToReadable()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
