// Package commands implements the config printer.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/config"
)

func init() {
	RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print fiber.toml config",
	Long: `Prints the default fiber.toml configuration file.

This provides a starting point for creating custom fibercoin configurations.
You can redirect the output to a file and then customize it:

Then edit mycoin.toml to customize your fibercoin parameters.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		// Read the embedded fiber.toml
		data, err := config.FS.ReadFile(config.FiberToml)
		if err != nil {
			return fmt.Errorf("failed to read embedded fiber.toml: %w", err)
		}

		// Print to stdout
		fmt.Print(string(data))
		return nil
	},
}
