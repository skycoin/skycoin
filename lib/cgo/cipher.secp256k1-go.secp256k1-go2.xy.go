package main

import (
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_secp256k1go_XY_Print
func SKY_secp256k1go_XY_Print(_xy *C.XY, _lab string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	lab := _lab
	xy.Print(lab)
	return
}

// export SKY_secp256k1go_XY_ParsePubkey
func SKY_secp256k1go_XY_ParsePubkey(_xy *C.XY, _pub *C.GoSlice_, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	pub := *(*[]byte)(unsafe.Pointer(_pub))
	__arg1 := xy.ParsePubkey(pub)
	*_arg1 = __arg1
	return
}

// export SKY_secp256k1go_XY_Bytes
func SKY_secp256k1go_XY_Bytes(_xy *C.XY, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := *(*XY)(unsafe.Pointer(_xy))
	__arg0 := xy.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_secp256k1go_XY_BytesUncompressed
func SKY_secp256k1go_XY_BytesUncompressed(_xy *C.XY, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	__arg0 := xy.BytesUncompressed()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_secp256k1go_XY_SetXY
func SKY_secp256k1go_XY_SetXY(_xy *C.XY, _X, _Y *C.Field) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	X := (*Field)(unsafe.Pointer(_X))
	Y := (*Field)(unsafe.Pointer(_Y))
	xy.SetXY(X, Y)
	return
}

// export SKY_secp256k1go_XY_IsValid
func SKY_secp256k1go_XY_IsValid(_xy *C.XY, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	__arg0 := xy.IsValid()
	*_arg0 = __arg0
	return
}

// export SKY_secp256k1go_XY_SetXYZ
func SKY_secp256k1go_XY_SetXYZ(_xy *C.XY, _a *C.XYZ) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	a := (*XYZ)(unsafe.Pointer(_a))
	xy.SetXYZ(a)
	return
}

// export SKY_secp256k1go_XY_Neg
func SKY_secp256k1go_XY_Neg(_xy *C.XY, _r *C.XY) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	r := (*XY)(unsafe.Pointer(_r))
	xy.Neg(r)
	return
}

// export SKY_secp256k1go_XY_SetXO
func SKY_secp256k1go_XY_SetXO(_xy *C.XY, _X *C.Field, _odd bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	X := (*Field)(unsafe.Pointer(_X))
	odd := _odd
	xy.SetXO(X, odd)
	return
}

// export SKY_secp256k1go_XY_AddXY
func SKY_secp256k1go_XY_AddXY(_xy *C.XY, _a *C.XY) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	a := (*XY)(unsafe.Pointer(_a))
	xy.AddXY(a)
	return
}

// export SKY_secp256k1go_XY_GetPublicKey
func SKY_secp256k1go_XY_GetPublicKey(_xy *C.XY, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	xy := (*XY)(unsafe.Pointer(_xy))
	__arg0 := xy.GetPublicKey()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
