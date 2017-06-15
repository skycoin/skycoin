package cli

import (
	"bytes"
	"errors"
	"fmt"

	"encoding/json"

	"os"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/util"
	gcli "github.com/urfave/cli"
)

// Commands all cmds that we support

var (
	walletExt = ".wlt"
	cfg       Config
)

var (
	errConnectNodeFailed = errors.New("connect to node failed")
	errWalletName        = fmt.Errorf("error wallet file name, must has %v extension", walletExt)
	errAddress           = errors.New("invalidate address")
	errReadResponse      = errors.New("read response body failed")
	errJSONMarshal       = errors.New("json marshal failed")
	errJSONUnmarshal     = errors.New("json unmarshal failed")
)

var (
	commandHelpTemplate = `USAGE:
		{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}
		
CATEGORY:
		{{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
		{{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
		{{range .VisibleFlags}}{{.}}
		{{end}}{{end}}
	`
)

func stringPtr(v string) *string {
	return &v
}

func httpGet(url string, v interface{}) error {
	return nil
}

func init() {
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
	cfg Config
}

// Config cli's configuration struct
type Config struct {
	RPCAddress        string
	WalletDir         string
	DefaultWalletName string
	Coin              string
}

// Option Init argument type
type Option func(app *App)

// NewApp creates an app instance
func NewApp(ops ...Option) *App {
	home := util.UserHome()
	app := &App{
		App: *gcli.NewApp(),
		cfg: Config{
			RPCAddress:        "127.0.0.1:6430",
			WalletDir:         home + "/." + os.Args[0] + "/wallets",
			DefaultWalletName: fmt.Sprintf("%s_cli.wlt", os.Args[0]),
			Coin:              "skycoin",
		},
	}

	for _, op := range ops {
		op(app)
	}

	// init the global rpcAddr variable
	cfg = app.cfg

	commands := []gcli.Command{
		addPrivateKeyCMD(),
		blocksCMD(),
		broadcastTxCMD(),
		walletBalanceCMD(),
		walletOutputsCMD(),
		addressBalanceCMD(),
		addressOutputsCMD(),
		createRawTxCMD(),
		generateAddrsCMD(),
		generateWalletCMD(),
		lastBlocksCMD(),
		listAddressesCMD(),
		listWalletsCMD(),
		sendCMD(),
		statusCMD(),
		transactionCMD(),
		versionCMD(),
		walletDirCMD(),
		walletHisCMD(),
	}

	app.Usage = fmt.Sprintf("the %s command line interface", app.cfg.Coin)
	app.Version = "0.1"
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

	return app
}

// Run starts the app
func (app *App) Run(args []string) error {
	return app.App.Run(args)
}

// RPCAddr sets rpc address
func RPCAddr(addr string) Option {
	return func(app *App) {
		app.cfg.RPCAddress = addr
	}
}

// WalletDir sets wallet dir
func WalletDir(wltDir string) Option {
	return func(app *App) {
		app.cfg.WalletDir = wltDir
	}
}

// DefaultWltName sets default wallet name
func DefaultWltName(wltName string) Option {
	return func(app *App) {
		app.cfg.DefaultWalletName = wltName
	}
}

// Coin sets the coin name
func Coin(coin string) Option {
	return func(app *App) {
		app.cfg.Coin = coin
	}
}

func getUnspent(addrs []string) (unspentOutSet, error) {
	req, err := webrpc.NewRequest("get_outputs", addrs, "1")
	if err != nil {
		return unspentOutSet{}, fmt.Errorf("create webrpc request failed:%v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return unspentOutSet{}, fmt.Errorf("do rpc request failed:%v", err)
	}

	if rsp.Error != nil {
		return unspentOutSet{}, fmt.Errorf("rpc request failed, %+v", *rsp.Error)
	}

	var rlt webrpc.OutputsResult
	if err := json.NewDecoder(bytes.NewBuffer(rsp.Result)).Decode(&rlt); err != nil {
		return unspentOutSet{}, errJSONUnmarshal
	}

	return unspentOutSet{rlt.Outputs}, nil
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
