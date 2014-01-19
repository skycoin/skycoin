// +build daemon

// Command line interface arguments & definitions for daemon
package cli

import (
    "flag"
    "github.com/op/go-logging"
)

var (
    DisableGUI   = true
    DisableCoind = false
    // DHT uses this port for UDP; gnet uses this for TCP incoming and outgoing
    Port = 5798
    // Remote web interface
    EnableWebInterface = false
    WebInterfacePort   = 6402
    // Data directory holds app data -- defaults to ~/.skycoin
    DataDirectory = ""
    // Logging
    LogLevel = logging.NOTICE
    ColorLog = false
    logLevel = "notice"

    /* Developer options (don't parse these) */

    // Enable cpu profiling
    ProfileCPU = false
    // Where the file is written to
    ProfileCPUFile = "skycoin.prof"
    // HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
    HTTPProf = false
    // Will force it to connect to this ip:port, instead of waiting for it
    // to show up as a peer
    ConnectTo = ""
)

func RegisterArgs() {
    flag.BoolVar(&DisableCoind, "disable-daemon", DisableCoind,
        "disable the coin daemon")
    flag.BoolVar(&EnableWebInterface, "enable-web-interface",
        EnableWebInterface, "enable the web interface")
    flag.IntVar(&WebInterfacePort, "web-interface-port",
        WebInterfacePort, "port to serve web interface on")
    flag.IntVar(&Port, "port", Port,
        "Port to run application on")
    flag.StringVar(&DataDirectory, "data-dir", DataDirectory,
        "directory to store app data (defaults to ~/.skycoin)")
    flag.StringVar(&logLevel, "log-level", logLevel,
        "Choices are: debug, info, notice, warning, error, critical")
    flag.BoolVar(&ColorLog, "color-log", ColorLog,
        "Add terminal colors to log output")
}
