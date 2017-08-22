package cli

import (
	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/cipher"
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
			gcli.StringFlag{
				Name:  "seed",
				Usage: "Seed for deterministic key generation. Will use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if not provided.",
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
				// generate a new seed, as hex string
				seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
			}

			w, err := wallet.CreateAddresses(coinType, seed, c.Int("count"), c.Bool("hide-secret"))
			if err != nil {
				return err
			}

			return printJson(w)
		},
	}
}
