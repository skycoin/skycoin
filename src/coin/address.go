package sb_coin

import (
	//"encoding/hex"
	"fmt"
	"log"
)

import "lib/base58"

type Address struct {
	Version byte
	Value   [20]byte //ripemd160 of sha256 of pubkey
	//CheckSum [4]byte
}

func (g Address) Print() []byte {
	return AddressPrintable(g)
}

func (g Address) String() string {
	return string(AddressPrintable(g))
}

//get address struct from pubkey
func AddressFromPubkey(pubkey PubKey) Address {
	var ret Address
	ret.Version = 0x0f

	if len(pubkey.Value) != 33 {
		fmt.Printf("len= %v \n", len(pubkey.Value))
		log.Panic()
	}
	s := Sha256_func(pubkey.Value[0:33])
	r := Ripmd160_func(s[:])
	copy(ret.Value[0:20], r[0:20])

	b := append([]byte{ret.Version}, r[:]...) //add version prefix

	//4 byte checksum
	r2 := Sha256_func(b)
	r3 := Sha256_func(r2[:])

	r4 := r3[0:4] //first 1 bytes (error correction code)
	b2 := append(b[:], r4...)

	if len(b2) != 25 {
		fmt.Printf("len(b)= %v, len(b2)= %v, len(r)= %v, len(r4)= %v \n", len(b), len(b2), len(r), len(r4))
		log.Panic()
	}

	return ret
}

func AddressFromRawPubkey(pubkeyraw []byte) Address {
	var pubkey PubKey
	pubkey.Set(pubkeyraw)
	return AddressFromPubkey(pubkey)
}

//returns base 58 of Address
func AddressPrintable(a Address) []byte {
	b1 := append([]byte{a.Version}, a.Value[0:20]...) //add version prefix

	r1 := Sha256_func(b1)
	r2 := Sha256_func(r1[:])
	r3 := r2[0:4] // 4 bytes error correction code
	b2 := append(b1[:], r3...)

	if len(b2) != 25 {
		log.Panic()
	}
	var en base58.Base58 = base58.Hex2Base58(a.Value[:]) //encode as base 58
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
