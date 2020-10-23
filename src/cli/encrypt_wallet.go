package cli

import (
	"github.com/spf13/cobra"

	"github.com/SkycoinProject/skycoin/src/wallet"
)

func encryptWalletCmd() *cobra.Command {
	encryptWalletCmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Short: "Encrypt wallet",
		Use:   "encryptWallet [wallet]",
		Long: `Encrypt a decrypted wallet. The encrypted wallet file
    will be written on the filesystem in place of the decrypted wallet.

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			w := args[0]
			pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))

			return encryptWallet(w, pr)
		},
	}

	encryptWalletCmd.Flags().StringP("password", "p", "", "wallet password")
	return encryptWalletCmd
}

func encryptWallet(id string, pr PasswordReader) error {
	wlt, err := apiClient.Wallet(id)
	if err != nil {
		return err
	}

	if wlt.Meta.Encrypted {
		return wallet.ErrWalletEncrypted
	}

	if pr == nil {
		return wallet.ErrMissingPassword
	}

	pwd, err := pr.Password()
	if err != nil {
		return err
	}

	_, err = apiClient.EncryptWallet(id, string(pwd))
	return err
}
