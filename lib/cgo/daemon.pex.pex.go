package main

import (
	pex "github.com/skycoin/skycoin/src/daemon/pex"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_pex_NewPeer
func SKY_pex_NewPeer(_address string, _arg1 *C.pex__Peer) (____error_code uint32) {
	____error_code = 0
	address := _address
	__arg1 := pex.NewPeer(address)
	copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofPeer))
	return
}

// export SKY_pex_Peer_Seen
func SKY_pex_Peer_Seen(_peer *C.pex__Peer) (____error_code uint32) {
	____error_code = 0
	peer := (*pex.Peer)(unsafe.Pointer(_peer))
	peer.Seen()
	return
}

// export SKY_pex_Peer_IncreaseRetryTimes
func SKY_pex_Peer_IncreaseRetryTimes(_peer *C.pex__Peer) (____error_code uint32) {
	____error_code = 0
	peer := (*pex.Peer)(unsafe.Pointer(_peer))
	peer.IncreaseRetryTimes()
	return
}

// export SKY_pex_Peer_ResetRetryTimes
func SKY_pex_Peer_ResetRetryTimes(_peer *C.pex__Peer) (____error_code uint32) {
	____error_code = 0
	peer := (*pex.Peer)(unsafe.Pointer(_peer))
	peer.ResetRetryTimes()
	return
}

// export SKY_pex_Peer_CanTry
func SKY_pex_Peer_CanTry(_peer *C.pex__Peer, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	peer := (*pex.Peer)(unsafe.Pointer(_peer))
	__arg0 := peer.CanTry()
	*_arg0 = __arg0
	return
}

// export SKY_pex_Peer_String
func SKY_pex_Peer_String(_peer *C.pex__Peer, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	peer := (*pex.Peer)(unsafe.Pointer(_peer))
	__arg0 := peer.String()
	copyString(__arg0, _arg0)
	return
}
