package cli

import gcli "gopkg.in/urfave/cli.v1"

func init() {
	cmd := gcli.Command{
		Name:  "walletDir",
		Usage: "Displays wallet folder address.",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
