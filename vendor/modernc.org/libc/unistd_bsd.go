// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build netbsd || openbsd

package libc // import "modernc.org/libc"

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// int execve(const char *path, char *const argv[], char *const envp[]);
func Xexecve(tls *TLS, path uintptr, argv uintptr, envp uintptr) (r int32) {
	if __ccgo_strace {
		trc("tls=%v path=%v argv=%v envp=%v, (%v:)", tls, path, argv, envp, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	goPath := GoString(path)
	var goArgv, goEnvp []string
	for p := *(*uintptr)(unsafe.Pointer(argv)); p != 0; p = *(*uintptr)(unsafe.Pointer(argv)) {
		goArgv = append(goArgv, GoString(p))
		argv += unsafe.Sizeof(uintptr(0))
	}
	for p := *(*uintptr)(unsafe.Pointer(envp)); p != 0; p = *(*uintptr)(unsafe.Pointer(envp)) {
		goEnvp = append(goEnvp, GoString(p))
		envp += unsafe.Sizeof(uintptr(0))
	}
	if err := unix.Exec(goPath, goArgv, goEnvp); err != nil {
		tls.setErrno(err)
		return -1
	}

	panic("unreachable")
}

// int setuid(uid_t uid);
func Xsetuid(tls *TLS, uid uint32) (r int32) {
	if __ccgo_strace {
		trc("tls=%v uid=%v, (%v:)", tls, uid, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	if err := unix.Setuid(int(uid)); err != nil {
		tls.setErrno(err)
		return -1
	}

	return 0
}

// int setgid(gid_t gid);
func Xsetgid(tls *TLS, gid uint32) (r int32) {
	if __ccgo_strace {
		trc("tls=%v gid=%v, (%v:)", tls, gid, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	if err := unix.Setgid(int(gid)); err != nil {
		tls.setErrno(err)
		return -1
	}

	return 0
}

// int dup(int oldfd);
func Xdup(tls *TLS, fd int32) (r int32) {
	if __ccgo_strace {
		trc("tls=%v fd=%v, (%v:)", tls, fd, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	nfd, err := unix.Dup(int(fd))
	if err != nil {
		tls.setErrno(err)
		return -1
	}

	return int32(nfd)
}
