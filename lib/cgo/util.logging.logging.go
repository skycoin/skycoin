package main

import logging "github.com/skycoin/skycoin/src/util/logging"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_logging_EnableColors
func SKY_logging_EnableColors() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logging.EnableColors()
	return
}

//export SKY_logging_DisableColors
func SKY_logging_DisableColors() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logging.DisableColors()
	return
}

//export SKY_logging_Disable
func SKY_logging_Disable() (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	logging.Disable()
	return
}
