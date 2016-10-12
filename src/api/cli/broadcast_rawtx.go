package cli

import (
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:        "broadcastTransaction",
		Description: "Broadcast a raw transaction to the network.",
		Usage:       "skycoin broadcastTransaction [transaction] ",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
