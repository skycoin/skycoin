package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

func listAddressesCMD() gcli.Command {
	name := "listAddresses"
	return gcli.Command{
		Name:         name,
		Usage:        "Lists all addresses in a given wallet",
		ArgsUsage:    "[walletName]",
		OnUsageError: onCommandUsageError(name),
		Action:       listAddresses,
	}
	// Commands = append(Commands, cmd)
}

func listAddresses(c *gcli.Context) error {
	// get wallet name
	w := c.Args().First()
	if w == "" {
		w = filepath.Join(cfg.WalletDir, cfg.DefaultWalletName)
	}

	if !strings.HasSuffix(w, walletExt) {
		return errWalletName
	}

	if filepath.Base(w) == w {
		w = filepath.Join(cfg.WalletDir, w)
	}

	wlt, err := wallet.Load(w)
	if err != nil {
		return err
	}

	addrs := wlt.GetAddresses()

	s, err := FormatAddressesAsJson(addrs)
	if err != nil {
		return err
	}

	fmt.Println(s)

	return nil
}
