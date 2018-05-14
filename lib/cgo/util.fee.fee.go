package main

import (
	coin "github.com/skycoin/skycoin/src/coin"
	fee "github.com/skycoin/skycoin/src/util/fee"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_fee_VerifyTransactionFee
func SKY_fee_VerifyTransactionFee(_t *C.coin__Transaction, _fee uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	t := (*coin.Transaction)(unsafe.Pointer(_t))
	____return_err := fee.VerifyTransactionFee(t, _fee)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_fee_VerifyTransactionFeeForHours
func SKY_fee_VerifyTransactionFeeForHours(_hours, _fee uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hours := _hours
	____return_err := fee.VerifyTransactionFeeForHours(hours, _fee)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_fee_RequiredFee
func SKY_fee_RequiredFee(_hours uint64, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hours := _hours
	__arg1 := fee.RequiredFee(hours)
	*_arg1 = __arg1
	return
}

//export SKY_fee_RemainingHours
func SKY_fee_RemainingHours(_hours uint64, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hours := _hours
	__arg1 := fee.RemainingHours(hours)
	*_arg1 = __arg1
	return
}

//export SKY_fee_TransactionFee
func SKY_fee_TransactionFee(_tx *C.coin__Transaction, _headTime uint64, _inUxs *C.coin__UxArray, _arg3 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	tx := (*coin.Transaction)(unsafe.Pointer(_tx))
	headTime := _headTime
	inUxs := *(*coin.UxArray)(unsafe.Pointer(_inUxs))
	__arg3, ____return_err := fee.TransactionFee(tx, headTime, inUxs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg3 = __arg3
	}
	return
}
