package cli

import (
	"errors"
	"fmt"
	"strings"

	"encoding/json"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "status",
		ArgsUsage: "check the status of current skycoin node.",
		Usage:     "[options]",
		Action: func(c *gcli.Context) error {
			var status = struct {
				RPCAddress string `json:"webrpc_address"`
				Running    bool   `json:"running"`
			}{
				RPCAddress: rpcAddress,
			}

			req := webrpc.NewRequest("get_status", nil, "1")
			rsp, err := webrpc.Do(req, rpcAddress)
			if err != nil {
				return errors.New("do request webrpc failed")
			}

			if rsp.Error != nil {
				return fmt.Errorf("webrpc request failed, code:%d, message:%s", rsp.Error.Code, rsp.Error.Message)
			}

			if strings.Contains(rsp.Result, "true") {
				status.Running = true
			}

			d, err := json.MarshalIndent(status, "", "    ")
			if err != nil {
				return errJSONMarshal
			}
			fmt.Println(string(d))
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
