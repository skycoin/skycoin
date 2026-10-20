package eth

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

// Address is an Ethereum account address.
type Address [20]byte

// ErrInvalidAddress is returned for text that is not an address, or a
// mixed-case address whose EIP-55 checksum fails.
var ErrInvalidAddress = errors.New("invalid Ethereum address")

// AddressFromPubKey is the low 20 bytes of the Keccak-256 of the 64-byte
// uncompressed public key.
func AddressFromPubKey(pk cipher.PubKey) (Address, error) {
	uncompressed := secp256k1.UncompressPubkey(pk[:])
	if len(uncompressed) != 65 {
		return Address{}, errors.New("cannot decompress public key")
	}
	h := Keccak256(uncompressed[1:])
	var a Address
	copy(a[:], h[12:])
	return a, nil
}

// String is the EIP-55 mixed-case checksum form.
// https://eips.ethereum.org/EIPS/eip-55
func (a Address) String() string {
	lower := hex.EncodeToString(a[:])
	h := Keccak256([]byte(lower))
	out := []byte("0x" + lower)
	for i := 0; i < len(lower); i++ {
		nibble := h[i/2] >> 4
		if i%2 == 1 {
			nibble = h[i/2] & 0x0f
		}
		if lower[i] >= 'a' && nibble >= 8 {
			out[i+2] = lower[i] - 'a' + 'A'
		}
	}
	return string(out)
}

// ParseAddress reads an address. All-lowercase and all-uppercase hex carry no
// checksum and are accepted; mixed case must match EIP-55 exactly.
func ParseAddress(s string) (Address, error) {
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return Address{}, ErrInvalidAddress
	}
	body := s[2:]
	if len(body) != 40 {
		return Address{}, ErrInvalidAddress
	}
	b, err := hex.DecodeString(body)
	if err != nil {
		return Address{}, ErrInvalidAddress
	}
	var a Address
	copy(a[:], b)
	if body != strings.ToLower(body) && body != strings.ToUpper(body) && a.String()[2:] != body {
		return Address{}, ErrInvalidAddress
	}
	return a, nil
}

// ErrAddressPubKey is returned when a public key is not behind an address.
var ErrAddressPubKey = errors.New("public key does not match the Ethereum address")

// Bytes is the 20 address bytes.
func (a Address) Bytes() []byte {
	return append([]byte{}, a[:]...)
}

// Null reports whether a is the zero address.
func (a Address) Null() bool {
	return a == Address{}
}

// Checksum is the first 4 bytes of the Keccak-256 of the address bytes.
func (a Address) Checksum() cipher.Checksum {
	h := Keccak256(a[:])
	var c cipher.Checksum
	copy(c[:], h[:4])
	return c
}

// Verify checks that pk is the public key behind a.
func (a Address) Verify(pk cipher.PubKey) error {
	got, err := AddressFromPubKey(pk)
	if err != nil {
		return err
	}
	if got != a {
		return ErrAddressPubKey
	}
	return nil
}
