// Copyright 2023 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !libc.membrk && !libc.memgrind && linux && (amd64 || arm64 || loong64 || ppc64le || s390x || riscv64 || 386 || arm)

package libc // import "modernc.org/libc"

import (
	"math"
	mbits "math/bits"
	"sync"
	"sync/atomic"
	"time"
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
// its cache, whichever TLS allocated it. The cache has a bin for each size
// class of modernc.org/memory up to tcacheMaxSize, see tcacheClasses. A bin of
// slots of size s holds at most tcacheBinBytes/s of them, but at least
// tcacheMinSlots. When it is empty, malloc takes up to half that many at once,
// and when it is full, free gives back half. Close gives back all. Slots in
// caches count as allocated in MemStat, slots in the depot do not.
//
// A TLS is not safe for concurrent use, so its cache needs no lock. Callers
// with no TLS of their own pass nil and use the depot.
//
// Between the caches and the allocator sits the depot: for each bin, a stack
// of lists of free slots under its own lock. A cache gives slots back by
// pushing them as one list and refills by popping one, each in constant time
// under the lock, so that short-lived TLSs pass memory to each other without
// the allocator. The depot holds at most depotMaxBytes in all. Past that,
// slots go back to the allocator.
//
// A slot in the depot keeps its page from emptying, so the depot gives back
// what it does not use. While it holds slots, depotTrim runs every
// depotTrimInterval and frees, in each bin, as many lists as the bin held at
// the least since the last trim: lists that no cache needed for a whole
// interval. A depot no one uses is empty after two intervals. MemTrim empties
// it at once.
const (
	tcacheBinBytes    = 4 << 10
	tcacheMinSlots    = 2
	tcacheBins        = len(TLS{}.tcache)
	tcacheMaxSize     = 8 << 10                       // includes the 5 KiB slots of SQLite's 4 KiB pages
	mallocAlign       = 2 * unsafe.Sizeof(uintptr(0)) // as in modernc.org/memory
	depotMaxBytes     = 16 << 20
	depotTrimInterval = time.Second
)

// tcacheClasses maps requests to bins. A request for 0 < n <= tcacheMaxSize
// bytes goes to bin i = bin[(n-1)/mallocAlign], which caches slots of size[i]
// bytes, at most max[i] of them. The sizes come from memory.UsableSizeFor, so
// size[i] is the size class modernc.org/memory serves the request from, and
// a slot whose usable size is size[i] goes back to bin i when freed.
//
// tcacheClasses is a variable with an initializer, not one filled by init, so
// that it is ready before any variable initializer of this package that calls
// malloc runs.
var tcacheClasses = newTcacheClasses()

type tcacheClassTable struct {
	bin  [tcacheMaxSize / mallocAlign]uint8
	size [tcacheBins]uint32
	max  [tcacheBins]uint32
}

func newTcacheClasses() (t tcacheClassTable) {
	n := 0
	for i := range t.bin {
		rq := (i + 1) * int(mallocAlign)
		s := memory.UsableSizeFor(rq)
		if s < rq || s > tcacheMaxSize {
			panic("internal error: unexpected size class")
		}

		if n == 0 || s != int(t.size[n-1]) {
			if n == tcacheBins {
				panic("internal error: more size classes than tcache bins")
			}

			t.size[n] = uint32(s)
			t.max[n] = tcacheMinSlots
			if m := tcacheBinBytes / s; m > tcacheMinSlots {
				t.max[n] = uint32(m)
			}
			n++
		}
		t.bin[i] = uint8(n - 1)
	}
	return t
}

// The masks in tcacheBinFor and on bin indexes are no-ops on valid indexes and
// spare the bounds checks.
const (
	tcacheBinMask   = tcacheBins - 1
	tcacheQuantMask = uintptr(len(tcacheClassTable{}.bin) - 1)
)

// tcacheBinFor returns the bin of a request for 0 < n <= tcacheMaxSize bytes.
func tcacheBinFor(n uintptr) int {
	return int(tcacheClasses.bin[(n-1)/mallocAlign&tcacheQuantMask]) & tcacheBinMask
}

// tcacheBinOf returns the bin that caches slots of us bytes, or -1 if none
// does. Every usable size up to tcacheMaxSize has one, the check guards the
// caches against a modernc.org/memory for which that does not hold.
func tcacheBinOf(us int) int {
	if us > tcacheMaxSize {
		return -1
	}

	i := tcacheBinFor(uintptr(us))
	if int(tcacheClasses.size[i]) != us {
		return -1
	}

	return i
}

var (
	// depotBytes is the size of the slots in the depot.
	depotBytes atomic.Int64

	// depotSlots is the number of slots in the depot, including those a
	// trim has taken out and not yet freed, which MemStat does not count
	// as allocated. It is kept only under the memory.counters build tag,
	// see memCounters.
	depotSlots atomic.Int64

	// depotTrimMu serializes trims. Tests hold it to keep the depot still.
	depotTrimMu sync.Mutex

	// depotTrimArmed is set while a trim is scheduled.
	depotTrimArmed atomic.Bool
)

// depotBin holds the free slots of one bin that the depot keeps.
type depotBin struct {
	mu    sync.Mutex
	lists []tcacheBin
	low   int // the fewest lists held since the last trim

	// The padding keeps the locks of two bins off one cache line.
	_ [64 - unsafe.Sizeof(sync.Mutex{}) - unsafe.Sizeof([]tcacheBin{}) - unsafe.Sizeof(int(0))]byte
}

var depot [tcacheBins]depotBin

// depotPut adds the list head...tail of n > 0 slots of bin i to the depot and
// reports whether it did. It does not when the depot is full.
func depotPut(i int, head, tail uintptr, n int) bool {
	d := &depot[i&tcacheBinMask]
	d.mu.Lock()

	defer d.mu.Unlock()

	size := int64(n) * int64(tcacheClasses.size[i&tcacheBinMask])
	if depotBytes.Add(size) > depotMaxBytes {
		depotBytes.Add(-size)
		return false
	}

	if memCounters {
		depotSlots.Add(int64(n))
	}
	depotArm()

	// Single slots from callers with no TLS collect into lists of the size
	// a refill takes.
	if k := len(d.lists); k != 0 && d.lists[k-1].n+n <= int(tcacheClasses.max[i&tcacheBinMask])/2 {
		l := &d.lists[k-1]
		*(*uintptr)(unsafe.Pointer(tail)) = l.head
		l.head = head
		l.n += n
		return true
	}

	*(*uintptr)(unsafe.Pointer(tail)) = 0
	d.lists = append(d.lists, tcacheBin{head, n})
	return true
}

// depotGet removes a list of slots of bin i from the depot and returns it. The
// list is empty if the depot has none.
func depotGet(i int) (r tcacheBin) {
	d := &depot[i&tcacheBinMask]
	d.mu.Lock()

	defer d.mu.Unlock()

	if k := len(d.lists); k != 0 {
		r = d.lists[k-1]
		d.lists = d.lists[:k-1]
		d.low = min(d.low, k-1)
		depotBytes.Add(-int64(r.n) * int64(tcacheClasses.size[i&tcacheBinMask]))
		if memCounters {
			depotSlots.Add(-int64(r.n))
		}
	}
	return r
}

// depotGetOne removes a slot of bin i from the depot and returns it, or returns
// zero if the depot has none.
func depotGetOne(i int) (r uintptr) {
	d := &depot[i&tcacheBinMask]
	d.mu.Lock()

	defer d.mu.Unlock()

	k := len(d.lists)
	if k == 0 {
		return 0
	}

	l := &d.lists[k-1]
	r = l.head
	l.head = *(*uintptr)(unsafe.Pointer(r))
	depotBytes.Add(-int64(tcacheClasses.size[i&tcacheBinMask]))
	if memCounters {
		depotSlots.Add(-1)
	}
	if l.n--; l.n == 0 {
		d.lists = d.lists[:k-1]
		d.low = min(d.low, k-1)
	}
	return r
}

// depotArm schedules a trim unless one is scheduled.
//
// A timer made inside a testing/synctest bubble belongs to the bubble and never
// fires once the bubble has ended. If a trim is scheduled from one and still
// pending then, depotTrimArmed stays set and no trim runs in that process
// again. The depot is then bounded by depotMaxBytes only, and MemTrim still
// empties it.
func depotArm() {
	if !depotTrimArmed.Load() && depotTrimArmed.CompareAndSwap(false, true) {
		time.AfterFunc(depotTrimInterval, depotTrimLater)
	}
}

// depotTrimLater trims the depot and schedules the next trim if the depot still
// holds slots. A depotPut that saw a trim scheduled added its slots before
// depotTrimArmed is cleared here, so they are not left without one.
func depotTrimLater() {
	depotTrimMu.Lock()
	depotTrim()
	depotTrimMu.Unlock()
	depotTrimArmed.Store(false)
	if depotBytes.Load() != 0 {
		depotArm()
	}
}

// depotTrim frees, in each bin of the depot, as many lists as the bin held at
// the least since the last trim, oldest first. Those lists sat in the depot
// unused all that time. The caller holds depotTrimMu.
func depotTrim() {
	var lists []tcacheBin
	for i := range depot {
		d := &depot[i]
		d.mu.Lock()
		lists = depotTake(i, d.low, lists[:0])
		d.low = len(d.lists)
		d.mu.Unlock()
		depotFree(lists)
	}
}

// depotDrain frees every slot in the depot. The caller holds depotTrimMu.
func depotDrain() {
	var lists []tcacheBin
	for i := range depot {
		d := &depot[i]
		d.mu.Lock()
		lists = depotTake(i, len(d.lists), lists[:0])
		d.low = 0
		d.mu.Unlock()
		depotFree(lists)
	}
}

// depotTake removes the k oldest lists from bin i of the depot, which the
// caller has locked, and returns them appended to lists.
func depotTake(i, k int, lists []tcacheBin) []tcacheBin {
	if k == 0 {
		return lists
	}

	d := &depot[i&tcacheBinMask]
	var n int64
	for _, l := range d.lists[:k] {
		n += int64(l.n)
		lists = append(lists, l)
	}
	d.lists = d.lists[:copy(d.lists, d.lists[k:])]
	depotBytes.Add(-n * int64(tcacheClasses.size[i&tcacheBinMask]))
	return lists
}

// depotFree gives the slots of lists, taken out of the depot, back to the
// allocator, so that their pages can empty.
func depotFree(lists []tcacheBin) {
	for _, l := range lists {
		allocatorMu.Lock()
		p := l.head
		for j := 0; j < l.n; j++ {
			next := *(*uintptr)(unsafe.Pointer(p))
			allocator.UintptrFree(p)
			p = next
		}
		if memCounters {
			// Under allocatorMu, so that MemStat sees the slots
			// either in the depot or freed.
			depotSlots.Add(-int64(l.n))
		}
		allocatorMu.Unlock()
	}
}

// tcacheMalloc returns a slot for 0 < n <= tcacheMaxSize bytes, or zero with
// errno set to ENOMEM.
func (tls *TLS) tcacheMalloc(n Tsize_t) (r uintptr) {
	i := tcacheBinFor(uintptr(n))
	b := &tls.tcache[i]
	if r = b.head; r != 0 {
		b.head = *(*uintptr)(unsafe.Pointer(r))
		b.n--
		return r
	}

	if r = tls.tcacheRefill(b, i); r == 0 {
		tls.setErrno(ENOMEM)
	}
	return r
}

// tcacheRefill returns a new slot for bin i and puts more into b, which is bin
// i and empty: a list from the depot or, if it has none, up to half the
// capacity of b from the allocator. It returns zero if the allocator is out of
// memory.
func (tls *TLS) tcacheRefill(b *tcacheBin, i int) (r uintptr) {
	if l := depotGet(i); l.n != 0 {
		r = l.head
		b.head = *(*uintptr)(unsafe.Pointer(r))
		b.n = l.n - 1
		return r
	}

	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	size := int(tcacheClasses.size[i&tcacheBinMask])
	r, err := allocator.UintptrMalloc(size)
	if err != nil {
		return 0
	}

	// Link the slots in the order the allocator hands them out, which is
	// the order they then come out of b.
	var tail uintptr
	for n := int(tcacheClasses.max[i&tcacheBinMask]) / 2; b.n < n; b.n++ {
		p, err := allocator.UintptrMalloc(size)
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

// tcacheFree puts p, a slot of us bytes, into the cache of tls if us is the
// slot size of a bin and reports whether it did.
func (tls *TLS) tcacheFree(p uintptr, us int) bool {
	i := tcacheBinOf(us)
	if i < 0 {
		return false
	}

	b := &tls.tcache[i&tcacheBinMask]
	if b.n == int(tcacheClasses.max[i&tcacheBinMask]) {
		tls.tcacheFlush(b, i, b.n/2)
	}
	*(*uintptr)(unsafe.Pointer(p)) = b.head
	b.head = p
	b.n++
	return true
}

// tcacheFlush gives n > 0 slots from b, which is bin i, back to the depot or,
// if it is full, to the allocator.
func (tls *TLS) tcacheFlush(b *tcacheBin, i, n int) {
	tail := b.head
	for j := 1; j < n; j++ {
		tail = *(*uintptr)(unsafe.Pointer(tail))
	}
	if next := *(*uintptr)(unsafe.Pointer(tail)); depotPut(i, b.head, tail, n) {
		b.head = next
		b.n -= n
		return
	}

	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	for ; n != 0; n-- {
		p := b.head
		b.head = *(*uintptr)(unsafe.Pointer(p))
		b.n--
		allocator.UintptrFree(p)
	}
}

// tcacheFlushAll gives every slot in the cache of tls back to the depot or the
// allocator.
func (tls *TLS) tcacheFlushAll() {
	for i := range tls.tcache {
		if b := &tls.tcache[i]; b.n != 0 {
			tls.tcacheFlush(b, i, b.n)
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

	if n <= tcacheMaxSize {
		if r = depotGetOne(tcacheBinFor(uintptr(n))); r != 0 {
			return r
		}
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

	if rq <= tcacheMaxSize {
		if r = depotGetOne(tcacheBinFor(uintptr(rq))); r != 0 {
			clear((*RawMem)(unsafe.Pointer(r))[:rq:rq])
			return r
		}
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
	// This is what memory.UintptrRealloc does, but through Xmalloc and
	// Xfree, so that it uses the caches and takes no lock when the block
	// already fits.
	switch {
	case p == 0:
		return Xmalloc(tls, n)
	case n == 0:
		Xfree(tls, p)
		return 0
	}

	// The slot of a live block does not change size, so its usable size
	// can be read without the lock, as Xfree does.
	us := memory.UintptrUsableSize(p)
	if uint64(us) >= uint64(n) {
		return p
	}

	if r = Xmalloc(tls, n); r == 0 {
		return 0
	}

	copy((*RawMem)(unsafe.Pointer(r))[:us:us], (*RawMem)(unsafe.Pointer(p))[:us:us])
	Xfree(tls, p)
	return r
}

func Xfree(tls *TLS, p uintptr) {
	if __ccgo_strace {
		trc("tls=%v p=%v, (%v:)", tls, p, origin(2))
	}
	if p == 0 {
		return
	}

	if tls != nil {
		if tls.tcacheFree(p, memory.UintptrUsableSize(p)) {
			return
		}
	} else if i := tcacheBinOf(memory.UintptrUsableSize(p)); i >= 0 && depotPut(i, p, p, 1) {
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
	// The slot of a live block does not change size, so this needs no
	// lock.
	return Tsize_t(memory.UintptrUsableSize(p))
}

func MemAudit() (r []*MemAuditError) {
	return nil
}

func UsableSize(p uintptr) Tsize_t {
	return Tsize_t(memory.UintptrUsableSize(p))
}

type MemAllocatorStat struct {
	Allocs int
	Bytes  int
	Mmaps  int
}

// MemStat returns the global memory allocator statistics.
// should be compiled with the memory.counters build tag for the data to be available.
//
// Allocs does not count the free slots in the depot, but it does count those
// in the cache of each TLS until the TLS is closed.
func MemStat() MemAllocatorStat {
	allocatorMu.Lock()
	defer allocatorMu.Unlock()

	allocs := allocator.Allocs
	if memCounters {
		allocs -= int(depotSlots.Load())
	}
	return MemAllocatorStat{
		Allocs: allocs,
		Bytes:  allocator.Bytes,
		Mmaps:  allocator.Mmaps,
	}
}

// MemTrim gives the free slots in the depot back to the allocator and then
// returns to the OS the empty pages the allocator retains for reuse. Pages that
// still hold a live block, or a free slot in the cache of a TLS, stay mapped.
//
// MemTrim is never necessary for correctness: the depot frees what it does not
// use within two seconds. MemTrim trades a smaller resident set now for the
// cost of mapping the pages again later.
func MemTrim() error {
	depotTrimMu.Lock()
	depotDrain()
	depotTrimMu.Unlock()

	allocatorMu.Lock()

	defer allocatorMu.Unlock()

	return allocator.Trim()
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
