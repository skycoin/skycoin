package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_logging_TextFormatter_SetColorScheme
func SKY_logging_TextFormatter_SetColorScheme(_f *C.TextFormatter, _colorScheme *C.ColorScheme) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	f := (*cipher.TextFormatter)(unsafe.Pointer(_f))
	colorScheme := (*cipher.ColorScheme)(unsafe.Pointer(_colorScheme))
	f.SetColorScheme(colorScheme)
	return
}

// export SKY_logging_TextFormatter_Format
func SKY_logging_TextFormatter_Format(_f *C.TextFormatter, _entry *C.Entry, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	f := (*cipher.TextFormatter)(unsafe.Pointer(_f))
	__arg1, ____return_err := f.Format(entry)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}
