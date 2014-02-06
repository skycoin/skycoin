// Command line interface arguments
package cli

import (
    "flag"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/util"
    "log"
    "path/filepath"
)

type Args interface {
    register()
    postProcess()
    getConfig() *Config
}

type Config struct {
    DisableGUI    bool
    DisableDaemon bool
    DisableDHT    bool
    // DHT uses this port for UDP; gnet uses this for TCP incoming and outgoing
    Port int
    // How often to make outgoing connections, in seconds
    OutgoingConnectionsRate uint64
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
    WalletFile     string
    WalletSizeMin  int
    BlockchainFile string
    BlockSigsFile  string
    // Is allowed to make outgoing transactions
    CanSpend bool

    // Centralized network configuration
    MasterPublic     string
    MasterChain      bool
    MasterKeys       string
    GenesisAddress   string
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
    if self.WalletFile == "" {
        self.WalletFile = filepath.Join(self.DataDirectory, "wallet.json")
    }
    if self.BlockchainFile == "" {
        self.BlockchainFile = filepath.Join(self.DataDirectory, "blockchain.bin")
    }
    if self.BlockSigsFile == "" {
        self.BlockSigsFile = filepath.Join(self.DataDirectory, "blockchain.sigs")
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
