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

// export SKY_logging_TextFormatter_SetColorScheme
func SKY_logging_TextFormatter_SetColorScheme(_f *C.TextFormatter, _colorScheme *C.ColorScheme) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	f := (*TextFormatter)(unsafe.Pointer(_f))
	colorScheme := (*ColorScheme)(unsafe.Pointer(_colorScheme))
	f.SetColorScheme(colorScheme)
	return
}

// export SKY_logging_TextFormatter_Format
func SKY_logging_TextFormatter_Format(_f *C.TextFormatter, _entry *C.Entry, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	f := (*TextFormatter)(unsafe.Pointer(_f))
	__arg1, ____return_err := f.Format(entry)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}
