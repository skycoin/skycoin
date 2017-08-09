package cli

import (
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
			status, err := GetStatus()
			if err != nil {
				return err
			}

			return printJson(struct {
				webrpc.StatusResult
				RPCAddress string `json:"webrpc_address"`
			}{
				StatusResult: *status,
				RPCAddress:   cfg.RPCAddress,
			})
		},
	}
}

func GetStatus() (*webrpc.StatusResult, error) {
	status := webrpc.StatusResult{}
	if err := DoRpcRequest(&status, "get_status", nil, "1"); err != nil {
		return nil, err
	}

	return &status, nil
}
