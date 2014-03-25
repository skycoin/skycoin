package main

import (
	"github.com/skycoin/sync/src/hashchain"

	"fmt"
	"log"
)

func main() {

	_, seckey := hashchain.GenerateDeterministicKeyPair([]byte("seed"))
	bc := hashchain.NewBlockChain(seckey)

	for i := 0; i < 256; i++ {

		//write this data to the block
		s := fmt.Sprintf("test data: %v", i)

		block := bc.NewBlock(seckey, uint64(time.Now().Unix()), []byte(s))

		err := bc.ApplyBlock(block)
		if err != nil {
			log.Panic(err)
		}
	}

	_ = bc
}
