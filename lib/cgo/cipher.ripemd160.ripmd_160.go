package main

import (
	ripemd160 "github.com/skycoin/skycoin/src/cipher/ripemd160"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_ripemd160_digest_Reset
func SKY_ripemd160_digest_Reset(_d digest) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	d := (*digest)(unsafe.Pointer(_d))
	d.Reset()
	return
}

// export SKY_ripemd160_New
func SKY_ripemd160_New(_arg0 *C.Hash) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := ripemd160.New()
	return
}

// export SKY_ripemd160_digest_Size
func SKY_ripemd160_digest_Size(_d digest, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	d := (*digest)(unsafe.Pointer(_d))
	__arg0 := d.Size()
	*_arg0 = __arg0
	return
}

// export SKY_ripemd160_digest_BlockSize
func SKY_ripemd160_digest_BlockSize(_d digest, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	d := (*digest)(unsafe.Pointer(_d))
	__arg0 := d.BlockSize()
	*_arg0 = __arg0
	return
}

// export SKY_ripemd160_digest_Write
func SKY_ripemd160_digest_Write(_d digest, _p *C.GoSlice_, _arg1 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	d := (*digest)(unsafe.Pointer(_d))
	p := *(*[]byte)(unsafe.Pointer(_p))
	__arg1, ____return_err := d.Write(p)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_ripemd160_digest_Sum
func SKY_ripemd160_digest_Sum(_d0 digest, _in *C.GoSlice_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	d0 := (*digest)(unsafe.Pointer(_d0))
	in := *(*[]byte)(unsafe.Pointer(_in))
	__arg1 := d0.Sum(in)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}
