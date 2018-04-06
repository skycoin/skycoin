package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	visor "github.com/skycoin/skycoin/src/visor"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_NewErrTxnViolatesHardConstraint
func SKY_visor_NewErrTxnViolatesHardConstraint(_err error) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	err := *(*cipher.error)(unsafe.Pointer(_err))
	____return_err := visor.NewErrTxnViolatesHardConstraint(err)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_ErrTxnViolatesHardConstraint_Error
func SKY_visor_ErrTxnViolatesHardConstraint_Error(_e *C.ErrTxnViolatesHardConstraint, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	e := *(*cipher.ErrTxnViolatesHardConstraint)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_visor_NewErrTxnViolatesSoftConstraint
func SKY_visor_NewErrTxnViolatesSoftConstraint(_err error) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	err := *(*cipher.error)(unsafe.Pointer(_err))
	____return_err := visor.NewErrTxnViolatesSoftConstraint(err)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_ErrTxnViolatesSoftConstraint_Error
func SKY_visor_ErrTxnViolatesSoftConstraint_Error(_e *C.ErrTxnViolatesSoftConstraint, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	e := *(*cipher.ErrTxnViolatesSoftConstraint)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_visor_VerifySingleTxnSoftConstraints
func SKY_visor_VerifySingleTxnSoftConstraints(_txn *C.Transaction, _headTime uint64, _uxIn *C.UxArray, _maxSize int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	headTime := _headTime
	maxSize := _maxSize
	____return_err := visor.VerifySingleTxnSoftConstraints(txn, headTime, uxIn, maxSize)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_VerifySingleTxnHardConstraints
func SKY_visor_VerifySingleTxnHardConstraints(_txn *C.Transaction, _head *C.SignedBlock, _uxIn *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	____return_err := visor.VerifySingleTxnHardConstraints(txn, head, uxIn)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_VerifyBlockTxnConstraints
func SKY_visor_VerifyBlockTxnConstraints(_txn *C.Transaction, _head *C.SignedBlock, _uxIn *C.UxArray) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	____return_err := visor.VerifyBlockTxnConstraints(txn, head, uxIn)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
