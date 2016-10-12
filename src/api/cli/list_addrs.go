package cli

import (
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:        "listAddresses",
		Description: "Lists all addresses in a given wallet. All results returned in JSON format.",
		Usage:       "skycoin listAddresses [walletName] ",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
