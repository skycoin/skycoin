package cli

import gcli "gopkg.in/urfave/cli.v1"

func init() {
	cmd := gcli.Command{
		Name:        "listWallets",
		Usage:       "Lists all wallets stored in the default wallet directory [display directory].",
		ArgsUsage:   "[options]",
		Description: "All results returned in JSON format",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
