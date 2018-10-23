package cli

import (
	"fmt"
	"path/filepath"

	gcli "github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

func addPrivateKeyCmd() *gcli.Command {
	addPrivateKeyCmd := &gcli.Command{
		Short: "Add a private key to specific wallet",
		Use: "addPrivateKey [flags] [private key]",
		Long: fmt.Sprintf(`Add a private key to specific wallet, the default
    wallet (%s) will be
    used if the wallet file or path is not specified

    Use caution when using the "-p" command. If you have command
    history enabled your wallet encryption password can be recovered from the
    history log. If you do not include the "-p" option you will be prompted to
    enter your password after you enter your command.`, cliConfig.FullWalletPath()),
        SilenceUsage: true,
		Args: gcli.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(c *gcli.Command, args []string) error {
			// get private key
			skStr := args[0]
			if skStr == "" {
				return c.Help()
			}

			// get wallet file path
			w, err := resolveWalletPath(cliConfig, walletFile)
			if err != nil {
				return err
			}

			pr := NewPasswordReader([]byte(password))

			err = AddPrivateKeyToFile(w, skStr, pr)

			switch err.(type) {
			case nil:
				fmt.Println("success")
				return nil
			case WalletLoadError:
				printHelp(c)
				return err
			default:
				return err
			}
			return nil
		},
	}

	addPrivateKeyCmd.Flags().StringVarP(&walletFile, "wallet-file", "f", "", "wallet file or path. If no path is specified your default wallet path will be used.")
	addPrivateKeyCmd.Flags().StringVarP(&password, "password", "p", "", "Wallet password")

	return addPrivateKeyCmd
}

// AddPrivateKey adds a private key to a *wallet.Wallet. Caller should save the wallet afterwards
func AddPrivateKey(wlt *wallet.Wallet, key string) error {
	sk, err := cipher.SecKeyFromHex(key)
	if err != nil {
		return fmt.Errorf("invalid private key: %s, must be a hex string of length 64", key)
	}

	pk, err := cipher.PubKeyFromSecKey(sk)
	if err != nil {
		return err
	}

	addr := cipher.AddressFromPubKey(pk)

	entry := wallet.Entry{
		Address: addr,
		Public:  pk,
		Secret:  sk,
	}

	return wlt.AddEntry(entry)
}

// AddPrivateKeyToFile adds a private key to a wallet based on filename.  Will save the wallet after modifying.
func AddPrivateKeyToFile(walletFile, key string, pr PasswordReader) error {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return WalletLoadError{err}
	}

	switch pr.(type) {
	case nil:
		if wlt.IsEncrypted() {
			return wallet.ErrMissingPassword
		}
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return err
		}

		if !wlt.IsEncrypted() && len(p) != 0 {
			return wallet.ErrWalletNotEncrypted
		}
	}

	addKey := func(w *wallet.Wallet, key string) error {
		return AddPrivateKey(w, key)
	}

	if wlt.IsEncrypted() {
		addKey = func(w *wallet.Wallet, key string) error {
			password, err := pr.Password()
			if err != nil {
				return err
			}

			return w.GuardUpdate(password, func(wlt *wallet.Wallet) error {
				return AddPrivateKey(wlt, key)
			})
		}
	}

	if err := addKey(wlt, key); err != nil {
		return err
	}

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return err
	}

	if err := wlt.Save(dir); err != nil {
		return WalletSaveError{err}
	}

	return nil
}
