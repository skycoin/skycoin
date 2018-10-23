package cli

import (
	"fmt"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/spf13/cobra"
)

func listAddressesCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Lists all addresses in a given wallet",
		Use:   "listAddresses [walletName]",
        Args:  gcli.MaximumNArgs(1),
        DisableFlagsInUseLine: true,
        SilenceUsage: true,
		RunE:  listAddresses,
	}
}

func listAddresses(c *gcli.Command, args []string) error {
    var wltPath string
    if len(args) > 0 {
        wltPath = args[0]
    }

	// get wallet name
	w, err := resolveWalletPath(cliConfig, wltPath)
	if err != nil {
		return err
	}

	wlt, err := wallet.Load(w)
	if err != nil {
		return WalletLoadError{err}
	}

	addrs := wlt.GetAddresses()

	s, err := FormatAddressesAsJSON(addrs)
	if err != nil {
		return err
	}

	fmt.Println(s)

	return nil
}
