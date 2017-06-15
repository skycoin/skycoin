package gui

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

var (
	logger   = util.MustGetLogger("gui")
	listener net.Listener
	quit     chan struct{}
)

const (
	resourceDir = "dist/"
	devDir      = "dev/"
	indexPage   = "index.html"
)

// LaunchWebInterface begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, daemon *daemon.Daemon) error {
	quit = make(chan struct{})
	logger.Info("Starting web interface on http://%s", host)
	logger.Warning("HTTPS not in use!")
	appLoc, err := util.DetermineResourcePath(staticDir, resourceDir, devDir)
	if err != nil {
		return err
	}
	logger.Info("Web resources directory: %s", appLoc)

	listener, err = net.Listen("tcp", host)
	if err != nil {
		return err
	}

	// Runs http.Serve() in a goroutine
	serve(listener, NewGUIMux(appLoc, daemon), quit)
	return nil
}

// LaunchWebInterfaceHTTPS begins listening on https://$host, for enabling remote web access
// Uses HTTPS
func LaunchWebInterfaceHTTPS(host, staticDir string, daemon *daemon.Daemon, certFile, keyFile string) error {
	quit = make(chan struct{})
	logger.Info("Starting web interface on https://%s", host)
	logger.Info("Using %s for the certificate", certFile)
	logger.Info("Using %s for the key", keyFile)
	logger.Info("Web resources directory: %s", staticDir)

	appLoc, err := util.DetermineResourcePath(staticDir, devDir, resourceDir)
	if err != nil {
		return err
	}

	certs := make([]tls.Certificate, 1)
	if certs[0], err = tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return err
	}

	listener, err = tls.Listen("tcp", host, &tls.Config{Certificates: certs})
	if err != nil {
		return err
	}

	// Runs http.Serve() in a goroutine
	serve(listener, NewGUIMux(appLoc, daemon), quit)
	return nil
}

func serve(listener net.Listener, mux *http.ServeMux, q chan struct{}) {
	go func() {
		for {
			if err := http.Serve(listener, mux); err != nil {
				select {
				case <-q:
					return
				default:
				}
				continue
			}
		}
	}()
}

// Shutdown close http service
func Shutdown() {
	if quit != nil {
		// must close quit first
		close(quit)
		listener.Close()
		listener = nil
	}
}

// NewGUIMux creates an http.ServeMux with handlers registered
func NewGUIMux(appLoc string, daemon *daemon.Daemon) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", newIndexHandler(appLoc))

	fileInfos, _ := ioutil.ReadDir(appLoc)
	for _, fileInfo := range fileInfos {
		route := fmt.Sprintf("/%s", fileInfo.Name())
		if fileInfo.IsDir() {
			route = route + "/"
		}
		mux.Handle(route, http.FileServer(http.Dir(appLoc)))
	}

	// Wallet interface
	RegisterWalletHandlers(mux, daemon.Gateway)
	// Blockchain interface
	RegisterBlockchainHandlers(mux, daemon.Gateway)
	// Network stats interface
	RegisterNetworkHandlers(mux, daemon.Gateway)
	// Network API handler
	RegisterAPIHandlers(mux, daemon.Gateway)
	// Transaction handler
	RegisterTxHandlers(mux, daemon.Gateway)
	// UxOUt api handler
	RegisterUxOutHandlers(mux, daemon.Gateway)
	// expplorer handler
	RegisterExplorerHandlers(mux, daemon.Gateway)
	return mux
}

// Returns a http.HandlerFunc for index.html, where index.html is in appLoc
func newIndexHandler(appLoc string) http.HandlerFunc {
	// Serves the main page
	return func(w http.ResponseWriter, r *http.Request) {
		page := filepath.Join(appLoc, indexPage)
		logger.Debug("Serving index page: %s", page)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, page)
		} else {
			wh.Error404(w)
		}
	}
}
