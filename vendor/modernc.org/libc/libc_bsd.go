// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || netbsd || openbsd

package libc // import "modernc.org/libc"

import (
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
	"modernc.org/libc/locale"
	"modernc.org/libc/sys/types"
)

// isalnumC reports whether c is a letter or a digit in the C locale.
func isalnumC(c int32) bool {
	return c >= '0' && c <= '9' || c|32 >= 'a' && c|32 <= 'z'
}

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
	return Bool32(c > ' ' && c < 0x7f && !isalnumC(c))
}

// charMax returns CHAR_MAX: char is unsigned on arm and arm64.
func charMax() byte {
	switch runtime.GOARCH {
	case "arm", "arm64":
		return 0xff
	}
	return 0x7f
}

// setChars sets the char fields p of an lconv to v. Their Go type is int8 or
// uint8, depending on the target.
func setChars[T ~int8 | ~uint8](v byte, p ...*T) {
	for _, q := range p {
		*q = T(v)
	}
}

var (
	lconvOnce sync.Once
	lconvC    locale.Lconv // The lconv of the C locale.
	// The strings of lconvC: ".", "" and, on FreeBSD, the grouping string
	// {CHAR_MAX}. NetBSD and OpenBSD have "" for grouping.
	lconvStrings [4]byte
)

// struct lconv *localeconv(void);
//
// It returns the lconv of the C locale, as the system libc has it.
func Xlocaleconv(t *TLS) uintptr {
	if __ccgo_strace {
		trc("t=%v, (%v:)", t, origin(2))
	}
	lconvOnce.Do(func() {
		lconvStrings = [...]byte{'.', 0, charMax(), 0}
		dot := uintptr(unsafe.Pointer(&lconvStrings[0]))
		empty := dot + 1
		grouping := empty
		if runtime.GOOS == "freebsd" {
			grouping = dot + 2
		}
		l := &lconvC
		l.Fdecimal_point = dot
		l.Fthousands_sep = empty
		l.Fgrouping = grouping
		l.Fint_curr_symbol = empty
		l.Fcurrency_symbol = empty
		l.Fmon_decimal_point = empty
		l.Fmon_thousands_sep = empty
		l.Fmon_grouping = grouping
		l.Fpositive_sign = empty
		l.Fnegative_sign = empty
		setChars(charMax(),
			&l.Fint_frac_digits, &l.Ffrac_digits,
			&l.Fp_cs_precedes, &l.Fp_sep_by_space, &l.Fn_cs_precedes, &l.Fn_sep_by_space,
			&l.Fp_sign_posn, &l.Fn_sign_posn,
			&l.Fint_p_cs_precedes, &l.Fint_p_sep_by_space, &l.Fint_n_cs_precedes, &l.Fint_n_sep_by_space,
			&l.Fint_p_sign_posn, &l.Fint_n_sign_posn,
		)
	})
	return uintptr(unsafe.Pointer(&lconvC))
}

// time_t time(time_t *tloc);
func Xtime(t *TLS, tloc uintptr) types.Time_t {
	if __ccgo_strace {
		trc("t=%v tloc=%v, (%v:)", t, tloc, origin(2))
	}
	var tv unix.Timeval
	if err := unix.Gettimeofday(&tv); err != nil {
		t.setErrno(err)
		return -1
	}

	if tloc != 0 {
		*(*types.Time_t)(unsafe.Pointer(tloc)) = types.Time_t(tv.Sec)
	}
	return types.Time_t(tv.Sec)
}
