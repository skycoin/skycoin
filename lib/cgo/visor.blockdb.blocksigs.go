package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_blockdb_blockSigs_Get
func SKY_blockdb_blockSigs_Get(_bs blockSigs, _hash *C.SHA256, _arg1 *C.Sig, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bs := *(*cipher.blockSigs)(unsafe.Pointer(_bs))
	__arg1, __arg2, ____return_err := bs.Get(hash)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = __arg2
	}
	return
}

// export SKY_blockdb_blockSigs_AddWithTx
func SKY_blockdb_blockSigs_AddWithTx(_bs blockSigs, _tx *C.Tx, _hash *C.SHA256, _sig *C.Sig) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	bs := (*cipher.blockSigs)(unsafe.Pointer(_bs))
	____return_err := bs.AddWithTx(tx, hash, sig)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
