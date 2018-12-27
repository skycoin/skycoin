package cli

import (
	"fmt"
	"strconv"

	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/api"
)

func richlistCmd() gcli.Command {
	name := "richlist"
	return gcli.Command{
		Name:         name,
		Usage:        "Returns top N address (default 20) balances (based on unspent outputs). Optionally include distribution addresses (exluded by default).",
		ArgsUsage:    "[top N addresses (20 default)] [include distribution addresses (false default)]",
		OnUsageError: onCommandUsageError(name),
		Action:       getRichlist,
	}
}

func getRichlist(c *gcli.Context) error {
	client := APIClientFromContext(c)

	num := c.Args().First()
	if num == "" {
		num = "20" // default to 20 addresses
	}

	dist := c.Args().Get(1)
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

	richlist, err := client.Richlist(params)
	if err != nil {
		return err
	}

	return printJSON(richlist)
}
