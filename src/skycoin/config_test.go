package skycoin

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/wallet"
)

// TODO(therealssj): write better tests
func TestNewConfig(t *testing.T) {
	coinConfig, err := NewConfig("test.fiber.toml", "./testdata")
	require.NoError(t, err)
	require.Equal(t, Config{
		Blockchain: BlockchainConfig{
			GenesisSignatureStr: "eb10468d10054d15f2b6f8946cd46797779aa20a7617ceb4be884189f219bc9a164e56a5b9f7bec392a804ff3740210348d73db77a37adb542a8e08d429ac92700",
			GenesisAddressStr:   "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6",
			BlockchainPubkeyStr: "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a",
			BlockchainSeckeyStr: "",
			GenesisTimestamp:    1426562704,
			GenesisCoinVolume:   100e12,
			DefaultConnections: []string{
				"118.178.135.93:6000",
				"47.88.33.156:6000",
				"121.41.103.148:6000",
				"120.77.69.188:6000",
				"104.237.142.206:6000",
				"176.58.126.224:6000",
				"172.104.85.6:6000",
				"139.162.7.132:6000",
			},
		},
		Node: NodeConfig{
			DisablePEX:                 false,
			DisableOutgoingConnections: false,
			DisableIncomingConnections: false,
			DisableNetworking:          false,
			EnableWalletAPI:            false,
			EnableGUI:                  false,
			EnableUnversionedAPI:       false,
			EnableSeedAPI:              false,
			DisableCSRF:                false,
			LocalhostOnly:              false,
			Address:                    "",
			Port:                       6000,
			MaxOutgoingConnections:            16,
			MaxDefaultPeerOutgoingConnections: 1,
			DownloadPeerList:                  false,
			PeerListURL:                       "https://downloads.skycoin.net/blockchain/peers.txt",
			OutgoingConnectionsRate:           time.Second * 5,
			PeerlistSize:                      65535,
			WebInterface:                      true,
			WebInterfacePort:                  6420,
			WebInterfaceAddr:                  "127.0.0.1",
			WebInterfaceCert:                  "",
			WebInterfaceKey:                   "",
			WebInterfaceHTTPS:                 false,
			PrintWebInterfaceAddress:          false,
			RPCInterface:                      true,
			LaunchBrowser:                     true,
			DataDirectory:                     "~/.skycoin",
			GUIDirectory:                      "./src/gui/static/",
			ColorLog:                          true,
			LogLevel:                          "INFO",
			LogToFile:                         false,
			DisablePingPong:                   false,
			VerifyDB:                          true,
			ResetCorruptDB:                    false,
			WalletDirectory:                   "",
			WalletCryptoType:                  string(wallet.CryptoTypeScryptChacha20poly1305),
			ReadTimeout:                       10 * time.Second,
			WriteTimeout:                      60 * time.Second,
			IdleTimeout:                       120 * time.Second,
			RunMaster:                         false,
			ProfileCPU:                        false,
			ProfileCPUFile:                    "skycoin.prof",
			HTTPProf:                          false,
		},
		Build: BuildConfig{
			Version: "0.23.1-rc2",
			Commit:  "0aab9bf7730827d6fd11beb0d02096b40cea1872",
			Branch:  "test-branch",
		},
	}, coinConfig)
}
