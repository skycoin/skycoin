package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"os"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/util/file"
	gcli "github.com/urfave/cli"
)

// Commands all cmds that we support

const (
	Version           = "0.19.0"
	walletExt         = ".wlt"
	defaultCoin       = "skycoin"
	defaultWalletName = "skycoin_cli.wlt"
	defaultRpcAddress = "127.0.0.1:6430"
)

var (
	envVarsHelp = fmt.Sprintf(`ENVIRONMENT VARIABLES:
    RPC_ADDR: Address of RPC node. Default %s
    WALLET_DIR: Directory where wallets are stored. Default $HOME/.%s/wallets
    WALLET_NAME: Name of wallet file (without path). Default %s`, defaultRpcAddress, defaultCoin, defaultWalletName)

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

// CLI cmds should call this from there init() method
func Init() {
	gcli.AppHelpTemplate = appHelpTemplate
	gcli.SubcommandHelpTemplate = commandHelpTemplate
	gcli.CommandHelpTemplate = commandHelpTemplate
	gcli.HelpFlag = gcli.BoolFlag{
		Name:  "help,h",
		Usage: "show help, can also be used to show subcommand help",
	}
}

// App Wraps the app so that main package won't use the raw App directly,
// which will cause import issue
type App struct {
	gcli.App
}

// Config cli's configuration struct
type Config struct {
	WalletDir  string
	WalletName string
	Coin       string
	RpcAddress string
}

func LoadConfig() (Config, error) {
	// get coin name from env
	coin := os.Getenv("COIN")
	if coin == "" {
		coin = defaultCoin
	}

	// get rpc address from env
	rpcAddr := os.Getenv("RPC_ADDR")
	if rpcAddr == "" {
		rpcAddr = defaultRpcAddress
	}

	// get wallet dir from env
	wltDir := os.Getenv("WALLET_DIR")
	if wltDir == "" {
		home := file.UserHome()
		wltDir = fmt.Sprintf("%s/.%s/wallets", home, coin)
	}

	// get wallet name from env
	wltName := os.Getenv("WALLET_NAME")
	if wltName == "" {
		wltName = defaultWalletName
	}

	if !strings.HasSuffix(wltName, walletExt) {
		return Config{}, ErrWalletName
	}

	return Config{
		WalletDir:  wltDir,
		WalletName: wltName,
		Coin:       coin,
		RpcAddress: rpcAddr,
	}, nil
}

func (c Config) FullWalletPath() string {
	return filepath.Join(c.WalletDir, c.WalletName)
}

// Returns a full wallet path based on cfg and optional cli arg specifying wallet file
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

// NewApp creates an app instance
func NewApp(cfg Config) *App {
	gcliApp := gcli.NewApp()
	app := &App{
		App: *gcliApp,
	}

	commands := []gcli.Command{
		addPrivateKeyCmd(cfg),
		blocksCmd(),
		broadcastTxCmd(),
		walletBalanceCmd(cfg),
		walletOutputsCmd(cfg),
		addressBalanceCmd(),
		addressOutputsCmd(),
		createRawTxCmd(cfg),
		generateAddrsCmd(cfg),
		generateWalletCmd(cfg),
		lastBlocksCmd(),
		listAddressesCmd(),
		listWalletsCmd(),
		sendCmd(),
		statusCmd(),
		transactionCmd(),
		versionCmd(),
		walletDirCmd(),
		walletHisCmd(),
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

	gcliApp.Metadata = map[string]interface{}{
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
	return ConfigFromContext(c)
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
