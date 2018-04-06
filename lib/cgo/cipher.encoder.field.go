package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	encoder "github.com/skycoin/skycoin/src/encoder"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_encoder_StructField_String
func SKY_encoder_StructField_String(_s *C.StructField, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := (*cipher.StructField)(unsafe.Pointer(_s))
	__arg0 := s.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_encoder_DeserializeField
func SKY_encoder_DeserializeField(_in *C.GoSlice_, _fields *C.GoSlice_, _fieldName string, _field interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	fields := *(*[]cipher.StructField)(unsafe.Pointer(_fields))
	fieldName := _fieldName
	____return_err := encoder.DeserializeField(in, fields, fieldName, field)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_encoder_ParseFields
func SKY_encoder_ParseFields(_in *C.GoSlice_, _fields *C.GoSlice_, _arg2 map[string]string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	fields := *(*[]cipher.StructField)(unsafe.Pointer(_fields))
	__arg2 := encoder.ParseFields(in, fields)
	return
}
