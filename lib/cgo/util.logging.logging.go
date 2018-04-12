package main

import logging "github.com/skycoin/skycoin/src/util/logging"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_logging_Disable
func SKY_logging_Disable() (____error_code uint32) {
	____error_code = 0
	logging.Disable()
	return
}
