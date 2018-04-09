package main

import (
	visor "github.com/skycoin/skycoin/src/visor"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_visor_NewErrTxnViolatesHardConstraint
func SKY_visor_NewErrTxnViolatesHardConstraint(_err error) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	err := *(*error)(unsafe.Pointer(_err))
	____return_err := visor.NewErrTxnViolatesHardConstraint(err)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_ErrTxnViolatesHardConstraint_Error
func SKY_visor_ErrTxnViolatesHardConstraint_Error(_e *C.ErrTxnViolatesHardConstraint, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := *(*ErrTxnViolatesHardConstraint)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_visor_NewErrTxnViolatesSoftConstraint
func SKY_visor_NewErrTxnViolatesSoftConstraint(_err error) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	err := *(*error)(unsafe.Pointer(_err))
	____return_err := visor.NewErrTxnViolatesSoftConstraint(err)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_ErrTxnViolatesSoftConstraint_Error
func SKY_visor_ErrTxnViolatesSoftConstraint_Error(_e *C.ErrTxnViolatesSoftConstraint, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := *(*ErrTxnViolatesSoftConstraint)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_visor_VerifySingleTxnSoftConstraints
func SKY_visor_VerifySingleTxnSoftConstraints(_txn *C.Transaction, _headTime uint64, _uxIn *C.UxArray, _maxSize int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	headTime := _headTime
	maxSize := _maxSize
	____return_err := visor.VerifySingleTxnSoftConstraints(txn, headTime, uxIn, maxSize)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_VerifySingleTxnHardConstraints
func SKY_visor_VerifySingleTxnHardConstraints(_txn *C.Transaction, _head *C.SignedBlock, _uxIn *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	____return_err := visor.VerifySingleTxnHardConstraints(txn, head, uxIn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_VerifyBlockTxnConstraints
func SKY_visor_VerifyBlockTxnConstraints(_txn *C.Transaction, _head *C.SignedBlock, _uxIn *C.UxArray) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	____return_err := visor.VerifyBlockTxnConstraints(txn, head, uxIn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
