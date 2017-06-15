package cli

import (
	"encoding/json"
	"fmt"
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

func listWalletsCMD() gcli.Command {
	name := "listWallets"
	return gcli.Command{
		Name:         name,
		Usage:        "Lists all wallets stored in the default wallet directory",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
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
						return err
					}
					wlts.Wallets = append(wlts.Wallets, walletEntry{
						Name:       name,
						Label:      w.GetLabel(),
						AddressNum: len(w.Entries),
					})
				}
			}
			d, err := json.MarshalIndent(wlts, "", "    ")
			if err != nil {
				return errJSONMarshal
			}
			fmt.Println(string(d))
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}
