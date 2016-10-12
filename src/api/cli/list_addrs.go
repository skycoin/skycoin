package cli

import gcli "gopkg.in/urfave/cli.v1"

func init() {
	cmd := gcli.Command{
		Name:        "listAddresses",
		Usage:       "Lists all addresses in a given wallet.",
		Description: "All results returned in JSON format.",
		ArgsUsage:   "[walletName]",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
