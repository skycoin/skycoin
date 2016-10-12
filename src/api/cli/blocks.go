package cli

import gcli "gopkg.in/urfave/cli.v1"

func init() {
	cmd := gcli.Command{
		Name:      "blocks",
		Usage:     "Lists the content of a single block of a range of blocks. Block results are always in JSON format.",
		ArgsUsage: "[option] [starting block or single block] [ending block]]",
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
