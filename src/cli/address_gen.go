package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/cipher"
	bip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/wallet"
)

func addressGenCmd() gcli.Command {
	name := "addressGen"
	return gcli.Command{
		Name:        name,
		Usage:       "Generate skycoin or bitcoin addresses",
		Description: "",
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name:  "num,n",
				Value: 1,
				Usage: "Number of addresses to generate",
			},
			gcli.StringFlag{
				Name:  "coin,c",
				Value: "sky",
				Usage: "Coin type. Must be sky or btc",
			},
			// gcli.BoolFlag{
			// 	Name:  "hide-secret,s",
			// 	Usage: "Hide the secret key from the output",
			// },
			gcli.BoolFlag{
				Name:  "hex,x",
				Usage: "Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if not seed is not provided",
			},
			gcli.StringFlag{
				Name:  "seed",
				Usage: "Seed for deterministic key generation. Will use bip39 as the seed if not provided.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Subcommands: []gcli.Command{
			addressGenFile,
			addressGenPrint,
		},
		// 	Action: func(c *gcli.Context) error {
		// 		var coinType wallet.CoinType
		// 		if c.Bool("bitcoin") {
		// 			coinType = wallet.CoinTypeBitcoin
		// 		} else {
		// 			coinType = wallet.CoinTypeSkycoin
		// 		}

		// 		seed := c.String("seed")
		// 		if seed == "" {
		// 			hex := c.Bool("hex")
		// 			if hex {
		// 				// generate a new seed, as hex string
		// 				seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
		// 			} else {
		// 				var err error
		// 				seed, err = bip39.NewDefaultMnemonic()
		// 				if err != nil {
		// 					return err
		// 				}
		// 			}
		// 		}

		// 		w, err := wallet.CreateAddresses(coinType, seed, c.Int("count"), c.Bool("hide-secret"))
		// 		if err != nil {
		// 			return err
		// 		}

		// 		if !c.Bool("only-addr") {
		// 			return printJSON(w)
		// 		}

		// 		for _, e := range w.Entries {
		// 			fmt.Println(e.Address)
		// 		}
		// 		return nil
		// 	},
	}
}

var addressGenFile = gcli.Command{
	Name:         "file",
	Usage:        "Writes addressGen output to file(s)",
	Flags:        []gcli.Flag{},
	OnUsageError: onCommandUsageError("addressGen file"),
	Action: func(c *gcli.Context) error {
	},
}

var addressGenPrint = gcli.Command{
	Name:  "print",
	Usage: "Prints addressGen output to stdout",
	Flags: []gcli.Flag{
		gcli.BoolFlag{
			Name:  "hide-secrets,h",
			Value: false,
			Usage: "Hide the secret key and seed from the output when printing a JSON wallet file",
		},
		gcli.StringFlag{
			Name:  "mode,m",
			Value: "addresses",
			Usage: "Print mode. Options are json (prints a full JSON wallet), addresses (prints addresses in plain text), secrets (prints secret keys in plain text)",
		},
	},
	OnUsageError: onCommandUsageError("addressGen print"),
	Action: func(c *gcli.Context) error {
		coinType, err := wallet.ResolveCoinType(c.GlobalString("coin"))
		if err != nil {
			return err
		}

		seed, err := resolveSeed(c)
		if err != nil {
			return err
		}

		hideSecrets := c.BoolFlag("hide-secrets")
		mode := c.StringFlag("mode")

		switch strings.ToLower(mode) {
		case "json", "wallet":
			if hideSecrets {
				w.Erase()
			}

			output, err := json.MarshalIndent(w, "", "    ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))
		case "addrs", "addresses":
			for _, e := range w.Entries {
				fmt.Println(e.Address)
			}
		case "secrets":
			if hideSecrets {
				return errors.New("secrets mode selected but hideSecrets enabled")
			}
			for _, e := range w.Entries {
				fmt.Println(e.Secret)
			}
		default:
			return errors.New("invalid mode")
		}
	},
}

func resolveSeed(c gcli.Context) (string, error) {
	seed := c.GlobalString("seed")
	useHex := c.GlobalBool("hex")
	strict := c.GlobalBool("strict-seed")

	if seed != "" {
		if strict && !bip39.IsMnemonicValid(seed) {
			return "", errors.New("seed is not a valid bip39 seed")
		}

		return seed, nil
	}

	if useHex {
		seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	} else {
		var err error
		seed, err = bip39.NewDefaultMnemonic()
		if err != nil {
			return "", err
		}
	}

	return seed, nil
}
