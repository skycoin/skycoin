package cli

import gcli "gopkg.in/urfave/cli.v1"

func init() {
	cmd := gcli.Command{
		Name:      "version",
		Usage:     "List the current version of Skycoin components.",
		ArgsUsage: "[options]",
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
