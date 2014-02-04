package keyring

/*
   Start a local Blockchain service

*/

import (
    //"encoding/hex"
    //"errors"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "time"
    //"github.com/skycoin/skycoin/src/keyring"

    "log"
    "math/rand"
    //"encoding/hex"
)

/*
Creates a new block every 15 seconds

*/
type BlockchainService struct {
    PendingBlock        coin.Block
    BC                  *coin.Blockchain
    PendingTransactions []coin.Transaction
}

func (self *BlockchainService) Run() {
    //TODO, set genesis address

    seckey_hex := "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d"
    seckey := coin.SecKeyFromHex(seckey_hex)
    pubkey := coin.PubKeyFromSecKey(seckey)
    address := coin.AddressFromPubKey(pubkey) //genesis address

    self.BC = coin.NewBlockchain(address)

    go func() {

        block := self.BC.NewBlock()
        for true {
            time.Sleep(250 * time.Millisecond)
            if self.BC.Head.Header.Time > uint64(time.Now().Unix()) {
                continue
            }
        }

        //pull some transactions out
        for _, t := range self.PendingTransactions {
            if rand.Int()%2 == 0 {
                continue
            }
            err := self.BC.AppendTransaction(&block, t)
            if err == nil {
                continue
            }
        }
        //execute the transactions
        fmt.Printf("New Block!")
        err := self.BC.ExecuteBlock(block)
        if err != nil {
            log.Panic()
        }
        block = self.BC.NewBlock()

    }()
}

func (self *BlockchainService) InsertTransaction(transaction coin.Transaction) {
    self.PendingTransactions = append(self.PendingTransactions, transaction)
}
