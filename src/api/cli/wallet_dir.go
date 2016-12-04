package cli

import (
	"encoding/json"
	"fmt"

	gcli "github.com/urfave/cli"
)

func walletDirCMD() gcli.Command {
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
			jsonFmt := c.Bool("json")
			if jsonFmt {
				var rlt = struct {
					WltDir string `json:"walletDir"`
				}{
					cfg.WalletDir,
				}
				d, err := json.MarshalIndent(rlt, "", "    ")
				if err != nil {
					return errJSONMarshal
				}
				fmt.Println(string(d))
				return nil
			}

			fmt.Println(cfg.WalletDir)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}
