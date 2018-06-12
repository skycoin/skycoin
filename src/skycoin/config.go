package skycoin

import (
	"flag"
	"path/filepath"
	"strings"
	"time"

	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	help = false
)

// Config records skycoin node and build config
type Config struct {
	Node  NodeConfig
	Build visor.BuildInfo
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
	// Enable wallet API
	EnableWalletAPI bool
	// Enable GUI
	EnableGUI bool
	// Disable CSRF check in the wallet API
	DisableCSRF bool
	// Enable /api/v1/wallet/seed API endpoint
	EnableSeedAPI bool
	// Enable unversioned API endpoints (without the /api/v1 prefix)
	EnableUnversionedAPI bool

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
	WebInterface      bool
	WebInterfacePort  int
	WebInterfaceAddr  string
	WebInterfaceCert  string
	WebInterfaceKey   string
	WebInterfaceHTTPS bool

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

	RunMaster bool

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool

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
	genesisTimestamp uint64
	genesisAddress   cipher.Address

	blockchainPubkey cipher.PubKey
	blockchainSeckey cipher.SecKey
}

// NewNodeConfig returns a new node config instance
func NewNodeConfig(mode string, node NodeParameters) *NodeConfig {
	nodeConfig := &NodeConfig{
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

		RPCInterface: true,

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
		ProfileCPUFile: node.ProfileCPUFile,
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf: false,
	}

	nodeConfig.applyConfigMode(mode)

	return nodeConfig
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

func (n *NodeConfig) applyConfigMode(configMode string) {
	switch configMode {
	case "":
	case "STANDALONE_CLIENT":
		n.EnableWalletAPI = true
		n.EnableGUI = true
		n.EnableSeedAPI = true
		n.LaunchBrowser = true
		n.DisableCSRF = false
		n.DownloadPeerList = true
		n.RPCInterface = false
		n.WebInterface = true
		n.LogToFile = false
		n.ColorLog = true
		n.ResetCorruptDB = true
		n.WebInterfacePort = 0 // randomize web interface port
	default:
		panic("Invalid ConfigMode")
	}
}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}
