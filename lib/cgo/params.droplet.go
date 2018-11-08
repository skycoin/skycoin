package main

import (
	params "github.com/skycoin/skycoin/src/params"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export
func SKY_params_MaxDropletDivisor() uint64 {
	return params.MaxDropletDivisor()
}

//export
func SKY_params_DropletPrecisionCheck(amount uint64) uint32 {
	return libErrorCode(params.DropletPrecisionCheck(amount))
}
