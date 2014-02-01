// Command line interface for daemon (skycoind)
package cli

import (
    "flag"
    "github.com/op/go-logging"
)

type DaemonConfig struct {
    Config
}

var DaemonArgs = DaemonConfig{Config{
    DisableGUI:   true,
    DisableCoind: false,
    // DHT uses this port for UDP; gnet uses this for TCP incoming and outgoing
    Port: 5798,
    // Remote web interface
    WebInterface:      false,
    WebInterfacePort:  6402,
    WebInterfaceAddr:  "127.0.0.1",
    WebInterfaceCert:  "",
    WebInterfaceKey:   "",
    WebInterfaceHTTPS: true,
    // Data directory holds app data -- defaults to ~/.skycoin
    DataDirectory: "",
    // GUI directory contains assets for the html gui
    GUIDirectory: "./static/",
    // Logging
    LogLevel: logging.NOTICE,
    ColorLog: false,
    logLevel: "notice",

    // Centralized network configuration
    MasterPublic:   "02eeb60535f8f22ee49b8cfadafb0e84621e283deb6572637f7b2aa70a003c5983",
    GenesisAddress: "Db8Hc2ZVzx4P7K9jXaTpzeGRxTM3gJRxaC",
    MasterChain:    false,
    MasterKeys:     "",

    /* Developer options (don't parse these) */

    // Enable cpu profiling
    ProfileCPU: false,
    // Where the file is written to
    ProfileCPUFile: "skycoin.prof",
    // HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
    HTTPProf: false,
    // Will force it to connect to this ip:port, instead of waiting for it
    // to show up as a peer
    ConnectTo: "",
}}

func (self *DaemonConfig) register() {
    flag.BoolVar(&self.DisableCoind, "disable-daemon", self.DisableCoind,
        "disable the coin daemon")
    flag.BoolVar(&self.WebInterface, "web-interface",
        self.WebInterface, "enable the web interface")
    flag.IntVar(&self.WebInterfacePort, "web-interface-port",
        self.WebInterfacePort, "port to serve web interface on")
    flag.StringVar(&self.WebInterfaceAddr, "web-interface-addr",
        self.WebInterfaceAddr, "addr to serve web interface on")
    flag.StringVar(&self.WebInterfaceCert, "web-interface-cert",
        self.WebInterfaceCert, "cert.pem file for web interface HTTPS. "+
            "If not provided, will use cert.pem in -data-directory")
    flag.StringVar(&self.WebInterfaceKey, "web-interface-key",
        self.WebInterfaceKey, "key.pem file for web interface HTTPS. "+
            "If not provided, will use key.pem in -data-directory")
    flag.BoolVar(&self.WebInterfaceHTTPS, "web-interface-https",
        self.WebInterfaceHTTPS, "enable HTTPS for web interface")
    flag.IntVar(&self.Port, "port", self.Port,
        "Port to run application on")
    flag.StringVar(&self.DataDirectory, "data-dir", self.DataDirectory,
        "directory to store app data (defaults to ~/.skycoin)")
    flag.StringVar(&self.logLevel, "log-level", self.logLevel,
        "Choices are: debug, info, notice, warning, error, critical")
    flag.BoolVar(&self.ColorLog, "color-log", self.ColorLog,
        "Add terminal colors to log output")
    flag.StringVar(&self.GUIDirectory, "gui-dir", self.GUIDirectory,
        "static content directory for the html gui")
    flag.StringVar(&self.MasterPublic, "master-public-key", self.MasterPublic,
        "public key of the master chain")
    flag.StringVar(&self.GenesisAddress, "genesis-address", self.GenesisAddress,
        "blockchain genesis address")
}
