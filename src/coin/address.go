package coin

import (
    "bytes"
    "encoding/hex"
    "github.com/skycoin/skycoin/src/lib/base58"
    "log"
)

type Address struct {
    Version byte
    Key     [20]byte //sha256(sha256(ridmd160(pubkey)))
    ChkSum  [4]byte
}

//get address as string
func (g *Address) String() string {
    return string(g.Base58())
}

// Returns address as raw bytes, containing version and then key
func (g *Address) Bytes() []byte {
    b := make([]byte, 25)
    b[0] = g.Version
    copy(b[1:21], g.Key[0:20])
    copy(b[21:25], g.ChkSum[0:4])
    return b
    //return append([]byte{g.Version}, g.Key[:]...)
}

// Returns address base58-encoded
func (g *Address) Base58() []byte {
    return []byte(base58.Hex2Base58(g.Key[:]))
}

func (g Address) Equals(other Address) bool {
    return g == other
}

// Returns address checksum
// 4 byte checksum
func (g *Address) Checksum() []byte {
    r1 := append([]byte{g.Version}, g.Key[:]...)
    r2 := SumSHA256(r1[:])
    return r2[0:4] //4 bytes
}

func (g *Address) SetChecksum() {
    copy(g.ChkSum[0:4], g.Checksum())
}

//r3 := SumSHA256(r2[:])
//r4 := HashRipemd160(r3[:])

func (g *Address) ChecksumVerify() int {
    chksum := g.Checksum()
    if len(chksum) != 4 {
        log.Panic("Invalid address checksum")
    }
    if !bytes.Equal(chksum[0:4], g.ChkSum[0:4]) {
        return 0
    }

    return 1
}

// Creates Address from PubKey
// sha256(sha256(ridmd160(pubkey)))
func AddressFromPubKey(pubkey PubKey) Address {
    r1 := SumSHA256(pubkey[:])
    r2 := SumSHA256(r1[:])
    r3 := HashRipemd160(r2[:])
    addr := Address{Version: 0x0f, Key: r3}
    addr.SetChecksum()
    return addr
}

// Creates Address from []byte
func AddressFromRawPubKey(pubkeyraw []byte) Address {
    pubkey := NewPubKey(pubkeyraw)
    return AddressFromPubKey(pubkey)
}

func init() {
    a := "02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8"
    b, err := hex.DecodeString(a)
    if err != nil {
        log.Panic(err)
    }
    addr := AddressFromRawPubKey(b)
    if !bytes.Equal(addr.ChkSum[:], addr.Checksum()) {
        log.Panic("Checksum fail")
    }
    if len(addr.Bytes()) != 25 {
        log.Panic("Address length invalid")
    }
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
