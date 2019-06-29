package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/wallet"
)

func showSeedCmd() *cobra.Command {
	showSeedCmd := &cobra.Command{
		Use:   "showSeed",
		Short: "Show wallet seed",
		Long: fmt.Sprintf(`The default wallet (%s) will be used if no wallet was specified.

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`, cliConfig.FullWalletPath()),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, _ []string) error {
			w, err := resolveWalletPath(cliConfig, "")
			if err != nil {
				return err
			}

			password, err := c.Flags().GetString("password")
			if err != nil {
				return err
			}

			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			pr := NewPasswordReader([]byte(password))
			seed, err := getSeed(w, pr)
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
					Seed string `json:"seed"`
				}{
					Seed: seed,
				}

				return printJSON(v)
			}

			fmt.Println(seed)
			return nil
		},
	}

	showSeedCmd.Flags().StringP("password", "p", "", "Wallet password")
	showSeedCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")

	return showSeedCmd
}

func getSeed(walletFile string, pr PasswordReader) (string, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return "", WalletLoadError{err}
	}

	switch pr.(type) {
	case nil:
		if wlt.IsEncrypted() {
			return "", wallet.ErrWalletEncrypted
		}
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return "", err
		}

		if !wlt.IsEncrypted() && len(p) != 0 {
			return "", wallet.ErrWalletNotEncrypted
		}
	}

	if !wlt.IsEncrypted() {
		return wlt.Seed(), nil
	}

	password, err := pr.Password()
	if err != nil {
		return "", err
	}

	var seed string
	if err := wallet.GuardView(wlt, password, func(w wallet.Wallet) error {
		seed = w.Seed()
		return nil
	}); err != nil {
		return "", err
	}

	return seed, nil
}
