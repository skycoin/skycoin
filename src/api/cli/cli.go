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
	commands  []gcli.Command
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

func stringPtr(v string) *string {
	return &v
}

func httpGet(url string, v interface{}) error {
	return nil
}

// Config cli's configuration struct
type Config struct {
	RPCAddress        string
	WalletDir         string
	DefaultWalletName string
}

// Option Init argument type
type Option func(cfg *Config)

// Init initialize the cli's configuration
func Init(ops ...Option) {
	for _, op := range ops {
		op(&cfg)
	}

	if cfg.RPCAddress == "" {
		cfg.RPCAddress = "127.0.0.1:6422"
	}

	if cfg.WalletDir == "" {
		home := util.UserHome()
		cfg.WalletDir = home + "/." + os.Args[0] + "/wallets"
	}

	if cfg.DefaultWalletName == "" {
		cfg.DefaultWalletName = fmt.Sprintf("%s_cli.wlt", os.Args[0])
	}

	commands = append(commands,
		addPrivateKeyCMD(),
		blocksCMD(),
		broadcastTxCMD(),
		checkBalanceCMD(),
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
		walletHisCMD())
}

// RPCAddr sets rpc address
func RPCAddr(addr string) Option {
	return func(cfg *Config) {
		cfg.RPCAddress = addr
	}
}

// WalletDir sets wallet dir
func WalletDir(wltDir string) Option {
	return func(cfg *Config) {
		cfg.WalletDir = wltDir
	}
}

// DefaultWltName sets default wallet name
func DefaultWltName(wltName string) Option {
	return func(cfg *Config) {
		cfg.DefaultWalletName = wltName
	}
}

func Commands() []gcli.Command {
	return commands
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
