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

	"github.com/skycoin/skycoin/src/daemon"
)

var (
	logger = logging.MustGetLogger("skycoin.gui")
)

const (
	resourceDir = "dist/"
	indexPage   = "index.html"
)

// Begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, daemon *daemon.Daemon) error {
	logger.Info("Starting web interface on http://%s", host)
	logger.Warning("HTTPS not in use!")

	//appLoc := filepath.Join(staticDir, resourceDir)
	appLoc, err := determineResourcePath(staticDir)
	if err != nil {
		return err
	}
	web_interface_active := make(chan bool, 1) //do not return until webserver is running
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	// Runs http.Serve() in a goroutine
	serve(listener, NewGUIMux(appLoc, daemon))
	return nil
}

// Begins listening on https://$host, for enabling remote web access
// Uses HTTPS
func LaunchWebInterfaceHTTPS(host, staticDir string, daemon *daemon.Daemon, certFile, keyFile string) error {
	logger.Info("Starting web interface on https://%s", host)
	logger.Info("Using %s for the certificate", certFile)
	logger.Info("Using %s for the key", keyFile)

	//appLoc := filepath.Join(staticDir, resourceDir)
	appLoc, err := determineResourcePath(staticDir)
	if err != nil {
		return err
	}

	mux := NewGUIMux(appLoc, daemon)

	certs := make([]tls.Certificate, 1)
	if certs[0], err = tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return err
	}

	listener, err := tls.Listen("tcp", host, &tls.Config{Certificates: certs})
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

func determineResourcePath(staticDir string) (string, error) {
	appLoc := filepath.Join(staticDir, resourceDir)
	if strings.HasPrefix(appLoc, "/") {
		return appLoc, nil
	}

	// Prepend the binary's directory path if appLoc is relative
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, appLoc), nil
}

// Creates an http.ServeMux with handlers registered
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
