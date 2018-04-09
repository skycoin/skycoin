package main

import require "github.com/skycoin/skycoin/src/testutil/require"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_require_PanicsWithCondition
func SKY_require_PanicsWithCondition(_t *C.TestingT, _condition *C.TestValuePredicate, _f *C.PanicTestFunc, _msgAndArgs ...interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	msgAndArgs := _msgAndArgs
	require.PanicsWithCondition(t, condition, f, msgAndArgs)
	return
}

// export SKY_require_PanicsWithLogMessage
func SKY_require_PanicsWithLogMessage(_t *C.TestingT, _expectedMessage string, _f *C.PanicTestFunc, _msgAndArgs ...interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	expectedMessage := _expectedMessage
	msgAndArgs := _msgAndArgs
	require.PanicsWithLogMessage(t, expectedMessage, f, msgAndArgs)
	return
}
