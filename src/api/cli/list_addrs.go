package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "gopkg.in/urfave/cli.v1"
)

func init() {
	cmd := gcli.Command{
		Name:        "listAddresses",
		Usage:       "Lists all addresses in a given wallet.",
		Description: "All results returned in JSON format.",
		ArgsUsage:   "[walletName]",
		Action:      listAddresses,
	}
	Commands = append(Commands, cmd)
}

func listAddresses(c *gcli.Context) error {
	// get wallet name
	w := c.Args().First()
	if w == "" {
		w = filepath.Join(walletDir, defaultWalletName)
	}

	// check if the wallet does exist
	if _, err := os.Stat(w); os.IsNotExist(err) {
		return err
	}

	wlt := wallet.Wallet{
		Meta: map[string]string{
			"filename": filepath.Base(w),
		},
	}
	fp, err := filepath.Abs(filepath.Dir(w))
	if err != nil {
		return err
	}
	if err := wlt.Load(fp); err != nil {
		return err
	}
	addrs := wlt.GetAddresses()
	var rlt = struct {
		Addresses []string `json:"addresses"`
	}{
		make([]string, len(addrs)),
	}

	for i, a := range addrs {
		rlt.Addresses[i] = a.String()
	}

	d, err := json.MarshalIndent(rlt, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(d))
	return nil
}
