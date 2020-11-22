package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

func addPrivateKeyCmd() *cobra.Command {
	// TODO -- allow private key to be entered privately, same as the password can be
	addPrivateKeyCmd := &cobra.Command{
		Short: "Add a private key to wallet",
		Use:   "addPrivateKey [wallet] [private key]",
		Long: `Add a private key to wallet.

    This method only works on "collection" type wallets.
    Use "skycoin-cli walletCreate -t collection" to create a "collection" type wallet.

    Use caution when using this from your shell. The private key will be recorded
    if your shell's history file, unless you disable the shell history.

    Use caution when using the "-p" command. If you have command
    history enabled your wallet encryption password can be recovered from the
    history log. If you do not include the "-p" option you will be prompted to
    enter your password after you enter your command.`,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		RunE: func(c *cobra.Command, args []string) error {
			walletFile := args[0]
			skStr := args[1]

			password, err := c.Flags().GetString("password")
			if err != nil {
				return err
			}
			pr := NewPasswordReader([]byte(password))

			err = AddPrivateKeyToFile(walletFile, skStr, pr)

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
