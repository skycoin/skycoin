package skycoin

import (
	"flag"
	"path/filepath"
	"strings"
	"time"

	"log"

	"fmt"

	"github.com/spf13/viper"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	help = false
)

// Config records fiber coin configs
type Config struct {
	Node  NodeConfig  `mapstructure:"node"`
	Build BuildConfig `mapstructure:"build"`
}

// BuildConfig records build info
type BuildConfig struct {
	Version string `mapstructure:"version"` // version number
	Commit  string `mapstructure:"commit"`  // git commit id
	Branch  string `mapstructure:"branch"`  // git branch name
}

// NodeConfig records the node's configuration
type NodeConfig struct {
	// Disable peer exchange
	DisablePEX bool `mapstructure:"disable_pex"`
	// Download peer list
	DownloadPeerList bool `mapstructure:"download_peer_list"`
	// Download the peers list from this URL
	PeerListURL string `mapstructure:"peer_list_url"`
	// Don't make any outgoing connections
	DisableOutgoingConnections bool `mapstructure:"disable_outgoing_connections"`
	// Don't allowing incoming connections
	DisableIncomingConnections bool `mapstructure:"disable_incoming_connections"`
	// Disables networking altogether
	DisableNetworking bool `mapstructure:"disable_networking"`
	// Enable wallet API
	EnableWalletAPI bool `mapstructure:"enable_wallet_api"`
	// Enable GUI
	EnableGUI bool `mapstructure:"enable_gui"`
	// Disable CSRF check in the wallet API
	DisableCSRF bool `mapstructure:"disable_csrf"`
	// Enable /api/v1/wallet/seed API endpoint
	EnableSeedAPI bool `mapstructure:"enable_seed_api"`
	// Enable unversioned API endpoints (without the /api/v1 prefix)
	EnableUnversionedAPI bool `mapstructure:"enable_unversioned_api"`

	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool `mapstructure:"localhost_only"`
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string `mapstructure:"address"`
	// gnet uses this for TCP incoming and outgoing
	Port int `mapstructure:"port"`
	// Maximum outgoing connections to maintain
	MaxOutgoingConnections int `mapstructure:"max_outgoing_connections"`
	// Maximum default outgoing connections
	MaxDefaultPeerOutgoingConnections int `mapstructure:"max_default_peer_outgoing_connections"`
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration `mapstructure:"outgoing_connections_rate"`
	// PeerlistSize represents the maximum number of peers that the pex would maintain
	PeerlistSize int `mapstructure:"peerlist_size"`
	// Wallet Address Version
	//AddressVersion string
	// Remote web interface
	WebInterface      bool   `mapstructure:"web_interface"`
	WebInterfacePort  int    `mapstructure:"web_interface_port"`
	WebInterfaceAddr  string `mapstructure:"web_interface_addr"`
	WebInterfaceCert  string `mapstructure:"web_interface_cert"`
	WebInterfaceKey   string `mapstructure:"web_interface_key"`
	WebInterfaceHTTPS bool   `mapstructure:"web_interface_https"`

	RPCInterface bool `mapstructure:"rpc_interface"`

	// Launch System Default Browser after client startup
	LaunchBrowser bool `mapstructure:"launch_browser"`

	// If true, print the configured client web interface address and exit
	PrintWebInterfaceAddress bool `mapstructure:"print_web_interface_address"`

	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string `mapstructure:"data_directory"`
	// GUI directory contains assets for the HTML interface
	GUIDirectory string `mapstructure:"gui_directory"`

	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`

	// Logging
	ColorLog bool `mapstructure:"color_log"`
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string `mapstructure:"log_level"`
	// Disable "Reply to ping", "Received pong" log messages
	DisablePingPong bool `mapstructure:"disable_ping_pong"`

	// Verify the database integrity after loading
	VerifyDB bool `mapstructure:"verify_db"`
	// Reset the database if integrity checks fail, and continue running
	ResetCorruptDB bool `mapstructure:"reset_corrupt_db"`

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string `mapstructure:"wallet_directory"`
	// Wallet crypto type
	WalletCryptoType string `mapstructure:"wallet_crypto_type"`

	RunMaster bool `mapstructure:"run_master"`

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool `mapstructure:"profile_cpu"`
	// Where the file is written to
	ProfileCPUFile string `mapstructure:"profile_cpu_file"`
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool `mapstructure:"http_prof"`

	DBPath      string `mapstructure:"db_path"`
	DBReadOnly  bool   `mapstructure:"db_read_only"`
	Arbitrating bool   `mapstructure:"arbitrating"`
	LogToFile   bool   `mapstructure:"log_to_file"`
	Version     bool   `mapstructure:"version"` // show node version

	GenesisSignatureStr string   `mapstructure:"genesis_signature_str"`
	GenesisAddressStr   string   `mapstructure:"genesis_address_str"`
	BlockchainPubkeyStr string   `mapstructure:"blockchain_pubkey_str"`
	BlockchainSeckeyStr string   `mapstructure:"blockchain_seckey_str"`
	GenesisTimestamp    uint64   `mapstructure:"genesis_timestamp"`
	GenesisCoinVolume   uint64   `mapstructure:"genesis_coin_volume"`
	DefaultConnections  []string `mapstructure:"default_connections"`

	genesisSignature cipher.Sig
	genesisTimestamp uint64
	genesisAddress   cipher.Address

	blockchainPubkey cipher.PubKey
	blockchainSeckey cipher.SecKey
}

func setDefaults() {
	// node defaults
	viper.SetDefault("node.genesis_coin_volume", 100e12)
	viper.SetDefault("node.disable_pex", false)
	viper.SetDefault("node.outgoing_connections", false)
	viper.SetDefault("node.disable_outgoing_connections", false)
	viper.SetDefault("node.disable_incoming_connections", false)
	viper.SetDefault("node.disable_networking", false)
	viper.SetDefault("node.enable_wallet_api", false)
	viper.SetDefault("node.enable_gui", false)
	viper.SetDefault("node.disable_csrf", false)
	viper.SetDefault("node.enable_seed_api", false)
	viper.SetDefault("node.enable_unversioned_api", false)
	viper.SetDefault("node.localhost_only", false)
	viper.SetDefault("node.address", "")
	viper.SetDefault("node.port", 6000)
	viper.SetDefault("node.max_outgoing_connections", 16)
	viper.SetDefault("node.max_default_peer_outgoing_connections", 1)
	viper.SetDefault("node.outgoing_connections_rate", time.Second*5)
	viper.SetDefault("node.peerlist_size", 65535)
	viper.SetDefault("node.web_interface", true)
	viper.SetDefault("node.web_interface_port", 6420)
	viper.SetDefault("node.web_interface_addr", "127.0.0.1")
	viper.SetDefault("node.web_interface_cert", "")
	viper.SetDefault("node.web_interface_key", "")
	viper.SetDefault("node.web_interface_https", false)
	viper.SetDefault("node.print_web_interface_address", false)
	viper.SetDefault("node.rpc_interface", true)
	viper.SetDefault("node.launch_browser", false)
	viper.SetDefault("node.data_directory", "$HOME/.skycoin")
	viper.SetDefault("node.gui_directory", "./src/gui/static/")
	viper.SetDefault("node.read_timeout", time.Second*10)
	viper.SetDefault("node.write_timeout", time.Second*60)
	viper.SetDefault("node.idle_timeout", time.Second*120)
	viper.SetDefault("node.color_log", true)
	viper.SetDefault("node.log_level", "INFO")
	viper.SetDefault("node.disable_ping_pong", false)
	viper.SetDefault("node.verify_db", true)
	viper.SetDefault("node.reset_corrupt_db", false)
	viper.SetDefault("node.wallet_directory", "")
	viper.SetDefault("node.wallet_crypto_type", string(wallet.CryptoTypeScryptChacha20poly1305))
	viper.SetDefault("node.run_master", false)
	viper.SetDefault("node.profile_cpu", false)
	viper.SetDefault("node.profile_cpu_file", "skycoin.prof")
	viper.SetDefault("node.http_prof", false)
	viper.SetDefault("node.db_path", "")
	viper.SetDefault("node.db_read_only", false)
	viper.SetDefault("node.arbitrating", false)
	viper.SetDefault("node.log_to_file", false)
	viper.SetDefault("node.version", false)

	// build defaults
	viper.SetDefault("build.commit", "")
	viper.SetDefault("build.branch", "")
}

// NewConfig loads blockchain config parameters from a config file
// default file is: fiber.toml in the project root
// JSON, toml or yaml file can be used (toml preferred).
func NewConfig(configName, appDir string) (Config, error) {
	// set viper parameters
	// check that file is of supported type
	confNameSplit := strings.Split(configName, ".")
	fileType := confNameSplit[len(confNameSplit)-1]
	switch fileType {
	case "toml", "json", "yaml", "yml":
		viper.SetConfigType(confNameSplit[len(confNameSplit)-1])
	default:
		return Config{}, fmt.Errorf("invalid blockchain config file type: %s", fileType)
	}

	configName = configName[:len(configName)-(len(fileType)+1)]
	viper.SetConfigName(configName)

	viper.AddConfigPath(appDir)
	viper.AddConfigPath(".")

	// set defaults
	setDefaults()

	cfg := Config{}

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	// TODO(therealssj): validate the config values
	return cfg, nil
}

func (c *Config) postProcess() {
	var err error
	if c.Node.GenesisSignatureStr != "" {
		c.Node.genesisSignature, err = cipher.SigFromHex(c.Node.GenesisSignatureStr)
		panicIfError(err, "Invalid Signature")
	}

	if c.Node.GenesisAddressStr != "" {
		c.Node.genesisAddress, err = cipher.DecodeBase58Address(c.Node.GenesisAddressStr)
		panicIfError(err, "Invalid Address")
	}
	if c.Node.BlockchainPubkeyStr != "" {
		c.Node.blockchainPubkey, err = cipher.PubKeyFromHex(c.Node.BlockchainPubkeyStr)
		panicIfError(err, "Invalid Pubkey")
	}
	if c.Node.BlockchainSeckeyStr != "" {
		c.Node.blockchainSeckey, err = cipher.SecKeyFromHex(c.Node.BlockchainSeckeyStr)
		panicIfError(err, "Invalid Seckey")
		c.Node.BlockchainSeckeyStr = ""
	}
	if c.Node.BlockchainSeckeyStr != "" {
		c.Node.blockchainSeckey = cipher.SecKey{}
	}

	home := file.UserHome()
	c.Node.DataDirectory, err = file.InitDataDir(replaceHome(c.Node.DataDirectory, home))
	panicIfError(err, "Invalid DataDirectory")

	if c.Node.WebInterfaceCert == "" {
		c.Node.WebInterfaceCert = filepath.Join(c.Node.DataDirectory, "cert.pem")
	} else {
		c.Node.WebInterfaceCert = replaceHome(c.Node.WebInterfaceCert, home)
	}

	if c.Node.WebInterfaceKey == "" {
		c.Node.WebInterfaceKey = filepath.Join(c.Node.DataDirectory, "key.pem")
	} else {
		c.Node.WebInterfaceKey = replaceHome(c.Node.WebInterfaceKey, home)
	}

	if c.Node.WalletDirectory == "" {
		c.Node.WalletDirectory = filepath.Join(c.Node.DataDirectory, "wallets")
	} else {
		c.Node.WalletDirectory = replaceHome(c.Node.WalletDirectory, home)
	}

	if c.Node.DBPath == "" {
		c.Node.DBPath = filepath.Join(c.Node.DataDirectory, "data.db")
	} else {
		c.Node.DBPath = replaceHome(c.Node.DBPath, home)
	}

	if c.Node.RunMaster {
		// Run in arbitrating mode if the node is master
		c.Node.Arbitrating = true
	}

	// Don't open browser to load wallets if wallet apis are disabled.
	if !c.Node.EnableWalletAPI {
		c.Node.EnableGUI = false
		c.Node.LaunchBrowser = false
	}

	if c.Node.EnableGUI {
		c.Node.GUIDirectory = file.ResolveResourceDirectory(c.Node.GUIDirectory)
	}
}

func (c *Config) register() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&c.Node.DisablePEX, "disable-pex", c.Node.DisablePEX, "disable PEX peer discovery")
	flag.BoolVar(&c.Node.DownloadPeerList, "download-peerlist", c.Node.DownloadPeerList, "download a peers.txt from -peerlist-url")
	flag.StringVar(&c.Node.PeerListURL, "peerlist-url", c.Node.PeerListURL, "with -download-peerlist=true, download a peers.txt file from this url")
	flag.BoolVar(&c.Node.DisableOutgoingConnections, "disable-outgoing", c.Node.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.Node.DisableIncomingConnections, "disable-incoming", c.Node.DisableIncomingConnections, "Don't make incoming connections")
	flag.BoolVar(&c.Node.DisableNetworking, "disable-networking", c.Node.DisableNetworking, "Disable all network activity")
	flag.BoolVar(&c.Node.EnableWalletAPI, "enable-wallet-api", c.Node.EnableWalletAPI, "Enable the wallet API")
	flag.BoolVar(&c.Node.EnableGUI, "enable-gui", c.Node.EnableGUI, "Enable GUI")
	flag.BoolVar(&c.Node.EnableUnversionedAPI, "enable-unversioned-api", c.Node.EnableUnversionedAPI, "Enable the deprecated unversioned API endpoints without /api/v1 prefix")
	flag.BoolVar(&c.Node.DisableCSRF, "disable-csrf", c.Node.DisableCSRF, "disable CSRF check")
	flag.BoolVar(&c.Node.EnableSeedAPI, "enable-seed-api", c.Node.EnableSeedAPI, "enable /api/v1/wallet/seed api")
	flag.StringVar(&c.Node.Address, "address", c.Node.Address, "IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Node.Port, "port", c.Node.Port, "Port to run application on")

	flag.BoolVar(&c.Node.WebInterface, "web-interface", c.Node.WebInterface, "enable the web interface")
	flag.IntVar(&c.Node.WebInterfacePort, "web-interface-port", c.Node.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.Node.WebInterfaceAddr, "web-interface-addr", c.Node.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.Node.WebInterfaceCert, "web-interface-cert", c.Node.WebInterfaceCert, "cert.pem file for web interface HTTPS. If not provided, will use cert.pem in -data-directory")
	flag.StringVar(&c.Node.WebInterfaceKey, "web-interface-key", c.Node.WebInterfaceKey, "key.pem file for web interface HTTPS. If not provided, will use key.pem in -data-directory")
	flag.BoolVar(&c.Node.WebInterfaceHTTPS, "web-interface-https", c.Node.WebInterfaceHTTPS, "enable HTTPS for web interface")

	flag.BoolVar(&c.Node.RPCInterface, "rpc-interface", c.Node.RPCInterface, "enable the rpc interface")

	flag.BoolVar(&c.Node.LaunchBrowser, "launch-browser", c.Node.LaunchBrowser, "launch system default webbrowser at client startup")
	flag.BoolVar(&c.Node.PrintWebInterfaceAddress, "print-web-interface-address", c.Node.PrintWebInterfaceAddress, "print configured web interface address and exit")
	flag.StringVar(&c.Node.DataDirectory, "data-dir", c.Node.DataDirectory, "directory to store app data (defaults to ~/.skycoin)")
	flag.StringVar(&c.Node.DBPath, "db-path", c.Node.DBPath, "path of database file (defaults to ~/.skycoin/data.db)")
	flag.BoolVar(&c.Node.DBReadOnly, "db-read-only", c.Node.DBReadOnly, "open bolt db read-only")
	flag.BoolVar(&c.Node.ProfileCPU, "profile-cpu", c.Node.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.Node.ProfileCPUFile, "profile-cpu-file", c.Node.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.Node.HTTPProf, "http-prof", c.Node.HTTPProf, "Run the http profiling interface")
	flag.StringVar(&c.Node.LogLevel, "log-level", c.Node.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.Node.ColorLog, "color-log", c.Node.ColorLog, "Add terminal colors to log output")
	flag.BoolVar(&c.Node.DisablePingPong, "no-ping-log", c.Node.DisablePingPong, `disable "reply to ping" and "received pong" debug log messages`)
	flag.BoolVar(&c.Node.LogToFile, "logtofile", c.Node.LogToFile, "log to file")
	flag.StringVar(&c.Node.GUIDirectory, "gui-dir", c.Node.GUIDirectory, "static content directory for the HTML interface")

	flag.BoolVar(&c.Node.VerifyDB, "verify-db", c.Node.VerifyDB, "check the database for corruption")
	flag.BoolVar(&c.Node.ResetCorruptDB, "reset-corrupt-db", c.Node.ResetCorruptDB, "reset the database if corrupted, and continue running instead of exiting")

	// Key Configuration Data
	flag.BoolVar(&c.Node.RunMaster, "master", c.Node.RunMaster, "run the daemon as blockchain master server")

	flag.StringVar(&c.Node.BlockchainPubkeyStr, "master-public-key", c.Node.BlockchainPubkeyStr, "public key of the master chain")
	flag.StringVar(&c.Node.BlockchainSeckeyStr, "master-secret-key", c.Node.BlockchainSeckeyStr, "secret key, set for master")

	flag.StringVar(&c.Node.GenesisAddressStr, "genesis-address", c.Node.GenesisAddressStr, "genesis address")
	flag.StringVar(&c.Node.GenesisSignatureStr, "genesis-signature", c.Node.GenesisSignatureStr, "genesis block signature")
	flag.Uint64Var(&c.Node.GenesisTimestamp, "genesis-timestamp", c.Node.GenesisTimestamp, "genesis block timestamp")

	flag.StringVar(&c.Node.WalletDirectory, "wallet-dir", c.Node.WalletDirectory, "location of the wallet files. Defaults to ~/.skycoin/wallet/")
	flag.IntVar(&c.Node.MaxOutgoingConnections, "max-outgoing-connections", c.Node.MaxOutgoingConnections, "The maximum outgoing connections allowed")
	flag.IntVar(&c.Node.MaxDefaultPeerOutgoingConnections, "max-default-peer-outgoing-connections", c.Node.MaxDefaultPeerOutgoingConnections, "The maximum default peer outgoing connections allowed")
	flag.IntVar(&c.Node.PeerlistSize, "peerlist-size", c.Node.PeerlistSize, "The peer list size")
	flag.DurationVar(&c.Node.OutgoingConnectionsRate, "connection-rate", c.Node.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.BoolVar(&c.Node.LocalhostOnly, "localhost-only", c.Node.LocalhostOnly, "Run on localhost and only connect to localhost peers")
	flag.BoolVar(&c.Node.Arbitrating, "arbitrating", c.Node.Arbitrating, "Run node in arbitrating mode")
	flag.StringVar(&c.Node.WalletCryptoType, "wallet-crypto-type", c.Node.WalletCryptoType, "wallet crypto type. Can be sha256-xor or scrypt-chacha20poly1305")
	flag.BoolVar(&c.Node.Version, "version", false, "show node version")
}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}
