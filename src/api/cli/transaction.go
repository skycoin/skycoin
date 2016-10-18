package cli

import gcli "gopkg.in/urfave/cli.v1"

func init() {
	cmd := gcli.Command{
		Name:      "transaction",
		Usage:     "Lists details of specific transaction",
		ArgsUsage: "[option] [transaction id]",
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
