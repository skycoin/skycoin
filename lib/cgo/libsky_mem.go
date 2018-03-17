package main

import (
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher"
)

/*
  #include "../../include/skytypes.h"
*/
import "C"

type Handle uint64

var (
	handleMap = make(map[Handle]interface{})
)

func openHandle(obj interface{}) Handle {
	ptr := &obj
	handle := *(*Handle)(unsafe.Pointer(&ptr))
	handleMap[handle] = obj
	return handle
}

func lookupHandleObj(handle Handle) (interface{}, bool) {
	obj, ok := handleMap[handle]
	return obj, ok
}

func closeHandle(handle Handle) {
	delete(handleMap, handle)
}

func inplaceByteArray(p unsafe.Pointer, length int) *[]byte {
	// Create slice without copying data
	// TODO: Memory efficiency
	slice := (*[1 << 30]byte)(p)[:length:length]
	return &slice
}

func inplacePubKey(p *C.PubKey) *cipher.PubKey {
	return (*cipher.PubKey)(unsafe.Pointer(inplaceByteArray(unsafe.Pointer(p), 33)))
}

func inplacePubKeySlice(p *C.PubKeySlice) *cipher.PubKeySlice {
	// Create slice without copying data
	slice := (*[1 << 30]cipher.PubKey)(p.data)[:p.len:p.len]
	return (*cipher.PubKeySlice)(unsafe.Pointer(&slice))
}

func inplaceSecKey(p *C.SecKey) *cipher.SecKey {
	return (*cipher.SecKey)(unsafe.Pointer(inplaceByteArray(unsafe.Pointer(p), 32)))
}

func inplaceSig(p *C.Sig) *cipher.Sig {
	return (*cipher.Sig)(unsafe.Pointer(inplaceByteArray(unsafe.Pointer(p), 64 + 1)))
}

func inplaceChecksum(p *C.Checksum) *cipher.Checksum {
	return (*cipher.Checksum)(unsafe.Pointer(inplaceByteArray(unsafe.Pointer(p), 4)))
}

func inplaceRipemd160(p *C.Ripemd160) *cipher.Ripemd160 {
	return (*cipher.Ripemd160)(unsafe.Pointer(inplaceByteArray(unsafe.Pointer(p), 20)))
}

func inplaceSHA256(p *C.SHA256) *cipher.SHA256 {
	return (*cipher.SHA256)(unsafe.Pointer(inplaceByteArray(unsafe.Pointer(p), 20)))
}

func inplaceAddress(p *C.Address) *cipher.Address {
	return (*cipher.Address)(unsafe.Pointer(p))
}

