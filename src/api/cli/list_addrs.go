package cli

import (
	"fmt"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

func listAddressesCmd() gcli.Command {
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
	cfg := ConfigFromContext(c)

	// get wallet name
	w, err := resolveWalletPath(cfg, c.Args().First())
	if err != nil {
		return err
	}

	wlt, err := wallet.Load(w)
	if err != nil {
		return WalletLoadError(err)
	}

	addrs := wlt.GetAddresses()

	s, err := FormatAddressesAsJson(addrs)
	if err != nil {
		return err
	}

	fmt.Println(s)

	return nil
}
