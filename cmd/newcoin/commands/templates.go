package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/templates"
)

func init() {
	RootCmd.AddCommand(templatesCmd)
}

var templatesCmd = &cobra.Command{
	Use:   "templates [output-dir]",
	Short: "Export embedded templates to a directory",
	Long: `Exports the embedded newcoin templates to a directory for customization.

If no output directory is specified, templates are printed to stdout.
If an output directory is specified, template files are written there.

After exporting, you can edit the templates and use them with createcoin:

    skycoin newcoin templates ./my-templates
    # edit templates in ./my-templates/
    skycoin newcoin createcoin --coin mycoin --template-dir ./my-templates`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		templateFiles := []string{
			templates.CoinTemplate,
			templates.CommandTemplate,
			templates.CoinTestTemplate,
			templates.ParamsTemplate,
		}

		// No output dir — print all templates to stdout
		if len(args) == 0 {
			for _, name := range templateFiles {
				data, err := templates.FS.ReadFile(name)
				if err != nil {
					return fmt.Errorf("failed to read embedded template %s: %w", name, err)
				}
				fmt.Printf("--- %s ---\n", name)
				fmt.Print(string(data))
				fmt.Println()
			}
			return nil
		}

		// Write templates to output directory
		outDir := args[0]
		if err := os.MkdirAll(outDir, 0750); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		for _, name := range templateFiles {
			data, err := templates.FS.ReadFile(name)
			if err != nil {
				return fmt.Errorf("failed to read embedded template %s: %w", name, err)
			}

			outPath := filepath.Join(outDir, name)
			if err := os.WriteFile(outPath, data, 0644); err != nil { //nolint:gosec
				return fmt.Errorf("failed to write %s: %w", outPath, err)
			}
			fmt.Printf("Wrote %s\n", outPath)
		}

		return nil
	},
}
