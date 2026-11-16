// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !(linux && (amd64 || arm64 || loong64 || ppc64le || s390x || riscv64 || 386 || arm))

package libc // import "modernc.org/libc"

// ReleaseFileLocks releases the stdio FILE locks that the C thread of tls
// still holds after a panic unwound its frames. On this target stdio does not
// lock FILEs, so there is nothing to release. See the linux version for
// details.
func (tls *TLS) ReleaseFileLocks() {}
