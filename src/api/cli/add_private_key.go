package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

func addPrivateKeyCMD() gcli.Command {
	name := "addPrivateKey"
	return gcli.Command{
		Name:      name,
		Usage:     "Add a private key to specific wallet",
		ArgsUsage: "[private key]",
		Description: fmt.Sprintf(`Add a private key to specific wallet, the default
		wallet(%s/%s) will be
		used if the wallet file or path is not specified`,
			cfg.WalletDir, cfg.DefaultWalletName),
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] private key will be added to this wallet",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			// get private key
			skStr := c.Args().First()
			if skStr == "" {
				gcli.ShowSubcommandHelp(c)
				return nil
			}

			// get wallet file path
			w := c.String("f")
			if w == "" {
				w = filepath.Join(cfg.WalletDir, cfg.DefaultWalletName)
			}

			if !strings.HasSuffix(w, walletExt) {
				return errWalletName
			}

			// only wallet file name, no path.
			if filepath.Base(w) == w {
				w = filepath.Join(cfg.WalletDir, w)
			}

			err := AddPrivateKeyToFile(w, skStr)

			switch err.(type) {
			case nil:
				fmt.Println("success")
				return nil
			case WalletLoadError:
				errorWithHelp(c, err)
				return nil
			case WalletSaveError:
				return errors.New("Save wallet failed")
			default:
				return err
			}
		},
	}
	// Commands = append(Commands, cmd)
}

// PUBLIC

type WalletLoadError error
type WalletSaveError error

// Adds a private key to a *wallet.Wallet. Caller should save the wallet afterwards
func AddPrivateKey(wlt *wallet.Wallet, key string) error {
	sk, err := cipher.SecKeyFromHex(key)
	if err != nil {
		return fmt.Errorf("invalid private key: %s, must be an hex string of length 64", key)
	}

	pk := cipher.PubKeyFromSecKey(sk)
	addr := cipher.AddressFromPubKey(pk)

	entry := wallet.Entry{
		Address: addr,
		Public:  pk,
		Secret:  sk,
	}

	return wlt.AddEntry(entry)
}

// Adds a private key to a wallet based on filename.  Will save the wallet after modifying.
func AddPrivateKeyToFile(walletFile, key string) error {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return WalletLoadError(err)
	}

	if err := AddPrivateKey(wlt, key); err != nil {
		return err
	}

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return err
	}

	if err := wlt.Save(dir); err != nil {
		return WalletSaveError(err)
	}

	return nil
}
