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

// export SKY_logging_NewModuleLogHook
func SKY_logging_NewModuleLogHook(_moduleName string, _arg1 *C.logging__ModuleLogHook) (____error_code uint32) {
	____error_code = 0
	moduleName := _moduleName
	__arg1 := logging.NewModuleLogHook(moduleName)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofModuleLogHook))
	return
}
