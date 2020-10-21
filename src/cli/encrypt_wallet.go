package cli

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
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
			cryptoType, err := crypto.CryptoTypeFromString(c.Flag("crypto-type").Value.String())
			if err != nil {
				printHelp(c)
				return err
			}

			pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))

			_, err = encryptWallet(w, pr, cryptoType)
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

	encryptWalletCmd.Flags().StringP("password", "p", "", "wallet password")
	encryptWalletCmd.Flags().StringP("crypto-type", "x", "scrypt-chacha20poly1305", "The crypto type for wallet encryption, can be scrypt-chacha20poly1305 or sha256-xor")
	return encryptWalletCmd
}

func encryptWallet(walletFile string, pr PasswordReader, cryptoType crypto.CryptoType) (wallet.Wallet, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, WalletLoadError{err}
	}

	if wlt.IsEncrypted() {
		return nil, wallet.ErrWalletEncrypted
	}

	if pr == nil {
		return nil, wallet.ErrMissingPassword
	}

	password, err := pr.Password()
	if err != nil {
		return nil, err
	}

	if wlt.CryptoType() != cryptoType {
		wlt.SetCryptoType(cryptoType)
	}

	if err := wlt.Lock(password); err != nil {
		return nil, err
	}

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return nil, err
	}

	// save the wallet
	if err := wallet.Save(wlt, dir); err != nil {
		return nil, WalletLoadError{err}
	}

	return wlt, nil
}
