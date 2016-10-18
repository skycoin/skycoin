package cli

import gcli "gopkg.in/urfave/cli.v1"

func init() {
	cmd := gcli.Command{
		Name:      "broadcastTransaction",
		Usage:     "Broadcast a raw transaction to the network.",
		ArgsUsage: "[transaction]",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
