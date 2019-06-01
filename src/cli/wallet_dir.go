package cli

import (
	"fmt"

	cobra "github.com/spf13/cobra"
)

func walletDirCmd() *cobra.Command {
	walletDirCmd := &cobra.Command{
		Use:          "walletDir",
		Short:        "Displays wallet folder address",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, _ []string) error {
			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			if jsonOutput {
				return printJSON(struct {
					WltDir string `json:"walletDir"`
				}{
					WltDir: cliConfig.WalletDir,
				})
			}

			fmt.Println(cliConfig.WalletDir)
			return nil
		},
	}

	walletDirCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")
	return walletDirCmd
}
