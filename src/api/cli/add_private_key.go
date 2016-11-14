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

func init() {
	cmd := gcli.Command{
		Name:      "addPrivateKey",
		Usage:     "Add a private key to specific wallet",
		ArgsUsage: "[private key]",
		Description: `Add a private key to specific wallet, the default
		wallet file will be used if the wallet file or path is not specified`,
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] private key will be added to this wallet",
			},
		},
		Action: func(c *gcli.Context) error {
			// get wallet file path
			w := c.String("f")

			// get private key
			skStr := c.Args().First()
			if skStr == "" {
				return errors.New("private key value is empty")
			}

			if w == "" {
				w = filepath.Join(walletDir, defaultWalletName)
			}

			if !strings.HasSuffix(w, walletExt) {
				return errWalletName
			}

			// only wallet file name, no path.
			if filepath.Base(w) == w {
				w = filepath.Join(walletDir, w)
			}

			wlt, err := wallet.Load(w)
			if err != nil {
				return errLoadWallet
			}

			sk, err := cipher.SecKeyFromHex(skStr)
			if err != nil {
				return fmt.Errorf("invalid private key, %v", err)
			}

			pk := cipher.PubKeyFromSecKey(sk)
			addr := cipher.AddressFromPubKey(pk)

			entry := wallet.WalletEntry{
				Address: addr,
				Public:  pk,
				Secret:  sk,
			}

			if err := wlt.AddEntry(entry); err != nil {
				return err
			}

			dir, err := filepath.Abs(filepath.Dir(w))
			if err != nil {
				return err
			}

			if err := wlt.Save(dir); err != nil {
				return errors.New("save wallet failed")
			}

			fmt.Println("success")

			return nil
		},
	}
	Commands = append(Commands, cmd)
}
