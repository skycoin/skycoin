package skycoin

import (
	"strings"
	"github.com/spf13/viper"
	"fmt"
)

type BlockchainConfig struct {
	GenesisSignatureStr string   `mapstructure:"genesis_signature_str"`
	GenesisAddressStr   string   `mapstructure:"genesis_address_str"`
	BlockchainPubkeyStr string   `mapstructure:"blockchain_pubkey_str"`
	BlockchainSeckeyStr string   `mapstructure:"blockchain_seckey_str"`
	GenesisTimestamp    uint64   `mapstructure:"genesis_timestamp"`
	GenesisCoinVolume   uint64   `mapstructure:"genesis_coin_volume"`
	DefaultConnections  []string `mapstructure:"default_connections"`
	DataDirectory       string   `mapstructure:"data_directory"`
}

func setDefaults() {
	viper.SetDefault("genesis_coin_volume", 100e12)
	viper.SetDefault("data_directory", "~/.skycoin")
}

// NewBlockchainConfig loads blockchain config parameters from a config file
// default file is: fiber.toml in the project root
// JSON, toml or yaml file can be used (toml preferred).
func NewBlockchainConfig(configName, appDir string) (BlockchainConfig, error) {
	// set viper parameters

	// check that file is of supported type
	confNameSplit := strings.Split(configName, ".")
	fileType := confNameSplit[len(confNameSplit)-1]
	switch fileType {
	case "toml", "json", "yaml", "yml":
		viper.SetConfigType(confNameSplit[len(confNameSplit)-1])
	default:
		return BlockchainConfig{}, fmt.Errorf("invalid blockchain config file type: %s", fileType)
	}

	configName = configName[:len(configName)-(len(fileType)+1)]
	viper.SetConfigName(configName)

	viper.AddConfigPath(appDir)
	viper.AddConfigPath(".")

	// set defaults
	setDefaults()

	cfg := BlockchainConfig{}

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	// TODO(therealssj): validate the config values

	return cfg, nil
}
