package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/util/logging"
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
	logging.Disable()

	genCount := flag.Int("n", 1, "Number of addresses to generate")
	hideSecKey := flag.Bool("s", false, "Hide the secret key from the output")
	isBitcoin := flag.Bool("b", false, "Print address as a bitcoin address")
	hexSeed := flag.Bool("x", false, "Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if seed is not provided")
	hideSecrets := flag.Bool("hide-secrets", false, "Hide seed and secret key")
	seed := flag.String("seed", "", "Seed for deterministic key generation. Will use bip39 as the seed if not provided")
	secKeysList := flag.Bool("sec-keys-list", false, "only print a list of secret keys")
	addrsList := flag.Bool("addrs-list", false, "only print a list of addresses")
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
			mnemonic, err := bip39.NewDefaultMnemonic()
			if err != nil {
				fmt.Printf("bip39.NewDefaultMnemonic failed: %v\n", err)
				os.Exit(1)
			}

			*seed = mnemonic
		}
	}

	// w, err := wallet.NewWallet(coinType, *seed, *genCount, *hideSecKey)
	w, err := wallet.NewWallet("a.wlt", wallet.Options{
		Type:      "deterministic",
		Coin:      coinType,
		Seed:      *seed,
		GenerateN: uint64(*genCount),
		Encrypt:   *hideSecKey,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *hideSecrets && *secKeysList {
		fmt.Println("-hide-secrets and -sec-keys-list can't be combined")
		os.Exit(1)
	}

	if *addrsList && *secKeysList {
		fmt.Println("-addrs-list and -sec-keys-list can't be combined")
		os.Exit(1)
	}

	if *addrsList {
		for _, a := range w.GetAddresses() {
			fmt.Println(a)
		}
		// for _, e := range w.Entries {
		// }
	} else if *secKeysList {
		// for _, e := range w.Entries {
		// 	fmt.Println(e.Secret)
		// }
		for _, e := range w.GetEntries() {
			fmt.Println(e.Secret.Hex())
		}
	} else {
		if *hideSecrets {
			w.Erase()
		}
		output, err := json.MarshalIndent(w.ToReadable(), "", "    ")
		if err != nil {
			fmt.Println("Error formating wallet to JSON. Error:", err)
			os.Exit(1)
		}

		fmt.Println(string(output))
	}
}
