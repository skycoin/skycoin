// +build cgo

package chacha

// TODO: Get gcc compiling this
// TODO: Decide using C Preprocessor or other method if to use ref or krovetz version

// #include "chacha20_krovetz.h"
import "C"

import (
	"unsafe"
)

// XORKeyStream crypts bytes from in to out using the given key, initialisation vector,
// constant and number of ChaCha20 rounds to perform.
//
// In and out may be the same slice but otherwise should not overlap. Counter
// contains the raw salsa20 counter bytes (both nonce and block counter).
// This calls Ted Krovetz's C chacha20 implementation (with one or two minor
// modifications), the speed boost over the go ref version is around 5x.
func XORKeyStream(out, in []byte, iv *[8]byte, constant *[16]byte, key *[32]byte, rounds int) {
	var constUint [4]uint32
	constUint[0] = uint32(constant[0]) | uint32(constant[1])<<8 | uint32(constant[2])<<16 | uint32(constant[3])<<24
	constUint[1] = uint32(constant[4]) | uint32(constant[5])<<8 | uint32(constant[6])<<16 | uint32(constant[7])<<24
	constUint[2] = uint32(constant[8]) | uint32(constant[9])<<8 | uint32(constant[10])<<16 | uint32(constant[11])<<24
	constUint[3] = uint32(constant[12]) | uint32(constant[13])<<8 | uint32(constant[14])<<16 | uint32(constant[15])<<24

	var cKey *C.uchar = (*C.uchar)(unsafe.Pointer(&key[0]))
	var cIv *C.uchar = (*C.uchar)(unsafe.Pointer(&iv[0]))
	var cIn *C.uchar = (*C.uchar)(unsafe.Pointer(&in[0]))
	var cOut *C.uchar = (*C.uchar)(unsafe.Pointer(&out[0]))
	var cConst *C.uint = (*C.uint)(unsafe.Pointer(&constUint[0]))
	var cRounds C.uint = C.uint(rounds)

	C.xor_key_stream(cOut, cIn, C.ulonglong(len(out)), cIv, cKey, cConst, cRounds)

	copy(out, C.GoBytes(unsafe.Pointer(cOut), C.int(len(out))))
}
