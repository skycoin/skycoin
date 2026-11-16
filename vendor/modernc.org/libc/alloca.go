// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libc // import "modernc.org/libc"

// alloca(3) memory comes from a per-TLS bump region made of allocaChunkSize
// chunks, which are kept for reuse until Close. AllocaEntry saves the region
// position in an allocaMark and AllocaExit returns to it, so the memory lives
// as long as it would if it were malloc'd and freed by AllocaExit. Requests
// over allocaMaxBump, and all requests in the libc.membrk and libc.memgrind
// builds, are malloc'd, recorded in TLS.allocas and freed by AllocaExit.
// Memory from the bump region is allocaAlign aligned, malloc'd memory has the
// alignment of malloc.
//
// The code is shared by the musl targets and the hand-written ports, whose
// TLS types both have the fields used here.
const (
	allocaAlign     = 16
	allocaChunkSize = 64 << 10
	allocaMaxBump   = 16 << 10
)

type allocaChunk struct {
	mem uintptr // as returned by malloc
	p   uintptr // mem aligned up to allocaAlign
}

type allocaMark struct {
	n    int // len(TLS.allocas)
	next int // TLS.allocaNext
	p    uintptr
	end  uintptr
	seq  uintptr // TLS.allocaSeq of the region, see TLSCheckpoint
}

// alloca implements alloca(3), see allocaChunk.
func (tls *TLS) alloca(n Tsize_t) (r uintptr) {
	if allocaBump && n-1 < allocaMaxBump { // 0 < n <= allocaMaxBump
		sz := (uintptr(n) + allocaAlign - 1) &^ (allocaAlign - 1)
		if r = tls.allocaP; sz <= tls.allocaEnd-r {
			tls.allocaP = r + sz
			return r
		}
	}

	return tls.allocaSlow(n)
}

func (tls *TLS) allocaSlow(n Tsize_t) (r uintptr) {
	if !allocaBump || n > allocaMaxBump {
		r = tls.allocaMalloc(n)
		tls.allocas = append(tls.allocas, r)
		return r
	}

	sz := (uintptr(n) + allocaAlign - 1) &^ (allocaAlign - 1)
	if sz == 0 {
		sz = allocaAlign // Unique pointers, as malloc(0) returns.
	}
	if r = tls.allocaP; sz <= tls.allocaEnd-r {
		tls.allocaP = r + sz
		return r
	}

	// Move to the next chunk, allocating it if needed.
	if tls.allocaNext == len(tls.allocaChunks) {
		mem := tls.allocaMalloc(allocaChunkSize + allocaAlign - 1)
		p := (mem + allocaAlign - 1) &^ (allocaAlign - 1)
		tls.allocaChunks = append(tls.allocaChunks, allocaChunk{mem: mem, p: p})
		if len(tls.allocaChunks) == 1 {
			// Marks saved before the first chunk existed now mean its start,
			// so that returning to them does not leave the bump region and
			// every outermost region does not start in allocaSlow.
			for i := range tls.allocaStack {
				if m := &tls.allocaStack[i]; m.end == 0 {
					m.next, m.p, m.end = 1, p, p+allocaChunkSize
				}
			}
		}
	}
	r = tls.allocaChunks[tls.allocaNext].p
	tls.allocaNext++
	tls.allocaP = r + sz
	tls.allocaEnd = r + allocaChunkSize
	return r
}

func (tls *TLS) allocaMalloc(n Tsize_t) (r uintptr) {
	if r = Xmalloc(tls, n); r == 0 {
		panic(todo("OOM"))
	}

	return r
}

// AllocaEntry must be called early on function entry when the function calls
// or may call alloca(3).
func (tls *TLS) AllocaEntry() {
	tls.allocaSeq++
	tls.allocaStack = append(tls.allocaStack, allocaMark{len(tls.allocas), tls.allocaNext, tls.allocaP, tls.allocaEnd, tls.allocaSeq})
}

// AllocaExit must be defer-called on function exit when the function calls or
// may call alloca(3).
func (tls *TLS) AllocaExit() {
	n := len(tls.allocaStack) - 1
	m := &tls.allocaStack[n]
	if len(tls.allocas) != m.n {
		tls.allocaFree(m.n)
	}
	tls.allocaNext, tls.allocaP, tls.allocaEnd = m.next, m.p, m.end
	tls.allocaStack = tls.allocaStack[:n]
}

// allocaFree frees the malloc'd allocas from index n on.
func (tls *TLS) allocaFree(n int) {
	for _, v := range tls.allocas[n:] {
		Xfree(tls, v)
	}
	tls.allocas = tls.allocas[:n]
}

// allocaClose frees every alloca and every chunk of the bump region, whether
// or not their regions were exited, and empties the region.
func (tls *TLS) allocaClose() {
	tls.allocaFree(0)
	for _, v := range tls.allocaChunks {
		Xfree(tls, v.mem)
	}
	tls.allocaStack = nil
	tls.allocas = nil
	tls.allocaChunks = nil
	tls.allocaNext, tls.allocaP, tls.allocaEnd = 0, 0, 0
}
