// Package commands implements the newcoin template generator.
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/calvin"

	"github.com/skycoin/skycoin/src/fiber"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/useragent"
	"github.com/skycoin/skycoin/templates"
)

const (
	// Version is the CLI version
	Version = "0.2"
)

var (
	log                  = logging.MustGetLogger("newcoin")
	coinName             string
	templateDir          string
	coinTemplateFile     string
	commandTemplateFile  string
	coinTestTemplateFile string
	paramsTemplateFile   string
	configDir            string
	configFile           string
)

func init() {
	createCoinCmd.Flags().SortFlags = false
	createCoinCmd.Flags().StringVarP(&coinName, "coin", "c", "skycoin", "name of the coin to create")
	createCoinCmd.Flags().StringVarP(&templateDir, "template-dir", "d", "", "template directory path (uses embedded templates if empty)")
	createCoinCmd.Flags().StringVarP(&coinTemplateFile, "coin-template-file", "e", templates.CoinTemplate, "coin template file (importable)")
	createCoinCmd.Flags().StringVarP(&commandTemplateFile, "command-template-file", "f", templates.CommandTemplate, "command template file (executable)")
	createCoinCmd.Flags().StringVarP(&coinTestTemplateFile, "coin-test-template-file", "g", templates.CoinTestTemplate, "coin test template file")
	createCoinCmd.Flags().StringVarP(&paramsTemplateFile, "params-template-file", "i", templates.ParamsTemplate, "params template file")
	createCoinCmd.Flags().StringVarP(&configDir, "config-dir", "j", "./", "config directory path")
	createCoinCmd.Flags().StringVarP(&configFile, "config-file", "k", "config/fiber.toml", "config file path")
	RootCmd.AddCommand(createCoinCmd)
}

// RootCmd represents the base command for the application
var RootCmd = &cobra.Command{
	Use:   "newcoin",
	Short: "create a new fibercoin",
	Long: calvin.AsciiFont("newcoin") + `
create a new fibercoin`,
}

var createCoinCmd = &cobra.Command{
	Use:   "createcoin",
	Short: "Create a new coin from a template file",
	RunE: func(_ *cobra.Command, _ []string) error {
		if err := validateCoinName(coinName); err != nil {
			return err
		}
		configFilepath := filepath.Join(configDir, configFile)
		if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
			return err
		}
		config, err := fiber.NewConfig(configFile, configDir)
		if err != nil {
			log.Errorf("failed to create new fiber coin config")
			return err
		}
		coinDir := fmt.Sprintf("./cmd/%s", coinName)
		err = os.MkdirAll(coinDir, 0750)
		if err != nil {
			log.Errorf("failed to create new coin directory %s", coinDir)
			return err
		}
		commandsDir := fmt.Sprintf("./cmd/%s/commands", coinName)
		err = os.MkdirAll(commandsDir, 0750)
		if err != nil {
			log.Errorf("failed to create new coin directory %s", coinDir)
			return err
		}
		coinFilePath := fmt.Sprintf("./cmd/%[1]s/%[1]s.go", coinName)
		coinFile, err := os.Create(coinFilePath) //nolint:gosec
		if err != nil {
			log.Errorf("failed to create new coin file %s", coinFilePath)
			return err
		}
		defer coinFile.Close() //nolint
		commandFilePath := fmt.Sprintf("./cmd/%[1]s/commands/root.go", coinName)
		commandFile, err := os.Create(commandFilePath) //nolint:gosec
		if err != nil {
			log.Errorf("failed to create new coin file %s", coinFilePath)
			return err
		}
		defer commandFile.Close() //nolint
		coinTestFilePath := fmt.Sprintf("./cmd/%[1]s/%[1]s_test.go", coinName)
		coinTestFile, err := os.Create(coinTestFilePath) //nolint:gosec
		if err != nil {
			log.Errorf("failed to create new coin test file %s", coinTestFilePath)
			return err
		}
		defer coinTestFile.Close() //nolint
		paramsFilePath := "./src/params/params.go"
		paramsFile, err := os.Create(paramsFilePath)
		if err != nil {
			log.Errorf("failed to create new file %s", paramsFilePath)
			return err
		}
		defer paramsFile.Close() //nolint

		// Read template files — from --template-dir if provided, otherwise from embedded FS
		readTemplate := func(name string) ([]byte, error) {
			if templateDir != "" {
				path := filepath.Join(templateDir, name)
				data, err := os.ReadFile(path) //nolint:gosec
				if err != nil {
					return nil, fmt.Errorf("failed to read template %s from %s: %w", name, templateDir, err)
				}
				return data, nil
			}
			return templates.FS.ReadFile(name)
		}

		coinTemplateContent, err := readTemplate(coinTemplateFile)
		if err != nil {
			log.Errorf("failed to read coin template: %v", err)
			return err
		}
		commandTemplateContent, err := readTemplate(commandTemplateFile)
		if err != nil {
			log.Errorf("failed to read command template: %v", err)
			return err
		}
		coinTestTemplateContent, err := readTemplate(coinTestTemplateFile)
		if err != nil {
			log.Errorf("failed to read coin test template: %v", err)
			return err
		}
		paramsTemplateContent, err := readTemplate(paramsTemplateFile)
		if err != nil {
			log.Errorf("failed to read params template: %v", err)
			return err
		}

		// Parse templates from embedded content
		t := template.New(coinTemplateFile)
		t, err = t.Parse(string(coinTemplateContent))
		if err != nil {
			log.Errorf("failed to parse coin template: %v", err)
			return err
		}
		t, err = t.New(commandTemplateFile).Parse(string(commandTemplateContent))
		if err != nil {
			log.Errorf("failed to parse command template: %v", err)
			return err
		}
		t, err = t.New(coinTestTemplateFile).Parse(string(coinTestTemplateContent))
		if err != nil {
			log.Errorf("failed to parse coin test template: %v", err)
			return err
		}
		t, err = t.New(paramsTemplateFile).Parse(string(paramsTemplateContent))
		if err != nil {
			log.Errorf("failed to parse params template: %v", err)
			return err
		}

		config.Node.CoinName = coinName
		config.Node.CoinAscii = calvin.AsciiFont(coinName)
		config.Node.DataDirectory = "$HOME/." + coinName
		err = t.ExecuteTemplate(commandFile, coinTemplateFile, config.Node)
		if err != nil {
			log.Error("failed to parse coin template variables")
			return err
		}
		err = t.ExecuteTemplate(coinFile, commandTemplateFile, config.Node)
		if err != nil {
			log.Error("failed to parse command template variables")
			return err
		}
		if _, err := coinFile.WriteString(helpTemplate); err != nil {
			log.Error("failed to append help constant to command")
			return err
		}
		err = t.ExecuteTemplate(coinTestFile, coinTestTemplateFile, config.Node)
		if err != nil {
			log.Error("failed to parse coin test template variables")
			return err
		}
		err = t.ExecuteTemplate(paramsFile, paramsTemplateFile, config.Params)
		if err != nil {
			log.Error("failed to parse params template variables")
			return err
		}
		return nil
	},
}

func validateCoinName(s string) error {
	x := regexp.MustCompile(fmt.Sprintf(`^%s$`, useragent.NamePattern))
	if !x.MatchString(s) {
		return fmt.Errorf("invalid coin name. must only contain the characters %s", useragent.NamePattern)
	}
	return nil
}

const helpTemplate = `
const help = "{{if .HasAvailableSubCommands}}{{end}} {{if gt (len .Aliases) 0}}\r\n\r\n" +
	"{{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}" +
	"Available Commands:{{range .Commands}}  {{if and (ne .Name \"completion\") .IsAvailableCommand}}\r\n  " +
	"{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}\r\n\r\n" +
	"Flags:\r\n" +
	"{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}\r\n\r\n" +
	"Global Flags:\r\n" +
	"{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}\r\n\r\n"
`
