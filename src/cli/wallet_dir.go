package cli

import (
	"fmt"

	gcli "github.com/spf13/cobra"
)

func walletDirCmd() *gcli.Command {
    walletDirCmd := &gcli.Command{
		Use:   "walletDir",
		Short: "Displays wallet folder address",
        Args: gcli.NoArgs,
        SilenceUsage: true,
		RunE: func(c *gcli.Command, args []string) error {
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

    walletDirCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Returns the results in JSON format.")
    return walletDirCmd
}
