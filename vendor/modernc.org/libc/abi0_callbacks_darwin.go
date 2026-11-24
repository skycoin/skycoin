// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm64

package libc // import "modernc.org/libc"

import (
	"unsafe"

	"modernc.org/libc/signal"
)

// The types of the C function pointer parameters of qsort, bsearch and ftw, in
// the form ccgo gives them in generated code. qbecc --abi0wrap recognizes a C
// function pointer parameter only by such an alias, named after the function
// and the index of the parameter after the TLS. With it, the generated Y
// wrapper converts the ABI0 code address its caller passes to a func value for
// the duration of the call. Without it, libc would call the code address as a
// func value and crash. The X functions written by hand on darwin have no
// aliases of their own. fts_open ignores its compar, and the functions that
// keep a C function pointer have hand-written Y functions below.
type (
	t__ccgo_fp__Xbsearch_4 = func(*TLS, uintptr, uintptr) int32        // int (*compar)(const void *, const void *)
	t__ccgo_fp__Xftw_1     = func(*TLS, uintptr, uintptr, int32) int32 // int (*fn)(const char *, const struct stat *, int)
	t__ccgo_fp__Xqsort_3   = func(*TLS, uintptr, uintptr) int32        // int (*compar)(const void *, const void *)
)

// abi0Handler is abi0Callback for a signal handler, which may be one of the
// sentinels SIG_DFL, SIG_IGN or SIG_ERR, ((void (*)(int))-1).
func abi0Handler(h uintptr, sig abi0Sig) uintptr {
	switch h {
	case signal.SIG_DFL, signal.SIG_IGN, ^uintptr(0):
		return h
	}

	return abi0Callback(h, sig)
}

// The hand-written Y functions, see abi0_callbacks.go. Each TEXT stub in
// abi0_callbacks_darwin_arm64.s calls the lowercase Go function of the same
// name, which converts the C function pointers and calls X. darwin has no
// at_quick_exit, _pthread_cleanup_push or fopencookie.

//go:noescape
func Yatexit(tls *TLS, function uintptr) (r int32)

func yatexit(tls *TLS, f uintptr) int32 {
	return Xatexit(tls, abi0Callback(f, abi0SigV))
}

//go:noescape
func Ysignal(tls *TLS, signum int32, handler uintptr) (r uintptr)

func ysignal(tls *TLS, signum int32, handler uintptr) uintptr {
	return abi0Entry(Xsignal(tls, signum, abi0Handler(handler, abi0SigI)))
}

//go:noescape
func Ysigaction(tls *TLS, signum int32, act, oldact uintptr) (r int32)

// ysigaction converts the handler in a copy of *act, so the caller's struct is
// left alone and may be *oldact as well.
func ysigaction(tls *TLS, signum int32, act, oldact uintptr) (r int32) {
	if act != 0 {
		n := int(unsafe.Sizeof(signal.Sigaction{}))
		p := tls.Alloc(n)

		defer tls.Free(n)

		*(*signal.Sigaction)(unsafe.Pointer(p)) = *(*signal.Sigaction)(unsafe.Pointer(act))
		h := &(*signal.Sigaction)(unsafe.Pointer(p)).F__sigaction_u.F__sa_handler // also sa_sigaction
		sig := abi0SigI
		if (*signal.Sigaction)(unsafe.Pointer(p)).Fsa_flags&signal.SA_SIGINFO != 0 {
			sig = abi0SigIPP
		}
		*h = abi0Handler(*h, sig)
		act = p
	}
	if r = Xsigaction(tls, signum, act, oldact); r == 0 && oldact != 0 {
		h := &(*signal.Sigaction)(unsafe.Pointer(oldact)).F__sigaction_u.F__sa_handler
		*h = abi0Entry(*h)
	}
	return r
}

//go:noescape
func Ypthread_create(tls *TLS, pThread, pAttr, startRoutine, arg uintptr) (r int32)

func ypthread_create(tls *TLS, pThread, pAttr, startRoutine, arg uintptr) int32 {
	return Xpthread_create(tls, pThread, pAttr, abi0Callback(startRoutine, abi0SigP_P), arg)
}

//go:noescape
func Ypthread_key_create(tls *TLS, pKey, destructor uintptr) (r int32)

func ypthread_key_create(tls *TLS, pKey, destructor uintptr) int32 {
	return Xpthread_key_create(tls, pKey, abi0Callback(destructor, abi0SigP))
}

//go:noescape
func Ypthread_atfork(tls *TLS, prepare, parent, child uintptr) (r int32)

func ypthread_atfork(tls *TLS, prepare, parent, child uintptr) int32 {
	return Xpthread_atfork(tls, abi0Callback(prepare, abi0SigV), abi0Callback(parent, abi0SigV), abi0Callback(child, abi0SigV))
}
