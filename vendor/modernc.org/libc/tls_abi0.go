// Copyright 2025 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && (amd64 || arm64)) || (darwin && arm64)

package libc // import "modernc.org/libc"

// The ABI0 entry points of the TLS, for ABI0 code like qbe's goabi0 output, in
// tls_$GOOS_$GOARCH.s. On linux the assembly inlines the fast paths of Alloc
// and Free, on darwin it only calls the Go functions below.

// TLSAlloc is (*TLS).Alloc for ABI0 callers. On linux the slot reuse fast path
// is inlined in assembly, anything else calls tlsAlloc.
//
//go:noescape
func TLSAlloc(p0 *TLS, p1 int) uintptr

// TLSFree is (*TLS).Free for ABI0 callers. On linux, unless signals must be
// checked, it is inlined in assembly, otherwise it calls tlsFree.
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
