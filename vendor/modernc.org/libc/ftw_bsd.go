// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build netbsd || openbsd

package libc // import "modernc.org/libc"

import (
	"os"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
	"modernc.org/libc/errno"
)

// The ftw(3) type flags of <ftw.h>, the same on NetBSD and OpenBSD.
const (
	ftwF   = 0 // FTW_F, not a directory
	ftwD   = 1 // FTW_D, a directory
	ftwDNR = 2 // FTW_DNR, a directory that cannot be read
	ftwNS  = 4 // FTW_NS, stat failed
	ftwSL  = 5 // FTW_SL, a symbolic link
)

// int ftw(const char *path, int (*fn)(const char *, const struct stat *, int), int nfds);
//
// Like the ftw of the OS, it visits directories before their entries, in
// directory order, and reports a directory it cannot read first as FTW_D, then
// as FTW_DNR. It follows a symbolic link at the root and reports a dangling one
// there as FTW_NS with its lstat data. Below the root, OpenBSD's ftw follows
// symbolic links the same way, and a directory that is its own ancestor ends
// the walk with -1 and ELOOP. NetBSD's reports them as FTW_SL with their lstat
// data instead (ftwFollowLinks). nfds must be in [1, OPEN_MAX] (ftwMaxFDs) or
// ftw fails with EINVAL, but otherwise does not matter: at most one directory
// is open at a time.
func Xftw(tls *TLS, path, fn uintptr, nfds int32) (r int32) {
	if __ccgo_strace {
		trc("tls=%v path=%v fn=%v nfds=%v, (%v:)", tls, path, fn, nfds, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	if nfds < 1 || nfds > ftwMaxFDs {
		tls.setErrno(errno.EINVAL)
		return -1
	}

	statp := tls.Alloc(int(unsafe.Sizeof(unix.Stat_t{})))

	defer tls.Free(int(unsafe.Sizeof(unix.Stat_t{})))

	w := &ftwWalker{
		tls:   tls,
		fn:    *(*func(*TLS, uintptr, uintptr, int32) int32)(unsafe.Pointer(&struct{ uintptr }{fn})),
		statp: statp,
	}
	return w.walk(GoString(path), true)
}

type ftwWalker struct {
	tls       *TLS
	fn        func(*TLS, uintptr, uintptr, int32) int32
	statp     uintptr
	ancestors [][2]uint64 // dev, ino of the directories being walked
}

// call passes path, the stat data at w.statp and flag to the callback.
func (w *ftwWalker) call(path string, flag int32) int32 {
	cs, err := CString(path)
	if err != nil {
		w.tls.setErrno(errno.ENOMEM)
		return -1
	}

	defer Xfree(w.tls, cs)

	return w.fn(w.tls, cs, w.statp, flag)
}

func (w *ftwWalker) walk(path string, root bool) int32 {
	st := (*unix.Stat_t)(unsafe.Pointer(w.statp))
	stat := unix.Stat
	if !root && !ftwFollowLinks {
		stat = unix.Lstat
	}
	if err := stat(path, st); err != nil {
		if unix.Lstat(path, st) != nil {
			*st = unix.Stat_t{}
		}
		return w.call(path, ftwNS)
	}

	switch st.Mode & unix.S_IFMT {
	case unix.S_IFDIR:
		// Walked below.
	case unix.S_IFLNK:
		// Only seen when not following links.
		return w.call(path, ftwSL)
	default:
		return w.call(path, ftwF)
	}

	id := [2]uint64{uint64(st.Dev), uint64(st.Ino)}
	for _, v := range w.ancestors {
		if v == id {
			w.tls.setErrno(errno.ELOOP)
			return -1
		}
	}

	if r := w.call(path, ftwD); r != 0 {
		return r
	}

	var names []string
	f, err := os.Open(path)
	if err == nil {
		var entries []os.DirEntry
		entries, err = f.ReadDir(-1) // Directory order.
		f.Close()
		for _, v := range entries {
			names = append(names, v.Name())
		}
	}
	if err != nil {
		// *st still describes the directory.
		return w.call(path, ftwDNR)
	}

	prefix := path
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	w.ancestors = append(w.ancestors, id)

	defer func() { w.ancestors = w.ancestors[:len(w.ancestors)-1] }()

	for _, name := range names {
		if r := w.walk(prefix+name, false); r != 0 {
			return r
		}
	}
	return 0
}
