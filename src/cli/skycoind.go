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
    EnableWebInterface: false,
    WebInterfacePort:   6402,
    // Data directory holds app data -- defaults to ~/.skycoin
    DataDirectory: "",
    // Logging
    LogLevel: logging.NOTICE,
    ColorLog: false,
    logLevel: "notice",

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
    flag.BoolVar(&self.EnableWebInterface, "enable-web-interface",
        self.EnableWebInterface, "enable the web interface")
    flag.IntVar(&self.WebInterfacePort, "web-interface-port",
        self.WebInterfacePort, "port to serve web interface on")
    flag.IntVar(&self.Port, "port", self.Port,
        "Port to run application on")
    flag.StringVar(&self.DataDirectory, "data-dir", self.DataDirectory,
        "directory to store app data (defaults to ~/.skycoin)")
    flag.StringVar(&self.logLevel, "log-level", self.logLevel,
        "Choices are: debug, info, notice, warning, error, critical")
    flag.BoolVar(&self.ColorLog, "color-log", self.ColorLog,
        "Add terminal colors to log output")
}
