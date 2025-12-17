/*
skycoin daemon
*/
package commands

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
edit the template instead at template/coin.template
and then run cmd/newcoin
*/

import (
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/calvin"
	"github.com/spf13/cobra"

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
	// DB VERSION - WARNING: CHANGING WILL MAKE THE SYNCED BLOCKCHAIN NOT WORK WITH OLDER SOFTWARE VERSIONS
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
	CoinName = "skycoin"

	// GenesisSignatureStr hex string of genesis signature
	GenesisSignatureStr = "eb10468d10054d15f2b6f8946cd46797779aa20a7617ceb4be884189f219bc9a164e56a5b9f7bec392a804ff3740210348d73db77a37adb542a8e08d429ac92700"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1426562704
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 100000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"139.162.121.185:6000",
		"172.104.164.147:6000",
		"139.162.248.183:6000",
		"45.56.109.228:6000",
		"173.230.130.174:6000",
		"172.104.99.241:6000",
		"172.104.57.147:6000",
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
		PeerListURL:         "https://downloads.skycoin.com/blockchain/peers.txt",
		Port:                6000,
		WebInterfacePort:    6420,
		DataDirectory:       "$HOME/.skycoin",

		UnconfirmedBurnFactor:          10,
		UnconfirmedMaxTransactionSize:  32768,
		UnconfirmedMaxDropletPrecision: 3,
		CreateBlockBurnFactor:          10,
		CreateBlockMaxTransactionSize:  32768,
		CreateBlockMaxDropletPrecision: 3,
		MaxBlockTransactionsSize:       32768,

		DisplayName:           "Skycoin",
		Ticker:                "SKY",
		CoinHoursName:         "Coin Hours",
		CoinHoursNameSingular: "Coin Hour",
		CoinHoursTicker:       "SCH",
		QrURIPrefix:           "skycoin",
		ExplorerURL:           "https://explorer.skycoin.com",
		VersionURL:            "https://version.skycoin.com/skycoin/version.txt",
		Bip44Coin:             8000,
	})
)

func init() {
	// Check for FIBER_TOML environment variable to load custom fibercoin config
	// This sets the default values before flag registration, so CLI flags naturally override
	if fiberTomlPath := os.Getenv("FIBER_TOML"); fiberTomlPath != "" {
		if err := nodeConfig.LoadFromFiberConfig(fiberTomlPath); err != nil {
			logger.Errorf("Failed to load FIBER_TOML config: %v", err)
			os.Exit(1)
		}
		logger.Infof("Loaded fiber config from FIBER_TOML: %s", fiberTomlPath)
	}

	// Check for GENESIS environment variable to load genesis wallet credentials
	// This takes precedence over fiber.toml values (address/pubkey/seckey)
	if genesisWalletPath := os.Getenv("GENESIS"); genesisWalletPath != "" {
		if err := nodeConfig.LoadFromGenesisWallet(genesisWalletPath); err != nil {
			logger.Errorf("Failed to load GENESIS wallet: %v", err)
			os.Exit(1)
		}
		logger.Infof("Loaded genesis credentials from GENESIS: %s", genesisWalletPath)
	}

	// Set dynamic Long description based on loaded config
	// Use DisplayName if available (from fiber.toml), otherwise CoinName (from template)
	coinName := nodeConfig.Fiber.DisplayName
	if coinName == "" {
		coinName = nodeConfig.Fiber.Name
	}
	if coinName == "" {
		coinName = "skycoin"
	}
	// Use lowercase for ASCII art and wallet text
	coinNameLower := strings.ToLower(coinName)
	RootCmd.Long = calvin.AsciiFont(coinNameLower) + "\n " + coinNameLower + " wallet"

	nodeConfig.RegisterFlags(RootCmd)
}

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "skycoin",
	Short: "skycoin wallet",
	Long:  "", // Set dynamically in init()
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}
