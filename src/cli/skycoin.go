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
    DisableGUI:    false,
    DisableDaemon: false,
    DisableDHT:    false,
    // DHT uses this port for UDP; gnet uses this for TCP incoming and outgoing
    Port: 5798,
    // How often to make outgoing connections, in seconds
    OutgoingConnectionsRate: time.Second * 5,
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
    MasterPublic:     "03c2fc73628b77512dc14c123ea741e72d27dc455e6d01141a8e7a0a83fff1fb23",
    MasterChain:      false,
    MasterKeys:       "",
    GenesisAddress:   "CL9nba1DqVADzqH6HAGC4oJzf2pRtXKEyT",
    GenesisTimestamp: 1391649057,
    GenesisSignature: "a1a09bee02a92fddaf34856aedde9c1ef626caaf31ada221fc2acc9212493e61064b32d4cfd92f38948e799f231f8c42428086405bbf42f9e913a149c0ca743f00",

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
    flag.BoolVar(&self.DisableDaemon, "disable-daemon", self.DisableDaemon,
        "disable the coin daemon")
    flag.BoolVar(&self.DisableDHT, "disable-dht", self.DisableDHT,
        "disable DHT peer discovery")
    flag.IntVar(&self.Port, "port", self.Port,
        "Port to run application on")
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
