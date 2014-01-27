package main

import (
	//"encoding/hex"
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
	genesisAddress := genesisWallet.Addresses[0]
	var BC *coin.BlockChain = coin.NewBlockChain(genesisAddress.Address)

	genesisWallet.RefeshUnspentOutputs(BC)
	//create 16 wallets with 64 addresses
	wn := 16 //number of wallets to create
	var WA []keyring.Wallet
	for i := 0; i < wn; i++ {
		WA = append(WA, keyring.NewWallet(64))
	}

	B := BC.NewBlock()

	var T coin.Transaction

	if true {
		//create transaction by hand
		var ti coin.TransactionInput
		ti.SigIdx = uint16(0)
		ti.UxOut = genesisWallet.Outputs[0].Hash()
		T.TxIn = append(T.TxIn, ti)

		var to coin.TransactionOutput
		to.DestinationAddress = genesisWallet.Addresses[0].Address
		to.Coins = uint64(100*1000000 - wn*1000)
		T.TxOut = append(T.TxOut, to)

		for i := 0; i < wn; i++ {
			var to coin.TransactionOutput
			a := WA[i].GetRandomAddress()
			to.DestinationAddress = a.Address
			to.Coins = 1000
			to.Hours = 1024
			T.TxOut = append(T.TxOut, to)
		}

		var sec coin.SecKey
		sec.Set(genesisAddress.Seckey)
		T.SetSig(0, sec)

	} else {
		T.PushInput(genesisWallet.Outputs[0].Hash())
		T.PushOutput(genesisWallet.Addresses[0].Address, uint64(100*1000000-wn*1000), 0)

		for i := 0; i < wn; i++ {
			a := WA[i].GetRandomAddress()
			T.PushOutput(a.Address, uint64(1000), 1024*1024)
		}

		var sec coin.SecKey
		sec.Set(genesisAddress.Seckey)
		T.SetSig(0, sec)
	}
	T.UpdateHeader() //sets hash

	fmt.Printf("genesis transaction: \n")

	err := BC.AppendTransaction(B, &T)

	if err != nil {
		log.Panic(err)
	}

	keyring.PrintWalletBalances(BC, WA)
	err = BC.ExecuteBlock(B)
	if err != nil {
		log.Panic(err)
	}

	keyring.PrintWalletBalances(BC, WA)
}

func main() {
	tests()
}
