package gui

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/op/go-logging.v1"

	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

var (
	logger = logging.MustGetLogger("skywire.gui")
)

const (
	resourceDir = "dist/"
	devDir      = "dev/"
	indexPage   = "index.html"
)

// Begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, nodemanager *nodemanager.NodeManager) error {
	logger.Info("Starting web interface on http://%s", host)
	logger.Warning("HTTPS not in use!")
	logger.Info("Web resources directory: %s", staticDir)

	appLoc, err := util.DetermineResourcePath(staticDir, devDir, resourceDir)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	// Runs http.Serve() in a goroutine
	serve(listener, NewGUIMux(appLoc, daemon))
	return nil
}

func serve(listener net.Listener, mux *http.ServeMux) {
	// http.Serve() blocks
	// Minimize the chance of http.Serve() not being ready before the
	// function returns and the browser opens
	ready := make(chan struct{})
	go func() {
		ready <- struct{}{}
		if err := http.Serve(listener, mux); err != nil {
			log.Panic(err)
		}
	}()
	<-ready
}

//move to util

// Creates an http.ServeMux with handlers registered
func NewGUIMux(appLoc string, daemon *nodemanager.NodeManager) *http.ServeMux {
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
	RegisterApiHandlers(mux, daemon.Gateway)
	// Transaction interface
	RegisterTxHandlers(mux, daemon.Gateway)
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
			Error404(w)
		}
	}
}
