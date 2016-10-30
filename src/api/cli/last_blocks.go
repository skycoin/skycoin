package cli

import (
	"fmt"

	"encoding/json"
	"strings"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/visor"
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:        "lastBlocks",
		ArgsUsage:   "Displays the content of the most recently N generated blocks.",
		Usage:       "[numberOfBlocks]",
		Description: "All results returned in JSON format.",
		Action:      getLastBlocks,
	}
	Commands = append(Commands, cmd)
}

func getLastBlocks(c *gcli.Context) error {
	num := c.Args().First()
	if num == "" {
		num = "1"
	}

	params := map[string]string{
		"num": num,
	}

	req := webrpc.NewRequest("get_lastblocks", params, "1")
	rsp, err := webrpc.Do(req, rpcAddress)
	if err != nil {
		return err
	}

	if rsp.Error != nil {
		return fmt.Errorf("rpc error code:%v, message:%v", rsp.Error.Code, rsp.Error.Message)
	}

	var blocks visor.ReadableBlocks
	if err := json.NewDecoder(strings.NewReader(rsp.Result)).Decode(&blocks); err != nil {
		return errJSONMarshal
	}

	d, err := json.MarshalIndent(blocks, "", "    ")
	if err != nil {
		return errJSONMarshal
	}

	fmt.Println(string(d))
	return nil
}
