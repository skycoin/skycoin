package cipher

import (
	"bytes"
	"errors"
	"log"

	"github.com/SkycoinProject/skycoin/src/cipher/base58"
)

var (
	// ErrInvalidLength Unexpected size of string or bytes buffer
	ErrInvalidLength = errors.New("Invalid length")
	// ErrBitcoinWIFInvalidFirstByte Unexpected value (!= 0x80) of first byte in Bitcoin Wallet Import Format
	ErrBitcoinWIFInvalidFirstByte = errors.New("Bitcoin WIF: First byte invalid")
	// ErrBitcoinWIFInvalidSuffix Unexpected value (!= 0x01) of 33rd byte in Bitcoin Wallet Import Format
	ErrBitcoinWIFInvalidSuffix = errors.New("Bitcoin WIF: Invalid 33rd byte")
	// ErrBitcoinWIFInvalidChecksum Invalid Checksum in Bitcoin WIF address
	ErrBitcoinWIFInvalidChecksum = errors.New("Bitcoin WIF: Checksum fail")
)

// BitcoinAddress is a bitcoin address
type BitcoinAddress struct {
	Version byte      // 1 byte
	Key     Ripemd160 // 20 byte pubkey hash
}

// BitcoinPubKeyRipemd160 returns ripemd160(sha256(key))
func BitcoinPubKeyRipemd160(pubKey PubKey) Ripemd160 {
	r1 := SumSHA256(pubKey[:])
	return HashRipemd160(r1[:])
}

// BitcoinAddressFromPubKey creates a mainnet (version 0) BitcoinAddress from PubKey as ripemd160(sha256(pubkey)))
func BitcoinAddressFromPubKey(pubKey PubKey) BitcoinAddress {
	return BitcoinAddress{
		Version: 0,
		Key:     BitcoinPubKeyRipemd160(pubKey),
	}
}

// BitcoinAddressFromSecKey generates a BitcoinAddress from SecKey
func BitcoinAddressFromSecKey(secKey SecKey) (BitcoinAddress, error) {
	p, err := PubKeyFromSecKey(secKey)
	if err != nil {
		return BitcoinAddress{}, err
	}
	return BitcoinAddressFromPubKey(p), nil
}

// MustBitcoinAddressFromSecKey generates a BitcoinAddress from SecKey, panics on error
func MustBitcoinAddressFromSecKey(secKey SecKey) BitcoinAddress {
	return BitcoinAddressFromPubKey(MustPubKeyFromSecKey(secKey))
}

// DecodeBase58BitcoinAddress creates a BitcoinAddress from its base58 encoding
func DecodeBase58BitcoinAddress(addr string) (BitcoinAddress, error) {
	b, err := base58.Decode(addr)
	if err != nil {
		return BitcoinAddress{}, err
	}
	return BitcoinAddressFromBytes(b)
}

// MustDecodeBase58BitcoinAddress creates a BitcoinAddress from its base58 encoding, panics on error
func MustDecodeBase58BitcoinAddress(addr string) BitcoinAddress {
	a, err := DecodeBase58BitcoinAddress(addr)
	if err != nil {
		log.Panicf("Invalid bitcoin address %s: %v", addr, err)
	}
	return a
}

// BitcoinAddressFromBytes converts []byte to a BitcoinAddress. Only supports mainnet (version 0) addresses.
func BitcoinAddressFromBytes(b []byte) (BitcoinAddress, error) {
	if len(b) != 20+1+4 {
		return BitcoinAddress{}, ErrAddressInvalidLength
	}
	a := BitcoinAddress{}
	copy(a.Key[0:20], b[1:21])
	a.Version = b[0]

	var checksum [4]byte
	copy(checksum[0:4], b[21:25])

	if checksum != a.Checksum() {
		return BitcoinAddress{}, ErrAddressInvalidChecksum
	}

	// BitcoinAddress only supports mainnet addresses for now
	if a.Version != 0 {
		return BitcoinAddress{}, ErrAddressInvalidVersion
	}

	return a, nil
}

// MustBitcoinAddressFromBytes converts []byte to a BitcoinAddress, panics on error
func MustBitcoinAddressFromBytes(b []byte) BitcoinAddress {
	addr, err := BitcoinAddressFromBytes(b)
	if err != nil {
		log.Panic(err)
	}

	return addr
}

// Null returns true if the address is null (0x0000....)
func (addr BitcoinAddress) Null() bool {
	return addr == BitcoinAddress{}
}

// Bytes returns bitcoin address as byte slice
func (addr BitcoinAddress) Bytes() []byte {
	b := make([]byte, 20+1+4)
	b[0] = addr.Version
	copy(b[1:21], addr.Key[0:20])
	chksum := addr.Checksum()
	copy(b[21:25], chksum[0:4])
	return b
}

// Verify checks that the bitcoin address appears valid for the public key
func (addr BitcoinAddress) Verify(key PubKey) error {
	// BitcoinAddress only supports mainnet addresses for now
	if addr.Version != 0x00 {
		return ErrAddressInvalidVersion
	}
	if addr.Key != BitcoinPubKeyRipemd160(key) {
		return ErrAddressInvalidPubKey
	}
	return nil
}

// String convert bitcoin address to hex string
func (addr BitcoinAddress) String() string {
	return string(base58.Encode(addr.Bytes()))
}

// Checksum returns a bitcoin address Checksum which is the first 4 bytes of sha256(sha256(version+key))
func (addr BitcoinAddress) Checksum() Checksum {
	r1 := append([]byte{addr.Version}, addr.Key[:]...)
	r2 := DoubleSHA256(r1[:])
	c := Checksum{}
	copy(c[:], r2[:len(c)])
	return c
}

// BitcoinWalletImportFormatFromSeckey exports seckey in wallet import format
// key must be compressed
func BitcoinWalletImportFormatFromSeckey(seckey SecKey) string {
	b1 := append([]byte{byte(0x80)}, seckey[:]...)
	b2 := append(b1[:], []byte{0x01}...)
	b3 := DoubleSHA256(b2) //checksum
	b4 := append(b2, b3[0:4]...)
	return string(base58.Encode(b4))
}

// SecKeyFromBitcoinWalletImportFormat extracts a seckey from the bitcoin wallet import format
func SecKeyFromBitcoinWalletImportFormat(input string) (SecKey, error) {
	b, err := base58.Decode(input)
	if err != nil {
		return SecKey{}, err
	}

	//1+32+1+4
	if len(b) != 38 {
		return SecKey{}, ErrInvalidLength
	}
	if b[0] != 0x80 {
		return SecKey{}, ErrBitcoinWIFInvalidFirstByte
	}

	if b[1+32] != 0x01 {
		return SecKey{}, ErrBitcoinWIFInvalidSuffix
	}

	b2 := DoubleSHA256(b[0:34])
	chksum := b[34:38]

	if !bytes.Equal(chksum, b2[0:4]) {
		return SecKey{}, ErrBitcoinWIFInvalidChecksum
	}

	return NewSecKey(b[1:33])
}

// MustSecKeyFromBitcoinWalletImportFormat extracts a seckey from the bitcoin wallet import format, panics on error
func MustSecKeyFromBitcoinWalletImportFormat(input string) SecKey {
	seckey, err := SecKeyFromBitcoinWalletImportFormat(input)
	if err != nil {
		log.Panicf("MustSecKeyFromBitcoinWalletImportFormat, invalid seckey, %v", err)
	}
	return seckey
}
