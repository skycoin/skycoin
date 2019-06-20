package cli

import (
	"github.com/spf13/cobra"
)

func addresscountCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Get the count of addresses with unspent outputs (coins).",
		Long:                  "Returns the count of all addresses that currently have unspent outputs (coins) associated with them.",
		Use:                   "addresscount",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getAddresscount,
	}
}

func getAddresscount(_ *cobra.Command, _ []string) error {
	addresscount, err := apiClient.AddressCount()
	if err != nil {
		return err
	}

	return printJSON(addresscount)
}
