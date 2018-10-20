package cli

import (
	"fmt"
	"strconv"

	gcli "github.com/spf13/cobra"
)

func blocksCmd() *gcli.Command {
	blocksCmd := &gcli.Command{
		Short: "Lists the content of a single block or a range of blocks",
		Use:   "blocks [starting block or single block seq] [ending block seq]",
		Args:  gcli.RangeArgs(1, 2),
		RunE:  getBlocks,
	}

	return blocksCmd
}

func getBlocks(c *gcli.Command, args []string) error {
	start := args[0]
	end := args[1]

	if end == "" {
		end = start
	}

	s, err := strconv.ParseUint(start, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid block seq: %v, must be unsigned integer", start)
	}

	e, err := strconv.ParseUint(end, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid block seq: %v, must be unsigned integer", end)
	}

	rlt, err := apiClient.BlocksInRange(s, e)
	if err != nil {
		return err
	}

	return printJSON(rlt)
}
