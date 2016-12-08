package nodemanager

import (
	"encoding/json"
	"os"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/transport"
	"github.com/skycoin/skycoin/src/mesh/transport/physical"
)

type ConfigData struct {
	PubKey          cipher.PubKey    `json:"pubkey"`
	ExternalAddress string           `json:"external_address"`
	Transports      []*TransportData `json:"transports"`
}

type TransportData struct {
	IncomingPort    int    `json:"incoming_port"`
	OutgoingAddress string `json:"outgoing_address"`
	OutgoingPort    int    `json:"outgoing_port"`
}

func loadConfig(configFileName string) ([]*ConfigData, error) {

	configData := []*ConfigData{}

	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&configData)
	if err != nil {
		return nil, err
	}

	return configData, nil
}

func TestConfigsFromFile(configFileName string) []*TestConfig {
	configDatas, err := loadConfig(configFileName)
	if err != nil {
		panic(err)
	}

	testConfigs := []*TestConfig{}

	for _, configData := range configDatas {
		testConfig := createTestConfigFromData(configData)
		testConfigs = append(testConfigs, testConfig)
	}

	return testConfigs
}

func createTestConfigFromData(configData *ConfigData) *TestConfig {
	testConfig := &TestConfig{}
	testConfig.ExternalAddress = configData.ExternalAddress
	testConfig.NodeConfig = NewNodeConfig()
	testConfig.NodeConfig.PubKey = configData.PubKey
	testConfig.TransportConfig = transport.CreateTransportConfig(testConfig.NodeConfig.PubKey)
	testConfig.UDPConfigs = []physical.UDPConfig{}
	testConfig.PeerToPeers = map[string]*Peer{}

	testConfig.AddPeersToConnectNew(configData)

	return testConfig
}
