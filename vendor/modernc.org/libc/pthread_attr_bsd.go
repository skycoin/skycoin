// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || openbsd

package libc // import "modernc.org/libc"

import (
	"unsafe"

	"modernc.org/libc/errno"
)

// pthread_attr_t is a pointer on FreeBSD and OpenBSD, 4 or 8 bytes. Its first
// 4 bytes hold the detach state, the only attribute kept.
func attrDetachState(a uintptr) *int32 { return (*int32)(unsafe.Pointer(a)) }

// int pthread_attr_init(pthread_attr_t *attr);
func Xpthread_attr_init(tls *TLS, a uintptr) int32 {
	if __ccgo_strace {
		trc("tls=%v a=%v, (%v:)", tls, a, origin(2))
	}
	*(*uintptr)(unsafe.Pointer(a)) = 0
	return 0
}

func X__ccgo_pthreadAttrGetDetachState(tls *TLS, a uintptr) int32 {
	if __ccgo_strace {
		trc("tls=%v a=%v, (%v:)", tls, a, origin(2))
	}
	return *attrDetachState(a)
}

// int pthread_attr_getdetachstate(const pthread_attr_t *attr, int *detachstate);
func Xpthread_attr_getdetachstate(tls *TLS, a uintptr, state uintptr) int32 {
	if __ccgo_strace {
		trc("tls=%v a=%v state=%v, (%v:)", tls, a, state, origin(2))
	}
	*(*int32)(unsafe.Pointer(state)) = *attrDetachState(a)
	return 0
}

// int pthread_attr_setdetachstate(pthread_attr_t *attr, int detachstate);
func Xpthread_attr_setdetachstate(tls *TLS, a uintptr, state int32) int32 {
	if __ccgo_strace {
		trc("tls=%v a=%v state=%v, (%v:)", tls, a, state, origin(2))
	}
	if uint32(state) > 1 {
		return errno.EINVAL
	}
	*attrDetachState(a) = state
	return 0
}
