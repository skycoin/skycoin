package main

import (
	"hash"
	"reflect"
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/http"
)

/*
	#include <string.h>

  #include "skytypes.h"

	void eos(char *str, int len) {
		str[len] = 0;
	}

*/
import "C"

const (
	SizeofRipemd160         = unsafe.Sizeof(C.cipher__Ripemd160{})
	SizeOfAddress           = unsafe.Sizeof(C.cipher__Address{})
	SizeofPubKey            = unsafe.Sizeof(C.cipher__PubKey{})
	SizeofPubKeySlice       = unsafe.Sizeof(C.cipher__PubKeySlice{})
	SizeofSecKey            = unsafe.Sizeof(C.cipher__SecKey{})
	SizeofSig               = unsafe.Sizeof(C.cipher__Sig{})
	SizeofChecksum          = unsafe.Sizeof(C.cipher__Checksum{})
	SizeofSendAmount        = unsafe.Sizeof(C.cli__SendAmount{})
	SizeofSHA256            = unsafe.Sizeof(C.cipher__SHA256{})
	SizeofTransactionOutput = unsafe.Sizeof(C.coin__TransactionOutput{})
	SizeofTransaction       = unsafe.Sizeof(C.coin__Transaction{})
	SizeofWallet            = unsafe.Sizeof(C.wallet__Wallet{})
	SizeofEntry             = unsafe.Sizeof(C.wallet__Entry{})
	SizeofUxBalance         = unsafe.Sizeof(C.wallet__UxBalance{})
)

/**
 * Inplace memory references
 */

func inplacePubKeySlice(p *C.cipher__PubKeySlice) *cipher.PubKeySlice {
	return (*cipher.PubKeySlice)(unsafe.Pointer(p))
}

func inplaceAddress(p *C.cipher__Address) *cipher.Address {
	return (*cipher.Address)(unsafe.Pointer(p))
}

func inplaceHttpHelperAddress(p *C.httphelper__Address) *httphelper.Address {
	return (*httphelper.Address)(unsafe.Pointer(p))
}

/**
 * Copy helpers
 */

func copyString(src string, dest *C.GoString_) {
	srcLen := len(src)
	dest.p = (*C.char)(C.malloc(C.size_t(srcLen + 1)))
	strAddr := (*C.GoString_)(unsafe.Pointer(&src))
	result := C.memcpy(unsafe.Pointer(dest.p), unsafe.Pointer(strAddr.p), C.size_t(srcLen))
	if result != nil {
		C.eos(dest.p, C.int(srcLen))
		dest.n = C.GoInt_(srcLen)
	}
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
	if C.memcpy(dest, srcAddr, C.size_t(n)*elemSize) != nil {
	}
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
	result := C.memcpy(dest.data, srcAddr, C.size_t(n)*elemSize)
	if result != nil {
		// Do not modify slice metadata until memory is actually copied
		if overflow {
			dest.len = dest.cap - C.GoInt_(srcLen)
		} else {
			dest.len = C.GoInt_(srcLen)
		}
	}
}

func convertToInterface(a *C.GoInterface_) interface{} {
	//TODO: Implement
	return nil
}

func copyToFunc(f C.Handle) func() hash.Hash {
	//TODO: Implement
	return nil
}

func copyToStringMap(gomap map[string]string, dest *C.GoStringMap_) {
	*dest = (C.GoStringMap_)(registerHandle(gomap))
}

func splitCliArgs(args string) (result []string) {
	prevSep := -1
	quoted := false
	var i int
	for i = 0; i < len(args); i++ {
		if args[i] == '"' {
			quoted = !quoted
			if !quoted {
				result = append(result, args[prevSep+1:i])
			}
			prevSep = i
		} else if !quoted && args[i] == ' ' {
			if prevSep+1 < i {
				result = append(result, args[prevSep+1:i])
			}
			prevSep = i
		}
	}
	if len(args) > 0 {
		if prevSep+1 < i {
			result = append(result, args[prevSep+1:i])
		}
	}
	return
}
