package skycoin

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Parameters records fiber coin parameters
type Parameters struct {
	Node   NodeParameters   `mapstructure:"node"`
	Params ParamsParameters `mapstructure:"params"`
}

// NodeParameters records the node's configurable parameters
type NodeParameters struct {
	CoinName            string   `mapstructure:"coin_name"`
	PeerListURL         string   `mapstructure:"peer_list_url"`
	Port                int      `mapstructure:"port"`
	WebInterfacePort    int      `mapstructure:"web_interface_port"`
	GenesisSignatureStr string   `mapstructure:"genesis_signature_str"`
	GenesisAddressStr   string   `mapstructure:"genesis_address_str"`
	BlockchainPubkeyStr string   `mapstructure:"blockchain_pubkey_str"`
	BlockchainSeckeyStr string   `mapstructure:"blockchain_seckey_str"`
	GenesisTimestamp    uint64   `mapstructure:"genesis_timestamp"`
	GenesisCoinVolume   uint64   `mapstructure:"genesis_coin_volume"`
	DefaultConnections  []string `mapstructure:"default_connections"`

	DataDirectory string
}

// ParamsParameters are the parameters used to generate config/blockchain.go
type ParamsParameters struct {
	// MaxCoinSupply is the maximum supply of coins
	MaxCoinSupply uint64 `mapstructure:"max_coin_supply"`

	// DistributionAddressesTotal is the number of distribution addresses
	DistributionAddressesTotal uint64 `mapstructure:"distribution_addresses_total"`

	// DistributionAddressInitialBalance is the initial balance of each distribution address
	DistributionAddressInitialBalance uint64

	// InitialUnlockedCount is the initial number of unlocked addresses
	InitialUnlockedCount uint64 `mapstructure:"initial_unlocked_count"`

	// UnlockAddressRate is the number of addresses to unlock per unlock time interval
	UnlockAddressRate uint64 `mapstructure:"unlock_address_rate"`

	// UnlockTimeInterval is the distribution address unlock time interval, measured in seconds
	// Once the InitialUnlockedCount is exhausted,
	// UnlockAddressRate addresses will be unlocked per UnlockTimeInterval
	UnlockTimeInterval uint64 `mapstructure:"unlock_time_interval"`

	// MaxDropletPrecision represents the decimal precision of droplets
	MaxDropletPrecision uint64 `mapstructure:"max_droplet_precision"`

	//DefaultMaxBlockSize is max block size
	DefaultMaxBlockSize int `mapstructure:"default_max_block_size"`

	// Addresses that received coins from the genesis address in the first block,
	// used to calculate current and max supply and do distribution timelocking
	DistributionAddresses []string `mapstructure:"distribution_addresses"`

	// BurnFactor inverse fraction of coinhours that must be burned (can be overridden with COINHOUR_BURN_FACTOR env var)
	CoinHourBurnFactor uint64 `mapstructure:"coinhour_burn_factor"`
}

// NewParameters loads blockchain config parameters from a config file
// default file is: fiber.toml in the project root
// JSON, toml or yaml file can be used (toml preferred).
func NewParameters(configName, appDir string) (Parameters, error) {
	// set viper parameters
	// check that file is of supported type
	confNameSplit := strings.Split(configName, ".")
	fileType := confNameSplit[len(confNameSplit)-1]
	switch fileType {
	case "toml", "json", "yaml", "yml":
		viper.SetConfigType(confNameSplit[len(confNameSplit)-1])
	default:
		return Parameters{}, fmt.Errorf("invalid blockchain config file type: %s", fileType)
	}

	configName = configName[:len(configName)-(len(fileType)+1)]
	viper.SetConfigName(configName)

	viper.AddConfigPath(appDir)
	viper.AddConfigPath(".")

	// set defaults
	setDefaults()

	params := Parameters{}

	if err := viper.ReadInConfig(); err != nil {
		return params, err
	}

	if err := viper.Unmarshal(&params); err != nil {
		return params, err
	}

	return params, nil
}

func setDefaults() {
	// node defaults
	viper.SetDefault("node.genesis_coin_volume", 100e12)
	viper.SetDefault("node.port", 6000)
	viper.SetDefault("node.web_interface_port", 6420)

	// build defaults
	viper.SetDefault("build.commit", "")
	viper.SetDefault("build.branch", "")

	// params defaults
	viper.SetDefault("params.max_coin_supply", 1e8)
	viper.SetDefault("params.distribution_addresses_total", 100)
	viper.SetDefault("params.initial_unlocked_count", 25)
	viper.SetDefault("params.unlock_address_rate", 5)
	viper.SetDefault("params.unlock_time_interval", 60*60*24*365)
	viper.SetDefault("params.max_droplet_precision", 3)
	viper.SetDefault("params.default_max_block_size", 32*1024)
	viper.SetDefault("params.coinhour_burn_factor", 2)
}
