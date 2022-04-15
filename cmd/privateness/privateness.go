/*
skycoin daemon
*/
package main

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

import (
	"flag"
	_ "net/http/pprof"
	"os"

	"github.com/skycoin/skycoin/src/fiber"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/skycoin"
	"github.com/skycoin/skycoin/src/util/logging"

	// register the supported wallets
	_ "github.com/skycoin/skycoin/src/wallet/bip44wallet"
	_ "github.com/skycoin/skycoin/src/wallet/collection"
	_ "github.com/skycoin/skycoin/src/wallet/deterministic"
	_ "github.com/skycoin/skycoin/src/wallet/xpubwallet"
)

var (
	// Version of the node. Can be set by -ldflags
	Version = "0.27.1"
	// Commit ID. Can be set by -ldflags
	Commit = ""
	// Branch name. Can be set by -ldflags
	Branch = ""
	// ConfigMode (possible values are "", "STANDALONE_CLIENT").
	// This is used to change the default configuration.
	// Can be set by -ldflags
	ConfigMode = ""

	logger = logging.MustGetLogger("main")

	// CoinName name of coin
	CoinName = "privateness"

	// GenesisSignatureStr hex string of genesis signature
	GenesisSignatureStr = "052805bdf985499bd223ae7a45ec938c085dbdf604a0a86bcdc2cbcaf14583e00271efd72717acda4977a6039536eeafd66962bd42620813315d9580d784c4c900"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "24GJTLPMoz61sV4J4qg1n14x5qqDwXqyJJy"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "02933015bd2fa1e0a885c05fb08eb7c647bf8c3188ed5120b51d0d09ccaf525036"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1650046003
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 165000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"154.16.118.97:6006",
		"151.80.37.6:6006",
		"179.61.232.155:6006",
	}

	nodeConfig = skycoin.NewNodeConfig(ConfigMode, fiber.NodeConfig{
		CoinName:            CoinName,
		GenesisSignatureStr: GenesisSignatureStr,
		GenesisAddressStr:   GenesisAddressStr,
		GenesisCoinVolume:   GenesisCoinVolume,
		GenesisTimestamp:    GenesisTimestamp,
		BlockchainPubkeyStr: BlockchainPubkeyStr,
		BlockchainSeckeyStr: BlockchainSeckeyStr,
		DefaultConnections:  DefaultConnections,
		PeerListURL:         "http://nodes.privateness.network/blockchain/peers2.txt",
		Port:                6006,
		WebInterfacePort:    6660,
		DataDirectory:       "$HOME/.privateness",

		UnconfirmedBurnFactor:          5,
		UnconfirmedMaxTransactionSize:  32768,
		UnconfirmedMaxDropletPrecision: 6,
		CreateBlockBurnFactor:          5,
		CreateBlockMaxTransactionSize:  32768,
		CreateBlockMaxDropletPrecision: 6,
		MaxBlockTransactionsSize:       32768,

		DisplayName:           "Privateness",
		Ticker:                "NESS",
		CoinHoursName:         "Coin Hours",
		CoinHoursNameSingular: "Coin Hour",
		CoinHoursTicker:       "NCH",
		QrURIPrefix:           "privateness",
		ExplorerURL:           "https://explorer.privateness.network",
		VersionURL:            "https://nodes.privateness.network/blockchain/version.txt",
		Bip44Coin:             8000,
	})

	parseFlags = true
)

func init() {
	nodeConfig.RegisterFlags()
}

func main() {
	if parseFlags {
		flag.Parse()
	}

	// create a new fiber coin instance
	coin := skycoin.NewCoin(skycoin.Config{
		Node: nodeConfig,
		Build: readable.BuildInfo{
			Version: Version,
			Commit:  Commit,
			Branch:  Branch,
		},
	}, logger)

	// parse config values
	if err := coin.ParseConfig(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// run fiber coin node
	if err := coin.Run(); err != nil {
		os.Exit(1)
	}
}
