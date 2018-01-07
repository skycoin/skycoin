package base58

import (
	"errors"
	"fmt"
	"math/big"
)

// An Encoding is a radix 58 encoding/decoding scheme.
type Encoding struct {
	decode [256]int64
	encode [58]byte
}

// It panics if the passed string is not 58 bytes long or isn't valid ASCII.
func encoding(s string) *Encoding {
	if len(s) != 58 {
		panic("base58 alphabets must be 58 bytes long")
	}
	ret := new(Encoding)
	copy(ret.encode[:], s)
	for i := range ret.decode {
		ret.decode[i] = -1
	}
	for i, b := range ret.encode {
		ret.decode[b] = int64(i)
	}
	return ret
}

// encmap is the encoding scheme used for Bitcoin addresses.
const BTCALPAHBET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var encmap = encoding(BTCALPAHBET)

var (
	bn0   = big.NewInt(0)
	radix = big.NewInt(58)
	zero  = encmap.encode[0]
)

// Encode encodes the given bytes using bitcoin base58 encoding
func Encode(src []byte) (string, error) {
	if len(src) == 0 {
		return "", errors.New("empty string given for encoding")
	}

	idx := len(src)*138/100 + 1 // log(256) / log(58), rounded up.
	buf := make([]byte, idx)
	mod := new(big.Int)

	// Get the integer value of string to encoded
	bn := new(big.Int).SetBytes(src)

	// divide the integer value till you hit zero
	for {
		switch bn.Cmp(bn0) {
		case 1:
			bn, mod = bn.DivMod(bn, radix, new(big.Int))
			idx--
			// store the mod value
			buf[idx] = encmap.encode[mod.Int64()]
		case 0:
			// check for zeros
			for i := range src {
				if src[i] != 0 {
					break
				}
				idx--
				buf[idx] = zero
			}
			return string(buf[idx:]), nil
		default:
			return "", fmt.Errorf("expecting a positive number in base58 encoding but got %q", bn)
		}
	}
}

// Decode decodes the base58 encoded bytes using the bitcoin base58 encoding
func Decode(str string) ([]byte, error) {
	zero := encmap.encode[0]

	// zero count
	var zeros int
	for i := 0; i < len(str) && str[i] == zero; i++ {
		zeros++
	}
	leading := make([]byte, zeros)

	var padChar rune = -1
	src := []byte(str)
	j := 0
	for ; j < len(src) && src[j] == byte(padChar); j++ {
	}

	n := new(big.Int)
	for i := range src[j:] {
		c := encmap.decode[src[i]]
		if c == -1 {
			return nil, errors.New("Invalid base58 character")
		}
		n.Mul(n, radix)
		n.Add(n, big.NewInt(int64(c)))
	}
	return append(leading, n.Bytes()...), nil
}
