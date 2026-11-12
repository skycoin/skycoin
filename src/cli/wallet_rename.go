// Package cli src/cli/wallet_rename.go c3-coin-cli
package cli

import (
	"github.com/spf13/cobra"
)

func walletRenameCmd() *cobra.Command {
	return &cobra.Command{
		Args:  cobra.ExactArgs(2),
		Short: "Change the label of a wallet",
		Use:   "walletRename [wallet] [label]",
		Long: `Change the label of a wallet the node holds. The node saves the new
    label in the wallet file. The wallet is named by its file name, as
    listWallets prints it.`,
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return walletRename(args[0], args[1])
		},
	}
}

func walletRename(id, label string) error {
	if err := apiClient.UpdateWallet(id, label); err != nil {
		return err
	}

	wlt, err := apiClient.Wallet(id)
	if err != nil {
		return err
	}

	return printJSON(wlt.Meta)
}
