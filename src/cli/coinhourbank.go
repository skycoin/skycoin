package cli

import (
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	chb "github.com/skycoin/skycoin/src/coinhourbank"
	"github.com/spf13/cobra"
)

func coinhourBalanceCmd() *cobra.Command {
	coinhourBalanceCmd := &cobra.Command{
		Use: "coinhourBalance",
		Short: "Get balance of coinhour bank account",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			// error is already checked when bank client is initialized
			address, _ := c.Flags().GetString("address")
			if _, err := cipher.DecodeBase58Address(address); err != nil {
				return err
			}

			balance, err := bankClient.Balance(chb.Account(address))
			if err != nil {
				return err
			}

			fmt.Printf("%s balance: %v\n", address, balance)

			return nil
		},
	}

	coinhourBalanceCmd.Flags().StringP("wallet-file", "f", "", "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.")
	coinhourBalanceCmd.Flags().StringP("password", "p", "", "wallet password")
	coinhourBalanceCmd.Flags().StringP("address", "a", "", "wallet address to take coinhours from")
	coinhourBalanceCmd.Flags().StringP("nodeURL", "n", "http://localhost:6420", "skycoin node url")
	coinhourBalanceCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return coinhourBalanceCmd
}

func depositCoinhoursCmd() *cobra.Command {
	depositCoinhoursCmd := &cobra.Command{
		Use:   "depositCoinhours [hours amount]",
		Short: "Sends coinhours to a coinhour bank account.",
		Long: `Deposits coinhours into a coinhour bank account which a skycoin address you want to deposit hours into.
		Once hours are into coinhour bank they can be transferred to other addresses without paying transaction fee.`,
		SilenceUsage: true,
		DisableFlagsInUseLine: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			address, _ := c.Flags().GetString("address")
			if _, err := cipher.DecodeBase58Address(address); err != nil {
				return err
			}

			coinhours, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			return bankClient.DepositHours(chb.CoinHours(coinhours), chb.Account(address))
		},
	}

	depositCoinhoursCmd.Flags().StringP("wallet-file", "f", "", "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.")
	depositCoinhoursCmd.Flags().StringP("password", "p", "", "wallet password")
	depositCoinhoursCmd.Flags().StringP("address", "a", "", "wallet address to take coinhours from")
	depositCoinhoursCmd.Flags().StringP("nodeURL", "n", "http://localhost:6420", "skycoin node url")
	depositCoinhoursCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return depositCoinhoursCmd
}

func transferCoinhoursCmd() *cobra.Command {
	transferCoinhoursCmd := &cobra.Command{
		Use:          "transferCoinhours [destination address] [hours amount]",
		Short:        "Transfer coinhours from one coinhour bank account to another",
		Long:         `Transferring coinhours from one account to another does not require any transaction fee.`,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			address, _ := c.Flags().GetString("address")
			if _, err := cipher.DecodeBase58Address(address); err != nil {
				return err
			}

			coinhours, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			return bankClient.TransferHours(chb.Account(address), chb.Account(args[0]), chb.CoinHours(coinhours))
		},
	}

	transferCoinhoursCmd.Flags().StringP("wallet-file", "f", "", "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.")
	transferCoinhoursCmd.Flags().StringP("password", "p", "", "wallet password")
	transferCoinhoursCmd.Flags().StringP("address", "a", "", "coinhour bank to take coinhours from")
	transferCoinhoursCmd.Flags().StringP("nodeURL", "n", "http://localhost:6420", "skycoin node url")
	transferCoinhoursCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return transferCoinhoursCmd
}

func withdrawCoinhoursCmd() *cobra.Command {
	withdrawCoinhoursCmd := &cobra.Command{
		Use:          "withdrawCoinhours [hours amount]",
		Short:        "Withdraws coinhours from coinhour bank.",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			address, _ := c.Flags().GetString("address")
			if _, err := cipher.DecodeBase58Address(address); err != nil {
				return err
			}

			coinhours, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			return bankClient.WithdrawHours(chb.CoinHours(coinhours), chb.Account(address))
		},
	}

	withdrawCoinhoursCmd.Flags().StringP("wallet-file", "f", "", "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.")
	withdrawCoinhoursCmd.Flags().StringP("password", "p", "", "wallet password")
	withdrawCoinhoursCmd.Flags().StringP("address", "a", "", "wallet address to take coinhours from")
	withdrawCoinhoursCmd.Flags().StringP("nodeURL", "n", "http://localhost:6420", "skycoin node url")
	withdrawCoinhoursCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return withdrawCoinhoursCmd
}

func getCoinhourBankClient(c *cobra.Command) (*chb.HourBankClient, error) {
	walletFile, err := c.Flags().GetString("wallet-file")
	if err != nil {
		return nil, err
	}

	wlt, err := resolveWalletPath(cliConfig, walletFile)
	if err != nil {
		return nil, err
	}

	wltPassword, err := c.Flags().GetString("password")
	if err != nil {
		return nil, err
	}

	nodeURL, err := c.Flags().GetString("nodeURL")
	if err != nil {
		return nil, err
	}

	bankURL, err := c.Flags().GetString("bankURL")
	if err != nil {
		return nil, err
	}

	address, err := c.Flags().GetString("address")
	if err != nil {
		return nil, err
	}

	bankClient, err := chb.NewHourBankClient(nodeURL, bankURL, wlt, []byte(wltPassword), chb.SourceAddress(address))
	if err != nil {
		return nil, err
	}

	return bankClient, nil
}
