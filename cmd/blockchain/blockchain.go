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

func (self *PendingTransactions) Push(t coin.Transaction) {
    self.A = append(self.A, t)
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
    var bc *coin.BlockChain = coin.NewBlockChain(genesisAddress.Address)

    genesisWallet.RefeshUnspentOutputs(bc)
    //create 16 wallets with 64 addresses
    wn := 16 //number of wallets to create
    var wa []keyring.Wallet
    for i := 0; i < wn; i++ {
        wa = append(wa, keyring.NewWallet(64))
    }

    b := bc.NewBlock()

    var t coin.Transaction

    if true {
        //create transaction by hand
        var ti coin.TransactionInput
        ti.SigIdx = uint16(0)
        ti.UxOut = genesisWallet.Outputs[0].Hash()
        t.TxIn = append(t.TxIn, ti)

        var to coin.TransactionOutput
        to.DestinationAddress = genesisWallet.Addresses[0].Address
        to.Coins = uint64(100*1e6 - wn*1000)
        t.TxOut = append(t.TxOut, to)

        for i := 0; i < wn; i++ {
            var to coin.TransactionOutput
            a := wa[i].GetRandomAddress()
            to.DestinationAddress = a.Address
            to.Coins = 1000
            to.Hours = 1024
            t.TxOut = append(t.TxOut, to)
        }

        var sec coin.SecKey
        sec.Set(genesisAddress.SecKey[:])
        t.SetSig(0, sec)

    } else {
        t.PushInput(genesisWallet.Outputs[0].Hash())
        t.PushOutput(genesisWallet.Addresses[0].Address,
            uint64(100*1e6-wn*1000), 0)

        for i := 0; i < wn; i++ {
            a := wa[i].GetRandomAddress()
            t.PushOutput(a.Address, uint64(1000), 1024*1024)
        }

        var sec coin.SecKey
        sec.Set(genesisAddress.SecKey[:])
        t.SetSig(0, sec)
    }
    t.UpdateHeader() //sets hash

    fmt.Printf("genesis transaction: \n")

    err := bc.AppendTransaction(&b, t)
    if err != nil {
        log.Panic(err)
    }

    keyring.PrintWalletBalances(bc, wa)
    err = bc.ExecuteBlock(b)
    if err != nil {
        log.Panic(err)
    }

    fmt.Printf("after execution: \n")
    keyring.PrintWalletBalances(bc, wa)
}

func main() {
    tests()
}
