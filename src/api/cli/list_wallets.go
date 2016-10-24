package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "gopkg.in/urfave/cli.v1"
)

type walletEntry struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	AddressNum int    `json:"address_num"`
}

func init() {
	cmd := gcli.Command{
		Name:        "listWallets",
		Usage:       "Lists all wallets stored in the default wallet directory [display directory].",
		ArgsUsage:   "[options]",
		Description: "All results returned in JSON format",
		Action: func(c *gcli.Context) error {
			var wlts struct {
				Wallets []walletEntry `json:"wallets"`
			}

			entries, err := ioutil.ReadDir(walletDir)
			if err != nil {
				return err
			}

			for _, e := range entries {
				if e.Mode().IsRegular() {
					name := e.Name()
					if !strings.HasSuffix(name, walletExt) {
						continue
					}

					path := filepath.Join(walletDir, name)
					w, err := wallet.Load(path)
					if err != nil {
						return fmt.Errorf("load wallet %v failed", name)
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
				return err
			}
			fmt.Println(string(d))
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
