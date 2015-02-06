package cipher

import (
	"errors"
	"log"

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

type Checksum [4]byte

//version is after Key to enable better vanity address generation
//Address stuct is a 25 byte with a 20 byte publickey hash, 1 byte address
//type and 4 byte checksum.
type Address struct {
	Version byte      //1 byte
	Key     Ripemd160 //20 byte pubkey hash
}

// Creates Address from PubKey as ripemd160(sha256(sha256(pubkey)))
func AddressFromPubKey(pubKey PubKey) Address {
	addr := Address{
		Version: 0,
		Key:     pubKey.ToAddressHash(),
	}
	return addr
}

func AddressFromSecKey(secKey SecKey) Address {
	return AddressFromPubKey(PubKeyFromSecKey(secKey))
}

// Creates an Address from its base58 encoding.  Will panic if the addr is
// invalid
func MustDecodeBase58Address(addr string) Address {
	a, err := DecodeBase58Address(addr)
	if err != nil {
		log.Panicf("Invalid address %s: %v", addr, err)
	}
	return a
}

// Creates an Address from its base58 encoding
func DecodeBase58Address(addr string) (Address, error) {
	b, err := base58.Base582Hex(addr)
	if err != nil {
		return Address{}, err
	}
	return addressFromBytes(b)
}

// Returns an address given an Address.Bytes()
func addressFromBytes(b []byte) (Address, error) {
	if len(b) != 20+1+4 {
		return Address{}, errors.New("Invalid address bytes")
	}
	a := Address{}
	copy(a.Key[0:20], b[0:20])
	a.Version = b[20]
	if a.Version != 0 {
		return Address{}, errors.New("Invalid Version")
	}

	chksum := a.Checksum()
	var checksum [4]byte
	copy(checksum[0:4], b[21:25])

	if checksum != chksum {
		return Address{}, errors.New("Invalid checksum")
	}

	return a, nil
}

//return address as a byte slice
func (self *Address) Bytes() []byte {
	b := make([]byte, 20+1+4)
	copy(b[0:20], self.Key[0:20])
	b[20] = self.Version
	chksum := self.Checksum()
	copy(b[21:25], chksum[0:4])
	return b
}

// Checks that the address appears valid for the public key
func (self Address) Verify(key PubKey) error {
	if self.Version != 0x00 {
		return errors.New("Address version invalid")
	}
	if self.Key != key.ToAddressHash() {
		return errors.New("Public key invalid for address")
	}
	return nil
}

// Address as Base58 encoded string
// Returns address as printable
// version is first byte in binary format
// in printed address its key, version, checksum
func (self Address) String() string {
	return string(base58.Hex2Base58(self.Bytes()))
}

// Returns Address Checksum which is the first 4 bytes of sha256(key+version)
func (self *Address) Checksum() Checksum {
	// Version comes after the address to support vanity addresses
	r1 := append(self.Key[:], []byte{self.Version}...)
	r2 := SumSHA256(r1[:])
	c := Checksum{}
	copy(c[:], r2[:len(c)])
	return c
}

/*
Bitcoin Functions
*/

/*
//prints the bitcoin address for a seckey
func BitcoinAddressFromSeckey(seckey SecKey) string {

}

//exports seckey in wallet import format
//key must be compressed
func WalletImportFormat(seckey SecKey) string {

}

func MustSecKeyFromWalletImportFormat(intput string) SecKey {
	seckey, err := SecKeyFromWalletImportFormat(input)
	if err != nil {
		log.Panic("MustSecKeyFromWalletImportFormat, invalid seckey")
	}
	return seckey
}

//extracts a seckey from wallet import format
func SecKeyFromWalletImportFormat(input string) (SecKey, errors) {

}
*/
