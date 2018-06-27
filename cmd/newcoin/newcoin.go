package main

import (
	"fmt"

	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/skycoin"
	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	// Version is the cli version
	Version = "0.1"
)

// CoinTemplateParameters represents parameters used to generate the new coin files.
type CoinTemplateParameters struct {
	Version             string
	PeerListURL         string
	Port                int
	WebInterfacePort    int
	DataDirectory       string
	ProfileCPUFile      string
	GenesisSignatureStr string
	GenesisAddressStr   string
	BlockchainPubkeyStr string
	BlockchainSeckeyStr string
	GenesisTimestamp    uint64
	GenesisCoinVolume   uint64
	DefaultConnections  []string
}

var (
	app = cli.NewApp()
	log = logging.MustGetLogger("newcoin")
)

func init() {
	app.Name = "newcoin"
	app.Usage = "newcoin is a helper tool for creating new fiber coins"
	app.Version = Version
	commands := cli.Commands{
		createCoinCommand(),
	}

	app.Commands = commands
	app.EnableBashCompletion = true
	app.OnUsageError = func(context *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(context.App.Writer, "error: %v\n\n", err)
		cli.ShowAppHelp(context)
		return nil
	}
	app.CommandNotFound = func(context *cli.Context, command string) {
		tmp := fmt.Sprintf("{{.HelpName}}: '%s' is not a {{.HelpName}} "+
			"command. See '{{.HelpName}} --help'. \n", command)
		cli.HelpPrinter(app.Writer, tmp, app)
	}
}

func createCoinCommand() cli.Command {
	name := "createcoin"
	return cli.Command{
		Name:  name,
		Usage: "Create a new coin from a template file",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "coin",
				Usage: "name of the coin to create",
				Value: "skycoin",
			},
			cli.StringFlag{
				Name:  "template-dir, td",
				Usage: "template directory path",
				Value: "./template",
			},
			cli.StringFlag{
				Name:  "coin-template-file, ct",
				Usage: "coin template file",
				Value: "coin.template",
			},
			cli.StringFlag{
				Name:  "visor-template-file, vt",
				Usage: "visor template file",
				Value: "visor.template",
			},
			cli.StringFlag{
				Name:  "config-dir, cd",
				Usage: "config directory path",
				Value: "./",
			},
			cli.StringFlag{
				Name:  "config-file, cf",
				Usage: "config file path",
				Value: "fiber.toml",
			},
		},
		Action: func(c *cli.Context) error {
			// -- parse flags -- //

			templateDir := c.String("template-dir")

			coinTemplateFile := c.String("coin-template-file")
			visorTemplateFile := c.String("visor-template-file")

			// check that the coin template file exists
			if _, err := os.Stat(filepath.Join(templateDir, coinTemplateFile)); os.IsNotExist(err) {
				return err
			}
			// check that the visor template file exists
			if _, err := os.Stat(filepath.Join(templateDir, visorTemplateFile)); os.IsNotExist(err) {
				return err
			}

			configFile := c.String("config-file")
			configDir := c.String("config-dir")

			configFilepath := filepath.Join(configDir, configFile)
			// check that the config file exists
			if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
				return err
			}

			coinName := c.String("coin")

			// -- parse template and create new coin.go and visor parameters.go -- //

			config, err := skycoin.NewParameters(configFile, configDir)
			if err != nil {
				log.Errorf("failed to create new fiber coin config")
				return err
			}

			coinDir := fmt.Sprintf("./cmd/%s", coinName)
			// create new coin directory
			// MkdirAll does not error out if the directory already exists
			err = os.MkdirAll(coinDir, 0755)
			if err != nil {
				log.Errorf("failed to create new coin directory %s", coinDir)
				return err
			}

			// we have to always create a new file otherwise the templating gives an error
			coinFilePath := fmt.Sprintf("./cmd/%[1]s/%[1]s.go", coinName)
			coinFile, err := os.Create(coinFilePath)
			if err != nil {
				log.Errorf("failed to create new coin file %s", coinFilePath)
				return err
			}

			visorParamsFile, err := os.Create("./src/visor/parameters.go")
			if err != nil {
				log.Errorf("failed to create new visor parameters.go")
				return err
			}

			// change dir so that text/template can parse the file
			err = os.Chdir(templateDir)
			if err != nil {
				log.Errorf("failed to change directory to %s", templateDir)
				return err
			}

			t := template.New(coinTemplateFile)
			t, err = t.ParseFiles(coinTemplateFile, visorTemplateFile)
			if err != nil {
				log.Errorf("failed to parse template file: %s", coinTemplateFile)
				return err
			}

			err = t.ExecuteTemplate(coinFile, coinTemplateFile, CoinTemplateParameters{
				Version:             config.Build.Version,
				PeerListURL:         config.Node.PeerListURL,
				Port:                config.Node.Port,
				WebInterfacePort:    config.Node.WebInterfacePort,
				DataDirectory:       "$HOME/." + coinName,
				ProfileCPUFile:      coinName + ".prof",
				GenesisSignatureStr: config.Node.GenesisSignatureStr,
				GenesisAddressStr:   config.Node.GenesisAddressStr,
				BlockchainPubkeyStr: config.Node.BlockchainPubkeyStr,
				BlockchainSeckeyStr: config.Node.BlockchainSeckeyStr,
				GenesisTimestamp:    config.Node.GenesisTimestamp,
				GenesisCoinVolume:   config.Node.GenesisCoinVolume,
				DefaultConnections:  config.Node.DefaultConnections,
			})
			if err != nil {
				log.Errorf("failed to parse coin template variables")
				return err
			}

			err = t.ExecuteTemplate(visorParamsFile, visorTemplateFile, config.Visor)
			if err != nil {
				log.Errorf("failed to parse visor params template variables")
				return err
			}

			return nil
		},
	}
}

func main() {
	if e := app.Run(os.Args); e != nil {
		log.Fatal(e)
	}
}
