package eth

import (
	"math/big"
)

// rlpItem is a byte string or a list, the two kinds of RLP item.
// https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/
type rlpItem interface {
	encode() []byte
}

type rlpString []byte

type rlpList []rlpItem

func rlpUint(v *big.Int) rlpString {
	if v.Sign() < 0 {
		panic("RLP integers are unsigned")
	}
	return rlpString(v.Bytes())
}

func (s rlpString) encode() []byte {
	if len(s) == 1 && s[0] < 0x80 {
		return []byte{s[0]}
	}
	return append(rlpHeader(0x80, len(s)), s...)
}

func (l rlpList) encode() []byte {
	var body []byte
	for _, item := range l {
		body = append(body, item.encode()...)
	}
	return append(rlpHeader(0xc0, len(body)), body...)
}

// rlpHeader is the length prefix, short form under 56 bytes.
func rlpHeader(offset byte, n int) []byte {
	if n < 56 {
		return []byte{offset + byte(n)} //nolint:gosec // n < 56
	}
	length := big.NewInt(int64(n)).Bytes()
	return append([]byte{offset + 55 + byte(len(length))}, length...) //nolint:gosec // an int is at most 8 bytes
}
