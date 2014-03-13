// Command line interface arguments for client (skycoin)
package cli

import (
    "flag"
    "github.com/op/go-logging"
    "time"
)

type ClientConfig struct {
    Config
}

var ClientArgs = ClientConfig{Config{
    DisableGUI: false,
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
    GUIDirectory: "",
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
    MasterPublic:     "02b0333bd8f1910663b8b1f60fb2e154b70436a2c19efb79cdbdf09bf9bb2056dc",
    MasterChain:      false,
    MasterKeys:       "",
    GenesisTimestamp: 1394689119,
    GenesisSignature: "173e1cdf628e78ae4946af4415f070e2aad5a1f4273b77971f8d42a6eb7ff3af68d0d7a3360460e96123f93decf43c28abbc02a65ffb243e525131ba357f21d800",

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
    flag.IntVar(&self.Port, "port", self.Port, "Port to run application on")
    flag.StringVar(&self.Address, "address", self.Address,
        "IP Address to run application on. Leave empty to default to a public interface")
    flag.StringVar(&self.DataDirectory, "data-dir", self.DataDirectory,
        "directory to store app data (defaults to ~/.skycoin)")
    flag.StringVar(&self.logLevel, "log-level", self.logLevel,
        "Choices are: debug, info, notice, warning, error, critical")
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
