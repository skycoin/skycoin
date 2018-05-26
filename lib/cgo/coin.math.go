package main

import coin "github.com/skycoin/skycoin/src/coin"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_coin_AddUint64
func SKY_coin_AddUint64(_a, _b uint64, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	a := _a
	b := _b
	__arg1, ____return_err := coin.AddUint64(a, b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

//export SKY_coin_Uint64ToInt64
func SKY_coin_Uint64ToInt64(_a uint64, _arg1 *int64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	a := _a
	__arg1, ____return_err := coin.Uint64ToInt64(a)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

//export SKY_coin_Int64ToUint64
func SKY_coin_Int64ToUint64(_a int64, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	a := _a
	__arg1, ____return_err := coin.Int64ToUint64(a)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}
