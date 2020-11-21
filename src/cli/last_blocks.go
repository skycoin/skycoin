package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func lastBlocksCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Displays the content of the most recently N generated blocks",
		Use:                   "lastBlocks [numberOfBlocks]",
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getLastBlocks,
	}
}

func getLastBlocks(_ *cobra.Command, args []string) error {
	n := uint64(1)
	if len(args) > 0 {
		var err error
		n, err = strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid block number, %s", err)
		}
	}

	blocks, err := apiClient.LastBlocks(n)
	if err != nil {
		return err
	}

	return printJSON(blocks)
}
