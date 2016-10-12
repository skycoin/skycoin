package cli

import (
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:        "walletDir",
		Description: "Displays wallet folder address. ",
		Usage:       "skycoin walletDir",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
