package main

import (
	daemon "github.com/skycoin/skycoin/src/daemon"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_daemon_NewConfig
func SKY_daemon_NewConfig(_arg0 *C.Config) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := daemon.NewConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofConfig))
	return
}

// export SKY_daemon_NewDaemonConfig
func SKY_daemon_NewDaemonConfig(_arg0 *C.DaemonConfig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := daemon.NewDaemonConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofDaemonConfig))
	return
}

// export SKY_daemon_NewDaemon
func SKY_daemon_NewDaemon(_config *C.Config, _db *C.DB, _defaultConns *C.GoSlice_, _arg3 *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	config := *(*Config)(unsafe.Pointer(_config))
	defaultConns := *(*[]string)(unsafe.Pointer(_defaultConns))
	__arg3, ____return_err := daemon.NewDaemon(config, db, defaultConns)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofDaemon))
	}
	return
}

// export SKY_daemon_Daemon_Shutdown
func SKY_daemon_Daemon_Shutdown(_dm *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dm := (*Daemon)(unsafe.Pointer(_dm))
	dm.Shutdown()
	return
}

// export SKY_daemon_Daemon_Run
func SKY_daemon_Daemon_Run(_dm *C.Daemon) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dm := (*Daemon)(unsafe.Pointer(_dm))
	____return_err := dm.Run()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Daemon_GetListenPort
func SKY_daemon_Daemon_GetListenPort(_dm *C.Daemon, _addr string, _arg1 *uint16) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dm := (*Daemon)(unsafe.Pointer(_dm))
	addr := _addr
	__arg1 := dm.GetListenPort(addr)
	*_arg1 = __arg1
	return
}
