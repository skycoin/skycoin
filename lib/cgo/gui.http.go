package main

import (
	gui "github.com/skycoin/skycoin/src/gui"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gui_Create
func SKY_gui_Create(_host string, _c *C.Config, _daemon *C.Daemon, _arg3 *C.Server) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	host := _host
	c := *(*Config)(unsafe.Pointer(_c))
	__arg3, ____return_err := gui.Create(host, c, daemon)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofServer))
	}
	return
}

// export SKY_gui_CreateHTTPS
func SKY_gui_CreateHTTPS(_host string, _c *C.Config, _daemon *C.Daemon, _certFile, _keyFile string, _arg4 *C.Server) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	host := _host
	c := *(*Config)(unsafe.Pointer(_c))
	certFile := _certFile
	keyFile := _keyFile
	__arg4, ____return_err := gui.CreateHTTPS(host, c, daemon, certFile, keyFile)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg4)[:]), unsafe.Pointer(_arg4), uint(SizeofServer))
	}
	return
}

// export SKY_gui_Server_Serve
func SKY_gui_Server_Serve(_s *C.Server) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	s := (*Server)(unsafe.Pointer(_s))
	____return_err := s.Serve()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gui_Server_Shutdown
func SKY_gui_Server_Shutdown(_s *C.Server) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	s := (*Server)(unsafe.Pointer(_s))
	s.Shutdown()
	return
}
