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

// export SKY_logging_NewLogger
func SKY_logging_NewLogger(_priorityKey, _criticalPriority string, _arg1 *C.Logger) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	priorityKey := _priorityKey
	criticalPriority := _criticalPriority
	__arg1 := logging.NewLogger(priorityKey, criticalPriority)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofLogger))
	return
}

// export SKY_logging_LoggerForModules
func SKY_logging_LoggerForModules(_priorityKey, _criticalPriority string, _enabledModules *C.GoSlice_, _arg2 *C.Logger) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	priorityKey := _priorityKey
	criticalPriority := _criticalPriority
	enabledModules := *(*[]string)(unsafe.Pointer(_enabledModules))
	__arg2 := logging.LoggerForModules(priorityKey, criticalPriority, enabledModules)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofLogger))
	return
}

// export SKY_logging_Logger_MustGetLogger
func SKY_logging_Logger_MustGetLogger(_logger *C.Logger, _moduleName string, _arg1 *C.Logger) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	moduleName := _moduleName
	__arg1 := logger.MustGetLogger(moduleName)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofLogger))
	return
}

// export SKY_logging_Logger_DisableAllModules
func SKY_logging_Logger_DisableAllModules(_logger *C.Logger) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	logger.DisableAllModules()
	return
}

// export SKY_logging_Logger_EnableModules
func SKY_logging_Logger_EnableModules(_logger *C.Logger, _modules *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	modules := *(*[]string)(unsafe.Pointer(_modules))
	logger.EnableModules(modules)
	return
}

// export SKY_logging_Logger_Criticalf
func SKY_logging_Logger_Criticalf(_logger *C.Logger, _format string, _args ...interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	format := _format
	args := _args
	logger.Criticalf(format, args)
	return
}

// export SKY_logging_Logger_Critical
func SKY_logging_Logger_Critical(_logger *C.Logger, _args ...interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	args := _args
	logger.Critical(args)
	return
}

// export SKY_logging_Logger_Criticalln
func SKY_logging_Logger_Criticalln(_logger *C.Logger, _args ...interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	args := _args
	logger.Criticalln(args)
	return
}

// export SKY_logging_Logger_Noticef
func SKY_logging_Logger_Noticef(_logger *C.Logger, _format string, _args ...interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	format := _format
	args := _args
	logger.Noticef(format, args)
	return
}

// export SKY_logging_Logger_Notice
func SKY_logging_Logger_Notice(_logger *C.Logger, _args ...interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	args := _args
	logger.Notice(args)
	return
}

// export SKY_logging_Logger_Noticeln
func SKY_logging_Logger_Noticeln(_logger *C.Logger, _args ...interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	args := _args
	logger.Noticeln(args)
	return
}

// export SKY_logging_Logger_Disable
func SKY_logging_Logger_Disable(_logger *C.Logger) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logger := (*Logger)(unsafe.Pointer(_logger))
	logger.Disable()
	return
}
