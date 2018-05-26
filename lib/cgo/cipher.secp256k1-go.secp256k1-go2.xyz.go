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

//export SKY_secp256k1go_XYZ_Print
func SKY_secp256k1go_XYZ_Print(_xyz *C.secp256k1go__XYZ, _lab string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := *(*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	lab := _lab
	xyz.Print(lab)
	return
}

//export SKY_secp256k1go_XYZ_SetXY
func SKY_secp256k1go_XYZ_SetXY(_xyz *C.secp256k1go__XYZ, _a *C.secp256k1go__XY) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	a := (*secp256k1go2.XY)(unsafe.Pointer(_a))
	xyz.SetXY(a)
	return
}

//export SKY_secp256k1go_XYZ_IsInfinity
func SKY_secp256k1go_XYZ_IsInfinity(_xyz *C.secp256k1go__XYZ, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	__arg0 := xyz.IsInfinity()
	*_arg0 = __arg0
	return
}

//export SKY_secp256k1go_XYZ_IsValid
func SKY_secp256k1go_XYZ_IsValid(_xyz *C.secp256k1go__XYZ, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	__arg0 := xyz.IsValid()
	*_arg0 = __arg0
	return
}

//export SKY_secp256k1go_XYZ_Normalize
func SKY_secp256k1go_XYZ_Normalize(_xyz *C.secp256k1go__XYZ) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	xyz.Normalize()
	return
}

//export SKY_secp256k1go_XYZ_Equals
func SKY_secp256k1go_XYZ_Equals(_xyz *C.secp256k1go__XYZ, _b *C.secp256k1go__XYZ, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	b := (*secp256k1go2.XYZ)(unsafe.Pointer(_b))
	__arg1 := xyz.Equals(b)
	*_arg1 = __arg1
	return
}

//export SKY_secp256k1go_XYZ_ECmult
func SKY_secp256k1go_XYZ_ECmult(_xyz *C.secp256k1go__XYZ, _r *C.secp256k1go__XYZ, _na, _ng *C.Number) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	r := (*secp256k1go2.XYZ)(unsafe.Pointer(_r))
	na := (*secp256k1go2.Number)(unsafe.Pointer(_na))
	ng := (*secp256k1go2.Number)(unsafe.Pointer(_ng))
	xyz.ECmult(r, na, ng)
	return
}

//export SKY_secp256k1go_XYZ_Neg
func SKY_secp256k1go_XYZ_Neg(_xyz *C.secp256k1go__XYZ, _r *C.secp256k1go__XYZ) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	r := (*secp256k1go2.XYZ)(unsafe.Pointer(_r))
	xyz.Neg(r)
	return
}

//export SKY_secp256k1go_XYZ_Double
func SKY_secp256k1go_XYZ_Double(_xyz *C.secp256k1go__XYZ, _r *C.secp256k1go__XYZ) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	r := (*secp256k1go2.XYZ)(unsafe.Pointer(_r))
	xyz.Double(r)
	return
}

//export SKY_secp256k1go_XYZ_AddXY
func SKY_secp256k1go_XYZ_AddXY(_xyz *C.secp256k1go__XYZ, _r *C.secp256k1go__XYZ, _b *C.secp256k1go__XY) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	r := (*secp256k1go2.XYZ)(unsafe.Pointer(_r))
	b := (*secp256k1go2.XY)(unsafe.Pointer(_b))
	xyz.AddXY(r, b)
	return
}

//export SKY_secp256k1go_XYZ_Add
func SKY_secp256k1go_XYZ_Add(_xyz *C.secp256k1go__XYZ, _r, _b *C.secp256k1go__XYZ) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xyz := (*secp256k1go2.XYZ)(unsafe.Pointer(_xyz))
	r := (*secp256k1go2.XYZ)(unsafe.Pointer(_r))
	b := (*secp256k1go2.XYZ)(unsafe.Pointer(_b))
	xyz.Add(r, b)
	return
}

//export SKY_secp256k1go_ECmultGen
func SKY_secp256k1go_ECmultGen(_r *C.secp256k1go__XYZ, _a *C.Number) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	r := (*secp256k1go2.XYZ)(unsafe.Pointer(_r))
	a := (*secp256k1go2.Number)(unsafe.Pointer(_a))
	secp256k1go2.ECmultGen(r, a)
	return
}
