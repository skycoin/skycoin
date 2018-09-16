package skycoin

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/apputil"
	"github.com/skycoin/skycoin/src/util/browser"
	"github.com/skycoin/skycoin/src/util/cert"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/wallet"
)

// Coin represents a fiber coin instance
type Coin struct {
	config Config
	logger *logging.Logger
}

// Run starts the node
func (c *Coin) Run() error {
	var db *dbutil.DB
	var d *daemon.Daemon
	var webInterface *api.Server
	var retErr error
	errC := make(chan error, 10)

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

	// creates blockchain instance
	dconf := c.ConfigureDaemon()

	c.logger.Infof("Opening database %s", dconf.Visor.DBPath)
	db, err = visor.OpenDB(dconf.Visor.DBPath, c.config.Node.DBReadOnly)
	if err != nil {
		c.logger.Errorf("Database failed to open: %v. Is another skycoin instance running?", err)
		return err
	}

	if c.config.Node.ResetCorruptDB {
		// Check the database integrity and recreate it if necessary
		c.logger.Info("Checking database and resetting if corrupted")
		if newDB, err := visor.RepairCorruptDB(db, c.config.Node.blockchainPubkey, quit); err != nil {
			if err != visor.ErrVerifyStopped {
				c.logger.Errorf("visor.ResetCorruptDB failed: %v", err)
				retErr = err
			}
			goto earlyShutdown
		} else {
			db = newDB
		}
	} else if c.config.Node.VerifyDB {
		c.logger.Info("Checking database")
		if err := visor.CheckDatabase(db, c.config.Node.blockchainPubkey, quit); err != nil {
			if err != visor.ErrVerifyStopped {
				c.logger.Errorf("visor.CheckDatabase failed: %v", err)
				retErr = err
			}
			goto earlyShutdown
		}
	}

	d, err = daemon.NewDaemon(dconf, db)
	if err != nil {
		c.logger.Error(err)
		retErr = err
		goto earlyShutdown
	}

	if c.config.Node.WebInterface {
		webInterface, err = c.createGUI(d, host)
		if err != nil {
			c.logger.Error(err)
			retErr = err
			goto earlyShutdown
		}
	}

	fullAddress = fmt.Sprintf("%s://%s", scheme, webInterface.Addr())
	c.logger.Critical().Infof("Full address: %s", fullAddress)
	if c.config.Node.PrintWebInterfaceAddress {
		fmt.Println(fullAddress)
	}

	if err := d.Init(); err != nil {
		c.logger.Error(err)
		retErr = err
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

			if err := webInterface.Serve(); err != nil {
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
					if err := browser.Open(fullAddress); err != nil {
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

	if webInterface != nil {
		c.logger.Info("Closing web interface")
		webInterface.Shutdown()
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

// ConfigureDaemon sets the daemon config values
func (c *Coin) ConfigureDaemon() daemon.Config {
	dc := daemon.NewConfig()

	dc.Pool.DefaultConnections = c.config.Node.DefaultConnections
	dc.Pool.MaxDefaultPeerOutgoingConnections = c.config.Node.MaxDefaultPeerOutgoingConnections

	dc.Pex.DataDirectory = c.config.Node.DataDirectory
	dc.Pex.Disabled = c.config.Node.DisablePEX
	dc.Pex.Max = c.config.Node.PeerlistSize
	dc.Pex.DownloadPeerList = c.config.Node.DownloadPeerList
	dc.Pex.PeerListURL = c.config.Node.PeerListURL
	dc.Pex.DisableTrustedPeers = c.config.Node.DisableDefaultPeers
	dc.Pex.CustomPeersFile = c.config.Node.CustomPeersFile
	dc.Pex.DefaultConnections = c.config.Node.DefaultConnections

	dc.Daemon.DefaultConnections = c.config.Node.DefaultConnections
	dc.Daemon.DisableOutgoingConnections = c.config.Node.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = c.config.Node.DisableIncomingConnections
	dc.Daemon.DisableNetworking = c.config.Node.DisableNetworking
	dc.Daemon.Port = c.config.Node.Port
	dc.Daemon.Address = c.config.Node.Address
	dc.Daemon.LocalhostOnly = c.config.Node.LocalhostOnly
	dc.Daemon.OutgoingMax = c.config.Node.MaxOutgoingConnections
	dc.Daemon.DataDirectory = c.config.Node.DataDirectory
	dc.Daemon.LogPings = !c.config.Node.DisablePingPong
	dc.Daemon.BlockchainPubkey = c.config.Node.blockchainPubkey

	if c.config.Node.OutgoingConnectionsRate == 0 {
		c.config.Node.OutgoingConnectionsRate = time.Millisecond
	}
	dc.Daemon.OutgoingRate = c.config.Node.OutgoingConnectionsRate
	dc.Visor.IsMaster = c.config.Node.RunMaster

	dc.Visor.BlockchainPubkey = c.config.Node.blockchainPubkey
	dc.Visor.BlockchainSeckey = c.config.Node.blockchainSeckey

	dc.Visor.GenesisAddress = c.config.Node.genesisAddress
	dc.Visor.GenesisSignature = c.config.Node.genesisSignature
	dc.Visor.GenesisTimestamp = c.config.Node.GenesisTimestamp
	dc.Visor.GenesisCoinVolume = c.config.Node.GenesisCoinVolume
	dc.Visor.DBPath = c.config.Node.DBPath
	dc.Visor.Arbitrating = c.config.Node.Arbitrating
	dc.Visor.WalletDirectory = c.config.Node.WalletDirectory
	_, dc.Visor.EnableWalletAPI = c.config.Node.enabledAPISets[api.EndpointsWallet]
	_, dc.Visor.EnableSeedAPI = c.config.Node.enabledAPISets[api.EndpointsWalletSeed]

	_, dc.Gateway.EnableWalletAPI = c.config.Node.enabledAPISets[api.EndpointsWallet]
	_, dc.Gateway.EnableSpendMethod = c.config.Node.enabledAPISets[api.EndpointsDeprecatedWalletSpend]

	// Initialize wallet default crypto type
	cryptoType, err := wallet.CryptoTypeFromString(c.config.Node.WalletCryptoType)
	if err != nil {
		log.Panic(err)
	}

	dc.Visor.WalletCryptoType = cryptoType

	return dc
}

func (c *Coin) createGUI(d *daemon.Daemon, host string) (*api.Server, error) {
	var s *api.Server
	var err error

	config := api.Config{
		StaticDir:            c.config.Node.GUIDirectory,
		DisableCSRF:          c.config.Node.DisableCSRF,
		DisableCSP:           c.config.Node.DisableCSP,
		EnableJSON20RPC:      c.config.Node.RPCInterface,
		EnableGUI:            c.config.Node.EnableGUI,
		EnableUnversionedAPI: c.config.Node.EnableUnversionedAPI,
		ReadTimeout:          c.config.Node.ReadTimeout,
		WriteTimeout:         c.config.Node.WriteTimeout,
		IdleTimeout:          c.config.Node.IdleTimeout,
		EnabledAPISets:       c.config.Node.enabledAPISets,
		BuildInfo: readable.BuildInfo{
			Version: c.config.Build.Version,
			Commit:  c.config.Build.Commit,
			Branch:  c.config.Build.Branch,
		},
	}

	if c.config.Node.WebInterfaceHTTPS {
		// Verify cert/key parameters, and if neither exist, create them
		if err := cert.CreateCertIfNotExists(host, c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey, "Skycoind"); err != nil {
			c.logger.Errorf("cert.CreateCertIfNotExists failure: %v", err)
			return nil, err
		}

		s, err = api.CreateHTTPS(host, config, d.Gateway, c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey)
	} else {
		s, err = api.Create(host, config, d.Gateway)
	}
	if err != nil {
		c.logger.Errorf("Failed to start web GUI: %v", err)
		return nil, err
	}

	return s, nil
}

// ParseConfig prepare the config
func (c *Coin) ParseConfig() {
	c.config.postProcess()
}

// InitTransaction creates the initialize transaction
func InitTransaction(UxID string, genesisSecKey cipher.SecKey) coin.Transaction {
	var tx coin.Transaction

	output := cipher.MustSHA256FromHex(UxID)
	tx.PushInput(output)

	addrs := visor.GetDistributionAddresses()

	if len(addrs) != 100 {
		log.Panic("Should have 100 distribution addresses")
	}

	// 1 million per address, measured in droplets
	if visor.DistributionAddressInitialBalance != 1e6 {
		log.Panic("visor.DistributionAddressInitialBalance expected to be 1e6*1e6")
	}

	for i := range addrs {
		addr := cipher.MustDecodeBase58Address(addrs[i])
		tx.PushOutput(addr, visor.DistributionAddressInitialBalance*1e6, 1)
	}

	seckeys := make([]cipher.SecKey, 1)
	seckey := genesisSecKey.Hex()
	seckeys[0] = cipher.MustSecKeyFromHex(seckey)
	tx.SignInputs(seckeys)

	tx.UpdateHeader()

	err := tx.Verify()

	if err != nil {
		log.Panic(err)
	}

	log.Printf("signature= %s", tx.Sigs[0].Hex())
	return tx
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0750)
}
