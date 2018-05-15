package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/apputil"
	"github.com/skycoin/skycoin/src/util/browser"
	"github.com/skycoin/skycoin/src/util/cert"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
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

	help = false

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
	GenesisCoinVolume uint64 = 100e12

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

// Config records the node's configuration
type Config struct {
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
	// Enable /wallet/seed API endpoint
	EnableSeedAPI bool

	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	//gnet uses this for TCP incoming and outgoing
	Port int
	//max outgoing connections to maintain
	MaxOutgoingConnections int
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

	GenesisSignature cipher.Sig
	GenesisTimestamp uint64
	GenesisAddress   cipher.Address

	BlockchainPubkey cipher.PubKey
	BlockchainSeckey cipher.SecKey

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
}

func (c *Config) register() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX, "disable PEX peer discovery")
	flag.BoolVar(&c.DownloadPeerList, "download-peerlist", c.DownloadPeerList, "download a peers.txt from -peerlist-url")
	flag.StringVar(&c.PeerListURL, "peerlist-url", c.PeerListURL, "with -download-peerlist=true, download a peers.txt file from this url")
	flag.BoolVar(&c.DisableOutgoingConnections, "disable-outgoing", c.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.DisableIncomingConnections, "disable-incoming", c.DisableIncomingConnections, "Don't make incoming connections")
	flag.BoolVar(&c.DisableNetworking, "disable-networking", c.DisableNetworking, "Disable all network activity")
	flag.BoolVar(&c.EnableWalletAPI, "enable-wallet-api", c.EnableWalletAPI, "Enable the wallet API")
	flag.BoolVar(&c.EnableGUI, "enable-gui", c.EnableGUI, "Enable GUI")
	flag.BoolVar(&c.DisableCSRF, "disable-csrf", c.DisableCSRF, "disable CSRF check")
	flag.BoolVar(&c.EnableSeedAPI, "enable-seed-api", c.EnableSeedAPI, "enable /wallet/seed api")
	flag.StringVar(&c.Address, "address", c.Address, "IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")

	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterfaceCert, "web-interface-cert", c.WebInterfaceCert, "cert.pem file for web interface HTTPS. If not provided, will use cert.pem in -data-directory")
	flag.StringVar(&c.WebInterfaceKey, "web-interface-key", c.WebInterfaceKey, "key.pem file for web interface HTTPS. If not provided, will use key.pem in -data-directory")
	flag.BoolVar(&c.WebInterfaceHTTPS, "web-interface-https", c.WebInterfaceHTTPS, "enable HTTPS for web interface")

	flag.BoolVar(&c.RPCInterface, "rpc-interface", c.RPCInterface, "enable the rpc interface")

	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
	flag.BoolVar(&c.PrintWebInterfaceAddress, "print-web-interface-address", c.PrintWebInterfaceAddress, "print configured web interface address and exit")
	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, "directory to store app data (defaults to ~/.skycoin)")
	flag.StringVar(&c.DBPath, "db-path", c.DBPath, "path of database file (defaults to ~/.skycoin/data.db)")
	flag.BoolVar(&c.DBReadOnly, "db-read-only", c.DBReadOnly, "open bolt db read-only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "Run the http profiling interface")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.BoolVar(&c.DisablePingPong, "no-ping-log", c.DisablePingPong, `disable "reply to ping" and "received pong" debug log messages`)
	flag.BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory, "static content directory for the HTML interface")

	flag.BoolVar(&c.VerifyDB, "verify-db", c.VerifyDB, "check the database for corruption")
	flag.BoolVar(&c.ResetCorruptDB, "reset-corrupt-db", c.ResetCorruptDB, "reset the database if corrupted, and continue running instead of exiting")

	// Key Configuration Data
	flag.BoolVar(&c.RunMaster, "master", c.RunMaster, "run the daemon as blockchain master server")

	flag.StringVar(&BlockchainPubkeyStr, "master-public-key", BlockchainPubkeyStr, "public key of the master chain")
	flag.StringVar(&BlockchainSeckeyStr, "master-secret-key", BlockchainSeckeyStr, "secret key, set for master")

	flag.StringVar(&GenesisAddressStr, "genesis-address", GenesisAddressStr, "genesis address")
	flag.StringVar(&GenesisSignatureStr, "genesis-signature", GenesisSignatureStr, "genesis block signature")
	flag.Uint64Var(&c.GenesisTimestamp, "genesis-timestamp", c.GenesisTimestamp, "genesis block timestamp")

	flag.StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory, "location of the wallet files. Defaults to ~/.skycoin/wallet/")
	flag.IntVar(&c.MaxOutgoingConnections, "max-outgoing-connections", c.MaxOutgoingConnections, "The maximum outgoing connections allowed")
	flag.IntVar(&c.PeerlistSize, "peerlist-size", c.PeerlistSize, "The peer list size")
	flag.DurationVar(&c.OutgoingConnectionsRate, "connection-rate", c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly, "Run on localhost and only connect to localhost peers")
	flag.BoolVar(&c.Arbitrating, "arbitrating", c.Arbitrating, "Run node in arbitrating mode")
	flag.StringVar(&c.WalletCryptoType, "wallet-crypto-type", c.WalletCryptoType, "wallet crypto type. Can be sha256-xor or scrypt-chacha20poly1305")
	flag.BoolVar(&c.Version, "version", false, "show node version")
}

var home = file.UserHome()

var devConfig = Config{
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
	MaxOutgoingConnections: 16,
	DownloadPeerList:       false,
	PeerListURL:            "https://downloads.skycoin.net/blockchain/peers.txt",
	// How often to make outgoing connections, in seconds
	OutgoingConnectionsRate: time.Second * 5,
	PeerlistSize:            65535,
	// Wallet Address Version
	//AddressVersion: "test",
	// Remote web interface
	WebInterface:             true,
	WebInterfacePort:         6420,
	WebInterfaceAddr:         "127.0.0.1",
	WebInterfaceCert:         "",
	WebInterfaceKey:          "",
	WebInterfaceHTTPS:        false,
	PrintWebInterfaceAddress: false,

	RPCInterface: true,

	LaunchBrowser: false,
	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory: filepath.Join(home, ".skycoin"),
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
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 60 * time.Second,
	IdleTimeout:  120 * time.Second,

	// Centralized network configuration
	RunMaster:        false,
	BlockchainPubkey: cipher.PubKey{},
	BlockchainSeckey: cipher.SecKey{},

	GenesisAddress:   cipher.Address{},
	GenesisTimestamp: GenesisTimestamp,
	GenesisSignature: cipher.Sig{},

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
	default:
		panic("Invalid ConfigMode")
	}
}

// Parse prepare the config
func (c *Config) Parse() {
	c.register()
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	c.postProcess()
}

func (c *Config) postProcess() {
	var err error
	if GenesisSignatureStr != "" {
		c.GenesisSignature, err = cipher.SigFromHex(GenesisSignatureStr)
		panicIfError(err, "Invalid Signature")
	}
	if GenesisAddressStr != "" {
		c.GenesisAddress, err = cipher.DecodeBase58Address(GenesisAddressStr)
		panicIfError(err, "Invalid Address")
	}
	if BlockchainPubkeyStr != "" {
		c.BlockchainPubkey, err = cipher.PubKeyFromHex(BlockchainPubkeyStr)
		panicIfError(err, "Invalid Pubkey")
	}
	if BlockchainSeckeyStr != "" {
		c.BlockchainSeckey, err = cipher.SecKeyFromHex(BlockchainSeckeyStr)
		panicIfError(err, "Invalid Seckey")
		BlockchainSeckeyStr = ""
	}
	if BlockchainSeckeyStr != "" {
		c.BlockchainSeckey = cipher.SecKey{}
	}

	c.DataDirectory, err = file.InitDataDir(c.DataDirectory)
	panicIfError(err, "Invalid DataDirectory")

	if c.WebInterfaceCert == "" {
		c.WebInterfaceCert = filepath.Join(c.DataDirectory, "cert.pem")
	}
	if c.WebInterfaceKey == "" {
		c.WebInterfaceKey = filepath.Join(c.DataDirectory, "key.pem")
	}

	if c.WalletDirectory == "" {
		c.WalletDirectory = filepath.Join(c.DataDirectory, "wallets")
	}

	if c.DBPath == "" {
		c.DBPath = filepath.Join(c.DataDirectory, "data.db")
	}

	if c.RunMaster {
		// Run in arbitrating mode if the node is master
		c.Arbitrating = true
	}

	// Don't open browser to load wallets if wallet APIs are disabled.
	if !c.EnableWalletAPI {
		c.EnableGUI = false
		c.LaunchBrowser = false
	}

	if c.EnableGUI {
		c.GUIDirectory = file.ResolveResourceDirectory(c.GUIDirectory)
	}
}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func createGUI(c *Config, d *daemon.Daemon, host string) (*api.Server, error) {
	var s *api.Server
	var err error

	config := api.Config{
		StaticDir:       c.GUIDirectory,
		DisableCSRF:     c.DisableCSRF,
		EnableWalletAPI: c.EnableWalletAPI,
		EnableJSON20RPC: c.RPCInterface,
		EnableGUI:       c.EnableGUI,
		ReadTimeout:     c.ReadTimeout,
		WriteTimeout:    c.WriteTimeout,
		IdleTimeout:     c.IdleTimeout,
	}

	if c.WebInterfaceHTTPS {
		// Verify cert/key parameters, and if neither exist, create them
		if err := cert.CreateCertIfNotExists(host, c.WebInterfaceCert, c.WebInterfaceKey, "Skycoind"); err != nil {
			logger.Errorf("cert.CreateCertIfNotExists failure: %v", err)
			return nil, err
		}

		s, err = api.CreateHTTPS(host, config, d.Gateway, c.WebInterfaceCert, c.WebInterfaceKey)
	} else {
		s, err = api.Create(host, config, d.Gateway)
	}
	if err != nil {
		logger.Errorf("Failed to start web GUI: %v", err)
		return nil, err
	}

	return s, nil
}

func initLogFile(dataDir string) (*os.File, error) {
	logDir := filepath.Join(dataDir, "logs")
	if err := createDirIfNotExist(logDir); err != nil {
		logger.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
		return nil, fmt.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
	}

	// open log file
	tf := "2006-01-02-030405"
	logfile := filepath.Join(logDir, fmt.Sprintf("%s-v%s.log", time.Now().Format(tf), Version))

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		logger.Errorf("os.OpenFile(%s) failed: %v", logfile, err)
		return nil, err
	}

	hook := logging.NewWriteHook(f)
	logging.AddHook(hook)

	return f, nil
}

func initProfiling(httpProf, profileCPU bool, profileCPUFile string) {
	if profileCPU {
		f, err := os.Create(profileCPUFile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if httpProf {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}

func configureDaemon(c *Config) daemon.Config {
	//cipher.SetAddressVersion(c.AddressVersion)
	dc := daemon.NewConfig()
	dc.Pex.DataDirectory = c.DataDirectory
	dc.Pex.Disabled = c.DisablePEX
	dc.Pex.Max = c.PeerlistSize
	dc.Pex.DownloadPeerList = c.DownloadPeerList
	dc.Pex.PeerListURL = c.PeerListURL
	dc.Daemon.DisableOutgoingConnections = c.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = c.DisableIncomingConnections
	dc.Daemon.DisableNetworking = c.DisableNetworking
	dc.Daemon.Port = c.Port
	dc.Daemon.Address = c.Address
	dc.Daemon.LocalhostOnly = c.LocalhostOnly
	dc.Daemon.OutgoingMax = c.MaxOutgoingConnections
	dc.Daemon.DataDirectory = c.DataDirectory
	dc.Daemon.LogPings = !c.DisablePingPong

	if c.OutgoingConnectionsRate == 0 {
		c.OutgoingConnectionsRate = time.Millisecond
	}
	dc.Daemon.OutgoingRate = c.OutgoingConnectionsRate
	dc.Visor.Config.IsMaster = c.RunMaster

	dc.Visor.Config.BlockchainPubkey = c.BlockchainPubkey
	dc.Visor.Config.BlockchainSeckey = c.BlockchainSeckey

	dc.Visor.Config.GenesisAddress = c.GenesisAddress
	dc.Visor.Config.GenesisSignature = c.GenesisSignature
	dc.Visor.Config.GenesisTimestamp = c.GenesisTimestamp
	dc.Visor.Config.GenesisCoinVolume = GenesisCoinVolume
	dc.Visor.Config.DBPath = c.DBPath
	dc.Visor.Config.Arbitrating = c.Arbitrating
	dc.Visor.Config.EnableWalletAPI = c.EnableWalletAPI
	dc.Visor.Config.WalletDirectory = c.WalletDirectory
	dc.Visor.Config.BuildInfo = visor.BuildInfo{
		Version: Version,
		Commit:  Commit,
		Branch:  Branch,
	}
	dc.Visor.Config.EnableSeedAPI = c.EnableSeedAPI

	dc.Gateway.EnableWalletAPI = c.EnableWalletAPI

	// Initialize wallet default crypto type
	cryptoType, err := wallet.CryptoTypeFromString(c.WalletCryptoType)
	if err != nil {
		log.Panic(err)
	}

	dc.Visor.Config.WalletCryptoType = cryptoType

	return dc
}

// Run starts the skycoin node
func Run(c *Config) {
	defer func() {
		// try catch panic in main thread
		if r := recover(); r != nil {
			logger.Errorf("recover: %v\nstack:%v", r, string(debug.Stack()))
		}
	}()

	var db *dbutil.DB
	var d *daemon.Daemon
	var webInterface *api.Server
	errC := make(chan error, 10)

	if c.Version {
		fmt.Println(Version)
		return
	}

	logLevel, err := logging.LevelFromString(c.LogLevel)
	if err != nil {
		logger.Error("Invalid -log-level:", err)
		return
	}

	logging.SetLevel(logLevel)

	if c.ColorLog {
		logging.EnableColors()
	} else {
		logging.DisableColors()
	}

	var logFile *os.File
	if c.LogToFile {
		var err error
		logFile, err = initLogFile(c.DataDirectory)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	scheme := "http"
	if c.WebInterfaceHTTPS {
		scheme = "https"
	}
	host := fmt.Sprintf("%s:%d", c.WebInterfaceAddr, c.WebInterfacePort)
	fullAddress := fmt.Sprintf("%s://%s", scheme, host)
	logger.Critical().Infof("Full address: %s", fullAddress)
	if c.PrintWebInterfaceAddress {
		fmt.Println(fullAddress)
	}

	initProfiling(c.HTTPProf, c.ProfileCPU, c.ProfileCPUFile)

	var wg sync.WaitGroup

	quit := make(chan struct{})

	// Catch SIGINT (CTRL-C) (closes the quit channel)
	go apputil.CatchInterrupt(quit)

	// Catch SIGUSR1 (prints runtime stack to stdout)
	go apputil.CatchDebug()

	// creates blockchain instance
	dconf := configureDaemon(c)

	logger.Infof("Opening database %s", dconf.Visor.Config.DBPath)
	db, err = visor.OpenDB(dconf.Visor.Config.DBPath, c.DBReadOnly)
	if err != nil {
		logger.Errorf("Database failed to open: %v. Is another skycoin instance running?", err)
		return
	}

	if c.ResetCorruptDB {
		// Check the database integrity and recreate it if necessary
		logger.Info("Checking database and resetting if corrupted")
		if newDB, err := visor.ResetCorruptDB(db, c.BlockchainPubkey, quit); err != nil {
			if err != visor.ErrVerifyStopped {
				logger.Errorf("visor.ResetCorruptDB failed: %v", err)
			}
			goto earlyShutdown
		} else {
			db = newDB
		}
	} else if c.VerifyDB {
		logger.Info("Checking database")
		if err := visor.CheckDatabase(db, c.BlockchainPubkey, quit); err != nil {
			if err != visor.ErrVerifyStopped {
				logger.Errorf("visor.CheckDatabase failed: %v", err)
			}
			goto earlyShutdown
		}
	}

	d, err = daemon.NewDaemon(dconf, db, DefaultConnections)
	if err != nil {
		logger.Error(err)
		goto earlyShutdown
	}

	if c.WebInterface {
		webInterface, err = createGUI(c, d, host)
		if err != nil {
			logger.Error(err)
			goto earlyShutdown
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := d.Run(); err != nil {
			logger.Error(err)
			errC <- err
		}
	}()

	if c.WebInterface {
		cancelLaunchBrowser := make(chan struct{})

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := webInterface.Serve(); err != nil {
				close(cancelLaunchBrowser)
				logger.Error(err)
				errC <- err
			}
		}()

		if c.LaunchBrowser {
			go func() {
				select {
				case <-cancelLaunchBrowser:
					logger.Warning("Browser launching cancelled")

				// Wait a moment just to make sure the http interface is up
				case <-time.After(time.Millisecond * 100):
					logger.Infof("Launching System Browser with %s", fullAddress)
					if err := browser.Open(fullAddress); err != nil {
						logger.Error(err)
					}
				}
			}()
		}
	}

	/*
	   time.Sleep(5)
	   tx := InitTransaction()
	   _ = tx
	   err, _ = d.Visor.Visor.InjectTransaction(tx)
	   if err != nil {
	       log.Panic(err)
	   }
	*/

	/*
	   //first transaction
	   if c.RunMaster == true {
	       go func() {
	           for d.Visor.Visor.Blockchain.Head().Seq() < 2 {
	               time.Sleep(5)
	               tx := InitTransaction()
	               err, _ := d.Visor.Visor.InjectTransaction(tx)
	               if err != nil {
	                   //log.Panic(err)
	               }
	           }
	       }()
	   }
	*/

	select {
	case <-quit:
	case err := <-errC:
		logger.Error(err)
	}

	logger.Info("Shutting down...")

	if webInterface != nil {
		logger.Info("Closing web interface")
		webInterface.Shutdown()
	}

	logger.Info("Closing daemon")
	d.Shutdown()

	logger.Info("Waiting for goroutines to finish")
	wg.Wait()

earlyShutdown:
	if db != nil {
		logger.Info("Closing database")
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close DB")
		}
	}

	logger.Info("Goodbye")

	if logFile != nil {
		if err := logFile.Close(); err != nil {
			fmt.Println("Failed to close log file")
		}
	}
}

func main() {
	devConfig.Parse()
	Run(&devConfig)
}

// InitTransaction creates the initialize transaction
func InitTransaction() coin.Transaction {
	var tx coin.Transaction

	output := cipher.MustSHA256FromHex("043836eb6f29aaeb8b9bfce847e07c159c72b25ae17d291f32125e7f1912e2a0")
	tx.PushInput(output)

	addrs := visor.GetDistributionAddresses()

	if len(addrs) != 100 {
		log.Panic("Should have 100 distribution addresses")
	}

	// 1 million per address, measured in droplets
	if visor.DistributionAddressInitialBalance != 1e6 {
		log.Panic("visor.DistributionAddressInitialBalance expected to be 1e6*1e6")
	}

	for i := range addrs {
		addr := cipher.MustDecodeBase58Address(addrs[i])
		tx.PushOutput(addr, visor.DistributionAddressInitialBalance*1e6, 1)
	}
	/*
		seckeys := make([]cipher.SecKey, 1)
		seckey := ""
		seckeys[0] = cipher.MustSecKeyFromHex(seckey)
		tx.SignInputs(seckeys)
	*/

	txs := make([]cipher.Sig, 1)
	sig := "ed9bd7a31fe30b9e2d53b35154233dfdf48aaaceb694a07142f84cdf4f5263d21b723f631817ae1c1f735bea13f0ff2a816e24a53ccb92afae685fdfc06724de01"
	txs[0] = cipher.MustSigFromHex(sig)
	tx.Sigs = txs

	tx.UpdateHeader()

	err := tx.Verify()

	if err != nil {
		log.Panic(err)
	}

	log.Printf("signature= %s", tx.Sigs[0].Hex())
	return tx
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0777)
}
