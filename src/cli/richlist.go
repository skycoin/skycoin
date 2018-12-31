package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/api"
)

func richlistCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Get skycoin richlist",
		Long:                  "Returns top N address (default 20) balances (based on unspent outputs). Optionally include distribution addresses (exluded by default).",
		Use:                   "richlist [top N addresses (20 default)] [include distribution addresses (false default)]",
		Args:                  cobra.MaximumNArgs(2),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getRichlist,
	}
}

func getRichlist(_ *cobra.Command, args []string) error {
	// default values
	num := "20"
	dist := "false"

	switch len(args) {
	case 1:
		num = args[0]
	case 2:
		num = args[0]
		dist = args[1]
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
