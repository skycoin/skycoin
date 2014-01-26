package main

import (
	"encoding/hex"
	//"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/keyring"

	"log"
	"math/rand"
)

/*
   Transaction set
*/
type PendingTransactions struct {
	A []coin.Transaction
}

func (self *PendingTransactions) Push(T coin.Transaction) {
	self.A = append(self.A, T)
}

func (self *PendingTransactions) Rand() coin.Transaction {
	i := rand.Int() % len(self.A) //random element
	element := self.A[i]
	self.A = append(self.A[:i], self.A[i+1:]...) //remove ith element
	return element
}

func tests() {

	genesisWallet := keyring.NewWallet(1)
	genesisAddress := genesisWallet.AA[0]
	var BC *coin.BlockChain = coin.NewBlockChain(genesisAddress.address)

	genesisWallet.RefeshUnspentOutputs(BC)
	//create 16 wallets with 64 addresses
	wn := 16 //number of wallets to create
	var WA []keyring.Wallet
	for i := 0; i < wn; i++ {
		WA = append(WA, NewWallet(64))
	}

/*
	B := BC.NewBlock()

	var T coin.Transaction

	if true {
		//create transaction by hand
		var ti coin.TransactionInput
		ti.SigIdx = uint16(0)
		ti.UxOut = genesisWallet.Outputs[0].Hash()
		T.TI = append(T.TI, ti)

		var to coin.TransactionOutput
		to.DestinationAddress = genesisWallet.AA[0].address
		to.Value1 = uint64(100*1000000 - wn*1000)
		T.TO = append(T.TO, to)

		for i := 0; i < wn; i++ {
			var to coin.TransactionOutput
			a := WA[i].GetRandom()
			to.DestinationAddress = a.address
			to.Value1 = 1000
			to.Value2 = 1024
			T.TO = append(T.TO, to)
		}

		var sec coin.SecKey
		sec.Set(genesisAddress.seckey)
		T.SetSig(0, sec)

	} else {
		T.PushInput(genesisWallet.Outputs[0].Hash())
		T.PushOutput(genesisWallet.AA[0].address, uint64(100*1000000-wn*1000), 0)

		for i := 0; i < wn; i++ {
			a := WA[i].GetRandom()
			T.PushOutput(a.address, uint64(1000), 1024*1024)
		}

		var sec coin.SecKey
		sec.Set(genesisAddress.seckey)
		T.SetSig(0, sec)
	}
	T.UpdateHeader() //sets hash

	fmt.Printf("genesis transaction: \n")

	err := BC.AppendTransaction(B, &T)

	if err != nil {
		log.Panic(err)
	}

	WalletBalances(BC, WA)
	err = BC.ExecuteBlock(B)
	if err != nil {
		log.Panic(err)
	}

	WalletBalances(BC, WA)
	*/
}

func main() {
	tests()
}
