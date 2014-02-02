// Command line interface arguments for client (skycoin)
package cli

import (
    "flag"
    "github.com/op/go-logging"
)

type ClientConfig struct {
    Config
}

var ClientArgs = ClientConfig{Config{
    DisableGUI:   false,
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
    GUIDirectory: "",
    // Logging
    LogLevel: logging.NOTICE,
    ColorLog: false,
    logLevel: "notice",

    // Centralized network configuration
    MasterPublic:   "0223f1cd8652e64f0b2b6960e25c5799426220d44d50d016a4c64ecefb5b0043db",
    GenesisAddress: "2bGDcaLJH8Ve7JBgqgSjNbHGSSSagEWtBrJ",
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

func (self *ClientConfig) register() {
    flag.BoolVar(&self.DisableCoind, "disable-daemon", self.DisableCoind,
        "disable the coin daemon")
    flag.IntVar(&self.Port, "port", self.Port,
        "Port to run application on")
    flag.StringVar(&self.DataDirectory, "data-dir", self.DataDirectory,
        "directory to store app data (defaults to ~/.skycoin)")
    flag.StringVar(&self.logLevel, "log-level", self.logLevel,
        "Choices are: debug, info, notice, warning, error, critical")
    flag.StringVar(&self.MasterPublic, "master-public-key", self.MasterPublic,
        "public key of the master chain")
    flag.StringVar(&self.GenesisAddress, "genesis-address", self.GenesisAddress,
        "blockchain genesis address")
}
