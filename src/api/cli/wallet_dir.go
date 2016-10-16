package cli

import (
	"encoding/json"
	"fmt"

	gcli "gopkg.in/urfave/cli.v1"
)

func init() {
	cmd := gcli.Command{
		Name:  "walletDir",
		Usage: "Displays wallet folder address.",
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
					walletDir,
				}
				d, err := json.MarshalIndent(rlt, "", "    ")
				if err != nil {
					return err
				}
				fmt.Println(string(d))
				return nil
			}

			fmt.Println(walletDir)
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
