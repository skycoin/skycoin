package cipher

import (
	"errors"
	"log"

	"github.com/SkycoinProject/skycoin/src/cipher/base58"
)

var (
	// ErrAddressInvalidLength Unexpected size of address bytes buffer
	ErrAddressInvalidLength = errors.New("Invalid address length")
	// ErrAddressInvalidChecksum Computed checksum did not match expected value
	ErrAddressInvalidChecksum = errors.New("Invalid checksum")
	// ErrAddressInvalidVersion Unsupported address version value
	ErrAddressInvalidVersion = errors.New("Address version invalid")
	// ErrAddressInvalidPubKey Public key invalid for address
	ErrAddressInvalidPubKey = errors.New("Public key invalid for address")
	// ErrAddressInvalidFirstByte Invalid first byte in wallet import format string
	ErrAddressInvalidFirstByte = errors.New("first byte invalid")
	// ErrAddressInvalidLastByte 33rd byte in wallet import format string is invalid
	ErrAddressInvalidLastByte = errors.New("invalid 33rd byte")
)

/*
Addresses are the Ripemd160 of the double SHA256 of the public key
- public key must be in compressed format

In the block chain the address is 20+1 bytes
- the first byte is the version byte
- the next twenty bytes are RIPMD160(SHA256(SHA256(pubkey)))

In base 58 format the address is 20+1+4 bytes
- the first 20 bytes are RIPMD160(SHA256(SHA256(pubkey))).
-- this is to allow for any prefix in vanity addresses
- the next byte is the version byte
- the next 4 bytes are a checksum
-- the first 4 bytes of the SHA256 of the 21 bytes that come before

*/

// Checksum 4 bytes
type Checksum [4]byte

// Addresser defines an interface for cryptocurrency addresses
type Addresser interface {
	Bytes() []byte
	String() string
	Checksum() Checksum
	Verify(PubKey) error
	Null() bool
}

// PubKeyRipemd160 returns ripemd160(sha256(sha256(pubkey)))
func PubKeyRipemd160(pubKey PubKey) Ripemd160 {
	r1 := SumSHA256(pubKey[:])
	r2 := SumSHA256(r1[:])
	return HashRipemd160(r2[:])
}

// Address version is after Key to enable better vanity address generation
// Address struct is a 25 byte with a 20 byte public key hash, 1 byte address
// type and 4 byte checksum.
type Address struct {
	Version byte      //1 byte
	Key     Ripemd160 //20 byte pubkey hash
}

// AddressFromPubKey creates Address from PubKey as ripemd160(sha256(sha256(pubkey)))
func AddressFromPubKey(pubKey PubKey) Address {
	return Address{
		Version: 0,
		Key:     PubKeyRipemd160(pubKey),
	}
}

// AddressFromSecKey generates address from secret key
func AddressFromSecKey(secKey SecKey) (Address, error) {
	p, err := PubKeyFromSecKey(secKey)
	if err != nil {
		return Address{}, err
	}
	return AddressFromPubKey(p), nil
}

// MustAddressFromSecKey generates address from secret key, panics on error
func MustAddressFromSecKey(secKey SecKey) Address {
	return AddressFromPubKey(MustPubKeyFromSecKey(secKey))
}

// DecodeBase58Address creates an Address from its base58 encoding
func DecodeBase58Address(addr string) (Address, error) {
	b, err := base58.Decode(addr)
	if err != nil {
		return Address{}, err
	}
	return AddressFromBytes(b)
}

// MustDecodeBase58Address creates an Address from its base58 encoding, panics on error
func MustDecodeBase58Address(addr string) Address {
	a, err := DecodeBase58Address(addr)
	if err != nil {
		log.Panicf("Invalid address %s: %v", addr, err)
	}
	return a
}

// AddressFromBytes converts []byte to an Address
func AddressFromBytes(b []byte) (Address, error) {
	if len(b) != 20+1+4 {
		return Address{}, ErrAddressInvalidLength
	}
	a := Address{}
	copy(a.Key[0:20], b[0:20])
	a.Version = b[20]

	chksum := a.Checksum()
	var checksum [4]byte
	copy(checksum[0:4], b[21:25])

	if checksum != chksum {
		return Address{}, ErrAddressInvalidChecksum
	}

	if a.Version != 0 {
		return Address{}, ErrAddressInvalidVersion
	}

	return a, nil
}

// MustAddressFromBytes converts []byte to an Address, panics on error
func MustAddressFromBytes(b []byte) Address {
	addr, err := AddressFromBytes(b)
	if err != nil {
		log.Panic(err)
	}

	return addr
}

// Null returns true if the address is null (0x0000....)
func (addr Address) Null() bool {
	return addr == Address{}
}

// Bytes return address as a byte slice
func (addr Address) Bytes() []byte {
	b := make([]byte, 20+1+4)
	copy(b[0:20], addr.Key[0:20])
	b[20] = addr.Version
	chksum := addr.Checksum()
	copy(b[21:25], chksum[0:4])
	return b
}

// Verify checks that the address appears valid for the public key
func (addr Address) Verify(pubKey PubKey) error {
	if addr.Version != 0x00 {
		return ErrAddressInvalidVersion
	}

	if addr.Key != PubKeyRipemd160(pubKey) {
		return ErrAddressInvalidPubKey
	}

	return nil
}

// String address as Base58 encoded string
func (addr Address) String() string {
	return string(base58.Encode(addr.Bytes()))
}

// Checksum returns Address Checksum which is the first 4 bytes of sha256(key+version)
func (addr Address) Checksum() Checksum {
	r1 := append(addr.Key[:], []byte{addr.Version}...)
	r2 := SumSHA256(r1[:])
	c := Checksum{}
	copy(c[:], r2[:len(c)])
	return c
}
