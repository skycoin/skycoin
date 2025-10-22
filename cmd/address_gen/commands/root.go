// package commands cmd/address_gen/commands/root.go
/*
 generate public keys and addresses
*/
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/deterministic"
)

var (
	genCount    int
	isBitcoin   bool
	hexSeed     bool
	hideSecrets bool
	seed        string
	secKeysList bool
	addrsList   bool
)

func init() {
	logging.Disable()
	RootCmd.Flags().IntVarP(&genCount, "number", "n", 1, "Number of addresses to generate")
	RootCmd.Flags().BoolVarP(&isBitcoin, "bitcoin", "b", false, "Print address as a bitcoin address")
	RootCmd.Flags().BoolVarP(&hexSeed, "hex-seed", "x", false, "Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if seed is not provided")
	RootCmd.Flags().BoolVar(&hideSecrets, "hide-secrets", false, "Hide seed and secret key")
	RootCmd.Flags().StringVar(&seed, "seed", "", "Seed for deterministic key generation. Will use bip39 as the seed if not provided")
	RootCmd.Flags().BoolVar(&secKeysList, "sec-keys-list", false, "Only print a list of secret keys")
	RootCmd.Flags().BoolVar(&addrsList, "addrs-list", false, "Only print a list of addresses")
}

var RootCmd = &cobra.Command{
	Use:   "address-gen",
	Short: "generate public keys and addresses",
	Long: `
	┌─┐┌┬┐┌┬┐┬─┐┌─┐┌─┐┌─┐   ┌─┐┌─┐┌┐┌
	├─┤ ││ ││├┬┘├┤ └─┐└─┐───│ ┬├┤ │││
	┴ ┴─┴┘─┴┘┴└─└─┘└─┘└─┘   └─┘└─┘┘└┘
	generate public keys and addresses`,
	Run: func(_ *cobra.Command, _ []string) {
		var coinType wallet.CoinType
		if isBitcoin {
			coinType = wallet.CoinTypeBitcoin
		} else {
			coinType = wallet.CoinTypeSkycoin
		}

		if seed == "" {
			if hexSeed {
				// generate a new seed, as hex string
				seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
			} else {
				mnemonic, err := bip39.NewDefaultMnemonic()
				if err != nil {
					fmt.Printf("bip39.NewDefaultMnemonic failed: %v\n", err)
					os.Exit(1)
				}
				seed = mnemonic
			}
		}

		w, err := wallet.NewWallet("a.wlt", "", seed, wallet.Options{
			Type:      deterministic.WalletType,
			Coin:      coinType,
			GenerateN: uint64(genCount),
		})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if hideSecrets {
			w.Erase()
		}

		if hideSecrets && secKeysList {
			fmt.Println("-hide-secrets and -sec-keys-list can't be combined")
			os.Exit(1)
		}

		if addrsList && secKeysList {
			fmt.Println("-addrs-list and -sec-keys-list can't be combined")
			os.Exit(1)
		}

		if addrsList {
			addrs, err := w.GetAddresses()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, a := range addrs {
				fmt.Println(a)
			}
		} else if secKeysList {
			es, err := w.GetEntries()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, e := range es {
				fmt.Println(e.Secret.Hex())
			}
		} else {
			if hideSecrets {
				w.Erase()
			}
			output, err := w.Serialize()
			if err != nil {
				fmt.Println("Error formatting wallet to JSON. Error:", err)
				os.Exit(1)
			}
			fmt.Println(string(output))
		}
	},
}
