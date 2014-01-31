package keyring

/*
	Start a local blockchain service

*/

import (
    //"encoding/hex"
    //"errors"
    "fmt"
    "time"
    "github.com/skycoin/skycoin/src/coin"
    //"github.com/skycoin/skycoin/src/keyring"

    //"log"
    //"math/rand"
    //"encoding/hex"
)

/*
Creates a new block every 15 seconds

*/
type BlockChainService struct {
	PendingBlock coin.Block
	BC *coin.BlockChain
	PendingTransactions []coin.Transaction
}


func (self *BlockChainService) Run() {
	//TODO, set genesis address

	seckey_hex := "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d"
   	seckey := coin.SecKeyFromHex(seckey_hex)
    pubkey := coin.PubKeyFromSecKey(seckey)
    address := coin.AddressFromPubKey(pubkey) //genesis address

	self.BC = coin.NewBlockChain(address)

	go func(){
		for true {
			time.Sleep(250*time.Millisecond)	
			if self.BC.Head.Header.Time > uint64(time.Now().Unix()) {
				continue
			}
		}

		fmt.Printf("New Block!")

	}()
}

func (self *BlockChainService) InsertTransaction(transaction coin.Transaction) {
	self.PendingTransactions = append(self.PendingTransactions, transaction)
}


