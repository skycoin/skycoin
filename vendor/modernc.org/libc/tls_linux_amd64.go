// Copyright 2025 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libc // import "modernc.org/libc"

import (
	"unsafe"
)

// The TLSAlloc fast path in tls_linux_amd64.s indexes tls.stack by shifting
// sp left by 4.
var _ = [1]struct{}{}[unsafe.Sizeof(tlsStackSlot{})-16]

// TLSAlloc is (*TLS).Alloc for ABI0 callers. The slot reuse fast path is
// inlined in assembly, anything else calls tlsAlloc.
//
//go:noescape
func TLSAlloc(p0 *TLS, p1 int) uintptr

// TLSFree is (*TLS).Free for ABI0 callers. Unless signals must be checked it
// is inlined in assembly, otherwise it calls tlsFree.
//
//go:noescape
func TLSFree(p0 *TLS, p1 int)

//go:noescape
func TLSAllocaEntry(p0 *TLS)

//go:noescape
func TLSAllocaExit(p0 *TLS)

func tlsAlloc(tls *TLS, n int) uintptr {
	return tls.Alloc(n)
}

func tlsFree(tls *TLS, n int) {
	tls.Free(n)
}

func tlsAllocaEntry(tls *TLS) {
	tls.AllocaEntry()
}

func tlsAllocaExit(tls *TLS) {
	tls.AllocaExit()
}
