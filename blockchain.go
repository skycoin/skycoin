package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
)

import (
	//"./src/cli/"
	"./src/coin/"
	"./src/util"
	//"./src/daemon/"
	//"./src/gui/"
)

func Run() {
	tests()
}

func ToHex(b []byte) string {
	return hex.EncodeToString(b)
}

func FromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return b
}

/*
   genesis address
   pub: 02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8
   sec: f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0

*/

type Wallet struct {
	AA []Address
}

func (self *Wallet) GetRandom() Address {
	index := rand.Int() % len(self.AA)
	return self.AA[index]
}

func (self *Wallet) Sign(address sb_coin.Address, hash []byte) []byte {
	for _, a := range self.AA {
		if a.address == address {
			//func GenerateSignature(seckey []byte, msg []byte) []byte
			sig := sb_coin.GenerateSignature(a.seckey, hash)
			return sig
		}

	}
	return nil
}

func NewWallet(int n) Wallet {
	for i := 0; i < n; i++ {
		self.AA = append(self.AA, GenerateAddress())
	}
}

type Address struct {
	pubkey  []byte
	seckey  []byte
	address sb_coin.Address
}

func GenerateAddress() Address {
	var A Address
	A.pub, A.sec = util.GenerateKeyPair()
	A.add = sb_coin.AddressFromRawPubkey(A.pub)
	return A
}

func AddressArray(n int) []Address {
	aa := make([]Address, n)
	for i := 0; i < n; i++ {
		aa[i] = GenerateAddress()
	}
	return aa
}

func tests() {

	aa := AddressArray(16)
	_ = aa

	pub, sec := util.GenerateKeyPair()
	pub = FromHex("02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8")
	sec = FromHex("f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0")

	fmt.Printf("sec: %s \n", ToHex(sec))
	fmt.Printf("pub: %s \n", ToHex(pub))

	var BC *sb_coin.BlockChain = sb_coin.NewBlockChain()

	if false {
		fmt.Printf("l= %v\n", len(BC.Blocks))
	}

	B := BC.NewBlock()

	var T sb_coin.Transaction

	T.UpdateHeader() //sets hash

	/*
	   Need to add input
	*/
	err := BC.AppendTransaction(B, &T)

	if err != nil {
		log.Panic(err)
	}

}

func main() {
	tests()
}
