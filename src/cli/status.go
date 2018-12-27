package cli

import (
	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/api"
)

// StatusResult is printed by cli status command
type StatusResult struct {
	Status api.HealthResponse `json:"status"`
	Config ConfigStatus       `json:"cli_config"`
}

// ConfigStatus contains the configuration parameters loaded by the cli
type ConfigStatus struct {
	RPCAddress string `json:"webrpc_address"`
}

func statusCmd() gcli.Command {
	name := "status"
	return gcli.Command{
		Name:         name,
		Usage:        "Check the status of current skycoin node",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			client := APIClientFromContext(c)
			status, err := client.Health()
			if err != nil {
				return err
			}

			cfg := ConfigFromContext(c)

			return printJSON(StatusResult{
				Status: *status,
				Config: ConfigStatus{
					RPCAddress: cfg.RPCAddress,
				},
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
