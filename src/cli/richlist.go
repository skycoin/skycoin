package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/api"
)

func richlistCmd() cobra.Command {
	name := "richlist"
	return cobra.Command{
		Short:                 name,
		Long:                  "Returns top N address (default 20) balances (based on unspent outputs). Optionally include distribution addresses (exluded by default).",
		Use:                   "[top N addresses (20 default)] [include distribution addresses (false default)]",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getRichlist,
	}
}

func getRichlist(_ *cobra.Command, args []string) error {
	num := args[0]
	if num == "" {
		num = "20" // default to 20 addresses
	}

	dist := args[1]
	if dist == "" {
		dist = "false" // default to false
	}

	n, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("invalid number of addresses, %s", err)
	}

	d, err := strconv.ParseBool(dist)
	if err != nil {
		return fmt.Errorf("invalid (bool) flag for include distribution addresses, %s", err)
	}

	params := &api.RichlistParams{
		N:                   n,
		IncludeDistribution: d,
	}

	richlist, err := apiClient.Richlist(params)
	if err != nil {
		return err
	}

	return printJSON(richlist)
}
