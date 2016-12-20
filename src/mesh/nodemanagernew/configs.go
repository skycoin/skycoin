package nodemanager

import (
	"encoding/json"
	"os"

	"github.com/skycoin/skycoin/src/cipher"
	//	"github.com/skycoin/skycoin/src/mesh/transport"
	//	"github.com/skycoin/skycoin/src/mesh/transport/physical"
)

type ConfigData struct {
	PubKey          cipher.PubKey `json:"pubkey"`
	ExternalAddress string        `json:"external_address"`
	StartPort       int           `json:"start_port"`
}

type TransportData struct {
	PubKey1 cipher.PubKey `json:"pubkey_1"`
	PubKey2 cipher.PubKey `json:"pubkey_2"`
}

func loadConfigs(configIndex string) ([]*ConfigData, error) {

	configData := []*ConfigData{}

	configFile, err := os.Open(configIndex + "_nodes.cfg")
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

func loadTransports(transportIndex string) ([]*TransportData, error) {

	transportData := []*TransportData{}

	transportFile, err := os.Open(transportIndex + "_transports.cfg")
	if err != nil {
		return nil, err
	}
	defer transportFile.Close()

	decoder := json.NewDecoder(transportFile)
	err = decoder.Decode(&transportData)
	if err != nil {
		return nil, err
	}

	return transportData, nil
}

func (nm *NodeManager) GetFromFile(configIndex string) {
	configDatas, err := loadConfigs(configIndex)
	if err != nil {
		panic(err)
	}

	transportDatas, err := loadTransports(configIndex)
	if err != nil {
		panic(err)
	}

	nm.ConfigList = testConfigsFromData(configDatas)
	nm.connectConfigs(transportDatas)
}

func testConfigsFromData(configDatas []*ConfigData) map[cipher.PubKey]*TestConfig {

	testConfigs := map[cipher.PubKey]*TestConfig{}

	for _, configData := range configDatas {
		testConfig := createTestConfigFromData(configData)
		testConfigs[testConfig.NodeConfig.PubKey] = testConfig
	}

	return testConfigs
}

func (nm *NodeManager) connectConfigs(transportDatas []*TransportData) {

	configs := nm.ConfigList

	for _, transportData := range transportDatas {
		config1 := configs[transportData.PubKey1]
		config2 := configs[transportData.PubKey2]
		ConnectNodeToNode(config1, config2)
	}
}

func createTestConfigFromData(configData *ConfigData) *TestConfig {
	config := CreateConfig(configData.ExternalAddress, configData.StartPort)
	config.NodeConfig.PubKey = configData.PubKey

	return config
}
