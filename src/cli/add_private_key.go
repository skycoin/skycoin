package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

func addPrivateKeyCmd() *cobra.Command {
	addPrivateKeyCmd := &cobra.Command{
		Short: "Add a private key to specific wallet",
		Use:   "addPrivateKey [flags] [private key]",
		Long: fmt.Sprintf(`Add a private key to specific wallet, the default
    wallet (%s) will be used if the wallet file or path is not specified.

    This method only works on "collection" type wallets.
    Use "skycoin-cli walletCreate -t collection" to create a "collection" type wallet.

    Use caution when using the "-p" command. If you have command
    history enabled your wallet encryption password can be recovered from the
    history log. If you do not include the "-p" option you will be prompted to
    enter your password after you enter your command.`, cliConfig.FullWalletPath()),
		SilenceUsage:          true,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(c *cobra.Command, args []string) error {
			// get private key
			skStr := args[0]
			if skStr == "" {
				return c.Help()
			}

			// get wallet file path
			walletFile, err := c.Flags().GetString("wallet-file")
			if err != nil {
				return err
			}

			w, err := resolveWalletPath(cliConfig, walletFile)
			if err != nil {
				return err
			}

			password, err := c.Flags().GetString("password")
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
		},
	}

	addPrivateKeyCmd.Flags().StringP("wallet-file", "f", "", "wallet file or path. If no path is specified your default wallet path will be used.")
	addPrivateKeyCmd.Flags().StringP("password", "p", "", "wallet password")

	return addPrivateKeyCmd
}

// AddPrivateKey adds a private key to a wallet.Wallet. Caller should save the wallet afterwards
func AddPrivateKey(wlt *wallet.CollectionWallet, key string) error {
	sk, err := cipher.SecKeyFromHex(key)
	if err != nil {
		return fmt.Errorf("invalid private key: %s, must be a hex string of length 64", key)
	}

	pk, err := cipher.PubKeyFromSecKey(sk)
	if err != nil {
		return err
	}

	addr := wlt.AddressConstructor()(pk)

	entry := wallet.Entry{
		Address: addr,
		Public:  pk,
		Secret:  sk,
	}

	return wlt.AddEntry(entry)
}

// AddPrivateKeyToFile adds a private key to a wallet based on filename.
// Will save the wallet after modifying.
func AddPrivateKeyToFile(walletFile, key string, pr PasswordReader) error {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return WalletLoadError{err}
	}

	if wlt.Type() != wallet.WalletTypeCollection {
		return fmt.Errorf("only %q type wallets can have keypairs added manually", wallet.WalletTypeCollection)
	}

	if pr == nil && wlt.IsEncrypted() {
		return wallet.ErrMissingPassword
	}

	addKey := func(w *wallet.CollectionWallet, key string) error {
		return AddPrivateKey(w, key)
	}

	if wlt.IsEncrypted() {
		addKey = func(w *wallet.CollectionWallet, key string) error {
			password, err := pr.Password()
			if err != nil {
				return err
			}

			return wallet.GuardUpdate(w, password, func(wlt wallet.Wallet) error {
				return AddPrivateKey(wlt.(*wallet.CollectionWallet), key)
			})
		}
	}

	if err := addKey(wlt.(*wallet.CollectionWallet), key); err != nil {
		return err
	}

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return err
	}

	if err := wallet.Save(wlt, dir); err != nil {
		return WalletSaveError{err}
	}

	return nil
}
