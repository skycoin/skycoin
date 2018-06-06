package main

import (
	_ "net/http/pprof"
	"time"

	"github.com/skycoin/skycoin/src/skycoin"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	// Version of the node. Can be set by -ldflags
	Version = "0.23.1-rc2"
	// Commit ID. Can be set by -ldflags
	Commit = ""
	// Branch name. Can be set by -ldflags
	Branch = ""
	// ConfigMode (possible values are "", "STANDALONE_CLIENT").
	// This is used to change the default configuration.
	// Can be set by -ldflags
	ConfigMode = ""

	logger = logging.MustGetLogger("main")

	// GenesisSignatureStr hex string of genesis signature
	GenesisSignatureStr = "eb10468d10054d15f2b6f8946cd46797779aa20a7617ceb4be884189f219bc9a164e56a5b9f7bec392a804ff3740210348d73db77a37adb542a8e08d429ac92700"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1426562704
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 100000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"118.178.135.93:6000",
		"47.88.33.156:6000",
		"121.41.103.148:6000",
		"120.77.69.188:6000",
		"104.237.142.206:6000",
		"176.58.126.224:6000",
		"172.104.85.6:6000",
		"139.162.7.132:6000",
	}
)

var devConfig = skycoin.NodeConfig{
	GenesisSignatureStr: GenesisSignatureStr,
	GenesisAddressStr:   GenesisAddressStr,
	GenesisCoinVolume:   GenesisCoinVolume,
	GenesisTimestamp:    GenesisTimestamp,
	BlockchainPubkeyStr: BlockchainPubkeyStr,
	BlockchainSeckeyStr: BlockchainSeckeyStr,
	DefaultConnections:  DefaultConnections,
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
	Port: 6000,
	// MaxOutgoingConnections is the maximum outgoing connections allowed.
	MaxOutgoingConnections: 8,
	// MaxDefaultOutgoingConnections is the maximum default outgoing connections allowed.
	MaxDefaultPeerOutgoingConnections: 1,
	DownloadPeerList:                  true,
	PeerListURL:                       "https://downloads.skycoin.net/blockchain/peers.txt",
	// How often to make outgoing connections, in seconds
	OutgoingConnectionsRate: time.Second * 5,
	PeerlistSize:            65535,
	// Wallet Address Version
	//AddressVersion: "test",
	// Remote web interface
	WebInterface:      true,
	WebInterfacePort:  6420,
	WebInterfaceAddr:  "127.0.0.1",
	WebInterfaceCert:  "",
	WebInterfaceKey:   "",
	WebInterfaceHTTPS: false,

	RPCInterface: true,

	LaunchBrowser: false,
	// Data directory holds app data -- defaults to $HOME/.skycoin
	DataDirectory: "$HOME/.skycoin",
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
	ProfileCPUFile: "skycoin.prof",
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf: false,
}

func init() {
	applyConfigMode()
}

func applyConfigMode() {
	switch ConfigMode {
	case "":
	case "STANDALONE_CLIENT":
		devConfig.EnableWalletAPI = true
		devConfig.EnableGUI = true
		devConfig.EnableSeedAPI = true
		devConfig.LaunchBrowser = true
		devConfig.DisableCSRF = false
		devConfig.DownloadPeerList = true
		devConfig.RPCInterface = false
		devConfig.WebInterface = true
		devConfig.LogToFile = false
		devConfig.ColorLog = true
		devConfig.ResetCorruptDB = true
		devConfig.WebInterfacePort = 0 // randomize web interface port
	default:
		panic("Invalid ConfigMode")
	}
}

func main() {
	// create a new fiber coin instance
	coin := skycoin.NewCoin(
		skycoin.Config{
			Node: devConfig,
			Build: visor.BuildInfo{
				Version: Version,
				Commit:  Commit,
				Branch:  Branch,
			},
		},
		logger,
	)

	// parse config values
	coin.ParseConfig()

	// run fiber coin node
	coin.Run()
}
