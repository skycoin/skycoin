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

// export SKY_daemon_NewPoolConfig
func SKY_daemon_NewPoolConfig(_arg0 *C.PoolConfig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := daemon.NewPoolConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofPoolConfig))
	return
}

// export SKY_daemon_NewPool
func SKY_daemon_NewPool(_c *C.PoolConfig, _d *C.Daemon, _arg2 *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	c := *(*PoolConfig)(unsafe.Pointer(_c))
	d := (*Daemon)(unsafe.Pointer(_d))
	__arg2 := daemon.NewPool(c, d)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofPool))
	return
}

// export SKY_daemon_Pool_Shutdown
func SKY_daemon_Pool_Shutdown(_pool *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pool := (*Pool)(unsafe.Pointer(_pool))
	pool.Shutdown()
	return
}

// export SKY_daemon_Pool_Run
func SKY_daemon_Pool_Run(_pool *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pool := (*Pool)(unsafe.Pointer(_pool))
	____return_err := pool.Run()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Pool_RunOffline
func SKY_daemon_Pool_RunOffline(_pool *C.Pool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pool := (*Pool)(unsafe.Pointer(_pool))
	____return_err := pool.RunOffline()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
