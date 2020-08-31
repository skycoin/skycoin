package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/SkycoinProject/skycoin/src/wallet"
)

func listAddressesCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Lists all addresses in a given wallet",
		Use:                   "listAddresses [wallet]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  listAddresses,
	}
}

func listAddresses(_ *cobra.Command, args []string) error {
	wlt, err := wallet.Load(args[0])
	if err != nil {
		return WalletLoadError{err}
	}

	addrs, err := wlt.GetAddresses()
	if err != nil {
		return err
	}

	s, err := FormatAddressesAsJSON(addrs)
	if err != nil {
		return err
	}

	fmt.Println(s)

	return nil
}
