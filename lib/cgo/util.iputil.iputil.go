package main

import iputil "github.com/skycoin/skycoin/src/iputil"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_iputil_LocalhostIP
func SKY_iputil_LocalhostIP(_arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0, ____return_err := iputil.LocalhostIP()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg0, _arg0)
	}
	return
}

// export SKY_iputil_IsLocalhost
func SKY_iputil_IsLocalhost(_addr string, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	addr := _addr
	__arg1 := iputil.IsLocalhost(addr)
	*_arg1 = __arg1
	return
}

// export SKY_iputil_SplitAddr
func SKY_iputil_SplitAddr(_addr string, _arg1 *C.GoString_, _arg2 *uint16) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	addr := _addr
	__arg1, __arg2, ____return_err := iputil.SplitAddr(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
		*_arg2 = __arg2
	}
	return
}
