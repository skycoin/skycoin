package cli

import (
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/api/webrpc"
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

	req, err := webrpc.NewRequest("get_blocks", param, "1")
	if err != nil {
		return nil, fmt.Errorf("create rpc request failed: %v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return nil, fmt.Errorf("do rpc request failed: %v", err)
	}

	if rsp.Error != nil {
		return nil, fmt.Errorf("rpc response error: %+v", *rsp.Error)
	}

	var rlt visor.ReadableBlocks
	if err := decodeJson(rsp.Result, &rlt); err != nil {
		return nil, err
	}

	return &rlt, nil
}

func GetBlocksBySeq(ss []uint64) (*visor.ReadableBlocks, error) {
	req, err := webrpc.NewRequest("get_blocks_by_seq", ss, "1")
	if err != nil {
		return nil, fmt.Errorf("create rpc request failed: %v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return nil, fmt.Errorf("do rpc request failed: %v", err)
	}

	if rsp.Error != nil {
		return nil, fmt.Errorf("rpc response error: %+v", *rsp.Error)
	}

	blks := visor.ReadableBlocks{}
	if err := decodeJson(rsp.Result, &blks); err != nil {
		return nil, err
	}

	return &blks, nil
}
