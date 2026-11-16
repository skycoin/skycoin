// Copyright 2023 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !libc.membrk && !libc.memgrind && linux && (amd64 || arm64 || loong64 || ppc64le || s390x || riscv64 || 386 || arm)

package libc // import "modernc.org/libc"

import (
	"math"
	mbits "math/bits"
	"unsafe"

	"modernc.org/memory"
)

const (
	// allocaBump selects the alloca(3) bump region, see allocaChunk.
	allocaBump = true
	isMemBrk   = false
)

// Each TLS caches free slots of up to tcacheMaxSize bytes, so that most mallocs
// and frees of small blocks take no lock. A slot freed through a TLS goes to
// its cache, whichever TLS allocated it. A cache of slots of size s holds at
// most tcacheBinBytes/s of them. When it is empty, malloc takes half that many
// from the allocator at once, and when it is full, free gives back half. Close
// gives back all. Slots in caches count as allocated in MemStat.
//
// A TLS is not safe for concurrent use, so its cache needs no lock. Callers
// with no TLS of their own pass nil and use the allocator directly.
const (
	tcacheBinBytes = 4 << 10
	tcacheMaxLog   = 10
	tcacheMask     = uint(len(TLS{}.tcache) - 1)
	tcacheMaxSize  = 1 << tcacheMaxLog
	mallocAlign    = 2 * unsafe.Sizeof(uintptr(0)) // as in modernc.org/memory
)

// tcacheLog returns the log of the size of the slot malloc uses for 0 < n <=
// tcacheMaxSize bytes, as modernc.org/memory picks it.
func tcacheLog(n Tsize_t) uint {
	return uint(mbits.Len(uint((uintptr(n)+mallocAlign-1)&^(mallocAlign-1) - 1)))
}

// tcacheMalloc returns a slot for 0 < n <= tcacheMaxSize bytes, or zero with
// errno set to ENOMEM.
func (tls *TLS) tcacheMalloc(n Tsize_t) (r uintptr) {
	log := tcacheLog(n)
	b := &tls.tcache[log&tcacheMask]
	if r = b.head; r != 0 {
		b.head = *(*uintptr)(unsafe.Pointer(r))
		b.n--
		return r
	}

	if r = tls.tcacheRefill(b, log); r == 0 {
		tls.setErrno(ENOMEM)
	}
	return r
}

// tcacheRefill returns a new slot of 1<<log bytes and puts up to half the
// capacity of b more into b, which is empty. It returns zero if the allocator
// is out of memory.
func (tls *TLS) tcacheRefill(b *tcacheBin, log uint) (r uintptr) {
	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	r, err := allocator.UintptrMalloc(1 << log)
	if err != nil {
		return 0
	}

	// Link the slots in the order the allocator hands them out, which is
	// the order they then come out of b.
	var tail uintptr
	for n := tcacheBinBytes >> (log & 63) / 2; b.n < n; b.n++ {
		p, err := allocator.UintptrMalloc(1 << log)
		if err != nil {
			break
		}

		if tail == 0 {
			b.head = p
		} else {
			*(*uintptr)(unsafe.Pointer(tail)) = p
		}
		tail = p
	}
	if tail != 0 {
		*(*uintptr)(unsafe.Pointer(tail)) = 0
	}
	return r
}

// tcacheFree puts p, a slot of us bytes, into the cache of tls if us is a power
// of two not above tcacheMaxSize and reports whether it did.
func (tls *TLS) tcacheFree(p uintptr, us int) bool {
	if us > tcacheMaxSize || us&(us-1) != 0 {
		return false
	}

	log := uint(mbits.TrailingZeros(uint(us)))
	b := &tls.tcache[log&tcacheMask]
	if b.n == tcacheBinBytes>>(log&63) {
		tls.tcacheFlush(b, b.n/2)
	}
	*(*uintptr)(unsafe.Pointer(p)) = b.head
	b.head = p
	b.n++
	return true
}

// tcacheFlush gives n slots from b back to the allocator.
func (tls *TLS) tcacheFlush(b *tcacheBin, n int) {
	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	for ; n != 0; n-- {
		p := b.head
		b.head = *(*uintptr)(unsafe.Pointer(p))
		b.n--
		allocator.UintptrFree(p)
	}
}

// tcacheFlushAll gives every slot in the cache of tls back to the allocator.
func (tls *TLS) tcacheFlushAll() {
	for i := range tls.tcache {
		if b := &tls.tcache[i]; b.n != 0 {
			tls.tcacheFlush(b, b.n)
		}
	}
}

func Xmalloc(tls *TLS, n Tsize_t) (r uintptr) {
	if __ccgo_strace {
		trc("tls=%v n=%v, (%v:)", tls, n, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	if n == 0 {
		// malloc(0) should return unique pointers
		// (often expected and gnulib replaces malloc if malloc(0) returns 0)
		n = 1
	}
	if n <= tcacheMaxSize && tls != nil {
		return tls.tcacheMalloc(n)
	}

	if n > math.MaxInt {
		tls.setErrno(ENOMEM)
		return 0
	}

	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	var err error
	if r, err = allocator.UintptrMalloc(int(n)); err != nil {
		r = 0
		tls.setErrno(ENOMEM)
	}
	return r
}

func Xcalloc(tls *TLS, m Tsize_t, n Tsize_t) (r uintptr) {
	if __ccgo_strace {
		trc("tls=%v m=%v n=%v, (%v:)", tls, m, n, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	hi, rq := mbits.Mul(uint(m), uint(n))
	if hi != 0 || rq > math.MaxInt {
		tls.setErrno(ENOMEM)
		return 0
	}

	if rq == 0 {
		rq = 1
	}
	if rq <= tcacheMaxSize && tls != nil {
		if r = tls.tcacheMalloc(Tsize_t(rq)); r != 0 {
			clear((*RawMem)(unsafe.Pointer(r))[:rq:rq])
		}
		return r
	}

	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	var err error
	if r, err = allocator.UintptrCalloc(int(rq)); err != nil {
		r = 0
		tls.setErrno(ENOMEM)
	}
	return r
}

func Xrealloc(tls *TLS, p uintptr, n Tsize_t) (r uintptr) {
	if __ccgo_strace {
		trc("tls=%v p=%v n=%v, (%v:)", tls, p, n, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	var err error
	if r, err = allocator.UintptrRealloc(p, int(n)); err != nil {
		r = 0
		tls.setErrno(ENOMEM)
	}
	return r
}

func Xfree(tls *TLS, p uintptr) {
	if __ccgo_strace {
		trc("tls=%v p=%v, (%v:)", tls, p, origin(2))
	}
	if p == 0 || tls != nil && tls.tcacheFree(p, memory.UintptrUsableSize(p)) {
		return
	}

	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	allocator.UintptrFree(p)
}

func Xmalloc_usable_size(tls *TLS, p uintptr) (r Tsize_t) {
	if __ccgo_strace {
		trc("tls=%v p=%v, (%v:)", tls, p, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	if p == 0 {
		return 0
	}

	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	return Tsize_t(memory.UintptrUsableSize(p))
}

func MemAudit() (r []*MemAuditError) {
	return nil
}

func UsableSize(p uintptr) Tsize_t {
	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	return Tsize_t(memory.UintptrUsableSize(p))
}

type MemAllocatorStat struct {
	Allocs int
	Bytes  int
	Mmaps  int
}

// MemStat returns the global memory allocator statistics.
// should be compiled with the memory.counters build tag for the data to be available.
func MemStat() MemAllocatorStat {
	allocatorMu.Lock()
	defer allocatorMu.Unlock()

	return MemAllocatorStat{
		Allocs: allocator.Allocs,
		Bytes:  allocator.Bytes,
		Mmaps:  allocator.Mmaps,
	}
}

// MemAuditStart locks the memory allocator, initializes and enables memory
// auditing. Finaly it unlocks the memory allocator.
//
// Some memory handling errors, like double free or freeing of unallocated
// memory, will panic when memory auditing is enabled.
//
// This memory auditing functionality has to be enabled using the libc.memgrind
// build tag.
//
// It is intended only for debug/test builds. It slows down memory allocation
// routines and it has additional memory costs.
func MemAuditStart() {}

// MemAuditReport locks the memory allocator, reports memory leaks, if any.
// Finally it disables memory auditing and unlocks the memory allocator.
//
// This memory auditing functionality has to be enabled using the libc.memgrind
// build tag.
//
// It is intended only for debug/test builds. It slows down memory allocation
// routines and it has additional memory costs.
func MemAuditReport() error { return nil }
