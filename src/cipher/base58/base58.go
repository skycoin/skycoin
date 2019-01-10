package base58

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidChar Invalid base58 character
	ErrInvalidChar = errors.New("Invalid base58 character")
	// ErrInvalidString Invalid base58 string
	ErrInvalidString = errors.New("Invalid base58 string")
)

// Alphabet is a a b58 alphabet.
type Alphabet struct {
	decode [128]int8
	encode [58]byte
}

// NewAlphabet creates a new alphabet from the passed string.
//
// It panics if the passed string is not 58 bytes long or isn't valid ASCII.
func NewAlphabet(s string) *Alphabet {
	if len(s) != 58 {
		panic("base58 alphabets must be 58 bytes long")
	}

	ret := &Alphabet{}

	copy(ret.encode[:], s)

	for i := range ret.decode {
		ret.decode[i] = -1
	}
	for i, b := range ret.encode {
		ret.decode[b] = int8(i)
	}

	return ret
}

// btcAlphabet is the bitcoin base58 alphabet.
var btcAlphabet = NewAlphabet("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Encode encodes the passed bytes into a base58 encoded string.
func Encode(bin []byte) string {
	return fastBase58EncodingAlphabet(bin, btcAlphabet)
}

// fastBase58EncodingAlphabet encodes the passed bytes into a base58 encoded
// string with the passed alphabet.
func fastBase58EncodingAlphabet(bin []byte, alphabet *Alphabet) string {
	binsz := len(bin)
	var i, j, zcount, high int
	var carry uint32

	for zcount < binsz && bin[zcount] == 0 {
		zcount++
	}

	size := (binsz-zcount)*138/100 + 1
	var buf = make([]uint32, size)

	high = size - 1
	for i = zcount; i < binsz; i++ {
		j = size - 1
		for carry = uint32(bin[i]); j > high || carry != 0; j-- {
			carry += buf[j] << 8
			buf[j] = carry % 58
			carry /= 58
		}
		high = j
	}

	for j = 0; j < size && buf[j] == 0; j++ {
	}

	var b58 = make([]byte, size-j+zcount)

	if zcount != 0 {
		for i = 0; i < zcount; i++ {
			b58[i] = '1'
		}
	}

	for i = zcount; j < size; i++ {
		b58[i] = alphabet.encode[buf[j]]
		j++
	}

	return string(b58)
}

// Decode decodes the base58 encoded bytes.
func Decode(str string) ([]byte, error) {
	return fastBase58DecodingAlphabet(str, btcAlphabet)
}

// fastBase58DecodingAlphabet decodes the base58 encoded bytes using the given
// b58 alphabet.
func fastBase58DecodingAlphabet(str string, alphabet *Alphabet) ([]byte, error) {
	if len(str) == 0 {
		return nil, ErrInvalidString
	}

	var (
		t, c   uint64
		zmask  uint32
		zcount int

		b58u  = []rune(str)
		b58sz = len(b58u)

		outisz    = (b58sz + 3) >> 2
		binu      = make([]byte, (b58sz+3)*3)
		bytesleft = b58sz & 3
	)

	if bytesleft > 0 {
		zmask = 0xffffffff << uint32(bytesleft*8)
	} else {
		bytesleft = 4
	}

	var outi = make([]uint32, outisz)

	for i := 0; i < b58sz && b58u[i] == '1'; i++ {
		zcount++
	}

	for _, r := range b58u {
		if r > 127 {
			return nil, ErrInvalidChar
		}
		if alphabet.decode[r] == -1 {
			return nil, ErrInvalidChar
		}

		c = uint64(alphabet.decode[r])

		for j := outisz - 1; j >= 0; j-- {
			t = uint64(outi[j])*58 + c
			c = (t >> 32) & 0x3f
			outi[j] = uint32(t & 0xffffffff)
		}

		// Neither of these should occur because the buffer is allocated ourselves
		if c > 0 {
			return nil, fmt.Errorf("output number too big (carry to the next int32)")
		}

		if outi[0]&zmask != 0 {
			return nil, fmt.Errorf("output number too big (last int32 filled too far)")
		}
	}

	var j, cnt int
	for j, cnt = 0, 0; j < outisz; j++ {
		for mask := byte(bytesleft-1) * 8; mask <= 0x18; mask, cnt = mask-8, cnt+1 {
			binu[cnt] = byte(outi[j] >> mask)
		}
		if j == 0 {
			bytesleft = 4 // because it could be less than 4 the first time through
		}
	}

	for n, v := range binu {
		if v > 0 {
			start := n - zcount
			if start < 0 {
				start = 0
			}
			return binu[start:cnt], nil
		}
	}
	return binu[:cnt], nil
}
