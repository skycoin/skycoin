package main

import (
	"reflect"
	"unsafe"

	http "github.com/skycoin/skycoin/src/util/http"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_httphelper_Address_UnmarshalJSON
func SKY_httphelper_Address_UnmarshalJSON(_a *C.httphelper__Address, _b []byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	a := inplaceHttpHelperAddress(_a)
	b := *(*[]byte)(unsafe.Pointer(&_b))
	____return_err := a.UnmarshalJSON(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_httphelper_Address_MarshalJSON
func SKY_httphelper_Address_MarshalJSON(_a *C.httphelper__Address, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	a := *inplaceHttpHelperAddress(_a)
	__arg0, ____return_err := a.MarshalJSON()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

//export SKY_httphelper_Coins_UnmarshalJSON
func SKY_httphelper_Coins_UnmarshalJSON(_c *C.httphelper__Coins, _b []byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := (*http.Coins)(unsafe.Pointer(_c))
	b := *(*[]byte)(unsafe.Pointer(&_b))
	____return_err := c.UnmarshalJSON(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_httphelper_Coins_MarshalJSON
func SKY_httphelper_Coins_MarshalJSON(_c *C.httphelper__Coins, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*http.Coins)(unsafe.Pointer(_c))
	__arg0, ____return_err := c.MarshalJSON()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

//export SKY_httphelper_Coins_Value
func SKY_httphelper_Coins_Value(_c *C.httphelper__Coins, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*http.Coins)(unsafe.Pointer(_c))
	__arg0 := c.Value()
	*_arg0 = __arg0
	return
}

//export SKY_httphelper_Hours_UnmarshalJSON
func SKY_httphelper_Hours_UnmarshalJSON(_h *C.httphelper__Hours, _b []byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	h := (*http.Hours)(unsafe.Pointer(_h))
	b := *(*[]byte)(unsafe.Pointer(&_b))
	____return_err := h.UnmarshalJSON(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_httphelper_Hours_MarshalJSON
func SKY_httphelper_Hours_MarshalJSON(_h *C.httphelper__Hours, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	h := *(*http.Hours)(unsafe.Pointer(_h))
	__arg0, ____return_err := h.MarshalJSON()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

//export SKY_httphelper_Hours_Value
func SKY_httphelper_Hours_Value(_h *C.httphelper__Hours, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	h := *(*http.Hours)(unsafe.Pointer(_h))
	__arg0 := h.Value()
	*_arg0 = __arg0
	return
}
