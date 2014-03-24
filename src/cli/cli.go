// Command line interface arguments
package cli

import (
    "flag"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/util"
    "log"
    "path/filepath"
    "time"
)

type Args interface {
    register()
    postProcess()
    getConfig() *Config
}

type Config struct {
    DisableGUI bool
    // Disable DHT peer discovery
    DisableDHT bool
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
    // DHT uses this port for UDP; gnet uses this for TCP incoming and outgoing
    Port int
    // How often to make outgoing connections
    OutgoingConnectionsRate time.Duration
    // Wallet Address Version
    AddressVersion string
    // Remote web interface
    WebInterface      bool
    WebInterfacePort  int
    WebInterfaceAddr  string
    WebInterfaceCert  string
    WebInterfaceKey   string
    WebInterfaceHTTPS bool
    // Data directory holds app data -- defaults to ~/.skycoin
    DataDirectory string
    // GUI directory contains assets for the html gui
    GUIDirectory string
    // Logging
    LogLevel logging.Level
    ColorLog bool
    // This is the value registered with flag, it is converted to LogLevel
    // after parsing
    logLevel string

    // Wallets
    WalletDirectory string
    BlockchainFile  string
    BlockSigsFile   string
    // Is allowed to make outgoing transactions
    CanSpend bool

    // Centralized network configuration
    MasterPublic     string
    MasterChain      bool
    MasterKeys       string
    GenesisSignature string
    GenesisTimestamp uint64

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
}

func (self *Config) register() {
    log.Panic("Config.register must be overridden")
}

func (self *Config) postProcess() {
    self.DataDirectory = util.InitDataDir(self.DataDirectory)
    if self.WebInterfaceCert == "" {
        self.WebInterfaceCert = filepath.Join(self.DataDirectory, "cert.pem")
    }
    if self.WebInterfaceKey == "" {
        self.WebInterfaceKey = filepath.Join(self.DataDirectory, "key.pem")
    }
    if self.MasterKeys == "" {
        self.MasterKeys = filepath.Join(self.DataDirectory, "master.keys")
    }
    if self.BlockchainFile == "" {
        self.BlockchainFile = filepath.Join(self.DataDirectory, "blockchain.bin")
    }
    if self.BlockSigsFile == "" {
        self.BlockSigsFile = filepath.Join(self.DataDirectory, "blockchain.sigs")
    }
    if self.WalletDirectory == "" {
        self.WalletDirectory = filepath.Join(self.DataDirectory, "wallets/")
    }
    ll, err := logging.LogLevel(self.logLevel)
    if err != nil {
        log.Panic("Invalid -log-level %s: %v\n", self.logLevel, err)
    }
    self.LogLevel = ll
}

func (self *Config) getConfig() *Config {
    return self
}

// Parses arguments defined in a struct that satisfies Config interface
func ParseArgs(args Args) *Config {
    args.register()
    flag.Parse()
    args.postProcess()
    return args.getConfig()
}
