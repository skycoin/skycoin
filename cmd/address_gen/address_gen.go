package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

// Note: Address_gen generates public keys and addresses
// address, pubkey, privatekey
// -n=5 for number of addresses
// -seed to set wallet seed. Prompt will ask
// for seed to prevent seed from being stored in .bashrc

// -json for json output
// -add option to password the secret key
// -let people add the key from the command line

func main() {
	genCount := flag.Int("n", 1, "Number of addresses to generate")
	hideSecKey := flag.Bool("s", false, "Hide the secret key from the output")
	isBitcoin := flag.Bool("b", false, "Print address as a bitcoin address")
	seed := flag.String("seed", "", "Seed for deterministic key generation. Will generate a random 1024-byte CSPRNG-generated seed if not provided.")
	flag.Parse()

	var coinType wallet.CoinType
	if *isBitcoin {
		coinType = wallet.CoinTypeBitcoin
	} else {
		coinType = wallet.CoinTypeSkycoin
	}

	if *seed == "" {
		// generate a new seed, as hex string
		*seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	}

	w, err := wallet.CreateAddresses(coinType, *seed, *genCount, *hideSecKey)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		fmt.Println("Error formating wallet to JSON. Error:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
