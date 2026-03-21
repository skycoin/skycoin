package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func haltCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "halt",
		Short: "Shut down the running node",
		Long:  "Sends a shutdown request to the running node via the RPC API.",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := apiClient.Shutdown(); err != nil {
				return fmt.Errorf("failed to shutdown node: %v", err)
			}
			fmt.Println("Node shutdown initiated")
			return nil
		},
	}
}
