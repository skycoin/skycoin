// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

package libc // import "modernc.org/libc"

import (
	"syscall"
	"unsafe"
)

// shellPath is the shell system and popen run, as in musl and the BSD libcs.
const shellPath = "/bin/sh"

// errnoOf returns the errno err carries, EINVAL if it carries none.
func errnoOf(err error) int32 {
	if e, ok := err.(syscall.Errno); ok {
		return int32(e)
	}

	return int32(syscall.EINVAL)
}

// environStrings returns the C environment, which setenv and putenv change,
// not the Go one.
func environStrings() (r []string) {
	for p := Environ(); p != 0; p += unsafe.Sizeof(uintptr(0)) {
		s := *(*uintptr)(unsafe.Pointer(p))
		if s == 0 {
			break
		}

		r = append(r, GoString(s))
	}
	return r
}

// spawnShell starts sh -c command in the C environment with files as its
// descriptors 0, 1 and 2 and returns its pid. The child also inherits every
// other descriptor that is not close-on-exec, as a C child does.
func spawnShell(command uintptr, files []uintptr) (pid int, err error) {
	return syscall.ForkExec(shellPath, []string{"sh", "-c", GoString(command)}, &syscall.ProcAttr{
		Env:   environStrings(),
		Files: files,
	})
}

// waitStatus waits for the child pid to terminate and returns its wait status,
// the value WIFEXITED, WEXITSTATUS, WIFSIGNALED and WTERMSIG decode.
func waitStatus(pid int) (int32, error) {
	for {
		var ws syscall.WaitStatus
		_, err := syscall.Wait4(pid, &ws, 0, nil)
		switch err {
		case nil:
			return int32(ws), nil
		case syscall.EINTR:
			// try again
		default:
			return -1, err
		}
	}
}

// int system(const char *command);
//
// Unlike musl and the BSD libcs, system does not ignore SIGINT and SIGQUIT or
// block SIGCHLD in the caller while the command runs. Signal dispositions
// belong to the Go runtime and os/signal.
func Xsystem(t *TLS, command uintptr) (r int32) {
	if __ccgo_strace {
		trc("t=%v command=%v, (%v:)", t, command, origin(2))
		defer func() { trc("-> %v", r) }()
	}
	if command == 0 {
		return 1 // POSIX: nonzero, a shell is always available.
	}

	pid, err := spawnShell(command, []uintptr{0, 1, 2})
	if err != nil {
		t.setErrno(errnoOf(err))
		return -1
	}

	ws, err := waitStatus(pid)
	if err != nil {
		t.setErrno(errnoOf(err))
		return -1
	}

	return ws
}

// popenMode returns the end of the pipe the caller of popen keeps: 0 for a
// mode starting with 'r', reading the output of the command, and 1 for 'w',
// writing its input. Like musl it ignores the rest of mode.
func popenMode(t *TLS, mode uintptr) (op int, ok bool) {
	switch *(*byte)(unsafe.Pointer(mode)) {
	case 'r':
		return 0, true
	case 'w':
		return 1, true
	}

	t.setErrno(int32(syscall.EINVAL))
	return 0, false
}

// popen starts sh -c command with one end of a new pipe as its descriptor
// 1-op. It returns the FILE fdopen makes of the other end and the pid of the
// child.
//
// The end the caller keeps stays close-on-exec, as with the 'e' mode flag of
// glibc and the BSDs, so no later child inherits it. That is how POSIX's rule
// that a popen child does not inherit the streams of earlier popen calls is
// kept: syscall.ForkExec cannot close them in the child the way musl does.
func popen(t *TLS, command uintptr, op int, fdopen func(fd int32) uintptr) (f uintptr, pid int) {
	var p [2]int
	syscall.ForkLock.RLock()
	err := syscall.Pipe(p[:])
	if err == nil {
		syscall.CloseOnExec(p[0])
		syscall.CloseOnExec(p[1])
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		t.setErrno(errnoOf(err))
		return 0, 0
	}

	if f = fdopen(int32(p[op])); f == 0 {
		syscall.Close(p[0])
		syscall.Close(p[1])
		return 0, 0
	}

	files := []uintptr{0, 1, 2}
	files[1-op] = uintptr(p[1-op])
	pid, err = spawnShell(command, files)
	syscall.Close(p[1-op])
	if err != nil {
		Xfclose(t, f)
		t.setErrno(errnoOf(err))
		return 0, 0
	}

	return f, pid
}
