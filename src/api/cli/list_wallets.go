package cli

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

type walletEntry struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	AddressNum int    `json:"address_num"`
}

func listWalletsCmd() gcli.Command {
	name := "listWallets"
	return gcli.Command{
		Name:         name,
		Usage:        "Lists all wallets stored in the wallet directory",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Action:       listWallets,
	}
	// Commands = append(Commands, cmd)
}

func listWallets(c *gcli.Context) error {
	cfg := ConfigFromContext(c)

	var wlts struct {
		Wallets []walletEntry `json:"wallets"`
	}

	entries, err := ioutil.ReadDir(cfg.WalletDir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, walletExt) {
				continue
			}

			path := filepath.Join(cfg.WalletDir, name)
			w, err := wallet.Load(path)
			if err != nil {
				return WalletLoadError(err)
			}
			wlts.Wallets = append(wlts.Wallets, walletEntry{
				Name:       name,
				Label:      w.GetLabel(),
				AddressNum: len(w.Entries),
			})
		}
	}

	return printJson(wlts)
}
