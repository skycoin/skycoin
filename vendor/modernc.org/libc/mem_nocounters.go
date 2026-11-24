// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !memory.counters

package libc // import "modernc.org/libc"

// memCounters reports whether modernc.org/memory counts allocations, which it
// does under the memory.counters build tag.
const memCounters = false
