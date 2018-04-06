package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	daemon "github.com/skycoin/skycoin/src/daemon"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_daemon_NewConfig
func SKY_daemon_NewConfig(_arg0 *C.Config) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofConfig))
	return
}

// export SKY_daemon_NewDaemonConfig
func SKY_daemon_NewDaemonConfig(_arg0 *C.DaemonConfig) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewDaemonConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofDaemonConfig))
	return
}

// export SKY_daemon_NewDaemon
func SKY_daemon_NewDaemon(_config *C.Config, _db *C.DB, _defaultConns *C.GoSlice_, _arg3 *C.Daemon) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	config := *(*cipher.Config)(unsafe.Pointer(_config))
	defaultConns := *(*[]string)(unsafe.Pointer(_defaultConns))
	__arg3, ____return_err := daemon.NewDaemon(config, db, defaultConns)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofDaemon))
	}
	return
}

// export SKY_daemon_Daemon_Shutdown
func SKY_daemon_Daemon_Shutdown(_dm *C.Daemon) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dm := (*cipher.Daemon)(unsafe.Pointer(_dm))
	dm.Shutdown()
	return
}

// export SKY_daemon_Daemon_Run
func SKY_daemon_Daemon_Run(_dm *C.Daemon) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dm := (*cipher.Daemon)(unsafe.Pointer(_dm))
	____return_err := dm.Run()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Daemon_GetListenPort
func SKY_daemon_Daemon_GetListenPort(_dm *C.Daemon, _addr string, _arg1 *uint16) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dm := (*cipher.Daemon)(unsafe.Pointer(_dm))
	addr := _addr
	__arg1 := dm.GetListenPort(addr)
	*_arg1 = __arg1
	return
}
