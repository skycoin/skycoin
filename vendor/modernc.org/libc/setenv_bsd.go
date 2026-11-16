// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build netbsd || openbsd

package libc // import "modernc.org/libc"

import (
	"sync"
	"unsafe"

	"modernc.org/libc/errno"
	"modernc.org/libc/sys/types"
)

// setenv and unsetenv change environ, the array getenv reads, and so does the
// local time zone of localtime and mktime when the name is TZ. Like musl, they
// free only the strings setenv allocated and only the environ array they
// allocated themselves. The initial environment and an array the program put
// in environ are left alone.

var (
	envMu    sync.Mutex
	envArray uintptr                  // the environ array setenv allocated
	envOwned = map[uintptr]struct{}{} // the strings setenv allocated
)

// envName returns the length of name, or false if name is NULL, empty or
// contains '='.
func envName(name uintptr) (int, bool) {
	if name == 0 {
		return 0, false
	}

	n := 0
	for ; ; n++ {
		switch *(*byte)(unsafe.Pointer(name + uintptr(n))) {
		case 0:
			return n, n != 0
		case '=':
			return 0, false
		}
	}
}

// envMatch reports whether the environ entry e is "name=..." with name n bytes long.
func envMatch(e, name uintptr, n int) bool {
	for i := uintptr(0); i < uintptr(n); i++ {
		if *(*byte)(unsafe.Pointer(e + i)) != *(*byte)(unsafe.Pointer(name + i)) {
			return false
		}
	}
	return *(*byte)(unsafe.Pointer(e + uintptr(n))) == '='
}

// envFree frees the environ entry e if setenv allocated it.
func envFree(t *TLS, e uintptr) {
	if _, ok := envOwned[e]; ok {
		delete(envOwned, e)
		Xfree(t, e)
	}
}

// int setenv(const char *name, const char *value, int overwrite);
func Xsetenv(t *TLS, name, value uintptr, overwrite int32) int32 {
	if __ccgo_strace {
		trc("t=%v name=%v value=%v overwrite=%v, (%v:)", t, name, value, overwrite, origin(2))
	}
	n, ok := envName(name)
	if !ok {
		t.setErrno(errno.EINVAL)
		return -1
	}

	envMu.Lock()

	defer envMu.Unlock()

	p := Environ()
	i := 0
	for ; ; i++ {
		e := *(*uintptr)(unsafe.Pointer(p + uintptr(i)*uintptrSize))
		if e == 0 {
			break
		}

		if envMatch(e, name, n) {
			if overwrite == 0 {
				return 0
			}

			break
		}
	}
	vn := int(Xstrlen(t, value))
	s := Xmalloc(t, types.Size_t(n+1+vn+1))
	if s == 0 {
		return -1
	}

	copy(unsafe.Slice((*byte)(unsafe.Pointer(s)), n), unsafe.Slice((*byte)(unsafe.Pointer(name)), n))
	*(*byte)(unsafe.Pointer(s + uintptr(n))) = '='
	copy(unsafe.Slice((*byte)(unsafe.Pointer(s+uintptr(n)+1)), vn+1), unsafe.Slice((*byte)(unsafe.Pointer(value)), vn+1))
	if slot := p + uintptr(i)*uintptrSize; *(*uintptr)(unsafe.Pointer(slot)) != 0 {
		old := *(*uintptr)(unsafe.Pointer(slot))
		*(*uintptr)(unsafe.Pointer(slot)) = s
		envFree(t, old)
		envOwned[s] = struct{}{}
		return 0
	}

	// Append: i entries plus the terminator become i+2 slots.
	size := types.Size_t(i+2) * types.Size_t(uintptrSize)
	var np uintptr
	if p == envArray {
		np = Xrealloc(t, p, size)
	} else if np = Xmalloc(t, size); np != 0 {
		copy(unsafe.Slice((*uintptr)(unsafe.Pointer(np)), i), unsafe.Slice((*uintptr)(unsafe.Pointer(p)), i))
	}
	if np == 0 {
		Xfree(t, s)
		return -1
	}

	*(*uintptr)(unsafe.Pointer(np + uintptr(i)*uintptrSize)) = s
	*(*uintptr)(unsafe.Pointer(np + uintptr(i+1)*uintptrSize)) = 0
	envArray = np
	Xenviron = np
	envOwned[s] = struct{}{}
	return 0
}

// int unsetenv(const char *name);
func Xunsetenv(t *TLS, name uintptr) int32 {
	if __ccgo_strace {
		trc("t=%v name=%v, (%v:)", t, name, origin(2))
	}
	n, ok := envName(name)
	if !ok {
		t.setErrno(errno.EINVAL)
		return -1
	}

	envMu.Lock()

	defer envMu.Unlock()

	// Remove every entry of name, keeping the order of the others.
	p := Environ()
	w := p
	for r := p; ; r += uintptrSize {
		e := *(*uintptr)(unsafe.Pointer(r))
		if e == 0 {
			*(*uintptr)(unsafe.Pointer(w)) = 0
			return 0
		}

		if envMatch(e, name, n) {
			envFree(t, e)
			continue
		}

		*(*uintptr)(unsafe.Pointer(w)) = e
		w += uintptrSize
	}
}
