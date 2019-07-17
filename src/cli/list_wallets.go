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
		Short: "Lists all wallets stored in the wallet directory",
		Use:   "listWallets [wallet dir]",
		Long: `Lists all wallets stored in the wallet directory

    The [wallet dir] argument is optional. If not provided, defaults to $DATA_DIR/wallets`,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  gcli.MaximumNArgs(1),
		RunE:                  listWallets,
	}
}

func listWallets(_ *gcli.Command, args []string) error {
	var wlts struct {
		Directory string        `json:"directory"`
		Wallets   []WalletEntry `json:"wallets"`
	}

	wlts.Wallets = []WalletEntry{}

	dir := filepath.Join(cliConfig.DataDir, "wallets")
	if len(args) > 0 {
		dir = args[0]
	}

	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}

	wlts.Directory = dir

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, walletExt) {
				continue
			}

			path := filepath.Join(dir, name)
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
