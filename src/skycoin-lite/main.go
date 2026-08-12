// Package main provides a lightweight Skycoin client implementation
package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/skycoin-lite/liteclient"
)

// For testing purposes. This file is not part of the library.
//
// This used to import github.com/skycoin/skycoin-lite, an external module last
// published in 2019, rather than the package next to it — so it demonstrated a
// seven-year-old copy of the client instead of this one.
func main() {
	seed := hex.EncodeToString([]byte("nest*"))

	addrs, err := liteclient.GenerateAddresses(seed, 3)
	if err != nil {
		fmt.Fprintln(os.Stderr, "generating addresses:", err)
		os.Exit(1)
	}
	fmt.Println(addrs)

	next := seed
	for i := 0; i < 3; i++ {
		address, err := liteclient.GenerateAddress(next)
		if err != nil {
			fmt.Fprintln(os.Stderr, "generating address:", err)
			os.Exit(1)
		}

		fmt.Println("----")
		fmt.Println(address)
		next = address.NextSeed
	}
}
