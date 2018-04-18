package main

import (
	logging "github.com/skycoin/skycoin/src/util/logging"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_logging_NewModuleLogHook
func SKY_logging_NewModuleLogHook(_moduleName string, _arg1 *C.logging__ModuleLogHook) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	moduleName := _moduleName
	__arg1 := logging.NewModuleLogHook(moduleName)
	*_arg1 = *(*C.logging__ModuleLogHook)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_logging_ModuleLogHook_Levels
func SKY_logging_ModuleLogHook_Levels(_h *C.logging__ModuleLogHook, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	h := *(*logging.ModuleLogHook)(unsafe.Pointer(_h))
	__arg0 := h.Levels()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
