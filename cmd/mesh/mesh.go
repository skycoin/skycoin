package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"time"

	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/gui"
	mesh "github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/node/connection"
	"github.com/skycoin/skycoin/src/mesh/nodemanager"
	"github.com/skycoin/skycoin/src/util"

	"gopkg.in/op/go-logging.v1"
)

var (
	logger     = logging.MustGetLogger("skywire.main")
	logFormat  = "[%{module}:%{level}] %{message}"
	logModules = []string{
		"skywire.main",
	}
)

type Config struct {
	// Remote web interface
	WebInterface      bool
	WebInterfacePort  int
	WebInterfaceAddr  string
	WebInterfaceCert  string
	WebInterfaceKey   string
	WebInterfaceHTTPS bool

	// Launch System Default Browser after client startup
	LaunchBrowser bool
	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string
	// GUI directory contains assets for the html gui
	GUIDirectory string

	// Logging
	LogLevel logging.Level
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	logLevel string
}

func (c *Config) Parse() {
	c.register()
	flag.Parse()
	c.postProcess()
}

func (c *Config) register() {

	//flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")
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
	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser,
		"launch system default webbrowser at client startup")

	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory,
		"directory to store app data (defaults to ~/.skycoin)")

	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory,
		"static content directory for the html gui")

	flag.StringVar(&c.logLevel, "log-level", c.logLevel,
		"Choices are: debug, info, notice, warning, error, critical")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog,
		"Add terminal colors to log output")
}

func (c *Config) postProcess() {

	c.DataDirectory = util.InitDataDir(c.DataDirectory)
	if c.WebInterfaceCert == "" {
		c.WebInterfaceCert = filepath.Join(c.DataDirectory, "cert.pem")
	}
	if c.WebInterfaceKey == "" {
		c.WebInterfaceKey = filepath.Join(c.DataDirectory, "key.pem")
	}

	ll, err := logging.LogLevel(c.logLevel)
	panicIfError(err, "Invalid -log-level %s", c.logLevel)
	c.LogLevel = ll

}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

var devConfig Config = Config{
	// Remote web interface
	WebInterface:      true,
	WebInterfacePort:  6480,
	WebInterfaceAddr:  "127.0.0.1",
	WebInterfaceCert:  "",
	WebInterfaceKey:   "",
	WebInterfaceHTTPS: false,

	LaunchBrowser: true,

	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory: ".skywire",
	// Web GUI static resources
	GUIDirectory: "./src/mesh/gui/static/",

	// Logging
	LogLevel: logging.DEBUG,
	ColorLog: true,
	logLevel: "DEBUG",
}

//move to util
func catchInterrupt(quit chan<- int) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	quit <- 1
}

//move to util
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

func Run(c *Config) {

	c.GUIDirectory = util.ResolveResourceDirectory(c.GUIDirectory)

	// If the user Ctrl-C's, shutdown properly
	quit := make(chan int)
	go catchInterrupt(quit)

	scheme := "http"
	host := fmt.Sprintf("%s:%d", c.WebInterfaceAddr, c.WebInterfacePort)
	fullAddress := fmt.Sprintf("%s://%s", scheme, host)
	logger.Critical("Full address: %s", fullAddress)

	initLogging(c.LogLevel, c.ColorLog)

	//state node_manager
	fmt.Fprintln(os.Stdout, "Starting Node Manager Service...")
	var nm_config = nodemanager.NodeManagerConfig{}
	var nm *nodemanager.NodeManager = nodemanager.NewNodeManager(&nm_config)

	if c.WebInterface {
		var err error

		/*
			if c.WebInterfaceHTTPS {
				// Verify cert/key parameters, and if neither exist, create them
				errs := util.CreateCertIfNotExists(host, c.WebInterfaceCert, c.WebInterfaceKey, "Skywired")
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
		*/

		//pass the node manager instance to the http server
		//no https
		err = gui.LaunchWebInterface(host, c.GUIDirectory, nm)

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

	stopDaemon := make(chan int)
	//start the node manager
	go nm.Start(stopDaemon)

	<-quit
	stopDaemon <- 1

	logger.Info("Shutting down")
	nm.Shutdown()
	logger.Info("Goodbye")

}

func main() {
	devConfig.Parse()
	Run(&devConfig)
}

//
// OLD
//

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
					fmt.Fprintf(os.Stdout, "Message received: %v\nReplyTo: %+v\n", string(msgRecvd.Contents), msgRecvd.ReplyTo)
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
