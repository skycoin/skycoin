/*
Package fiber provides configuration definitions and utilities for managing fiber coins
*/
package fiber

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/skycoin/skycoin/src/cipher/bip44"
)

// Config records fiber coin parameters
type Config struct {
	Node   NodeConfig   `mapstructure:"node"`
	Params ParamsConfig `mapstructure:"params"`
}

// NodeConfig configures the default CLI options for the skycoin node.
// These parameters are loaded via cmd/skycoin/skycoin.go into src/skycoin/skycoin.go.
type NodeConfig struct {
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
	UnconfirmedBurnFactor uint32 `mapstructure:"unconfirmed_burn_factor"`
	// UnconfirmedMaxTransactionSize is the maximum size of an unconfirmed transaction
	UnconfirmedMaxTransactionSize uint32 `mapstructure:"unconfirmed_max_transaction_size"`
	// UnconfirmedMaxDropletPrecision is the maximum number of decimals allowed in an unconfirmed transaction
	UnconfirmedMaxDropletPrecision uint8 `mapstructure:"unconfirmed_max_decimals"`
	// CreateBlockBurnFactor is the burn factor to apply to transactions when publishing blocks
	CreateBlockBurnFactor uint32 `mapstructure:"create_block_burn_factor"`
	// CreateBlockMaxTransactionSize is the maximum size of an transaction when publishing blocks
	CreateBlockMaxTransactionSize uint32 `mapstructure:"create_block_max_transaction_size"`
	// CreateBlockMaxDropletPrecision is the maximum number of decimals allowed in a transaction when publishing blocks
	CreateBlockMaxDropletPrecision uint8 `mapstructure:"create_block_max_decimals"`
	// MaxBlockTransactionsSize is the maximum total size of transactions in a block when publishing a block
	MaxBlockTransactionsSize uint32 `mapstructure:"max_block_transactions_size"`

	// DisplayName is the display name of the coin in the wallet e.g. Skycoin
	DisplayName string `mapstructure:"display_name"`
	// Ticker is the coin's price ticker, e.g. SKY
	Ticker string `mapstructure:"ticker"`
	// CoinHoursName is the name of the coinhour asset type, e.g. Coin Hours
	CoinHoursName string `mapstructure:"coin_hours_display_name"`
	// CoinHoursNameSingular is the singular form of the name of the coinhour asset type, e.g. Coin Hour
	CoinHoursNameSingular string `mapstructure:"coin_hours_display_name_singular"`
	// CoinHoursTicker is the name of the coinhour asset type's price ticker, e.g. SCH (Skycoin Coin Hours)
	CoinHoursTicker string `mapstructure:"coin_hours_ticker"`
	// ExplorerURL is the URL of the public explorer
	ExplorerURL string `mapstructure:"explorer_url"`
	// VersionURL is the URL for wallet to check the latest version number
	VersionURL string `mapstructure:"version_url"`
	// Bip44Coin is the default "coin" value of the bip44 path
	Bip44Coin bip44.CoinType `mapstructure:"bip44_coin"`

	// These fields are set by cmd/newcoin and are not configured in the fiber.toml file
	CoinName      string
	DataDirectory string
}

// ParamsConfig are the parameters used to generate params/params.go.
// These parameters are exposed in an importable package `params` because they
// may need to be imported by libraries that would not know the node's configured CLI options.
type ParamsConfig struct {
	// MaxCoinSupply is the maximum supply of coins
	MaxCoinSupply uint64 `mapstructure:"max_coin_supply"`
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

// NewConfig loads blockchain config parameters from a config file
// default file is: fiber.toml in the project root
// JSON, toml or yaml file can be used (toml preferred).
func NewConfig(configName, appDir string) (Config, error) {
	// set viper parameters
	// check that file is of supported type
	confNameSplit := strings.Split(configName, ".")
	fileType := confNameSplit[len(confNameSplit)-1]
	switch fileType {
	case "toml", "json", "yaml", "yml":
		viper.SetConfigType(confNameSplit[len(confNameSplit)-1])
	default:
		return Config{}, fmt.Errorf("invalid blockchain config file type: %s", fileType)
	}

	configName = configName[:len(configName)-(len(fileType)+1)]
	viper.SetConfigName(configName)

	viper.AddConfigPath(appDir)
	viper.AddConfigPath(".")

	// set defaults
	setDefaults()

	params := Config{}

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
	viper.SetDefault("node.unconfirmed_burn_factor", 10)
	viper.SetDefault("node.unconfirmed_max_transaction_size", 32*1024)
	viper.SetDefault("node.unconfirmed_max_decimals", 3)
	viper.SetDefault("node.create_block_burn_factor", 10)
	viper.SetDefault("node.create_block_max_transaction_size", 32*1024)
	viper.SetDefault("node.create_block_max_decimals", 3)
	viper.SetDefault("node.max_block_transactions_size", 32*1024)
	viper.SetDefault("node.display_name", "Skycoin")
	viper.SetDefault("node.ticker", "SKY")
	viper.SetDefault("node.coin_hours_display_name", "Coin Hours")
	viper.SetDefault("node.coin_hours_display_name_singular", "Coin Hour")
	viper.SetDefault("node.coin_hours_ticker", "SCH")
	viper.SetDefault("node.explorer_url", "https://explorer.skycoin.com")
	viper.SetDefault("node.version_url", "https://version.skycoin.com/skycoin/version.txt")
	viper.SetDefault("node.bip44_coin", bip44.CoinTypeSkycoin)

	// build defaults
	viper.SetDefault("build.commit", "")
	viper.SetDefault("build.branch", "")

	// params defaults
	viper.SetDefault("params.max_coin_supply", 1e8)
	viper.SetDefault("params.initial_unlocked_count", 25)
	viper.SetDefault("params.unlock_address_rate", 5)
	viper.SetDefault("params.unlock_time_interval", 60*60*24*365)
	viper.SetDefault("params.user_max_decimals", 3)
	viper.SetDefault("params.user_burn_factor", 10)
	viper.SetDefault("params.user_max_transaction_size", 32*1024)
}
