package cli

import (
	"fmt"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/spf13/cobra"
)

func listAddressesCmd() *gcli.Command {
	return &gcli.Command{
		Short:                 "Lists all addresses in a given wallet",
		Use:                   "listAddresses [walletName]",
		Args:                  gcli.ExactArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  listAddresses,
	}
}

func listAddresses(_ *gcli.Command, args []string) error {
	wlt, err := wallet.Load(args[0])
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
