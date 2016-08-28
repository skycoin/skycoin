package main

import (
	"fmt"
	"log"
	"time"

	"github.com/skycoin/skycoin/src/aether/hashchain"
	"github.com/skycoin/skycoin/src/cipher"
)

//Blockchain sync example
/*
	Todo
	- publish chain
	- get chain public key hash
	- look up peers through DHT via hash
	- peers download the chain and apply blocks

*/

//creates chain and hosts
func runChain() {

	_, seckey := cipher.GenerateDeterministicKeyPair([]byte("seed"))
	bc := hashchain.NewBlockChain(seckey)

	for i := 0; i < 256; i++ {

		//write this data to the block
		s := fmt.Sprintf("test data: %v", i)

		block := bc.NewBlock(seckey, uint64(time.Now().Unix()), []byte(s))

		err := bc.ApplyBlock(block)
		if err != nil {
			log.Panic(err)
		}

		//do something with the data, if block is valid
	}

	_ = bc
}

func main() {
	runChain()
}
