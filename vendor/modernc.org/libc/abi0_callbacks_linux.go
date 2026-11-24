// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && (amd64 || arm64)

package libc // import "modernc.org/libc"

import (
	"unsafe"
)

// abi0Handler is abi0Callback for a signal handler, which may be one of the
// sentinels.
func abi0Handler(h uintptr, sig abi0Sig) uintptr {
	switch h {
	case SIG_DFL, SIG_IGN, SIG_ERR:
		return h
	}

	return abi0Callback(h, sig)
}

// The hand-written Y functions, see abi0_callbacks.go. Each TEXT stub in
// abi0_callbacks_linux_$GOARCH.s calls the lowercase Go function of the same
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
