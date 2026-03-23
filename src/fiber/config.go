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

// NodeConfig records the node's configuration from fiber.toml
type NodeConfig struct {
	// GenesisSignatureStr is the signature of the genesis block
	GenesisSignatureStr string `mapstructure:"genesis_signature_str"`
	// GenesisAddressStr is the skycoin address string of the genesis block
	GenesisAddressStr string `mapstructure:"genesis_address_str"`
	// BlockchainPubkeyStr is the public key of the blockchain
	BlockchainPubkeyStr string `mapstructure:"blockchain_pubkey_str"`
	// BlockchainSeckeyStr is the secret key of the blockchain
	BlockchainSeckeyStr string `mapstructure:"blockchain_seckey_str"`
	// GenesisTimestamp is the genesis block creation unix timestamp
	GenesisTimestamp uint64 `mapstructure:"genesis_timestamp"`
	// GenesisCoinVolume is the total number of coins in genesis block
	GenesisCoinVolume uint64 `mapstructure:"genesis_coin_volume"`
	// DefaultConnections are the default peer connections
	DefaultConnections []string `mapstructure:"default_connections"`
	// PeerListURL is the URL to download the peer list from
	PeerListURL string `mapstructure:"peer_list_url"`
	// Port is the port that the wire protocol listens on
	Port int `mapstructure:"port"`
	// WebInterfacePort is the port the web/api interface listens on
	WebInterfacePort int `mapstructure:"web_interface_port"`

	// UnconfirmedBurnFactor is the burn factor to apply when verifying unconfirmed transactions
	UnconfirmedBurnFactor uint32 `mapstructure:"unconfirmed_burn_factor"`
	// UnconfirmedMaxTransactionSize is the maximum size of an unconfirmed transaction
	UnconfirmedMaxTransactionSize uint32 `mapstructure:"unconfirmed_max_transaction_size"`
	// UnconfirmedMaxDropletPrecision is the maximum number of decimal places for an unconfirmed transaction
	UnconfirmedMaxDropletPrecision uint8 `mapstructure:"unconfirmed_max_decimals"`
	// CreateBlockBurnFactor is the burn factor to apply when creating blocks
	CreateBlockBurnFactor uint32 `mapstructure:"create_block_burn_factor"`
	// CreateBlockMaxTransactionSize is the maximum size of a transaction when creating a block
	CreateBlockMaxTransactionSize uint32 `mapstructure:"create_block_max_transaction_size"`
	// CreateBlockMaxDropletPrecision is the maximum number of decimal places when creating a block
	CreateBlockMaxDropletPrecision uint8 `mapstructure:"create_block_max_decimals"`
	// MaxBlockTransactionsSize is the maximum total size of transactions in a block
	MaxBlockTransactionsSize uint32 `mapstructure:"max_block_transactions_size"`

	// DisplayName is the display name of the coin in the wallet e.g. Skycoin
	DisplayName string `mapstructure:"display_name"`
	// Ticker is the coin's ticker e.g. SKY
	Ticker string `mapstructure:"ticker"`
	// CoinHoursName is the display name of coin hours e.g. "Coin Hours"
	CoinHoursName string `mapstructure:"coin_hours_display_name"`
	// CoinHoursNameSingular is the singular form of the coin hours display name e.g. "Coin Hour"
	CoinHoursNameSingular string `mapstructure:"coin_hours_display_name_singular"`
	// CoinHoursTicker is the ticker of coin hours e.g. SCH
	CoinHoursTicker string `mapstructure:"coin_hours_ticker"`
	// QrURIPrefix is the prefix for QR code URIs
	QrURIPrefix string `mapstructure:"qr_uri_prefix"`
	// ExplorerURL is the URL of the blockchain explorer
	ExplorerURL string `mapstructure:"explorer_url"`
	// VersionURL is the URL for wallet to check the latest version number
	VersionURL string `mapstructure:"version_url"`
	// Bip44Coin is the default "coin" value of the bip44 path
	Bip44Coin bip44.CoinType `mapstructure:"bip44_coin"`

	// PriceTickerID is the coin ID used for price lookups (e.g. "sky-skycoin" for CoinPaprika, "aixexchange" for CoinGecko)
	PriceTickerID string `mapstructure:"price_ticker_id"`
	// PriceTickerSource is the price API source: "coinpaprika" or "coingecko"
	PriceTickerSource string `mapstructure:"price_ticker_source"`

	// These fields are set by cmd/newcoin and are not configured in the fiber.toml file
	CoinName string
	// Ascii Font rendering of CoinName
	// CoinAscii is the ASCII art representation of the coin
	CoinAscii     string //nolint:revive
	DataDirectory string
}

// ParamsConfig are the parameters used to generate params/params.go.
// These parameters are exposed in an importable package `params` because they
// may need to be imported by libraries that would not know the node's configured CLI options.
type ParamsConfig struct {
	// MaxCoinSupply is the maximum supply of coins
	MaxCoinSupply uint64 `mapstructure:"max_coin_supply"`
	// InitialUnlockedCount is the initial number of unlocked distribution addresses
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

// NewConfig loads blockchain config parameters from a config file.
// Uses an isolated viper instance to avoid polluting global state.
// default file is: fiber.toml in the project root
// JSON, toml or yaml file can be used (toml preferred).
func NewConfig(configName, appDir string) (Config, error) {
	v := viper.New()

	// check that file is of supported type
	confNameSplit := strings.Split(configName, ".")
	fileType := confNameSplit[len(confNameSplit)-1]
	switch fileType {
	case "toml", "json", "yaml", "yml":
		v.SetConfigType(confNameSplit[len(confNameSplit)-1])
	default:
		return Config{}, fmt.Errorf("invalid blockchain config file type: %s", fileType)
	}

	configName = configName[:len(configName)-(len(fileType)+1)]
	v.SetConfigName(configName)

	v.AddConfigPath(appDir)
	v.AddConfigPath(".")

	// set defaults
	setDefaultsOn(v)

	params := Config{}

	if err := v.ReadInConfig(); err != nil {
		return params, err
	}

	if err := v.Unmarshal(&params); err != nil {
		return params, err
	}

	return params, nil
}

func setDefaultsOn(v *viper.Viper) {
	// node defaults
	v.SetDefault("node.genesis_coin_volume", 100e12)
	v.SetDefault("node.port", 6000)
	v.SetDefault("node.web_interface_port", 6420)
	v.SetDefault("node.unconfirmed_burn_factor", 10)
	v.SetDefault("node.unconfirmed_max_transaction_size", 32*1024)
	v.SetDefault("node.unconfirmed_max_decimals", 3)
	v.SetDefault("node.create_block_burn_factor", 10)
	v.SetDefault("node.create_block_max_transaction_size", 32*1024)
	v.SetDefault("node.create_block_max_decimals", 3)
	v.SetDefault("node.max_block_transactions_size", 32*1024)
	v.SetDefault("node.display_name", "Skycoin")
	v.SetDefault("node.ticker", "SKY")
	v.SetDefault("node.coin_hours_display_name", "Coin Hours")
	v.SetDefault("node.coin_hours_display_name_singular", "Coin Hour")
	v.SetDefault("node.coin_hours_ticker", "SCH")
	v.SetDefault("node.qr_uri_prefix", "skycoin")
	v.SetDefault("node.explorer_url", "https://explorer.skycoin.com")
	v.SetDefault("node.version_url", "https://version.skycoin.com/skycoin/version.txt")
	v.SetDefault("node.bip44_coin", bip44.CoinTypeSkycoin)
	v.SetDefault("node.price_ticker_id", "sky-skycoin")
	v.SetDefault("node.price_ticker_source", "coinpaprika")

	// build defaults
	v.SetDefault("build.commit", "")
	v.SetDefault("build.branch", "")

	// params defaults
	v.SetDefault("params.max_coin_supply", 1e8)
	v.SetDefault("params.initial_unlocked_count", 25)
	v.SetDefault("params.unlock_address_rate", 5)
	v.SetDefault("params.unlock_time_interval", 60*60*24*365)
	v.SetDefault("params.user_max_decimals", 3)
	v.SetDefault("params.user_burn_factor", 10)
	v.SetDefault("params.user_max_transaction_size", 32*1024)
}
