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

// export SKY_daemon_NewPoolConfig
func SKY_daemon_NewPoolConfig(_arg0 *C.PoolConfig) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewPoolConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofPoolConfig))
	return
}

// export SKY_daemon_NewPool
func SKY_daemon_NewPool(_c *C.PoolConfig, _d *C.Daemon, _arg2 *C.Pool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := *(*cipher.PoolConfig)(unsafe.Pointer(_c))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	__arg2 := daemon.NewPool(c, d)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofPool))
	return
}

// export SKY_daemon_Pool_Shutdown
func SKY_daemon_Pool_Shutdown(_pool *C.Pool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.Pool)(unsafe.Pointer(_pool))
	pool.Shutdown()
	return
}

// export SKY_daemon_Pool_Run
func SKY_daemon_Pool_Run(_pool *C.Pool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.Pool)(unsafe.Pointer(_pool))
	____return_err := pool.Run()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_daemon_Pool_RunOffline
func SKY_daemon_Pool_RunOffline(_pool *C.Pool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.Pool)(unsafe.Pointer(_pool))
	____return_err := pool.RunOffline()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
