package main

import (
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher"
)

/*
  #include "../../include/skytypes.h"
*/
import "C"

const (
	SizeofRipemd160         = unsafe.Sizeof(C.Ripemd160{})
	SizeOfAddress           = unsafe.Sizeof(C.Address{})
	SizeofPubKey            = unsafe.Sizeof(C.PubKey{})
	SizeofPubKeySlice       = unsafe.Sizeof(C.PubKeySlice{})
	SizeofSecKey            = unsafe.Sizeof(C.SecKey{})
	SizeofSig               = unsafe.Sizeof(C.Sig{})
	SizeofChecksum          = unsafe.Sizeof(C.Checksum{})
	SizeofSendAmount        = unsafe.Sizeof(C.SendAmount{})
	SizeofSHA256            = unsafe.Sizeof(C.SHA256{})
	SizeofTransactionOutput = unsafe.Sizeof(C.TransactionOutput{})
	SizeofTransaction       = unsafe.Sizeof(C.Transaction{})
	SizeofWallet            = unsafe.Sizeof(C.Wallet{})
	SizeofEntry             = unsafe.Sizeof(C.Entry{})
	SizeofUxBalance         = unsafe.Sizeof(C.UxBalance{})
)

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

func inplacePubKeySlice(p *C.PubKeySlice) *cipher.PubKeySlice {
	return (*cipher.PubKeySlice)(unsafe.Pointer(p))
}

func inplaceAddress(p *C.Address) *cipher.Address {
	return (*cipher.Address)(unsafe.Pointer(p))
}
