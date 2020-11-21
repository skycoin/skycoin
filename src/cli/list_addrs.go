package cli

import (
	"fmt"

	"github.com/spf13/cobra"
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
	addrs, err := getWalletAddresses(args[0])
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

func getWalletAddresses(id string) ([]string, error) {
	wlt, err := apiClient.Wallet(id)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, e := range wlt.Entries {
		addrs = append(addrs, e.Address)
	}
	return addrs, nil
}
