package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_daemon_RPC_GetConnection
func SKY_daemon_RPC_GetConnection(_rpc *C.RPC, _d *C.Daemon, _addr string, _arg2 *C.Connection) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	addr := _addr
	__arg2 := rpc.GetConnection(d, addr)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofConnection))
	return
}

// export SKY_daemon_RPC_GetConnections
func SKY_daemon_RPC_GetConnections(_rpc *C.RPC, _d *C.Daemon, _arg1 *C.Connections) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	__arg1 := rpc.GetConnections(d)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofConnections))
	return
}

// export SKY_daemon_RPC_GetDefaultConnections
func SKY_daemon_RPC_GetDefaultConnections(_rpc *C.RPC, _d *C.Daemon, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	__arg1 := rpc.GetDefaultConnections(d)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_daemon_RPC_GetTrustConnections
func SKY_daemon_RPC_GetTrustConnections(_rpc *C.RPC, _d *C.Daemon, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	__arg1 := rpc.GetTrustConnections(d)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_daemon_RPC_GetAllExchgConnections
func SKY_daemon_RPC_GetAllExchgConnections(_rpc *C.RPC, _d *C.Daemon, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	d := (*cipher.Daemon)(unsafe.Pointer(_d))
	__arg1 := rpc.GetAllExchgConnections(d)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_daemon_RPC_GetBlockchainProgress
func SKY_daemon_RPC_GetBlockchainProgress(_rpc *C.RPC, _v *C.Visor, _arg1 *C.BlockchainProgress) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	__arg1 := rpc.GetBlockchainProgress(v)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofBlockchainProgress))
	return
}

// export SKY_daemon_RPC_ResendTransaction
func SKY_daemon_RPC_ResendTransaction(_rpc *C.RPC, _v *C.Visor, _p *C.Pool, _txHash *C.SHA256, _arg3 *C.ResendResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	p := (*cipher.Pool)(unsafe.Pointer(_p))
	__arg3 := rpc.ResendTransaction(v, p, txHash)
	copyToBuffer(reflect.ValueOf((*__arg3)[:]), unsafe.Pointer(_arg3), uint(SizeofResendResult))
	return
}

// export SKY_daemon_RPC_ResendUnconfirmedTxns
func SKY_daemon_RPC_ResendUnconfirmedTxns(_rpc *C.RPC, _v *C.Visor, _p *C.Pool, _arg2 *C.ResendResult) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rpc := *(*cipher.RPC)(unsafe.Pointer(_rpc))
	v := (*cipher.Visor)(unsafe.Pointer(_v))
	p := (*cipher.Pool)(unsafe.Pointer(_p))
	__arg2 := rpc.ResendUnconfirmedTxns(v, p)
	copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofResendResult))
	return
}
