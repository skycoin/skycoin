package cipher

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher/base58"
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

// Address version is after Key to enable better vanity address generation
// Address stuct is a 25 byte with a 20 byte publickey hash, 1 byte address
// type and 4 byte checksum.
type Address struct {
	Version byte      //1 byte
	Key     Ripemd160 //20 byte pubkey hash
}

// AddressFromPubKey creates Address from PubKey as ripemd160(sha256(sha256(pubkey)))
func AddressFromPubKey(pubKey PubKey) Address {
	addr := Address{
		Version: 0,
		Key:     pubKey.ToAddressHash(),
	}
	return addr
}

// AddressFromSecKey generates address from secret key
func AddressFromSecKey(secKey SecKey) Address {
	return AddressFromPubKey(PubKeyFromSecKey(secKey))
}

// DecodeBase58Address creates an Address from its base58 encoding
func DecodeBase58Address(addr string) (Address, error) {
	b, err := base58.Base582Hex(addr)
	if err != nil {
		return Address{}, err
	}
	return addressFromBytes(b)
}

// MustDecodeBase58Address creates an Address from its base58 encoding.  Will panic if the addr is
// invalid
func MustDecodeBase58Address(addr string) Address {
	a, err := DecodeBase58Address(addr)
	if err != nil {
		logger.Panicf("Invalid address %s: %v", addr, err)
	}
	return a
}

// BitcoinDecodeBase58Address decode bitcoin address from string
func BitcoinDecodeBase58Address(addr string) (Address, error) {
	b, err := base58.Base582Hex(addr)
	if err != nil {
		return Address{}, err
	}
	return BitcoinAddressFromBytes(b)
}

// BitcoinMustDecodeBase58Address must decodes bitcoin address from string
func BitcoinMustDecodeBase58Address(addr string) Address {
	a, err := BitcoinDecodeBase58Address(addr)
	if err != nil {
		logger.Panicf("Invalid address %s: %v", addr, err)
	}
	return a
}

// Returns an address given an Address.Bytes()
func addressFromBytes(b []byte) (addr Address, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if len(b) != 20+1+4 {
		return Address{}, errors.New("Invalid address length")
	}
	a := Address{}
	copy(a.Key[0:20], b[0:20])
	a.Version = b[20]
	if a.Version != 0 {
		return Address{}, errors.New("Invalid version")
	}

	chksum := a.Checksum()
	var checksum [4]byte
	copy(checksum[0:4], b[21:25])

	if checksum != chksum {
		return Address{}, errors.New("Invalid checksum")
	}

	return a, nil
}

// Bytes return address as a byte slice
func (addr *Address) Bytes() []byte {
	b := make([]byte, 20+1+4)
	copy(b[0:20], addr.Key[0:20])
	b[20] = addr.Version
	chksum := addr.Checksum()
	copy(b[21:25], chksum[0:4])
	return b
}

// BitcoinBytes returns bitcoin address as byte slice
func (addr *Address) BitcoinBytes() []byte {
	b := make([]byte, 20+1+4)
	b[0] = addr.Version
	copy(b[1:21], addr.Key[0:20])
	// b[20] = self.Version
	chksum := addr.BitcoinChecksum()
	copy(b[21:25], chksum[0:4])
	return b
}

// Verify checks that the address appears valid for the public key
func (addr Address) Verify(key PubKey) error {
	if addr.Version != 0x00 {
		return errors.New("Address version invalid")
	}
	if addr.Key != key.ToAddressHash() {
		return errors.New("Public key invalid for address")
	}
	return nil
}

// String address as Base58 encoded string
// Returns address as printable
// version is first byte in binary format
// in printed address its key, version, checksum
func (addr Address) String() string {
	return string(base58.Hex2Base58(addr.Bytes()))
}

// BitcoinString convert bitcoin address to hex string
func (addr Address) BitcoinString() string {
	return string(base58.Hex2Base58(addr.BitcoinBytes()))
}

// Checksum returns Address Checksum which is the first 4 bytes of sha256(key+version)
func (addr *Address) Checksum() Checksum {
	// Version comes after the address to support vanity addresses
	r1 := append(addr.Key[:], []byte{addr.Version}...)
	r2 := SumSHA256(r1[:])
	c := Checksum{}
	copy(c[:], r2[:len(c)])
	return c
}

// BitcoinChecksum bitcoin checksum
func (addr *Address) BitcoinChecksum() Checksum {
	// Version comes after the address to support vanity addresses
	r1 := append([]byte{addr.Version}, addr.Key[:]...)
	r2 := DoubleSHA256(r1[:])
	c := Checksum{}
	copy(c[:], r2[:len(c)])
	return c
}

/*
Bitcoin Functions
*/

// BitcoinAddressFromPubkey prints the bitcoin address for a seckey
func BitcoinAddressFromPubkey(pubkey PubKey) string {
	b1 := SumSHA256(pubkey[:])
	b2 := HashRipemd160(b1[:])
	b3 := append([]byte{byte(0)}, b2[:]...)
	b4 := DoubleSHA256(b3)
	b5 := append(b3, b4[0:4]...)
	return string(base58.Hex2Base58(b5))
	// return Address{
	// 	Version: 0,
	// 	Key:     b2,
	// }
}

// BitcoinWalletImportFormatFromSeckey exports seckey in wallet import format
// key must be compressed
func BitcoinWalletImportFormatFromSeckey(seckey SecKey) string {
	b1 := append([]byte{byte(0x80)}, seckey[:]...)
	b2 := append(b1[:], []byte{0x01}...)
	b3 := DoubleSHA256(b2) //checksum
	b4 := append(b2, b3[0:4]...)
	return string(base58.Hex2Base58(b4))
}

// BitcoinAddressFromBytes Returns an address given an Address.Bytes()
func BitcoinAddressFromBytes(b []byte) (Address, error) {
	if len(b) != 20+1+4 {
		return Address{}, errors.New("Invalid address length")
	}
	a := Address{}
	copy(a.Key[0:20], b[1:21])
	a.Version = b[0]
	if a.Version != 0 {
		return Address{}, errors.New("Invalid version")
	}

	chksum := a.BitcoinChecksum()
	var checksum [4]byte
	copy(checksum[0:4], b[21:25])

	if checksum != chksum {
		return Address{}, errors.New("Invalid checksum")
	}

	return a, nil
}

// SecKeyFromWalletImportFormat extracts a seckey from wallet import format
func SecKeyFromWalletImportFormat(input string) (SecKey, error) {
	b, err := base58.Base582Hex(input)
	if err != nil {
		return SecKey{}, err
	}

	//1+32+1+4
	if len(b) != 38 {
		//log.Printf("len= %v ", len(b))
		return SecKey{}, errors.New("invalid length")
	}
	if b[0] != 0x80 {
		return SecKey{}, errors.New("first byte invalid")
	}

	if b[1+32] != 0x01 {
		return SecKey{}, errors.New("invalid 33rd byte")
	}

	b2 := DoubleSHA256(b[0:34])
	chksum := b[34:38]

	if !bytes.Equal(chksum, b2[0:4]) {
		return SecKey{}, errors.New("checksum fail")
	}

	seckey := b[1:33]
	if len(seckey) != 32 {
		logger.Panic("...")
	}
	return NewSecKey(b[1:33]), nil
}

// MustSecKeyFromWalletImportFormat SecKeyFromWalletImportFormat or panic
func MustSecKeyFromWalletImportFormat(input string) SecKey {
	seckey, err := SecKeyFromWalletImportFormat(input)
	if err != nil {
		logger.Panicf("MustSecKeyFromWalletImportFormat, invalid seckey, %v", err)
	}
	return seckey
}
