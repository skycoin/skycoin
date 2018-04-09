package main

import (
	assert "github.com/skycoin/skycoin/src/testutil/assert"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_testutil_PanicsWithCondition
func SKY_testutil_PanicsWithCondition(_t *C.TestingT, _condition *C.TestValuePredicate, _f *C.PanicTestFunc, _msgAndArgs ...interface{}, _arg4 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	condition := *(*TestValuePredicate)(unsafe.Pointer(_condition))
	msgAndArgs := _msgAndArgs
	__arg4 := assert.PanicsWithCondition(t, condition, f, msgAndArgs)
	*_arg4 = __arg4
	return
}

// export SKY_testutil_PanicsWithLogMessage
func SKY_testutil_PanicsWithLogMessage(_t *C.TestingT, _expectedMessage string, _f *C.PanicTestFunc, _msgAndArgs ...interface{}, _arg4 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	expectedMessage := _expectedMessage
	msgAndArgs := _msgAndArgs
	__arg4 := assert.PanicsWithLogMessage(t, expectedMessage, f, msgAndArgs)
	*_arg4 = __arg4
	return
}
