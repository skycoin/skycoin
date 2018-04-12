package main

import (
	coin "github.com/skycoin/skycoin/src/coin"
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
	err := _err
	____return_err := visor.NewErrTxnViolatesHardConstraint(err)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_ErrTxnViolatesHardConstraint_Error
func SKY_visor_ErrTxnViolatesHardConstraint_Error(_e *C.visor__ErrTxnViolatesHardConstraint, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	e := *(*visor.ErrTxnViolatesHardConstraint)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_visor_NewErrTxnViolatesSoftConstraint
func SKY_visor_NewErrTxnViolatesSoftConstraint(_err error) (____error_code uint32) {
	____error_code = 0
	err := _err
	____return_err := visor.NewErrTxnViolatesSoftConstraint(err)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_ErrTxnViolatesSoftConstraint_Error
func SKY_visor_ErrTxnViolatesSoftConstraint_Error(_e *C.visor__ErrTxnViolatesSoftConstraint, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	e := *(*visor.ErrTxnViolatesSoftConstraint)(unsafe.Pointer(_e))
	__arg0 := e.Error()
	copyString(__arg0, _arg0)
	return
}

// export SKY_visor_VerifySingleTxnSoftConstraints
func SKY_visor_VerifySingleTxnSoftConstraints(_txn *C.coin__Transaction, _headTime uint64, _uxIn *C.coin__UxArray, _maxSize int) (____error_code uint32) {
	____error_code = 0
	txn := *(*coin.Transaction)(unsafe.Pointer(_txn))
	headTime := _headTime
	uxIn := *(*coin.UxArray)(unsafe.Pointer(_uxIn))
	maxSize := _maxSize
	____return_err := visor.VerifySingleTxnSoftConstraints(txn, headTime, uxIn, maxSize)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_VerifySingleTxnHardConstraints
func SKY_visor_VerifySingleTxnHardConstraints(_txn *C.coin__Transaction, _head *C.coin__SignedBlock, _uxIn *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	txn := *(*coin.Transaction)(unsafe.Pointer(_txn))
	head := (*coin.SignedBlock)(unsafe.Pointer(_head))
	uxIn := *(*coin.UxArray)(unsafe.Pointer(_uxIn))
	____return_err := visor.VerifySingleTxnHardConstraints(txn, head, uxIn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_visor_VerifyBlockTxnConstraints
func SKY_visor_VerifyBlockTxnConstraints(_txn *C.coin__Transaction, _head *C.coin__SignedBlock, _uxIn *C.coin__UxArray) (____error_code uint32) {
	____error_code = 0
	txn := *(*coin.Transaction)(unsafe.Pointer(_txn))
	head := (*coin.SignedBlock)(unsafe.Pointer(_head))
	uxIn := *(*coin.UxArray)(unsafe.Pointer(_uxIn))
	____return_err := visor.VerifyBlockTxnConstraints(txn, head, uxIn)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
