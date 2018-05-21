/*
Package cli implements an interface for creating a CLI application.
Includes methods for manipulating wallets files and interacting with the
webrpc API to query a skycoin node's status.
*/
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"os"

	gcli "github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/util/file"
)

const (
	// Version is the CLI Version
	Version           = "0.23.1-rc2"
	walletExt         = ".wlt"
	defaultCoin       = "skycoin"
	defaultWalletName = "$COIN_cli" + walletExt
	defaultWalletDir  = "$DATA_DIR/wallets"
	defaultRPCAddress = "http://127.0.0.1:6420"
	defaultDataDir    = "$HOME/.$COIN/"
)

var (
	envVarsHelp = fmt.Sprintf(`ENVIRONMENT VARIABLES:
    RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "%s"
    COIN: Name of the coin. Default "%s"
    USE_CSRF: Set to 1 or true if the remote node has CSRF enabled. Default false (unset)
    WALLET_DIR: Directory where wallets are stored. This value is overriden by any subcommand flag specifying a wallet filename, if that filename includes a path. Default "%s"
    WALLET_NAME: Name of wallet file (without path). This value is overriden by any subcommand flag specifying a wallet filename. Default "%s"
    DATA_DIR: Directory where everything is stored. Default "%s"`, defaultRPCAddress, defaultCoin, defaultWalletDir, defaultWalletName, defaultDataDir)

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

	// ErrWalletName is returned if the wallet file name is invalid
	ErrWalletName = fmt.Errorf("error wallet file name, must have %s extension", walletExt)
	// ErrAddress is returned if an address is invalid
	ErrAddress = errors.New("invalid address")
	// ErrJSONMarshal is returned if JSON marshaling failed
	ErrJSONMarshal = errors.New("json marshal failed")
)

// App Wraps the app so that main package won't use the raw App directly,
// which will cause import issue
type App struct {
	gcli.App
}

// Config cli's configuration struct
type Config struct {
	WalletDir  string `json:"wallet_directory"`
	WalletName string `json:"wallet_name"`
	DataDir    string `json:"data_directory"`
	Coin       string `json:"coin"`
	RPCAddress string `json:"rpc_address"`
	UseCSRF    bool   `json:"use_csrf"`
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

	if _, err := url.Parse(rpcAddr); err != nil {
		return Config{}, errors.New("RPC_ADDR must be in scheme://host format")
	}

	home := file.UserHome()

	// get data dir dir from env
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = filepath.Join(home, fmt.Sprintf(".%s", coin))
	}

	// get wallet dir from env
	wltDir := os.Getenv("WALLET_DIR")
	if wltDir == "" {
		wltDir = filepath.Join(dataDir, "wallets")
	}

	// get wallet name from env
	wltName := os.Getenv("WALLET_NAME")
	if wltName == "" {
		wltName = fmt.Sprintf("%s_cli%s", coin, walletExt)
	}

	if !strings.HasSuffix(wltName, walletExt) {
		return Config{}, ErrWalletName
	}
	var useCSRF bool
	useCSRFStr := os.Getenv("USE_CSRF")
	if useCSRFStr != "" {
		var err error
		useCSRF, err = strconv.ParseBool(useCSRFStr)
		if err != nil {
			return Config{}, errors.New("Invalid USE_CSRF value, must be interpretable as a boolean e.g. 0, 1, true, false")
		}
	}

	return Config{
		WalletDir:  wltDir,
		WalletName: wltName,
		DataDir:    dataDir,
		Coin:       coin,
		RPCAddress: rpcAddr,
		UseCSRF:    useCSRF,
	}, nil
}

// FullWalletPath returns the joined wallet dir and wallet name path
func (c Config) FullWalletPath() string {
	return filepath.Join(c.WalletDir, c.WalletName)
}

// FullDBPath returns the joined data directory and db file name path
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
func NewApp(cfg Config) (*App, error) {
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
		showConfigCmd(),
		statusCmd(),
		transactionCmd(),
		verifyAddressCmd(),
		versionCmd(),
		walletBalanceCmd(cfg),
		walletDirCmd(),
		walletHisCmd(),
		walletOutputsCmd(cfg),
		encryptWalletCmd(cfg),
		decryptWalletCmd(cfg),
		showSeedCmd(cfg),
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
		gcli.OsExiter(1)
	}

	rpcClient, err := webrpc.NewClient(cfg.RPCAddress)
	if err != nil {
		return nil, err
	}
	rpcClient.UseCSRF = cfg.UseCSRF

	app.Metadata = map[string]interface{}{
		"config":   cfg,
		"rpc":      rpcClient,
		"quitChan": make(chan struct{}),
	}

	return app, nil
}

// Run starts the app
func (app *App) Run(args []string) error {
	return app.App.Run(args)
}

// RPCClientFromContext returns a webrpc.Client from a urfave/cli Context
func RPCClientFromContext(c *gcli.Context) *webrpc.Client {
	return c.App.Metadata["rpc"].(*webrpc.Client)
}

// ConfigFromContext returns a Config from a urfave/cli Context
func ConfigFromContext(c *gcli.Context) Config {
	return c.App.Metadata["config"].(Config)
}

// QuitChanFromContext returns a chan struct{} from a urfave/cli Context
func QuitChanFromContext(c *gcli.Context) chan struct{} {
	return c.App.Metadata["quitChan"].(chan struct{})
}

func onCommandUsageError(command string) gcli.OnUsageErrorFunc {
	return func(c *gcli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "Error: %v\n\n", err)
		gcli.ShowCommandHelp(c, command)
		return nil
	}
}

func errorWithHelp(c *gcli.Context, err error) {
	fmt.Fprintf(c.App.Writer, "Error: %v. See '%s %s --help'\n\n", err, c.App.HelpName, c.Command.Name)
}

func formatJSON(obj interface{}) ([]byte, error) {
	d, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return nil, ErrJSONMarshal
	}
	return d, nil
}

func printJSON(obj interface{}) error {
	d, err := formatJSON(obj)
	if err != nil {
		return err
	}

	fmt.Println(string(d))

	return nil
}

// readPasswordFromTerminal promotes user to enter password and read it.
func readPasswordFromTerminal() ([]byte, error) {
	// Promotes to enter the wallet password
	fmt.Fprint(os.Stdout, "enter password:")
	bp, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, "")
	return bp, nil
}

// PUBLIC

// WalletLoadError is returned if a wallet could not be loaded
type WalletLoadError struct {
	error
}

// WalletSaveError is returned if a wallet could not be saved
type WalletSaveError struct {
	error
}

// PasswordReader is an interface for getting password
type PasswordReader interface {
	Password() ([]byte, error)
}

// PasswordFromBytes represents an implementation of PasswordReader,
// which reads password from the bytes itself.
type PasswordFromBytes []byte

// Password implements the PasswordReader's Password method
func (p PasswordFromBytes) Password() ([]byte, error) {
	return []byte(p), nil
}

// PasswordFromTerm reads password from terminal
type PasswordFromTerm struct{}

// Password implements the PasswordReader's Password method
func (p PasswordFromTerm) Password() ([]byte, error) {
	v, err := readPasswordFromTerminal()
	if err != nil {
		return nil, err
	}

	return v, nil
}

// NewPasswordReader creats a PasswordReader instance,
// reads password from the input bytes first, if it's empty, then read from terminal.
func NewPasswordReader(p []byte) PasswordReader {
	if len(p) != 0 {
		return PasswordFromBytes(p)
	}

	return PasswordFromTerm{}
}
