// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
// SIG_ERR are never converted.

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
		f = func(tls *TLS, a, b uintptr, c Tsize_t) Tssize_t { return abi0CallPPS_S(tls, a, b, c, entry) }
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

// abi0Handler is abi0Callback for a signal handler, which may be one of the
// sentinels.
func abi0Handler(h uintptr, sig abi0Sig) uintptr {
	switch h {
	case SIG_DFL, SIG_IGN, SIG_ERR:
		return h
	}

	return abi0Callback(h, sig)
}

// The calls of an ABI0 function at entry, see abi0_callbacks_linux_amd64.s.

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
func abi0CallPPS_S(tls *TLS, a, b uintptr, c Tsize_t, entry uintptr) (r Tssize_t)

//go:noescape
func abi0CallPPI_I(tls *TLS, a, b uintptr, c int32, entry uintptr) (r int32)

// The hand-written Y functions. Each TEXT stub in
// abi0_callbacks_linux_amd64.s calls the lowercase Go function of the same
// name, which converts the C function pointers and calls X.

//go:noescape
func Yatexit(tls *TLS, func_ uintptr) (r int32)

func yatexit(tls *TLS, f uintptr) int32 {
	return Xatexit(tls, abi0Callback(f, abi0SigV))
}

//go:noescape
func Yat_quick_exit(tls *TLS, __ccgo_fp_func uintptr) (r1 int32)

func yat_quick_exit(tls *TLS, f uintptr) int32 {
	return Xat_quick_exit(tls, abi0Callback(f, abi0SigV))
}

//go:noescape
func Ysignal(tls *TLS, signum int32, handler uintptr) (r uintptr)

func ysignal(tls *TLS, signum int32, handler uintptr) uintptr {
	return abi0Entry(Xsignal(tls, signum, abi0Handler(handler, abi0SigI)))
}

//go:noescape
func Ysigaction(tls *TLS, sig int32, sa uintptr, old uintptr) (r int32)

// ysigaction converts the handler in a copy of *sa, so the caller's struct is
// left alone and may be *old as well.
func ysigaction(tls *TLS, sig int32, sa, old uintptr) (r int32) {
	if sa != 0 {
		n := int(unsafe.Sizeof(Tsigaction{}))
		p := tls.Alloc(n)

		defer tls.Free(n)

		*(*Tsigaction)(unsafe.Pointer(p)) = *(*Tsigaction)(unsafe.Pointer(sa))
		h := &(*Tsigaction)(unsafe.Pointer(p)).F__sa_handler.Fsa_handler // also sa_sigaction
		sig := abi0SigI
		if (*Tsigaction)(unsafe.Pointer(p)).Fsa_flags&SA_SIGINFO != 0 {
			sig = abi0SigIPP
		}
		*h = abi0Handler(*h, sig)
		sa = p
	}
	if r = Xsigaction(tls, sig, sa, old); r == 0 && old != 0 {
		h := &(*Tsigaction)(unsafe.Pointer(old)).F__sa_handler.Fsa_handler
		*h = abi0Entry(*h)
	}
	return r
}

//go:noescape
func Ypthread_create(tls *TLS, res, attrp, entry, arg uintptr) (_2 int32)

func ypthread_create(tls *TLS, res, attrp, entry, arg uintptr) int32 {
	return Xpthread_create(tls, res, attrp, abi0Callback(entry, abi0SigP_P), arg)
}

//go:noescape
func Ypthread_key_create(tls *TLS, k uintptr, dtor uintptr) (_3 int32)

func ypthread_key_create(tls *TLS, k, dtor uintptr) int32 {
	return Xpthread_key_create(tls, k, abi0Callback(dtor, abi0SigP))
}

//go:noescape
func Y_pthread_cleanup_push(tls *TLS, cb, f, x uintptr)

func y_pthread_cleanup_push(tls *TLS, cb, f, x uintptr) {
	X_pthread_cleanup_push(tls, cb, abi0Callback(f, abi0SigP), x)
}

//go:noescape
func Yfopencookie(tls *TLS, cookie uintptr, mode uintptr, iofuncs Tcookie_io_functions_t) (r uintptr)

func yfopencookie(tls *TLS, cookie, mode uintptr, iofuncs Tcookie_io_functions_t) uintptr {
	iofuncs.Fread = abi0Callback(iofuncs.Fread, abi0SigPPS_S)
	iofuncs.Fwrite = abi0Callback(iofuncs.Fwrite, abi0SigPPS_S)
	iofuncs.Fseek = abi0Callback(iofuncs.Fseek, abi0SigPPI_I)
	iofuncs.Fclose1 = abi0Callback(iofuncs.Fclose1, abi0SigP_I)
	return Xfopencookie(tls, cookie, mode, iofuncs)
}

//go:noescape
func Ypthread_atfork(tls *TLS, prepare, parent, child uintptr) (_2 int32)

func ypthread_atfork(tls *TLS, prepare, parent, child uintptr) int32 {
	return Xpthread_atfork(tls, abi0Callback(prepare, abi0SigV), abi0Callback(parent, abi0SigV), abi0Callback(child, abi0SigV))
}
