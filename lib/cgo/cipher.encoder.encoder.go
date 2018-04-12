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

// export SKY_encoder_EncodeInt
func SKY_encoder_EncodeInt(_b *C.GoSlice_, _data interface{}) (____error_code uint32) {
	____error_code = 0
	b := *(*[]byte)(unsafe.Pointer(_b))
	encoder.EncodeInt(b, data)
	return
}

// export SKY_encoder_DecodeInt
func SKY_encoder_DecodeInt(_in *C.GoSlice_, _data interface{}) (____error_code uint32) {
	____error_code = 0
	in := *(*[]byte)(unsafe.Pointer(_in))
	encoder.DecodeInt(in, data)
	return
}

// export SKY_encoder_DeserializeAtomic
func SKY_encoder_DeserializeAtomic(_in *C.GoSlice_, _data interface{}) (____error_code uint32) {
	____error_code = 0
	in := *(*[]byte)(unsafe.Pointer(_in))
	encoder.DeserializeAtomic(in, data)
	return
}

// export SKY_encoder_DeserializeRaw
func SKY_encoder_DeserializeRaw(_in *C.GoSlice_, _data interface{}) (____error_code uint32) {
	____error_code = 0
	in := *(*[]byte)(unsafe.Pointer(_in))
	____return_err := encoder.DeserializeRaw(in, data)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_encoder_SerializeAtomic
func SKY_encoder_SerializeAtomic(_data interface{}, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	__arg1 := encoder.SerializeAtomic(data)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_encoder_Serialize
func SKY_encoder_Serialize(_data interface{}, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	__arg1 := encoder.Serialize(data)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_encoder_Size
func SKY_encoder_Size(_v interface{}, _arg1 *int) (____error_code uint32) {
	____error_code = 0
	__arg1 := encoder.Size(v)
	*_arg1 = __arg1
	return
}
