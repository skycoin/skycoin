package hashchain

import (
	"bytes"
	"errors"
	"github.com/skycoin/skywire/src/lib/base58"
	"log"
)

//Move address version checks to visor/transaction validation?
var (
	addressVersions = map[string]byte{
		"main":    0x0F, //main network
		"obelisk": 0x10, // obelisk node
		"test":    0x1F, //test network

	}
	// Address version is a global default version used for all address
	// creation and checking
	addressVersion = addressVersions["obelisk"]
)

// Returns the named address version and whether it is a known version
func VersionByName(name string) (byte, bool) {
	v, ok := addressVersions[name]
	return v, ok
}

// Returns the named address version, panics if unknown name
func MustVersionByName(name string) byte {
	v, ok := VersionByName(name)
	if !ok {
		log.Panicf("Invalid version name: %s", name)
	}
	return v
}

// Sets the address version used for all address creation and checking
func SetAddressVersion(name string) {
	addressVersion = MustVersionByName(name)
	//logger.Debug("Set address version to %s", name)
}

type Checksum [4]byte

//version is after Key to enable better vanity address generation
//Address stuct is a 25 byte with a 20 byte publickey hash, 1 byte address
//type and 4 byte checksum.
type Address struct {
	Key      Ripemd160 //20 byte pubkey hash
	Version  byte      //1 byte
	Checksum Checksum  //4 byte checksum, first 4 bytes of sha256 of key+version
}

// Creates Address from PubKey as ripemd160(sha256(sha256(pubkey)))
func AddressFromPubKey(pubKey PubKey) Address {
	addr := Address{
		Version: addressVersion,
		Key:     pubKey.ToAddressHash(),
	}
	addr.setChecksum()
	return addr
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
	a := Address{}
	keyLen := len(a.Key)
	if len(b) != keyLen+len(a.Checksum)+1 {
		return a, errors.New("Invalid address bytes")
	}
	copy(a.Key[:], b[:keyLen])
	a.Version = b[keyLen]
	copy(a.Checksum[:], b[keyLen+1:])
	if !a.HasValidChecksum() {
		return a, errors.New("Invalid checksum")
	} else {
		return a, nil
	}
}

// Checks that the address appears valid for the public key
func (self *Address) Verify(key PubKey) error {
	if self.Key != key.ToAddressHash() {
		return errors.New("Public key invalid for address")
	}
	if self.Version != addressVersion {
		return errors.New("Invalid address version")
	}
	if !self.HasValidChecksum() {
		return errors.New("Invalid address checksum")
	}
	return nil
}

// Address as Base58 encoded string
func (self *Address) String() string {
	return string(base58.Hex2Base58(self.Bytes()))
}

// Returns address as raw bytes, containing version and then key
func (self *Address) Bytes() []byte {
	keyLen := len(self.Key)
	b := make([]byte, keyLen+len(self.Checksum)+1)
	copy(b[:keyLen], self.Key[:])
	b[keyLen] = self.Version
	copy(b[keyLen+1:], self.Checksum[:])
	return b
}

// Returns Address Checksum which is the first 4 bytes of sha256(key+version)
func (self *Address) CreateChecksum() Checksum {
	// Version comes after the address to support vanity addresses
	r1 := append(self.Key[:], []byte{self.Version}...)
	r2 := SumSHA256(r1[:])
	c := Checksum{}
	copy(c[:], r2[:len(c)])
	return c
}

// Returns whether the checksum on address is valid for its key
func (self *Address) HasValidChecksum() bool {
	c := self.CreateChecksum()
	return bytes.Equal(c[:], self.Checksum[:])
}

func (self *Address) setChecksum() {
	self.Checksum = self.CreateChecksum()
}
