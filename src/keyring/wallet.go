package keyring

import (
	"encoding/hex"
	//"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/coin"
	"log"
	"math/rand"
)

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

type Address struct {
	Pubkey  []byte
	Address coin.Address
	seckey  []byte //keep secret
}

func GenerateAddress() Address {
	var A Address
	A.Pubkey, A.seckey = coin.GenerateKeyPair()
	A.Address = coin.AddressFromRawPubkey(A.Pubkey)
	return A
}


/*
   genesis address
   pub: 02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8
   sec: f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0

*/

type Wallet struct {
	AA      []Address    //address array
	Outputs []coin.UxOut //unspent outputs
}

func (self *Wallet) GetRandom() Address {
	index := rand.Int() % len(self.AA)
	return self.AA[index]
}

//signs a hash for a given address. nil on failure
func (self *Wallet) Sign(address coin.Address, hash []byte) []byte {
	for _, a := range self.AA {
		if a.Address == address {
			//func GenerateSignature(seckey []byte, msg []byte) []byte
			sig := coin.GenerateSignature(a.seckey, hash)
			return sig
		}

	}
	return nil
}

//refresh the unspent outputs for the wallet
func (self *Wallet) RefeshUnspentOutputs(bc *coin.BlockChain) {
	var outputs []coin.UxOut
	for _, a := range self.AA {
		//get unspent outputs for the address
		unspentOutputs := bc.GetUnspentOutputs(a.Address)
		outputs = append(outputs, unspentOutputs...)
	}
	self.Outputs = outputs
}

func (self *Wallet) Balance(bc *coin.BlockChain) (uint64, uint64) {
	self.RefeshUnspentOutputs(bc)

	var balance1 uint64 = 0
	var balance2 uint64 = 0

	var time uint64 = bc.Head.Header.Time

	for _, ux := range self.Outputs {
		balance1 += ux.Body.Coins
		balance2 += ux.CoinHours(time)
	}
	return balance1, balance2
}

/*
func (self *Wallet) NewTransaction(bc *coin.BlockChain, Address coin.Address, amt1 uint64, amt2 uint64) (coin.Transaction, error) {
	self.RefeshUnspentOutputs(bc)
	bal1, bal2 := self.Balance()

	if bal1 < amt1 {
		return coin.Transaction{}, errors.New("insufficient coin balance")
	}
	if bal2 < amt2 {
		return coin.Transaction{}, errors.New("insufficient coinhour balance")
	}

	//decide which outputs get spent

	var ti coin.TransactionInput
	ti.SigIdx = uint16(0)
	ti.UxOut = genesisWallet.Outputs[0].Hash()
	T.TI = append(T.TI, ti)

	var to coin.TransactionOutput
	to.DestinationAddress = genesisWallet.AA[0].Address
	to.Value1 = uint64(100*1000000 - wn*1000)
	T.TO = append(T.TO, to)

	for i := 0; i < wn; i++ {
		var to coin.TransactionOutput
		a := WA[i].GetRandom()
		to.DestinationAddress = a.Address
		to.Value1 = 1000
		to.Value2 = 1024
		T.TO = append(T.TO, to)
	}

	var sec coin.SecKey
	sec.Set(genesisAddress.seckey)
	T.SetSig(0, sec)

}
*/

func NewWallet(n int) Wallet {
	var w Wallet
	for i := 0; i < n; i++ {
		w.AA = append(w.AA, GenerateAddress())
	}
	return w
}

//func (self *Address) GetOutputs(bc coin.BlockChain) []coin.UxOut {
//	ux := bc.GetUnspentOutputs(*self.Address)
//	return ux
//}

//pub, sec := util.GenerateKeyPair()
//pub := FromHex("02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8")
//sec := FromHex("f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0")
//fmt.Printf("sec: %s \n", ToHex(sec))
//fmt.Printf("pub: %s \n", ToHex(pub))

func WalletBalances(bc *coin.BlockChain, wallets []Wallet) {
	for i, w := range wallets {
		b1, b2 := w.Balance(bc)
		fmt.Printf("%v: %v %v \n", i, b1, b2)
	}

	for i, ux := range bc.Unspent {
		fmt.Printf("%v: %v \n", i, ux.String())
	}
}
