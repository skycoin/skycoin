package cli

import (
	"fmt"

	"strconv"

	gcli "github.com/spf13/cobra"
)

func lastBlocksCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Displays the content of the most recently N generated blocks",
		Use:   "lastBlocks [numberOfBlocks]",
        Args:  gcli.MaximumNArgs(1),
        DisableFlagsInUseLine: true,
        SilenceUsage: true,
		RunE:  getLastBlocks,
	}
}

func getLastBlocks(c *gcli.Command, args []string) error {
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
