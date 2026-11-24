// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows && !(linux && (amd64 || arm64 || loong64 || ppc64le || s390x || riscv64 || 386 || arm))

package libc // import "modernc.org/libc"

import (
	"unsafe"
)

// setErrnoValue sets errno to the C errno value e.
func (t *TLS) setErrnoValue(e int32) {
	*(*int32)(unsafe.Pointer(t.errnop)) = e
}

// setErrnoSyscall sets errno for the system error e. On unix it is the errno
// value already.
func (t *TLS) setErrnoSyscall(e syscallErrno) {
	t.setErrnoValue(int32(e))
}
