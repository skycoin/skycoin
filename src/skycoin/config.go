package skycoin

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"log"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/collections"
	"github.com/skycoin/skycoin/src/util/file"
	flagutil "github.com/skycoin/skycoin/src/util/flag"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	help = false
)

// Config records skycoin node and build config
type Config struct {
	Node  NodeConfig
	Build readable.BuildInfo
}

// NodeConfig records the node's configuration
type NodeConfig struct {
	// Disable peer exchange
	DisablePEX bool
	// Download peer list
	DownloadPeerList bool
	// Download the peers list from this URL
	PeerListURL string
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// DEPRECATED Enable wallet API
	EnableWalletAPI bool
	// Enable GUI
	EnableGUI bool
	// Disable CSRF check in the wallet API
	DisableCSRF bool
	// DEPRECATED Enable /api/v1/wallet/seed API endpoint
	EnableSeedAPI bool
	// Enable unversioned API endpoints (without the /api/v1 prefix)
	EnableUnversionedAPI bool
	// Disable CSP disable content-security-policy in http response
	DisableCSP bool
	// Enable all API endpoints
	EnableAllAPISets bool
	// Disable all API endpoints
	DisableAllAPISets bool
	// List of API sets enabled to be exported via remote web interface
	EnabledAPISets collections.StringSet
	// List of API sets enabled to be exported via remote web interface
	DisabledAPISets collections.StringSet

	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	// gnet uses this for TCP incoming and outgoing
	Port int
	// Maximum outgoing connections to maintain
	MaxOutgoingConnections int
	// Maximum default outgoing connections
	MaxDefaultPeerOutgoingConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
	// PeerlistSize represents the maximum number of peers that the pex would maintain
	PeerlistSize int
	// Wallet Address Version
	//AddressVersion string
	// Remote web interface
	WebInterface bool
	// Remote web interface port
	WebInterfacePort int
	// Remote web interface address
	WebInterfaceAddr string
	// Remote web interface certificate
	WebInterfaceCert string
	// Remote web interface key
	WebInterfaceKey string
	// Remote web interface HTTPS support
	WebInterfaceHTTPS bool

	// Enable the deprecated JSON 2.0 RPC interface
	RPCInterface bool

	// Launch System Default Browser after client startup
	LaunchBrowser bool

	// If true, print the configured client web interface address and exit
	PrintWebInterfaceAddress bool

	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string
	// GUI directory contains assets for the HTML interface
	GUIDirectory string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Disable "Reply to ping", "Received pong" log messages
	DisablePingPong bool

	// Verify the database integrity after loading
	VerifyDB bool
	// Reset the database if integrity checks fail, and continue running
	ResetCorruptDB bool

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string
	// Wallet crypto type
	WalletCryptoType string

	// Disable the hardcoded default peers
	DisableDefaultPeers bool
	// Load custom peers from disk
	CustomPeersFile string

	RunMaster bool

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// Enable HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Expose HTTP profiling on this interface
	HTTPProfHost string

	DBPath      string
	DBReadOnly  bool
	Arbitrating bool
	LogToFile   bool
	Version     bool // show node version

	GenesisSignatureStr string
	GenesisAddressStr   string
	BlockchainPubkeyStr string
	BlockchainSeckeyStr string
	GenesisTimestamp    uint64
	GenesisCoinVolume   uint64
	DefaultConnections  []string

	genesisSignature cipher.Sig
	genesisAddress   cipher.Address

	blockchainPubkey cipher.PubKey
	blockchainSeckey cipher.SecKey
}

// NewNodeConfig returns a new node config instance
func NewNodeConfig(mode string, node NodeParameters) NodeConfig {
	nodeConfig := NodeConfig{
		GenesisSignatureStr: node.GenesisSignatureStr,
		GenesisAddressStr:   node.GenesisAddressStr,
		GenesisCoinVolume:   node.GenesisCoinVolume,
		GenesisTimestamp:    node.GenesisTimestamp,
		BlockchainPubkeyStr: node.BlockchainPubkeyStr,
		BlockchainSeckeyStr: node.BlockchainSeckeyStr,
		DefaultConnections:  node.DefaultConnections,
		// Disable peer exchange
		DisablePEX: false,
		// Don't make any outgoing connections
		DisableOutgoingConnections: false,
		// Don't allowing incoming connections
		DisableIncomingConnections: false,
		// Disables networking altogether
		DisableNetworking: false,
		// Enable wallet API
		EnableWalletAPI: false,
		// Enable GUI
		EnableGUI: false,
		// Enable unversioned API
		EnableUnversionedAPI: false,
		// Enable seed API
		EnableSeedAPI: false,
		// Disable CSRF check in the wallet API
		DisableCSRF: false,
		// DisableCSP disable content-security-policy in http response
		DisableCSP: false,
		// Only run on localhost and only connect to others on localhost
		LocalhostOnly: false,
		// Which address to serve on. Leave blank to automatically assign to a
		// public interface
		Address: "",
		//gnet uses this for TCP incoming and outgoing
		Port: node.Port,
		// MaxOutgoingConnections is the maximum outgoing connections allowed.
		MaxOutgoingConnections: 8,
		// MaxDefaultOutgoingConnections is the maximum default outgoing connections allowed.
		MaxDefaultPeerOutgoingConnections: 1,
		DownloadPeerList:                  true,
		PeerListURL:                       node.PeerListURL,
		// How often to make outgoing connections, in seconds
		OutgoingConnectionsRate: time.Second * 5,
		PeerlistSize:            65535,
		// Wallet Address Version
		//AddressVersion: "test",
		// Remote web interface
		WebInterface:      true,
		WebInterfacePort:  node.WebInterfacePort,
		WebInterfaceAddr:  "127.0.0.1",
		WebInterfaceCert:  "",
		WebInterfaceKey:   "",
		WebInterfaceHTTPS: false,
		EnabledAPISets:    collections.NewStringSet(),
		DisabledAPISets:   collections.NewStringSet(),

		RPCInterface: false,

		LaunchBrowser: false,
		// Data directory holds app data
		DataDirectory: node.DataDirectory,
		// Web GUI static resources
		GUIDirectory: "./src/gui/static/",
		// Logging
		ColorLog:        true,
		LogLevel:        "INFO",
		LogToFile:       false,
		DisablePingPong: false,

		VerifyDB:       true,
		ResetCorruptDB: false,

		// Wallets
		WalletDirectory:  "",
		WalletCryptoType: string(wallet.CryptoTypeScryptChacha20poly1305),

		// Timeout settings for http.Server
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 120,

		// Centralized network configuration
		RunMaster: false,
		/* Developer options */

		// Enable cpu profiling
		ProfileCPU: false,
		// Where the file is written to
		ProfileCPUFile: "cpu.prof",
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf:     false,
		HTTPProfHost: "localhost:6060",
	}

	nodeConfig.applyConfigMode(mode)

	return nodeConfig
}

func exitWithMessage(errorCode int, message string) {
	fmt.Println(message)
	os.Exit(errorCode)
}

var allAPISets = strings.Split(api.EndpointsAll, ",")

func (c *Config) postProcess() {
	if help {
		flag.Usage()
		os.Exit(0)
	}

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

	// Post process API set flags
	if c.Node.EnableAllAPISets && c.Node.DisableAllAPISets {
		exitWithMessage(1, "Impossible to enable and disable all API sets at the same time")
	}
	conflicts := collections.StringSetIntersection(c.Node.EnabledAPISets, c.Node.DisabledAPISets)
	if len(conflicts) > 0 {
		exitWithMessage(1, "Impossible to enable and disable "+conflicts.String()+" endpoints at the same time")
	}
	originalEnabledAPISets := c.Node.EnabledAPISets
	// step 3 - If -disable-all-apisets, turn off all API sets.
	// We do this instead of following the exact order mentioned in
	// https://github.com/skycoin/skycoin/issues/1897
	// because there is no point in enabling default API sets at the
	// beginning if them all are meant to be removed later
	if c.Node.DisableAllAPISets {
		c.Node.EnabledAPISets = collections.NewStringSet()
	} else {
		// step 2 - If -enable-all-apisets, turn on all API sets.
		// We do this instead of following the exact order mentioned in
		// https://github.com/skycoin/skycoin/issues/1897
		// because there is no point in enabling default API sets at the
		// beginning if they are meant to be enabled altogether.
		if c.Node.EnableAllAPISets {
			c.Node.EnabledAPISets = collections.NewStringSet(allAPISets...)
		} else {
			// step 1 - Starting from the default configuration
			// FIXME: ... which is different for each run mode
			c.Node.EnabledAPISets = collections.NewStringSet(api.EndpointsRead)
			if c.Node.EnableWalletAPI {
				c.Node.EnabledAPISets[api.EndpointsWallet] = struct{}{}
			}
			if c.Node.EnableSeedAPI {
				c.Node.EnabledAPISets[api.EndpointsWalletSeed] = struct{}{}
			}
		}
	}
	// step 4 - For each API set in -disable-apiset, disable that set.
	for apiSetName := range c.Node.DisabledAPISets {
		c.Node.EnabledAPISets.Remove(apiSetName)
	}
	// step 5 - For each API set in -enable-apiset, enable that set.
	for apiSetName := range originalEnabledAPISets {
		c.Node.EnabledAPISets[apiSetName] = struct{}{}
	}

	// Accept the deprecated -enable options
	if c.Node.EnableWalletAPI {
		c.Node.EnabledAPISets[api.EndpointsWallet] = struct{}{}
	}
	if c.Node.EnableSeedAPI {
		c.Node.EnabledAPISets[api.EndpointsWalletSeed] = struct{}{}
	}
	_, enableWalletAPI := c.Node.EnabledAPISets[api.EndpointsWallet]
	_, enableWalletSeedAPI := c.Node.EnabledAPISets[api.EndpointsWalletSeed]
	c.Node.EnableWalletAPI = c.Node.EnableWalletAPI || enableWalletAPI
	c.Node.EnableSeedAPI = c.Node.EnableSeedAPI || enableWalletSeedAPI

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
	if _, ok := c.Node.EnabledAPISets[api.EndpointsWallet]; !ok {
		c.Node.EnableGUI = false
		c.Node.LaunchBrowser = false
	}

	if c.Node.EnableGUI {
		c.Node.GUIDirectory = file.ResolveResourceDirectory(c.Node.GUIDirectory)
	}

	if c.Node.DisableDefaultPeers {
		c.Node.DefaultConnections = nil
	}
}

// RegisterFlags binds CLI flags to config values
func (c *NodeConfig) RegisterFlags() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX, "disable PEX peer discovery")
	flag.BoolVar(&c.DownloadPeerList, "download-peerlist", c.DownloadPeerList, "download a peers.txt from -peerlist-url")
	flag.StringVar(&c.PeerListURL, "peerlist-url", c.PeerListURL, "with -download-peerlist=true, download a peers.txt file from this url")
	flag.BoolVar(&c.DisableOutgoingConnections, "disable-outgoing", c.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.DisableIncomingConnections, "disable-incoming", c.DisableIncomingConnections, "Don't make incoming connections")
	flag.BoolVar(&c.DisableNetworking, "disable-networking", c.DisableNetworking, "Disable all network activity")
	flag.BoolVar(&c.EnableWalletAPI, "enable-wallet-api", c.EnableWalletAPI, "[DEPRECATED, use -enable-api-set] Enable the wallet API")
	flag.BoolVar(&c.EnableGUI, "enable-gui", c.EnableGUI, "Enable GUI")
	flag.BoolVar(&c.EnableUnversionedAPI, "enable-unversioned-api", c.EnableUnversionedAPI, "Enable the deprecated unversioned API endpoints without /api/v1 prefix")
	flag.BoolVar(&c.DisableCSRF, "disable-csrf", c.DisableCSRF, "disable CSRF check")
	flag.BoolVar(&c.EnableSeedAPI, "enable-seed-api", c.EnableSeedAPI, "[DEPRECATED, use -enable-api-set] enable /api/v1/wallet/seed api")
	flag.BoolVar(&c.DisableCSP, "disable-csp", c.DisableCSP, "disable content-security-policy in http response")
	flag.StringVar(&c.Address, "address", c.Address, "IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")

	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterfaceCert, "web-interface-cert", c.WebInterfaceCert, "cert.pem file for web interface HTTPS. If not provided, will use cert.pem in -data-directory")
	flag.StringVar(&c.WebInterfaceKey, "web-interface-key", c.WebInterfaceKey, "key.pem file for web interface HTTPS. If not provided, will use key.pem in -data-directory")
	flag.BoolVar(&c.WebInterfaceHTTPS, "web-interface-https", c.WebInterfaceHTTPS, "enable HTTPS for web interface")
	flag.BoolVar(&c.EnableAllAPISets, "enable-all-api-sets", c.EnableAllAPISets, "Export all endpoints via web interface. Option -disable-api-set can override this.")
	flag.BoolVar(&c.DisableAllAPISets, "disable-all-api-sets", c.DisableAllAPISets, "Forbid access to all web interface endpoints. Option -enable-api-set can override this.")

	flagutil.StringSetVar(&c.EnabledAPISets, allAPISets, false, "enable-api-set", "enable API set. Options are "+api.EndpointsAll+". Multiple comma-separated values may be specified")
	flagutil.StringSetVar(&c.DisabledAPISets, allAPISets, false, "disable-api-set", "disable API set. Options are "+api.EndpointsAll+". Multiple comma separated values may be specified")

	flag.BoolVar(&c.RPCInterface, "rpc-interface", c.RPCInterface, "enable the deprecated JSON 2.0 RPC interface")

	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
	flag.BoolVar(&c.PrintWebInterfaceAddress, "print-web-interface-address", c.PrintWebInterfaceAddress, "print configured web interface address and exit")
	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, "directory to store app data (defaults to ~/.skycoin)")
	flag.StringVar(&c.DBPath, "db-path", c.DBPath, "path of database file (defaults to ~/.skycoin/data.db)")
	flag.BoolVar(&c.DBReadOnly, "db-read-only", c.DBReadOnly, "open bolt db read-only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "run the HTTP profiling interface")
	flag.StringVar(&c.HTTPProfHost, "http-prof-host", c.HTTPProfHost, "hostname to bind the HTTP profiling interface to")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.BoolVar(&c.DisablePingPong, "no-ping-log", c.DisablePingPong, `disable "reply to ping" and "received pong" debug log messages`)
	flag.BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory, "static content directory for the HTML interface")

	flag.BoolVar(&c.VerifyDB, "verify-db", c.VerifyDB, "check the database for corruption")
	flag.BoolVar(&c.ResetCorruptDB, "reset-corrupt-db", c.ResetCorruptDB, "reset the database if corrupted, and continue running instead of exiting")

	flag.BoolVar(&c.DisableDefaultPeers, "disable-default-peers", c.DisableDefaultPeers, "disable the hardcoded default peers")
	flag.StringVar(&c.CustomPeersFile, "custom-peers-file", c.CustomPeersFile, "load custom peers from a newline separate list of ip:port in a file. Note that this is different from the peers.json file in the data directory")

	// Key Configuration Data
	flag.BoolVar(&c.RunMaster, "master", c.RunMaster, "run the daemon as blockchain master server")

	flag.StringVar(&c.BlockchainPubkeyStr, "master-public-key", c.BlockchainPubkeyStr, "public key of the master chain")
	flag.StringVar(&c.BlockchainSeckeyStr, "master-secret-key", c.BlockchainSeckeyStr, "secret key, set for master")

	flag.StringVar(&c.GenesisAddressStr, "genesis-address", c.GenesisAddressStr, "genesis address")
	flag.StringVar(&c.GenesisSignatureStr, "genesis-signature", c.GenesisSignatureStr, "genesis block signature")
	flag.Uint64Var(&c.GenesisTimestamp, "genesis-timestamp", c.GenesisTimestamp, "genesis block timestamp")

	flag.StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory, "location of the wallet files. Defaults to ~/.skycoin/wallet/")
	flag.IntVar(&c.MaxOutgoingConnections, "max-outgoing-connections", c.MaxOutgoingConnections, "The maximum outgoing connections allowed")
	flag.IntVar(&c.MaxDefaultPeerOutgoingConnections, "max-default-peer-outgoing-connections", c.MaxDefaultPeerOutgoingConnections, "The maximum default peer outgoing connections allowed")
	flag.IntVar(&c.PeerlistSize, "peerlist-size", c.PeerlistSize, "The peer list size")
	flag.DurationVar(&c.OutgoingConnectionsRate, "connection-rate", c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly, "Run on localhost and only connect to localhost peers")
	flag.BoolVar(&c.Arbitrating, "arbitrating", c.Arbitrating, "Run node in arbitrating mode")
	flag.StringVar(&c.WalletCryptoType, "wallet-crypto-type", c.WalletCryptoType, "wallet crypto type. Can be sha256-xor or scrypt-chacha20poly1305")
	flag.BoolVar(&c.Version, "version", false, "show node version")
}

func (c *NodeConfig) applyConfigMode(configMode string) {
	if runtime.GOOS == "windows" {
		c.ColorLog = false
	}
	switch configMode {
	case "":
	case "STANDALONE_CLIENT":
		c.EnabledAPISets = collections.NewStringSet(allAPISets...)
		c.EnableGUI = true
		c.LaunchBrowser = true
		c.DisableCSRF = false
		c.DisableCSP = false
		c.DownloadPeerList = true
		c.RPCInterface = false
		c.WebInterface = true
		c.LogToFile = false
		c.ResetCorruptDB = true
		c.WebInterfacePort = 0 // randomize web interface port
	default:
		panic("Invalid ConfigMode")
	}
}

func panicIfError(err error, msg string, args ...interface{}) { // nolint: unparam
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}
