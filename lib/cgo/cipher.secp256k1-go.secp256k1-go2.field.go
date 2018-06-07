package main

import (
	"unsafe"

	secp256k1go2 "github.com/skycoin/skycoin/src/cipher/secp256k1-go/secp256k1-go2"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_secp256k1go_Field_String
func SKY_secp256k1go_Field_String(_fd *C.secp256k1go__Field, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	__arg0 := fd.String()
	copyString(__arg0, _arg0)
	return
}

//export SKY_secp256k1go_Field_Print
func SKY_secp256k1go_Field_Print(_fd *C.secp256k1go__Field, _lab string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	lab := _lab
	fd.Print(lab)
	return
}

//export SKY_secp256k1go_Field_SetB32
func SKY_secp256k1go_Field_SetB32(_fd *C.secp256k1go__Field, _a []byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	a := *(*[]byte)(unsafe.Pointer(&_a))
	fd.SetB32(a)
	return
}

//export SKY_secp256k1go_Field_SetBytes
func SKY_secp256k1go_Field_SetBytes(_fd *C.secp256k1go__Field, _a []byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	a := *(*[]byte)(unsafe.Pointer(&_a))
	fd.SetBytes(a)
	return
}

//export SKY_secp256k1go_Field_SetHex
func SKY_secp256k1go_Field_SetHex(_fd *C.secp256k1go__Field, _s string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	s := _s
	fd.SetHex(s)
	return
}

//export SKY_secp256k1go_Field_IsOdd
func SKY_secp256k1go_Field_IsOdd(_fd *C.secp256k1go__Field, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	__arg0 := fd.IsOdd()
	*_arg0 = __arg0
	return
}

//export SKY_secp256k1go_Field_IsZero
func SKY_secp256k1go_Field_IsZero(_fd *C.secp256k1go__Field, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	__arg0 := fd.IsZero()
	*_arg0 = __arg0
	return
}

//export SKY_secp256k1go_Field_SetInt
func SKY_secp256k1go_Field_SetInt(_fd *C.secp256k1go__Field, _a uint32) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	a := _a
	fd.SetInt(a)
	return
}

//export SKY_secp256k1go_Field_Normalize
func SKY_secp256k1go_Field_Normalize(_fd *C.secp256k1go__Field) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	fd.Normalize()
	return
}

//export SKY_secp256k1go_Field_GetB32
func SKY_secp256k1go_Field_GetB32(_fd *C.secp256k1go__Field, _r []byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	r := *(*[]byte)(unsafe.Pointer(&_r))
	fd.GetB32(r)
	return
}

//export SKY_secp256k1go_Field_Equals
func SKY_secp256k1go_Field_Equals(_fd *C.secp256k1go__Field, _b *C.secp256k1go__Field, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	b := (*secp256k1go2.Field)(unsafe.Pointer(_b))
	__arg1 := fd.Equals(b)
	*_arg1 = __arg1
	return
}

//export SKY_secp256k1go_Field_SetAdd
func SKY_secp256k1go_Field_SetAdd(_fd *C.secp256k1go__Field, _a *C.secp256k1go__Field) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	a := (*secp256k1go2.Field)(unsafe.Pointer(_a))
	fd.SetAdd(a)
	return
}

//export SKY_secp256k1go_Field_MulInt
func SKY_secp256k1go_Field_MulInt(_fd *C.secp256k1go__Field, _a uint32) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	a := _a
	fd.MulInt(a)
	return
}

//export SKY_secp256k1go_Field_Negate
func SKY_secp256k1go_Field_Negate(_fd *C.secp256k1go__Field, _r *C.secp256k1go__Field, _m uint32) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	r := (*secp256k1go2.Field)(unsafe.Pointer(_r))
	m := _m
	fd.Negate(r, m)
	return
}

//export SKY_secp256k1go_Field_Inv
func SKY_secp256k1go_Field_Inv(_fd *C.secp256k1go__Field, _r *C.secp256k1go__Field) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	r := (*secp256k1go2.Field)(unsafe.Pointer(_r))
	fd.Inv(r)
	return
}

//export SKY_secp256k1go_Field_Sqrt
func SKY_secp256k1go_Field_Sqrt(_fd *C.secp256k1go__Field, _r *C.secp256k1go__Field) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	r := (*secp256k1go2.Field)(unsafe.Pointer(_r))
	fd.Sqrt(r)
	return
}

//export SKY_secp256k1go_Field_InvVar
func SKY_secp256k1go_Field_InvVar(_fd *C.secp256k1go__Field, _r *C.secp256k1go__Field) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	r := (*secp256k1go2.Field)(unsafe.Pointer(_r))
	fd.InvVar(r)
	return
}

//export SKY_secp256k1go_Field_Mul
func SKY_secp256k1go_Field_Mul(_fd *C.secp256k1go__Field, _r, _b *C.secp256k1go__Field) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	r := (*secp256k1go2.Field)(unsafe.Pointer(_r))
	b := (*secp256k1go2.Field)(unsafe.Pointer(_b))
	fd.Mul(r, b)
	return
}

//export SKY_secp256k1go_Field_Sqr
func SKY_secp256k1go_Field_Sqr(_fd *C.secp256k1go__Field, _r *C.secp256k1go__Field) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	fd := (*secp256k1go2.Field)(unsafe.Pointer(_fd))
	r := (*secp256k1go2.Field)(unsafe.Pointer(_r))
	fd.Sqr(r)
	return
}
