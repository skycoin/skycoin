package cli

import (
	"fmt"
	"strconv"

	gcli "github.com/urfave/cli"
)

func blocksCmd() gcli.Command {
	name := "blocks"
	return gcli.Command{
		Name:         name,
		Usage:        "Lists the content of a single block or a range of blocks",
		ArgsUsage:    "[starting block or single block seq] [ending block seq]",
		Action:       getBlocks,
		OnUsageError: onCommandUsageError(name),
	}
	// Commands = append(Commands, cmd)
}

func getBlocks(c *gcli.Context) error {
	rpcClient := RpcClientFromContext(c)

	// get start
	start := c.Args().Get(0)
	end := c.Args().Get(1)
	if end == "" {
		end = start
	}

	if start == "" {
		gcli.ShowSubcommandHelp(c)
		return nil
	}

	s, err := strconv.ParseUint(start, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid block seq: %v, must be unsigned integer", start)
	}

	e, err := strconv.ParseUint(end, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid block seq: %v, must be unsigned integer", end)
	}

	rlt, err := rpcClient.GetBlocks(s, e)
	if err != nil {
		return err
	}

	return printJson(rlt)
}
