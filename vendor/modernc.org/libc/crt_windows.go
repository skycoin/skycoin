// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// C standard library functions of the Windows port, as the UCRT does them.

package libc // import "modernc.org/libc"

import (
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	gotime "time"
	"unsafe"

	"golang.org/x/sys/windows"
	"modernc.org/libc/errno"
	"modernc.org/libc/time"
)

// fopenFlags returns the os.Open flags for an fopen mode. All streams are
// binary, so a 'b' or 't' changes nothing.
func fopenFlags(mode string) (flags int, ok bool) {
	switch strings.NewReplacer("b", "", "t", "").Replace(mode) {
	case "r":
		return os.O_RDONLY, true
	case "r+":
		return os.O_RDWR, true
	case "w":
		return os.O_WRONLY | os.O_CREATE | os.O_TRUNC, true
	case "w+":
		return os.O_RDWR | os.O_CREATE | os.O_TRUNC, true
	case "a":
		return os.O_WRONLY | os.O_CREATE | os.O_APPEND, true
	case "a+":
		return os.O_RDWR | os.O_CREATE | os.O_APPEND, true
	}
	return 0, false
}

// fopen opens pathname for an fopen mode and returns a new stream. On failure
// it sets errno and returns 0.
func fopen(t *TLS, pathname, mode uintptr) uintptr {
	h, ok := fopenHandle(t, pathname, mode)
	if !ok {
		return 0
	}

	p, _ := wrapFdHandle(h)
	return p
}

// fopenHandle opens pathname for an fopen mode. On failure it sets errno.
func fopenHandle(t *TLS, pathname, mode uintptr) (windows.Handle, bool) {
	flags, ok := fopenFlags(GoString(mode))
	if !ok {
		t.setErrno(errno.EINVAL)
		return 0, false
	}

	h, err := windows.Open(GoString(pathname), flags, 0o666)
	if err != nil {
		t.setErrno(err)
		return 0, false
	}

	return h, true
}

// FILE *freopen(const char *pathname, const char *mode, FILE *stream);
//
// freopen closes the file of stream and opens pathname in its place, for the
// same stream. If that fails, stream is closed. A NULL pathname, which would
// change only the mode, is EINVAL, as in the UCRT. For stdin, stdout and
// stderr the new file also becomes the process's standard handle, which child
// processes inherit.
func Xfreopen(t *TLS, pathname, mode, stream uintptr) (r uintptr) {
	if __ccgo_strace {
		trc("t=%v pathname=%v mode=%v stream=%v, (%v:)", t, pathname, mode, stream, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	f, ok := winGetObject(stream).(*file)
	if !ok {
		t.setErrno(errno.EBADF)
		return 0
	}

	if _, ok := fopenFlags(GoString(mode)); !ok || pathname == 0 {
		t.setErrno(errno.EINVAL)
		return 0
	}

	std := stdHandleOf(f._fd)
	if std == 0 {
		windows.Close(f.Handle) // C: "Failure to close the file is ignored."
	}
	h, ok := fopenHandle(t, pathname, mode)
	if !ok {
		remFile(f)
		return 0
	}

	f.Handle, f.hadErr, f.eof, f.unget, f.popen = h, false, false, -1, nil
	if std != 0 {
		windows.SetStdHandle(std, h)
	}
	return stream
}

// stdHandleOf returns the GetStdHandle id of fd 0, 1 or 2, else 0.
func stdHandleOf(fd int32) uint32 {
	switch fd {
	case 0:
		return windows.STD_INPUT_HANDLE
	case 1:
		return windows.STD_OUTPUT_HANDLE
	case 2:
		return windows.STD_ERROR_HANDLE
	}
	return 0
}

// tmpName returns a random name in the temporary directory, prefix and six
// letters and digits.
func tmpName(prefix string) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789" // file names ignore case
	b := []byte(prefix + "XXXXXX")
	for i := len(prefix); i < len(b); i++ {
		b[i] = chars[rand.IntN(len(chars))]
	}
	return filepath.Join(os.TempDir(), string(b))
}

// FILE *tmpfile(void);
//
// tmpfile creates a file in the temporary directory, open "w+b", which the
// system deletes when the stream is closed or the process ends.
func Xtmpfile(t *TLS) (r uintptr) {
	if __ccgo_strace {
		trc("t=%v, (%v:)", t, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	for range 100 {
		name, err := windows.UTF16PtrFromString(tmpName("tmpfile_"))
		if err != nil {
			t.setErrno(errno.EINVAL)
			return 0
		}

		h, err := windows.CreateFile(name,
			windows.GENERIC_READ|windows.GENERIC_WRITE,
			windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
			nil,
			windows.CREATE_NEW,
			windows.FILE_ATTRIBUTE_TEMPORARY|windows.FILE_FLAG_DELETE_ON_CLOSE,
			0,
		)
		if err == windows.ERROR_FILE_EXISTS {
			continue
		}

		if err != nil {
			t.setErrno(err)
			return 0
		}

		p, _ := wrapFdHandle(h)
		return p
	}
	t.setErrno(errno.EEXIST)
	return 0
}

// tmpnamBuf is the static buffer tmpnam uses when its argument is NULL. Its
// size is L_tmpnam of the UCRT headers, MAX_PATH.
var tmpnamBuf [260]byte

// char *tmpnam(char *s);
//
// tmpnam returns, as the UCRT's does, the full path of a file in the temporary
// directory that does not exist: tmpnam_ and six letters and digits. It writes
// it to s, which must hold L_tmpnam (260) bytes, or, if s is NULL, to a static
// buffer. msvcrt's tmpnam returned names like "\s1a4." instead, and its
// headers make L_tmpnam 14. tmpnam returns NULL if the path is too long for
// L_tmpnam or 100 names all exist.
func Xtmpnam(t *TLS, s uintptr) (r uintptr) {
	if __ccgo_strace {
		trc("t=%v s=%v, (%v:)", t, s, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	for range 100 {
		name := tmpName("tmpnam_")
		if len(name) >= len(tmpnamBuf) {
			t.setErrno(errno.ENAMETOOLONG)
			return 0
		}

		if _, err := os.Lstat(name); !os.IsNotExist(err) {
			continue
		}

		if s == 0 {
			s = uintptr(unsafe.Pointer(&tmpnamBuf[0]))
		}
		b := unsafe.Slice((*byte)(unsafe.Pointer(s)), len(name)+1)
		copy(b, name)
		b[len(name)] = 0
		return s
	}
	t.setErrno(errno.EEXIST)
	return 0
}

// int remove(const char *pathname);
//
// remove deletes a file. As in the UCRT, it does not delete directories: that
// fails with EACCES.
func Xremove(t *TLS, pathname uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v pathname=%v, (%v:)", t, pathname, origin(2))
	}
	p, err := windows.UTF16PtrFromString(GoString(pathname))
	if err != nil {
		t.setErrno(errno.EINVAL)
		return -1
	}

	if err := windows.DeleteFile(p); err != nil {
		t.setErrno(err)
		return -1
	}

	return 0
}

// clockStart is when the process started, for clock.
var clockStart = gotime.Now()

// clock_t clock(void);
//
// clock returns, as the UCRT's does, the wall time since the process started,
// in CLOCKS_PER_SEC (1000) units, or -1 once that does not fit clock_t.
func Xclock(t *TLS) (r time.Clock_t) {
	if __ccgo_strace {
		trc("t=%v, (%v:)", t, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	ms := gotime.Since(clockStart).Milliseconds()
	if ms > math.MaxInt32 {
		return -1
	}

	return time.Clock_t(ms)
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

// The ctype functions of the C locale, ASCII only, as in the UCRT.

// int iscntrl(int c);
func Xiscntrl(t *TLS, c int32) int32 {
	if __ccgo_strace {
		trc("t=%v c=%v, (%v:)", t, c, origin(2))
	}
	return Bool32(c >= 0 && c < ' ' || c == 0x7f)
}

// int isgraph(int c);
func Xisgraph(t *TLS, c int32) int32 {
	if __ccgo_strace {
		trc("t=%v c=%v, (%v:)", t, c, origin(2))
	}
	return Bool32(c > ' ' && c < 0x7f)
}

// int ispunct(int c);
func Xispunct(t *TLS, c int32) int32 {
	if __ccgo_strace {
		trc("t=%v c=%v, (%v:)", t, c, origin(2))
	}
	return Bool32(c > ' ' && c < 0x7f && !(c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z'))
}

// ucrtLconv is struct lconv of the UCRT headers. libc/locale.Lconv, generated from
// the msvcrt headers, ends before the wide _W_ fields; this one has both
// layouts as its prefix.
type ucrtLconv struct {
	decimal_point     uintptr
	thousands_sep     uintptr
	grouping          uintptr
	int_curr_symbol   uintptr
	currency_symbol   uintptr
	mon_decimal_point uintptr
	mon_thousands_sep uintptr
	mon_grouping      uintptr
	positive_sign     uintptr
	negative_sign     uintptr
	int_frac_digits   int8
	frac_digits       int8
	p_cs_precedes     int8
	p_sep_by_space    int8
	n_cs_precedes     int8
	n_sep_by_space    int8
	p_sign_posn       int8
	n_sign_posn       int8
	w                 [8]uintptr // _W_decimal_point ... _W_negative_sign
}

var (
	lconvDot   = [2]byte{'.', 0}
	lconvNone  = [1]byte{}
	lconvWDot  = [2]uint16{'.', 0}
	lconvWNone = [1]uint16{}
	cLconv     ucrtLconv
)

func init() {
	none := uintptr(unsafe.Pointer(&lconvNone[0]))
	cLconv = ucrtLconv{
		decimal_point:     uintptr(unsafe.Pointer(&lconvDot[0])),
		thousands_sep:     none,
		grouping:          none,
		int_curr_symbol:   none,
		currency_symbol:   none,
		mon_decimal_point: none,
		mon_thousands_sep: none,
		mon_grouping:      none,
		positive_sign:     none,
		negative_sign:     none,
		int_frac_digits:   math.MaxInt8, // CHAR_MAX
		frac_digits:       math.MaxInt8,
		p_cs_precedes:     math.MaxInt8,
		p_sep_by_space:    math.MaxInt8,
		n_cs_precedes:     math.MaxInt8,
		n_sep_by_space:    math.MaxInt8,
		p_sign_posn:       math.MaxInt8,
		n_sign_posn:       math.MaxInt8,
	}
	cLconv.w[0] = uintptr(unsafe.Pointer(&lconvWDot[0]))
	for i := 1; i < len(cLconv.w); i++ {
		cLconv.w[i] = uintptr(unsafe.Pointer(&lconvWNone[0]))
	}
}

// struct lconv *localeconv(void);
//
// localeconv returns the values of the C locale, the only one, as the UCRT's.
func Xlocaleconv(t *TLS) uintptr {
	if __ccgo_strace {
		trc("t=%v, (%v:)", t, origin(2))
	}
	return uintptr(unsafe.Pointer(&cLconv))
}
