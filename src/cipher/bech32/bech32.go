// Package bech32 implements BIP173 bech32 encoding for segwit addresses.
package bech32

import (
	"fmt"
	"strings"
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var charsetRev = [128]int8{
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	15, -1, 10, 17, 21, 20, 26, 30, 7, 5, -1, -1, -1, -1, -1, -1,
	-1, 29, -1, 24, 13, 25, 9, 8, 23, -1, 18, 22, 31, 27, 19, -1,
	1, 0, 3, 16, 11, 28, 12, 14, 6, 4, 2, -1, -1, -1, -1, -1,
	-1, 29, -1, 24, 13, 25, 9, 8, 23, -1, 18, 22, 31, 27, 19, -1,
	1, 0, 3, 16, 11, 28, 12, 14, 6, 4, 2, -1, -1, -1, -1, -1,
}

func polymod(values []int) uint32 {
	gen := [5]uint32{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}
	chk := uint32(1)
	for _, v := range values {
		b := chk >> 25
		chk = ((chk & 0x1ffffff) << 5) ^ uint32(v) //nolint:gosec // v is a 5-bit value (0-31)
		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 { //nolint:gosec // i is 0-4
				chk ^= gen[i]
			}
		}
	}
	return chk
}

func hrpExpand(hrp string) []int {
	ret := make([]int, 0, len(hrp)*2+1)
	for _, c := range hrp {
		ret = append(ret, int(c>>5))
	}
	ret = append(ret, 0)
	for _, c := range hrp {
		ret = append(ret, int(c&31))
	}
	return ret
}

func verifyChecksum(hrp string, data []int) bool {
	values := append(hrpExpand(hrp), data...)
	return polymod(values) == 1
}

func createChecksum(hrp string, data []int) []int {
	values := append(hrpExpand(hrp), data...)
	values = append(values, 0, 0, 0, 0, 0, 0)
	mod := polymod(values) ^ 1
	ret := make([]int, 6)
	for i := 0; i < 6; i++ {
		ret[i] = int((mod >> uint(5*(5-i))) & 31) //nolint:gosec // result fits in int
	}
	return ret
}

// Encode encodes a bech32 string from HRP and 5-bit data values.
func Encode(hrp string, data []int) (string, error) {
	chk := createChecksum(hrp, data)
	combined := append(data, chk...)

	var ret strings.Builder
	ret.WriteString(hrp)
	ret.WriteByte('1')
	for _, d := range combined {
		if d < 0 || d >= 32 {
			return "", fmt.Errorf("invalid data value: %d", d)
		}
		ret.WriteByte(charset[d])
	}
	return ret.String(), nil
}

// Decode decodes a bech32 string into its HRP and 5-bit data values.
func Decode(s string) (string, []int, error) {
	if len(s) > 90 {
		return "", nil, fmt.Errorf("bech32 string too long: %d", len(s))
	}

	// Find separator
	pos := strings.LastIndex(s, "1")
	if pos < 1 || pos+7 > len(s) {
		return "", nil, fmt.Errorf("invalid bech32 separator position")
	}

	hrp := strings.ToLower(s[:pos])
	dataStr := strings.ToLower(s[pos+1:])

	data := make([]int, len(dataStr))
	for i, c := range dataStr {
		if c >= 128 || charsetRev[c] == -1 {
			return "", nil, fmt.Errorf("invalid bech32 character: %c", c)
		}
		data[i] = int(charsetRev[c])
	}

	if !verifyChecksum(hrp, data) {
		return "", nil, fmt.Errorf("invalid bech32 checksum")
	}

	return hrp, data[:len(data)-6], nil
}

// ConvertBits converts a byte slice from one bit grouping to another.
func ConvertBits(data []byte, fromBits, toBits uint8, pad bool) ([]byte, error) {
	acc := 0
	bits := uint8(0)
	maxv := (1 << toBits) - 1
	var ret []byte

	for _, value := range data {
		if value>>fromBits != 0 {
			return nil, fmt.Errorf("invalid data: value %d exceeds %d bits", value, fromBits)
		}
		acc = (acc << fromBits) | int(value)
		bits += fromBits
		for bits >= toBits {
			bits -= toBits
			ret = append(ret, byte((acc>>bits)&maxv))
		}
	}

	if pad {
		if bits > 0 {
			ret = append(ret, byte((acc<<(toBits-bits))&maxv))
		}
	} else if bits >= fromBits {
		return nil, fmt.Errorf("invalid padding")
	} else if (acc<<(toBits-bits))&maxv != 0 {
		return nil, fmt.Errorf("non-zero padding")
	}

	return ret, nil
}

// SegwitEncode encodes a segwit address from HRP, witness version, and witness program.
func SegwitEncode(hrp string, version byte, program []byte) (string, error) {
	if version > 16 {
		return "", fmt.Errorf("invalid witness version: %d", version)
	}
	if len(program) < 2 || len(program) > 40 {
		return "", fmt.Errorf("invalid witness program length: %d", len(program))
	}
	if version == 0 && len(program) != 20 && len(program) != 32 {
		return "", fmt.Errorf("invalid witness v0 program length: %d", len(program))
	}

	conv, err := ConvertBits(program, 8, 5, true)
	if err != nil {
		return "", err
	}

	data := append([]int{int(version)}, intsFromBytes(conv)...)
	return Encode(hrp, data)
}

// SegwitDecode decodes a segwit address into witness version and program.
func SegwitDecode(hrp, addr string) (byte, []byte, error) {
	decHRP, data, err := Decode(addr)
	if err != nil {
		return 0, nil, err
	}
	if decHRP != hrp {
		return 0, nil, fmt.Errorf("HRP mismatch: expected %q, got %q", hrp, decHRP)
	}
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("empty data")
	}

	version := byte(data[0])
	if version > 16 {
		return 0, nil, fmt.Errorf("invalid witness version: %d", version)
	}

	program, err := ConvertBits(bytesFromInts(data[1:]), 5, 8, false)
	if err != nil {
		return 0, nil, err
	}

	if len(program) < 2 || len(program) > 40 {
		return 0, nil, fmt.Errorf("invalid witness program length: %d", len(program))
	}
	if version == 0 && len(program) != 20 && len(program) != 32 {
		return 0, nil, fmt.Errorf("invalid witness v0 program length: %d", len(program))
	}

	return version, program, nil
}

func intsFromBytes(b []byte) []int {
	r := make([]int, len(b))
	for i, v := range b {
		r[i] = int(v)
	}
	return r
}

func bytesFromInts(ints []int) []byte {
	b := make([]byte, len(ints))
	for i, v := range ints {
		b[i] = byte(v)
	}
	return b
}
