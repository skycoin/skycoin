package fiber

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/bip44"
)

// TODO(therealssj): write better tests
func TestNewConfig(t *testing.T) {
	coinConfig, err := NewConfig("test.fiber.toml", "./testdata")
	require.NoError(t, err)
	require.Equal(t, Config{
		Node: NodeConfig{
			GenesisSignatureStr: "eb10468d10054d15f2b6f8946cd46797779aa20a7617ceb4be884189f219bc9a164e56a5b9f7bec392a804ff3740210348d73db77a37adb542a8e08d429ac92700",
			GenesisAddressStr:   "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6",
			BlockchainPubkeyStr: "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a",
			BlockchainSeckeyStr: "",
			GenesisTimestamp:    1426562704,
			GenesisCoinVolume:   100e12,
			DefaultConnections: []string{
				"118.178.135.93:6000",
				"47.88.33.156:6000",
				"104.237.142.206:6000",
				"176.58.126.224:6000",
				"172.104.85.6:6000",
				"139.162.7.132:6000",
			},
			Port:                           6000,
			PeerListURL:                    "https://downloads.skycoin.com/blockchain/peers.txt",
			WebInterfacePort:               6420,
			UnconfirmedBurnFactor:          10,
			UnconfirmedMaxTransactionSize:  777,
			UnconfirmedMaxDropletPrecision: 3,
			CreateBlockBurnFactor:          9,
			CreateBlockMaxTransactionSize:  1234,
			CreateBlockMaxDropletPrecision: 4,
			MaxBlockTransactionsSize:       1111,
			DisplayName:                    "Testcoin",
			Ticker:                         "TST",
			CoinHoursName:                  "Testcoin Hours",
			CoinHoursNameSingular:          "Testcoin Hour",
			CoinHoursTicker:                "TCH",
			QrURIPrefix:                    "skycoin",
			ExplorerURL:                    "https://explorer.testcoin.com",
			VersionURL:                     "https://version.testcoin.com/testcoin/version.txt",
			Bip44Coin:                      bip44.CoinTypeSkycoin,
		},
		Params: ParamsConfig{
			MaxCoinSupply:           1e8,
			UnlockAddressRate:       5,
			InitialUnlockedCount:    25,
			UnlockTimeInterval:      60 * 60 * 24 * 365,
			UserBurnFactor:          3,
			UserMaxTransactionSize:  999,
			UserMaxDropletPrecision: 2,
		},
	}, coinConfig)
}
