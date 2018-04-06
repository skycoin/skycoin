package main

import droplet "github.com/skycoin/skycoin/src/droplet"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_droplet_FromString
func SKY_droplet_FromString(_b string, _arg1 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := _b
	__arg1, ____return_err := droplet.FromString(b)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_droplet_ToString
func SKY_droplet_ToString(_n uint64, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	n := _n
	__arg1, ____return_err := droplet.ToString(n)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}
