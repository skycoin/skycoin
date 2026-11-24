// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && (amd64 || arm64)) || (darwin && arm64)

package libc // import "modernc.org/libc"

import (
	"sync"
	"unsafe"
)

// The Y functions here are written by hand, not by 'qbecc --abi0wrap', because
// libc keeps the C function pointers passed to them after they return. Code
// calling the Y functions, like C compiled at run time, passes a C function
// pointer as the ABI0 entry address of the function. libc calls a C function
// pointer as a Go func value. A generated Y function converts such an address
// to a closure that lives until it returns, which is enough only for callbacks
// used during the call, like qsort's.
//
// The Y functions here convert the address to a permanent closure instead, one
// per address and signature, see abi0Callback. The same address therefore
// always becomes the same func value, which keeps C function pointers
// comparable, and abi0Entry maps it back for the function pointers libc
// returns, the old handlers of signal and sigaction. The values 0, SIG_IGN and
// SIG_ERR are never converted, see abi0Handler.
//
// The Y functions are in abi0_callbacks_$GOOS.go and
// abi0_callbacks_$GOOS_$GOARCH.s, the calls of ABI0 functions in
// abi0_calls_$GOARCH.s.

// abi0Sig is the C signature of a callback, as libc calls it.
type abi0Sig int

const (
	abi0SigV     abi0Sig = iota // void f(void)
	abi0SigI                    // void f(int)
	abi0SigIPP                  // void f(int, void*, void*)
	abi0SigP                    // void f(void*)
	abi0SigP_P                  // void *f(void*)
	abi0SigP_I                  // int f(void*)
	abi0SigPPS_S                // ssize_t f(void*, char*, size_t)
	abi0SigPPI_I                // int f(void*, off_t*, int)
)

type abi0Key struct {
	entry uintptr
	sig   abi0Sig
}

var (
	abi0CallbacksMu sync.Mutex
	abi0Callbacks   = map[abi0Key]any{}     // The closures, kept alive here.
	abi0Entries     = map[uintptr]uintptr{} // func value -> ABI0 entry
)

// abi0Callback returns the permanent func value that calls the ABI0 function
// at entry, which has the signature sig.
func abi0Callback(entry uintptr, sig abi0Sig) (r uintptr) {
	if entry == 0 {
		return 0
	}

	abi0CallbacksMu.Lock()

	defer abi0CallbacksMu.Unlock()

	k := abi0Key{entry, sig}
	if f, ok := abi0Callbacks[k]; ok {
		return abi0FuncValue(f)
	}

	var f any
	switch sig {
	case abi0SigV:
		f = func(tls *TLS) { abi0CallV(tls, entry) }
	case abi0SigI:
		f = func(tls *TLS, a int32) { abi0CallI(tls, a, entry) }
	case abi0SigIPP:
		f = func(tls *TLS, a int32, b, c uintptr) { abi0CallIPP(tls, a, b, c, entry) }
	case abi0SigP:
		f = func(tls *TLS, a uintptr) { abi0CallP(tls, a, entry) }
	case abi0SigP_P:
		f = func(tls *TLS, a uintptr) uintptr { return abi0CallP_P(tls, a, entry) }
	case abi0SigP_I:
		f = func(tls *TLS, a uintptr) int32 { return abi0CallP_I(tls, a, entry) }
	case abi0SigPPS_S:
		f = func(tls *TLS, a, b uintptr, c uint64) int64 { return abi0CallPPS_S(tls, a, b, c, entry) }
	case abi0SigPPI_I:
		f = func(tls *TLS, a, b uintptr, c int32) int32 { return abi0CallPPI_I(tls, a, b, c, entry) }
	default:
		panic(todo("", sig))
	}
	abi0Callbacks[k] = f
	r = abi0FuncValue(f)
	abi0Entries[r] = entry
	return r
}

// abi0FuncValue returns the func value in f, as a C function pointer, the way
// ccgo represents it.
func abi0FuncValue(f any) uintptr {
	return (*[2]uintptr)(unsafe.Pointer(&f))[1]
}

// abi0Entry returns the ABI0 entry address abi0Callback converted to fn. Any
// other fn, a sentinel or a func value not made by abi0Callback, is returned
// as is.
func abi0Entry(fn uintptr) uintptr {
	abi0CallbacksMu.Lock()

	defer abi0CallbacksMu.Unlock()

	if r, ok := abi0Entries[fn]; ok {
		return r
	}

	return fn
}

// The calls of an ABI0 function at entry, see abi0_calls_$GOARCH.s. size_t and
// ssize_t are uint64 and int64 on every target here.

//go:noescape
func abi0CallV(tls *TLS, entry uintptr)

//go:noescape
func abi0CallI(tls *TLS, a int32, entry uintptr)

//go:noescape
func abi0CallIPP(tls *TLS, a int32, b, c, entry uintptr)

//go:noescape
func abi0CallP(tls *TLS, a, entry uintptr)

//go:noescape
func abi0CallP_P(tls *TLS, a, entry uintptr) (r uintptr)

//go:noescape
func abi0CallP_I(tls *TLS, a, entry uintptr) (r int32)

//go:noescape
func abi0CallPPS_S(tls *TLS, a, b uintptr, c uint64, entry uintptr) (r int64)

//go:noescape
func abi0CallPPI_I(tls *TLS, a, b uintptr, c int32, entry uintptr) (r int32)
