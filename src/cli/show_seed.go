package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func showSeedCmd() *cobra.Command {
	showSeedCmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "showSeed [wallet]",
		Short: "Show wallet seed and seed passphrase",
		Long: `Print the seed and seed passphrase from a wallet.

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			w := args[0]

			password, err := c.Flags().GetString("password")
			if err != nil {
				return err
			}

			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			pr := NewPasswordReader([]byte(password))
			seed, seedPassphrase, err := getSeed(w, pr)
			switch err.(type) {
			case nil:
			case WalletLoadError:
				printHelp(c)
				return err
			default:
				return err
			}

			if jsonOutput {
				v := struct {
					Seed           string `json:"seed"`
					SeedPassphrase string `json:"seed_passphrase,omitempty"`
				}{
					Seed:           seed,
					SeedPassphrase: seedPassphrase,
				}

				return printJSON(v)
			}

			fmt.Println(seed)
			if seedPassphrase != "" {
				fmt.Println(seedPassphrase)
			}
			return nil
		},
	}

	showSeedCmd.Flags().StringP("password", "p", "", "Wallet password")
	showSeedCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")

	return showSeedCmd
}

func getSeed(walletID string, pr PasswordReader) (string, string, error) {
	wlt, err := apiClient.Wallet(walletID)
	if err != nil {
		return "", "", err
	}

	var pwd []byte
	if wlt.Meta.Encrypted {
		pwd, err = pr.Password()
		if err != nil {
			return "", "", err
		}
	}

	sr, err := apiClient.WalletSeed(walletID, string(pwd))
	if err != nil {
		return "", "", err
	}

	return sr.Seed, sr.SeedPassphrase, nil
}
