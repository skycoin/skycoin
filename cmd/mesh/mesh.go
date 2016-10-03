package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin/src/mesh/gui"
	"github.com/skycoin/skycoin/src/mesh/nodemanager"
	"github.com/skycoin/skycoin/src/util"
)

const (
	logModuleName = "skywire.main"

	// logFormat is same as src/util/logging.go#defaultLogFormat
	// logFormat     = "[%{module}:%{level}] %{message}"
)

var (
	logger     = util.MustGetLogger(logModuleName)
	logModules = []string{logModuleName}
)

// default configurations
const (
	webInterfaceEnable = true
	webInterfacePort   = 6480
	webInterfaceAddr   = "127.0.0.1"
	webInterfaceHttps  = false
	launchBrowser      = true
	guiDirectory       = "./src/mesh/gui/static/"
	dataDirectory      = ".skywire"
)

// TODO:
//    Q: to move WebInterfaceConfig to its related package
//       (github.com/skycoin/skycoin/src/mesh/gui) ?
// Remote web interface
type WebInterfaceConfig struct {
	Enable bool
	Port   int
	Addr   string
	Cert   string
	Key    string
	HTTPS  bool
	// Launch system default browser after client startup
	LaunchBrowser bool
	// static htmls + assets location
	GUIDirectory string
}

type Config struct {
	// WebInterface configs
	WebInterface WebInterfaceConfig
	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string
	// logging configs
	Log *util.LogConfig
}

// TODO: defaultConfig or may be (developConfig + productConfig) ?

func defaultConfig() *Config {
	cfg := &Config{
		WebInterface: WebInterfaceConfig{
			Enable: webInterfaceEnable,
			Port:   webInterfacePort,
			Addr:   webInterfaceAddr,
			// Cert: "",
			// Key: "",
			HTTPS:         webInterfaceHttps,
			LaunchBrowser: launchBrowser,
			GUIDirectory:  guiDirectory,
		},
		// Data directory holds app data -- defaults to ~/.skycoin
		DataDirectory: dataDirectory,
		// TODO: dev/prod vs default?
		//       see src/util/logging.go 'TODO' note before DevLogConfig()
		//       for details
		Log: util.DevLogConfig(logModules),
	}
	return cfg
}

func (c *Config) Parse() {
	// obtain values from flags
	c.fromFlags()
	// post process
	c.DataDirectory = util.InitDataDir(c.DataDirectory)
	// if HTTPS is turned off then cerk/key are never used
	if c.WebInterface.HTTPS == true {
		if c.WebInterface.Cert == "" {
			c.WebInterface.Cert = filepath.Join(c.DataDirectory, "cert.pem")
		}
		if c.WebInterface.Key == "" {
			c.WebInterface.Key = filepath.Join(c.DataDirectory, "key.pem")
		}
	}

	// initialize logger
	c.Log.InitLogger()
}

func (c *Config) fromFlags() {
	flag.BoolVar(&c.WebInterface.Enable, "web-interface", c.WebInterface.Enable,
		"enable the web interface")
	flag.IntVar(&c.WebInterface.Port, "web-interface-port",
		c.WebInterface.Port, "port to serve web interface on")
	flag.StringVar(&c.WebInterface.Addr, "web-interface-addr",
		c.WebInterface.Addr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterface.Cert, "web-interface-cert",
		c.WebInterface.Cert, "cert.pem file for web interface HTTPS. "+
			"If not provided, will use cert.pem in -data-directory")
	flag.StringVar(&c.WebInterface.Key, "web-interface-key",
		c.WebInterface.Key, "key.pem file for web interface HTTPS. "+
			"If not provided, will use key.pem in -data-directory")
	flag.BoolVar(&c.WebInterface.HTTPS, "web-interface-https",
		c.WebInterface.HTTPS, "enable HTTPS for web interface")
	flag.BoolVar(&c.WebInterface.LaunchBrowser, "launch-browser",
		c.WebInterface.LaunchBrowser,
		"launch system default webbrowser at client startup")

	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory,
		"directory to store app data (defaults to ~/.skycoin)")

	flag.StringVar(&c.WebInterface.GUIDirectory, "gui-dir",
		c.WebInterface.GUIDirectory,
		"static content directory for the html gui")

	flag.StringVar(&c.Log.Level, "log-level", c.Log.Level,
		"Choices are: debug, info, notice, warning, error, critical")
	flag.BoolVar(&c.Log.Colors, "color-log", c.Log.Colors,
		"Add terminal colors to log output")

	flag.Parse()
}

// returns "http" ot "https";
// result depends on c.WebInterface.HTTPS
func (c *Config) scheme() string {
	if c.WebInterface.HTTPS == true {
		return "https"
	}
	return "http"
}

// returns c.WebInterfaceAddr:c.WebInterface.Port string
func (c *Config) host() string {
	return fmt.Sprintf("%s:%d", c.WebInterface.Addr, c.WebInterface.Port)
}

// returns full address: scheme://address:port
func (c *Config) fullAddress() string {
	return fmt.Sprintf("%s://%s", c.scheme(), c.host())
}

// launches browser if it's enabled
func (c *Config) launchBrowser() {
	if c.WebInterface.LaunchBrowser == true {
		go func() {
			// TODO: wait is BS, is it really needed?
			//
			// Wait a moment just to make sure the http interface is up
			time.Sleep(time.Millisecond * 100)
			//

			logger.Info("Launching System Browser with %s", c.fullAddress())
			if err := util.OpenBrowser(c.fullAddress()); err != nil {
				logger.Error(err.Error())
			}
		}()
	}
}

// subscribe to SIGINT
func catchInterrupt() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	return sig
}

func Run(c *Config) {
	c.WebInterface.GUIDirectory = util.ResolveResourceDirectory(
		c.WebInterface.GUIDirectory)

	// print this message regardless of the log level
	fmt.Printf("Full address: %s\n", c.fullAddress())

	// start node_manager
	fmt.Printf("Starting Node Manager Service...\n")
	// TODO: empty config? what the hell?
	// *nodemanager.NodeManager
	nm := nodemanager.NewNodeManager(&nodemanager.NodeManagerConfig{})

	if c.WebInterface.Enable == true {
		var err error
		if c.WebInterface.HTTPS == true {
			// TODO

			log.Panic("HTTPS support is not implemented yet")

			//
			// errs := util.CreateCertIfNotExists(host, c.WebInterfaceCert,
			//                                   c.WebInterfaceKey, "Skywired")
			// if len(errs) != 0 {
			// 	for _, err := range errs {
			// 		logger.Error(err.Error())
			// 	}
			// 	logger.Error("gui.CreateCertIfNotExists failure")
			// 	os.Exit(1)
			// }
			//
			// err = gui.LaunchWebInterfaceHTTPS(host, c.GUIDirectory, d,
			//                           c.WebInterfaceCert, c.WebInterfaceKey)
			//
		} else {
			err = gui.LaunchWebInterface(c.host(), c.WebInterface.GUIDirectory,
				nm)
		}

		if err != nil {
			logger.Error(err.Error())
			logger.Error("Failed to start web GUI")
			os.Exit(1)
		}

		c.launchBrowser()
	}

	// subscribe to SIGINT (Ctrl+C)
	sigint := catchInterrupt()
	// start the node manager
	//     don't use 'chan int' to stop the node manager;
	//     there is Shutdown() method
	go nm.Start()

	// waiting for SIGINT (Ctrl+C)
	logger.Info("Got signal %q, shutting down...", <-sigint)

	// shutdown the node manager
	nm.Shutdown()
	logger.Info("Goodbye")
}

func main() {
	cfg := defaultConfig()
	cfg.Parse()
	Run(cfg)
}

//
// TODO: what is the stuff below ?
//

// Q: to omit or not to omit?

/*

func main_old() {
	fmt.Fprintln(os.Stdout, "Starting Node Manager Service...")

	config := nodemanager.ServerConfig
	fmt.Fprintln(os.Stdout, "PubKey:", config.Node.PubKey)
	//fmt.Fprintln(os.Stdout, "ChaCha20Key:", config.Node.ChaCha20Key)
	fmt.Fprintln(os.Stdout, "Port:", config.Udp.ListenPortMin)
	node := nodemanager.CreateNode(*config)

	node.AddTransportToNode(*config)

	received := make(chan mesh.MeshMessage, 10)
	node.SetReceiveChannel(received)

	isActiveService := true

	for isActiveService {
		select {
		case msgRecvd, ok := <-received:
			{
				if ok {
					fmt.Fprintf(os.Stdout,
						"Message received: %v\nReplyTo: %+v\n",
						string(msgRecvd.Contents), msgRecvd.ReplyTo)
					go filterMessages(msgRecvd)
				}
			}
		}
	}

}

func filterMessages(msg mesh.MeshMessage) bool {
	v, err := connection.ConnectionManager.DeserializeMessage(msg.Contents)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	msg_type := reflect.TypeOf(v)

	fmt.Fprintln(os.Stdout, "msg-type", msg_type)

	if msg_type == reflect.TypeOf(domain.AddNodeMessage{}) {
		addNodeMsg := v.(domain.AddNodeMessage)

		config := mesh.TestConfig{}
		err := json.Unmarshal(addNodeMsg.Content, &config)
		if err != nil {
			return false
		}
		fmt.Fprintf(os.Stdout, "TestConfig %+v\n", config)

		return true
	}
	return false
}

*/
