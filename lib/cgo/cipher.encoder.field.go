package main

import (
	encoder "github.com/skycoin/skycoin/src/cipher/encoder"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_encoder_StructField_String
func SKY_encoder_StructField_String(_s *C.encoder__StructField, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	s := (*encoder.StructField)(unsafe.Pointer(_s))
	__arg0 := s.String()
	copyString(__arg0, _arg0)
	return
}

//export SKY_encoder_DeserializeField
func SKY_encoder_DeserializeField(_in *C.GoSlice_, _fields *C.GoSlice_, _fieldName string, _field *C.GoInterface_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	fields := *(*[]encoder.StructField)(unsafe.Pointer(_fields))
	fieldName := _fieldName
	field := copyToInterface(_field)
	____return_err := encoder.DeserializeField(in, fields, fieldName, field)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_encoder_ParseFields
func SKY_encoder_ParseFields(_in *C.GoSlice_, _fields *C.GoSlice_, _arg2 *C.GoMap_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	fields := *(*[]encoder.StructField)(unsafe.Pointer(_fields))
	__arg2 := encoder.ParseFields(in, fields)
	copyStringMap(__arg2, _arg2)
	return
}
