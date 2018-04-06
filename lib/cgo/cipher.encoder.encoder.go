package main

import (
	encoder "github.com/skycoin/skycoin/src/encoder"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_encoder_EncodeInt
func SKY_encoder_EncodeInt(_b *C.GoSlice_, _data interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(_b))
	encoder.EncodeInt(b, data)
	return
}

// export SKY_encoder_DecodeInt
func SKY_encoder_DecodeInt(_in *C.GoSlice_, _data interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	encoder.DecodeInt(in, data)
	return
}

// export SKY_encoder_DeserializeAtomic
func SKY_encoder_DeserializeAtomic(_in *C.GoSlice_, _data interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	encoder.DeserializeAtomic(in, data)
	return
}

// export SKY_encoder_DeserializeRaw
func SKY_encoder_DeserializeRaw(_in *C.GoSlice_, _data interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	____return_err := encoder.DeserializeRaw(in, data)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_encoder_Deserialize
func SKY_encoder_Deserialize(_r *C.Reader, _dsize int, _data interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dsize := _dsize
	____return_err := encoder.Deserialize(r, dsize, data)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_encoder_CanDeserialize
func SKY_encoder_CanDeserialize(_in *C.GoSlice_, _dst *C.Value, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	__arg2 := encoder.CanDeserialize(in, dst)
	*_arg2 = __arg2
	return
}

// export SKY_encoder_DeserializeRawToValue
func SKY_encoder_DeserializeRawToValue(_in *C.GoSlice_, _dst *C.Value, _arg2 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	in := *(*[]byte)(unsafe.Pointer(_in))
	__arg2, ____return_err := encoder.DeserializeRawToValue(in, dst)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = __arg2
	}
	return
}

// export SKY_encoder_DeserializeToValue
func SKY_encoder_DeserializeToValue(_r *C.Reader, _dsize int, _dst *C.Value) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dsize := _dsize
	____return_err := encoder.DeserializeToValue(r, dsize, dst)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_encoder_SerializeAtomic
func SKY_encoder_SerializeAtomic(_data interface{}, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1 := encoder.SerializeAtomic(data)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_encoder_Serialize
func SKY_encoder_Serialize(_data interface{}, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1 := encoder.Serialize(data)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_encoder_Size
func SKY_encoder_Size(_v interface{}, _arg1 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg1 := encoder.Size(v)
	*_arg1 = __arg1
	return
}
