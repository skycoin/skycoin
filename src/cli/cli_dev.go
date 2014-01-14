// +build dev

// Command line interface arguments & definitions for development
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
    LogLevel = logging.DEBUG
    logLevel = "DEBUG"

    /* Developer options */

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
    flag.BoolVar(&DisableGUI, "disable-gui", DisableGUI,
        "disable the gui")
    flag.BoolVar(&EnableWebInterface, "enable-web-interface",
        EnableWebInterface, "enable the web interface")
    flag.IntVar(&WebInterfacePort, "web-interface-port",
        WebInterfacePort, "port to serve web interface on")
    flag.IntVar(&Port, "port", Port,
        "Port to run application on")
    flag.StringVar(&DataDirectory, "data-dir", DataDirectory,
        "directory to store app data (defaults to ~/.skycoin)")
    flag.StringVar(&ConnectTo, "connect-to", ConnectTo,
        "connect to this ip only")
    flag.BoolVar(&ProfileCPU, "profile-cpu", ProfileCPU,
        "enable cpu profiling")
    flag.StringVar(&ProfileCPUFile, "profile-cpu-file", ProfileCPUFile,
        "where to write the cpu profile file")
    flag.BoolVar(&HTTPProf, "http-prof", HTTPProf,
        "Run the http profiling interface")
    flag.StringVar(&logLevel, "log-level", logLevel,
        "Choices are: debug, info, notice, warning, error, critical")
}
