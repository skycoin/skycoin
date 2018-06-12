package main

import apputil "github.com/skycoin/skycoin/src/util/apputil"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_apputil_CatchInterruptPanic
func SKY_apputil_CatchInterruptPanic() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	apputil.CatchInterruptPanic()
	return
}

//export SKY_apputil_CatchDebug
func SKY_apputil_CatchDebug() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	apputil.CatchDebug()
	return
}

//export SKY_apputil_PrintProgramStatus
func SKY_apputil_PrintProgramStatus() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	apputil.PrintProgramStatus()
	return
}
