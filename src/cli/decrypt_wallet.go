package cli

import (
	"path/filepath"

	gcli "github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/wallet"
)

func decryptWalletCmd() *gcli.Command {
	decryptWalletCmd := &gcli.Command{
		Args:  gcli.ExactArgs(1),
		Use:   "decryptWallet [wallet]",
		Short: "Decrypt a wallet",
		Long: `Decrypt an encrypted wallet. The decrypted wallet will be written
    on the filesystem in place of the encrypted wallet.

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		RunE: func(c *gcli.Command, args []string) error {
			w := args[0]
			pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))

			_, err := decryptWallet(w, pr)
			switch err.(type) {
			case nil:
			case WalletLoadError:
				printHelp(c)
				return err
			default:
				return err
			}

			return nil
		},
	}

	decryptWalletCmd.Flags().StringP("password", "p", "", "wallet password")

	return decryptWalletCmd
}

func decryptWallet(walletFile string, pr PasswordReader) (wallet.Wallet, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, WalletLoadError{err}
	}

	if !wlt.IsEncrypted() {
		return nil, wallet.ErrWalletNotEncrypted
	}

	if pr == nil {
		return nil, wallet.ErrMissingPassword
	}

	wltPassword, err := pr.Password()
	if err != nil {
		return nil, err
	}

	unlockedWlt, err := wallet.Unlock(wlt, wltPassword)
	if err != nil {
		return nil, err
	}

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return nil, err
	}

	// save the wallet
	if err := wallet.Save(unlockedWlt, dir); err != nil {
		return nil, WalletLoadError{err}
	}

	return unlockedWlt, nil
}
