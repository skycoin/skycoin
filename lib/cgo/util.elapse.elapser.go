package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	elapse "github.com/skycoin/skycoin/src/elapse"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_elapse_NewElapser
func SKY_elapse_NewElapser(_elapsedThreshold *C.Duration, _logger *C.Logger, _arg2 *C.Elapser) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg2 := elapse.NewElapser(elapsedThreshold, logger)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofElapser))
	return
}

// export SKY_elapse_Elapser_CheckForDone
func SKY_elapse_Elapser_CheckForDone(_e *C.Elapser) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	e := (*cipher.Elapser)(unsafe.Pointer(_e))
	e.CheckForDone()
	return
}

// export SKY_elapse_Elapser_Register
func SKY_elapse_Elapser_Register(_e *C.Elapser, _name string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	e := (*cipher.Elapser)(unsafe.Pointer(_e))
	name := _name
	e.Register(name)
	return
}

// export SKY_elapse_Elapser_ShowCurrentTime
func SKY_elapse_Elapser_ShowCurrentTime(_e *C.Elapser, _step string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	e := (*cipher.Elapser)(unsafe.Pointer(_e))
	step := _step
	e.ShowCurrentTime(step)
	return
}

// export SKY_elapse_Elapser_Elapsed
func SKY_elapse_Elapser_Elapsed(_e *C.Elapser) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	e := (*cipher.Elapser)(unsafe.Pointer(_e))
	e.Elapsed()
	return
}
