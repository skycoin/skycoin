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

// export SKY_logging_NewReplayHook
func SKY_logging_NewReplayHook(_logger *C.Logger, _arg1 *C.ReplayHook) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	__arg1 := logging.NewReplayHook(logger)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReplayHook))
	return
}

// export SKY_logging_ReplayHook_Levels
func SKY_logging_ReplayHook_Levels(_h *C.ReplayHook, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	h := *(*ReplayHook)(unsafe.Pointer(_h))
	__arg0 := h.Levels()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_logging_ReplayHook_Fire
func SKY_logging_ReplayHook_Fire(_h *C.ReplayHook, _entry *C.Entry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	h := *(*ReplayHook)(unsafe.Pointer(_h))
	____return_err := h.Fire(entry)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_logging_NewModuleLogHook
func SKY_logging_NewModuleLogHook(_moduleName string, _arg1 *C.ModuleLogHook) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	moduleName := _moduleName
	__arg1 := logging.NewModuleLogHook(moduleName)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofModuleLogHook))
	return
}

// export SKY_logging_ModuleLogHook_Levels
func SKY_logging_ModuleLogHook_Levels(_h *C.ModuleLogHook, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	h := *(*ModuleLogHook)(unsafe.Pointer(_h))
	__arg0 := h.Levels()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

// export SKY_logging_ModuleLogHook_Fire
func SKY_logging_ModuleLogHook_Fire(_h *C.ModuleLogHook, _entry *C.Entry) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	h := *(*ModuleLogHook)(unsafe.Pointer(_h))
	____return_err := h.Fire(entry)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
