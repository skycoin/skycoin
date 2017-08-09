package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"
)

func walletDirCmd() gcli.Command {
	name := "walletDir"
	return gcli.Command{
		Name:         name,
		Usage:        "Displays wallet folder address",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Flags: []gcli.Flag{
			gcli.BoolFlag{
				Name:  "j,json",
				Usage: "Returns the results in JSON format.",
			},
		},
		Action: func(c *gcli.Context) error {
			cfg := ConfigFromContext(c)
			jsonFmt := c.Bool("json")
			if jsonFmt {
				return printJson(struct {
					WltDir string `json:"walletDir"`
				}{
					WltDir: cfg.WalletDir,
				})
			}

			fmt.Println(cfg.WalletDir)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}
