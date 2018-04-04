package main

import (
	"reflect"
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher"
)

/*
	#include <string.h>

  #include "../../include/skytypes.h"

	void eos(char *str, int len) {
		str[len] = 0;
	}

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

/**
 * Inplace memory references
 */

func inplacePubKeySlice(p *C.PubKeySlice) *cipher.PubKeySlice {
	return (*cipher.PubKeySlice)(unsafe.Pointer(p))
}

func inplaceAddress(p *C.Address) *cipher.Address {
	return (*cipher.Address)(unsafe.Pointer(p))
}

/**
 * Copy helpers
 */

func copyString(src string, dest *C.GoString_) bool {
	srcLen := len(src)
	dest.p = (*C.char)(C.malloc(C.size_t(srcLen + 1)))
	strAddr := (*C.GoString_)(unsafe.Pointer(&src))
	C.memcpy(unsafe.Pointer(dest.p), unsafe.Pointer(strAddr.p), C.size_t(srcLen))
	C.eos(dest.p, C.int(srcLen))
	dest.n = C.GoInt_(srcLen)
	return true
}

// Determine the memory address of a slice buffer and the
// size of its underlaying element type
func getBufferData(src reflect.Value) (bufferAddr unsafe.Pointer, elemSize C.size_t) {
	firstElem := src.Index(0)
	elemSize = C.size_t(firstElem.Type().Size())
	bufferAddr = unsafe.Pointer(src.Pointer())
	return
}

// Copy n items in source slice/array/string onto C-managed memory buffer
//
// This function takes for granted that all values in src
// will be instances of the same type, and that src and dest
// element types will be aligned exactly the same
// in memory of the same size
func copyToBuffer(src reflect.Value, dest unsafe.Pointer, n uint) {
	srcLen := src.Len()
	if srcLen == 0 {
		return
	}
	srcAddr, elemSize := getBufferData(src)
	C.memcpy(dest, srcAddr, C.size_t(n)*elemSize)
}

// Copy source slice/array/string onto instance of C.GSlice struct
//
// This function takes for granted that all values in src
// will be instances of the same type, and that src and dest
// element types will be aligned exactly the same
// in memory of the same size
func copyToGoSlice(src reflect.Value, dest *C.GoSlice_) {
	srcLen := src.Len()
	if srcLen == 0 {
		dest.len = 0
		return
	}
	srcAddr, elemSize := getBufferData(src)
	if dest.cap == 0 {
		dest.data = C.malloc(C.size_t(srcLen) * elemSize)
		dest.cap = C.GoInt_(srcLen)
	}
	n, overflow := srcLen, srcLen > int(dest.cap)
	if overflow {
		n = int(dest.cap)
	}
	C.memcpy(dest.data, srcAddr, C.size_t(n)*elemSize)
	// Do not modify slice metadata until memory is actually copied
	if overflow {
		dest.len = dest.cap - C.GoInt_(srcLen)
	} else {
		dest.len = C.GoInt_(srcLen)
	}
}
