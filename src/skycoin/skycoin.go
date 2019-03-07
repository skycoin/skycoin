/*
Package skycoin implements the main daemon cmd's configuration and setup
*/
package skycoin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/toqueteos/webbrowser"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/apputil"
	"github.com/skycoin/skycoin/src/util/certutil"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	// DBVerifyCheckpointVersion is a checkpoint for determining if DB verification should be run.
	// Any DB upgrading from less than this version to equal or higher than this version will be forced to verify.
	// Update this version checkpoint if a newer version requires a new verification run.
	DBVerifyCheckpointVersion       = "0.25.0"
	dbVerifyCheckpointVersionParsed semver.Version
)

// Coin represents a fiber coin instance
type Coin struct {
	config Config
	logger *logging.Logger
}

// Run starts the node
func (c *Coin) Run() error {
	var (
		db *dbutil.DB
		w  *wallet.Service
		v  *visor.Visor
		d  *daemon.Daemon
		gw *api.Gateway
		wi *api.Server

		errC = make(chan error, 10)

		retErr error
	)

	if c.config.Node.Version {
		fmt.Println(c.config.Build.Version)
		return nil
	}

	logLevel, err := logging.LevelFromString(c.config.Node.LogLevel)
	if err != nil {
		err = fmt.Errorf("Invalid -log-level: %v", err)
		c.logger.Error(err)
		return err
	}

	logging.SetLevel(logLevel)

	if c.config.Node.ColorLog {
		logging.EnableColors()
	} else {
		logging.DisableColors()
	}

	var logFile *os.File
	if c.config.Node.LogToFile {
		var err error
		logFile, err = c.initLogFile()
		if err != nil {
			c.logger.Error(err)
			return err
		}
	}

	var fullAddress string
	scheme := "http"
	if c.config.Node.WebInterfaceHTTPS {
		scheme = "https"
	}
	host := fmt.Sprintf("%s:%d", c.config.Node.WebInterfaceAddr, c.config.Node.WebInterfacePort)

	if c.config.Node.ProfileCPU {
		f, err := os.Create(c.config.Node.ProfileCPUFile)
		if err != nil {
			c.logger.Error(err)
			return err
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			c.logger.Error(err)
			return err
		}
		defer pprof.StopCPUProfile()
	}

	if c.config.Node.HTTPProf {
		go func() {
			if err := http.ListenAndServe(c.config.Node.HTTPProfHost, nil); err != nil {
				c.logger.WithError(err).Errorf("Listen on HTTP profiling interface %s failed", c.config.Node.HTTPProfHost)
			}
		}()
	}

	var wg sync.WaitGroup

	quit := make(chan struct{})

	// Catch SIGINT (CTRL-C) (closes the quit channel)
	go apputil.CatchInterrupt(quit)

	// Catch SIGUSR1 (prints runtime stack to stdout)
	go apputil.CatchDebug()

	// Parse the current app version
	appVersion, err := c.config.Build.Semver()
	if err != nil {
		c.logger.WithError(err).Errorf("Version %s is not a valid semver", c.config.Build.Version)
		return err
	}

	c.logger.Infof("App version: %s", appVersion)
	c.logger.Infof("OS: %s", runtime.GOOS)
	c.logger.Infof("Arch: %s", runtime.GOARCH)

	c.logger.Infof("Coinhour burn factor for user transactions is %d", params.UserVerifyTxn.BurnFactor)
	c.logger.Infof("Max transaction size for user transactions is %d", params.UserVerifyTxn.MaxTransactionSize)
	c.logger.Infof("Max decimals for user transactions is %d", params.UserVerifyTxn.MaxDropletPrecision)

	if db, retErr = c.initDB(appVersion, quit); retErr != nil {
		c.logger.Error(err)
		goto earlyShutdown
	}

	if w, retErr = c.initWallet(); retErr != nil {
		c.logger.Error(retErr)
		goto earlyShutdown
	}

	if v, retErr = c.initVisor(db, w); retErr != nil {
		c.logger.Error(retErr)
		goto earlyShutdown
	}

	if d, retErr = c.initDaemon(v); retErr != nil {
		c.logger.Error(retErr)
		goto earlyShutdown
	}

	gw = api.NewGateway(d, v, w)

	if c.config.Node.WebInterface {
		if wi, retErr = c.createGUI(gw, host); retErr != nil {
			c.logger.Error(retErr)
			goto earlyShutdown
		}

		fullAddress = fmt.Sprintf("%s://%s", scheme, wi.Addr())
		c.logger.Critical().Infof("Full address: %s", fullAddress)
		if c.config.Node.PrintWebInterfaceAddress {
			fmt.Println(fullAddress)
		}
	}

	if retErr = v.Init(); retErr != nil {
		c.logger.Error(retErr)
		goto earlyShutdown
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := d.Run(); err != nil {
			c.logger.Error(err)
			errC <- err
		}
	}()

	if c.config.Node.WebInterface {
		cancelLaunchBrowser := make(chan struct{})

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := wi.Serve(); err != nil {
				close(cancelLaunchBrowser)
				c.logger.Error(err)
				errC <- err
			}
		}()

		if c.config.Node.LaunchBrowser {
			go func() {
				select {
				case <-cancelLaunchBrowser:
					c.logger.Warning("Browser launching cancelled")

					// Wait a moment just to make sure the http interface is up
				case <-time.After(time.Millisecond * 100):
					c.logger.Infof("Launching System Browser with %s", fullAddress)
					if err := webbrowser.Open(fullAddress); err != nil {
						c.logger.Error(err)
					}
				}
			}()
		}
	}

	select {
	case <-quit:
	case retErr = <-errC:
		c.logger.Error(retErr)
	}

	c.logger.Info("Shutting down...")

	if wi != nil {
		c.logger.Info("Closing web interface")
		wi.Shutdown()
	}

	c.logger.Info("Closing daemon")
	d.Shutdown()

	c.logger.Info("Waiting for goroutines to finish")
	wg.Wait()

earlyShutdown:
	if db != nil {
		c.logger.Info("Closing database")
		if err := db.Close(); err != nil {
			c.logger.WithError(err).Error("Failed to close DB")
		}
	}

	c.logger.Info("Goodbye")

	if logFile != nil {
		if err := logFile.Close(); err != nil {
			fmt.Println("Failed to close log file")
		}
	}

	return retErr
}

// NewCoin returns a new fiber coin instance
func NewCoin(config Config, logger *logging.Logger) *Coin {
	return &Coin{
		config: config,
		logger: logger,
	}
}

func (c *Coin) initLogFile() (*os.File, error) {
	logDir := filepath.Join(c.config.Node.DataDirectory, "logs")
	if err := createDirIfNotExist(logDir); err != nil {
		c.logger.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
		return nil, fmt.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
	}

	// open log file
	tf := "2006-01-02-030405"
	logfile := filepath.Join(logDir, fmt.Sprintf("%s-v%s.log", time.Now().Format(tf), c.config.Build.Version))

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		c.logger.Errorf("os.OpenFile(%s) failed: %v", logfile, err)
		return nil, err
	}

	hook := logging.NewWriteHook(f)
	logging.AddHook(hook)

	return f, nil
}

// initDB creates and initlializes new instance of *dbutil.DB
func (c *Coin) initDB(
	appVersion *semver.Version, quit chan struct{}) (*dbutil.DB, error) {

	// Open the database
	c.logger.Infof("Opening database %s", c.config.Node.DBPath)
	db, err := visor.OpenDB(c.config.Node.DBPath, c.config.Node.DBReadOnly)
	if err != nil {
		c.logger.Errorf("Database failed to open: %v. Is another skycoin instance running?", err)
		return nil, err
	}

	// Look for saved app version
	dbVersion, err := visor.GetDBVersion(db)
	if err != nil {
		c.logger.WithError(err).Error("visor.GetDBVersion failed")
		return nil, err
	}

	if dbVersion == nil {
		c.logger.Info("DB version not found in DB")
	} else {
		c.logger.Infof("DB version: %s", dbVersion)
	}

	c.logger.Infof("DB verify checkpoint version: %s", DBVerifyCheckpointVersion)

	// If the saved DB version is higher than the app version, abort.
	// Otherwise DB corruption could occur.
	if dbVersion != nil && dbVersion.GT(*appVersion) {
		err = fmt.Errorf("Cannot use newer DB version=%v with older software version=%v", dbVersion, appVersion)
		c.logger.WithError(err).Error()
		return nil, err
	}

	// Verify the DB if the version detection says to, or if it was requested on the command line
	if shouldVerifyDB(appVersion, dbVersion) || c.config.Node.VerifyDB {
		if c.config.Node.ResetCorruptDB {
			// Check the database integrity and recreate it if necessary
			c.logger.Info("Checking database and resetting if corrupted")
			newDB, err := visor.ResetCorruptDB(db, c.config.Node.blockchainPubkey, quit)
			if err != nil {
				if err != visor.ErrVerifyStopped {
					c.logger.Errorf("visor.ResetCorruptDB failed: %v", err)
				}
				return nil, err
			}
			db = newDB
		} else {
			c.logger.Info("Checking database")
			if err := visor.CheckDatabase(db, c.config.Node.blockchainPubkey, quit); err != nil {
				if err != visor.ErrVerifyStopped {
					c.logger.Errorf("visor.CheckDatabase failed: %v", err)
				}
				return nil, err
			}
		}
	}

	// Update the DB version
	if !db.IsReadOnly() {
		if err := visor.SetDBVersion(db, *appVersion); err != nil {
			c.logger.WithError(err).Error("visor.SetDBVersion failed")
			return nil, err
		}
	}

	return db, nil
}

func shouldVerifyDB(appVersion, dbVersion *semver.Version) bool {
	// If the dbVersion is not set, verify
	if dbVersion == nil {
		return true
	}

	// If the dbVersion is less than the verification checkpoint version
	// and the appVersion is greater than or equal to the checkpoint version,
	// verify
	if dbVersion.LT(dbVerifyCheckpointVersionParsed) && appVersion.GTE(dbVerifyCheckpointVersionParsed) {
		return true
	}

	return false
}

// initWallet creates new instance of *wallets.Servcie
func (c *Coin) initWallet() (*wallet.Service, error) {
	cfg := c.configureWallet()

	s, err := wallet.NewService(cfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// configureWallet sets the wallet config values
func (c *Coin) configureWallet() wallet.Config {
	wc := wallet.NewConfig()

	wc.WalletDir = c.config.Node.WalletDirectory
	_, wc.EnableWalletAPI = c.config.Node.enabledAPISets[api.EndpointsWallet]
	_, wc.EnableSeedAPI = c.config.Node.enabledAPISets[api.EndpointsInsecureWalletSeed]

	// Initialize wallet default crypto type
	cryptoType, err := wallet.CryptoTypeFromString(c.config.Node.WalletCryptoType)
	if err != nil {
		log.Panic(err)
	}

	wc.CryptoType = cryptoType

	return wc
}

// initVisor creates new instance of *visor.Visor
func (c *Coin) initVisor(
	db *dbutil.DB, ws *wallet.Service) (*visor.Visor, error) {

	cfg := c.configureVisor()

	if !db.IsReadOnly() {
		if err := visor.CreateBuckets(db); err != nil {
			c.logger.WithError(err).Error("CreateBuckets failed")
			return nil, err
		}
	}

	bc, err := visor.NewBlockchain(db, visor.BlockchainConfig{
		Pubkey:      cfg.BlockchainPubkey,
		Arbitrating: cfg.Arbitrating,
	})
	if err != nil {
		return nil, err
	}

	history := historydb.New()

	if !db.IsReadOnly() {
		if err := db.Update("build unspent indexes and init history", func(tx *dbutil.Tx) error {
			headSeq, _, err := bc.HeadSeq(tx)
			if err != nil {
				return err
			}

			if err := bc.Unspent().MaybeBuildIndexes(tx, headSeq); err != nil {
				return err
			}

			return visor.InitHistory(tx, bc, history)
		}); err != nil {
			return nil, err
		}
	}

	utp, err := visor.NewUnconfirmedTransactionPool(db)
	if err != nil {
		return nil, err
	}

	v, err := visor.NewWith(cfg, db, utp, bc, history, ws)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// configureVisor sets the visor config values
func (c *Coin) configureVisor() visor.Config {
	vc := visor.NewConfig()

	vc.IsBlockPublisher = c.config.Node.RunBlockPublisher

	vc.BlockchainPubkey = c.config.Node.blockchainPubkey
	vc.BlockchainSeckey = c.config.Node.blockchainSeckey

	vc.UnconfirmedVerifyTxn = c.config.Node.UnconfirmedVerifyTxn
	vc.CreateBlockVerifyTxn = c.config.Node.CreateBlockVerifyTxn
	vc.MaxBlockTransactionsSize = c.config.Node.MaxBlockTransactionsSize

	vc.GenesisAddress = c.config.Node.genesisAddress
	vc.GenesisSignature = c.config.Node.genesisSignature
	vc.GenesisTimestamp = c.config.Node.GenesisTimestamp
	vc.GenesisCoinVolume = c.config.Node.GenesisCoinVolume
	vc.Arbitrating = c.config.Node.Arbitrating

	return vc
}

// initDaemon creates new instance of *daemon.Daemon
func (c *Coin) initDaemon(v *visor.Visor) (*daemon.Daemon, error) {
	cfg := c.configureDaemon()

	d, err := daemon.New(cfg, v)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// configureDaemon sets the daemon config values
func (c *Coin) configureDaemon() daemon.Config {
	dc := daemon.NewConfig()

	dc.Pool.DefaultConnections = c.config.Node.DefaultConnections
	dc.Pool.MaxDefaultPeerOutgoingConnections = c.config.Node.MaxDefaultPeerOutgoingConnections
	dc.Pool.MaxIncomingMessageLength = c.config.Node.MaxIncomingMessageLength
	dc.Pool.MaxOutgoingMessageLength = c.config.Node.MaxOutgoingMessageLength

	dc.Pex.DataDirectory = c.config.Node.DataDirectory
	dc.Pex.Disabled = c.config.Node.DisablePEX
	dc.Pex.NetworkDisabled = c.config.Node.DisableNetworking
	dc.Pex.Max = c.config.Node.PeerlistSize
	dc.Pex.DownloadPeerList = c.config.Node.DownloadPeerList
	dc.Pex.PeerListURL = c.config.Node.PeerListURL
	dc.Pex.DisableTrustedPeers = c.config.Node.DisableDefaultPeers
	dc.Pex.CustomPeersFile = c.config.Node.CustomPeersFile
	dc.Pex.DefaultConnections = c.config.Node.DefaultConnections

	dc.Daemon.MaxOutgoingMessageLength = uint64(c.config.Node.MaxOutgoingMessageLength)
	dc.Daemon.MaxIncomingMessageLength = uint64(c.config.Node.MaxIncomingMessageLength)
	dc.Daemon.MaxBlockTransactionsSize = c.config.Node.MaxBlockTransactionsSize
	dc.Daemon.DefaultConnections = c.config.Node.DefaultConnections
	dc.Daemon.DisableOutgoingConnections = c.config.Node.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = c.config.Node.DisableIncomingConnections
	dc.Daemon.DisableNetworking = c.config.Node.DisableNetworking
	dc.Daemon.Port = c.config.Node.Port
	dc.Daemon.Address = c.config.Node.Address
	dc.Daemon.LocalhostOnly = c.config.Node.LocalhostOnly
	dc.Daemon.MaxConnections = c.config.Node.MaxConnections
	dc.Daemon.MaxOutgoingConnections = c.config.Node.MaxOutgoingConnections
	dc.Daemon.DataDirectory = c.config.Node.DataDirectory
	dc.Daemon.LogPings = !c.config.Node.DisablePingPong
	dc.Daemon.BlockchainPubkey = c.config.Node.blockchainPubkey
	dc.Daemon.UserAgent = c.config.Node.userAgent
	dc.Daemon.UnconfirmedVerifyTxn = c.config.Node.UnconfirmedVerifyTxn

	if c.config.Node.OutgoingConnectionsRate == 0 {
		c.config.Node.OutgoingConnectionsRate = time.Millisecond
	}
	dc.Daemon.OutgoingRate = c.config.Node.OutgoingConnectionsRate

	return dc
}

func (c *Coin) createGUI(gw *api.Gateway, host string) (*api.Server, error) {
	config := api.Config{
		StaticDir:          c.config.Node.GUIDirectory,
		DisableCSRF:        c.config.Node.DisableCSRF,
		DisableHeaderCheck: c.config.Node.DisableHeaderCheck,
		DisableCSP:         c.config.Node.DisableCSP,
		EnableGUI:          c.config.Node.EnableGUI,
		ReadTimeout:        c.config.Node.HTTPReadTimeout,
		WriteTimeout:       c.config.Node.HTTPWriteTimeout,
		IdleTimeout:        c.config.Node.HTTPIdleTimeout,
		EnabledAPISets:     c.config.Node.enabledAPISets,
		HostWhitelist:      c.config.Node.hostWhitelist,
		Health: api.HealthConfig{
			BuildInfo: readable.BuildInfo{
				Version: c.config.Build.Version,
				Commit:  c.config.Build.Commit,
				Branch:  c.config.Build.Branch,
			},
			CoinName:        c.config.Node.CoinName,
			DaemonUserAgent: c.config.Node.userAgent,
		},
		Username: c.config.Node.WebInterfaceUsername,
		Password: c.config.Node.WebInterfacePassword,
	}

	var s *api.Server
	if c.config.Node.WebInterfaceHTTPS {
		// Verify cert/key parameters, and if neither exist, create them
		exists, err := checkCertFiles(c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey)
		if err != nil {
			c.logger.Errorf("checkCertFiles failed: %v", err)
			return nil, err
		}

		if !exists {
			c.logger.Infof("Autogenerating HTTP certificate and key files %s, %s", c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey)
			if err := createCertFiles(c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey); err != nil {
				c.logger.Errorf("createCertFiles failed: %v", err)
				return nil, err
			}

			c.logger.Infof("Created cert file %s", c.config.Node.WebInterfaceCert)
			c.logger.Infof("Created key file %s", c.config.Node.WebInterfaceKey)
		}

		s, err = api.CreateHTTPS(host, config, gw, c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey)
		if err != nil {
			c.logger.Errorf("Failed to start web GUI: %v", err)
			return nil, err
		}
	} else {
		var err error
		s, err = api.Create(host, config, gw)
		if err != nil {
			c.logger.Errorf("Failed to start web GUI: %v", err)
			return nil, err
		}
	}

	return s, nil
}

// checkCertFiles returns true if both cert and key files exist, false if neither exist,
// or returns an error if only one does not exist
func checkCertFiles(cert, key string) (bool, error) {
	doesFileExist := func(f string) (bool, error) {
		if _, err := os.Stat(f); err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}

	certExists, err := doesFileExist(cert)
	if err != nil {
		return false, err
	}

	keyExists, err := doesFileExist(key)
	if err != nil {
		return false, err
	}

	switch {
	case certExists && keyExists:
		return true, nil
	case !certExists && !keyExists:
		return false, nil
	case certExists && !keyExists:
		return false, fmt.Errorf("certfile %s exists but keyfile %s does not", cert, key)
	case !certExists && keyExists:
		return false, fmt.Errorf("keyfile %s exists but certfile %s does not", key, cert)
	default:
		log.Panic("unreachable code")
		return false, errors.New("unreachable code")
	}
}

func createCertFiles(certFile, keyFile string) error {
	org := "skycoin daemon autogenerated cert"
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	cert, key, err := certutil.NewTLSCertPair(org, validUntil, nil)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(certFile, cert, 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(keyFile, key, 0600); err != nil {
		os.Remove(certFile)
		return err
	}

	return nil
}

// ParseConfig prepare the config
func (c *Coin) ParseConfig() error {
	return c.config.postProcess()
}

// InitTransaction creates the genesis transaction
func InitTransaction(uxID string, genesisSecKey cipher.SecKey) coin.Transaction {
	var txn coin.Transaction

	output := cipher.MustSHA256FromHex(uxID)
	if err := txn.PushInput(output); err != nil {
		log.Panic(err)
	}

	addrs := params.GetDistributionAddresses()

	if len(addrs) != 100 {
		log.Panic("Should have 100 distribution addresses")
	}

	// 1 million per address, measured in droplets
	if params.DistributionAddressInitialBalance != 1e6 {
		log.Panic("params.DistributionAddressInitialBalance expected to be 1e6*1e6")
	}

	for i := range addrs {
		addr := cipher.MustDecodeBase58Address(addrs[i])
		if err := txn.PushOutput(addr, params.DistributionAddressInitialBalance*1e6, 1); err != nil {
			log.Panic(err)
		}
	}

	seckeys := make([]cipher.SecKey, 1)
	seckey := genesisSecKey.Hex()
	seckeys[0] = cipher.MustSecKeyFromHex(seckey)
	txn.SignInputs(seckeys)

	if err := txn.UpdateHeader(); err != nil {
		log.Panic(err)
	}

	if err := txn.Verify(); err != nil {
		log.Panic(err)
	}

	log.Printf("signature= %s", txn.Sigs[0].Hex())
	return txn
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0750)
}

func init() {
	dbVerifyCheckpointVersionParsed = semver.MustParse(DBVerifyCheckpointVersion)
}
