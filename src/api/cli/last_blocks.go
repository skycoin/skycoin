package cli

import (
	"errors"
	"fmt"

	"strconv"

	"github.com/skycoin/skycoin/src/api/webrpc"
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

	if n <= 0 {
		return errors.New("block number must >= 0")
	}

	param := []uint64{n}
	req, err := webrpc.NewRequest("get_lastblocks", param, "1")
	if err != nil {
		return fmt.Errorf("do rpc request failed: %v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return err
	}

	if rsp.Error != nil {
		return fmt.Errorf("do rpc request failed: %+v", *rsp.Error)
	}

	fmt.Println(string(rsp.Result))
	return nil
}
