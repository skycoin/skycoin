package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "blocks",
		Usage:     "Lists the content of a single block or a range of blocks",
		ArgsUsage: "[starting block or single block seq] [ending block seq]",
		Action:    getBlocks,
	}
	Commands = append(Commands, cmd)
}

func getBlocks(c *gcli.Context) error {
	// get start
	start := c.Args().Get(0)
	end := c.Args().Get(1)
	if end == "" {
		end = start
	}

	s, err := strconv.ParseUint(start, 10, 64)
	if err != nil {
		return errors.New("error block seq")
	}

	e, err := strconv.ParseUint(end, 10, 64)
	if err != nil {
		return errors.New("error block seq")
	}

	param := []uint64{s, e}

	req, err := webrpc.NewRequest("get_blocks", param, "1")
	if err != nil {
		return fmt.Errorf("create rpc request failed: %v", err)
	}

	rsp, err := webrpc.Do(req, rpcAddress)
	if err != nil {
		return fmt.Errorf("do rpc request failed: %v", err)
	}

	if rsp.Error != nil {
		return fmt.Errorf("do rpc request failed: %+v", *rsp.Error)
	}

	fmt.Println(string(rsp.Result))
	return nil
}
