package cli

import (
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:        "blocks",
		Description: "Lists the content of a single block of a range of blocks. Block results are always in JSON format.",
		Usage:       "blocks [option] [starting block or single block] [ending block]]",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
