package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util/browser"
	"github.com/skycoin/skycoin/src/util/cert"
	"github.com/skycoin/skycoin/src/util/file"
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
	// Disables wallet API
	EnableWalletAPI bool
	// Disable CSRF check in the wallet api
	DisableCSRF bool
	// Enable /wallet/seed api endpoint
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
	// GUI directory contains assets for the html gui
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
	// Will force it to connect to this ip:port, instead of waiting for it
	// to show up as a peer
	ConnectTo string

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
	flag.StringVar(&c.ConnectTo, "connect-to", c.ConnectTo, "connect to this ip only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "Run the http profiling interface")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.BoolVar(&c.DisablePingPong, "no-ping-log", c.DisablePingPong, `disable "reply to ping" and "received pong" debug log messages`)
	flag.BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory, "static content directory for the html gui")

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
	// Enable seed API
	EnableSeedAPI: false,
	// Disable CSRF check in the wallet api
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
	// Will force it to connect to this ip:port, instead of waiting for it
	// to show up as a peer
	ConnectTo: "",
}

func init() {
	applyConfigMode()
}

func applyConfigMode() {
	switch ConfigMode {
	case "":
	case "STANDALONE_CLIENT":
		devConfig.EnableWalletAPI = true
		devConfig.EnableSeedAPI = true
		devConfig.LaunchBrowser = true
		devConfig.DisableCSRF = false
		devConfig.DownloadPeerList = true
		devConfig.RPCInterface = false
		devConfig.WebInterface = true
		devConfig.LogToFile = false
		devConfig.ColorLog = true
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

	// Don't open browser to load wallets if wallet apis are disabled.
	if c.EnableWalletAPI {
		c.GUIDirectory = file.ResolveResourceDirectory(c.GUIDirectory)
	} else {
		c.LaunchBrowser = false
	}
}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func printProgramStatus() {
	p := pprof.Lookup("goroutine")
	if err := p.WriteTo(os.Stdout, 2); err != nil {
		fmt.Println("ERROR:", err)
		return
	}
}

// Catches SIGUSR1 and prints internal program state
func catchDebug() {
	sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGUSR1)
	signal.Notify(sigchan, syscall.Signal(0xa)) // SIGUSR1 = Signal(0xa)
	for {
		select {
		case <-sigchan:
			printProgramStatus()
		}
	}
}

func catchInterrupt(quit chan<- struct{}) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	close(quit)

	// If ctrl-c is called again, panic so that the program state can be examined.
	// Ctrl-c would be called again if program shutdown was stuck.
	go catchInterruptPanic()
}

// catchInterruptPanic catches os.Interrupt and panics
func catchInterruptPanic() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	printProgramStatus()
	panic("SIGINT")
}

func createGUI(c *Config, d *daemon.Daemon, host string) (*gui.Server, error) {
	var s *gui.Server
	var err error

	config := gui.Config{
		StaticDir:       c.GUIDirectory,
		DisableCSRF:     c.DisableCSRF,
		EnableWalletAPI: c.EnableWalletAPI,
		EnableJSON20RPC: c.RPCInterface,
		ReadTimeout:     c.ReadTimeout,
		WriteTimeout:    c.WriteTimeout,
		IdleTimeout:     c.IdleTimeout,
	}

	if c.WebInterfaceHTTPS {
		// Verify cert/key parameters, and if neither exist, create them
		if err := cert.CreateCertIfNotExists(host, c.WebInterfaceCert, c.WebInterfaceKey, "Skycoind"); err != nil {
			logger.Errorf("gui.CreateCertIfNotExists failure: %v", err)
			return nil, err
		}

		s, err = gui.CreateHTTPS(host, config, d, c.WebInterfaceCert, c.WebInterfaceKey)
	} else {
		s, err = gui.Create(host, config, d)
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
	dc.Visor.Config.DBReadOnly = c.DBReadOnly
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

	// If the user Ctrl-C's, shutdown properly
	quit := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		catchInterrupt(quit)
	}()

	// Watch for SIGUSR1
	wg.Add(1)
	func() {
		defer wg.Done()
		go catchDebug()
	}()

	// creates blockchain instance
	dconf := configureDaemon(c)

	logger.Infof("Opening database %s", dconf.Visor.Config.DBPath)
	db, err := visor.OpenDB(dconf.Visor.Config.DBPath, dconf.Visor.Config.DBReadOnly)
	if err != nil {
		logger.Errorf("Database failed to open: %v. Is another skycoin instance running?", err)
		return
	}

	d, err := daemon.NewDaemon(dconf, db, DefaultConnections)
	if err != nil {
		logger.Error(err)
		return
	}

	var webInterface *gui.Server
	if c.WebInterface {
		webInterface, err = createGUI(c, d, host)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	// Debug only - forces connection on start.  Violates thread safety.
	if c.ConnectTo != "" {
		if err := d.Pool.Pool.Connect(c.ConnectTo); err != nil {
			logger.Errorf("Force connect %s failed, %v", c.ConnectTo, err)
			return
		}
	}

	errC := make(chan error, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := d.Run(); err != nil {
			logger.Error(err)
			errC <- err
		}
	}()

	if c.WebInterface {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := webInterface.Serve(); err != nil {
				logger.Error(err)
				errC <- err
			}
		}()

		if c.LaunchBrowser {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Wait a moment just to make sure the http interface is up
				time.Sleep(time.Millisecond * 100)

				logger.Infof("Launching System Browser with %s", fullAddress)
				if err := browser.Open(fullAddress); err != nil {
					logger.Error(err)
					return
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
		webInterface.Shutdown()
	}
	d.Shutdown()
	wg.Wait()

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
