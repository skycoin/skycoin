// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && (amd64 || arm64 || loong64 || ppc64le || s390x || riscv64 || 386 || arm)

package libc // import "modernc.org/libc"

import (
	"fmt"
	"unsafe"
)

// ExitError is what exit, _Exit, _exit, quick_exit and abort panic with when
// they are called with a TLS that unwinds exits, see SetExitUnwind.
type ExitError struct {
	// Status is the status passed to exit, _Exit, _exit or quick_exit, or
	// 128+SIGABRT for abort. A process ending with it would report the low
	// 8 bits.
	Status int32
	// Abort reports whether abort was called.
	Abort bool
}

// Error implements error.
func (e ExitError) Error() string {
	if e.Abort {
		return "abort"
	}

	return fmt.Sprintf("exit status %d", e.Status)
}

// SetExitUnwind sets whether a C program ending through tls panics back to its
// Go caller instead of ending the process. It is meant for a host that runs C
// programs in its own process and recovers the panic. While it is on:
//
//   - atexit and at_quick_exit called with tls register the handler with tls,
//     not with the process, so programs on different TLSs never run each
//     other's handlers.
//   - exit runs the atexit handlers registered with tls, newest first, flushes
//     every stdio stream as fflush(NULL) does and panics with
//     ExitError{Status: status}. The handlers registered with the process
//     and the functions registered by AtExit belong to the host and do not
//     run. Unlike a real exit it leaves the streams usable.
//   - quick_exit runs the at_quick_exit handlers registered with tls, newest
//     first, and panics in the same way, without flushing.
//   - _Exit and _exit panic in the same way at once.
//   - abort panics with ExitError{Status: 128 + SIGABRT, Abort: true} without
//     raising SIGABRT, so a handler the program installed for it does not
//     run.
//
// A handler is removed before it runs, so one that calls exit or quick_exit
// itself goes on with the remaining ones and the panic carries the status of
// the innermost call. Returning from a C main is equivalent to calling exit,
// so the host calls Xexit with main's result to end the program the same way.
//
// Nothing else is undone. Memory, open files, stdio buffers, the environment,
// the locale, signal dispositions, threads and the libc state of the process
// stay as the program left them, and so do locks held where the panic
// started. Exits through other TLSs, for example from other threads of the
// program, still end the process.
//
// No code in libc recovers panics. Code between the C program and the host that
// does must re-panic an ExitError.
func (tls *TLS) SetExitUnwind(on bool) {
	tls.exitUnwind = on
}

// PromoteExitHandlers moves the atexit and at_quick_exit handlers registered
// with tls under SetExitUnwind to the process, in registration order, as if
// the C program had registered them with the process at the time of the call.
// So they run when the process exits, newest first among its handlers, and
// before the functions registered by AtExit. A host that runs a C program's
// constructors under SetExitUnwind uses it to keep the handlers they register
// once the constructors have succeeded.
//
// The process has room for 32 at_quick_exit handlers, as in musl. Handlers
// that do not fit are dropped and reported in the error, the others are
// promoted. Nothing is left on tls, so a second call has nothing to do.
func (tls *TLS) PromoteExitHandlers() error {
	on := tls.exitUnwind
	tls.exitUnwind = false

	defer func() { tls.exitUnwind = on }()

	for _, f := range tls.atExit {
		Xatexit(tls, f)
	}
	tls.atExit = nil
	dropped := 0
	for _, f := range tls.atQuickExit {
		if Xat_quick_exit(tls, f) != 0 {
			dropped++
		}
	}
	tls.atQuickExit = nil
	if dropped != 0 {
		return fmt.Errorf("PromoteExitHandlers: %d at_quick_exit handlers dropped, the process has room for 32", dropped)
	}

	return nil
}

// runExitHandlers calls the C functions in *list, newest first, removing each
// before calling it.
func (tls *TLS) runExitHandlers(list *[]uintptr) {
	for n := len(*list); n != 0; n = len(*list) {
		f := (*list)[n-1]
		*list = (*list)[:n-1]
		(*(*func(*TLS))(unsafe.Pointer(&struct{ uintptr }{f})))(tls)
	}
}

// exitUnwound is exit under SetExitUnwind.
func (tls *TLS) exitUnwound(code int32) {
	tls.runExitHandlers(&tls.atExit)
	Xfflush(tls, 0)
	panic(ExitError{Status: code})
}

// ___exit_unwind is called first by _Exit, see
// internal/overlay/musl/src/exit/_Exit.c.
func ___exit_unwind(tls *TLS, code int32) {
	if tls != nil && tls.exitUnwind {
		panic(ExitError{Status: code})
	}
}

// ___quick_exit_unwind is called first by quick_exit, see
// internal/overlay/musl/src/exit/quick_exit.c.
func ___quick_exit_unwind(tls *TLS, code int32) {
	if tls != nil && tls.exitUnwind {
		tls.runExitHandlers(&tls.atQuickExit)
		panic(ExitError{Status: code})
	}
}

// ___at_quick_exit_unwind is called first by at_quick_exit, see
// internal/overlay/musl/src/exit/at_quick_exit.c.
func ___at_quick_exit_unwind(tls *TLS, f uintptr) int32 {
	if tls != nil && tls.exitUnwind {
		tls.atQuickExit = append(tls.atQuickExit, f)
		return 1
	}

	return 0
}
