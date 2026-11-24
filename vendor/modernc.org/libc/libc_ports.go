// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || freebsd || illumos || netbsd || openbsd

// C functions shared by the hand-written unix ports.

package libc // import "modernc.org/libc"

import (
	"bufio"
	"io"
	"math/rand/v2"
	"runtime"
	"strconv"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
	"modernc.org/libc/errno"
	"modernc.org/libc/locale"
	"modernc.org/libc/sys/types"
)

// double difftime(time_t time1, time_t time0);
func Xdifftime(t *TLS, time1, time0 types.Time_t) float64 {
	if __ccgo_strace {
		trc("t=%v time1=%v time0=%v, (%v:)", t, time1, time0, origin(2))
	}
	return float64(time1) - float64(time0)
}

// int strcoll(const char *s1, const char *s2);
//
// The C locale is the only one, and in it strcoll compares like strcmp.
func Xstrcoll(t *TLS, s1, s2 uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v s1=%v s2=%v, (%v:)", t, s1, s2, origin(2))
	}
	return Xstrcmp(t, s1, s2)
}

// cLocaleName is the name of the C locale, the only one there is.
var cLocaleName = [2]byte{'C', 0}

// char *setlocale(int category, const char *locale);
//
// The C locale is the only one. setlocale accepts "C", "POSIX" and "", the
// environment's locale, which is the C locale too. Any other name fails and
// changes nothing. So does an invalid category, which also sets errno to EINVAL.
// Every category is always named "C".
func Xsetlocale(t *TLS, category int32, name uintptr) (r uintptr) {
	if __ccgo_strace {
		trc("t=%v category=%v name=%v, (%v:)", t, category, name, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	switch category {
	case locale.LC_ALL, locale.LC_COLLATE, locale.LC_CTYPE, locale.LC_MESSAGES,
		locale.LC_MONETARY, locale.LC_NUMERIC, locale.LC_TIME:
	default:
		t.setErrno(errno.EINVAL)
		return 0
	}

	if name != 0 {
		switch GoString(name) {
		case "", "C", "POSIX":
		default:
			return 0
		}
	}
	return uintptr(unsafe.Pointer(&cLocaleName[0]))
}

// sigsetLayout describes sigset_t: words 32-bit words with signal n in bit
// (n-1)%32 of word (n-1)/32. sigaddset and its siblings accept signals 1 to
// max, and sigfillset sets the lowest fill bits. darwin also accepts signal 0,
// which changes nothing and is not a member.
func sigsetLayout() (words int, max, fill int32) {
	switch runtime.GOOS {
	case "darwin":
		return 1, 32, 32
	case "freebsd":
		return 4, 128, 128
	case "illumos":
		return 4, 74, 74
	case "netbsd":
		return 4, 63, 128
	case "openbsd":
		return 1, 32, 32
	}
	panic(todo(""))
}

// sigsetBit returns the word of set that holds signo and the bit of signo in
// it, or a nil w for signal 0 on darwin. It sets errno to EINVAL if signo is
// not a signal number.
func sigsetBit(t *TLS, set uintptr, signo int32) (w *uint32, bit uint32, ok bool) {
	if signo == 0 && runtime.GOOS == "darwin" {
		return nil, 0, true
	}

	if _, max, _ := sigsetLayout(); signo <= 0 || signo > max {
		t.setErrno(errno.EINVAL)
		return nil, 0, false
	}

	n := uint32(signo - 1)
	return (*uint32)(unsafe.Pointer(set + uintptr(n/32)*4)), 1 << (n % 32), true
}

// int sigemptyset(sigset_t *set);
func Xsigemptyset(t *TLS, set uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v set=%v, (%v:)", t, set, origin(2))
	}
	words, _, _ := sigsetLayout()
	clear(unsafe.Slice((*uint32)(unsafe.Pointer(set)), words))
	return 0
}

// int sigfillset(sigset_t *set);
func Xsigfillset(t *TLS, set uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v set=%v, (%v:)", t, set, origin(2))
	}
	words, _, fill := sigsetLayout()
	s := unsafe.Slice((*uint32)(unsafe.Pointer(set)), words)
	for i := range s {
		switch n := fill - int32(i)*32; {
		case n >= 32:
			s[i] = ^uint32(0)
		case n > 0:
			s[i] = 1<<n - 1
		default:
			s[i] = 0
		}
	}
	return 0
}

// int sigaddset(sigset_t *set, int signum);
func Xsigaddset(t *TLS, set uintptr, signum int32) int32 {
	if __ccgo_strace {
		trc("t=%v set=%v signum=%v, (%v:)", t, set, signum, origin(2))
	}
	w, bit, ok := sigsetBit(t, set, signum)
	if !ok {
		return -1
	}

	if w != nil {
		*w |= bit
	}
	return 0
}

// int sigdelset(sigset_t *set, int signum);
func Xsigdelset(t *TLS, set uintptr, signum int32) int32 {
	if __ccgo_strace {
		trc("t=%v set=%v signum=%v, (%v:)", t, set, signum, origin(2))
	}
	w, bit, ok := sigsetBit(t, set, signum)
	if !ok {
		return -1
	}

	if w != nil {
		*w &^= bit
	}
	return 0
}

// int sigismember(const sigset_t *set, int signum);
func Xsigismember(t *TLS, set uintptr, signum int32) int32 {
	if __ccgo_strace {
		trc("t=%v set=%v signum=%v, (%v:)", t, set, signum, origin(2))
	}
	w, bit, ok := sigsetBit(t, set, signum)
	if !ok {
		return -1
	}

	return Bool32(w != nil && *w&bit != 0)
}

// fileLock is the lock flockfile takes on a FILE. It is recursive and belongs
// to the TLS that holds it.
type fileLock struct {
	owner *TLS
	n     int // How many times owner has locked it.
}

// The flockfile locks held, by FILE. Unlike in C, the stdio functions do not
// take them: a held lock keeps other threads out of flockfile and ftrylockfile,
// but not out of the stdio functions. Each stdio function is atomic by itself:
// on darwin it takes the FILE's own lock, and the other ports do not buffer.
var (
	fileLocksMu   sync.Mutex
	fileLocksCond = sync.NewCond(&fileLocksMu)
	fileLocks     = map[uintptr]*fileLock{}
)

// tryLockFile takes the flockfile lock of f for t and reports whether it did.
// fileLocksMu must be held.
func tryLockFile(t *TLS, f uintptr) bool {
	switch l := fileLocks[f]; {
	case l == nil:
		fileLocks[f] = &fileLock{owner: t, n: 1}
	case l.owner == t:
		l.n++
	default:
		return false
	}
	return true
}

// void flockfile(FILE *filehandle);
func Xflockfile(t *TLS, filehandle uintptr) {
	if __ccgo_strace {
		trc("t=%v filehandle=%v, (%v:)", t, filehandle, origin(2))
	}
	fileLocksMu.Lock()

	defer fileLocksMu.Unlock()

	for !tryLockFile(t, filehandle) {
		fileLocksCond.Wait()
	}
}

// int ftrylockfile(FILE *filehandle);
func Xftrylockfile(t *TLS, filehandle uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v filehandle=%v, (%v:)", t, filehandle, origin(2))
	}
	fileLocksMu.Lock()

	defer fileLocksMu.Unlock()

	return Bool32(!tryLockFile(t, filehandle))
}

// void funlockfile(FILE *filehandle);
func Xfunlockfile(t *TLS, filehandle uintptr) {
	if __ccgo_strace {
		trc("t=%v filehandle=%v, (%v:)", t, filehandle, origin(2))
	}
	fileLocksMu.Lock()

	defer fileLocksMu.Unlock()

	if l := fileLocks[filehandle]; l != nil && l.owner == t {
		if l.n--; l.n == 0 {
			delete(fileLocks, filehandle)
			fileLocksCond.Broadcast()
		}
	}
}

// int getc_unlocked(FILE *stream);
func Xgetc_unlocked(t *TLS, stream uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v stream=%v, (%v:)", t, stream, origin(2))
	}
	return Xfgetc(t, stream)
}

// int putc_unlocked(int c, FILE *stream);
func Xputc_unlocked(t *TLS, c int32, stream uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v c=%v stream=%v, (%v:)", t, c, stream, origin(2))
	}
	return Xfputc(t, c, stream)
}

// int fseek(FILE *stream, long offset, int whence);
func Xfseek(t *TLS, stream uintptr, offset long, whence int32) int32 {
	if __ccgo_strace {
		trc("t=%v stream=%v offset=%v whence=%v, (%v:)", t, stream, offset, whence, origin(2))
	}
	return fseeko(t, stream, types.Off_t(offset), whence)
}

// int fseeko(FILE *stream, off_t offset, int whence);
func Xfseeko(t *TLS, stream uintptr, offset types.Off_t, whence int32) int32 {
	if __ccgo_strace {
		trc("t=%v stream=%v offset=%v whence=%v, (%v:)", t, stream, offset, whence, origin(2))
	}
	return fseeko(t, stream, offset, whence)
}

// long ftell(FILE *stream);
func Xftell(t *TLS, stream uintptr) long {
	if __ccgo_strace {
		trc("t=%v stream=%v, (%v:)", t, stream, origin(2))
	}
	n := ftello(t, stream)
	if types.Off_t(long(n)) != n {
		t.setErrno(errno.EOVERFLOW)
		return -1
	}

	return long(n)
}

// off_t ftello(FILE *stream);
func Xftello(t *TLS, stream uintptr) types.Off_t {
	if __ccgo_strace {
		trc("t=%v stream=%v, (%v:)", t, stream, origin(2))
	}
	return ftello(t, stream)
}

// fopenFlags returns the open flags for an fopen mode: 'r', 'w' or 'a', then
// any of '+', 'b', 'e' (close-on-exec) and 'x' (fail if the file exists, with
// 'w' only). Other characters are ignored, as in musl and glibc.
func fopenFlags(mode uintptr) (flags int, ok bool) {
	s := GoString(mode)
	if s == "" {
		return 0, false
	}

	switch s[0] {
	case 'r':
		flags = unix.O_RDONLY
	case 'w':
		flags = unix.O_WRONLY | unix.O_CREAT | unix.O_TRUNC
	case 'a':
		flags = unix.O_WRONLY | unix.O_CREAT | unix.O_APPEND
	default:
		return 0, false
	}

	for _, c := range s[1:] {
		switch c {
		case '+':
			flags = flags&^unix.O_ACCMODE | unix.O_RDWR
		case 'e':
			flags |= unix.O_CLOEXEC
		case 'x':
			if s[0] == 'w' {
				flags |= unix.O_EXCL
			}
		}
	}
	return flags, true
}

// FILE *freopen(const char *pathname, const char *mode, FILE *stream);
//
// The stream keeps its descriptor: the new file is opened and moved onto it
// with dup2. With a NULL pathname, freopen changes the mode of the descriptor
// as FreeBSD does: the descriptor must allow the new access mode, "w" truncates
// the file, and the stream moves to the start of the file unless it appends.
// On failure freopen closes the stream.
func Xfreopen(t *TLS, pathname, mode, stream uintptr) (r uintptr) {
	if __ccgo_strace {
		trc("t=%v pathname=%v mode=%v stream=%v, (%v:)", t, pathname, mode, stream, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	var err error
	flags, ok := fopenFlags(mode)
	switch {
	case !ok:
		err = unix.EINVAL
	default:
		Xfflush(t, stream)
		fd := int(file(stream).fd())
		if pathname == 0 {
			err = reopenMode(fd, flags)
		} else {
			err = reopenPath(fd, GoString(pathname), flags)
		}
	}
	if err != nil {
		Xfclose(t, stream)
		t.setErrno(err)
		return 0
	}

	file(stream).reopened()
	return stream
}

// reopenMode changes the open flags of fd to those of flags, see Xfreopen.
func reopenMode(fd, flags int) error {
	fl, err := unix.FcntlInt(uintptr(fd), unix.F_GETFL, 0)
	if err != nil {
		return err
	}

	if acc := fl & unix.O_ACCMODE; acc != unix.O_RDWR && acc != flags&unix.O_ACCMODE {
		return unix.EBADF
	}

	if (fl^flags)&unix.O_APPEND != 0 {
		if _, err := unix.FcntlInt(uintptr(fd), unix.F_SETFL, fl^unix.O_APPEND); err != nil {
			return err
		}
	}
	if flags&unix.O_CLOEXEC != 0 {
		unix.CloseOnExec(fd)
	}
	if flags&unix.O_TRUNC != 0 {
		unix.Ftruncate(fd, 0)
	}
	if flags&unix.O_APPEND == 0 {
		unix.Seek(fd, 0, io.SeekStart)
	}
	return nil
}

// reopenPath opens path with flags and moves the new descriptor onto fd.
func reopenPath(fd int, path string, flags int) error {
	nfd, err := unix.Open(path, flags|unix.O_CLOEXEC, 0666)
	if err != nil {
		return err
	}

	if nfd != fd { // fd may have been closed before.
		err = unix.Dup2(nfd, fd)
		unix.Close(nfd)
		if err != nil {
			return err
		}
	}

	// dup2 does not copy close-on-exec, but if nfd is fd, it has it.
	if flags&unix.O_CLOEXEC != 0 {
		unix.CloseOnExec(fd)
	} else if _, err := unix.FcntlInt(uintptr(fd), unix.F_SETFD, 0); err != nil {
		return err
	}
	return nil
}

// int fscanf(FILE *stream, const char *format, ...);
func Xfscanf(t *TLS, stream, format, va uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v stream=%v format=%v va=%v, (%v:)", t, stream, format, va, origin(2))
	}
	return scanf(&byteScanner{t: t, stream: stream}, format, va)
}

// byteScanner reads a stream for scanf. UnreadByte pushes the last byte read
// back with ungetc, so scanf leaves the stream at the first byte it did not
// use.
type byteScanner struct {
	t      *TLS
	stream uintptr

	last byte
	ok   bool // last holds a byte read and not yet pushed back
}

func (s *byteScanner) ReadByte() (byte, error) {
	c := Xfgetc(s.t, s.stream)
	if c < 0 {
		s.ok = false
		return 0, io.EOF
	}

	s.last = byte(c)
	s.ok = true
	return byte(c), nil
}

func (s *byteScanner) UnreadByte() error {
	if !s.ok {
		return bufio.ErrInvalidUnreadByte
	}

	s.ok = false
	Xungetc(s.t, int32(s.last), s.stream)
	return nil
}

// tmpnamBuf is the static buffer tmpnam uses when its argument is NULL.
// L_tmpnam is at least 25 on these ports, more than the names need.
var tmpnamBuf [20]byte

// char *tmpnam(char *s);
//
// tmpnam makes names as musl does: /tmp/tmpnam_ and six random letters and
// digits. It returns the first name that does not exist, written to s or, if s
// is NULL, to a static buffer. It returns NULL if 100 names all exist.
func Xtmpnam(t *TLS, s uintptr) (r uintptr) {
	if __ccgo_strace {
		trc("t=%v s=%v, (%v:)", t, s, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	name := []byte("/tmp/tmpnam_XXXXXX\x00")
	for range 100 {
		for i := len(name) - 7; i < len(name)-1; i++ {
			name[i] = chars[rand.IntN(len(chars))]
		}
		var st unix.Stat_t
		if err := unix.Lstat(string(name[:len(name)-1]), &st); err != unix.ENOENT {
			continue
		}

		if s == 0 {
			s = uintptr(unsafe.Pointer(&tmpnamBuf[0]))
		}
		copy(unsafe.Slice((*byte)(unsafe.Pointer(s)), len(name)), name)
		return s
	}
	return 0
}

// strerrorBuf is the static buffer strerror returns. C lets strerror share it
// between threads; strerror_r does not use it.
var strerrorBuf [128]byte

// char *strerror(int errnum);
func Xstrerror(t *TLS, errnum int32) uintptr {
	if __ccgo_strace {
		trc("t=%v errnum=%v, (%v:)", t, errnum, origin(2))
	}
	s, _ := strerrorText(errnum)
	n := copy(strerrorBuf[:len(strerrorBuf)-1], s)
	strerrorBuf[n] = 0
	return uintptr(unsafe.Pointer(&strerrorBuf[0]))
}

// int strerror_r(int errnum, char *buf, size_t buflen);
//
// strerror_r is the XSI one of the BSDs, macOS and illumos. It writes what
// strerror returns to buf and returns 0. If buf is too small, it writes as
// much as fits and returns ERANGE. For an unknown errnum it returns EINVAL.
func Xstrerror_r(t *TLS, errnum int32, buf uintptr, buflen types.Size_t) (r int32) {
	if __ccgo_strace {
		trc("t=%v errnum=%v buf=%v buflen=%v, (%v:)", t, errnum, buf, buflen, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	s, known := strerrorText(errnum)
	if uint64(len(s)) >= uint64(buflen) {
		r = errno.ERANGE
		if buflen == 0 {
			return r
		}

		s = s[:buflen-1]
	}
	if !known {
		r = errno.EINVAL
	}
	b := unsafe.Slice((*byte)(unsafe.Pointer(buf)), len(s)+1)
	copy(b, s)
	b[len(s)] = 0
	return r
}

// strerrorText returns strerror(errnum) as the system's libc writes it, see
// strerrorTexts, and reports whether errnum is a known error number.
func strerrorText(errnum int32) (s string, known bool) {
	if errnum >= 0 && int(errnum) < len(strerrorTexts) {
		return strerrorTexts[errnum], true
	}

	if strerrorUnknownNumber {
		return strerrorUnknown + strconv.Itoa(int(errnum)), false
	}

	return strerrorUnknown, false
}
