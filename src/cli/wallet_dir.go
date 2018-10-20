package cli

import (
	"fmt"

	gcli "github.com/spf13/cobra"
)

func walletDirCmd() *gcli.Command {
	return &gcli.Command{
		Use:   "walletDir",
		Short: "Displays wallet folder address",
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
}
