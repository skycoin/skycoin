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

// export SKY_daemon_NewExpectIntroductions
func SKY_daemon_NewExpectIntroductions(_arg0 *C.ExpectIntroductions) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewExpectIntroductions()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofExpectIntroductions))
	return
}

// export SKY_daemon_ExpectIntroductions_Add
func SKY_daemon_ExpectIntroductions_Add(_ei *C.ExpectIntroductions, _addr string, _tm *C.Time) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ei := (*cipher.ExpectIntroductions)(unsafe.Pointer(_ei))
	addr := _addr
	ei.Add(addr, tm)
	return
}

// export SKY_daemon_ExpectIntroductions_Remove
func SKY_daemon_ExpectIntroductions_Remove(_ei *C.ExpectIntroductions, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ei := (*cipher.ExpectIntroductions)(unsafe.Pointer(_ei))
	addr := _addr
	ei.Remove(addr)
	return
}

// export SKY_daemon_ExpectIntroductions_CullInvalidConns
func SKY_daemon_ExpectIntroductions_CullInvalidConns(_ei *C.ExpectIntroductions, _f *C.CullMatchFunc, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ei := (*cipher.ExpectIntroductions)(unsafe.Pointer(_ei))
	f := *(*cipher.CullMatchFunc)(unsafe.Pointer(_f))
	__arg1, ____return_err := ei.CullInvalidConns(f)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	}
	return
}

// export SKY_daemon_ExpectIntroductions_Get
func SKY_daemon_ExpectIntroductions_Get(_ei *C.ExpectIntroductions, _addr string, _arg1 *C.Time, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ei := (*cipher.ExpectIntroductions)(unsafe.Pointer(_ei))
	addr := _addr
	__arg1, __arg2 := ei.Get(addr)
	*_arg2 = __arg2
	return
}

// export SKY_daemon_NewConnectionMirrors
func SKY_daemon_NewConnectionMirrors(_arg0 *C.ConnectionMirrors) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewConnectionMirrors()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofConnectionMirrors))
	return
}

// export SKY_daemon_ConnectionMirrors_Add
func SKY_daemon_ConnectionMirrors_Add(_cm *C.ConnectionMirrors, _addr string, _mirror uint32) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	cm := (*cipher.ConnectionMirrors)(unsafe.Pointer(_cm))
	addr := _addr
	mirror := _mirror
	cm.Add(addr, mirror)
	return
}

// export SKY_daemon_ConnectionMirrors_Get
func SKY_daemon_ConnectionMirrors_Get(_cm *C.ConnectionMirrors, _addr string, _arg1 *uint32, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	cm := (*cipher.ConnectionMirrors)(unsafe.Pointer(_cm))
	addr := _addr
	__arg1, __arg2 := cm.Get(addr)
	*_arg1 = __arg1
	*_arg2 = __arg2
	return
}

// export SKY_daemon_ConnectionMirrors_Remove
func SKY_daemon_ConnectionMirrors_Remove(_cm *C.ConnectionMirrors, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	cm := (*cipher.ConnectionMirrors)(unsafe.Pointer(_cm))
	addr := _addr
	cm.Remove(addr)
	return
}

// export SKY_daemon_NewOutgoingConnections
func SKY_daemon_NewOutgoingConnections(_max int, _arg1 *C.OutgoingConnections) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	max := _max
	__arg1 := daemon.NewOutgoingConnections(max)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofOutgoingConnections))
	return
}

// export SKY_daemon_OutgoingConnections_Add
func SKY_daemon_OutgoingConnections_Add(_oc *C.OutgoingConnections, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	oc := (*cipher.OutgoingConnections)(unsafe.Pointer(_oc))
	addr := _addr
	oc.Add(addr)
	return
}

// export SKY_daemon_OutgoingConnections_Remove
func SKY_daemon_OutgoingConnections_Remove(_oc *C.OutgoingConnections, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	oc := (*cipher.OutgoingConnections)(unsafe.Pointer(_oc))
	addr := _addr
	oc.Remove(addr)
	return
}

// export SKY_daemon_OutgoingConnections_Get
func SKY_daemon_OutgoingConnections_Get(_oc *C.OutgoingConnections, _addr string, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	oc := (*cipher.OutgoingConnections)(unsafe.Pointer(_oc))
	addr := _addr
	__arg1 := oc.Get(addr)
	*_arg1 = __arg1
	return
}

// export SKY_daemon_OutgoingConnections_Len
func SKY_daemon_OutgoingConnections_Len(_oc *C.OutgoingConnections, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	oc := (*cipher.OutgoingConnections)(unsafe.Pointer(_oc))
	__arg0 := oc.Len()
	*_arg0 = __arg0
	return
}

// export SKY_daemon_NewPendingConnections
func SKY_daemon_NewPendingConnections(_maxConn int, _arg1 *C.PendingConnections) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	maxConn := _maxConn
	__arg1 := daemon.NewPendingConnections(maxConn)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofPendingConnections))
	return
}

// export SKY_daemon_PendingConnections_Add
func SKY_daemon_PendingConnections_Add(_pc *C.PendingConnections, _addr string, _peer *C.Peer) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pc := (*cipher.PendingConnections)(unsafe.Pointer(_pc))
	addr := _addr
	pc.Add(addr, peer)
	return
}

// export SKY_daemon_PendingConnections_Get
func SKY_daemon_PendingConnections_Get(_pc *C.PendingConnections, _addr string, _arg1 *C.Peer, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pc := (*cipher.PendingConnections)(unsafe.Pointer(_pc))
	addr := _addr
	__arg1, __arg2 := pc.Get(addr)
	*_arg2 = __arg2
	return
}

// export SKY_daemon_PendingConnections_Remove
func SKY_daemon_PendingConnections_Remove(_pc *C.PendingConnections, _addr string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pc := (*cipher.PendingConnections)(unsafe.Pointer(_pc))
	addr := _addr
	pc.Remove(addr)
	return
}

// export SKY_daemon_PendingConnections_Len
func SKY_daemon_PendingConnections_Len(_pc *C.PendingConnections, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	pc := (*cipher.PendingConnections)(unsafe.Pointer(_pc))
	__arg0 := pc.Len()
	*_arg0 = __arg0
	return
}

// export SKY_daemon_NewMirrorConnections
func SKY_daemon_NewMirrorConnections(_arg0 *C.MirrorConnections) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewMirrorConnections()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofMirrorConnections))
	return
}

// export SKY_daemon_MirrorConnections_Add
func SKY_daemon_MirrorConnections_Add(_mc *C.MirrorConnections, _mirror uint32, _ip string, _port uint16) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	mc := (*cipher.MirrorConnections)(unsafe.Pointer(_mc))
	mirror := _mirror
	ip := _ip
	port := _port
	mc.Add(mirror, ip, port)
	return
}

// export SKY_daemon_MirrorConnections_Get
func SKY_daemon_MirrorConnections_Get(_mc *C.MirrorConnections, _mirror uint32, _ip string, _arg2 *uint16, _arg3 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	mc := (*cipher.MirrorConnections)(unsafe.Pointer(_mc))
	mirror := _mirror
	ip := _ip
	__arg2, __arg3 := mc.Get(mirror, ip)
	*_arg2 = __arg2
	*_arg3 = __arg3
	return
}

// export SKY_daemon_MirrorConnections_Remove
func SKY_daemon_MirrorConnections_Remove(_mc *C.MirrorConnections, _mirror uint32, _ip string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	mc := (*cipher.MirrorConnections)(unsafe.Pointer(_mc))
	mirror := _mirror
	ip := _ip
	mc.Remove(mirror, ip)
	return
}

// export SKY_daemon_NewIPCount
func SKY_daemon_NewIPCount(_arg0 *C.IPCount) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := daemon.NewIPCount()
	copyToBuffer(reflect.ValueOf((*__arg0)[:]), unsafe.Pointer(_arg0), uint(SizeofIPCount))
	return
}

// export SKY_daemon_IPCount_Increase
func SKY_daemon_IPCount_Increase(_ic *C.IPCount, _ip string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ic := (*cipher.IPCount)(unsafe.Pointer(_ic))
	ip := _ip
	ic.Increase(ip)
	return
}

// export SKY_daemon_IPCount_Decrease
func SKY_daemon_IPCount_Decrease(_ic *C.IPCount, _ip string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ic := (*cipher.IPCount)(unsafe.Pointer(_ic))
	ip := _ip
	ic.Decrease(ip)
	return
}

// export SKY_daemon_IPCount_Get
func SKY_daemon_IPCount_Get(_ic *C.IPCount, _ip string, _arg1 *int, _arg2 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	ic := (*cipher.IPCount)(unsafe.Pointer(_ic))
	ip := _ip
	__arg1, __arg2 := ic.Get(ip)
	*_arg1 = __arg1
	*_arg2 = __arg2
	return
}
