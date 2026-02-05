package skycoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/fiber"
	"github.com/skycoin/skycoin/src/kvstorage"

	"log"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/useragent"
)

var (
	help = false //nolint:unused
)

// Config records skycoin node and build config
type Config struct {
	Node  NodeConfig
	Build readable.BuildInfo
}

// NodeConfig records the node's configuration
type NodeConfig struct {
	// Name of the coin
	CoinName string

	// Disable peer exchange
	DisablePEX bool
	// Download peer list
	DownloadPeerList bool
	// Download the peers list from this URL
	PeerListURL string
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// Allow connections from legacy peers that don't send blockchain pubkey in introduction
	LegacyPeerCompat bool
	// Enable GUI
	EnableGUI bool
	// Disable CSRF check in the wallet API
	DisableCSRF bool
	// Disable Host, Origin and Referer header check in the wallet API
	DisableHeaderCheck bool
	// Disable CSP disable content-security-policy in http response
	DisableCSP bool
	// Comma separated list of API sets enabled on the remote web interface
	EnabledAPISets string
	// Comma separated list of API sets disabled on the remote web interface
	DisabledAPISets string
	// Enable all of API sets. Applies before disabling individual sets
	EnableAllAPISets bool

	enabledAPISets map[string]struct{}
	// Comma separate list of hostnames to accept in the Host header, used to bypass the Host header check which only applies to localhost addresses
	HostWhitelist string
	hostWhitelist []string

	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	// gnet uses this for TCP incoming and outgoing
	Port int
	// MaxConnections is the maximum number of total connections allowed
	MaxConnections int
	// Maximum outgoing connections to maintain
	MaxOutgoingConnections int
	// Maximum incoming connections to maintain
	MaxIncomingConnections int
	// Maximum default outgoing connections
	MaxDefaultPeerOutgoingConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
	// MaxOutgoingMessageLength maximum size of outgoing messages
	MaxOutgoingMessageLength int
	// MaxIncomingMessageLength maximum size of incoming messages
	MaxIncomingMessageLength int
	// MaxLastBlocksCount maximum number of blocks to response to API /api/v1/last_blocks
	MaxLastBlocksCount uint64
	// PeerlistSize represents the maximum number of peers that the pex would maintain
	PeerlistSize int
	// Wallet Address Version
	// AddressVersion string
	// Remote web interface
	WebInterface bool
	// Remote web interface port
	WebInterfacePort int
	// Remote web interface address
	WebInterfaceAddr string
	// Remote web interface certificate
	WebInterfaceCert string
	// Remote web interface key
	WebInterfaceKey string
	// Remote web interface HTTPS support
	WebInterfaceHTTPS bool
	// Remote web interface username and password
	WebInterfaceUsername string
	WebInterfacePassword string
	// Allow web interface auth without HTTPS
	WebInterfacePlaintextAuth bool

	// Launch System Default Browser after client startup
	LaunchBrowser bool

	// Data directory holds app data
	DataDirectory string
	// GUI directory contains assets for the HTML interface
	GUIDirectory string

	// Timeouts for the HTTP listener
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration

	// Remark to include in user agent sent in the wire protocol introduction
	UserAgentRemark string
	userAgent       useragent.Data

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Disable "Reply to ping", "Received pong" log messages
	DisablePingPong bool

	// Verify the database integrity after loading
	VerifyDB bool
	// Reset the database if integrity checks fail, and continue running
	ResetCorruptDB bool

	// Transaction verification parameters for unconfirmed transactions
	UnconfirmedVerifyTxn params.VerifyTxn
	// Transaction verification parameters for transactions when creating blocks
	CreateBlockVerifyTxn params.VerifyTxn
	// Maximum total size of transactions in a block
	MaxBlockTransactionsSize uint32

	unconfirmedBurnFactor          uint64
	maxUnconfirmedTransactionSize  uint64
	unconfirmedMaxDropletPrecision uint64
	createBlockBurnFactor          uint64
	createBlockMaxTransactionSize  uint64
	createBlockMaxDropletPrecision uint64
	maxBlockSize                   uint64

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string
	// Wallet crypto type
	WalletCryptoType string

	// Key-value storage
	// Default to ${DataDirectory}/data
	KVStorageDirectory  string
	EnabledStorageTypes []kvstorage.Type

	// Disable the hardcoded default peers
	DisableDefaultPeers bool
	// Load custom peers from disk
	CustomPeersFile string

	RunBlockPublisher bool

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// Enable HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Expose HTTP profiling on this interface
	HTTPProfHost string

	DBPath     string
	DBReadOnly bool
	LogToFile  bool
	Version    bool // show node version

	GenesisSignatureStr string
	GenesisAddressStr   string
	BlockchainPubkeyStr string
	BlockchainSeckeyStr string
	GenesisTimestamp    uint64
	GenesisCoinVolume   uint64
	DefaultConnections  []string

	genesisSignature cipher.Sig
	genesisAddress   cipher.Address
	genesisHash      cipher.SHA256

	blockchainPubkey cipher.PubKey
	blockchainSeckey cipher.SecKey

	Fiber readable.FiberConfig

	// Paths for auto-updating fiber.toml
	fiberTomlPath     string // Set from FIBER_TOML env
	genesisWalletPath string // Set from GENESIS env
	// Distribution addresses from fiber.toml [params] section
	distributionAddresses      []string // Loaded from fiber.Params.DistributionAddresses
	distributionMaxCoinSupply  uint64   // Loaded from fiber.Params.MaxCoinSupply
	distributionUnlockedCount  uint64   // Loaded from fiber.Params.InitialUnlockedCount
	distributionUnlockRate     uint64   // Loaded from fiber.Params.UnlockAddressRate
	distributionUnlockInterval uint64   // Loaded from fiber.Params.UnlockTimeInterval
}

// NewNodeConfig returns a new node config instance
func NewNodeConfig(mode string, node fiber.NodeConfig) NodeConfig {
	nodeConfig := NodeConfig{
		CoinName:            node.CoinName,
		GenesisSignatureStr: node.GenesisSignatureStr,
		GenesisAddressStr:   node.GenesisAddressStr,
		GenesisCoinVolume:   node.GenesisCoinVolume,
		GenesisTimestamp:    node.GenesisTimestamp,
		BlockchainPubkeyStr: node.BlockchainPubkeyStr,
		BlockchainSeckeyStr: node.BlockchainSeckeyStr,
		DefaultConnections:  node.DefaultConnections,
		// Disable peer exchange
		DisablePEX: false,
		// Don't make any outgoing connections
		DisableOutgoingConnections: false,
		// Don't allowing incoming connections
		DisableIncomingConnections: false,
		// Disables networking altogether
		DisableNetworking: false,
		// Enable GUI
		EnableGUI: false,
		// Disable CSRF check in the wallet API
		DisableCSRF: false,
		// Disable Host, Origin and Referer header check in the wallet API
		DisableHeaderCheck: false,
		// DisableCSP disable content-security-policy in http response
		DisableCSP: false,
		// Only run on localhost and only connect to others on localhost
		LocalhostOnly: false,
		// Which address to serve on. Leave blank to automatically assign to a
		// public interface
		Address: "",
		// gnet uses this for TCP incoming and outgoing
		Port: node.Port,
		// MaxConnections is the maximum number of total connections allowed
		MaxConnections: 128,
		// MaxOutgoingConnections is the maximum outgoing connections allowed
		MaxOutgoingConnections: 8,
		// MaxIncomingConnections is the maximum incoming connections allowed
		MaxIncomingConnections: 120,
		// MaxDefaultOutgoingConnections is the maximum default outgoing connections allowed
		MaxDefaultPeerOutgoingConnections: 2,
		DownloadPeerList:                  true,
		PeerListURL:                       node.PeerListURL,
		// How often to make outgoing connections, in seconds
		OutgoingConnectionsRate:  time.Second * 5,
		MaxOutgoingMessageLength: 256 * 1024,
		MaxIncomingMessageLength: 1024 * 1024,
		MaxLastBlocksCount:       256,
		PeerlistSize:             65535,
		// Wallet Address Version
		// AddressVersion: "test",
		// Remote web interface
		WebInterface:      true,
		WebInterfacePort:  node.WebInterfacePort,
		WebInterfaceAddr:  "127.0.0.1",
		WebInterfaceCert:  "",
		WebInterfaceKey:   "",
		WebInterfaceHTTPS: false,
		EnabledAPISets: strings.Join([]string{
			api.EndpointsRead,
			api.EndpointsTransaction,
		}, ","),
		DisabledAPISets:  "",
		EnableAllAPISets: false,

		LaunchBrowser: false,
		// Data directory holds app data
		DataDirectory: node.DataDirectory,
		// Web GUI static resources - now embedded
		GUIDirectory: "", //"./src/gui/static/",
		// Logging
		ColorLog:        true,
		LogLevel:        "INFO",
		LogToFile:       false,
		DisablePingPong: false,

		VerifyDB:       false,
		ResetCorruptDB: false,

		// Blockchain/transaction validation
		UnconfirmedVerifyTxn: params.VerifyTxn{
			BurnFactor:          node.UnconfirmedBurnFactor,
			MaxTransactionSize:  node.UnconfirmedMaxTransactionSize,
			MaxDropletPrecision: node.UnconfirmedMaxDropletPrecision,
		},
		CreateBlockVerifyTxn: params.VerifyTxn{
			BurnFactor:          node.CreateBlockBurnFactor,
			MaxTransactionSize:  node.CreateBlockMaxTransactionSize,
			MaxDropletPrecision: node.CreateBlockMaxDropletPrecision,
		},
		MaxBlockTransactionsSize: node.MaxBlockTransactionsSize,

		// Wallets
		WalletDirectory:  "",
		WalletCryptoType: string(crypto.DefaultCryptoType),

		// Key-value storage
		KVStorageDirectory: "",
		EnabledStorageTypes: []kvstorage.Type{
			kvstorage.TypeTxIDNotes,
			kvstorage.TypeGeneral,
		},

		// Timeout settings for http.Server
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		HTTPReadTimeout:  time.Second * 10,
		HTTPWriteTimeout: time.Second * 60,
		HTTPIdleTimeout:  time.Second * 120,

		RunBlockPublisher: false,

		// Enable cpu profiling
		ProfileCPU: false,
		// Where the file is written to
		ProfileCPUFile: "cpu.prof",
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf:     false,
		HTTPProfHost: "localhost:6060",

		Fiber: readable.FiberConfig{
			Name:                  node.CoinName,
			DisplayName:           node.DisplayName,
			Ticker:                node.Ticker,
			CoinHoursName:         node.CoinHoursName,
			CoinHoursNameSingular: node.CoinHoursNameSingular,
			CoinHoursTicker:       node.CoinHoursTicker,
			QrURIPrefix:           node.QrURIPrefix,
			ExplorerURL:           node.ExplorerURL,
			VersionURL:            node.VersionURL,
			Bip44Coin:             node.Bip44Coin,
		},
	}

	nodeConfig.applyConfigMode(mode)

	return nodeConfig
}

func (c *Config) postProcess() error {
	//	if help {
	//		flag.Usage()
	//		os.Exit(0)
	//	}

	var err error
	if c.Node.GenesisSignatureStr != "" {
		c.Node.genesisSignature, err = cipher.SigFromHex(c.Node.GenesisSignatureStr)
		panicIfError(err, "Invalid Signature")
	}

	if c.Node.GenesisAddressStr != "" {
		c.Node.genesisAddress, err = cipher.DecodeBase58Address(c.Node.GenesisAddressStr)
		panicIfError(err, "Invalid Address")
	}

	// Compute genesis block hash
	gb, err := coin.NewGenesisBlock(c.Node.genesisAddress, c.Node.GenesisCoinVolume, c.Node.GenesisTimestamp)
	if err != nil {
		panicIfError(err, "Create genesis hash failed")
	}
	c.Node.genesisHash = gb.HashHeader()

	if c.Node.BlockchainPubkeyStr != "" {
		c.Node.blockchainPubkey, err = cipher.PubKeyFromHex(c.Node.BlockchainPubkeyStr)
		panicIfError(err, "Invalid Pubkey")
	}
	if c.Node.BlockchainSeckeyStr != "" {
		c.Node.blockchainSeckey, err = cipher.SecKeyFromHex(c.Node.BlockchainSeckeyStr)
		panicIfError(err, "Invalid Seckey")
		c.Node.BlockchainSeckeyStr = ""
	}
	if c.Node.BlockchainSeckeyStr != "" {
		c.Node.blockchainSeckey = cipher.SecKey{}
	}

	home := file.UserHome()
	c.Node.DataDirectory, err = file.InitDataDir(replaceHome(c.Node.DataDirectory, home))
	panicIfError(err, "Invalid DataDirectory")

	if c.Node.WebInterfaceCert == "" {
		c.Node.WebInterfaceCert = filepath.Join(c.Node.DataDirectory, "skycoind.cert")
	} else {
		c.Node.WebInterfaceCert = replaceHome(c.Node.WebInterfaceCert, home)
	}

	if c.Node.WebInterfaceKey == "" {
		c.Node.WebInterfaceKey = filepath.Join(c.Node.DataDirectory, "skycoind.key")
	} else {
		c.Node.WebInterfaceKey = replaceHome(c.Node.WebInterfaceKey, home)
	}

	if c.Node.WalletDirectory == "" {
		c.Node.WalletDirectory = filepath.Join(c.Node.DataDirectory, "wallets")
	} else {
		c.Node.WalletDirectory = replaceHome(c.Node.WalletDirectory, home)
	}

	if c.Node.KVStorageDirectory == "" {
		c.Node.KVStorageDirectory = filepath.Join(c.Node.DataDirectory, "data")
	} else {
		c.Node.KVStorageDirectory = replaceHome(c.Node.KVStorageDirectory, home)
	}
	if len(c.Node.EnabledStorageTypes) == 0 {
		c.Node.EnabledStorageTypes = []kvstorage.Type{
			kvstorage.TypeGeneral,
			kvstorage.TypeTxIDNotes,
		}
	}

	if c.Node.DBPath == "" {
		c.Node.DBPath = filepath.Join(c.Node.DataDirectory, "data.db")
	} else {
		c.Node.DBPath = replaceHome(c.Node.DBPath, home)
	}

	userAgentData := useragent.Data{
		Coin:    c.Node.CoinName,
		Version: c.Build.Version,
		Remark:  c.Node.UserAgentRemark,
	}

	if _, err := userAgentData.Build(); err != nil {
		return err
	}

	c.Node.userAgent = userAgentData

	apiSets, err := buildAPISets(c.Node)
	if err != nil {
		return err
	}

	// Don't open browser to load wallets if wallet apis are disabled.
	c.Node.enabledAPISets = apiSets
	if _, ok := c.Node.enabledAPISets[api.EndpointsWallet]; !ok {
		c.Node.EnableGUI = false
		c.Node.LaunchBrowser = false
	}

	if c.Node.EnableGUI && c.Node.GUIDirectory != "" {
		c.Node.GUIDirectory = file.ResolveResourceDirectory(c.Node.GUIDirectory)
	}

	if c.Node.DisableDefaultPeers {
		c.Node.DefaultConnections = nil
	}

	if c.Node.HostWhitelist != "" {
		if c.Node.DisableHeaderCheck {
			return errors.New("host whitelist should be empty when header check is disabled")
		}
		c.Node.hostWhitelist = strings.Split(c.Node.HostWhitelist, ",")
	}

	httpAuthEnabled := c.Node.WebInterfaceUsername != "" || c.Node.WebInterfacePassword != ""
	if httpAuthEnabled && !c.Node.WebInterfaceHTTPS && !c.Node.WebInterfacePlaintextAuth {
		return errors.New("Web interface auth enabled but HTTPS is not enabled. Use -web-interface-plaintext-auth=true if this is desired")
	}

	if c.Node.MaxConnections < c.Node.MaxOutgoingConnections+c.Node.MaxIncomingConnections {
		return errors.New("-max-connections must be >= -max-outgoing-connections + -max-incoming-connections")
	}

	if c.Node.MaxOutgoingConnections > c.Node.MaxConnections {
		return errors.New("-max-outgoing-connections cannot be higher than -max-connections")
	}

	if c.Node.MaxIncomingConnections > c.Node.MaxConnections {
		return errors.New("-max-incoming-connections cannot be higher than -max-connections")
	}

	if c.Node.maxBlockSize > math.MaxUint32 {
		return errors.New("-max-block-size exceeds MaxUint32")
	}
	if c.Node.maxUnconfirmedTransactionSize > math.MaxUint32 {
		return errors.New("-max-txn-size-unconfirmed exceeds MaxUint32")
	}
	if c.Node.unconfirmedBurnFactor > math.MaxUint32 {
		return errors.New("-burn-factor-unconfirmed exceeds MaxUint32")
	}
	if c.Node.createBlockBurnFactor > math.MaxUint32 {
		return errors.New("-burn-factor-create-block exceeds MaxUint32")
	}

	if c.Node.unconfirmedMaxDropletPrecision > math.MaxUint8 {
		return errors.New("-max-decimals-unconfirmed exceeds MaxUint8")
	}
	if c.Node.createBlockMaxDropletPrecision > math.MaxUint8 {
		return errors.New("-max-decimals-create-block exceeds MaxUint8")
	}

	c.Node.UnconfirmedVerifyTxn.BurnFactor = uint32(c.Node.unconfirmedBurnFactor)                  //nolint:gosec
	c.Node.UnconfirmedVerifyTxn.MaxTransactionSize = uint32(c.Node.maxUnconfirmedTransactionSize)  //nolint:gosec
	c.Node.UnconfirmedVerifyTxn.MaxDropletPrecision = uint8(c.Node.unconfirmedMaxDropletPrecision) //nolint:gosec
	c.Node.CreateBlockVerifyTxn.BurnFactor = uint32(c.Node.createBlockBurnFactor)                  //nolint:gosec
	c.Node.CreateBlockVerifyTxn.MaxTransactionSize = uint32(c.Node.createBlockMaxTransactionSize)  //nolint:gosec // Config conversion
	c.Node.CreateBlockVerifyTxn.MaxDropletPrecision = uint8(c.Node.createBlockMaxDropletPrecision) //nolint:gosec
	c.Node.MaxBlockTransactionsSize = uint32(c.Node.maxBlockSize)                                  //nolint:gosec // Config conversion

	if c.Node.UnconfirmedVerifyTxn.MaxTransactionSize < params.MinTransactionSize {
		return fmt.Errorf("-max-txn-size-unconfirmed must be >= params.MinTransactionSize (%d)", params.MinTransactionSize)
	}
	if c.Node.UnconfirmedVerifyTxn.MaxTransactionSize < params.UserVerifyTxn.MaxTransactionSize {
		return fmt.Errorf("-max-txn-size-unconfirmed must be >= params.UserVerifyTxn.MaxTransactionSize (%d)", params.UserVerifyTxn.MaxTransactionSize)
	}
	if c.Node.CreateBlockVerifyTxn.MaxTransactionSize < params.MinTransactionSize {
		return fmt.Errorf("-max-txn-size-create-block must be >= params.MinTransactionSize (%d)", params.MinTransactionSize)
	}
	if c.Node.CreateBlockVerifyTxn.MaxTransactionSize < params.UserVerifyTxn.MaxTransactionSize {
		return fmt.Errorf("-max-txn-size-create-block must be >= params.UserVerifyTxn.MaxTransactionSize (%d)", params.UserVerifyTxn.MaxTransactionSize)
	}

	if c.Node.MaxBlockTransactionsSize < params.MinTransactionSize {
		return fmt.Errorf("-max-block-size must be >= params.MinTransactionSize (%d)", params.MinTransactionSize)
	}
	if c.Node.MaxBlockTransactionsSize < params.UserVerifyTxn.MaxTransactionSize {
		return fmt.Errorf("-max-block-size must be >= params.UserVerifyTxn.MaxTransactionSize (%d)", params.UserVerifyTxn.MaxTransactionSize)
	}
	if c.Node.MaxBlockTransactionsSize < c.Node.UnconfirmedVerifyTxn.MaxTransactionSize {
		return errors.New("-max-block-size must be >= -max-txn-size-unconfirmed")
	}
	if c.Node.MaxBlockTransactionsSize < c.Node.CreateBlockVerifyTxn.MaxTransactionSize {
		return errors.New("-max-block-size must be >= -max-txn-size-create-block")
	}

	if c.Node.UnconfirmedVerifyTxn.BurnFactor < params.MinBurnFactor {
		return fmt.Errorf("-burn-factor-unconfirmed must be >= params.MinBurnFactor (%d)", params.MinBurnFactor)
	}
	if c.Node.UnconfirmedVerifyTxn.BurnFactor < params.UserVerifyTxn.BurnFactor {
		return fmt.Errorf("-burn-factor-unconfirmed must be >= params.UserVerifyTxn.BurnFactor (%d)", params.UserVerifyTxn.BurnFactor)
	}

	if c.Node.CreateBlockVerifyTxn.BurnFactor < params.MinBurnFactor {
		return fmt.Errorf("-burn-factor-create-block must be >= params.MinBurnFactor (%d)", params.MinBurnFactor)
	}
	if c.Node.CreateBlockVerifyTxn.BurnFactor < params.UserVerifyTxn.BurnFactor {
		return fmt.Errorf("-burn-factor-create-block must be >= params.UserVerifyTxn.BurnFactor (%d)", params.UserVerifyTxn.BurnFactor)
	}

	if c.Node.UnconfirmedVerifyTxn.MaxDropletPrecision > droplet.Exponent {
		return fmt.Errorf("-max-decimals-unconfirmed must be <= droplet.Exponent (%d)", droplet.Exponent)
	}
	if c.Node.UnconfirmedVerifyTxn.MaxDropletPrecision < params.UserVerifyTxn.MaxDropletPrecision {
		return fmt.Errorf("-max-decimals-unconfirmed must be >= params.UserVerifyTxn.MaxDropletPrecision (%d)", params.UserVerifyTxn.MaxDropletPrecision)
	}

	if c.Node.CreateBlockVerifyTxn.MaxDropletPrecision > droplet.Exponent {
		return fmt.Errorf("-max-decimals-create-block must be <= droplet.Exponent (%d)", droplet.Exponent)
	}
	if c.Node.CreateBlockVerifyTxn.MaxDropletPrecision < params.UserVerifyTxn.MaxDropletPrecision {
		return fmt.Errorf("-max-decimals-create-block must be >= params.UserVerifyTxn.MaxDropletPrecision (%d)", params.UserVerifyTxn.MaxDropletPrecision)
	}

	return nil
}

// buildAPISets builds the set of enable APIs by the following rules:
// * If EnableAll, all API sets are added
// * For each api set in EnabledAPISets, add
// * For each api set in DisabledAPISets, remove
func buildAPISets(c NodeConfig) (map[string]struct{}, error) {
	enabledAPISets := strings.Split(c.EnabledAPISets, ",")
	if err := validateAPISets("-enable-api-sets", enabledAPISets); err != nil {
		return nil, err
	}

	disabledAPISets := strings.Split(c.DisabledAPISets, ",")
	if err := validateAPISets("-disable-api-sets", disabledAPISets); err != nil {
		return nil, err
	}

	apiSets := make(map[string]struct{})

	allAPISets := []string{
		api.EndpointsRead,
		api.EndpointsStatus,
		api.EndpointsWallet,
		api.EndpointsTransaction,
		api.EndpointsNetCtrl,
		api.EndpointsStorage,
		// Do not include insecure or deprecated API sets, they must always
		// be explicitly enabled through -enable-api-sets
	}

	if c.EnableAllAPISets {
		for _, s := range allAPISets {
			apiSets[s] = struct{}{}
		}
	}

	// Add the enabled API sets
	for _, k := range enabledAPISets {
		apiSets[k] = struct{}{}
	}

	// Remove the disabled API sets
	for _, k := range disabledAPISets {
		delete(apiSets, k)
	}

	return apiSets, nil
}

func validateAPISets(opt string, apiSets []string) error {
	for _, k := range apiSets {
		k = strings.ToUpper(strings.TrimSpace(k))
		switch k {
		case api.EndpointsRead,
			api.EndpointsStatus,
			api.EndpointsTransaction,
			api.EndpointsWallet,
			api.EndpointsInsecureWalletSeed,
			api.EndpointsNetCtrl,
			api.EndpointsStorage:
		case "":
			continue
		default:
			return fmt.Errorf("Invalid value in %s: %q", opt, k)
		}
	}
	return nil
}

// RegisterFlags binds CLI flags to config values using Cobra
func (c *NodeConfig) RegisterFlags(cmd *cobra.Command) {
	//	cmd.Flags().BoolVar(&help, "help", false, "Show help")
	cmd.Flags().BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX, "disable PEX peer discovery")
	cmd.Flags().BoolVar(&c.DownloadPeerList, "download-peerlist", c.DownloadPeerList, "download a peers.txt from the peerlist URL")
	cmd.Flags().StringVar(&c.PeerListURL, "peerlist-url", c.PeerListURL, "URL to download peers.txt from (requires peerlist download enabled)")
	cmd.Flags().BoolVar(&c.DisableOutgoingConnections, "disable-outgoing", c.DisableOutgoingConnections, "Don't make outgoing connections")
	cmd.Flags().BoolVar(&c.DisableIncomingConnections, "disable-incoming", c.DisableIncomingConnections, "Don't allow incoming connections")
	cmd.Flags().BoolVar(&c.DisableNetworking, "disable-networking", c.DisableNetworking, "Disable all network activity")
	cmd.Flags().BoolVar(&c.LegacyPeerCompat, "legacy-peer-compat", c.LegacyPeerCompat, "Allow connections from legacy peers that don't send blockchain pubkey")
	cmd.Flags().BoolVar(&c.EnableGUI, "enable-gui", c.EnableGUI, "Enable GUI")
	cmd.Flags().BoolVar(&c.DisableCSRF, "disable-csrf", c.DisableCSRF, "disable CSRF check")
	cmd.Flags().BoolVar(&c.DisableHeaderCheck, "disable-header-check", c.DisableHeaderCheck, "disables the host, origin and referer header checks.")
	cmd.Flags().BoolVar(&c.DisableCSP, "disable-csp", c.DisableCSP, "disable Content Security Policy in http response")
	cmd.Flags().StringVar(&c.Address, "address", c.Address, "IP Address to run application on. Leave empty to default to a public interface")
	cmd.Flags().IntVar(&c.Port, "port", c.Port, "Port to run application on")

	cmd.Flags().BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	cmd.Flags().IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	cmd.Flags().StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")
	cmd.Flags().StringVar(&c.WebInterfaceCert, "web-interface-cert", c.WebInterfaceCert, "skycoind.cert file for web interface HTTPS. If not provided, will autogenerate or use skycoind.cert in data dir")
	cmd.Flags().StringVar(&c.WebInterfaceKey, "web-interface-key", c.WebInterfaceKey, "skycoind.key file for web interface HTTPS. If not provided, will autogenerate or use skycoind.key in data dir")
	cmd.Flags().BoolVar(&c.WebInterfaceHTTPS, "web-interface-https", c.WebInterfaceHTTPS, "enable HTTPS for web interface")
	cmd.Flags().StringVar(&c.HostWhitelist, "host-whitelist", c.HostWhitelist, "Hostnames to whitelist in the Host header check. Only applies when the web interface is bound to localhost.")

	allAPISets := []string{
		api.EndpointsRead,
		api.EndpointsStatus,
		api.EndpointsWallet,
		api.EndpointsTransaction,
		api.EndpointsNetCtrl,
		api.EndpointsInsecureWalletSeed,
		api.EndpointsStorage,
	}
	cmd.Flags().StringVar(&c.EnabledAPISets, "enable-api-sets", c.EnabledAPISets, fmt.Sprintf("enable API set. Options are %s. Multiple values should be separated by comma", strings.Join(allAPISets, ", ")))
	cmd.Flags().StringVar(&c.DisabledAPISets, "disable-api-sets", c.DisabledAPISets, fmt.Sprintf("disable API set. Options are %s. Multiple values should be separated by comma", strings.Join(allAPISets, ", ")))
	cmd.Flags().BoolVar(&c.EnableAllAPISets, "enable-all-api-sets", c.EnableAllAPISets, "enable all API sets except deprecated or insecure. Applied before the disable API sets flag.")

	cmd.Flags().StringVar(&c.WebInterfaceUsername, "web-interface-username", c.WebInterfaceUsername, "username for the web interface")
	cmd.Flags().StringVar(&c.WebInterfacePassword, "web-interface-password", c.WebInterfacePassword, "password for the web interface")
	cmd.Flags().BoolVar(&c.WebInterfacePlaintextAuth, "web-interface-plaintext-auth", c.WebInterfacePlaintextAuth, "allow web interface auth without https")

	cmd.Flags().BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
	cmd.Flags().StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, fmt.Sprintf("directory to store app data (defaults to %s)", c.DataDirectory))
	cmd.Flags().StringVar(&c.DBPath, "db-path", c.DBPath, "path of database file")
	cmd.Flags().BoolVar(&c.DBReadOnly, "db-read-only", c.DBReadOnly, "open bolt db in read only mode")
	cmd.Flags().BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	cmd.Flags().StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	cmd.Flags().BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "run the HTTP profiling interface")
	cmd.Flags().StringVar(&c.HTTPProfHost, "http-prof-host", c.HTTPProfHost, "hostname to bind the HTTP profiling interface to")
	cmd.Flags().StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	cmd.Flags().BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	cmd.Flags().BoolVar(&c.DisablePingPong, "no-ping-log", c.DisablePingPong, `disable "reply to ping" and "received pong" debug log messages`)
	cmd.Flags().BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")
	cmd.Flags().StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory, "static content directory for the HTML interface")

	cmd.Flags().BoolVar(&c.VerifyDB, "verify-db", c.VerifyDB, "check the database for corruption")
	cmd.Flags().BoolVar(&c.ResetCorruptDB, "reset-corrupt-db", c.ResetCorruptDB, "reset the database if corrupted, and continue running instead of exiting")

	cmd.Flags().BoolVar(&c.DisableDefaultPeers, "disable-default-peers", c.DisableDefaultPeers, "disable the hardcoded default peers")
	cmd.Flags().StringVar(&c.CustomPeersFile, "custom-peers-file", c.CustomPeersFile, "load custom peers from a newline separate list of ip:port in a file. Note that this is different from the peers.json file in the data directory")

	cmd.Flags().StringVar(&c.UserAgentRemark, "user-agent-remark", c.UserAgentRemark, "additional remark to include in the user agent sent over the wire protocol")

	cmd.Flags().Uint64Var(&c.maxUnconfirmedTransactionSize, "max-txn-size-unconfirmed", uint64(c.UnconfirmedVerifyTxn.MaxTransactionSize), "maximum size of an unconfirmed transaction")
	cmd.Flags().Uint64Var(&c.unconfirmedBurnFactor, "burn-factor-unconfirmed", uint64(c.UnconfirmedVerifyTxn.BurnFactor), "coinhour burn factor applied to unconfirmed transactions")
	cmd.Flags().Uint64Var(&c.unconfirmedMaxDropletPrecision, "max-decimals-unconfirmed", uint64(c.UnconfirmedVerifyTxn.MaxDropletPrecision), "max number of decimal places applied to unconfirmed transactions")
	cmd.Flags().Uint64Var(&c.createBlockBurnFactor, "burn-factor-create-block", uint64(c.CreateBlockVerifyTxn.BurnFactor), "coinhour burn factor applied when creating blocks")
	cmd.Flags().Uint64Var(&c.createBlockMaxTransactionSize, "max-txn-size-create-block", uint64(c.CreateBlockVerifyTxn.MaxTransactionSize), "maximum size of a transaction applied when creating blocks")
	cmd.Flags().Uint64Var(&c.createBlockMaxDropletPrecision, "max-decimals-create-block", uint64(c.CreateBlockVerifyTxn.MaxDropletPrecision), "max number of decimal places applied when creating blocks")
	cmd.Flags().Uint64Var(&c.maxBlockSize, "max-block-size", uint64(c.MaxBlockTransactionsSize), "maximum total size of transactions in a block")
	cmd.Flags().Uint64Var(&c.MaxLastBlocksCount, "max-last-blocks-count", c.MaxLastBlocksCount, "Maximum number of blocks to response for API /api/v1/last_blocks")

	cmd.Flags().BoolVar(&c.RunBlockPublisher, "block-publisher", c.RunBlockPublisher, "run the daemon as a block publisher")
	cmd.Flags().StringVar(&c.BlockchainPubkeyStr, "blockchain-public-key", c.BlockchainPubkeyStr, "public key of the blockchain")
	cmd.Flags().StringVar(&c.BlockchainSeckeyStr, "blockchain-secret-key", c.BlockchainSeckeyStr, "secret key of the blockchain")

	cmd.Flags().StringVar(&c.GenesisAddressStr, "genesis-address", c.GenesisAddressStr, "genesis address")
	cmd.Flags().StringVar(&c.GenesisSignatureStr, "genesis-signature", c.GenesisSignatureStr, "genesis block signature")
	cmd.Flags().Uint64Var(&c.GenesisTimestamp, "genesis-timestamp", c.GenesisTimestamp, "genesis block timestamp")

	cmd.Flags().StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory, "location of the wallet files")
	cmd.Flags().StringVar(&c.KVStorageDirectory, "storage-dir", c.KVStorageDirectory, "location of the storage data files")
	cmd.Flags().IntVar(&c.MaxConnections, "max-connections", c.MaxConnections, "Maximum number of total connections allowed")
	cmd.Flags().IntVar(&c.MaxOutgoingConnections, "max-outgoing-connections", c.MaxOutgoingConnections, "Maximum number of outgoing connections allowed")
	cmd.Flags().IntVar(&c.MaxIncomingConnections, "max-incoming-connections", c.MaxIncomingConnections, "Maximum number of incoming connections allowed")
	cmd.Flags().IntVar(&c.MaxDefaultPeerOutgoingConnections, "max-default-peer-outgoing-connections", c.MaxDefaultPeerOutgoingConnections, "The maximum default peer outgoing connections allowed")
	cmd.Flags().IntVar(&c.PeerlistSize, "peerlist-size", c.PeerlistSize, "Max number of peers to track in peerlist")
	cmd.Flags().DurationVar(&c.OutgoingConnectionsRate, "connection-rate", c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	cmd.Flags().IntVar(&c.MaxOutgoingMessageLength, "max-out-msg-len", c.MaxOutgoingMessageLength, "Maximum length of outgoing wire messages")
	cmd.Flags().IntVar(&c.MaxIncomingMessageLength, "max-in-msg-len", c.MaxIncomingMessageLength, "Maximum length of incoming wire messages")
	cmd.Flags().BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly, "Run on localhost and only connect to localhost peers")
	cmd.Flags().StringVar(&c.WalletCryptoType, "wallet-crypto-type", c.WalletCryptoType, "wallet encryption type (see default for format options)\n\r")
	cmd.Flags().BoolVar(&c.Version, "version", false, "show node version")

	// Display/Branding flags
	cmd.Flags().StringVar(&c.Fiber.Name, "coin-name", c.Fiber.Name, "name of the coin")
	cmd.Flags().StringVar(&c.Fiber.Ticker, "ticker", c.Fiber.Ticker, "coin ticker symbol (e.g., SKY)")
	cmd.Flags().StringVar(&c.Fiber.DisplayName, "display-name", c.Fiber.DisplayName, "display name of the coin")
	cmd.Flags().StringVar(&c.Fiber.CoinHoursName, "coin-hours-name", c.Fiber.CoinHoursName, "display name for coin hours")
	cmd.Flags().StringVar(&c.Fiber.CoinHoursNameSingular, "coin-hours-name-singular", c.Fiber.CoinHoursNameSingular, "singular display name for coin hours")
	cmd.Flags().StringVar(&c.Fiber.CoinHoursTicker, "coin-hours-ticker", c.Fiber.CoinHoursTicker, "ticker symbol for coin hours")
	cmd.Flags().StringVar(&c.Fiber.QrURIPrefix, "qr-uri-prefix", c.Fiber.QrURIPrefix, "prefix for QR code URIs")
	cmd.Flags().StringVar(&c.Fiber.ExplorerURL, "explorer-url", c.Fiber.ExplorerURL, "URL of the block explorer")
	cmd.Flags().StringVar(&c.Fiber.VersionURL, "version-url", c.Fiber.VersionURL, "URL for version checking")
	cmd.Flags().Uint32Var((*uint32)(&c.Fiber.Bip44Coin), "bip44-coin", uint32(c.Fiber.Bip44Coin), "BIP44 coin type")
}

func (c *NodeConfig) applyConfigMode(configMode string) {
	if runtime.GOOS == "windows" {
		c.ColorLog = false
	}
	switch configMode {
	case "":
	case "STANDALONE_CLIENT":
		c.EnableAllAPISets = true
		c.EnabledAPISets = api.EndpointsInsecureWalletSeed
		c.EnableGUI = true
		c.LaunchBrowser = true
		c.DisableCSRF = false
		c.DisableHeaderCheck = false
		c.DisableCSP = false
		c.DownloadPeerList = true
		c.WebInterface = true
		c.LogToFile = false
		c.ResetCorruptDB = true
		c.WebInterfacePort = 0 // randomize web interface port
	default:
		panic("Invalid ConfigMode")
	}
}

// LoadFromFiberConfig loads configuration from a fiber.toml file
// and overrides the default values in NodeConfig
func (c *NodeConfig) LoadFromFiberConfig(configPath string) error {
	if configPath == "" {
		return nil // No config file specified
	}

	// Store path for later writing
	c.fiberTomlPath = configPath

	// Load fiber config
	fiberCfg, err := fiber.NewConfig(filepath.Base(configPath), filepath.Dir(configPath))
	if err != nil {
		return fmt.Errorf("failed to load fiber config: %w", err)
	}

	// Map fiber.NodeConfig to skycoin.NodeConfig
	c.applyFiberNodeConfig(fiberCfg.Node)

	// Store distribution parameters from fiber.ParamsConfig
	if len(fiberCfg.Params.DistributionAddresses) > 0 {
		c.distributionAddresses = fiberCfg.Params.DistributionAddresses
		c.distributionMaxCoinSupply = fiberCfg.Params.MaxCoinSupply
		c.distributionUnlockedCount = fiberCfg.Params.InitialUnlockedCount
		c.distributionUnlockRate = fiberCfg.Params.UnlockAddressRate
		c.distributionUnlockInterval = fiberCfg.Params.UnlockTimeInterval
	}

	return nil
}

// LoadFromGenesisWallet loads genesis credentials from a genesis wallet JSON file
// This takes precedence over fiber.toml values for address, pubkey, and seckey
func (c *NodeConfig) LoadFromGenesisWallet(walletPath string) error {
	if walletPath == "" {
		return nil
	}

	// Store path for later use
	c.genesisWalletPath = walletPath

	// Read the genesis wallet file
	data, err := os.ReadFile(walletPath) //nolint:gosec // G304: User-specified wallet path is intentional
	if err != nil {
		return fmt.Errorf("failed to read genesis wallet: %w", err)
	}

	// Parse the wallet JSON
	var wallet struct {
		Entries []struct {
			Address   string `json:"address"`
			PublicKey string `json:"public_key"`
			SecretKey string `json:"secret_key"`
		} `json:"entries"`
	}

	if err := json.Unmarshal(data, &wallet); err != nil {
		return fmt.Errorf("failed to parse genesis wallet JSON: %w", err)
	}

	if len(wallet.Entries) == 0 {
		return fmt.Errorf("genesis wallet has no entries")
	}

	// Use the first entry
	entry := wallet.Entries[0]

	// Set genesis address and blockchain keys
	c.GenesisAddressStr = entry.Address
	c.BlockchainPubkeyStr = entry.PublicKey
	c.BlockchainSeckeyStr = entry.SecretKey

	// Clear the genesis signature since it's not valid for this wallet
	// The signature will be generated when the genesis block is created
	c.GenesisSignatureStr = ""

	return nil
}

// applyFiberNodeConfig maps fiber.NodeConfig fields to NodeConfig
func (c *NodeConfig) applyFiberNodeConfig(node fiber.NodeConfig) {
	// Core blockchain parameters
	if node.CoinName != "" {
		c.CoinName = node.CoinName
		c.Fiber.Name = node.CoinName
	}
	if node.Port != 0 {
		c.Port = node.Port
	}
	if node.WebInterfacePort != 0 {
		c.WebInterfacePort = node.WebInterfacePort
	}
	if node.GenesisSignatureStr != "" {
		c.GenesisSignatureStr = node.GenesisSignatureStr
	}
	if node.GenesisAddressStr != "" {
		c.GenesisAddressStr = node.GenesisAddressStr
	}
	if node.BlockchainPubkeyStr != "" {
		c.BlockchainPubkeyStr = node.BlockchainPubkeyStr
	}
	if node.BlockchainSeckeyStr != "" {
		c.BlockchainSeckeyStr = node.BlockchainSeckeyStr
	}
	if node.GenesisTimestamp != 0 {
		c.GenesisTimestamp = node.GenesisTimestamp
	}
	if node.GenesisCoinVolume != 0 {
		c.GenesisCoinVolume = node.GenesisCoinVolume
	}
	if len(node.DefaultConnections) > 0 {
		c.DefaultConnections = node.DefaultConnections
	}
	if node.PeerListURL != "" {
		c.PeerListURL = node.PeerListURL
	}

	// Data directory - expand $HOME
	// If not explicitly set, derive from coin name or display name
	if node.DataDirectory != "" {
		dataDir := node.DataDirectory
		home := file.UserHome()
		dataDir = replaceHome(dataDir, home)
		c.DataDirectory = dataDir
	} else if c.DataDirectory == "$HOME/.skycoin" {
		// Auto-derive data directory from coin name or display name (matches newcoin behavior)
		home := file.UserHome()
		derivedName := ""
		if node.CoinName != "" {
			derivedName = node.CoinName
		} else if node.DisplayName != "" {
			derivedName = node.DisplayName
		}
		if derivedName != "" {
			c.DataDirectory = replaceHome("$HOME/."+strings.ToLower(derivedName), home)
		}
	}

	// Transaction verification params
	if node.UnconfirmedBurnFactor != 0 {
		c.UnconfirmedVerifyTxn.BurnFactor = node.UnconfirmedBurnFactor
		c.unconfirmedBurnFactor = uint64(node.UnconfirmedBurnFactor)
	}
	if node.UnconfirmedMaxTransactionSize != 0 {
		c.UnconfirmedVerifyTxn.MaxTransactionSize = node.UnconfirmedMaxTransactionSize
		c.maxUnconfirmedTransactionSize = uint64(node.UnconfirmedMaxTransactionSize)
	}
	if node.UnconfirmedMaxDropletPrecision != 0 {
		c.UnconfirmedVerifyTxn.MaxDropletPrecision = node.UnconfirmedMaxDropletPrecision
		c.unconfirmedMaxDropletPrecision = uint64(node.UnconfirmedMaxDropletPrecision)
	}
	if node.CreateBlockBurnFactor != 0 {
		c.CreateBlockVerifyTxn.BurnFactor = node.CreateBlockBurnFactor
		c.createBlockBurnFactor = uint64(node.CreateBlockBurnFactor)
	}
	if node.CreateBlockMaxTransactionSize != 0 {
		c.CreateBlockVerifyTxn.MaxTransactionSize = node.CreateBlockMaxTransactionSize
		c.createBlockMaxTransactionSize = uint64(node.CreateBlockMaxTransactionSize)
	}
	if node.CreateBlockMaxDropletPrecision != 0 {
		c.CreateBlockVerifyTxn.MaxDropletPrecision = node.CreateBlockMaxDropletPrecision
		c.createBlockMaxDropletPrecision = uint64(node.CreateBlockMaxDropletPrecision)
	}
	if node.MaxBlockTransactionsSize != 0 {
		c.MaxBlockTransactionsSize = node.MaxBlockTransactionsSize
		c.maxBlockSize = uint64(node.MaxBlockTransactionsSize)
	}

	// Display/Branding
	if node.Ticker != "" {
		c.Fiber.Ticker = node.Ticker
	}
	if node.DisplayName != "" {
		c.Fiber.DisplayName = node.DisplayName
	}
	if node.CoinHoursName != "" {
		c.Fiber.CoinHoursName = node.CoinHoursName
	}
	if node.CoinHoursNameSingular != "" {
		c.Fiber.CoinHoursNameSingular = node.CoinHoursNameSingular
	}
	if node.CoinHoursTicker != "" {
		c.Fiber.CoinHoursTicker = node.CoinHoursTicker
	}
	if node.QrURIPrefix != "" {
		c.Fiber.QrURIPrefix = node.QrURIPrefix
	}
	if node.ExplorerURL != "" {
		c.Fiber.ExplorerURL = node.ExplorerURL
	}
	if node.VersionURL != "" {
		c.Fiber.VersionURL = node.VersionURL
	}
	if node.Bip44Coin != 0 {
		c.Fiber.Bip44Coin = node.Bip44Coin
	}
}

// WriteFiberTomlGenesis writes genesis address, pubkey, and signature to fiber.toml
// This is called after the genesis block is created to persist the values
func (c *NodeConfig) WriteFiberTomlGenesis(signature string) error {
	if c.fiberTomlPath == "" {
		return nil // No fiber.toml path set, nothing to write
	}

	// Read the current fiber.toml
	data, err := os.ReadFile(c.fiberTomlPath)
	if err != nil {
		return fmt.Errorf("failed to read fiber.toml: %w", err)
	}

	// Parse as map
	var tomlMap map[string]interface{}
	if err := toml.Unmarshal(data, &tomlMap); err != nil {
		return fmt.Errorf("failed to parse fiber.toml: %w", err)
	}

	// Get or create [node] section
	nodeSection, ok := tomlMap["node"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing [node] section in fiber.toml")
	}

	// Update genesis fields
	nodeSection["genesis_address_str"] = c.GenesisAddressStr
	nodeSection["blockchain_pubkey_str"] = c.BlockchainPubkeyStr
	nodeSection["genesis_signature_str"] = signature

	// Write back to file
	updatedData, err := toml.Marshal(tomlMap)
	if err != nil {
		return fmt.Errorf("failed to marshal fiber.toml: %w", err)
	}

	if err := os.WriteFile(c.fiberTomlPath, updatedData, 0600); err != nil {
		return fmt.Errorf("failed to write fiber.toml: %w", err)
	}

	return nil
}

func panicIfError(err error, msg string, args ...interface{}) { //nolint:unparam
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}
