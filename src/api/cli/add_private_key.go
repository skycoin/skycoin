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

			wlt, err := wallet.Load(w)
			if err != nil {
				errorWithHelp(c, err)
				return nil
			}

			sk, err := cipher.SecKeyFromHex(skStr)
			if err != nil {
				return fmt.Errorf("invalid private key: %s, must be an hex string of length 64", skStr)
			}

			pk := cipher.PubKeyFromSecKey(sk)
			addr := cipher.AddressFromPubKey(pk)

			entry := wallet.Entry{
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
	// Commands = append(Commands, cmd)
}
