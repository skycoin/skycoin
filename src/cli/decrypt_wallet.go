package cli

import (
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/spf13/cobra"
)

func decryptWalletCmd() *cobra.Command {
	decryptWalletCmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "decryptWallet [wallet]",
		Short: "Decrypt a wallet",
		Long: `Decrypt an encrypted wallet. The decrypted wallet will be written
    on the filesystem in place of the encrypted wallet.

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			w := args[0]
			pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))

			return decryptWallet(w, pr)
		},
	}

	decryptWalletCmd.Flags().StringP("password", "p", "", "wallet password")

	return decryptWalletCmd
}

func decryptWallet(id string, pr PasswordReader) error {
	wlt, err := apiClient.Wallet(id)
	if err != nil {
		return err
	}

	if !wlt.Meta.Encrypted {
		return wallet.ErrWalletNotEncrypted
	}

	if pr == nil {
		return wallet.ErrMissingPassword
	}

	pwd, err := pr.Password()
	if err != nil {
		return err
	}

	wlt, err = apiClient.DecryptWallet(id, string(pwd))
	if err != nil {
		return err
	}

	return printJSON(wlt)
}
