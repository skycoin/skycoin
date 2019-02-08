package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/coinhourbank"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

func depositCoinhoursCmd() *cobra.Command {
	depositCoinhoursCmd := &cobra.Command{
		Use:   "depositCoinhours [hours amount]",
		Short: "Sends coinhours to a coinhour bank account.",
		Long: `Deposits coinhours into a coinhour bank account which a skycoin address you want to deposit hours into.
		Once hours are into coinhour bank they can be transferred to other addresses without paying transaction fee.`,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			wlt, err := getWallet(c)
			if err != nil {
				return err
			}
			defer wlt.Erase()

			address, _ := c.Flags().GetString("address") // nolint: errcheck

			coinhours, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			outputSet, err := apiClient.OutputsForAddresses([]string{string(address)})
			if err != nil {
				return err
			}

			uxa, err := outputSet.SpendableOutputs().ToUxArray()
			if err != nil {
				return err
			}

			return bankClient.DepositHours(coinhours, address, uxa, wlt)
		},
	}

	depositCoinhoursCmd.Flags().StringP("wallet-file", "f", "", "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.")
	depositCoinhoursCmd.Flags().StringP("password", "p", "", "wallet password")
	depositCoinhoursCmd.Flags().StringP("address", "a", "", "wallet address to take coinhours from")
	depositCoinhoursCmd.Flags().StringP("nodeURL", "n", "http://localhost:6420", "skycoin node url")
	depositCoinhoursCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return depositCoinhoursCmd
}

func getCoinhourBankClient(c *cobra.Command) (*coinhourbank.HourBankClient, error) {
	bankURL, err := c.Flags().GetString("bankURL")
	if err != nil {
		return nil, err
	}

	bankClient := coinhourbank.NewHourBankClient(bankURL)
	return bankClient, nil
}

func getWallet(c *cobra.Command) (*wallet.Wallet, error) {
	walletFile, err := c.Flags().GetString("wallet-file")
	if err != nil {
		return nil, nil
	}

	wltPath, err := resolveWalletPath(cliConfig, walletFile)
	if err != nil {
		return nil, err
	}

	wlt, err := wallet.Load(wltPath)
	if err != nil {
		return nil, err
	}

	address, err := c.Flags().GetString("address")
	if err != nil {
		return nil, err
	}

	sourceAddr, err := cipher.DecodeBase58Address(address)
	if err != nil {
		return nil, err
	}

	if _, ok := wlt.GetEntry(sourceAddr); !ok {
		return nil, fmt.Errorf("sender address not in wallet")
	}

	password, err := c.Flags().GetString("password")
	if err != nil {
		return nil, err
	}
	pr := NewPasswordReader([]byte(password))

	switch pr.(type) {
	case nil:
		if wlt.IsEncrypted() {
			return nil, wallet.ErrWalletEncrypted
		}
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return nil, err
		}

		if !wlt.IsEncrypted() && len(p) != 0 {
			return nil, wallet.ErrWalletNotEncrypted
		}
	}

	if wlt.IsEncrypted() {
		p, err := pr.Password()
		if err != nil {
			return nil, err
		}

		return wlt.Unlock(p)
	}

	return wlt, nil
}
