package cli

import (
	gcli "github.com/spf13/cobra"

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

func statusCmd() *gcli.Command {
	return &gcli.Command{
		Use:   "status",
		Short: "Check the status of current skycoin node",
		RunE: func(c *gcli.Command, args []string) error {
			status, err := apiClient.Health()
			if err != nil {
				return err
			}

			return printJSON(StatusResult{
				Status: *status,
				Config: ConfigStatus{
					RPCAddress: cliConfig.RPCAddress,
				},
			})
		},
	}
}

func showConfigCmd() *gcli.Command {
	return &gcli.Command{
		Use:   "showConfig",
		Short: "Show cli configuration",
		RunE: func(c *gcli.Command, args []string) error {
			return printJSON(cliConfig)
		},
	}
}
