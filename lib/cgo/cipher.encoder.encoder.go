package main

import (
	encoder "github.com/skycoin/skycoin/src/cipher/encoder"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_encoder_EncodeInt
func SKY_encoder_EncodeInt(_b *C.GoSlice_, _data *C.GoInterface_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(_b))
	data := convertToInterface(_data)
	encoder.EncodeInt(b, data)
	return
}

//export SKY_encoder_DecodeInt
func SKY_encoder_DecodeInt(_in *C.GoSlice_, _data *C.GoInterface_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	data := convertToInterface(_data)
	encoder.DecodeInt(in, data)
	return
}

//export SKY_encoder_DeserializeAtomic
func SKY_encoder_DeserializeAtomic(_in *C.GoSlice_, _data *C.GoInterface_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	data := convertToInterface(_data)
	encoder.DeserializeAtomic(in, data)
	return
}

//export SKY_encoder_DeserializeRaw
func SKY_encoder_DeserializeRaw(_in *C.GoSlice_, _data *C.GoInterface_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	data := convertToInterface(_data)
	____return_err := encoder.DeserializeRaw(in, data)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_encoder_SerializeAtomic
func SKY_encoder_SerializeAtomic(_data *C.GoInterface_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	data := convertToInterface(_data)
	__arg1 := encoder.SerializeAtomic(data)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_encoder_Serialize
func SKY_encoder_Serialize(_data *C.GoInterface_, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	data := convertToInterface(_data)
	__arg1 := encoder.Serialize(data)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_encoder_Size
func SKY_encoder_Size(_v *C.GoInterface_, _arg1 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	v := convertToInterface(_v)
	__arg1 := encoder.Size(v)
	*_arg1 = __arg1
	return
}
