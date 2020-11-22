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
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getLastBlocks,
	}
}

func getLastBlocks(_ *cobra.Command, args []string) error {
	num := args[0]
	if num == "" {
		num = "1"
	}

	n, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid block number, %s", err)
	}

	blocks, err := apiClient.LastBlocks(n)
	if err != nil {
		return err
	}

	return printJSON(blocks)
}
