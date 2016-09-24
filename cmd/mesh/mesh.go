package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/skycoin/skycoin/src/mesh/domain"
	mesh "github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/node/connection"
	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

var ()

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
}

func (c *Config) Parse() {
	c.register()
	flag.Parse()
	c.postProcess()
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
	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser,
		"launch system default webbrowser at client startup")

	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory,
		"directory to store app data (defaults to ~/.skycoin)")

	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory,
		"static content directory for the html gui")

}

func (c *Config) postProcess() {
	var err error

	c.DataDirectory = util.InitDataDir(c.DataDirectory)
	if c.WebInterfaceCert == "" {
		c.WebInterfaceCert = filepath.Join(c.DataDirectory, "cert.pem")
	}
	if c.WebInterfaceKey == "" {
		c.WebInterfaceKey = filepath.Join(c.DataDirectory, "key.pem")
	}

	//ll, err := logging.LogLevel(c.logLevel)
	//panicIfError(err, "Invalid -log-level %s", c.logLevel)
	//c.LogLevel = ll

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
}

func catchInterrupt(quit chan<- int) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	quit <- 1
}

func Run(c *Config) {

	c.GUIDirectory = util.ResolveResourceDirectory(c.GUIDirectory)

	// If the user Ctrl-C's, shutdown properly
	quit := make(chan int)
	go catchInterrupt(quit)

	if c.WebInterface {
		var err error
		if c.WebInterfaceHTTPS {
			// Verify cert/key parameters, and if neither exist, create them
			errs := gui.CreateCertIfNotExists(host, c.WebInterfaceCert, c.WebInterfaceKey)
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
}

func main() {
	fmt.Fprintln(os.Stdout, "Starting Node Manager Service...")

	config := nodemanager.ServerConfig
	fmt.Fprintln(os.Stdout, "PubKey:", config.Node.PubKey)
	fmt.Fprintln(os.Stdout, "ChaCha20Key:", config.Node.ChaCha20Key)
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
