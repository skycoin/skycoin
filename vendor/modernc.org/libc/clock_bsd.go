// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || freebsd || netbsd || openbsd

package libc // import "modernc.org/libc"

import (
	"golang.org/x/sys/unix"
	"modernc.org/libc/time"
)

// clock_t clock(void);
//
// Like the clock of these systems, it returns the processor time of the
// process, user and system, from getrusage in CLOCKS_PER_SEC units: 1000000 a
// second on Darwin, 128 on FreeBSD, 100 on NetBSD and OpenBSD. The conversion
// truncates the microseconds the same way they do.
func Xclock(t *TLS) (r time.Clock_t) {
	if __ccgo_strace {
		trc("t=%v, (%v:)", t, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	var ru unix.Rusage
	if err := unix.Getrusage(unix.RUSAGE_SELF, &ru); err != nil {
		t.setErrno(err)
		return ^time.Clock_t(0) // (clock_t)-1
	}

	ticks := func(tv unix.Timeval) int64 {
		return int64(tv.Sec)*time.CLOCKS_PER_SEC + int64(tv.Usec)/(1000000/time.CLOCKS_PER_SEC)
	}
	return time.Clock_t(ticks(ru.Utime) + ticks(ru.Stime))
}
