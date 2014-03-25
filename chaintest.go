package main

import (
	"github.com/skycoin/sync/src/hashchain"

	//"fmt"
	//"log"
)

func main() {

	_, seckey := hashchain.GenerateDeterministicKeyPair([]byte("seed"))
	bc := hashchain.NewBlockChain(seckey)

	_ = bc
}
