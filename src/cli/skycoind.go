// Command line interface for daemon (skycoind)
package cli

import (
    "flag"
    "github.com/op/go-logging"
    "time"
)

type DaemonConfig struct {
    Config
}

var DaemonArgs = DaemonConfig{Config{
    DisableGUI: true,
    // Disable DHT peer discovery
    DisableDHT: false,
    // Disable peer exchange
    DisablePEX: false,
    // Don't make any outgoing connections
    DisableOutgoingConnections: false,
    // Don't allowing incoming connections
    DisableIncomingConnections: false,
    // Disables networking altogether
    DisableNetworking: false,
    // Only run on localhost and only connect to others on localhost
    LocalhostOnly: false,
    // Which address to serve on. Leave blank to automatically assign to a
    // public interface
    Address: "",
    // DHT uses this port for UDP; gnet uses this for TCP incoming and outgoing
    Port: 5798,
    // How often to make outgoing connections, in seconds
    OutgoingConnectionsRate: time.Second * 5,
    // Wallet Address Version
    AddressVersion: "test",
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

    // Wallets
    WalletFile:     "",
    WalletSizeMin:  100,
    BlockchainFile: "",
    BlockSigsFile:  "",
    CanSpend:       true,

    // Centralized network configuration
    MasterPublic:     "02887355c1d32bfdff37d16de654b0e64704deaaa0c8db93acbfade81220e59727",
    MasterChain:      false,
    MasterKeys:       "",
    GenesisTimestamp: 1393645856,
    GenesisSignature: "a4feb6aac432b636e93e824d28f727c045fd7ba9cc7da92fd0a64db229b48dae46ebabc91284de55bed9a772f6450e46fa2743421d9c48fae813115f05413c9801",

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
    flag.BoolVar(&self.DisableDHT, "disable-dht", self.DisableDHT,
        "disable DHT peer discovery")
    flag.BoolVar(&self.DisablePEX, "disable-pex", self.DisablePEX,
        "disable PEX peer discovery")
    flag.BoolVar(&self.DisableOutgoingConnections, "disable-outgoing",
        self.DisableOutgoingConnections, "Don't make outgoing connections")
    flag.BoolVar(&self.DisableIncomingConnections, "disable-incoming",
        self.DisableIncomingConnections, "Don't make incoming connections")
    flag.BoolVar(&self.DisableNetworking, "disable-networking",
        self.DisableNetworking, "Disable all network activity")
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
    flag.IntVar(&self.Port, "port", self.Port, "Port to run application on")
    flag.StringVar(&self.Address, "address", self.Address,
        "IP Address to run application on. Leave empty to default to a public interface")
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
    flag.StringVar(&self.WalletFile, "wallet-file", self.WalletFile,
        "location of the wallet file. Defaults to ~/.skycoin/wallet.json")
    flag.StringVar(&self.BlockchainFile, "blockchain-file", self.BlockchainFile,
        "location of the blockchain file. Default to ~/.skycoin/blockchain.bin")
    flag.StringVar(&self.BlockSigsFile, "blocksigs-file", self.BlockSigsFile,
        "location of the block signatures file. Default to ~/.skycoin/blockchain.sigs")
    flag.BoolVar(&self.CanSpend, "can-spend", self.CanSpend,
        "is allowed to make outgoing transactions")
}
