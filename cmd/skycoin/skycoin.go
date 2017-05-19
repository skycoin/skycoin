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
	"runtime/pprof"
	"syscall"
	"time"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util"
)

var (
	logger     = util.MustGetLogger("main")
	logFormat  = "[skycoin.%{module}:%{level}] %{message}"
	logModules = []string{
		"main",
		"daemon",
		"coin",
		"gui",
		"util",
		"visor",
		"wallet",
		"gnet",
		"pex",
		"webrpc",
	}

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
	}
)

// Command line interface arguments

// Config records the node's configuration
type Config struct {
	// Disable peer exchange
	DisablePEX bool
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	//gnet uses this for TCP incoming and outgoing
	Port int
	//max connections to maintain
	MaxConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
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
	LogLevel logging.Level
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	logLevel string

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string

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
	Arbitrating bool
}

func (c *Config) register() {
	flag.BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX,
		"disable PEX peer discovery")
	flag.BoolVar(&c.DisableOutgoingConnections, "disable-outgoing",
		c.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.DisableIncomingConnections, "disable-incoming",
		c.DisableIncomingConnections, "Don't make incoming connections")
	flag.BoolVar(&c.DisableNetworking, "disable-networking",
		c.DisableNetworking, "Disable all network activity")
	flag.StringVar(&c.Address, "address", c.Address,
		"IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")
	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface,
		"enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port",
		c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr",
		c.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterfaceCert, "web-interface-cert",
		c.WebInterfaceCert, "cert.pem file for web interface HTTPS. "+
			"If not provided, will use cert.pem in -data-directory")
	flag.StringVar(&c.WebInterfaceKey, "web-interface-key",
		c.WebInterfaceKey, "key.pem file for web interface HTTPS. "+
			"If not provided, will use key.pem in -data-directory")
	flag.BoolVar(&c.WebInterfaceHTTPS, "web-interface-https",
		c.WebInterfaceHTTPS, "enable HTTPS for web interface")

	flag.BoolVar(&c.RPCInterface, "rpc-interface", c.RPCInterface,
		"enable the rpc interface")
	flag.IntVar(&c.RPCInterfacePort, "rpc-interface-port", c.RPCInterfacePort,
		"port to serve rpc interface on")
	flag.StringVar(&c.RPCInterfaceAddr, "rpc-interface-addr", c.RPCInterfaceAddr,
		"addr to serve rpc interface on")

	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser,
		"launch system default webbrowser at client startup")
	flag.BoolVar(&c.PrintWebInterfaceAddress, "print-web-interface-address",
		c.PrintWebInterfaceAddress, "print configured web interface address and exit")
	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory,
		"directory to store app data (defaults to ~/.skycoin)")
	flag.StringVar(&c.ConnectTo, "connect-to", c.ConnectTo,
		"connect to this ip only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU,
		"enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file",
		c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf,
		"Run the http profiling interface")
	flag.StringVar(&c.logLevel, "log-level", c.logLevel,
		"Choices are: debug, info, notice, warning, error, critical")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog,
		"Add terminal colors to log output")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory,
		"static content directory for the html gui")

	//Key Configuration Data
	flag.BoolVar(&c.RunMaster, "master", c.RunMaster,
		"run the daemon as blockchain master server")

	flag.StringVar(&BlockchainPubkeyStr, "master-public-key", BlockchainPubkeyStr,
		"public key of the master chain")
	flag.StringVar(&BlockchainSeckeyStr, "master-secret-key", BlockchainSeckeyStr,
		"secret key, set for master")

	flag.StringVar(&GenesisAddressStr, "genesis-address", GenesisAddressStr,
		"genesis address")
	flag.StringVar(&GenesisSignatureStr, "genesis-signature", GenesisSignatureStr,
		"genesis block signature")
	flag.Uint64Var(&c.GenesisTimestamp, "genesis-timestamp", c.GenesisTimestamp,
		"genesis block timestamp")

	flag.StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory,
		"location of the wallet files. Defaults to ~/.skycoin/wallet/")

	flag.DurationVar(&c.OutgoingConnectionsRate, "connection-rate",
		c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly,
		"Run on localhost and only connect to localhost peers")
	flag.BoolVar(&c.Arbitrating, "arbitrating", c.Arbitrating, "Run node in arbitrating mode")
	//flag.StringVar(&c.AddressVersion, "address-version", c.AddressVersion,
	//	"Wallet address version. Options are 'test' and 'main'")
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
	// Only run on localhost and only connect to others on localhost
	LocalhostOnly: false,
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address: "",
	//gnet uses this for TCP incoming and outgoing
	Port: 6000,

	MaxConnections: 16,
	// How often to make outgoing connections, in seconds
	OutgoingConnectionsRate: time.Second * 5,
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

	RPCInterface:     true,
	RPCInterfacePort: 6430,
	RPCInterfaceAddr: "127.0.0.1",

	LaunchBrowser: true,
	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory: ".skycoin",
	// Web GUI static resources
	GUIDirectory: "./src/gui/static/",
	// Logging
	LogLevel: logging.DEBUG,
	ColorLog: true,
	logLevel: "DEBUG",

	// Wallets
	WalletDirectory: "",

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

// Parse prepare the config
func (c *Config) Parse() {
	c.register()
	flag.Parse()
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

	c.DataDirectory = util.InitDataDir(c.DataDirectory)
	if c.WebInterfaceCert == "" {
		c.WebInterfaceCert = filepath.Join(c.DataDirectory, "cert.pem")
	}
	if c.WebInterfaceKey == "" {
		c.WebInterfaceKey = filepath.Join(c.DataDirectory, "key.pem")
	}

	if c.WalletDirectory == "" {
		c.WalletDirectory = filepath.Join(c.DataDirectory, "wallets/")
	}

	ll, err := logging.LogLevel(c.logLevel)
	panicIfError(err, "Invalid -log-level %s", c.logLevel)
	c.LogLevel = ll

	if c.DBPath == "" {
		c.DBPath = filepath.Join(c.DataDirectory, "data.db")
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

func catchInterrupt(quit chan<- int) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	quit <- 1
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

func initLogging(level logging.Level, color bool) {
	format := logging.MustStringFormatter(logFormat)
	logging.SetFormatter(format)
	for _, s := range logModules {
		logging.SetLevel(level, s)
	}
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	stdout.Color = color
	logging.SetBackend(stdout)
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
	dc.Peers.DataDirectory = c.DataDirectory
	dc.Peers.Disabled = c.DisablePEX
	dc.Daemon.DisableOutgoingConnections = c.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = c.DisableIncomingConnections
	dc.Daemon.DisableNetworking = c.DisableNetworking
	dc.Daemon.Port = c.Port
	dc.Daemon.Address = c.Address
	dc.Daemon.LocalhostOnly = c.LocalhostOnly
	dc.Daemon.OutgoingMax = c.MaxConnections

	daemon.DefaultConnections = DefaultConnections

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
	return dc
}

// Run starts the skycoin node
func Run(c *Config) {

	c.GUIDirectory = util.ResolveResourceDirectory(c.GUIDirectory)

	scheme := "http"
	if c.WebInterfaceHTTPS {
		scheme = "https"
	}
	host := fmt.Sprintf("%s:%d", c.WebInterfaceAddr, c.WebInterfacePort)
	fullAddress := fmt.Sprintf("%s://%s", scheme, host)
	logger.Critical("Full address: %s", fullAddress)

	if c.PrintWebInterfaceAddress {
		fmt.Println(fullAddress)
		return
	}

	initProfiling(c.HTTPProf, c.ProfileCPU, c.ProfileCPUFile)

	logCfg := util.DevLogConfig(logModules)
	logCfg.Format = logFormat
	logCfg.Colors = c.ColorLog
	logCfg.InitLogger()

	// If the user Ctrl-C's, shutdown properly
	quit := make(chan int)
	go catchInterrupt(quit)
	// Watch for SIGUSR1
	go catchDebug()

	gui.InitWalletRPC(c.WalletDirectory)

	dconf := configureDaemon(c)
	d := daemon.NewDaemon(dconf)

	stopDaemon := make(chan int)
	go d.Start(stopDaemon)

	// start the webrpc
	closingC := make(chan struct{})
	if c.RPCInterface {
		go webrpc.Start(
			fmt.Sprintf("%v:%v", c.RPCInterfaceAddr, c.RPCInterfacePort),
			webrpc.ChanBuffSize(1000),
			webrpc.ThreadNum(1000),
			webrpc.Gateway(d.Gateway),
			webrpc.Quit(closingC))
	}

	// Debug only - forces connection on start.  Violates thread safety.
	if c.ConnectTo != "" {
		if err := d.Pool.Pool.Connect(c.ConnectTo); err != nil {
			log.Panic(err)
		}
	}

	if c.WebInterface {
		var err error
		if c.WebInterfaceHTTPS {
			// Verify cert/key parameters, and if neither exist, create them
			errs := util.CreateCertIfNotExists(host, c.WebInterfaceCert, c.WebInterfaceKey, "Skycoind")
			if len(errs) != 0 {
				for _, err := range errs {
					logger.Error(err.Error())
				}
				logger.Error("gui.CreateCertIfNotExists failure")
				os.Exit(1)
			}

			err = gui.LaunchWebInterfaceHTTPS(host, c.GUIDirectory, d, c.WebInterfaceCert, c.WebInterfaceKey)
		} else {
			err = gui.LaunchWebInterface(host, c.GUIDirectory, d)
		}

		if err != nil {
			logger.Error(err.Error())
			logger.Error("Failed to start web GUI")
			os.Exit(1)
		}

		if c.LaunchBrowser {
			go func() {
				// Wait a moment just to make sure the http interface is up
				time.Sleep(time.Millisecond * 100)

				logger.Info("Launching System Browser with %s", fullAddress)
				if err := util.OpenBrowser(fullAddress); err != nil {
					logger.Error(err.Error())
				}
			}()
		}
	}

	/*
		time.Sleep(5)
		tx := InitTransaction()
		_ = tx
		err, _ = d.Visor.Visor.InjectTxn(tx)
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
					err, _ := d.Visor.Visor.InjectTxn(tx)
					if err != nil {
						//log.Panic(err)
					}
				}
			}()
		}
	*/

	<-quit
	stopDaemon <- 1

	logger.Info("Shutting down")
	gui.Shutdown()
	close(closingC)

	d.Shutdown()
	logger.Info("Goodbye")
}

func main() {
	devConfig.Parse()
	Run(&devConfig)
}

// AddrList for storage of coins
var AddrList = []string{
	"R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
	"2EYM4WFHe4Dgz6kjAdUkM6Etep7ruz2ia6h",
	"25aGyzypSA3T9K6rgPUv1ouR13efNPtWP5m",
	"ix44h3cojvN6nqGcdpy62X7Rw6Ahnr3Thk",
	"AYV8KEBEAPCg8a59cHgqHMqYHP9nVgQDyW",
	"2Nu5Jv5Wp3RYGJU1EkjWFFHnebxMx1GjfkF",
	"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
	"tWZ11Nvor9parjg4FkwxNVcby59WVTw2iL",
	"m2joQiJRZnj3jN6NsoKNxaxzUTijkdRoSR",
	"8yf8PAQqU2cDj8Yzgz3LgBEyDqjvCh2xR7",
	"sgB3n11ZPUYHToju6TWMpUZTUcKvQnoFMJ",
	"2UYPbDBnHUEc67e7qD4eXtQQ6zfU2cyvAvk",
	"wybwGC9rhm8ZssBuzpy5goXrAdE31MPdsj",
	"JbM25o7kY7hqJZt3WGYu9pHZFCpA9TCR6t",
	"2efrft5Lnwjtk7F1p9d7BnPd72zko2hQWNi",
	"Syzmb3MiMoiNVpqFdQ38hWgffHg86D2J4e",
	"2g3GUmTQooLrNHaRDhKtLU8rWLz36Beow7F",
	"D3phtGr9iv6238b3zYXq6VgwrzwvfRzWZQ",
	"gpqsFSuMCZmsjPc6Rtgy1FmLx424tH86My",
	"2EUF3GPEUmfocnUc1w6YPtqXVCy3UZA4rAq",
	"TtAaxB3qGz5zEAhhiGkBY9VPV7cekhvRYS",
	"2fM5gVpi7XaiMPm4i29zddTNkmrKe6TzhVZ",
	"ix3NDKgxfYYANKAb5kbmwBYXPrkAsha7uG",
	"2RkPshpFFrkuaP98GprLtgHFTGvPY5e6wCK",
	"Ak1qCDNudRxZVvcW6YDAdD9jpYNNStAVqm",
	"2eZYSbzBKJ7QCL4kd5LSqV478rJQGb4UNkf",
	"KPfqM6S96WtRLMuSy4XLfVwymVqivdcDoM",
	"5B98bU1nsedGJBdRD5wLtq7Z8t8ZXio8u5",
	"2iZWk5tmBynWxj2PpAFyiZzEws9qSnG3a6n",
	"XUGdPaVnMh7jtzPe3zkrf9FKh5nztFnQU5",
	"hSNgHgewJme8uaHrEuKubHYtYSDckD6hpf",
	"2DeK765jLgnMweYrMp1NaYHfzxumfR1PaQN",
	"orrAssY5V2HuQAbW9K6WktFrGieq2m23pr",
	"4Ebf4PkG9QEnQTm4MVvaZvJV6Y9av3jhgb",
	"7Uf5xJ3GkiEKaLxC2WmJ1t6SeekJeBdJfu",
	"oz4ytDKbCqpgjW3LPc52pW2CaK2gxCcWmL",
	"2ex5Z7TufQ5Z8xv5mXe53fSQRfUr35SSo7Q",
	"WV2ap7ZubTxeDdmEZ1Xo7ufGMkekLWikJu",
	"ckCTV4r1pNuz6j2VBRHhaJN9HsCLY7muLV",
	"MXJx96ZJVSjktgeYZpVK8vn1H3xWP8ooq5",
	"wyQVmno9aBJZmQ99nDSLoYWwp7YDJCWsrH",
	"2cc9wKxCsFNRkoAQDAoHke3ZoyL1mSV14cj",
	"29k9g3F5AYfVaa1joE1PpZjBED6hQXes8Mm",
	"2XPLzz4ZLf1A9ykyTCjW5gEmVjnWa8CuatH",
	"iH7DqqojTgUn2JxmY9hgFp165Nk7wKfan9",
	"RJzzwUs3c9C8Y7NFYzNfFoqiUKeBhBfPki",
	"2W2cGyiCRM4nwmmiGPgMuGaPGeBzEm7VZPn",
	"ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
	"tBaeg9zE2sgmw5ZQENaPPYd6jfwpVpGTzS",
	"2hdTw5Hk3rsgpZjvk8TyKcCZoRVXU5QVrUt",
	"A1QU6jKq8YgTP79M8fwZNHUZc7hConFKmy",
	"q9RkXoty3X1fuaypDDRUi78rWgJWYJMmpJ",
	"2Xvm6is5cAPA85xnSYXDuAqiRyoXiky5RaD",
	"4CW2CPJEzxhn2PS4JoSLoWGL5QQ7dL2eji",
	"24EG6uTzL7DHNzcwsygYGRR1nfu5kco7AZ1",
	"KghGnWw5fppTrqHSERXZf61yf7GkuQdCnV",
	"2WojewRA3LbpyXTP9ANy8CZqJMgmyNm3MDr",
	"2BsMfywmGV3M2CoDA112Rs7ZBkiMHfy9X11",
	"kK1Q4gPyYfVVMzQtAPRzL8qXMqJ67Y7tKs",
	"28J4mx8xfUtM92DbQ6i2Jmqw5J7dNivfroN",
	"gQvgyG1djgtftoCVrSZmsRxr7okD4LheKw",
	"3iFGBKapAWWzbiGFSr5ScbhrEPm6Esyvia",
	"NFW2akQH2vu7AqkQXxFz2P5vkXTWkSqrSm",
	"2MQJjLnWRp9eHh6MpCwpiUeshhtmri12mci",
	"2QjRQUMyL6iodtHP9zKmxCNYZ7k3jxtk49C",
	"USdfKy7B6oFNoauHWMmoCA7ND9rHqYw2Mf",
	"cA49et9WtptYHf6wA1F8qqVgH3kS5jJ9vK",
	"qaJT9TjcMi46sTKcgwRQU8o5Lw2Ea1gC4N",
	"22pyn5RyhqtTQu4obYjuWYRNNw4i54L8xVr",
	"22dkmukC6iH4FFLBmHne6modJZZQ3MC9BAT",
	"z6CJZfYLvmd41GRVE8HASjRcy5hqbpHZvE",
	"GEBWJ2KpRQDBTCCtvnaAJV2cYurgXS8pta",
	"oS8fbEm82cprmAeineBeDkaKd7QownDZQh",
	"rQpAs1LVQdphyj9ipEAuukAoj9kNpSP8cM",
	"6NSJKsPxmqipGAfFFhUKbkopjrvEESTX3j",
	"cuC68ycVXmD2EBzYFNYQ6akhKGrh3FGjSf",
	"bw4wtYU8toepomrhWP2p8UFYfHBbvEV425",
	"HvgNmDz5jD39Gwmi9VfDY1iYMhZUpZ8GKz",
	"SbApuZAYquWP3Q6iD51BcMBQjuApYEkRVf",
	"2Ugii5yxJgLzC59jV1vF8GK7UBZdvxwobeJ",
	"21N2iJ1qnQRiJWcEqNRxXwfNp8QcmiyhtPy",
	"9TC4RGs6AtFUsbcVWnSoCdoCpSfM66ALAc",
	"oQzn55UWG4iMcY9bTNb27aTnRdfiGHAwbD",
	"2GCdwsRpQhcf8SQcynFrMVDM26Bbj6sgv9M",
	"2NRFe7REtSmaM2qAgZeG45hC8EtVGV2QjeB",
	"25RGnhN7VojHUTvQBJA9nBT5y1qTQGULMzR",
	"26uCBDfF8E2PJU2Dzz2ysgKwv9m4BhodTz9",
	"Wkvima5cF7DDFdmJQqcdq8Syaq9DuAJJRD",
	"286hSoJYxvENFSHwG51ZbmKaochLJyq4ERQ",
	"FEGxF3HPoM2HCWHn82tyeh9o7vEQq5ySGE",
	"h38DxNxGhWGTq9p5tJnN5r4Fwnn85Krrb6",
	"2c1UU8J6Y3kL4cmQh21Tj8wkzidCiZxwdwd",
	"2bJ32KuGmjmwKyAtzWdLFpXNM6t83CCPLq5",
	"2fi8oLC9zfVVGnzzQtu3Y3rffS65Hiz6QHo",
	"TKD93RxFr2Am44TntLiJQus4qcEwTtvEEQ",
	"zMDywYdGEDtTSvWnCyc3qsYHWwj9ogws74",
	"25NbotTka7TwtbXUpSCQD8RMgHKspyDubXJ",
	"2ayCELBERubQWH5QxUr3cTxrYpidvUAzsSw",
	"RMTCwLiYDKEAiJu5ekHL1NQ8UKHi5ozCPg",
	"ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm",
}

// InitTransaction creates the initialize transaction
func InitTransaction() coin.Transaction {
	var tx coin.Transaction

	output := cipher.MustSHA256FromHex("043836eb6f29aaeb8b9bfce847e07c159c72b25ae17d291f32125e7f1912e2a0")
	tx.PushInput(output)

	for i := 0; i < 100; i++ {
		addr := cipher.MustDecodeBase58Address(AddrList[i])
		tx.PushOutput(addr, 1e12, 1) // 10e6*10e6
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
