package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cipher"
	bip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
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
	hexSeed := flag.Bool("x", false, "Use hex string as seed")
	verbose := flag.Bool("v", false, "Show verbose info of generated addresses")
	seed := flag.String("seed", "", "Seed for deterministic key generation. Will use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if not provided.")
	flag.Parse()

	var coinType wallet.CoinType
	if *isBitcoin {
		coinType = wallet.CoinTypeBitcoin
	} else {
		coinType = wallet.CoinTypeSkycoin
	}

	if *seed == "" {
		if *hexSeed {
			// generate a new seed, as hex string
			*seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
		} else {
			entropy, err := bip39.NewEntropy(128)
			if err != nil {
				fmt.Printf("new entropy failed when new wallet seed: %v\n", err)
				return
			}

			mnemonic, err := bip39.NewMnemonic(entropy)
			if err != nil {
				fmt.Printf("new mnemonic failed when new wallet seed: %v\n", err)
				return
			}

			*seed = mnemonic
		}
	}

	w, err := wallet.CreateAddresses(coinType, *seed, *genCount, *hideSecKey)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *verbose {
		output, err := json.MarshalIndent(w, "", "    ")
		if err != nil {
			fmt.Println("Error formating wallet to JSON. Error:", err)
			os.Exit(1)
		}

		fmt.Println(string(output))
		return
	}

	// only show the address, no more info
	addrs := make([]string, 0, len(w.Entries))

	for _, e := range w.Entries {
		addrs = append(addrs, e.Address)
	}

	v, err := json.MarshalIndent(addrs, "", "    ")
	if err != nil {
		fmt.Println("Error formating address list:", err)
		return
	}

	fmt.Println(string(v))
}
