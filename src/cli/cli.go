/*
Package cli implements the CLI cmd's methods.

Includes methods for manipulating wallets files and interacting with the
REST API to query a skycoin node's status.
*/
package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"path/filepath"
	"syscall"

	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	// Version is the CLI Version
	Version = "0.26.0"
)

const (
	walletExt         = "." + wallet.WalletExt
	defaultCoin       = "skycoin"
	defaultRPCAddress = "http://127.0.0.1:6420"
	defaultDataDir    = "$HOME/.$COIN/"
)

var (
	envVarsHelp = fmt.Sprintf(`ENVIRONMENT VARIABLES:
    RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "%s"
    RPC_USER: Username for RPC API, if enabled in the RPC.
    RPC_PASS: Password for RPC API, if enabled in the RPC.
    COIN: Name of the coin. Default "%s"
    DATA_DIR: Directory where everything is stored. Default "%s"`, defaultRPCAddress, defaultCoin, defaultDataDir)

	helpTemplate = fmt.Sprintf(`USAGE:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command] [flags] [arguments...]{{end}}{{with (or .Long .Short)}}

DESCRIPTION:
    {{. | trimTrailingWhitespaces}}{{end}}{{if .HasExample}}

EXAMPLES:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

COMMANDS:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

FLAGS:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

GLOBAL FLAGS:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

%s
`, envVarsHelp)

	// ErrWalletName is returned if the wallet file name is invalid
	ErrWalletName = fmt.Errorf("error wallet file name, must have %s extension", walletExt)
	// ErrAddress is returned if an address is invalid
	ErrAddress = errors.New("invalid address")
	// ErrJSONMarshal is returned if JSON marshaling failed
	ErrJSONMarshal = errors.New("json marshal failed")
)

var (
	cliConfig Config
	apiClient *api.Client
	quitChan  = make(chan struct{})
)

// Config cli's configuration struct
type Config struct {
	DataDir     string `json:"data_directory"`
	Coin        string `json:"coin"`
	RPCAddress  string `json:"rpc_address"`
	RPCUsername string `json:"-"`
	RPCPassword string `json:"-"`
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

	rpcUser := os.Getenv("RPC_USER")
	rpcPass := os.Getenv("RPC_PASS")

	home := file.UserHome()

	// get data dir dir from env
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = filepath.Join(home, fmt.Sprintf(".%s", coin))
	}

	if os.Getenv("WALLET_DIR") != "" {
		return Config{}, errors.New("the envvar WALLET_DIR is no longer recognized by the CLI tool. Please review the updated CLI docs to learn how to specify the wallet file for your desired action")
	}
	if os.Getenv("WALLET_NAME") != "" {
		return Config{}, errors.New("the envvar WALLET_NAME is no longer recognized by the CLI tool. Please review the updated CLI docs to learn how to specify the wallet file for your desired action")
	}

	return Config{
		DataDir:     dataDir,
		Coin:        coin,
		RPCAddress:  rpcAddr,
		RPCUsername: rpcUser,
		RPCPassword: rpcPass,
	}, nil
}

// FullDBPath returns the joined data directory and db file name path
func (c Config) FullDBPath() string {
	return filepath.Join(c.DataDir, "data.db")
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

// NewCLI creates a cli instance
func NewCLI(cfg Config) (*cobra.Command, error) {
	apiClient = api.NewClient(cfg.RPCAddress)
	apiClient.SetAuth(cfg.RPCUsername, cfg.RPCPassword)

	cliConfig = cfg

	skyCLI := &cobra.Command{
		Short: fmt.Sprintf("The %s command line interface", cfg.Coin),
		Use:   fmt.Sprintf("%s-cli", cfg.Coin),
	}

	commands := []*cobra.Command{
		addPrivateKeyCmd(),
		addressBalanceCmd(),
		addressGenCmd(),
		fiberAddressGenCmd(),
		addressOutputsCmd(),
		blocksCmd(),
		broadcastTxCmd(),
		checkDBCmd(),
		checkDBEncodingCmd(),
		createRawTxnCmd(),
		decodeRawTxnCmd(),
		encodeJSONTxnCmd(),
		decryptWalletCmd(),
		encryptWalletCmd(),
		lastBlocksCmd(),
		listAddressesCmd(),
		listWalletsCmd(),
		sendCmd(),
		showConfigCmd(),
		showSeedCmd(),
		statusCmd(),
		transactionCmd(),
		verifyTransactionCmd(),
		verifyAddressCmd(),
		versionCmd(),
		walletCreateCmd(),
		walletAddAddressesCmd(),
		walletKeyExportCmd(),
		walletBalanceCmd(),
		walletHisCmd(),
		walletOutputsCmd(),
		richlistCmd(),
		addressTransactionsCmd(),
		pendingTransactionsCmd(),
		addresscountCmd(),
		distributeGenesisCmd(),
	}

	skyCLI.Version = Version
	skyCLI.SuggestionsMinimumDistance = 1
	skyCLI.AddCommand(commands...)

	skyCLI.SetHelpTemplate(helpTemplate)
	skyCLI.SetUsageTemplate(helpTemplate)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	return skyCLI, nil
}

func printHelp(c *cobra.Command) {
	c.Printf("See '%s %s --help'\n", c.Parent().Name(), c.Name())
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
	bp, err := terminal.ReadPassword(int(syscall.Stdin)) //nolint:unconvert
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

func (e WalletLoadError) Error() string {
	return fmt.Sprintf("Load wallet failed: %v", e.error)
}

// WalletSaveError is returned if a wallet could not be saved
type WalletSaveError struct {
	error
}

func (e WalletSaveError) Error() string {
	return fmt.Sprintf("Save wallet failed: %v", e.error)
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
