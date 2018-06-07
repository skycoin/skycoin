package main

import (
	"unsafe"

	encoder "github.com/skycoin/skycoin/src/cipher/encoder"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
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

//export SKY_encoder_ParseFields
func SKY_encoder_ParseFields(_in []byte, _fields []C.encoder__StructField, _arg2 *C.GoStringMap_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(&_in))
	fields := *(*[]encoder.StructField)(unsafe.Pointer(&_fields))
	__arg2 := encoder.ParseFields(in, fields)
	copyToStringMap(__arg2, _arg2)
	return
}
