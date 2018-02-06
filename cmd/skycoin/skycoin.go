package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util/browser"
	"github.com/skycoin/skycoin/src/util/cert"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor"
)

var (
	// Version node version which will be set when build wallet by LDFLAGS
	Version = "0.21.1"
	// Commit id
	Commit = ""

	help = false

	logger     = logging.MustGetLogger("main")
	logModules = []string{
		"main",
		"daemon",
		"coin",
		"gui",
		"file",
		"visor",
		"wallet",
		"gnet",
		"pex",
		"webrpc",
	}

	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = ""
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckey = ""
	// Name of the file containing trusted peer list (one-by-line)
	TrustedPeerlistFileName = "connections.txt"
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
	DisableWalletApi bool
	// Disable CSRF check in the wallet api
	DisableCSRF bool

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

	RPCInterface     bool
	RPCInterfacePort int
	RPCInterfaceAddr string

	// Launch System Default Browser after client startup
	LaunchBrowser bool

	// If true, print the configured client web interface address and exit
	PrintWebInterfaceAddress bool

	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string
	// GUI directory contains assets for the html gui
	GUIDirectory string

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Disable "Reply to ping", "Received pong" log messages
	DisablePingPong bool

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string

	RunMaster bool

	GenesisSignature  cipher.Sig
	GenesisTimestamp  uint64
	GenesisCoinVolume uint64
	GenesisAddress    cipher.Address

	BlockchainPubkey cipher.PubKey
	BlockchainSeckey cipher.SecKey

	DefaultConnections []string

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

	DBPath       string
	Arbitrating  bool
	RPCThreadNum uint   // rpc number
	LogFmt       string // log format
	Logtofile    bool
	TestChain    bool
	Logtogui     bool
	LogBuffSize  int
}

func (c *Config) register() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX,
		"disable PEX peer discovery")
	flag.BoolVar(&c.DisableOutgoingConnections, "disable-outgoing",
		c.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.DisableIncomingConnections, "disable-incoming",
		c.DisableIncomingConnections, "Don't make incoming connections")
	flag.BoolVar(&c.DisableNetworking, "disable-networking",
		c.DisableNetworking, "Disable all network activity")
	flag.BoolVar(&c.DisableCSRF, "disable-csrf", c.DisableCSRF, "disable csrf check")
	flag.StringVar(&c.Address, "address", c.Address,
		"IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")

	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterfaceCert, "web-interface-cert", c.WebInterfaceCert, "cert.pem file for web interface HTTPS. If not provided, will use cert.pem in -data-directory")
	flag.StringVar(&c.WebInterfaceKey, "web-interface-key", c.WebInterfaceKey, "key.pem file for web interface HTTPS. If not provided, will use key.pem in -data-directory")
	flag.BoolVar(&c.WebInterfaceHTTPS, "web-interface-https", c.WebInterfaceHTTPS, "enable HTTPS for web interface")

	flag.BoolVar(&c.RPCInterface, "rpc-interface", c.RPCInterface, "enable the rpc interface")
	flag.IntVar(&c.RPCInterfacePort, "rpc-interface-port", c.RPCInterfacePort, "port to serve rpc interface on")
	flag.StringVar(&c.RPCInterfaceAddr, "rpc-interface-addr", c.RPCInterfaceAddr, "addr to serve rpc interface on")
	flag.UintVar(&c.RPCThreadNum, "rpc-thread-num", 5, "rpc thread number")

	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
	flag.BoolVar(&c.PrintWebInterfaceAddress, "print-web-interface-address", c.PrintWebInterfaceAddress, "print configured web interface address and exit")
	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, "directory to store app data (defaults to ~/.skycoin)")
	flag.StringVar(&c.ConnectTo, "connect-to", c.ConnectTo, "connect to this ip only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "Run the http profiling interface")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, notice, warning, error, critical")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.BoolVar(&c.DisablePingPong, "no-ping-log", false, `disable "reply to ping" and "received pong" log messages`)
	flag.BoolVar(&c.Logtofile, "logtofile", false, "log to file")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory, "static content directory for the html gui")

	// Key Configuration Data
	flag.BoolVar(&c.RunMaster, "master", c.RunMaster, "run the daemon as blockchain master server")
	flag.StringVar(&BlockchainPubkeyStr, "master-public-key", BlockchainPubkeyStr, "public key of the master chain")
	flag.StringVar(&BlockchainSeckeyStr, "master-secret-key", BlockchainSeckeyStr, "secret key, set for master")

	flag.StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory, "location of the wallet files. Defaults to ~/.skycoin/wallet/")
	flag.IntVar(&c.MaxOutgoingConnections, "max-outgoing-connections", 16, "The maximum outgoing connections allowed")
	flag.IntVar(&c.PeerlistSize, "peerlist-size", 65535, "The peer list size")
	flag.DurationVar(&c.OutgoingConnectionsRate, "connection-rate", c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly, "Run on localhost and only connect to localhost peers")
	flag.BoolVar(&c.Arbitrating, "arbitrating", c.Arbitrating, "Run node in arbitrating mode")
	flag.BoolVar(&c.TestChain, "testchain", false, "Run node in test chain")
	flag.BoolVar(&c.Logtogui, "logtogui", true, "log to gui")
	flag.IntVar(&c.LogBuffSize, "logbufsize", c.LogBuffSize, "Log size saved in memeory for gui show")
}

var devConfig = Config{
	// Disable peer exchange
	DisablePEX: false,
	// Don't make any outgoing connections
	DisableOutgoingConnections: false,
	// Don't allowing incoming connections
	DisableIncomingConnections: false,
	// Disables networking altogether
	DisableNetworking: false,
	// Disable wallet API
	DisableWalletApi: false,
	// Disable CSRF check in the wallet api
	DisableCSRF: false,
	// Only run on localhost and only connect to others on localhost
	LocalhostOnly: false,
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address: "",
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
	WebInterfaceAddr:         "127.0.0.1",
	WebInterfaceCert:         "",
	WebInterfaceKey:          "",
	WebInterfaceHTTPS:        false,
	PrintWebInterfaceAddress: false,

	RPCInterface:     true,
	RPCInterfaceAddr: "127.0.0.1",
	RPCThreadNum:     5,

	LaunchBrowser: true,

	// Web GUI static resources
	GUIDirectory: "./src/gui/static/",
	// Logging
	ColorLog: true,
	LogLevel: "DEBUG",

	// Wallets
	WalletDirectory: "",

	// Centralized network configuration
	RunMaster: false,

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU: false,
	// Where the file is written to
	ProfileCPUFile: "skycoin.prof",
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf: false,
	// Will force it to connect to this ip:port, instead of waiting for it
	// to show up as a peer
	ConnectTo:   "",
	LogBuffSize: 8388608, //1024*1024*8
}

// Parse prepare the config
func (c *Config) Parse() {
	c.register()
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if c.TestChain {
		c.postProcess(TestChainCfg)
		return
	}

	c.postProcess(MainChainCfg)
}

func (c *Config) postProcess(chaincfg ChainConfig) {
	var err error
	// if c.TestChain {
	if chaincfg.GenesisSignature != "" {
		c.GenesisSignature, err = cipher.SigFromHex(chaincfg.GenesisSignature)
		panicIfError(err, "Invalid Signature")
	}

	c.GenesisAddress, err = cipher.DecodeBase58Address(chaincfg.GenesisAddress)
	panicIfError(err, "Invalid address")

	c.BlockchainPubkey, err = cipher.PubKeyFromHex(chaincfg.BlockchainPubkey)
	panicIfError(err, "Invalid Pubkey")

	c.GenesisTimestamp = chaincfg.GenesisTimestamp
	c.GenesisCoinVolume = chaincfg.GenesisCoinVolume

	c.Port = TestChainCfg.Port
	c.WebInterfacePort = chaincfg.WebInterfacePort
	c.RPCInterfacePort = chaincfg.RPCInterfacePort

	if c.DataDirectory == "" {
		c.DataDirectory = chaincfg.DataDirectory
	}
	c.LogFmt = chaincfg.LogFmt
	// } else {
	// if GenesisSignatureStr != "" {
	// 	c.GenesisSignature, err = cipher.SigFromHex(GenesisSignatureStr)
	// 	panicIfError(err, "Invalid Signature")
	// }
	// if GenesisAddressStr != "" {
	// 	c.GenesisAddress, err = cipher.DecodeBase58Address(GenesisAddressStr)
	// 	panicIfError(err, "Invalid Address")
	// }
	// if BlockchainPubkeyStr != "" {
	// 	c.BlockchainPubkey, err = cipher.PubKeyFromHex(BlockchainPubkeyStr)
	// 	panicIfError(err, "Invalid Pubkey")
	// }
	// }

	if BlockchainSeckey != "" {
		c.BlockchainSeckey, err = cipher.SecKeyFromHex(BlockchainSeckey)
		panicIfError(err, "Invalid Seckey")
		BlockchainSeckey = ""
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
	if c.TestChain {
		// Never download peers list if running testnet
		c.DownloadPeerList = false
		c.PeerListURL = ""

		// Force load default connections from file in data dir
		c.DefaultConnections = loadDefaultConnections(c.DataDirectory)
		if len(c.DefaultConnections) == 0 {
			logger.Info("Unable to load dafault connections from %v", c.DataDirectory)
			c.DefaultConnections = chaincfg.DefaultConnections
		}
	} else {
		c.DefaultConnections = chaincfg.DefaultConnections
	}
}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func printProgramStatus() {
	fn := "goroutine.prof"
	logger.Debug("Writing goroutine profile to %s", fn)
	p := pprof.Lookup("goroutine")
	f, err := os.Create(fn)
	defer f.Close()
	if err != nil {
		logger.Error("%v", err)
		return
	}
	err = p.WriteTo(f, 2)
	if err != nil {
		logger.Error("%v", err)
		return
	}
}

func catchInterrupt(quit chan<- struct{}) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	close(quit)
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

func createGUI(c *Config, d *daemon.Daemon, host string, quit chan struct{}) (*gui.Server, error) {
	var s *gui.Server
	var err error

	config := gui.ServerConfig{
		StaticDir:   c.GUIDirectory,
		DisableCSRF: c.DisableCSRF,
	}

	if c.WebInterfaceHTTPS {
		// Verify cert/key parameters, and if neither exist, create them
		if err := cert.CreateCertIfNotExists(host, c.WebInterfaceCert, c.WebInterfaceKey, "Skycoind"); err != nil {
			logger.Error("gui.CreateCertIfNotExists failure: %v", err)
			return nil, err
		}

		s, err = gui.CreateHTTPS(host, config, d, c.WebInterfaceCert, c.WebInterfaceKey)
	} else {
		s, err = gui.Create(host, config, d)
	}
	if err != nil {
		logger.Error("Failed to start web GUI: %v", err)
		return nil, err
	}

	return s, nil
}

// init logging settings
func initLogging(dataDir string, level string, color bool, logfmt string, logtofile bool) (func(), error) {
	logCfg := logging.DevLogConfig(logModules)
	logCfg.Format = logfmt
	logCfg.Colors = color
	logCfg.Level = level

	var fd *os.File
	if logtofile {
		logDir := filepath.Join(dataDir, "logs")
		if err := createDirIfNotExist(logDir); err != nil {
			log.Println("initial logs folder failed", err)
			return nil, fmt.Errorf("init log folder fail, %v", err)
		}

		// open log file
		tf := "2006-01-02-030405"
		logfile := filepath.Join(logDir,
			fmt.Sprintf("%s-v%s.log", time.Now().Format(tf), Version))
		var err error
		fd, err = os.OpenFile(logfile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}

		logCfg.Output = io.MultiWriter(os.Stdout, fd)
	}

	logCfg.InitLogger()

	return func() {
		logger.Info("Log file closed")
		if fd != nil {
			fd.Close()
		}
	}, nil
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

func loadDefaultConnections(dataDirectory string) []string {
	connections := make([]string, 0)
	fp := filepath.Join(dataDirectory, TrustedPeerlistFileName)
	fo, err := os.Open(fp)
	if err != nil {
		logger.Warning("Unable to open default connections file from %v\n%v",
			fp, err)
		return connections
	}
	defer fo.Close()

	input := bufio.NewScanner(fo)
	for input.Scan() {
		strAddress := input.Text()
		// TODO: Validate addresses
		connections = append(connections, strAddress)
	}
	return connections
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
	dc.Visor.Config.GenesisCoinVolume = c.GenesisCoinVolume
	dc.Visor.Config.DBPath = c.DBPath
	dc.Visor.Config.Arbitrating = c.Arbitrating
	dc.Visor.Config.WalletDirectory = c.WalletDirectory
	dc.Visor.Config.BuildInfo = visor.BuildInfo{
		Version: Version,
		Commit:  Commit,
	}

	dc.Gateway.DisableWalletAPI = c.DisableWalletApi

	return dc
}

// Run starts the skycoin node
func Run(c *Config) {
	defer func() {
		// try catch panic in main thread
		if r := recover(); r != nil {
			logger.Error("recover: %v\nstack:%v", r, string(debug.Stack()))
		}
	}()

	c.GUIDirectory = file.ResolveResourceDirectory(c.GUIDirectory)

	scheme := "http"
	if c.WebInterfaceHTTPS {
		scheme = "https"
	}
	host := fmt.Sprintf("%s:%d", c.WebInterfaceAddr, c.WebInterfacePort)
	fullAddress := fmt.Sprintf("%s://%s", scheme, host)
	logger.Critical("Full address: %s", fullAddress)
	if c.PrintWebInterfaceAddress {
		fmt.Println(fullAddress)
	}

	initProfiling(c.HTTPProf, c.ProfileCPU, c.ProfileCPUFile)

	closelog, err := initLogging(c.DataDirectory, c.LogLevel, c.ColorLog, c.LogFmt, c.Logtofile)
	if err != nil {
		fmt.Println(err)
		return
	}

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

	db, err := visor.OpenDB(dconf.Visor.Config.DBPath)
	if err != nil {
		logger.Error("Database failed to open: %v. Is another skycoin instance running?", err)
		return
	}

	d, err := daemon.NewDaemon(dconf, db, c.DefaultConnections)
	if err != nil {
		logger.Error("%v", err)
		return
	}

	var rpc *webrpc.WebRPC
	// start the webrpc
	if c.RPCInterface {
		rpcAddr := fmt.Sprintf("%v:%v", c.RPCInterfaceAddr, c.RPCInterfacePort)
		rpc, err = webrpc.New(rpcAddr, d.Gateway)
		if err != nil {
			logger.Error("%v", err)
			return
		}
		rpc.ChanBuffSize = 1000
		rpc.WorkerNum = c.RPCThreadNum
	}

	var webInterface *gui.Server
	if c.WebInterface {
		webInterface, err = createGUI(c, d, host, quit)
		if err != nil {
			logger.Error("%v", err)
			return
		}
	}

	// Debug only - forces connection on start.  Violates thread safety.
	if c.ConnectTo != "" {
		if err := d.Pool.Pool.Connect(c.ConnectTo); err != nil {
			logger.Error("Force connect %s failed, %v", c.ConnectTo, err)
			return
		}
	}

	// POTENTIALLY UNSAFE CODE -- See https://github.com/skycoin/skycoin/issues/838
	// closelog, err := initLogging(c.DataDirectory, c.LogLevel, c.ColorLog, c.Logtofile, c.Logtogui, &d.LogBuff)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// if c.Logtogui {
	// 	go func(buf *bytes.Buffer, quit chan struct{}) {
	// 		for {
	// 			select {
	// 			case <-quit:
	// 				logger.Info("Logbuff service closed normally")
	// 				return
	// 			case <-time.After(1 * time.Second): //insure logbuff size not exceed required size, like lru
	// 				for buf.Len() > c.LogBuffSize {
	// 					_, err := buf.ReadString(byte('\n')) //discard one line
	// 					if err != nil {
	// 						continue
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}(&d.LogBuff, quit)
	// }

	errC := make(chan error, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := d.Run(); err != nil {
			logger.Error("%v", err)
			errC <- err
		}
	}()

	// start the webrpc
	if c.RPCInterface {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := rpc.Run(); err != nil {
				logger.Error("%v", err)
				errC <- err
			}
		}()
	}

	if c.WebInterface {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := webInterface.Serve(); err != nil {
				logger.Error("%v", err)
				errC <- err
			}
		}()

		if c.LaunchBrowser {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Wait a moment just to make sure the http interface is up
				time.Sleep(time.Millisecond * 100)

				logger.Info("Launching System Browser with %s", fullAddress)
				if err := browser.Open(fullAddress); err != nil {
					logger.Error(err.Error())
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
		logger.Error("%v", err)
	}

	logger.Info("Shutting down...")
	if rpc != nil {
		rpc.Shutdown()
	}
	if webInterface != nil {
		webInterface.Shutdown()
	}
	d.Shutdown()
	closelog()
	wg.Wait()
	logger.Info("Goodbye")
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

// ChainConfig blockchain config info
type ChainConfig struct {
	// GenesisSignature genesis signature
	GenesisSignature string
	// GenesisAddressStr genesis address
	GenesisAddress string
	// BlockchainPubkeyStr blockchain pubkey
	BlockchainPubkey string
	// BlockchainSeckey blockchain seckey
	BlockchainSeckey string
	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64
	// Port node port
	Port int
	// Web interface port http api service port
	WebInterfacePort int
	// RPC interface port
	RPCInterfacePort int
	// Data directory
	DataDirectory string
	// DefaultConnections the default trust node addresses
	DefaultConnections []string
	// LogFmt log format
	LogFmt string
}

// MainChainCfg main chain config info
var MainChainCfg = ChainConfig{
	GenesisSignature:  "eb10468d10054d15f2b6f8946cd46797779aa20a7617ceb4be884189f219bc9a164e56a5b9f7bec392a804ff3740210348d73db77a37adb542a8e08d429ac92700",
	GenesisAddress:    "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6",
	BlockchainPubkey:  "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a",
	BlockchainSeckey:  "",
	GenesisTimestamp:  1426562704,
	GenesisCoinVolume: 100e12,
	Port:              6000,
	WebInterfacePort:  6420,
	RPCInterfacePort:  6430,
	DataDirectory:     "~/.skycoin",
	LogFmt:            "[skycoin.%{module}:%{level}] %{message}",
	DefaultConnections: []string{
		"118.178.135.93:6000",
		"47.88.33.156:6000",
		"121.41.103.148:6000",
		"120.77.69.188:6000",
		"104.237.142.206:6000",
		"176.58.126.224:6000",
		"172.104.85.6:6000",
		"139.162.7.132:6000",
	},
}

// TestChainCfg test chain config info
var TestChainCfg = ChainConfig{
	GenesisSignature:  "07f46ce7502147a97f2fb32c7c1e66638af851c1cb532d893f1f360bb4ab1ccf0656f2f358695e8cb752e05080af69c8f44b0d72610bd11e3fb028ecdcfed2ea01",
	GenesisAddress:    "F5k1VyFHZGJgQADWpmMEW8Se2HNidFm9k3",
	BlockchainPubkey:  "03b2595c36f542bf4d3cf347327fef1e21cbe0600c281efed5f673eb0c77298e4c",
	GenesisTimestamp:  1505801448,
	GenesisCoinVolume: 100e12,
	Port:              16000,
	WebInterfacePort:  16420,
	RPCInterfacePort:  16430,
	DataDirectory:     "~/.skycoin-testnet",
	LogFmt:            "[skycoin.testnet.%{module}:%{level}] %{message}",
	DefaultConnections: []string{
		"139.162.33.154:16000",
	},
}
