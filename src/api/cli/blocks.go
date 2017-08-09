package cli

import (
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/visor"
	gcli "github.com/urfave/cli"
)

func blocksCMD() gcli.Command {
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

	rlt, err := GetBlocks(s, e)
	if err != nil {
		return err
	}

	return printJson(rlt)
}

// PUBLIC

func GetBlocks(start, end uint64) (*visor.ReadableBlocks, error) {
	param := []uint64{start, end}
	blocks := visor.ReadableBlocks{}

	if err := DoRpcRequest(&blocks, "get_blocks", param, "1"); err != nil {
		return nil, err
	}

	return &blocks, nil
}

func GetBlocksBySeq(ss []uint64) (*visor.ReadableBlocks, error) {
	blocks := visor.ReadableBlocks{}

	if err := DoRpcRequest(&blocks, "get_blocks_by_seq", ss, "1"); err != nil {
		return nil, err
	}

	return &blocks, nil
}
