package main

import (
	elapse "github.com/skycoin/skycoin/src/util/elapse"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_elapse_NewElapser
func SKY_elapse_NewElapser(_elapsedThreshold *C.Duration, _logger *C.Logger, _arg2 *C.Elapser) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg2 := elapse.NewElapser(elapsedThreshold, logger)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofElapser))
	return
}

// export SKY_elapse_Elapser_CheckForDone
func SKY_elapse_Elapser_CheckForDone(_e *C.Elapser) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := (*Elapser)(unsafe.Pointer(_e))
	e.CheckForDone()
	return
}

// export SKY_elapse_Elapser_Register
func SKY_elapse_Elapser_Register(_e *C.Elapser, _name string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := (*Elapser)(unsafe.Pointer(_e))
	name := _name
	e.Register(name)
	return
}

// export SKY_elapse_Elapser_ShowCurrentTime
func SKY_elapse_Elapser_ShowCurrentTime(_e *C.Elapser, _step string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := (*Elapser)(unsafe.Pointer(_e))
	step := _step
	e.ShowCurrentTime(step)
	return
}

// export SKY_elapse_Elapser_Elapsed
func SKY_elapse_Elapser_Elapsed(_e *C.Elapser) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	e := (*Elapser)(unsafe.Pointer(_e))
	e.Elapsed()
	return
}
