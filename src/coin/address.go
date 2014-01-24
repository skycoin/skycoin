package coin

import (
    "fmt"
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

func (g Address) String() string {
    return string(AddressPrintable(g))
}

//get address struct from pubkey
func AddressFromPubkey(pubkey PubKey) Address {
    s := SumSHA256(pubkey[:])
    addr := Address{Version: 0x0f, Key: HashRipemd160(s[:])}
    // add version prefix
    b := append([]byte{addr.Version}, addr.Key[:]...)

    //4 byte checksum
    r2 := SumSHA256(b)
    r3 := SumSHA256(r2[:])

    r4 := r3[:4] //first 1 bytes (error correction code)
    b2 := append(b[:], r4...)

    // TODO -- b2 is never used. what is it supposed to do

    if len(b2) != 25 {
        fmt.Printf("len(b)= %v, len(b2)= %v, len(r)= %v, len(r4)= %v \n",
            len(b), len(b2), len(addr.Key), len(r4))
        log.Panic("Invalid b2 length")
    }

    return addr
}

func AddressFromRawPubkey(pubkeyraw []byte) Address {
    pubkey := NewPubKey(pubkeyraw)
    return AddressFromPubkey(pubkey)
}

//returns base 58 of Address
func AddressPrintable(a Address) []byte {
    b1 := append([]byte{a.Version}, a.Key[:]...) //add version prefix

    r1 := SumSHA256(b1)
    r2 := SumSHA256(r1[:])
    r3 := r2[:4] // 4 bytes error correction code
    b2 := append(b1[:], r3...)

    // TODO -- b2 is never used. what is it supposed to do

    if len(b2) != 25 {
        log.Panic("Invalid b2 len")
    }
    var en base58.Base58 = base58.Hex2Base58(a.Key[:]) //encode as base 58
    //fmt.Printf("address= %v\n", en)
    return []byte(en)
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
