package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func blocksCmd() *cobra.Command {
	blocksCmd := &cobra.Command{
		Short:                 "Lists the content of a single block or a range of blocks",
		Use:                   "blocks [starting block or single block seq] [ending block seq]",
		Args:                  cobra.RangeArgs(1, 2),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getBlocks,
	}

	return blocksCmd
}

func getBlocks(_ *cobra.Command, args []string) error {
	var start, end string
	start = args[0]

	if len(args) == 1 {
		end = start
	} else {
		end = args[1]
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
