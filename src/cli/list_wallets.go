package cli

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/spf13/cobra"
)

// WalletEntry represents an enty in a wallet file
type WalletEntry struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	AddressNum int    `json:"address_num"`
}

func listWalletsCmd() *gcli.Command {
	return &gcli.Command{
		Short:                 "Lists all wallets stored in the wallet directory",
		Use:                   "listWallets",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  listWallets,
	}
}

func listWallets(_ *gcli.Command, _ []string) error {
	var wlts struct {
		Wallets []WalletEntry `json:"wallets"`
	}

	entries, err := ioutil.ReadDir(cliConfig.WalletDir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, walletExt) {
				continue
			}

			path := filepath.Join(cliConfig.WalletDir, name)
			w, err := wallet.Load(path)
			if err != nil {
				return WalletLoadError{err}
			}
			wlts.Wallets = append(wlts.Wallets, WalletEntry{
				Name:       name,
				Label:      w.Label(),
				AddressNum: w.EntriesLen(),
			})
		}
	}

	return printJSON(wlts)
}
