package main

import (
	"encoding/hex"
	"fmt"

	"github.com/skycoin/skycoin-lite/liteclient"
)

// For testing purposes. This file is not part of the library.
func main() {
	seed := hex.EncodeToString([]byte("nest*"))

	addrs := liteclient.GenerateAddresses(seed, 3)
	fmt.Println(addrs)

	address1 := liteclient.GenerateAddress(seed)
	fmt.Println("----")
	fmt.Println(address1)

	address2 := liteclient.GenerateAddress(address1.NextSeed)
	fmt.Println("----")
	fmt.Println(address2)

	address3 := liteclient.GenerateAddress(address2.NextSeed)
	fmt.Println("----")
	fmt.Println(address3)
}
