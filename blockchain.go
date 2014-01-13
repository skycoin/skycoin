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
func (self *Wallet) RefeshUnspentOutputs(bc sb_coin.Blockchain) {
	self.Outputs = new([]sb_coin.UxOuts)
	for _, a := range self.AA {
		//get unspent outputs for the address
		unspentOutputs := bc.GetUnspentOutputs(a.address)
		self.Outputs = append(self.Outputs, unspentOutputs)
	}
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

func (self *Address) GetOutputs(bc sb_coin.Blockchain) []sb_coin.UxOut {
	ux := bc.GetUnspentOutputs(*self)
	return ux
}

func tests() {

	genesisWallet := NewWallet(1)
	genesisAddress = genesisWallet.AA[0].address
	var BC *sb_coin.BlockChain = sb_coin.NewBlockChain(genesisAddress)

	genesisWallet.RefeshUnspentOutputs(BC)
	//create 16 wallets with 64 addresses
	wn := 16 //number of wallets to create
	WA := new([]Wallet)
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
		ti := sb_coin.TransactionInput
		ti.SigIdx = uint16(0)
		ti.UxOut = genesisWallet.Outputs[0].Hash()
		T.TI = append(T.TI, ti)

		to := sb_coin.TransactionOutput
		to.DestinationAddress = genesisWallet.AA[0].address
		to.Value1 = 100*1000000 - wn*1000
		T.TO = append(T.TO, to)

		for i := 0; i < wn; i++ {
			to := sb_coin.TransactionOutput
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
		T.PushOutput(genesisWallet.AA[0].address, 100*1000000-wn*1000, 0)
		for i := 0; i < wn; i++ {
			a := WA[i].GetRandom()
			T.PushOutput(A.address, wn, 1024)
		}

		var sec sb_coin.SecKey
		sec.Set(genesisAddress.seckey)
		T.SetSig(0, sec)
	}
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
