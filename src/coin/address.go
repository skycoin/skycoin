package coin

import (
    "bytes"
    "github.com/skycoin/skycoin/src/lib/base58"
    "log"
)

type Address struct {
    // ??
    Version byte
    // ripemd160 of sha256 of pubkey
    Key Ripemd160
    //CheckSum [4]byte
}

func (g *Address) String() string {
    return string(g.Base58())
}

// Returns address as raw bytes, containing version and then key
func (g *Address) Bytes() []byte {
    return append([]byte{g.Version}, g.Key[:]...)
}

// Returns address base58-encoded
func (g *Address) Base58() []byte {
    return []byte(base58.Hex2Base58(g.Key[:]))
}

func (g *Address) Equals(other *Address) bool {
    return g.Version == other.Version && bytes.Equal(g.Key[:], other.Key[:])
}

// Returns the address checksum
func (g *Address) Checksum() []byte {
    // TODO -- the comments here don't match the code and I have no idea
    // what this is supposed to be doing
    b := g.Bytes()
    //4 byte checksum
    r2 := SumSHA256(b)
    r3 := SumSHA256(r2[:])

    r4 := r3[:4] // first 1 bytes (error correction code)
    b2 := append(b[:], r4...)

    return b2
}

func (g *Address) MustChecksum() []byte {
    b := g.Checksum()
    if len(b) != 25 {
        log.Panic("Invalid address checksum")
    }
    return b
}

// Creates Address from PubKey
func AddressFromPubkey(pubkey PubKey) Address {
    s := SumSHA256(pubkey[:])
    addr := Address{Version: 0x0f, Key: HashRipemd160(s[:])}
    addr.MustChecksum()
    return addr
}

// Creates Address from []byte
func AddressFromRawPubkey(pubkeyraw []byte) Address {
    pubkey := NewPubKey(pubkeyraw)
    return AddressFromPubkey(pubkey)
}

/*
//set the genesis address pubkey
var GenesisAddress Address

func init() {
	a := "02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8"
	b, err := hex.DecodeString(a)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Printf("len= %v \n", len(b))
	GenesisAddress = AddressFromRawPubkey(b)
}
*/
