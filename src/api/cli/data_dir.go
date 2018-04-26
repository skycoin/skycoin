package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"
)

func dataDirCmd() gcli.Command {
	name := "dataDir"
	return gcli.Command{
		Name:         name,
		Usage:        "Displays address directory of the data",
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
				return printJSON(struct {
					dataDir string `json:"dataDir"`
				}{
					dataDir: cfg.DataDir,
				})
			}

			fmt.Println(cfg.DataDir)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}
