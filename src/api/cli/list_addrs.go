package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "listAddresses",
		Usage:     "Lists all addresses in a given wallet",
		ArgsUsage: "[walletName]",
		Action:    listAddresses,
	}
	Commands = append(Commands, cmd)
}

func listAddresses(c *gcli.Context) error {
	// get wallet name
	w := c.Args().First()
	if w == "" {
		w = filepath.Join(walletDir, defaultWalletName)
	}

	if !strings.HasSuffix(w, walletExt) {
		return errWalletName
	}

	if filepath.Base(w) == w {
		w = filepath.Join(walletDir, w)
	}

	wlt, err := wallet.Load(w)
	if err != nil {
		return errLoadWallet
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
		return errors.New("json marshal failed")
	}
	fmt.Println(string(d))
	return nil
}
