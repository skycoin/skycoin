// Copyright 2017 The Memory Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package memory implements a memory allocator.
//
// # Build status
//
// available at https://modern-c.appspot.com/-/builder/?importpath=modernc.org%2fmemory
//
// # Changelog
//
// 2017-10-03 Added alternative, unsafe.Pointer-based API.
//
// Package memory implements a memory allocator.
//
// # Changelog
//
// 2017-10-03 Added alternative, unsafe.Pointer-based API.
//
// # Page retention
//
// Memory is acquired from the OS in 64 KiB units or, for larger
// allocations, in units rounded up to a power of two on 64-bit hosts.
// An Allocator retains empty regions for reuse instead of unmapping
// them as soon as their last slot is freed: each size class keeps its
// current carve-target page, and beyond that, retired regions of any
// size go into a pool bounded by 4 MiB plus the high-water mark of the
// live mapping. An allocation pattern that repeatedly drains a size
// class, or repeatedly frees and reallocates large buffers, no longer
// pays a mmap/munmap round trip per turnover. Use Allocator.Trim to
// hand the retained regions back.
//
// # Size classes
//
// Allocations of up to 16 KiB share pages, in slots of a fixed set of size
// classes: on 64-bit 16 to 128 bytes in steps of 16, then four classes per
// power of two (160, 192, 224, 256, 320, ...). A request just above a power
// of two thus no longer takes almost twice its size. Every slot is aligned to
// 2*unsafe.Sizeof(uintptr(0)) bytes. UsableSize reports the class size, and
// UsableSizeFor tells it in advance, from the request size alone.
//
// # Benchmarks
//
//	jnml@3900x:~/src/modernc.org/memory$ date ; go version ; go test -run @ -bench . -benchmem |& tee log
//	Mon Sep 25 16:02:02 CEST 2023
//	go version go1.21.1 linux/amd64
//	goos: linux
//	goarch: amd64
//	pkg: modernc.org/memory
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkFree16-24             	123506772	         9.802 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFree32-24             	73853230	        15.08 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFree64-24             	43070334	        25.15 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkCalloc16-24           	59353304	        18.92 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkCalloc32-24           	39415004	        29.00 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkCalloc64-24           	35825725	        32.02 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkGoCalloc16-24         	38274313	        26.99 ns/op	      16 B/op	       1 allocs/op
//	BenchmarkGoCalloc32-24         	44590477	        33.06 ns/op	      32 B/op	       1 allocs/op
//	BenchmarkGoCalloc64-24         	44233016	        37.20 ns/op	      64 B/op	       1 allocs/op
//	BenchmarkMalloc16-24           	145736911	         7.720 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkMalloc32-24           	128898334	         7.887 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkMalloc64-24           	149569483	         7.994 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrFree16-24      	117043012	         9.205 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrFree32-24      	77399617	        14.20 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrFree64-24      	48770785	        25.04 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrCalloc16-24    	79257636	        15.44 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrCalloc32-24    	49644562	        23.62 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrCalloc64-24    	39854710	        28.22 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrMalloc16-24    	252987727	         4.525 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrMalloc32-24    	241423840	         4.433 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkUintptrMalloc64-24    	256450324	         4.669 ns/op	       0 B/op	       0 allocs/op
//	PASS
//	ok  	modernc.org/memory	93.178s
//	jnml@3900x:~/src/modernc.org/memory$
package memory // import "modernc.org/memory"

import (
	"fmt"
	"math/bits"
	"os"
	"unsafe"
)

const (
	headerSize     = unsafe.Sizeof(page{})
	mallocAllign   = 2 * unsafe.Sizeof(uintptr(0))
	maxSlotSize    = 1 << maxSlotSizeLog
	maxSlotSizeLog = pageSizeLog - 2
	pageAvail      = pageSize - headerSize
	pageMask       = pageSize - 1
	pageSize       = 1 << pageSizeLog
)

// Size classes. A request for up to maxSlotSize bytes is served from a slot of
// the smallest size class that fits it. The classes are linearClasses
// multiples of mallocAllign, then classesPerDoubling evenly spaced classes
// between each power of two and the next, up to maxSlotSize. On 64-bit:
//
//	16, 32, 48, 64, 80, 96, 112, 128,
//	160, 192, 224, 256,
//	320, 384, 448, 512,
//	...
//	10240, 12288, 14336, 16384
//
// A request wastes less than mallocAllign bytes up to 8*mallocAllign, and less
// than a fifth of its slot above that, where with power-of-two classes a
// request just above a power of two wasted almost half. Every power of two from
// mallocAllign to maxSlotSize remains a class, so a power-of-two request is
// served exactly. Slots sit back to back after the page header, so every slot
// is mallocAllign aligned, but no more than that.
const (
	linearClasses      = 8
	classesPerDoubling = 4

	// maxClasses bounds the number of size classes, including class 0, which
	// marks a dedicated page. It is a power of two: class indexes are masked
	// with maxClasses-1 to spare bounds checks.
	maxClasses = 64
)

// sizeClass describes the slots of a size class.
type sizeClass struct {
	size  uint32 // bytes per slot
	slots uint32 // slots per page
}

var (
	classes [maxClasses]sizeClass

	// sizeToClass[i] is the class serving requests for i*mallocAllign+1 to
	// (i+1)*mallocAllign bytes.
	sizeToClass [maxSlotSize / mallocAllign]uint8
)

func init() {
	if unsafe.Sizeof(page{})%mallocAllign != 0 {
		panic("internal error")
	}

	initClasses(linearClasses, classesPerDoubling)
}

// initClasses fills the size class tables: linear classes mallocAllign apart,
// starting at mallocAllign, then perDoubling classes per power of two, up to
// maxSlotSize. initClasses(1, 1) yields the power-of-two classes.
func initClasses(linear, perDoubling int) {
	classes = [maxClasses]sizeClass{}
	n := 0
	add := func(size int) {
		if n++; n == maxClasses {
			panic("internal error: too many size classes")
		}

		classes[n] = sizeClass{uint32(size), uint32(int(pageAvail) / size)}
	}
	q := int(mallocAllign)
	for i := 1; i <= linear; i++ {
		add(i * q)
	}
	for base := linear * q; base < maxSlotSize; base *= 2 {
		step := base / perDoubling
		if step%q != 0 || base%perDoubling != 0 {
			panic("internal error: misaligned size class")
		}

		for size := base + step; size <= 2*base; size += step {
			add(size)
		}
	}
	if classes[n].size != maxSlotSize {
		panic("internal error: size classes do not end at maxSlotSize")
	}

	c := 1
	for i := range sizeToClass {
		for int(classes[c].size) < (i+1)*q {
			c++
		}
		sizeToClass[i] = uint8(c)
	}
}

// if n%m != 0 { n += m-n%m }. m must be a power of 2.
func roundup(n, m int) int { return (n + m - 1) &^ (m - 1) }

type node struct {
	prev, next uintptr // *node
}

type page struct {
	brk   int  // shared page: slots carved so far
	class uint // shared page: size class of its slots; 0: dedicated page
	size  int  // bytes mapped
	used  int  // shared page: live slots
}

// Allocator allocates and frees memory. Its zero value is ready for use.  The
// exported counters are updated only when build tag memory.counters is
// present.
//
// An Allocator retains empty regions for reuse, see Trim.
type Allocator struct {
	Allocs    int                  // # of allocs.
	Bytes     int                  // Asked from OS.
	lists     [maxClasses]uintptr  // *node, free slots of a size class
	Mmaps     int                  // Asked from OS.
	pages     [maxClasses]uintptr  // *page, carve target of a size class
	regs      map[uintptr]struct{} // map[*page]struct{}
	freed     map[int][]uintptr    // empty regions retained for reuse, keyed by region size
	freedSize int                  // total bytes in freed
	cold      map[int]int          // the bottom cold[size] regions of freed[size] are decommitted
	live      int                  // bytes mapped and not in freed
	hiLive    int                  // high-water mark of live
}

// maxFreedSize returns the bound on the total bytes retained in
// a.freed: a fixed floor plus the high-water mark of the live mapping.
// The high-water mark rather than the current live size matters for
// large transient buffers: when one is freed, live has already shrunk
// by the buffer, and it is exactly the region most worth retaining.
func (a *Allocator) maxFreedSize() int { return 4<<20 + a.hiLive }

// addLive adds n bytes to the live mapping size, tracking its
// high-water mark.
func (a *Allocator) addLive(n int) {
	a.live += n
	if a.live > a.hiLive {
		a.hiLive = a.live
	}
}

// canRetain reports whether a.freed has room to retain size more bytes.
func (a *Allocator) canRetain(size int) bool { return a.freedSize+size <= a.maxFreedSize() }

// mmapSize returns the region size class that backs a request for size
// bytes: requests up to pageSize share the single pageSize class, and
// on 64-bit larger requests are rounded up to a power of two. The
// extra address space is never touched and so costs nothing, and
// collapsing sizes into classes makes reuse via a.freed much more
// likely.
//
// Classes are multiples of mmapGranularity - the granularity the platform
// mmap actually maps at - so the class of a request agrees with the
// page.size the mapping comes back with, which is what retain later keys
// the pool by. Rounding to anything finer describes a mapping the platform
// cannot make: on windows/386 a 1 MiB+1 request rounded to OS pages misses
// forever the 1 MiB+64 KiB region its own mapping put into the pool.
func mmapSize(size int) int {
	switch {
	case size <= pageSize && pageSize%osPageSize == 0:
		return pageSize
	case size > pageSize && bits.UintSize == 64:
		return roundup(1<<bits.Len(uint(size-1)), mmapGranularity)
	default:
		return roundup(size, mmapGranularity)
	}
}

// slabBatch is how many pageSize regions a single mapping acquires
// when a pageSize region is needed and the pool is empty.
const slabBatch = 16

func (a *Allocator) mmap(size int) (uintptr /* *page */, error) {
	size = mmapSize(size)
	if s := a.freed[size]; len(s) != 0 {
		p := s[len(s)-1]
		if canDecommit {
			if c := a.cold[size]; c >= len(s) { // popping from the decommitted prefix
				if err := recommit(p, size); err != nil {
					return 0, err // pool left untouched
				}

				a.cold[size] = len(s) - 1
			}
		}
		a.freed[size] = s[:len(s)-1]
		a.freedSize -= size
		a.addLive(size)
		pg := (*page)(unsafe.Pointer(p))
		pg.brk = 0
		pg.used = 0
		return p, nil
	}

	// Aligning a fresh mapping to a pageSize boundary costs up to two
	// munmap calls for the trimmed-off ends, so when a single pageSize
	// region is needed, acquire a batch in one mapping and put the
	// rest into the pool: three syscalls amortized over slabBatch
	// regions instead of per region.
	if canCarve && size == pageSize {
		if p, n, err := mmap(slabBatch * pageSize); err == nil {
			if n%pageSize == 0 && a.canRetain(n-pageSize) {
				for off := pageSize; off < n; off += pageSize {
					q := p + uintptr(off)
					a.reg(q, pageSize)
					a.retain(q, pageSize)
				}
				return a.reg(p, pageSize), nil
			}

			if err := unmap(p, n); err != nil {
				return 0, err
			}
		}
		// Fall through and map a single region.
	}

	p, size, err := mmap(size)
	if err != nil {
		return 0, err
	}

	//TODO(jnml) The returned size may now be nearly as twice as large as we asked
	//for. Use that extra capacity. For that we need to move the slot count of
	//the size class into the page struct so the page cap becomes dynamic.
	//
	// Related: This is a consequence of fixing the bigsort.test failures on
	// linux/s390x, see: https://gitlab.com/cznic/sqlite/-/issues/207
	return a.reg(p, size), nil
}

// reg registers the freshly mapped size-byte region at p.
func (a *Allocator) reg(p uintptr, size int) uintptr {
	a.addLive(size)
	if counters {
		a.Mmaps++
		a.Bytes += size
	}
	if a.regs == nil {
		a.regs = map[uintptr]struct{}{}
	}
	(*page)(unsafe.Pointer(p)).size = size
	a.regs[p] = struct{}{}
	return p
}

func (a *Allocator) newPage(size int) (uintptr /* *page */, error) {
	size += int(headerSize)
	p, err := a.mmap(size)
	if err != nil {
		return 0, err
	}

	(*page)(unsafe.Pointer(p)).class = 0
	return p, nil
}

func (a *Allocator) newSharedPage(class uint) (uintptr /* *page */, error) {
	c := classes[class]
	size := int(headerSize) + int(c.slots)*int(c.size)
	p, err := a.mmap(size)
	if err != nil {
		return 0, err
	}

	a.pages[class] = p
	(*page)(unsafe.Pointer(p)).class = class
	return p, nil
}

func (a *Allocator) unmap(p uintptr /* *page */) error {
	delete(a.regs, p)
	size := (*page)(unsafe.Pointer(p)).size
	if counters {
		a.Mmaps--
		a.Bytes -= size
	}
	return unmap(p, size)
}

// release retires the empty region at p: it is retained in a.freed for
// reuse, unless the retained total is at its bound, in which case the
// region is returned to the OS. Retained regions stay in a.regs and in
// the counters; Trim and Close return them to the OS.
func (a *Allocator) release(p uintptr /* *page */) error {
	size := (*page)(unsafe.Pointer(p)).size
	if a.canRetain(size) {
		a.retain(p, size)
		return nil
	}

	a.live -= size
	return a.unmap(p)
}

// hotBytes is how much of the retained pool stays committed per size class.
// 4 MiB costs nothing measurable against not decommitting at all (speedtest1
// and production simulations) while returning the hoard above it.
const hotBytes = 4 << 20

// retain puts the empty size-byte region at p into a.freed.
func (a *Allocator) retain(p uintptr, size int) {
	if a.freed == nil {
		a.freed = map[int][]uintptr{}
	}
	a.freed[size] = append(a.freed[size], p)
	a.freedSize += size
	a.live -= size
	if !canDecommit {
		return
	}

	// Keep the newest hotBytes worth of this class committed: those are the
	// regions a churning workload takes straight back, where decommitting only
	// buys a page fault. What falls out of that window is hoard - it is holding
	// a resident set the workload has stopped using - so hand its pages back
	// while keeping the mapping.
	//
	// a.cold counts the stack's already-decommitted bottom, so only a region
	// newly falling out of the hot window pays a syscall: a pool oscillating
	// around the boundary would otherwise re-madvise the same region on every
	// retain.
	hot := hotBytes / size
	if hot < 1 {
		hot = 1
	}
	if s := a.freed[size]; len(s)-a.cold[size] > hot {
		if a.cold == nil {
			a.cold = map[int]int{}
		}
		decommit(s[a.cold[size]], size)
		a.cold[size]++
	}
}

// UintptrCalloc is like Calloc except it returns an uintptr.
func (a *Allocator) UintptrCalloc(size int) (r uintptr, err error) {
	if trace {
		defer func() {
			fmt.Fprintf(os.Stderr, "Calloc(%#x) %#x, %v\n", size, r, err)
		}()
	}
	if r, err = a.UintptrMalloc(size); r == 0 || err != nil {
		return 0, err
	}
	b := ((*rawmem)(unsafe.Pointer(r)))[:size:size]
	for i := range b {
		b[i] = 0
	}
	return r, nil
}

// UintptrFree is like Free except its argument is an uintptr, which must have
// been acquired from UintptrCalloc or UintptrMalloc or UintptrRealloc.
func (a *Allocator) UintptrFree(p uintptr) (err error) {
	if trace {
		defer func() {
			fmt.Fprintf(os.Stderr, "Free(%#x) %v\n", p, err)
		}()
	}
	if p == 0 {
		return nil
	}

	if counters {
		a.Allocs--
	}
	pg := p &^ uintptr(pageMask)
	class := (*page)(unsafe.Pointer(pg)).class
	if class == 0 {
		return a.release(pg)
	}

	(*node)(unsafe.Pointer(p)).prev = 0
	(*node)(unsafe.Pointer(p)).next = a.lists[class]
	if next := (*node)(unsafe.Pointer(p)).next; next != 0 {
		(*node)(unsafe.Pointer(next)).prev = p
	}
	a.lists[class] = p
	(*page)(unsafe.Pointer(pg)).used--
	if (*page)(unsafe.Pointer(pg)).used != 0 {
		return nil
	}

	size := uintptr(classes[class].size)
	for i := 0; i < (*page)(unsafe.Pointer(pg)).brk; i++ {
		n := pg + headerSize + uintptr(i)*size
		next := (*node)(unsafe.Pointer(n)).next
		prev := (*node)(unsafe.Pointer(n)).prev
		switch {
		case prev == 0:
			a.lists[class] = next
			if next != 0 {
				(*node)(unsafe.Pointer(next)).prev = 0
			}
		case next == 0:
			(*node)(unsafe.Pointer(prev)).next = 0
		default:
			(*node)(unsafe.Pointer(prev)).next = next
			(*node)(unsafe.Pointer(next)).prev = prev
		}
	}

	// The page is now empty. Retain it as the carve target of its size class
	// rather than returning it to the OS, unless the class already has one.
	// Otherwise a pattern that repeatedly drains the class pays a mmap/munmap
	// round trip per turnover. Trim releases the pages retained this way.
	if a.pages[class] == 0 || a.pages[class] == pg {
		(*page)(unsafe.Pointer(pg)).brk = 0
		a.pages[class] = pg
		return nil
	}

	return a.release(pg)
}

// UintptrMalloc is like Malloc except it returns an uinptr.
func (a *Allocator) UintptrMalloc(size int) (r uintptr, err error) {
	if trace {
		defer func() {
			fmt.Fprintf(os.Stderr, "Malloc(%#x) %#x, %v\n", size, r, err)
		}()
	}
	if size < 0 {
		panic("invalid malloc size")
	}

	if size == 0 {
		return 0, nil
	}
	if size > int(^uint(0)>>1)-int(headerSize) {
		return 0, fmt.Errorf("memory: allocation size %d is too large", size)
	}

	if counters {
		a.Allocs++
	}
	if size > maxSlotSize {
		p, err := a.newPage(size)
		if err != nil {
			return 0, err
		}

		return p + headerSize, nil
	}

	// The mask is a no-op on the table's values and spares the bounds checks
	// on the per-class arrays below.
	class := uint(sizeToClass[uintptr(size-1)/mallocAllign]) & (maxClasses - 1)
	if a.lists[class] == 0 && a.pages[class] == 0 {
		if _, err := a.newSharedPage(class); err != nil {
			return 0, err
		}
	}

	if p := a.pages[class]; p != 0 {
		pg := (*page)(unsafe.Pointer(p))
		c := classes[class]
		r := p + headerSize + uintptr(pg.brk)*uintptr(c.size)
		pg.used++
		if pg.brk++; pg.brk == int(c.slots) {
			a.pages[class] = 0
		}
		return r, nil
	}

	n := a.lists[class]
	p := n &^ uintptr(pageMask)
	a.lists[class] = (*node)(unsafe.Pointer(n)).next
	if next := (*node)(unsafe.Pointer(n)).next; next != 0 {
		(*node)(unsafe.Pointer(next)).prev = 0
	}
	(*page)(unsafe.Pointer(p)).used++
	return n, nil
}

// UintptrRealloc is like Realloc except its first argument is an uintptr,
// which must have been returned from UintptrCalloc, UintptrMalloc or
// UintptrRealloc.
func (a *Allocator) UintptrRealloc(p uintptr, size int) (r uintptr, err error) {
	if trace {
		defer func() {
			fmt.Fprintf(os.Stderr, "UnsafeRealloc(%#x, %#x) %#x, %v\n", p, size, r, err)
		}()
	}
	switch {
	case p == 0:
		return a.UintptrMalloc(size)
	case size == 0 && p != 0:
		return 0, a.UintptrFree(p)
	}

	us := UintptrUsableSize(p)
	if us >= size {
		return p, nil
	}

	if r, err = a.UintptrMalloc(size); err != nil {
		return 0, err
	}

	if us < size {
		size = us
	}
	copy((*rawmem)(unsafe.Pointer(r))[:size:size], (*rawmem)(unsafe.Pointer(p))[:size:size])
	return r, a.UintptrFree(p)
}

// UintptrUsableSize is like UsableSize except its argument is an uintptr,
// which must have been returned from UintptrCalloc, UintptrMalloc or
// UintptrRealloc.
func UintptrUsableSize(p uintptr) (r int) {
	if trace {
		defer func() {
			fmt.Fprintf(os.Stderr, "UsableSize(%#x) %#x\n", p, r)
		}()
	}
	if p == 0 {
		return 0
	}

	return usableSize(p)
}

func usableSize(p uintptr) (r int) {
	pg := p &^ uintptr(pageMask)
	if class := (*page)(unsafe.Pointer(pg)).class; class != 0 {
		return int(classes[class&(maxClasses-1)].size)
	}

	return (*page)(unsafe.Pointer(pg)).size - int(headerSize)
}

// Calloc is like Malloc except the allocated memory is zeroed.
func (a *Allocator) Calloc(size int) (r []byte, err error) {
	p, err := a.UintptrCalloc(size)
	if err != nil {
		return nil, err
	}

	return (*rawmem)(unsafe.Pointer(p))[:size:usableSize(p)], nil
}

// Close releases all OS resources used by a and sets it to its zero value.
//
// It's not necessary to Close the Allocator when exiting a process.
func (a *Allocator) Close() (err error) {
	for p := range a.regs {
		if e := a.unmap(p); e != nil && err == nil {
			err = e
		}
	}
	*a = Allocator{}
	return err
}

// Free deallocates memory (as in C.free). The argument of Free must have been
// acquired from Calloc or Malloc or Realloc.
func (a *Allocator) Free(b []byte) (err error) {
	if b = b[:cap(b)]; len(b) == 0 {
		return nil
	}

	return a.UintptrFree(uintptr(unsafe.Pointer(&b[0])))
}

// Malloc allocates size bytes and returns a byte slice of the allocated
// memory. The memory is not initialized. Malloc panics for size < 0 and
// returns (nil, nil) for zero size.
//
// It's ok to reslice the returned slice but the result of appending to it
// cannot be passed to Free or Realloc as it may refer to a different backing
// array afterwards.
func (a *Allocator) Malloc(size int) (r []byte, err error) {
	p, err := a.UintptrMalloc(size)
	if p == 0 || err != nil {
		return nil, err
	}

	return (*rawmem)(unsafe.Pointer(p))[:size:usableSize(p)], nil
}

// Realloc changes the size of the backing array of b to size bytes or returns
// an error, if any.  The contents will be unchanged in the range from the
// start of the region up to the minimum of the old and new  sizes.   If the
// new size is larger than the old size, the added memory will not be
// initialized.  If b's backing array is of zero size, then the call is
// equivalent to Malloc(size), for all values of size; if size is equal to
// zero, and b's backing array is not of zero size, then the call is equivalent
// to Free(b).  Unless b's backing array is of zero size, it must have been
// returned by an earlier call to Malloc, Calloc or Realloc.  If the area
// pointed to was moved, a Free(b) is done.
func (a *Allocator) Realloc(b []byte, size int) (r []byte, err error) {
	var p uintptr
	if b = b[:cap(b)]; len(b) != 0 {
		p = uintptr(unsafe.Pointer(&b[0]))
	}
	if p, err = a.UintptrRealloc(p, size); p == 0 || err != nil {
		return nil, err
	}

	return (*rawmem)(unsafe.Pointer(p))[:size:usableSize(p)], nil
}

// Trim returns to the OS the empty regions a retains for reuse: each
// size class's empty carve-target page and the bounded pool of retired
// regions. Pages still having live allocations are not affected and a
// stays ready for use.
//
// Trim is never necessary for correctness. It trades a smaller resident set
// now for the cost of mapping those pages again later.
func (a *Allocator) Trim() (err error) {
	if trace {
		defer func() {
			fmt.Fprintf(os.Stderr, "Trim() %v\n", err)
		}()
	}
	for class := range a.pages {
		pg := a.pages[class]
		if pg == 0 || (*page)(unsafe.Pointer(pg)).used != 0 {
			continue
		}

		a.pages[class] = 0
		a.live -= (*page)(unsafe.Pointer(pg)).size
		if e := a.unmap(pg); e != nil && err == nil {
			err = e
		}
	}
	for _, s := range a.freed {
		for _, pg := range s {
			if e := a.unmap(pg); e != nil && err == nil {
				err = e
			}
		}
	}
	a.freed = nil
	a.freedSize = 0
	a.cold = nil
	return err
}

// UsableSize reports the size of the memory block allocated at p, which must
// point to the first byte of a slice returned from Calloc, Malloc or Realloc.
// The allocated memory block size can be larger than the size originally
// requested from Calloc, Malloc or Realloc.
func UsableSize(p *byte) (r int) { return UintptrUsableSize(uintptr(unsafe.Pointer(p))) }

// UsableSizeFor returns the usable size, as UsableSize reports it, of the
// block that Malloc or Calloc, or their Unsafe and Uintptr variants, return for
// a request of size bytes. For size up to 16 KiB it is the size of the size
// class serving the request, and it is exact: a cache of freed blocks keyed by
// usable size, like the per-thread one of modernc.org/libc, can thus serve a
// request without asking the Allocator. For the size s of a size class,
// UsableSizeFor(s) is s. A larger block can come out larger still, in the rare
// case the OS refuses to trim a fresh mapping. The result is never less than
// size. UsableSizeFor(0) is 0, and UsableSizeFor panics for a negative size,
// as Malloc does.
//
// UsableSizeFor depends on no Allocator, so unlike the methods of Allocator it
// is safe for concurrent use.
func UsableSizeFor(size int) int {
	if i := uintptr(size-1) / mallocAllign; i < uintptr(len(sizeToClass)) {
		return int(classes[sizeToClass[i]&(maxClasses-1)].size)
	}

	return usableSizeForBig(size)
}

// usableSizeForBig is UsableSizeFor for sizes that do not fit a slot.
func usableSizeForBig(size int) int {
	switch {
	case size < 0:
		panic("invalid malloc size")
	case size == 0:
		return 0
	case size > int(^uint(0)>>1)-int(headerSize):
		return size // Malloc fails
	}

	// A dedicated page, as newPage maps it.
	if n := mmapSize(size+int(headerSize)) - int(headerSize); n >= size {
		return n
	}

	return size // mmapSize overflows, Malloc fails
}

// UnsafeCalloc is like Calloc except it returns an unsafe.Pointer.
func (a *Allocator) UnsafeCalloc(size int) (r unsafe.Pointer, err error) {
	p, err := a.UintptrCalloc(size)
	if err != nil {
		return nil, err
	}

	return unsafe.Pointer(p), nil
}

// UnsafeFree is like Free except its argument is an unsafe.Pointer, which must
// have been acquired from UnsafeCalloc or UnsafeMalloc or UnsafeRealloc.
func (a *Allocator) UnsafeFree(p unsafe.Pointer) (err error) { return a.UintptrFree(uintptr(p)) }

// UnsafeMalloc is like Malloc except it returns an unsafe.Pointer.
func (a *Allocator) UnsafeMalloc(size int) (r unsafe.Pointer, err error) {
	p, err := a.UintptrMalloc(size)
	if err != nil {
		return nil, err
	}

	return unsafe.Pointer(p), nil
}

// UnsafeRealloc is like Realloc except its first argument is an
// unsafe.Pointer, which must have been returned from UnsafeCalloc,
// UnsafeMalloc or UnsafeRealloc.
func (a *Allocator) UnsafeRealloc(p unsafe.Pointer, size int) (r unsafe.Pointer, err error) {
	q, err := a.UintptrRealloc(uintptr(p), size)
	if err != nil {
		return nil, err
	}

	return unsafe.Pointer(q), nil
}

// UnsafeUsableSize is like UsableSize except its argument is an
// unsafe.Pointer, which must have been returned from UnsafeCalloc,
// UnsafeMalloc or UnsafeRealloc.
func UnsafeUsableSize(p unsafe.Pointer) (r int) { return UintptrUsableSize(uintptr(p)) }
