package coin

import (
    "bytes"
    "github.com/skycoin/skycoin/src/lib/base58"
)

const (
    mainAddressVersion = 0x0f
    testAddressVersion = 0xB0
)

type Checksum [4]byte

//version is after Key to enable better vanity address generation
type Address struct {
    Key     Ripemd160 //sha256(sha256(ridmd160(pubkey)))
    Version byte
    ChkSum  Checksum
}

// Creates Address from PubKey as
// sha256(sha256(ridmd160(pubkey)))
func AddressFromPubKey(pubkey PubKey) Address {
    // WARNING - DOES NOT MATCH DOCSTRING
    // THIS DOES ripemd160(sha256(sha256(pubkey)))
    r1 := SumSHA256(pubkey[:])
    r2 := SumSHA256(r1[:])
    r3 := HashRipemd160(r2[:])
    addr := Address{
        Version: mainAddressVersion,
        Key:     r3,
    }
    addr.setChecksum()
    return addr
}

// Creates an address for the test network
func AddressFromPubkeyTestNet(pubKey PubKey) Address {
    a := AddressFromPubKey(pubKey)
    a.Version = testAddressVersion
    a.setChecksum()
    return a
}

// Address as Base58 encoded string
func (self *Address) String() string {
    return string(base58.Hex2Base58(self.Key[:]))
}

// Returns address as raw bytes, containing version and then key
func (self *Address) Bytes() []byte {
    keyLen := len(self.Key)
    b := make([]byte, keyLen+len(self.ChkSum)+1)
    copy(b[:keyLen], self.Key[:keyLen])
    b[keyLen] = self.Version
    copy(b[keyLen+1:], self.ChkSum[:])
    return b
}

// Returns Address Checksum
func (self *Address) Checksum() Checksum {
    r1 := append(self.Key[:], []byte{self.Version}...)
    r2 := SumSHA256(r1[:])
    var c Checksum
    copy(c[:], r2[:len(c)])
    return c
}

// Returns whether the checksum on address is valid for its key
func (self *Address) IsValidChecksum() bool {
    c := self.Checksum()
    return bytes.Equal(c[:], self.ChkSum[:])
}

func (self *Address) setChecksum() {
    self.ChkSum = self.Checksum()
}
