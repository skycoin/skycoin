// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || illumos || netbsd || openbsd

// The stdio of these ports does not buffer: a FILE is a descriptor plus the
// EOF and error indicators. Each port has a file type with the methods id, fd,
// sysRead, eof, setEOF, clearEOF, err, setErr, clearErr and close. id returns
// the address that identifies the FILE, which on openbsd differs from the
// address of stdin, stdout or stderr. sysRead returns a nil error or a
// unix.Errno.

package libc // import "modernc.org/libc"

import (
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
	"modernc.org/libc/stdio"
	"modernc.org/libc/sys/types"
)

// The bytes ungetc pushes back onto a FILE are kept here, the last one pushed
// at the end. The functions that read a FILE return them, last first, before
// they read the descriptor, fseek drops them and ftell counts them.
var (
	pushbackMu sync.Mutex
	pushback   = map[uintptr][]byte{} // file.id() -> the bytes pushed back onto it
)

// takePushback moves up to len(b) bytes pushed back onto f to b, last pushed
// first, and returns their number.
func (f file) takePushback(b []byte) (n int) {
	pushbackMu.Lock()

	defer pushbackMu.Unlock()

	s := pushback[f.id()]
	for n < len(b) && len(s) != 0 {
		b[n] = s[len(s)-1]
		s = s[:len(s)-1]
		n++
	}
	if len(s) == 0 {
		delete(pushback, f.id())
	} else {
		pushback[f.id()] = s
	}
	return n
}

// pushedBack returns the number of bytes pushed back onto f.
func (f file) pushedBack() int {
	pushbackMu.Lock()

	defer pushbackMu.Unlock()

	return len(pushback[f.id()])
}

// dropPushback discards the bytes pushed back onto f.
func (f file) dropPushback() {
	pushbackMu.Lock()

	defer pushbackMu.Unlock()

	delete(pushback, f.id())
}

// reopened drops the bytes pushed back onto f and clears its EOF and error
// indicators, after freopen gave f a new file.
func (f file) reopened() {
	f.dropPushback()
	f.clearErr()
}

// read reads n bytes into the C memory at p: the bytes pushed back onto f
// first, then from the descriptor until it has n bytes, reaches the end of the
// file or fails. It sets the EOF and error indicators and returns the number of
// bytes read.
func (f file) read(t *TLS, p uintptr, n int) (r int) {
	if n <= 0 {
		return 0
	}

	r = f.takePushback(unsafe.Slice((*byte)(unsafe.Pointer(p)), n))
	for r < n {
		m, err := f.sysRead(p+uintptr(r), n-r)
		switch {
		case err == unix.EINTR:
			continue
		case err != nil:
			t.setErrno(err)
			f.setErr()
			return r
		case m == 0:
			f.setEOF()
			return r
		}

		r += m
	}
	return r
}

// int fgetc(FILE *stream);
func Xfgetc(t *TLS, stream uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v stream=%v, (%v:)", t, stream, origin(2))
	}
	p := t.Alloc(1)

	defer t.Free(1)

	if file(stream).read(t, p, 1) == 0 {
		return stdio.EOF
	}

	return int32(*(*byte)(unsafe.Pointer(p)))
}

// size_t fread(void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfread(t *TLS, ptr uintptr, size, nmemb types.Size_t, stream uintptr) types.Size_t {
	if __ccgo_strace {
		trc("t=%v ptr=%v size=%v nmemb=%v stream=%v, (%v:)", t, ptr, size, nmemb, stream, origin(2))
	}
	if size == 0 || nmemb == 0 {
		return 0
	}

	return types.Size_t(file(stream).read(t, ptr, int(size*nmemb))) / size
}

// int ungetc(int c, FILE *stream);
//
// As in the BSD libcs and glibc, more than one byte can be pushed back.
func Xungetc(t *TLS, c int32, stream uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v c=%v stream=%v, (%v:)", t, c, stream, origin(2))
	}
	if c == stdio.EOF {
		return stdio.EOF
	}

	f := file(stream)
	pushbackMu.Lock()

	defer pushbackMu.Unlock()

	pushback[f.id()] = append(pushback[f.id()], byte(c))
	f.clearEOF()
	return int32(byte(c))
}

// int feof(FILE *stream);
func Xfeof(t *TLS, stream uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v stream=%v, (%v:)", t, stream, origin(2))
	}
	return Bool32(file(stream).eof())
}

// void clearerr(FILE *stream);
func Xclearerr(t *TLS, stream uintptr) {
	if __ccgo_strace {
		trc("t=%v stream=%v, (%v:)", t, stream, origin(2))
	}
	file(stream).clearErr()
}

// int fclose(FILE *stream);
func Xfclose(t *TLS, stream uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v stream=%v, (%v:)", t, stream, origin(2))
	}
	file(stream).dropPushback()
	return file(stream).close(t)
}

// fseeko is fseek with an off_t offset.
func fseeko(t *TLS, stream uintptr, offset types.Off_t, whence int32) int32 {
	f := file(stream)
	if whence == stdio.SEEK_CUR {
		offset -= types.Off_t(f.pushedBack()) // Relative to the bytes pushed back.
	}
	if n := Xlseek(t, f.fd(), offset, whence); n < 0 {
		f.setErr()
		return -1
	}

	f.dropPushback()
	f.clearEOF()
	return 0
}

// ftello is ftell returning an off_t.
func ftello(t *TLS, stream uintptr) types.Off_t {
	f := file(stream)
	n := Xlseek(t, f.fd(), 0, stdio.SEEK_CUR)
	if n < 0 {
		f.setErr()
		return -1
	}

	// Below zero if more bytes are pushed back than were read. The C libcs
	// return -1 then.
	if n -= types.Off_t(f.pushedBack()); n < 0 {
		return -1
	}

	return n
}
