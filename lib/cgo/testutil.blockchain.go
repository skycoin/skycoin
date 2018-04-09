package main

import (
	testutil "github.com/skycoin/skycoin/src/testutil"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_testutil_PrepareDB
func SKY_testutil_PrepareDB(_t *C.T, _arg1 *C.DB, _arg2 C.Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1, __arg2 := testutil.PrepareDB(t)
	return
}

// export SKY_testutil_RequireError
func SKY_testutil_RequireError(_t *C.T, _err error, _msg string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	err := *(*error)(unsafe.Pointer(_err))
	msg := _msg
	testutil.RequireError(t, err, msg)
	return
}

// export SKY_testutil_MakeAddress
func SKY_testutil_MakeAddress(_arg0 *C.Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := testutil.MakeAddress()
	return
}

// export SKY_testutil_RandBytes
func SKY_testutil_RandBytes(_t *C.T, _n int, _arg2 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	n := _n
	__arg2 := testutil.RandBytes(t, n)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

// export SKY_testutil_RandSHA256
func SKY_testutil_RandSHA256(_t *C.T, _arg1 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg1 := testutil.RandSHA256(t)
	return
}

// export SKY_testutil_SHA256FromHex
func SKY_testutil_SHA256FromHex(_t *C.T, _hex string, _arg2 *C.SHA256) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	hex := _hex
	__arg2 := testutil.SHA256FromHex(t, hex)
	return
}
