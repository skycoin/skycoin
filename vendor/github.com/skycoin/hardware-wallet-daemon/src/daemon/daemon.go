package daemon

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"
	"time"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
	"github.com/skycoin/skycoin/src/util/apputil"
	"github.com/skycoin/skycoin/src/util/logging"

	"github.com/skycoin/hardware-wallet-daemon/src/api"
)

// Daemon represents a hardware wallet daemon instance
type Daemon struct {
	config Config
	logger *logging.Logger
}

// NewDaemon returns a new hardware wallet daemon instance
func NewDaemon(config Config, logger *logging.Logger) *Daemon {
	return &Daemon{
		config: config,
		logger: logger,
	}
}

// Run starts the daemon
func (d *Daemon) Run() error {
	var apiServer *api.Server
	var retErr error
	errC := make(chan error, 10)

	logLevel, err := logging.LevelFromString(d.config.App.LogLevel)
	if err != nil {
		err = fmt.Errorf("invalid -log-level: %v", err)
		d.logger.Error(err)
		return err
	}

	logging.SetLevel(logLevel)

	if d.config.App.ColorLog {
		logging.EnableColors()
	} else {
		logging.DisableColors()
	}

	var logFile *os.File
	if d.config.App.LogToFile {
		var err error
		logFile, err = d.initLogFile()
		if err != nil {
			d.logger.Error(err)
			return err
		}
	}

	host := fmt.Sprintf("%s:%d", d.config.App.WebInterfaceAddr, d.config.App.WebInterfacePort)

	if d.config.App.ProfileCPU {
		f, err := os.Create(d.config.App.ProfileCPUFile)
		if err != nil {
			d.logger.Error(err)
			return err
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			d.logger.Error(err)
			return err
		}
		defer pprof.StopCPUProfile()
	}

	if d.config.App.HTTPProf {
		go func() {
			if err := http.ListenAndServe(d.config.App.HTTPProfHost, nil); err != nil {
				d.logger.WithError(err).Errorf("Listen on HTTP profiling interface %s failed", d.config.App.HTTPProfHost)
			}
		}()
	}

	var wg sync.WaitGroup

	quit := make(chan struct{})

	// Catch SIGINT (CTRL-C) (closes the quit channel)
	go apputil.CatchInterrupt(quit)

	// Catch SIGUSR1 (prints runtime stack to stdout)
	go apputil.CatchDebug()

	apiServer, err = d.createServer(host, api.NewGateway(skyWallet.NewDevice(d.config.App.daemonMode)))
	if err != nil {
		d.logger.Error(err)
		retErr = err
		goto earlyShutdown
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := apiServer.Serve(); err != nil {
			d.logger.Error(err)
			errC <- err
		}
	}()

	select {
	case <-quit:
	case retErr = <-errC:
		d.logger.Error(retErr)
	}

	d.logger.Info("Shutting down...")

	if apiServer != nil {
		d.logger.Info("Closing api server")
		apiServer.Shutdown()
	}

	d.logger.Info("Waiting for goroutines to finish")
	wg.Wait()

earlyShutdown:
	d.logger.Info("Goodbye")

	if logFile != nil {
		if err := logFile.Close(); err != nil {
			fmt.Println("Failed to close log file")
		}
	}

	return retErr
}

func (d *Daemon) initLogFile() (*os.File, error) {
	logDir := filepath.Join(d.config.App.DataDirectory, "logs")
	if err := createDirIfNotExist(logDir); err != nil {
		d.logger.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
		return nil, fmt.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
	}

	// open log file
	tf := "2006-01-02-030405"
	logfile := filepath.Join(logDir, fmt.Sprintf("%s.log", time.Now().Format(tf)))

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		d.logger.Errorf("os.OpenFile(%s) failed: %v", logfile, err)
		return nil, err
	}

	hook := logging.NewWriteHook(f)
	logging.AddHook(hook)

	return f, nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0750)
}

func (d *Daemon) createServer(host string, gateway *api.Gateway) (*api.Server, error) {
	apiConfig := api.Config{
		EnableCSRF:         d.config.App.EnableCSRF,
		DisableHeaderCheck: d.config.App.DisableHeaderCheck,
		HostWhitelist:      d.config.App.hostWhitelist,
		Mode:               d.config.App.daemonMode,
		Build:              d.config.Build,
	}

	var s *api.Server

	var err error
	s, err = api.Create(host, apiConfig, gateway)
	if err != nil {
		d.logger.Errorf("Failed to start web GUI: %v", err)
		return nil, err
	}

	return s, nil
}

// ParseConfig prepare the config
func (d *Daemon) ParseConfig() error {
	return d.config.postProcess()
}
