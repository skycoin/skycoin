package cli

import gcli "github.com/urfave/cli"

func init() {
	cmd := gcli.Command{
		Name:        "checkBalance",
		Description: "Check the balance of a wallet or specific address.",
		Usage:       "skycoin checkBalance  [option] [wallet path or address]",
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "w",
				Usage: "[wallet file or path], List balance of all addresses in a wallet.",
			},
			gcli.StringFlag{
				Name:  "a",
				Usage: "[address] List balance of specific address.",
			},
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
