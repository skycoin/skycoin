package cli

import (
	"fmt"

	"strconv"

	gcli "github.com/urfave/cli"
)

func lastBlocksCmd() gcli.Command {
	name := "lastBlocks"
	return gcli.Command{
		Name:         name,
		Usage:        "Displays the content of the most recently N generated blocks",
		ArgsUsage:    "[numberOfBlocks]",
		OnUsageError: onCommandUsageError(name),
		Action:       getLastBlocks,
	}
	// Commands = append(Commands, cmd)
}

func getLastBlocks(c *gcli.Context) error {
	rpcClient := RpcClientFromContext(c)

	num := c.Args().First()
	if num == "" {
		num = "1"
	}

	n, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid block number, %s", err)
	}

	blocks, err := rpcClient.GetLastBlocks(n)

	if err != nil {
		return err
	}

	return printJson(blocks)
}
