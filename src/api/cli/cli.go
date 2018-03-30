/*
Implements an interface for creating a CLI application.
Includes methods for manipulating wallets files and interacting with the
webrpc API to query a skycoin node's status.
*/
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"os"

	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/util/file"
)

const (
	// Version is the CLI Version
	Version           = "0.22.0"
	walletExt         = ".wlt"
	defaultCoin       = "skycoin"
	defaultWalletName = "$COIN_cli" + walletExt
	defaultWalletDir  = "$HOME/.$COIN/wallets"
	defaultRPCAddress = "127.0.0.1:6430"
)

var (
	envVarsHelp = fmt.Sprintf(`ENVIRONMENT VARIABLES:
    RPC_ADDR: Address of RPC node. Default "%s"
    COIN: Name of the coin. Default "%s"
    WALLET_DIR: Directory where wallets are stored. This value is overriden by any subcommand flag specifying a wallet filename, if that filename includes a path. Default "%s"
    WALLET_NAME: Name of wallet file (without path). This value is overriden by any subcommand flag specifying a wallet filename. Default "%s"`, defaultRPCAddress, defaultCoin, defaultWalletDir, defaultWalletName)

	commandHelpTemplate = fmt.Sprintf(`USAGE:
        {{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}

CATEGORY:
        {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
        {{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
        {{range .VisibleFlags}}{{.}}
        {{end}}{{end}}
%s
`, envVarsHelp)

	appHelpTemplate = fmt.Sprintf(`NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
%s
`, envVarsHelp)

	ErrWalletName  = fmt.Errorf("error wallet file name, must have %s extension", walletExt)
	ErrAddress     = errors.New("invalid address")
	ErrJSONMarshal = errors.New("json marshal failed")
)

// App Wraps the app so that main package won't use the raw App directly,
// which will cause import issue
type App struct {
	gcli.App
}

// Config cli's configuration struct
type Config struct {
	WalletDir  string
	WalletName string
	DataDir    string
	Coin       string
	RpcAddress string
}

// LoadConfig loads config from environment, prior to parsing CLI flags
func LoadConfig() (Config, error) {
	// get coin name from env
	coin := os.Getenv("COIN")
	if coin == "" {
		coin = defaultCoin
	}

	// get rpc address from env
	rpcAddr := os.Getenv("RPC_ADDR")
	if rpcAddr == "" {
		rpcAddr = defaultRPCAddress
	}

	home := file.UserHome()

	// get wallet dir from env
	wltDir := os.Getenv("WALLET_DIR")
	if wltDir == "" {
		wltDir = fmt.Sprintf("%s/.%s/wallets", home, coin)
	}

	// get wallet name from env
	wltName := os.Getenv("WALLET_NAME")
	if wltName == "" {
		wltName = fmt.Sprintf("%s_cli%s", coin, walletExt)
	}

	if !strings.HasSuffix(wltName, walletExt) {
		return Config{}, ErrWalletName
	}

	dataDir := filepath.Join(home, fmt.Sprintf(".%s", coin))

	return Config{
		WalletDir:  wltDir,
		WalletName: wltName,
		DataDir:    dataDir,
		Coin:       coin,
		RpcAddress: rpcAddr,
	}, nil
}

func (c Config) FullWalletPath() string {
	return filepath.Join(c.WalletDir, c.WalletName)
}

func (c Config) FullDBPath() string {
	return filepath.Join(c.DataDir, "data.db")
}

// Returns a full wallet path based on cfg and optional cli arg specifying wallet file
// FIXME: A CLI flag for the wallet filename is redundant with the envvar. Remove the flags or the envvar.
func resolveWalletPath(cfg Config, w string) (string, error) {
	if w == "" {
		w = cfg.FullWalletPath()
	}

	if !strings.HasSuffix(w, walletExt) {
		return "", ErrWalletName
	}

	// If w is only the basename, use the default wallet directory
	if filepath.Base(w) == w {
		w = filepath.Join(cfg.WalletDir, w)
	}

	absW, err := filepath.Abs(w)
	if err != nil {
		return "", fmt.Errorf("Invalid wallet path %s: %v", w, err)
	}

	return absW, nil
}

func resolveDBPath(cfg Config, db string) (string, error) {
	if db == "" {
		db = cfg.FullDBPath()
	}

	// If db is only the basename, use the default data dir
	if filepath.Base(db) == db {
		db = filepath.Join(cfg.DataDir, db)
	}

	absDB, err := filepath.Abs(db)
	if err != nil {
		return "", fmt.Errorf("Invalid data path %s: %v", db, err)
	}
	return absDB, nil
}

// NewApp creates an app instance
func NewApp(cfg Config) *App {
	gcli.AppHelpTemplate = appHelpTemplate
	gcli.SubcommandHelpTemplate = commandHelpTemplate
	gcli.CommandHelpTemplate = commandHelpTemplate

	gcliApp := gcli.NewApp()
	app := &App{
		App: *gcliApp,
	}

	commands := []gcli.Command{
		addPrivateKeyCmd(cfg),
		addressBalanceCmd(),
		addressGenCmd(),
		addressOutputsCmd(),
		blocksCmd(),
		broadcastTxCmd(),
		checkdbCmd(),
		createRawTxCmd(cfg),
		decodeRawTxCmd(),
		generateAddrsCmd(cfg),
		generateWalletCmd(cfg),
		lastBlocksCmd(),
		listAddressesCmd(),
		listWalletsCmd(),
		sendCmd(),
		statusCmd(),
		transactionCmd(),
		verifyAddressCmd(),
		versionCmd(),
		walletBalanceCmd(cfg),
		walletDirCmd(),
		walletHisCmd(),
		walletOutputsCmd(cfg),
	}

	app.Name = fmt.Sprintf("%s-cli", cfg.Coin)
	app.Version = Version
	app.Usage = fmt.Sprintf("the %s command line interface", cfg.Coin)
	app.Commands = commands
	app.EnableBashCompletion = true
	app.OnUsageError = func(context *gcli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(context.App.Writer, "Error: %v\n\n", err)
		gcli.ShowAppHelp(context)
		return nil
	}
	app.CommandNotFound = func(ctx *gcli.Context, command string) {
		tmp := fmt.Sprintf("{{.HelpName}}: '%s' is not a {{.HelpName}} command. See '{{.HelpName}} --help'.\n", command)
		gcli.HelpPrinter(app.Writer, tmp, app)
	}

	app.Metadata = map[string]interface{}{
		"config": cfg,
		"rpc": &webrpc.Client{
			Addr: cfg.RpcAddress,
		},
	}

	return app
}

// Run starts the app
func (app *App) Run(args []string) error {
	return app.App.Run(args)
}

func RpcClientFromContext(c *gcli.Context) *webrpc.Client {
	return c.App.Metadata["rpc"].(*webrpc.Client)
}

func ConfigFromContext(c *gcli.Context) Config {
	return c.App.Metadata["config"].(Config)
}

func onCommandUsageError(command string) gcli.OnUsageErrorFunc {
	return func(c *gcli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "Error: %v\n\n", err)
		gcli.ShowCommandHelp(c, command)
		return nil
	}
}

func errorWithHelp(c *gcli.Context, err error) {
	fmt.Fprintf(c.App.Writer, "ERROR: %v. See '%s %s --help'\n\n", err, c.App.HelpName, c.Command.Name)
}

func formatJson(obj interface{}) ([]byte, error) {
	d, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return nil, ErrJSONMarshal
	}
	return d, nil
}

func printJson(obj interface{}) error {
	d, err := formatJson(obj)
	if err != nil {
		return err
	}

	fmt.Println(string(d))

	return nil
}
