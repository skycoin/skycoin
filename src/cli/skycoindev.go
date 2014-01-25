// Command line interface arguments for development (skycoindev)
package cli

import (
    "flag"
    "github.com/op/go-logging"
)

type DevConfig struct {
    Config
}

var DevArgs = DevConfig{Config{
    DisableGUI:   true,
    DisableCoind: false,
    // DHT uses this port for UDP; gnet uses this for TCP incoming and outgoing
    Port: 5798,
    // Remote web interface
    WebInterface:     false,
    WebInterfacePort: 6402,
    WebInterfaceAddr: "127.0.0.1",
    // Data directory holds app data -- defaults to ~/.skycoin
    DataDirectory: "",
    // Data directory holds app data -- defaults to ~/.skycoin
    GUIDirectory: "./static/",
    // Logging
    LogLevel: logging.DEBUG,
    ColorLog: true,
    logLevel: "DEBUG",

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
}}

func (self *DevConfig) register() {
    flag.BoolVar(&self.DisableCoind, "disable-daemon", self.DisableCoind,
        "disable the coin daemon")
    flag.BoolVar(&self.DisableGUI, "disable-gui", self.DisableGUI,
        "disable the gui")
    flag.BoolVar(&self.WebInterface, "web-interface",
        self.WebInterface, "enable the web interface")
    flag.IntVar(&self.WebInterfacePort, "web-interface-port",
        self.WebInterfacePort, "port to serve web interface on")
    flag.StringVar(&self.WebInterfaceAddr, "web-interface-addr",
        self.WebInterfaceAddr, "addr to serve web interface on")
    flag.IntVar(&self.Port, "port", self.Port,
        "Port to run application on")
    flag.StringVar(&self.DataDirectory, "data-dir", self.DataDirectory,
        "directory to store app data (defaults to ~/.skycoin)")
    flag.StringVar(&self.ConnectTo, "connect-to", self.ConnectTo,
        "connect to this ip only")
    flag.BoolVar(&self.ProfileCPU, "profile-cpu", self.ProfileCPU,
        "enable cpu profiling")
    flag.StringVar(&self.ProfileCPUFile, "profile-cpu-file",
        self.ProfileCPUFile, "where to write the cpu profile file")
    flag.BoolVar(&self.HTTPProf, "http-prof", self.HTTPProf,
        "Run the http profiling interface")
    flag.StringVar(&self.logLevel, "log-level", self.logLevel,
        "Choices are: debug, info, notice, warning, error, critical")
    flag.BoolVar(&self.ColorLog, "color-log", self.ColorLog,
        "Add terminal colors to log output")
    flag.StringVar(&self.GUIDirectory, "gui-dir", self.GUIDirectory,
        "static content directory for the html gui")
}
