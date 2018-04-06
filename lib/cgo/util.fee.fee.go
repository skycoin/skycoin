package main

import fee "github.com/skycoin/skycoin/src/fee"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_fee_VerifyTransactionFee
func SKY_fee_VerifyTransactionFee(_t *C.Transaction, _fee uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	fee := _fee
	____return_err := fee.VerifyTransactionFee(t, fee)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_fee_VerifyTransactionFeeForHours
func SKY_fee_VerifyTransactionFeeForHours(_hours, _fee uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hours := _hours
	fee := _fee
	____return_err := fee.VerifyTransactionFeeForHours(hours, fee)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_fee_RequiredFee
func SKY_fee_RequiredFee(_hours uint64, _arg1 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	hours := _hours
	__arg1 := fee.RequiredFee(hours)
	*_arg1 = __arg1
	return
}

// export SKY_fee_TransactionFee
func SKY_fee_TransactionFee(_tx *C.Transaction, _headTime uint64, _inUxs *C.UxArray, _arg3 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	headTime := _headTime
	__arg3, ____return_err := fee.TransactionFee(tx, headTime, inUxs)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg3 = __arg3
	}
	return
}
