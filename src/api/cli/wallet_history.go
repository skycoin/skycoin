package cli

import (
	"fmt"

	"path/filepath"

	"strings"

	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

func walletHisCMD() gcli.Command {
	name := "walletHistory"
	return gcli.Command{
		Name:         name,
		Usage:        "Display the transaction history of specific wallet",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.",
			},
		},
		Action: func(c *gcli.Context) error {
			f := c.String("f")
			if f == "" {
				f = filepath.Join(cfg.WalletDir, cfg.DefaultWalletName)
			}

			// check the file extension.
			if !strings.HasSuffix(f, walletExt) {
				return errWalletName
			}

			// check if file name contains path.
			if filepath.Base(f) != f {
				af, err := filepath.Abs(f)
				if err != nil {
					return fmt.Errorf("invalid wallet file:%v, err:%v", f, err)
				}
				f = af
			} else {
				f = filepath.Join(cfg.WalletDir, f)
			}

			addrs, err := getAddresses(f)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func getAddresses(f string) ([]string, error) {
	wlt, err := wallet.Load(f)
	if err != nil {
		return nil, err
	}

	addrs := make([]string, len(wlt.Entries))
	for i, entry := range wlt.Entries {
		addrs[i] = entry.Address.String()
	}
	return addrs, nil
}
