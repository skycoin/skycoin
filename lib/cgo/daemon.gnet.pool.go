package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	gnet "github.com/skycoin/skycoin/src/gnet"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_gnet_NewConfig
func SKY_gnet_NewConfig(_arg0 *C.Config) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := gnet.NewConfig()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofConfig))
	return
}

// export SKY_gnet_NewConnection
func SKY_gnet_NewConnection(_pool *C.ConnectionPool, _id int, _conn *C.Conn, _writeQueueSize int, _solicited bool, _arg5 *C.Connection) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	id := _id
	writeQueueSize := _writeQueueSize
	solicited := _solicited
	__arg5 := gnet.NewConnection(pool, id, conn, writeQueueSize, solicited)
	copyToBuffer(reflect.ValueOf((*__arg5)[:]), unsafe.Pointer(_arg5), uint(SizeofConnection))
	return
}

// export SKY_gnet_Connection_Addr
func SKY_gnet_Connection_Addr(_conn *C.Connection, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	conn := (*cipher.Connection)(unsafe.Pointer(_conn))
	__arg0 := conn.Addr()
	copyString(__arg0, _arg0)
	return
}

// export SKY_gnet_Connection_String
func SKY_gnet_Connection_String(_conn *C.Connection, _arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	conn := (*cipher.Connection)(unsafe.Pointer(_conn))
	__arg0 := conn.String()
	copyString(__arg0, _arg0)
	return
}

// export SKY_gnet_Connection_Close
func SKY_gnet_Connection_Close(_conn *C.Connection) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	conn := (*cipher.Connection)(unsafe.Pointer(_conn))
	conn.Close()
	return
}

// export SKY_gnet_NewConnectionPool
func SKY_gnet_NewConnectionPool(_c *C.Config, _state interface{}, _arg2 *C.ConnectionPool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	c := *(*cipher.Config)(unsafe.Pointer(_c))
	__arg2 := gnet.NewConnectionPool(c, state)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofConnectionPool))
	return
}

// export SKY_gnet_ConnectionPool_Run
func SKY_gnet_ConnectionPool_Run(_pool *C.ConnectionPool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	____return_err := pool.Run()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_ConnectionPool_RunOffline
func SKY_gnet_ConnectionPool_RunOffline(_pool *C.ConnectionPool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	____return_err := pool.RunOffline()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_ConnectionPool_Shutdown
func SKY_gnet_ConnectionPool_Shutdown(_pool *C.ConnectionPool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	pool.Shutdown()
	return
}

// export SKY_gnet_ConnectionPool_NewConnection
func SKY_gnet_ConnectionPool_NewConnection(_pool *C.ConnectionPool, _conn *C.Conn, _solicited bool, _arg2 *C.Connection) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	solicited := _solicited
	__arg2, ____return_err := pool.NewConnection(conn, solicited)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofConnection))
	}
	return
}

// export SKY_gnet_ConnectionPool_ListeningAddress
func SKY_gnet_ConnectionPool_ListeningAddress(_pool *C.ConnectionPool, _arg0 *C.Addr) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	__arg0, ____return_err := pool.ListeningAddress()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_ConnectionPool_IsConnExist
func SKY_gnet_ConnectionPool_IsConnExist(_pool *C.ConnectionPool, _addr string, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	addr := _addr
	__arg1, ____return_err := pool.IsConnExist(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_gnet_ConnectionPool_GetConnection
func SKY_gnet_ConnectionPool_GetConnection(_pool *C.ConnectionPool, _addr string, _arg1 *C.Connection) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	addr := _addr
	__arg1, ____return_err := pool.GetConnection(addr)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofConnection))
	}
	return
}

// export SKY_gnet_ConnectionPool_Connect
func SKY_gnet_ConnectionPool_Connect(_pool *C.ConnectionPool, _address string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	address := _address
	____return_err := pool.Connect(address)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_ConnectionPool_Disconnect
func SKY_gnet_ConnectionPool_Disconnect(_pool *C.ConnectionPool, _addr string, _r *C.DisconnectReason) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	addr := _addr
	r := *(*cipher.DisconnectReason)(unsafe.Pointer(_r))
	____return_err := pool.Disconnect(addr, r)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_ConnectionPool_GetConnections
func SKY_gnet_ConnectionPool_GetConnections(_pool *C.ConnectionPool, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	__arg0, ____return_err := pool.GetConnections()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_gnet_ConnectionPool_Size
func SKY_gnet_ConnectionPool_Size(_pool *C.ConnectionPool, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	__arg0, ____return_err := pool.Size()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg0 = __arg0
	}
	return
}

// export SKY_gnet_ConnectionPool_SendMessage
func SKY_gnet_ConnectionPool_SendMessage(_pool *C.ConnectionPool, _addr string, _msg *C.Message) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	addr := _addr
	msg := *(*cipher.Message)(unsafe.Pointer(_msg))
	____return_err := pool.SendMessage(addr, msg)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_ConnectionPool_BroadcastMessage
func SKY_gnet_ConnectionPool_BroadcastMessage(_pool *C.ConnectionPool, _msg *C.Message) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	msg := *(*cipher.Message)(unsafe.Pointer(_msg))
	____return_err := pool.BroadcastMessage(msg)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_ConnectionPool_SendPings
func SKY_gnet_ConnectionPool_SendPings(_pool *C.ConnectionPool, _rate *C.Duration, _msg *C.Message) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	msg := *(*cipher.Message)(unsafe.Pointer(_msg))
	____return_err := pool.SendPings(rate, msg)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_ConnectionPool_ClearStaleConnections
func SKY_gnet_ConnectionPool_ClearStaleConnections(_pool *C.ConnectionPool, _idleLimit *C.Duration, _reason *C.DisconnectReason) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pool := (*cipher.ConnectionPool)(unsafe.Pointer(_pool))
	reason := *(*cipher.DisconnectReason)(unsafe.Pointer(_reason))
	____return_err := pool.ClearStaleConnections(idleLimit, reason)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_gnet_Now
func SKY_gnet_Now(_arg0 *C.Time) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := gnet.Now()
	return
}
