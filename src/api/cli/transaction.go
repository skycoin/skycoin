package cli

import (
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:        "transaction",
		Description: "Lists details of specific transaction",
		Usage:       "skycointransaction  [option] [transaction id]",
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "j,json",
				Usage: "Returns the results in JSON format.",
			},
		},
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
