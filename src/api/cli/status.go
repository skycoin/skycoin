package cli

import gcli "github.com/urfave/cli"

func init() {
	cmd := gcli.Command{
		Name:      "status",
		ArgsUsage: "check the status of current skycoin node.",
		Usage:     "[options]",
		Action: func(c *gcli.Context) error {
			// var status struct {
			// 	NodeAddress string `json:"node_addr"`
			// 	Running     bool   `json:"running"`
			// }

			return nil
		},
	}
	Commands = append(Commands, cmd)
}
