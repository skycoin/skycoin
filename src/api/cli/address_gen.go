package cli

import (
	"fmt"

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
				Name:  "count,c",
				Value: 1,
				Usage: "Number of addresses to generate",
			},
			gcli.BoolFlag{
				Name:  "hide-secret,s",
				Usage: "Hide the secret key from the output",
			},
			gcli.BoolFlag{
				Name:  "bitcoin,b",
				Usage: "Output the addresses as bitcoin addresses instead of skycoin addresses",
			},
			gcli.BoolFlag{
				Name:  "hex,x",
				Usage: "Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if not seed is not provided",
			},
			gcli.BoolFlag{
				Name:  "only-addr,oa",
				Usage: "Only show generated address list, hide seed, secret key and public key",
			},
			gcli.StringFlag{
				Name:  "seed",
				Usage: "Seed for deterministic key generation. Will use bip39 as the seed if not provided.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			var coinType wallet.CoinType
			if c.Bool("bitcoin") {
				coinType = wallet.CoinTypeBitcoin
			} else {
				coinType = wallet.CoinTypeSkycoin
			}

			seed := c.String("seed")
			if seed == "" {
				hex := c.Bool("hex")
				if hex {
					// generate a new seed, as hex string
					seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
				} else {
					var err error
					var entropy []byte
					entropy, err = bip39.NewEntropy(128)
					if err != nil {
						return fmt.Errorf("new entropy failed when new wallet seed: %v", err)
					}

					seed, err = bip39.NewMnemonic(entropy)
				}
			}

			w, err := wallet.CreateAddresses(coinType, seed, c.Int("count"), c.Bool("hide-secret"))
			if err != nil {
				return err
			}

			if !c.Bool("only-addr") {
				return printJson(w)
			}

			for _, e := range w.Entries {
				fmt.Println(e.Address)
			}
			return nil
		},
	}
}
