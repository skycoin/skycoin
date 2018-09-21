package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	bip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/wallet"
)

type genInfo struct {
	Date      time.Time
	Seed      string
	Coin      string
	KeysCount int
}

func main() {
	genCount := flag.Int("n", 1, "Number of addresses to generate")
	hideSecKey := flag.Bool("s", false, "Hide the secret key from the output")
	coin := flag.String("coin", "sky", "address output type. options are sky, btc")
	hexSeed := flag.Bool("x", false, "Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if seed is not provided")
	hideSecrets := flag.Bool("hide-secrets", false, "Hide secret keys and seed from JSON output")
	printSeed := flag.Bool("print-seed", false, "print the seed used")
	seed := flag.String("seed", "", "Seed for deterministic key generation. Will use bip39 as the seed if not provided")
	strictSeed := flag.Bool("strict-seed", true, "Validate the seed as a bip39 seed")
	secKeysList := flag.Bool("sec-keys-list", false, "Only print a list of secret keys")
	addrsList := flag.Bool("addrs-list", false, "Only print a list of addresses")

	secFile := flag.String("sec-keys-file", "", "write secrets to this location")
	addrFile := flag.String("addrs-file", "", "write addresses to this location")
	infoFile := flag.String("info-file", "", "write metadata file with date of generation, coin, number of keys generated to this file")
	infoFileSeed := flag.Bool("info-file-seed", false, "include the seed in the infoFile")

	fiber := flag.Bool("fiber-addresses", false, "generate addresses in format used by fiber.toml")

	flag.Parse()

	var coinType wallet.CoinType
	switch *coin {
	case "skycoin", "sky":
		coinType = wallet.CoinTypeSkycoin
	case "bitcoin", "btc":
		coinType = wallet.CoinTypeBitcoin
	default:
		fmt.Println("invalid coin type")
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

	if !*hexSeed && *strictSeed {
		if !bip39.IsMnemonicValid(*seed) {
			fmt.Println("seed is not a valid bip39 seed")
			os.Exit(1)
		}
	}

	w, err := wallet.CreateAddresses(coinType, *seed, *genCount, *hideSecKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *secFile != "" {
		f, err := os.OpenFile(*secFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			fmt.Println("secrets file already exists")
			os.Exit(1)
		}
		for _, e := range w.Entries {
			fmt.Fprintln(f, e.Secret)
		}

		err = f.Close()
		fmt.Println(err)
		os.Exit(1)
	}

	if *addrFile != "" {
		f, err := os.OpenFile(*addrFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			fmt.Println("addresses file already exists")
			os.Exit(1)
		}

		for _, e := range w.Entries {
			if *fiber {
				fmt.Fprintln(f, fmt.Sprintf("\"%s\",", e.Address))
			} else {
				fmt.Fprintln(f, e.Address)
			}
		}

		err = f.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if *infoFile != "" {
		f, err := os.Create(*infoFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var info genInfo
		info.Coin = *coin
		info.Date = time.Now()
		info.KeysCount = *genCount
		if *infoFileSeed {
			info.Seed = *seed
		}

		infoJSON, err := json.MarshalIndent(info, "", "    ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = ioutil.WriteFile(*infoFile, infoJSON, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = f.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if *addrsList {
		for _, e := range w.Entries {
			fmt.Println(e.Address)
		}
	} else if *secKeysList {
		for _, e := range w.Entries {
			fmt.Println(e.Secret)
		}
	} else {
		if *hideSecrets {
			w.Erase()
		}

		output, err := json.MarshalIndent(w, "", "    ")
		if err != nil {
			fmt.Println("Error formating wallet to JSON. Error:", err)
			os.Exit(1)
		}

		fmt.Println(string(output))
	}

	if *printSeed && !*hideSecrets {
		fmt.Println("Seed:", *seed)
	}
}
