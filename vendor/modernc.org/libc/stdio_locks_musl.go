// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && (amd64 || arm64 || loong64 || ppc64le || s390x || riscv64 || 386 || arm)

package libc // import "modernc.org/libc"

import (
	"sync/atomic"
	"unsafe"
)

// vfprintfBufSize is the size of the temporary buffer musl's vfprintf puts in
// place of the buffer of an unbuffered stream while it runs.
const vfprintfBufSize = 80

// stdFiles are the standard streams with the buffers they start with.
var stdFiles = [...]struct {
	f   uintptr
	buf uintptr
}{
	{uintptr(unsafe.Pointer(&X__stdin_FILE)), X__stdin_FILE.Fbuf},
	{uintptr(unsafe.Pointer(&X__stdout_FILE)), X__stdout_FILE.Fbuf},
	{uintptr(unsafe.Pointer(&X__stderr_FILE)), X__stderr_FILE.Fbuf},
}

// ReleaseFileLocks releases the stdio FILE locks that the C thread of tls
// still holds after a panic unwound its frames, typically a fault in C code
// recovered under debug.SetPanicOnFault. A FILE a stdio call had locked when
// the panic struck stays locked, and every later use of it, including the
// flush at exit, waits forever.
//
// It releases both kinds of FILE lock: the one stdio functions take
// internally and the one of flockfile and ftrylockfile, whatever its count.
// It wakes the C threads waiting for them. It covers stdin, stdout, stderr and
// every FILE on the open file list, so every FILE that fopen, fdopen, popen,
// fopencookie, fmemopen or open_memstream returned and fclose has not closed.
// FILEs locked by other TLSs are left alone.
//
// A FILE's state is otherwise as the panic left it. Output buffered but not
// yet written is written later, or lost if the panic struck while a write was
// in progress. One state is repaired: printf and its siblings make an
// unbuffered stream, such as stderr, use a temporary buffer in their own frame
// while they run. If the panic struck then, ReleaseFileLocks gives the stream
// a buffer of its own again; writes would otherwise go to memory that is now
// reused.
//
// Call it on the goroutine that owns tls, after the panic is recovered and
// before tls.Close. No stdio call may be running on tls. It does not release
// other locks. A panic inside fopen, fclose or another function holding
// musl's lock of the open file list leaves that lock held, and
// ReleaseFileLocks then blocks.
func (tls *TLS) ReleaseFileLocks() {
	if tls.pthread == 0 {
		return
	}

	self := (*t__pthread)(unsafe.Pointer(tls.pthread))
	tid := self.Ftid
	// flockfile lists the FILEs it locks with the C thread. Release those
	// first, in case one is missing from the open file list.
	for f := self.Fstdio_locks; f != 0; {
		next := (*TFILE)(unsafe.Pointer(f)).Fnext_locked
		tls.releaseFile(tid, f)
		f = next
	}
	self.Fstdio_locks = 0
	for _, v := range &stdFiles {
		tls.releaseFile(tid, v.f)
	}
	head := X__ofl_lock(tls)
	for f := *(*uintptr)(unsafe.Pointer(head)); f != 0; f = (*TFILE)(unsafe.Pointer(f)).Fnext {
		tls.releaseFile(tid, f)
	}
	X__ofl_unlock(tls)
}

// releaseFile releases f if the C thread tid holds its lock, see
// ReleaseFileLocks.
func (tls *TLS) releaseFile(tid int32, f uintptr) {
	p := (*TFILE)(unsafe.Pointer(f))
	if atomic.LoadInt32(&p.Flock)&^int32(MAYBE_WAITERS) != tid {
		return
	}

	if p.Fbuf_size == vfprintfBufSize && tls.isStackSlot(p.Fbuf) {
		// vfprintf's temporary buffer, a TLS.Alloc slot the panic released.
		// Put back the stream's own buffer if it is a standard stream, else
		// allocate one, which only has to hold what ungetc pushes back.
		buf := uintptr(0)
		for _, v := range &stdFiles {
			if v.f == f {
				buf = v.buf
			}
		}
		if buf == 0 {
			buf = mustMalloc(UNGET) + UNGET
		}
		p.Fbuf, p.Fbuf_size = buf, 0
		p.Fwpos, p.Fwbase, p.Fwend = 0, 0, 0
	}
	p.Flockcount = 0
	p.Fprev_locked, p.Fnext_locked = 0, 0
	___unlockfile(tls, f)
}

// isStackSlot reports whether p is the start of one of the TLS.Alloc slots of
// tls.
func (tls *TLS) isStackSlot(p uintptr) bool {
	for _, v := range tls.stack {
		if v.p == p {
			return true
		}
	}
	return false
}
