package cli

import (
	"errors"
	"fmt"

	"strconv"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/visor"
	gcli "github.com/urfave/cli"
)

func lastBlocksCMD() gcli.Command {
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
	num := c.Args().First()
	if num == "" {
		num = "1"
	}

	n, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid block number, %s", err)
	}

	blocks, err := GetLastBlocks(n)

	if err != nil {
		return err
	}

	return printJson(blocks)
}

func GetLastBlocks(n uint64) (*visor.ReadableBlocks, error) {
	if n <= 0 {
		return nil, errors.New("block number must >= 0")
	}

	param := []uint64{n}
	req, err := webrpc.NewRequest("get_lastblocks", param, "1")
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
