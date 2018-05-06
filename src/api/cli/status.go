package cli

import (
	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/api/webrpc"
)

// StatusResult is printed by cli status command
type StatusResult struct {
	webrpc.StatusResult
	RPCAddress string `json:"webrpc_address"`
	UseCSRF    bool   `json:"use_csrf"`
}

func statusCmd() gcli.Command {
	name := "status"
	return gcli.Command{
		Name:         name,
		Usage:        "Check the status of current skycoin node",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rpcClient := RPCClientFromContext(c)
			status, err := rpcClient.GetStatus()
			if err != nil {
				return err
			}

			cfg := ConfigFromContext(c)

			return printJSON(StatusResult{
				StatusResult: *status,
				RPCAddress:   cfg.RPCAddress,
				UseCSRF:      cfg.UseCSRF,
			})
		},
	}
}

func showConfigCmd() gcli.Command {
	name := "showConfig"
	return gcli.Command{
		Name:         name,
		Usage:        "Show cli configuration",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			cfg := ConfigFromContext(c)
			return printJSON(cfg)
		},
	}
}
