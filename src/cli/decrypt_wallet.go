package cli

import (
	"fmt"
	"path/filepath"

	"github.com/skycoin/skycoin/src/wallet"
)

import (
	gcli "github.com/spf13/cobra"
)

func decryptWalletCmd() *gcli.Command {
	decryptWalletCmd := &gcli.Command{
		Use:   "decryptWallet",
		Short: "Decrypt wallet",
		Long: fmt.Sprintf(`The default wallet (%s) will be used if no wallet was specified.

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`, cliConfig.FullWalletPath()),
		SilenceUsage: true,
		RunE: func(c *gcli.Command, _ []string) error {
			w, err := resolveWalletPath(cliConfig, "")
			if err != nil {
				return err
			}

			pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))

			wlt, err := decryptWallet(w, pr)
			switch err.(type) {
			case nil:
			case WalletLoadError:
				printHelp(c)
				return err
			default:
				return err
			}

			return ""
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
