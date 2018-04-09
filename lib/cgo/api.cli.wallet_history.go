package main

import "unsafe"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_cli_byTime_Less
func SKY_cli_byTime_Less(_obt byTime, _i, _j int, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obt := *(*byTime)(unsafe.Pointer(_obt))
	i := _i
	j := _j
	__arg1 := obt.Less(i, j)
	*_arg1 = __arg1
	return
}

// export SKY_cli_byTime_Swap
func SKY_cli_byTime_Swap(_obt byTime, _i, _j int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obt := *(*byTime)(unsafe.Pointer(_obt))
	i := _i
	j := _j
	obt.Swap(i, j)
	return
}

// export SKY_cli_byTime_Len
func SKY_cli_byTime_Len(_obt byTime, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obt := *(*byTime)(unsafe.Pointer(_obt))
	__arg0 := obt.Len()
	*_arg0 = __arg0
	return
}
