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

// NodeParameters configures the default CLI options for the skycoin node.
// These parameters are loaded via cmd/skycoin/skycoin.go into src/skycoin/skycoin.go.
type NodeParameters struct {
	// Port is the default port that the wire protocol communicates over
	Port int `mapstructure:"port"`
	// WebInterfacePort is the default port that the web/gui interface serves on
	WebInterfacePort int `mapstructure:"web_interface_port"`
	// GenesisSignatureStr is a hex-encoded signature of the genesis block input
	GenesisSignatureStr string `mapstructure:"genesis_signature_str"`
	// GenesisAddressStr is the skycoin address that the genesis coins were sent to in the genesis block
	GenesisAddressStr string `mapstructure:"genesis_address_str"`
	// BlockchainPubkeyStr is a hex-encoded public key used to validate published blocks
	BlockchainPubkeyStr string `mapstructure:"blockchain_pubkey_str"`
	// BlockchainSeckey is a hex-encoded secret key required for block publishing.
	// It must correspond to BlockchainPubkeyStr
	BlockchainSeckeyStr string `mapstructure:"blockchain_seckey_str"`
	// GenesisTimestamp is the timestamp of the genesis block
	GenesisTimestamp uint64 `mapstructure:"genesis_timestamp"`
	// GenesisCoinVolume is the total number of coins in the genesis block
	GenesisCoinVolume uint64 `mapstructure:"genesis_coin_volume"`
	// DefaultConnections are the default "trusted" connections a node will try to connect to for bootstrapping
	DefaultConnections []string `mapstructure:"default_connections"`
	// PeerlistURL is a URL pointing to a newline-separated list of ip:ports that are used for bootstrapping (but they are not "trusted")
	PeerListURL string `mapstructure:"peer_list_url"`
	// UnconfirmedBurnFactor is the burn factor to apply when verifying unconfirmed transactions
	UnconfirmedBurnFactor uint64 `mapstructure:"unconfirmed_burn_factor"`
	// UnconfirmedMaxTransactionSize is the maximum size of an unconfirmed transaction
	UnconfirmedMaxTransactionSize int `mapstructure:"unconfirmed_max_transaction_size"`
	// UnconfirmedMaxDropletPrecision is the maximum number of decimals allowed in an unconfirmed transaction
	UnconfirmedMaxDropletPrecision uint8 `mapstructure:"unconfirmed_max_decimals"`
	// CreateBlockBurnFactor is the burn factor to apply to transactions when publishing blocks
	CreateBlockBurnFactor uint64 `mapstructure:"create_block_burn_factor"`
	// CreateBlockMaxTransactionSize is the maximum size of an transaction when publishing blocks
	CreateBlockMaxTransactionSize int `mapstructure:"create_block_max_transaction_size"`
	// CreateBlockMaxDropletPrecision is the maximum number of decimals allowed in a transaction when publishing blocks
	CreateBlockMaxDropletPrecision uint8 `mapstructure:"create_block_max_decimals"`
	// MaxBlockSize is the maximum size of blocks when publishing blocks
	MaxBlockSize int `mapstructure:"max_block_size"`

	// These fields are set by cmd/newcoin and are not configured in the fiber.toml file
	CoinName      string
	DataDirectory string
}

// ParamsParameters are the parameters used to generate params/params.go.
// These parameters are exposed in an importable package `params` because they
// may need to be imported by libraries that would not know the node's configured CLI options.
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
	// UnlockTimeInterval is the distribution address unlock time interval, measured in seconds.
	// Once the InitialUnlockedCount is exhausted, UnlockAddressRate addresses will be unlocked per UnlockTimeInterval
	UnlockTimeInterval uint64 `mapstructure:"unlock_time_interval"`
	// UserMaxDropletPrecision represents the decimal precision of droplets
	UserMaxDropletPrecision uint64 `mapstructure:"user_max_decimals"`
	// UserMaxTransactionSize is max size of a user-created transaction (typically equal to the max size of a block)
	UserMaxTransactionSize int `mapstructure:"user_max_transaction_size"`
	// DistributionAddresses are addresses that received coins from the genesis address in the first block,
	// used to calculate current and max supply and do distribution timelocking
	DistributionAddresses []string `mapstructure:"distribution_addresses"`
	// UserBurnFactor inverse fraction of coinhours that must be burned, this value is used when creating transactions
	UserBurnFactor uint64 `mapstructure:"user_burn_factor"`
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
	viper.SetDefault("node.unconfirmed_burn_factor", 2)
	viper.SetDefault("node.unconfirmed_max_transaction_size", 32*1024)
	viper.SetDefault("node.unconfirmed_max_decimals", 3)
	viper.SetDefault("node.create_block_burn_factor", 2)
	viper.SetDefault("node.create_block_max_transaction_size", 32*1024)
	viper.SetDefault("node.create_block_max_decimals", 3)
	viper.SetDefault("node.max_block_size", 32*1024)

	// build defaults
	viper.SetDefault("build.commit", "")
	viper.SetDefault("build.branch", "")

	// params defaults
	viper.SetDefault("params.max_coin_supply", 1e8)
	viper.SetDefault("params.distribution_addresses_total", 100)
	viper.SetDefault("params.initial_unlocked_count", 25)
	viper.SetDefault("params.unlock_address_rate", 5)
	viper.SetDefault("params.unlock_time_interval", 60*60*24*365)
	viper.SetDefault("params.user_max_decimals", 3)
	viper.SetDefault("params.user_burn_factor", 2)
	viper.SetDefault("params.user_max_transaction_size", 32*1024)
}
