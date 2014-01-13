package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
)

import (
	"./src/coin/"
	//"./src/util"
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
	AA      []Address       //address array
	Outputs []sb_coin.UxOut //unspent outputs
}

func (self *Wallet) GetRandom() Address {
	index := rand.Int() % len(self.AA)
	return self.AA[index]
}

//signs a hash for a given address. nil on failure
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

//refresh the unspent outputs for the wallet
func (self *Wallet) RefeshUnspentOutputs(bc *sb_coin.BlockChain) {
	var outputs []sb_coin.UxOut
	for _, a := range self.AA {
		//get unspent outputs for the address
		unspentOutputs := bc.GetUnspentOutputs(a.address)
		outputs = append(outputs, unspentOutputs...)
	}
	self.Outputs = outputs
}

func NewWallet(n int) Wallet {
	var w Wallet
	for i := 0; i < n; i++ {
		w.AA = append(w.AA, GenerateAddress())
	}
	return w
}

type Address struct {
	pubkey  []byte
	seckey  []byte
	address sb_coin.Address
}

func GenerateAddress() Address {
	var A Address
	A.pubkey, A.seckey = sb_coin.GenerateKeyPair()
	A.address = sb_coin.AddressFromRawPubkey(A.pubkey)
	return A
}

//func (self *Address) GetOutputs(bc sb_coin.BlockChain) []sb_coin.UxOut {
//	ux := bc.GetUnspentOutputs(*self.address)
//	return ux
//}

func tests() {

	genesisWallet := NewWallet(1)
	genesisAddress := genesisWallet.AA[0]
	var BC *sb_coin.BlockChain = sb_coin.NewBlockChain(genesisAddress.address)

	genesisWallet.RefeshUnspentOutputs(BC)
	//create 16 wallets with 64 addresses
	wn := 16 //number of wallets to create
	var WA []Wallet
	for i := 0; i < wn; i++ {
		WA = append(WA, NewWallet(64))
	}

	//pub, sec := util.GenerateKeyPair()
	//pub := FromHex("02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8")
	//sec := FromHex("f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0")
	//fmt.Printf("sec: %s \n", ToHex(sec))
	//fmt.Printf("pub: %s \n", ToHex(pub))

	B := BC.NewBlock()

	var T sb_coin.Transaction

	if true {
		//create transaction by hand
		var ti sb_coin.TransactionInput
		ti.SigIdx = uint16(0)
		ti.UxOut = genesisWallet.Outputs[0].Hash()
		T.TI = append(T.TI, ti)

		var to sb_coin.TransactionOutput
		to.DestinationAddress = genesisWallet.AA[0].address
		to.Value1 = uint64(100*1000000 - wn*1000)
		T.TO = append(T.TO, to)

		for i := 0; i < wn; i++ {
			var to sb_coin.TransactionOutput
			a := WA[i].GetRandom()
			to.DestinationAddress = a.address
			to.Value1 = 1000
			to.Value2 = 1024 * 1024
			T.TO = append(T.TO, to)
		}

		var sec sb_coin.SecKey
		sec.Set(genesisAddress.seckey)
		T.SetSig(0, sec)

	} else {
		T.PushInput(genesisWallet.Outputs[0].Hash())
		T.PushOutput(genesisWallet.AA[0].address, uint64(100*1000000-wn*1000), 0)
		for i := 0; i < wn; i++ {
			a := WA[i].GetRandom()
			T.PushOutput(a.address, uint64(wn), 1024)
		}

		var sec sb_coin.SecKey
		sec.Set(genesisAddress.seckey)
		T.SetSig(0, sec)
	}
	T.UpdateHeader() //sets hash

	fmt.Printf("genesis transaction: \n")
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
