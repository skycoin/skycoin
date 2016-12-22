package nodemanager

import (
	"encoding/json"
	"os"

	"github.com/skycoin/skycoin/src/cipher"
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

func (nm *NodeManager) GetFromFile(configIndex string) error {

	configDatas, err := loadConfigs(configIndex)
	if err != nil {
		return err
	}

	transportDatas, err := loadTransports(configIndex)
	if err != nil {
		return err
	}

	nm.ConfigList = testConfigsFromData(configDatas)
	for pubKey := range nm.ConfigList {
		nm.PubKeyList = append(nm.PubKeyList, pubKey)
	}
	nm.connectConfigs(transportDatas)
	return nil
}

func (nm *NodeManager) PutToFile(configIndex string) error {

	configDatas, transportDatas := nodesToConfigData(nm.ConfigList)

	err := saveConfigs(configIndex, configDatas)
	if err != nil {
		return err
	}

	err = saveTransports(configIndex, transportDatas)
	return err
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

func saveConfigs(configIndex string, configData []*ConfigData) error {

	configFile, err := os.Create(configIndex + "_nodes.cfg")
	if err != nil {
		return err
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(configData)

	return err
}

func saveTransports(transportIndex string, transportData []*TransportData) error {

	transportFile, err := os.Create(transportIndex + "_transports.cfg")
	if err != nil {
		return err
	}
	defer transportFile.Close()

	encoder := json.NewEncoder(transportFile)
	err = encoder.Encode(transportData)

	return err
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

func nodesToConfigData(configList map[cipher.PubKey]*TestConfig) ([]*ConfigData, []*TransportData) {
	configDatas := []*ConfigData{}
	transportDatas := []*TransportData{}

	for pubKeyFrom, config := range configList {
		configData := &ConfigData{pubKeyFrom, config.ExternalAddress, config.StartPort}
		configDatas = append(configDatas, configData)
		for _, peerToPeer := range config.PeerToPeers {
			pubKeyTo := peerToPeer.Peer
			transportData := &TransportData{pubKeyFrom, pubKeyTo}
			found := false
			for _, td := range transportDatas {
				if (td.PubKey1 == pubKeyFrom && td.PubKey2 == pubKeyTo) || (td.PubKey1 == pubKeyTo && td.PubKey2 == pubKeyFrom) {
					found = true
					break
				}
			}
			if !found {
				transportDatas = append(transportDatas, transportData)
			}
		}
	}
	return configDatas, transportDatas
}
