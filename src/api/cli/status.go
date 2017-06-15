package cli

import (
	"bytes"
	"fmt"

	"encoding/json"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func statusCMD() gcli.Command {
	name := "status"
	return gcli.Command{
		Name:         name,
		Usage:        "Check the status of current skycoin node",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			var status = struct {
				webrpc.StatusResult
				RPCAddress string `json:"webrpc_address"`
			}{
				RPCAddress: cfg.RPCAddress,
			}

			req, err := webrpc.NewRequest("get_status", nil, "1")
			if err != nil {
				return fmt.Errorf("create rpc request failed: %v", err)
			}

			rsp, err := webrpc.Do(req, cfg.RPCAddress)
			if err != nil {
				return fmt.Errorf("do rpc request failed: %v", err)
			}

			if rsp.Error != nil {
				return fmt.Errorf("do rpc request failed: %+v", *rsp.Error)
			}

			var rlt webrpc.StatusResult
			if err := json.NewDecoder(bytes.NewBuffer(rsp.Result)).Decode(&rlt); err != nil {
				return errJSONUnmarshal
			}

			status.StatusResult = rlt

			d, err := json.MarshalIndent(status, "", "    ")
			if err != nil {
				return errJSONMarshal
			}
			fmt.Println(string(d))
			return nil
		},
	}
}
