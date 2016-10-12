package cli

import (
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:        "listWallets",
		Description: "Lists all wallets stored in the default wallet directory [display directory]. All results returned in JSON format",
		Usage:       "skycoin listWallets [options]",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
