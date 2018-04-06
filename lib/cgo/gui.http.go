package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	gui "github.com/skycoin/skycoin/src/gui"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gui_Create
func SKY_gui_Create(_host string, _c *C.Config, _daemon *C.Daemon, _arg3 *C.Server) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	host := _host
	c := *(*cipher.Config)(unsafe.Pointer(_c))
	__arg3, ____return_err := gui.Create(host, c, daemon)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofServer))
	}
	return
}

// export SKY_gui_CreateHTTPS
func SKY_gui_CreateHTTPS(_host string, _c *C.Config, _daemon *C.Daemon, _certFile, _keyFile string, _arg4 *C.Server) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	host := _host
	c := *(*cipher.Config)(unsafe.Pointer(_c))
	certFile := _certFile
	keyFile := _keyFile
	__arg4, ____return_err := gui.CreateHTTPS(host, c, daemon, certFile, keyFile)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg4)[:]), unsafe.Pointer(_arg4), uint(SizeofServer))
	}
	return
}

// export SKY_gui_Server_Serve
func SKY_gui_Server_Serve(_s *C.Server) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := (*cipher.Server)(unsafe.Pointer(_s))
	____return_err := s.Serve()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Server_Shutdown
func SKY_gui_Server_Shutdown(_s *C.Server) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	s := (*cipher.Server)(unsafe.Pointer(_s))
	s.Shutdown()
	return
}
