package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	bip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/wallet"
)

type genInfo struct {
	Date      time.Time
	Seed      string
	Coin      string
	KeysCount int
}

func run() error {
	genCount := flag.Int("n", 1, "Number of addresses to generate")
	seed := flag.String("seed", "", "Seed for deterministic key generation. Will use bip39 as the seed if not provided")
	strict := flag.Bool("strict", true, "Checks if input is space separated list of words.")
	coin := flag.String("coin", "skycoin", "address output type: sky/btc")
	secfile := flag.String("secfile", "", "command for file to write the secret keys")
	addrOut := flag.String("addrfile", "addresses", "command for changing addresses output file")
	outputInfo := flag.String("infofile", "", "create file with date of generation, seed, coin, number of keys generated")
	flag.Parse()

	var coinType wallet.CoinType
	switch *coin {
	case "btc", "bitcoin":
		coinType = wallet.CoinTypeBitcoin
	case "sky", "skycoin":
		coinType = wallet.CoinTypeSkycoin
	default:
		return errors.New("unknown coin type")
	}

	if *seed != "" && *strict {
		if !bip39.IsMnemonicValid(*seed) {
			return errors.New("your seed isn't valid")
		}
	}

	if *seed == "" {
		entropy, err := bip39.NewEntropy(128)
		if err != nil {
			return err
		}

		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			return err
		}

		*seed = mnemonic
		fmt.Println(*seed)
	}

	w, err := wallet.CreateAddresses(coinType, *seed, *genCount, false)
	if err != nil {
		return err
	}

	if *secfile != "" {
		f, err := os.OpenFile(*secfile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return errors.New("secret file already exist")
		}
		for _, e := range w.Entries {
			fmt.Fprintln(f, e.Secret)
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	if *addrOut != "" {
		f, err := os.Create(*addrOut)
		if err != nil {
			return err
		}
		for _, e := range w.Entries {
			fmt.Fprintln(f, e.Address)
		}

		err = f.Close()
		if err != nil {
			return err
		}
	} else {
		return errors.New("file for addresses output doesn't specified")
	}

	if *outputInfo != "" {
		f, err := os.Create(*outputInfo)
		if err != nil {
			return err
		}

		var info genInfo
		info.Coin = *coin
		info.Date = time.Now()
		info.KeysCount = *genCount
		info.Seed = *seed

		infoJSON, err := json.MarshalIndent(info, "", "    ")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(*outputInfo, infoJSON, 0644)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
