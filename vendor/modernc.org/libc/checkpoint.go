// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libc // import "modernc.org/libc"

import (
	"errors"
	"unsafe"
)

// TLSCheckpointVersion is the version of TLSCheckpoint. Version 1 is stable. A
// change to the layout, size or meaning of TLSCheckpoint comes with a new
// version, and CheckRollback rejects checkpoints of other versions as
// malformed.
const TLSCheckpointVersion = 1

// TLSCheckpointSize is the size of a TLSCheckpoint in bytes. It depends on the
// target.
const TLSCheckpointSize = unsafe.Sizeof(TLSCheckpoint{})

// tlsCheckpointMagic is the first word of every TLSCheckpoint made by
// TLS.Checkpoint. It fits in 32 bits.
const tlsCheckpointMagic = 0x7c5c0000 | TLSCheckpointVersion

var (
	errTLSCheckpointMalformed = errors.New("libc: malformed TLS checkpoint")
	errTLSCheckpointOtherTLS  = errors.New("libc: TLS checkpoint of another TLS")
	errTLSCheckpointNotLive   = errors.New("libc: TLS checkpoint not live")
)

// TLSCheckpoint is the state of the stack-like allocations of a TLS: the depth
// of its TLS.Alloc stack, its alloca(3) memory and its alloca regions, see
// TLS.AllocaEntry. TLS.Rollback returns the TLS to that state, as longjmp
// returns to the stack setjmp saved. Memory from malloc is not stack-like and
// is left alone.
//
// A TLSCheckpoint holds only integers, no Go pointers, so it can be kept in C
// memory, for example in a jmp_buf. Code that does so sizes the space with
// TLSCheckpointSize, aligns it to a pointer and copies the checkpoint as a
// value. Its TLSCheckpointSize bytes start with a word that identifies
// TLSCheckpointVersion. The rest is private to libc. A checkpoint is valid only
// in the process and for the TLS that made it.
type TLSCheckpoint struct {
	magic    uintptr
	id       uintptr // TLS.ID
	stack    tlsCheckpointStack
	scopes   uintptr // len(TLS.allocaStack)
	scopeSeq uintptr // the seq of the innermost alloca region, zero if none
	allocas  uintptr // len(TLS.allocas)
	next     uintptr // TLS.allocaNext
	p        uintptr // TLS.allocaP
	end      uintptr // TLS.allocaEnd
}

// Checkpoint returns the current state of the stack-like allocations of tls,
// see TLSCheckpoint.
func (tls *TLS) Checkpoint() (r TLSCheckpoint) {
	r = TLSCheckpoint{
		magic:   tlsCheckpointMagic,
		id:      uintptr(uint32(tls.ID)),
		stack:   tls.stackCheckpoint(),
		scopes:  uintptr(len(tls.allocaStack)),
		allocas: uintptr(len(tls.allocas)),
		next:    uintptr(tls.allocaNext),
		p:       tls.allocaP,
		end:     tls.allocaEnd,
	}
	if n := len(tls.allocaStack); n != 0 {
		r.scopeSeq = tls.allocaStack[n-1].seq
	}
	return r
}

// CheckRollback reports whether TLS.Rollback can return tls to c. It changes
// nothing.
//
// It fails for a checkpoint that Checkpoint did not make, one made by another
// TLS, and one that is no longer live because tls has since left the
// allocations or the alloca region it was made in. A checkpoint is live as
// long as the C function that made it has not returned, as for a setjmp
// buffer, but CheckRollback cannot always tell when it has. When the function
// returned and another one reached the same depth, it can tell only if the
// innermost alloca region also changed.
func (tls *TLS) CheckRollback(c TLSCheckpoint) error {
	_, err := tls.checkRollback(&c)
	return err
}

// Rollback returns the stack-like allocations of tls to c, see TLSCheckpoint.
// It frees the alloca(3) memory allocated since c was made and leaves the
// alloca regions entered since then. TLS.Alloc memory is reused afterwards from
// where c was made.
//
// Rollback cannot fail when CheckRollback(c) reported no error and tls has not
// been used since. Otherwise, if CheckRollback(c) would report an error,
// Rollback panics with it before changing anything.
func (tls *TLS) Rollback(c TLSCheckpoint) {
	pos, err := tls.checkRollback(&c)
	if err != nil {
		panic(err)
	}

	tls.stackRollback(c.stack)
	tls.allocaFree(int(c.allocas))
	tls.allocaStack = tls.allocaStack[:c.scopes]
	tls.allocaNext, tls.allocaP, tls.allocaEnd = pos.next, pos.p, pos.end
}

// checkRollback implements CheckRollback. It returns the alloca position of c,
// normalized by allocaNormal.
func (tls *TLS) checkRollback(c *TLSCheckpoint) (pos allocaMark, err error) {
	switch {
	case c.magic != tlsCheckpointMagic:
		return pos, errTLSCheckpointMalformed
	case c.id != uintptr(uint32(tls.ID)):
		return pos, errTLSCheckpointOtherTLS
	}

	if err = tls.checkStackRollback(c.stack); err != nil {
		return pos, err
	}

	scopes := len(tls.allocaStack)
	switch {
	case c.scopes > uintptr(scopes), c.allocas > uintptr(len(tls.allocas)), c.next > uintptr(len(tls.allocaChunks)):
		return pos, errTLSCheckpointNotLive
	case c.scopes == 0 && c.scopeSeq != 0:
		return pos, errTLSCheckpointMalformed
	case c.scopes != 0 && tls.allocaStack[c.scopes-1].seq != c.scopeSeq:
		// The region c was made in was left.
		return pos, errTLSCheckpointNotLive
	}

	pos = tls.allocaNormal(allocaMark{n: int(c.allocas), next: int(c.next), p: c.p, end: c.end})
	if pos.end == 0 {
		if pos.next != 0 || pos.p != 0 {
			return pos, errTLSCheckpointMalformed
		}
	} else {
		if pos.next == 0 {
			return pos, errTLSCheckpointMalformed
		}

		if ch := tls.allocaChunks[pos.next-1]; pos.end != ch.p+allocaChunkSize || pos.p < ch.p || pos.p > pos.end {
			return pos, errTLSCheckpointMalformed
		}
	}

	// The positions of the region c was made in, of c, of the first region
	// entered after c and of the current one must be in this order.
	cur := tls.allocaNormal(allocaMark{n: len(tls.allocas), next: tls.allocaNext, p: tls.allocaP, end: tls.allocaEnd})
	if !allocaNotAfter(pos, cur) {
		return pos, errTLSCheckpointNotLive
	}

	if c.scopes != 0 {
		if m := tls.allocaNormal(tls.allocaStack[c.scopes-1]); !allocaNotAfter(m, pos) {
			return pos, errTLSCheckpointNotLive
		}
	}
	if c.scopes < uintptr(scopes) {
		if m := tls.allocaNormal(tls.allocaStack[c.scopes]); !allocaNotAfter(pos, m) {
			return pos, errTLSCheckpointNotLive
		}
	}
	return pos, nil
}

// allocaNormal returns m, but if m is a position from before the first chunk
// of the bump region existed and the chunk exists now, the start of that chunk,
// as allocaSlow sets the marks of such positions.
func (tls *TLS) allocaNormal(m allocaMark) allocaMark {
	if m.end == 0 && len(tls.allocaChunks) != 0 {
		p := tls.allocaChunks[0].p
		m.next, m.p, m.end = 1, p, p+allocaChunkSize
	}
	return m
}

// allocaNotAfter reports whether the alloca state a, both normalized by
// allocaNormal, is not later than b.
func allocaNotAfter(a, b allocaMark) bool {
	if a.n > b.n {
		return false
	}

	return a.next < b.next || a.next == b.next && a.p <= b.p
}
